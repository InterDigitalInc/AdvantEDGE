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

package replay

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	ce "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-client"
	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const defaultLoopInterval = 5000 //in ms
const basepath = "http://meep-ctrl-engine/v1"

type ReplayMgr struct {
	name             string
	currentFileName  string
	isStarted        bool
	ticker           *time.Ticker
	nextEventIndex   int
	eventIndexMax    int
	replayEventsList ceModel.Replay
	loop             bool
	client           *ce.APIClient
	timeToNextEvent  int
	timeRemaining    int
	timeStarted      time.Time
	loopStarted      time.Time
	ignoreInitEvent  bool
}

func createClient(path string) (*ce.APIClient, error) {
	// Create & store client for App REST API
	ceClientCfg := ce.NewConfiguration()
	ceClientCfg.BasePath = path
	ceClient := ce.NewAPIClient(ceClientCfg)
	if ceClient == nil {
		err := errors.New("Failed to create ctrl-engine REST API client")
		return nil, err
	}
	return ceClient, nil
}

func (r *ReplayMgr) IsStarted() bool {
	return r.isStarted
}

// NewReplayMgr - Create, Initialize and connect the replay manager
func NewReplayMgr(name string) (r *ReplayMgr, err error) {
	if name == "" {
		err = errors.New("Missing replay manager name")
		log.Error(err)
		return nil, err
	}

	r = new(ReplayMgr)
	r.name = name
	r.isStarted = false

	client, err := createClient(basepath)
	if err != nil {
		log.Error("Error creating client: ", err)
		return
	}
	r.client = client

	log.Debug("ReplayMgr created ", r.name)
	return r, nil
}

func (r *ReplayMgr) playEventByIndex() error {

	// Retrieve & Marshal next event
	index := r.nextEventIndex
	replayEvent := r.replayEventsList.Events[index]
	j, err := json.Marshal(&replayEvent.Event)
	if err != nil {
		log.Error(err)
		return err
	}

	// If first event, take timestamp for reference
	if index == 0 {
		r.loopStarted = time.Now()
	}

	// Process INIT event
	isInitEvent := false
	if replayEvent.Event.Type_ == "OTHER" && replayEvent.Event.Name == "Init" {

		// Skip to next event if INIT event should be ignored
		if r.ignoreInitEvent {
			index = index + 1
			replayEvent = r.replayEventsList.Events[index]

			j, err = json.Marshal(&replayEvent.Event)
			if err != nil {
				log.Error(err)
				return err
			}
		} else {
			isInitEvent = true
		}
	}

	// Send event (except INIT event)
	if !isInitEvent {
		vars := make(map[string]string)
		vars["type"] = replayEvent.Event.Type_
		var validEvent ce.Event

		err = json.Unmarshal(j, &validEvent)
		if err != nil {
			log.Error(err)
			return err
		}

		_, err = r.client.ScenarioExecutionApi.SendEvent(context.TODO(), replayEvent.Event.Type_, validEvent)
		if err != nil {
			log.Error(err)
		}
	}

	// Retrieve index of next event
	nextIndex := 0
	if index != r.eventIndexMax {
		nextIndex = index + 1
	} else if r.loop {
		nextIndex = 0
	} else {
		nextIndex = -1
	}

	// If necessary, create timer to wait before sending next event
	if nextIndex != -1 {
		nextReplayEvent := r.replayEventsList.Events[nextIndex]

		// Calculate time until next event
		var diff int32
		if nextIndex == 0 {
			diff = defaultLoopInterval
		} else {
			diff = nextReplayEvent.Time - replayEvent.Time
		}
		tickerExpiry := time.Duration(diff) * time.Millisecond
		log.Debug("next replay event (index ", nextReplayEvent.Index, ") in ", tickerExpiry)
		r.ticker = time.NewTicker(tickerExpiry)
		r.timeToNextEvent, r.timeRemaining = r.getTimesRemaining()
		r.nextEventIndex = nextIndex

		// Start timer
		go func() {
			for range r.ticker.C {
				r.ticker.Stop()
				_ = r.playEventByIndex()
			}
		}()
	} else {
		r.Completed()
	}

	return nil
}

// Start - starts replay execution
func (r *ReplayMgr) Start(fileName string, replay ceModel.Replay, loop bool, ignoreInitEvent bool) error {
	if !r.isStarted {
		r.timeStarted = time.Now()
		r.isStarted = true
		r.nextEventIndex = 0
		r.replayEventsList = replay
		r.eventIndexMax = len(replay.Events) - 1
		r.loop = loop
		r.currentFileName = fileName
		r.ignoreInitEvent = ignoreInitEvent
		//executing the events
		_ = r.playEventByIndex()
	} else {
		return errors.New("Replay already running, filename: " + r.currentFileName)
	}
	return nil
}

// ForceStop - forced stop on the current replay file
func (r *ReplayMgr) ForceStop() bool {
	if r.isStarted {
		r.ticker.Stop()
		r.Completed()
		return true
	}
	return false
}

// Stop - stops replay file
func (r *ReplayMgr) Stop(replayFileName string) bool {
	if r.isStarted && r.currentFileName == replayFileName {
		r.ticker.Stop()
		r.Completed()
		return true
	}
	return false
}

// Completed - successfully terminates replay file
func (r *ReplayMgr) Completed() {
	r.isStarted = false
	log.Debug("replay completed execution")
}

// GetStatus - Returns the Replay Execution status
func (r *ReplayMgr) GetStatus() (status ceModel.ReplayStatus, err error) {
	if !r.IsStarted() {
		err = errors.New("No replay file running")
		return
	}

	status.ReplayFileRunning = r.currentFileName
	status.MaxIndex = int32(r.eventIndexMax)
	nextIndex := r.nextEventIndex
	maxIndex := r.eventIndexMax
	lastIndexPlayed := 0

	// if next index is 0, it means it will loop so we do not remove one to find the current, the current is the last event
	if nextIndex != 0 {
		if nextIndex == -1 {
			lastIndexPlayed = maxIndex
		} else {
			lastIndexPlayed = nextIndex - 1
		}
	} else {
		if r.loop {
			lastIndexPlayed = maxIndex
		} else {
			lastIndexPlayed = nextIndex
		}
	}

	status.Index = int32(lastIndexPlayed)
	status.MaxIndex = int32(maxIndex)

	status.LoopMode = r.loop
	timeToNextEvent, timeRemaining := r.getTimesRemaining()
	status.TimeToNextEvent = int32(timeToNextEvent)
	status.TimeRemaining = int32(timeRemaining)

	return status, nil
}

// getTimesRemaining - returns time left to execute next event and the rest of the replay file
func (r *ReplayMgr) getTimesRemaining() (int, int) {
	var elapsedDuration time.Duration
	if r.loop {
		elapsedDuration = time.Since(r.loopStarted)
	} else {
		elapsedDuration = time.Since(r.timeStarted)
	}
	elapsed := int(elapsedDuration / time.Millisecond)

	if r.ignoreInitEvent {
		elapsed += int(r.replayEventsList.Events[1].Time)
	}

	nextEventTimeRemaining := int(r.replayEventsList.Events[r.nextEventIndex].Time) - elapsed
	totalTimeRemaining := int(r.replayEventsList.Events[r.eventIndexMax].Time) - elapsed
	if nextEventTimeRemaining < 0 {
		nextEventTimeRemaining = 0
	}
	if totalTimeRemaining < 0 {
		totalTimeRemaining = 0
	}
	return nextEventTimeRemaining, totalTimeRemaining
}
