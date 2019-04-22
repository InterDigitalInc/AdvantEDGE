/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

// Client-specific list of mappings of exposed port to internal service
type ClientServiceMap struct {

	// Unique external client identifier
	Client string `json:"client,omitempty"`

	ServiceMap []ServiceMap `json:"serviceMap,omitempty"`
}
