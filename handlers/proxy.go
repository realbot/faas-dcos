// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"io/ioutil"

	"github.com/gorilla/mux"
)

// MakeProxy creates a proxy for HTTP web requests which can be routed to a function.
func MakeProxy() http.HandlerFunc {
	proxyClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 0,
			}).DialContext,
			MaxIdleConns:          1,
			DisableKeepAlives:     true,
			IdleConnTimeout:       120 * time.Millisecond,
			ExpectContinueTimeout: 1500 * time.Millisecond,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		switch r.Method {
		case http.MethodGet,
			http.MethodPost:

			vars := mux.Vars(r)
			service := vars["name"]

			if len(service) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Provide a valid function route."))
				return
			}

			stamp := strconv.FormatInt(time.Now().Unix(), 10)

			defer func(when time.Time) {
				seconds := time.Since(when).Seconds()
				log.Printf("[%s] took %f seconds\n", stamp, seconds)
			}(time.Now())

			requestBody, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			// Using Minuteman (https://github.com/dcos/minuteman)
			url := fmt.Sprintf("http://%s.marathon.l4lb.thisdcos.directory:8080/", Function2Endpoint(service))

			request, _ := http.NewRequest("POST", url, bytes.NewReader(requestBody))

			copyHeaders(&request.Header, &r.Header)

			defer request.Body.Close()

			response, err := proxyClient.Do(request)
			if err != nil {
				log.Println(err.Error())
				writeHead(service, http.StatusInternalServerError, w)
				buf := bytes.NewBufferString("Can't reach service: " + service)
				w.Write(buf.Bytes())
				return
			}

			clientHeader := w.Header()
			copyHeaders(&clientHeader, &response.Header)

			writeHead(service, http.StatusOK, w)
			io.Copy(w, response.Body)
		}
	}
}

func writeHead(service string, code int, w http.ResponseWriter) {
	w.WriteHeader(code)
}

func copyHeaders(destination *http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(*destination)[k] = vClone
	}
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
