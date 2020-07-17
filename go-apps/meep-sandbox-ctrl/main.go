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

	server "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-sandbox-ctrl/server"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	ss "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions"

	"github.com/gorilla/handlers"
)

func init() {
	log.MeepTextLogInit("meep-sandbox-ctrl")
}

func main() {
	log.Info(os.Args)
	log.Info("Starting MEEP Sandbox Controller")

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
		// Initialize Sandbox Controller
		err := server.Init()
		if err != nil {
			log.Error("Failed to initialize Sandbox Controller")
			run = false
			return
		}

		// Start Sandbox Controller
		err = server.Run()
		if err != nil {
			log.Error("Failed to start Sandbox Controller")
			run = false
			return
		}

		var priSw string
		var altSw string
		userSw := os.Getenv("USER_SWAGGER_SANDBOX")
		if userSw == "" {
			priSw = "./swagger/"
			altSw = ""
		} else {
			priSw = "." + userSw + "/"
			altSw = "./swagger/"
		}

		// Start alternate REST API Server if needed
		altServ := os.Getenv("ALT_SERVICE_PORT")
		if altServ != "" {
			go func() {
				log.Info("Starting Alt-server on port " + altServ)
				log.Info("Alt-serving [sw:" + altSw)
				secAccessMap := map[string]string{}
				secRouter := server.NewRouter(altSw, secAccessMap)
				methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
				header := handlers.AllowedHeaders([]string{"content-type"})
				log.Fatal(http.ListenAndServe(":"+altServ, handlers.CORS(methods, header)(secRouter)))
				run = false
			}()
		}

		// Start primary REST API Server
		log.Info("Starting Primary-server on port 80")
		log.Info("Primary-serving [sw:" + priSw)
		priAccessMap := map[string]string{}
		if altServ != "" {
			priAccessMap = map[string]string{
				"Index":                            ss.AccessBlock,
				"ActivateScenario":                 ss.AccessVerify,
				"GetActiveNodeServiceMaps":         ss.AccessBlock,
				"GetActiveScenario":                ss.AccessVerify,
				"TerminateScenario":                ss.AccessVerify,
				"CreateReplayFile":                 ss.AccessBlock,
				"CreateReplayFileFromScenarioExec": ss.AccessBlock,
				"DeleteReplayFile":                 ss.AccessBlock,
				"DeleteReplayFileList":             ss.AccessBlock,
				"GetReplayFile":                    ss.AccessBlock,
				"GetReplayFileList":                ss.AccessBlock,
				"GetReplayStatus":                  ss.AccessBlock,
				"LoopReplay":                       ss.AccessBlock,
				"PlayReplayFile":                   ss.AccessBlock,
				"StopReplayFile":                   ss.AccessBlock,
				"SendEvent":                        ss.AccessVerify,
			}
		}
		router := server.NewRouter(priSw, priAccessMap)
		methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
		header := handlers.AllowedHeaders([]string{"content-type"})
		log.Fatal(http.ListenAndServe(":80", handlers.CORS(methods, header)(router)))
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
