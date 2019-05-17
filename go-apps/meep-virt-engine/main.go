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
	"os"
	"os/signal"
	"syscall"
	"time"

	server "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine/server"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func init() {
	log.MeepJSONLogInit("meep-virt-engine")
}

func main() {
	log.Info(os.Args)

	log.Info("Server started")

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
		err := server.VirtEngineInit()
		if err != nil {
			log.Error("Failed to initialize Virt. Engine")
			run = false
			return
		}

		// Start TC Engine Event Handler thread
		server.ListenEvents()
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
