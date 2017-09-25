// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/alexellis/faas/gateway/requests"
	marathon "github.com/gambol99/go-marathon"
)

// MakeDeleteHandler delete a function
func MakeDeleteHandler(client marathon.Marathon) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.DeleteFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(request.FunctionName) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		// This makes sure we don't delete non-labelled deployments
		v := url.Values{}
		v.Set("label", "faas_function")
		v.Set("id", Function2ID(request.FunctionName))
		applications, err := client.Applications(v)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if len(applications.Apps) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		deleteFunction(client, request, w)
	}
}

func deleteFunction(client marathon.Marathon, request requests.DeleteFunctionRequest, w http.ResponseWriter) {

	if _, err := client.DeleteApplication(Function2ID(request.FunctionName), true); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
