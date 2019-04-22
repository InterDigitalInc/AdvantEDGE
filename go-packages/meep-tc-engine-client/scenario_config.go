/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package client

// Scenario configuration
type ScenarioConfig struct {

	// Visualization configuration
	Visualization string `json:"visualization,omitempty"`

	// Other scenario configuration
	Other string `json:"other,omitempty"`
}
