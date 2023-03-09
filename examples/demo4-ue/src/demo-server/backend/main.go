/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
	"path/filepath"
	"syscall"
	"time"

	server "github.com/AdvantEDGE/examples/demo4-ue/src/demo-server/backend/server"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/handlers"
)

// Initalize customized logger
func init() {
	log.MeepTextLogInit("demo4-ue")
}

func main() {
	var (
		dir      string
		fileName string
		run      bool = true
		done     chan bool
	)

	if len(os.Args) < 2 {
		// no config argument
		log.Fatal("Missing parameter, require file path to configurations!")
	}

	// Read configuration file path in command line arugments
	configPath := os.Args[1]
	dir = filepath.Dir(configPath)
	fileName = filepath.Base(configPath)

	go func() {

		port, err := server.Init(dir, fileName)
		if err != nil {
			log.Fatal("Failed to initalize Demo 4 ", err)
		}
		log.Info("main: server.Init done")

		// Channel sync pass channel listen for app termination
		server.Run(done)

		// Start demo 4 server
		router := server.NewRouter()
		methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
		header := handlers.AllowedHeaders([]string{"content-type"})
		log.Info("main: Starting listener on port ", port)
		log.Fatal(http.ListenAndServe(port, handlers.CORS(methods, header)(router)))
		run = false
	}()

	// Listen for SIGKILL
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info("Waiting to shut down program !")
		run = false
	}()

	// Listen for demo 4 error exit program
	go func() {
		<-done
		log.Info("Waiting to shut down program !")
		run = false
	}()

	for {
		// Invoke graceful termination upon program kill
		if !run {
			log.Info("Invoking demo 4 graceful termination")
			server.Terminate()
			break
		}
		time.Sleep(time.Second)
	}
}
