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

func GetMeepSettings(w http.ResponseWriter, r *http.Request) {
	ceGetMeepSettings(w, r)
}

func SetMeepSettings(w http.ResponseWriter, r *http.Request) {
	ceSetMeepSettings(w, r)
}
