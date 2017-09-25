// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/alexellis/faas/gateway/requests"
	marathon "github.com/gambol99/go-marathon"
)

func getServiceList(client marathon.Marathon) ([]requests.Function, error) {
	var functions []requests.Function

	v := url.Values{}
	v.Set("label", "faas_function")
	applications, err := client.Applications(v)

	if err != nil {
		return nil, err
	}
	for _, item := range applications.Apps {
		function := requests.Function{
			Name:            ID2Function(item.ID),
			Replicas:        uint64(item.TasksRunning),
			Image:           item.Container.Docker.Image,
			InvocationCount: 0,
		}
		functions = append(functions, function)
	}
	return functions, nil
}

// MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
func MakeFunctionReader(client marathon.Marathon) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		functions, err := getServiceList(client)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}
