/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package server

// Physical location object
type PhysicalLocation struct {

	// Unique physical location ID
	Id string `json:"id,omitempty"`

	// Physical location name
	Name string `json:"name,omitempty"`

	// Physical location type
	Type_ string `json:"type,omitempty"`

	// true: Physical location is external to MEEP false: Physical location is internal to MEEP
	IsExternal bool `json:"isExternal,omitempty"`

	NetworkLocationsInRange []string `json:"networkLocationsInRange,omitempty"`

	Processes []Process `json:"processes,omitempty"`
}
