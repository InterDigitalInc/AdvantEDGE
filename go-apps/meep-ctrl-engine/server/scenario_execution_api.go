/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package server

import (
	"net/http"
)

func ActivateScenario(w http.ResponseWriter, r *http.Request) {
	ceActivateScenario(w, r)
}

func GetActiveScenario(w http.ResponseWriter, r *http.Request) {
	ceGetActiveScenario(w, r)
}

func GetActiveClientServiceMaps(w http.ResponseWriter, r *http.Request) {
	ceGetActiveClientServiceMaps(w, r)
}

func GetEventList(w http.ResponseWriter, r *http.Request) {
	ceGetEventList(w, r)
}

func SendEvent(w http.ResponseWriter, r *http.Request) {
	ceSendEvent(w, r)
}

func TerminateScenario(w http.ResponseWriter, r *http.Request) {
	ceTerminateScenario(w, r)
}
