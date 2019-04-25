/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// UE Mobility Event object
type EventUeMobility struct {

	// UE identifier
	Ue string `json:"ue,omitempty"`

	// Destination identifier
	Dest string `json:"dest,omitempty"`
}
