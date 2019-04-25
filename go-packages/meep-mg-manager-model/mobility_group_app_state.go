/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// Mobility Group Application State
type MobilityGroupAppState struct {

	// Mobility Group UE Identifier
	UeId string `json:"ueId,omitempty"`

	// Mobility Group Application State for provided UE
	UeState string `json:"ueState,omitempty"`
}
