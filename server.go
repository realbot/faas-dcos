// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"log"

	"github.com/realbot/faas-dcos/handlers"
	"github.com/realbot/faas-provider"

	// until health check is merged...
	bootTypes "github.com/realbot/faas-provider/types"

	//bootTypes "github.com/openfaas/faas-provider/types"

	marathon "github.com/gambol99/go-marathon"
)

func main() {
	config := marathon.NewDefaultConfig()
	config.URL = "http://master.mesos/service/marathon"
	client, err := marathon.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create a client for marathon, error: %s", err)
	}

	readConfig := ReadConfig{}
	osEnv := OsEnv{}
	cfg := readConfig.Read(osEnv)

	log.Printf("HTTP Read Timeout: %s\n", cfg.ReadTimeout)
	log.Printf("HTTP Write Timeout: %s\n", cfg.WriteTimeout)
	log.Printf("Function Readiness Probe Enabled: %v\n", cfg.EnableFunctionReadinessProbe)

	bootstrapHandlers := bootTypes.FaaSHandlers{
		FunctionProxy:  handlers.MakeProxy(),
		DeleteHandler:  handlers.MakeDeleteHandler(client),
		DeployHandler:  handlers.MakeDeployHandler(client),
		FunctionReader: handlers.MakeFunctionReader(client),
		ReplicaReader:  handlers.MakeReplicaReader(client),
		ReplicaUpdater: handlers.MakeReplicaUpdater(client),
		UpdateHandler:  handlers.MakeUpdateHandler(client),
		HealthHandler:  handlers.MakeHealthHandler(client),
	}

	var port int
	port = 8080
	bootstrapConfig := bootTypes.FaaSConfig{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		TCPPort:      &port,
	}

	log.Println("Starting faas-dcos")
	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)
}
