/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-mg-manager/log"
	server "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-mg-manager/server"

	"github.com/gorilla/handlers"
)

func init() {
	log.MeepJSONLogInit()
}

func main() {
	log.Info(os.Args)

	log.Info("Starting MG Manager")

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
		// Initialize MG Manager
		err := server.Init()
		if err != nil {
			log.Error("Failed to initialize MG Manager")
			run = false
			return
		}

		// Start MG Manager Event Handler thread
		go server.Run()

		// Start MG Manager REST API Server
		router := server.NewRouter()
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
