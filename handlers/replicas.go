// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	marathon "github.com/gambol99/go-marathon"
	"github.com/gorilla/mux"
	"github.com/openfaas/faas-netes/types"
	"github.com/openfaas/faas/gateway/requests"
)

// MakeReplicaUpdater updates desired count of replicas
func MakeReplicaUpdater(client marathon.Marathon) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Update replicas")

		vars := mux.Vars(r)
		functionName := vars["name"]

		req := types.ScaleServiceRequest{}
		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				msg := "Cannot parse request. Please pass valid JSON."
				w.Write([]byte(msg))
				log.Println(msg, marshalErr)
				return
			}
		}

		applicationID := Function2ID(functionName)
		if _, err := client.ScaleApplicationInstances(applicationID, int(req.Replicas), true); err != nil {
			reportScaleError(functionName, err, w)
			return
		}

		if err := client.WaitOnApplication(applicationID, 30*time.Second); err != nil {
			reportScaleError(functionName, err, w)
		}
	}
}

func reportScaleError(functionName string, err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("Unable to update function deployment " + functionName))
	log.Println(err)
}

// MakeReplicaReader reads the amount of replicas for a deployment
func MakeReplicaReader(client marathon.Marathon) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Update replicas")

		vars := mux.Vars(r)
		functionName := vars["name"]

		functions, err := getServiceList(client)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		var found *requests.Function
		for _, function := range functions {
			if function.Name == functionName {
				found = &function
				break
			}
		}

		if found == nil {
			w.WriteHeader(404)
			return
		}

		functionBytes, _ := json.Marshal(found)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}
