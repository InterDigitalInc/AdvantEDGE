/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
