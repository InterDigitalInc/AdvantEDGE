/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package bws

const MAX_THROUGHPUT = 9999999999
const THROUGHPUT_UNIT = 1000000 //convert from Mbps to bps
const DEFAULT_THROUGHPUT_LINK = 1000.0
const moduleMetrics string = "metrics"
const moduleCtrlEngine string = "ctrl-engine"
const bwSharingControls string = "bws-controls"
const defaultTickerPeriod int = 500
const typeActive string = "active"

const channelCtrlActive string = moduleCtrlEngine + "-" + typeActive
const channelBwSharingControls string = bwSharingControls


