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
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	server "github.com/InterDigitalInc/AdvantEDGE/example/demo4/src/onboarded-demo/server"
)

const (
	port = ":31124" // hardcoded in mepp-dai Chart template
)

func main() {
	err := server.Init()
	if err != nil {
		log.Fatal("Failed to initalize onboarded-demo4 ", err)
	}

	// Channel sync pass channel listen for app termination
	server.Run()

	FrontendApiService := server.NewFrontendApiService()
	FrontendApiController := server.NewFrontendApiController(FrontendApiService)

	router := server.NewRouter(FrontendApiController)

	log.Printf("Server started on port " + port)
	srv := &http.Server{Addr: port, Handler: router}
	go func() {
		//returns ErrServerClosed on graceful close
		log.Printf("Call ListenAndServe")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("ListenAndServe(): %s", err)
			return
		}
		log.Printf("Terminate goroutine")
	}()

	//	log.Fatal(http.ListenAndServe(port, router))
	log.Printf("Configure signals")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	log.Printf("Got %s signal. Aborting...\n", sig)
	server.Terminate()
	err = srv.Shutdown(context.TODO())
	if err != nil {
		log.Printf("Shutdown(): %s", err)
		return
	}
}
