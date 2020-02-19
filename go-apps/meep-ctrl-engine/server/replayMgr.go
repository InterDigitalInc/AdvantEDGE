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

package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	ceModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/mux"
)

const defaultLoopInterval = 1000 //in ms

type ReplayMgr struct {
	name              string
	currentFileName   string
	isStarted         bool
	ticker            *time.Ticker
	currentEventIndex int
	eventIndexMax     int
	replayEventsList  ceModel.Replay
	loop              bool
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

	log.Debug("ReplayMgr created ", r.name)
	return r, nil
}

func (r *ReplayMgr) PlayEventByIndex() error {

	index := r.currentEventIndex
	nextIndex := 0
	replayEvent := r.replayEventsList.Events[index]

	j, err := json.Marshal(&replayEvent.Event)
	if err != nil {
		log.Error(err)
		return err
	}

	vars := make(map[string]string)

	vars["type"] = replayEvent.Event.Type_

	err = r.sendRequest(http.MethodPost, "/events", bytes.NewBuffer(j), vars, nil, ceSendEvent)
	if err != nil {
		log.Error(err)
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
				_ = r.PlayEventByIndex()
			}
		}()
	} else {
		r.Completed()
	}

	return nil
}

// Start - starts replay execution
func (r *ReplayMgr) Start(fileName string, replay ceModel.Replay, loop bool) error {
	if !r.isStarted {
		r.isStarted = true
		r.currentEventIndex = 0
		r.replayEventsList = replay
		r.eventIndexMax = len(replay.Events) - 1
		r.loop = loop
		r.currentFileName = fileName
		//executing the events
		_ = r.PlayEventByIndex()
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

func (r *ReplayMgr) sendRequest(method string, url string, body io.Reader, vars map[string]string, query map[string]string, f http.HandlerFunc) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil || req == nil {
		return err
	}
	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}
	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(f)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	_ = rr.Result()
	// Check the status code if needed

	return nil
}
