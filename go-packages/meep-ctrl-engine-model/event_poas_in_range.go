/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package model

// POAs In Range Event object
type EventPoasInRange struct {

	// UE identifier
	Ue string `json:"ue,omitempty"`

	PoasInRange []string `json:"poasInRange,omitempty"`
}
