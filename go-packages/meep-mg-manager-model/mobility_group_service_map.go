/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// Mobility group service mapping
type MobilityGroupServiceMap struct {

	// Mobility group service name
	MgSvcName string `json:"mgSvcName,omitempty"`

	// Load balanced service instance name
	LbSvcName string `json:"lbSvcName,omitempty"`
}
