/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// Event object
type MobilityGroupEvent struct {

	// Mobility Group event name
	Name string `json:"name,omitempty"`

	// Mobility Group event type
	Type_ string `json:"type,omitempty"`

	// Mobility Group UE identifier
	UeId string `json:"ueId,omitempty"`

	AppState *MobilityGroupAppState `json:"appState,omitempty"`
}
