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
  DEFAULT_THROUGHPUT_INTER_DOMAIN,
  DEFAULT_PACKET_LOSS_INTER_DOMAIN,
  DEFAULT_LATENCY_INTER_ZONE,
  DEFAULT_LATENCY_JITTER_INTER_ZONE,
  DEFAULT_THROUGHPUT_INTER_ZONE,
  DEFAULT_PACKET_LOSS_INTER_ZONE,
  DEFAULT_LATENCY_INTRA_ZONE,
  DEFAULT_LATENCY_JITTER_INTRA_ZONE,
  DEFAULT_THROUGHPUT_INTRA_ZONE,
  DEFAULT_PACKET_LOSS_INTRA_ZONE,
  DEFAULT_LATENCY_TERMINAL_LINK,
  DEFAULT_LATENCY_JITTER_TERMINAL_LINK,
  DEFAULT_THROUGHPUT_TERMINAL_LINK,
  DEFAULT_PACKET_LOSS_TERMINAL_LINK,
  DEFAULT_LATENCY_LINK,
  DEFAULT_LATENCY_JITTER_LINK,
  DEFAULT_THROUGHPUT_LINK,
  DEFAULT_PACKET_LOSS_LINK,
  DEFAULT_LATENCY_APP,
  DEFAULT_LATENCY_JITTER_APP,
  DEFAULT_THROUGHPUT_APP,
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
export const FIELD_CHART_ENABLED = 'userChartEnabled';
export const FIELD_CHART_LOC = 'userChartLocation';
export const FIELD_CHART_VAL = 'userChartAlternateValues';
export const FIELD_CHART_GROUP = 'userChartGroup';
export const FIELD_INT_DOM_LATENCY = 'interDomainLatency';
export const FIELD_INT_DOM_LATENCY_VAR = 'interDomainLatencyVariation';
export const FIELD_INT_DOM_THROUGPUT = 'interDomainThroughput';
export const FIELD_INT_DOM_PKT_LOSS = 'interDomainPacketLoss';
export const FIELD_INT_ZONE_LATENCY = 'interZoneLatency';
export const FIELD_INT_ZONE_LATENCY_VAR = 'interZoneLatencyVariation';
export const FIELD_INT_ZONE_THROUGPUT = 'interZoneThroughput';
export const FIELD_INT_ZONE_PKT_LOSS = 'interZonePacketLoss';
export const FIELD_INTRA_ZONE_LATENCY = 'intraZoneLatency';
export const FIELD_INTRA_ZONE_LATENCY_VAR = 'intraZoneLatencyVariation';
export const FIELD_INTRA_ZONE_THROUGPUT = 'intraZoneThroughput';
export const FIELD_INTRA_ZONE_PKT_LOSS = 'intraZonePacketLoss';
export const FIELD_TERM_LINK_LATENCY = 'terminalLinkLatency';
export const FIELD_TERM_LINK_LATENCY_VAR = 'terminalLinkLatencyVariation';
export const FIELD_TERM_LINK_THROUGPUT = 'terminalLinkThroughput';
export const FIELD_TERM_LINK_PKT_LOSS = 'terminalLinkPacketLoss';
export const FIELD_LINK_LATENCY = 'linkLatency';
export const FIELD_LINK_LATENCY_VAR = 'linkLatencyVariation';
export const FIELD_LINK_THROUGPUT = 'linkThroughput';
export const FIELD_LINK_PKT_LOSS = 'linkPacketLoss';
export const FIELD_APP_LATENCY = 'appLatency';
export const FIELD_APP_LATENCY_VAR = 'appLatencyVariation';
export const FIELD_APP_THROUGPUT = 'appThroughput';
export const FIELD_APP_PKT_LOSS = 'appPacketLoss';

export const getElemFieldVal = (elem, field) => {
  return elem[field] ? elem[field].val : null;
};

export const setElemFieldVal = (elem, field, val) => {
  elem[field] = { val: val, err: null };
};

export const getElemFieldErr = (elem, field) => {
  return elem[field] ? elem[field].err : null;
};

export const setElemFieldErr = (elem, field, err) => {
  elem[field].err = err;
};

export const createElem = name => {
  var elem = {};
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
  setElemFieldVal(elem, FIELD_CHART_ENABLED, false);
  setElemFieldVal(elem, FIELD_CHART_LOC, '');
  setElemFieldVal(elem, FIELD_CHART_VAL, '');
  setElemFieldVal(elem, FIELD_CHART_GROUP, '');
  setElemFieldVal(elem, FIELD_INT_DOM_LATENCY, DEFAULT_LATENCY_INTER_DOMAIN);
  setElemFieldVal(
    elem,
    FIELD_INT_DOM_LATENCY_VAR,
    DEFAULT_LATENCY_JITTER_INTER_DOMAIN
  );
  setElemFieldVal(
    elem,
    FIELD_INT_DOM_THROUGPUT,
    DEFAULT_THROUGHPUT_INTER_DOMAIN
  );
  setElemFieldVal(
    elem,
    FIELD_INT_DOM_PKT_LOSS,
    DEFAULT_PACKET_LOSS_INTER_DOMAIN
  );
  setElemFieldVal(elem, FIELD_INT_ZONE_LATENCY, DEFAULT_LATENCY_INTER_ZONE);
  setElemFieldVal(
    elem,
    FIELD_INT_ZONE_LATENCY_VAR,
    DEFAULT_LATENCY_JITTER_INTER_ZONE
  );
  setElemFieldVal(
    elem,
    FIELD_INT_ZONE_THROUGPUT,
    DEFAULT_THROUGHPUT_INTER_ZONE
  );
  setElemFieldVal(
    elem,
    FIELD_INT_ZONE_PKT_LOSS,
    DEFAULT_PACKET_LOSS_INTER_ZONE
  );
  setElemFieldVal(elem, FIELD_INTRA_ZONE_LATENCY, DEFAULT_LATENCY_INTRA_ZONE);
  setElemFieldVal(
    elem,
    FIELD_INTRA_ZONE_LATENCY_VAR,
    DEFAULT_LATENCY_JITTER_INTRA_ZONE
  );
  setElemFieldVal(elem, FIELD_INTRA_ZONE_THROUGPUT, DEFAULT_THROUGHPUT_INTRA_ZONE);
  setElemFieldVal(elem, FIELD_INTRA_ZONE_PKT_LOSS, DEFAULT_PACKET_LOSS_INTRA_ZONE);
  setElemFieldVal(elem, FIELD_TERM_LINK_LATENCY, DEFAULT_LATENCY_TERMINAL_LINK);
  setElemFieldVal(
    elem,
    FIELD_TERM_LINK_LATENCY_VAR,
    DEFAULT_LATENCY_JITTER_TERMINAL_LINK
  );
  setElemFieldVal(
    elem,
    FIELD_TERM_LINK_THROUGPUT,
    DEFAULT_THROUGHPUT_TERMINAL_LINK
  );
  setElemFieldVal(
    elem,
    FIELD_TERM_LINK_PKT_LOSS,
    DEFAULT_PACKET_LOSS_TERMINAL_LINK
  );
  setElemFieldVal(elem, FIELD_LINK_LATENCY, DEFAULT_LATENCY_LINK);
  setElemFieldVal(elem, FIELD_LINK_LATENCY_VAR, DEFAULT_LATENCY_JITTER_LINK);
  setElemFieldVal(elem, FIELD_LINK_THROUGPUT, DEFAULT_THROUGHPUT_LINK);
  setElemFieldVal(elem, FIELD_LINK_PKT_LOSS, DEFAULT_PACKET_LOSS_LINK);
  setElemFieldVal(elem, FIELD_APP_LATENCY, DEFAULT_LATENCY_APP);
  setElemFieldVal(elem, FIELD_APP_LATENCY_VAR, DEFAULT_LATENCY_JITTER_APP);
  setElemFieldVal(elem, FIELD_APP_THROUGPUT, DEFAULT_THROUGHPUT_APP);
  setElemFieldVal(elem, FIELD_APP_PKT_LOSS, DEFAULT_PACKET_LOSS_APP);

  return elem;
};
