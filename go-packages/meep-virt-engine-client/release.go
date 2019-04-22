/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

type Release struct {

	// Release name
	Name string `json:"name,omitempty"`

	// Current release state
	State string `json:"state,omitempty"`
}
