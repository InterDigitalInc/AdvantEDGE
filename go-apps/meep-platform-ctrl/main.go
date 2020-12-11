/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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

	server "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-platform-ctrl/server"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"

	"github.com/gorilla/handlers"
)

func init() {
	log.MeepTextLogInit("meep-platform-ctrl")
}

func main() {
	log.Info(os.Args)
	log.Info("Starting MEEP Platform Controller")

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
		// Initialize Platform Controller
		err := server.Init()
		if err != nil {
			log.Error("Failed to initialize Platform Controller")
			run = false
			return
		}

		// Start Platform Controller
		err = server.Run()
		if err != nil {
			log.Error("Failed to start Platform Controller")
			run = false
			return
		}

		var priFe string
		var altFe string
		userFe := os.Getenv("USER_FRONTEND")
		if userFe == "" {
			priFe = "./frontend/"
			altFe = ""
		} else {
			priFe = "." + userFe + "/"
			altFe = "./frontend/"
		}

		var priSw string
		var altSw string
		userSw := os.Getenv("USER_SWAGGER")
		if userSw == "" {
			priSw = "./swagger/"
			altSw = ""
		} else {
			priSw = "." + userSw + "/"
			altSw = "./swagger/"
		}

		// Start primary REST API Server
		log.Info("Primary-serving [fe:" + priFe + ", sw:" + priSw)
		log.Info("Alt-serving [fe:" + altFe + ", sw:" + altSw)
		priRouter := server.NewRouter(priFe, priSw, altFe, altSw)
		methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
		header := handlers.AllowedHeaders([]string{"content-type"})
		log.Fatal(http.ListenAndServe(":80", handlers.CORS(methods, header)(priRouter)))
		run = false
	}()

	count := 0
	for {
		if !run {
			log.Info("Ran for ", count, " seconds")
			break
		}
		time.Sleep(time.Second)
		count++
	}
}
