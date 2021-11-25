/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

package gc

import (
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"
)

const redisTable = 0

type GarbageCollector struct {
	rc      *redis.Connector
	ticker  *time.Ticker
	mutex   sync.mutex
	baseKey string
}

// NewGarbageCollector - Creates and initialize a Garbage Collector instance
func NewGarbageCollector(redisAddr string) (gc *GarbageCollector, err error) {
	// Create new Garbage Collector instance
	gc = new(GarbageCollector)

	// Connect to Redis DB
	gc.rc, err = redis.NewConnector(redisAddr, redisTable)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return nil, err
	}
	log.Info("Connected to Redis DB")

	log.Info("Created Garbage Collector")
	return gc, nil
}

// Start - start garbage collection ticker with provided interval
func (gc *GarbageCollector) Start(interval time.Duration) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Debug("[Garbage Collection] Periodic ticker started with interval: ", interval.String())
	gc.ticker = time.NewTicker(interval)
	go func() {
		for range gc.ticker.C {
			gc.Execute()
		}
	}()
}

// Stop - stop garbage collection ticker if running
func (gc *GarbageCollector) Stop() error {
	mutex.Lock()
	defer mutex.Unlock()

	if gc.ticker != nil {
		gc.ticker.Stop()
		gc.ticker = nil
		log.Debug("[Garbage Collection] Periodic ticker stopped")
	}
}

// Execute
func (gc *GarbageCollector) Execute() error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Debug("[Garbage Collection] Execution starting: ", time.Now())

	log.Debug("[Garbage Collection] Execution complete: ", time.Now())
}
