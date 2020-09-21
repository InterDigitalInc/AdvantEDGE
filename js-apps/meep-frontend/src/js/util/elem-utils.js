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

import {
  // Network Characteristics default values
  DEFAULT_LATENCY_INTER_DOMAIN,
  DEFAULT_LATENCY_JITTER_INTER_DOMAIN,
  DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN,
  DEFAULT_THROUGHPUT_DL_INTER_DOMAIN,
  DEFAULT_THROUGHPUT_UL_INTER_DOMAIN,
  DEFAULT_PACKET_LOSS_INTER_DOMAIN,
  DEFAULT_LATENCY_INTER_ZONE,
  DEFAULT_LATENCY_JITTER_INTER_ZONE,
  DEFAULT_THROUGHPUT_DL_INTER_ZONE,
  DEFAULT_THROUGHPUT_UL_INTER_ZONE,
  DEFAULT_PACKET_LOSS_INTER_ZONE,
  DEFAULT_LATENCY_INTRA_ZONE,
  DEFAULT_LATENCY_JITTER_INTRA_ZONE,
  DEFAULT_THROUGHPUT_DL_INTRA_ZONE,
  DEFAULT_THROUGHPUT_UL_INTRA_ZONE,
  DEFAULT_PACKET_LOSS_INTRA_ZONE,
  DEFAULT_LATENCY_TERMINAL_LINK,
  DEFAULT_LATENCY_JITTER_TERMINAL_LINK,
  DEFAULT_THROUGHPUT_DL_TERMINAL_LINK,
  DEFAULT_THROUGHPUT_UL_TERMINAL_LINK,
  DEFAULT_PACKET_LOSS_TERMINAL_LINK,
  DEFAULT_LATENCY_LINK,
  DEFAULT_LATENCY_JITTER_LINK,
  DEFAULT_THROUGHPUT_DL_LINK,
  DEFAULT_THROUGHPUT_UL_LINK,
  DEFAULT_PACKET_LOSS_LINK,
  DEFAULT_LATENCY_APP,
  DEFAULT_LATENCY_JITTER_APP,
  DEFAULT_THROUGHPUT_DL_APP,
  DEFAULT_THROUGHPUT_UL_APP,
  DEFAULT_PACKET_LOSS_APP
} from '../meep-constants';

// Network Element Fields
export const FIELD_TYPE = 'elementType';
export const FIELD_PARENT = 'parent';
export const FIELD_NAME = 'name';
export const FIELD_IMAGE = 'image';
export const FIELD_PORT = 'port';
export const FIELD_PROTOCOL = 'protocol';
export const FIELD_GROUP = 'group';
export const FIELD_INGRESS_SVC_MAP = 'ingressServiceMap';
export const FIELD_EGRESS_SVC_MAP = 'egressServiceMap';
export const FIELD_GPU_COUNT = 'gpuCount';
export const FIELD_GPU_TYPE = 'gpuType';
export const FIELD_PLACEMENT_ID = 'placementId';
export const FIELD_ENV_VAR = 'envVar';
export const FIELD_CMD = 'cmd';
export const FIELD_CMD_ARGS = 'cmdArgs';
export const FIELD_EXT_PORT = 'externalPort';
export const FIELD_IS_EXTERNAL = 'isExternal';
export const FIELD_MCC = 'mcc';
export const FIELD_MNC = 'mnc';
export const FIELD_DEFAULT_CELL_ID = 'defaultCellId';
export const FIELD_CELL_ID = 'cellId';
export const FIELD_NR_CELL_ID = 'nrCellId';
export const FIELD_MAC_ID = 'macId';
export const FIELD_GEO_LOCATION = 'location';
export const FIELD_GEO_RADIUS = 'radius';
export const FIELD_GEO_PATH = 'path';
export const FIELD_GEO_EOP_MODE = 'eopMode';
export const FIELD_GEO_VELOCITY = 'velocity';
export const FIELD_CHART_ENABLED = 'userChartEnabled';
export const FIELD_CHART_LOC = 'userChartLocation';
export const FIELD_CHART_VAL = 'userChartAlternateValues';
export const FIELD_CHART_GROUP = 'userChartGroup';
export const FIELD_CONNECTED = 'connected';
export const FIELD_WIRELESS = 'wireless';
export const FIELD_WIRELESS_TYPE = 'wirelessType';
export const FIELD_INT_DOM_LATENCY = 'interDomainLatency';
export const FIELD_INT_DOM_LATENCY_VAR = 'interDomainLatencyVariation';
export const FIELD_INT_DOM_LATENCY_DIST = 'interDomainLatencyDistribution';
export const FIELD_INT_DOM_THROUGHPUT_DL = 'interDomainThroughputDl';
export const FIELD_INT_DOM_THROUGHPUT_UL = 'interDomainThroughputUl';
export const FIELD_INT_DOM_PKT_LOSS = 'interDomainPacketLoss';
export const FIELD_INT_ZONE_LATENCY = 'interZoneLatency';
export const FIELD_INT_ZONE_LATENCY_VAR = 'interZoneLatencyVariation';
export const FIELD_INT_ZONE_THROUGHPUT_DL = 'interZoneThroughputDl';
export const FIELD_INT_ZONE_THROUGHPUT_UL = 'interZoneThroughputUl';
export const FIELD_INT_ZONE_PKT_LOSS = 'interZonePacketLoss';
export const FIELD_INTRA_ZONE_LATENCY = 'intraZoneLatency';
export const FIELD_INTRA_ZONE_LATENCY_VAR = 'intraZoneLatencyVariation';
export const FIELD_INTRA_ZONE_THROUGHPUT_DL = 'intraZoneThroughputDl';
export const FIELD_INTRA_ZONE_THROUGHPUT_UL = 'intraZoneThroughputUl';
export const FIELD_INTRA_ZONE_PKT_LOSS = 'intraZonePacketLoss';
export const FIELD_TERM_LINK_LATENCY = 'terminalLinkLatency';
export const FIELD_TERM_LINK_LATENCY_VAR = 'terminalLinkLatencyVariation';
export const FIELD_TERM_LINK_THROUGHPUT_DL = 'terminalLinkThroughputDl';
export const FIELD_TERM_LINK_THROUGHPUT_UL = 'terminalLinkThroughputUl';
export const FIELD_TERM_LINK_PKT_LOSS = 'terminalLinkPacketLoss';
export const FIELD_LINK_LATENCY = 'linkLatency';
export const FIELD_LINK_LATENCY_VAR = 'linkLatencyVariation';
export const FIELD_LINK_THROUGHPUT_DL = 'linkThroughputDl';
export const FIELD_LINK_THROUGHPUT_UL = 'linkThroughputUl';
export const FIELD_LINK_PKT_LOSS = 'linkPacketLoss';
export const FIELD_APP_LATENCY = 'appLatency';
export const FIELD_APP_LATENCY_VAR = 'appLatencyVariation';
export const FIELD_APP_THROUGHPUT_DL = 'appThroughput_Dl';
export const FIELD_APP_THROUGHPUT_UL = 'appThroughput_Ul';
export const FIELD_APP_PKT_LOSS = 'appPacketLoss';
export const FIELD_META_DISPLAY_MAP_COLOR = 'metaDisplayMapColor';
export const FIELD_META_DISPLAY_MAP_ICON = 'metaDisplayMapIcon';


export const getElemFieldVal = (elem, field) => {
  return (elem && elem[field]) ? elem[field].val : null;
};

export const setElemFieldVal = (elem, field, val) => {
  if (elem) {
    elem[field] = { val: val, err: null };
  }
};

export const getElemFieldErr = (elem, field) => {
  return (elem && elem[field]) ? elem[field].err : null;
};

export const setElemFieldErr = (elem, field, err) => {
  if (elem) {
    elem[field].err = err;
  }
};

export const resetElem = (elem) => {
  if (elem) {
    elem.editColor = false;
  }
};

export const createElem = name => {
  var elem = {};
  // State
  resetElem(elem);

  // Fields
  setElemFieldVal(elem, FIELD_TYPE, '');
  setElemFieldVal(elem, FIELD_PARENT, '');
  setElemFieldVal(elem, FIELD_NAME, name);
  setElemFieldVal(elem, FIELD_IMAGE, '');
  setElemFieldVal(elem, FIELD_PORT, '');
  setElemFieldVal(elem, FIELD_PROTOCOL, '');
  setElemFieldVal(elem, FIELD_GROUP, '');
  setElemFieldVal(elem, FIELD_INGRESS_SVC_MAP, '');
  setElemFieldVal(elem, FIELD_EGRESS_SVC_MAP, '');
  setElemFieldVal(elem, FIELD_GPU_COUNT, '');
  setElemFieldVal(elem, FIELD_GPU_TYPE, '');
  setElemFieldVal(elem, FIELD_PLACEMENT_ID, '');
  setElemFieldVal(elem, FIELD_ENV_VAR, '');
  setElemFieldVal(elem, FIELD_CMD, '');
  setElemFieldVal(elem, FIELD_CMD_ARGS, '');
  setElemFieldVal(elem, FIELD_EXT_PORT, '');
  setElemFieldVal(elem, FIELD_IS_EXTERNAL, false);
  setElemFieldVal(elem, FIELD_MNC, '');
  setElemFieldVal(elem, FIELD_MCC, '');
  setElemFieldVal(elem, FIELD_DEFAULT_CELL_ID, '');
  setElemFieldVal(elem, FIELD_CELL_ID, '');
  setElemFieldVal(elem, FIELD_NR_CELL_ID, '');
  setElemFieldVal(elem, FIELD_MAC_ID, '');
  setElemFieldVal(elem, FIELD_GEO_LOCATION, '');
  setElemFieldVal(elem, FIELD_GEO_RADIUS, '');
  setElemFieldVal(elem, FIELD_GEO_PATH, '');
  setElemFieldVal(elem, FIELD_GEO_EOP_MODE, '');
  setElemFieldVal(elem, FIELD_GEO_VELOCITY, '');
  setElemFieldVal(elem, FIELD_CHART_ENABLED, false);
  setElemFieldVal(elem, FIELD_CHART_LOC, '');
  setElemFieldVal(elem, FIELD_CHART_VAL, '');
  setElemFieldVal(elem, FIELD_CHART_GROUP, '');
  setElemFieldVal(elem, FIELD_CONNECTED, true);
  setElemFieldVal(elem, FIELD_WIRELESS, false);
  setElemFieldVal(elem, FIELD_WIRELESS_TYPE, '');
  setElemFieldVal(elem, FIELD_INT_DOM_LATENCY, DEFAULT_LATENCY_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_LATENCY_VAR, DEFAULT_LATENCY_JITTER_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_LATENCY_DIST, DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_THROUGHPUT_DL, DEFAULT_THROUGHPUT_DL_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_THROUGHPUT_UL, DEFAULT_THROUGHPUT_UL_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_PKT_LOSS, DEFAULT_PACKET_LOSS_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_ZONE_LATENCY, DEFAULT_LATENCY_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_ZONE_LATENCY_VAR, DEFAULT_LATENCY_JITTER_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_ZONE_THROUGHPUT_DL, DEFAULT_THROUGHPUT_DL_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_ZONE_THROUGHPUT_UL, DEFAULT_THROUGHPUT_UL_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_ZONE_PKT_LOSS, DEFAULT_PACKET_LOSS_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INTRA_ZONE_LATENCY, DEFAULT_LATENCY_INTRA_ZONE);
  setElemFieldVal(elem, FIELD_INTRA_ZONE_LATENCY_VAR, DEFAULT_LATENCY_JITTER_INTRA_ZONE);
  setElemFieldVal(elem, FIELD_INTRA_ZONE_THROUGHPUT_DL, DEFAULT_THROUGHPUT_DL_INTRA_ZONE);
  setElemFieldVal(elem, FIELD_INTRA_ZONE_THROUGHPUT_UL, DEFAULT_THROUGHPUT_UL_INTRA_ZONE);
  setElemFieldVal(elem, FIELD_INTRA_ZONE_PKT_LOSS, DEFAULT_PACKET_LOSS_INTRA_ZONE);
  setElemFieldVal(elem, FIELD_TERM_LINK_LATENCY, DEFAULT_LATENCY_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_TERM_LINK_LATENCY_VAR, DEFAULT_LATENCY_JITTER_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_TERM_LINK_THROUGHPUT_DL, DEFAULT_THROUGHPUT_DL_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_TERM_LINK_THROUGHPUT_UL, DEFAULT_THROUGHPUT_UL_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_TERM_LINK_PKT_LOSS, DEFAULT_PACKET_LOSS_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_LINK_LATENCY, DEFAULT_LATENCY_LINK);
  setElemFieldVal(elem, FIELD_LINK_LATENCY_VAR, DEFAULT_LATENCY_JITTER_LINK);
  setElemFieldVal(elem, FIELD_LINK_THROUGHPUT_DL, DEFAULT_THROUGHPUT_DL_LINK);
  setElemFieldVal(elem, FIELD_LINK_THROUGHPUT_UL, DEFAULT_THROUGHPUT_UL_LINK);
  setElemFieldVal(elem, FIELD_LINK_PKT_LOSS, DEFAULT_PACKET_LOSS_LINK);
  setElemFieldVal(elem, FIELD_APP_LATENCY, DEFAULT_LATENCY_APP);
  setElemFieldVal(elem, FIELD_APP_LATENCY_VAR, DEFAULT_LATENCY_JITTER_APP);
  setElemFieldVal(elem, FIELD_APP_THROUGHPUT_DL, DEFAULT_THROUGHPUT_DL_APP);
  setElemFieldVal(elem, FIELD_APP_THROUGHPUT_UL, DEFAULT_THROUGHPUT_UL_APP);
  setElemFieldVal(elem, FIELD_APP_PKT_LOSS, DEFAULT_PACKET_LOSS_APP);
  setElemFieldVal(elem, FIELD_META_DISPLAY_MAP_COLOR, '');
  setElemFieldVal(elem, FIELD_META_DISPLAY_MAP_ICON, '');

  return elem;
};

export const createUniqueName = (entries, namePrefix) => {
  var increment = 1;
  var isUniqueName = false;
  var suggestedName = namePrefix + String(increment);
  while (!isUniqueName) {
    if (!entries[suggestedName]) {
      isUniqueName = true;
    } else {
      increment++;
      suggestedName = namePrefix + String(increment);
    }
  }
  return suggestedName;
};
