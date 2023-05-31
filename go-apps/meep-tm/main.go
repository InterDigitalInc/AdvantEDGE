/*
 * Copyright (c) 2022  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/handlers"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.MeepJSONLogInit("meep-tm")
}

func main() {
	log.Info(os.Args)

	log.Info("Starting Traffic Management Service")

	run := true
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info("Program killed !")
		// do last actions and wait for all write operations to end
		run = false
	}()

	go func() {
		// Initialize Traffic Management Service
		err := server.Init()
		if err != nil {
			log.Error("Failed to initialize Traffic Management Service")
			run = false
			return
		}

		// Start Traffic Management Service Service Event Handler thread
		err = server.Run()
		if err != nil {
			log.Error("Failed to start Traffic Management Service")
			run = false
			return
		}

		// Start Traffic Management Service REST API Server
		router := server.NewRouter()
		methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT", "PATCH"})
		header := handlers.AllowedHeaders([]string{"Content-Type"})
		log.Fatal(http.ListenAndServe(":80", handlers.CORS(methods, header)(router)))
		run = false
	}()

	go func() {
		// Initialize Metrics Endpoint
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9000", nil))
		run = false
	}()

	count := 0
	for {
		if !run {
			_ = server.Stop()
			log.Info("Ran for ", count, " seconds")
			break
		}
		time.Sleep(time.Second)
		count++
	}

}
