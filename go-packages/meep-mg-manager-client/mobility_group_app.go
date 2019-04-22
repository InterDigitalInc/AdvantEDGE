/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

// Mobility Group Application instance
type MobilityGroupApp struct {

	// Mobility Group Application Identifier
	Id string `json:"id,omitempty"`

	// Event handler url
	Url string `json:"url,omitempty"`
}
