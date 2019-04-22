/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// Mobility Group
type MobilityGroup struct {

	// Mobility Group name
	Name string `json:"name,omitempty"`

	// State Transfer mode
	StateTransferMode string `json:"stateTransferMode,omitempty"`

	// State Transfer trigger
	StateTransferTrigger string `json:"stateTransferTrigger,omitempty"`

	// Session Transfer mode
	SessionTransferMode string `json:"sessionTransferMode,omitempty"`

	// Load Balancing Algorithm
	LoadBalancingAlgorithm string `json:"loadBalancingAlgorithm,omitempty"`
}
