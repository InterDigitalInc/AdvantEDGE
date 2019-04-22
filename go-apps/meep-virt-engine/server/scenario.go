/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package server

// Scenario object
type Scenario struct {

	// Unique scenario name
	Name string `json:"name,omitempty"`

	Config *ScenarioConfig `json:"config,omitempty"`

	Deployment *Deployment `json:"deployment,omitempty"`
}
