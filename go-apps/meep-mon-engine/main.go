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

	log "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-mon-engine/log"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.MeepJSONLogInit()
}

func main() {
	log.Info(os.Args)

	log.Info("Starting Monitoring Engine")

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
		// Initialize Mon Engine
		err := Init()
		if err != nil {
			log.Error("Failed to initialize Mon Engine")
			run = false
			return
		}

		// Run Mon Engine
		err = Run()
		if err != nil {
			log.Error("Failed to run Mon Engine")
			run = false
			return
		}
		run = false
	}()

	count := 0
	for {
		if run == false {
			log.Info("Ran for", count, "seconds")
			break
		}
		time.Sleep(time.Second)
		count++
	}
}
