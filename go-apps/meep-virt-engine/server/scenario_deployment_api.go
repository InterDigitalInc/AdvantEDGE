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
	veActivateScenario(w, r)
}

func GetActiveScenario(w http.ResponseWriter, r *http.Request) {
	veGetActiveScenario(w, r)
}

func TerminateScenario(w http.ResponseWriter, r *http.Request) {
	veTerminateScenario(w, r)
}
