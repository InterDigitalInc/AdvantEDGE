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

func CreateScenario(w http.ResponseWriter, r *http.Request) {
	ceCreateScenario(w, r)
}

func DeleteScenario(w http.ResponseWriter, r *http.Request) {
	ceDeleteScenario(w, r)
}

func DeleteScenarioList(w http.ResponseWriter, r *http.Request) {
	ceDeleteScenarioList(w, r)
}

func GetScenario(w http.ResponseWriter, r *http.Request) {
	ceGetScenario(w, r)
}

func GetScenarioList(w http.ResponseWriter, r *http.Request) {
	ceGetScenarioList(w, r)
}

func SetScenario(w http.ResponseWriter, r *http.Request) {
	ceSetScenario(w, r)
}
