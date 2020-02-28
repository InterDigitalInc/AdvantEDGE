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
	"encoding/json"
	"errors"
	"time"
	"context"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	ce "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

const defaultLoopInterval = 1000 //in ms
const basepath = "http://meep-ctrl-engine/v1"

type ReplayMgr struct {
	name              string
	currentFileName   string
	isStarted         bool
	ticker            *time.Ticker
	currentEventIndex int
	eventIndexMax     int
	replayEventsList  ceModel.Replay
	loop              bool
        client            *ce.APIClient
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

func (r *ReplayMgr) PlayEventByIndex(ignoreOtherEvents bool) error {

	index := r.currentEventIndex
	nextIndex := 0
	replayEvent := r.replayEventsList.Events[index]

	j, err := json.Marshal(&replayEvent.Event)
	if err != nil {
		log.Error(err)
		return err
	}

	if replayEvent.Event.Type_ == "OTHER" && ignoreOtherEvents {
		index = index + 1
		replayEvent = r.replayEventsList.Events[index]

		j, err = json.Marshal(&replayEvent.Event)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	//only send events that mean something for the scenario
	if replayEvent.Event.Type_ != "OTHER" {

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

	//see if we have a next event, if we are done or if we loop
	if index != r.eventIndexMax {
		nextIndex = index + 1
	} else {
		if r.loop {
			nextIndex = 0
		} else {
			nextIndex = -1
		}
	}

	//sending the ticker prior to processing the current event to minimize time lost
	if nextIndex != -1 {
		nextReplayEvent := r.replayEventsList.Events[nextIndex]
		//durations are all relative to event0.... need to be updated based on current time
		//act otherwise if execution is a circular loop
		var diff int32
		if nextIndex == 0 {
			diff = defaultLoopInterval
		} else {
			diff = nextReplayEvent.Time - replayEvent.Time
		}
		r.currentEventIndex = nextIndex
		tickerExpiry := time.Duration(diff) * time.Millisecond
		log.Debug("next replay event (index ", nextReplayEvent.Index, ") in ", tickerExpiry)
		r.ticker = time.NewTicker(tickerExpiry)
		go func() {
			for range r.ticker.C {
				r.ticker.Stop()
				_ = r.PlayEventByIndex(ignoreOtherEvents)
			}
		}()
	}

        //only send events that mean something for the scenario
        if replayEvent.Event.Type_ != "OTHER" {

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

	if index == -1 {
		r.Completed()
	}

	return nil
}

// Start - starts replay execution
func (r *ReplayMgr) Start(fileName string, replay ceModel.Replay, loop bool, ignoreOtherEvents bool) error {
	if !r.isStarted {
		r.isStarted = true
		r.currentEventIndex = 0
		r.replayEventsList = replay
		r.eventIndexMax = len(replay.Events) - 1
		r.loop = loop
		r.currentFileName = fileName
		//executing the events
		_ = r.PlayEventByIndex(ignoreOtherEvents)
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
