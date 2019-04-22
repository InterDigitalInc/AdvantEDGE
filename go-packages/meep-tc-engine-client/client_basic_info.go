/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

// Client basic information object
type ClientBasicInfo struct {

	// Unique pod identifier
	PodId string `json:"podId,omitempty"`

	// IP address of the pod (client)
	Ip string `json:"ip,omitempty"`
}
