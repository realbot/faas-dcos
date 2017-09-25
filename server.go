// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"log"
	"time"

	"github.com/realbot/faas-dcos/handlers"

	// until health check is merged...
	myboot "github.com/realbot/faas-provider"

	bootTypes "github.com/alexellis/faas-provider/types"

	marathon "github.com/gambol99/go-marathon"
)

func main() {
	config := marathon.NewDefaultConfig()
	config.URL = "http://master.mesos/service/marathon"
	client, err := marathon.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create a client for marathon, error: %s", err)
	}

	bootstrapHandlers := bootTypes.FaaSHandlers{
		FunctionProxy:  handlers.MakeProxy(),
		DeleteHandler:  handlers.MakeDeleteHandler(client),
		DeployHandler:  handlers.MakeDeployHandler(client),
		FunctionReader: handlers.MakeFunctionReader(client),
		ReplicaReader:  handlers.MakeReplicaReader(client),
		ReplicaUpdater: handlers.MakeReplicaUpdater(client),
	}
	var port int
	port = 8080
	bootstrapConfig := bootTypes.FaaSConfig{
		ReadTimeout:  time.Second * 8,
		WriteTimeout: time.Second * 8,
		TCPPort:      &port,
	}

	log.Println("Starting faas-dcos")
	myboot.Serve(&bootstrapHandlers, &bootstrapConfig)
}
