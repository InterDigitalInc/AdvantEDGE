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
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.MeepJSONLogInit("meep-tc-engine")
}

func main() {
	log.Info(os.Args)
	log.Info("Starting TC Engine")

	// Start signal handler
	run := true
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Info("Program killed !")
		// do last actions and wait for all write operations to end
		run = false
	}()

	// Initialize & Start TC Engine
	go func() {
		err := Init()
		if err != nil {
			log.Error("Failed to initialize TC Engine")
			run = false
			return
		}

		err = Run()
		if err != nil {
			log.Error("Failed to start TC Engine")
			run = false
			return
		}
	}()

	// Main loop
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
