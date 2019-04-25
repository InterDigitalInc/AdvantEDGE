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

func GetStates(w http.ResponseWriter, r *http.Request) {
	ceGetStates(w, r)
}
