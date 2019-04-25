/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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
  DEFAULT_LATENCY_INTER_EDGE,
  DEFAULT_LATENCY_JITTER_INTER_EDGE,
  DEFAULT_THROUGHPUT_INTER_EDGE,
  DEFAULT_PACKET_LOSS_INTER_EDGE,
  DEFAULT_LATENCY_INTER_FOG,
  DEFAULT_LATENCY_JITTER_INTER_FOG,
  DEFAULT_THROUGHPUT_INTER_FOG,
  DEFAULT_PACKET_LOSS_INTER_FOG,
  DEFAULT_LATENCY_EDGE_FOG,
  DEFAULT_LATENCY_JITTER_EDGE_FOG,
  DEFAULT_THROUGHPUT_EDGE_FOG,
  DEFAULT_PACKET_LOSS_EDGE_FOG,
  DEFAULT_LATENCY_TERMINAL_LINK,
  DEFAULT_LATENCY_JITTER_TERMINAL_LINK,
  DEFAULT_THROUGHPUT_TERMINAL_LINK,
  DEFAULT_PACKET_LOSS_TERMINAL_LINK,
} from '../meep-constants';

// Network Element Fields
export const FIELD_TYPE = 'elementType';
export const FIELD_PARENT = 'parent';
export const FIELD_NAME = 'name';
export const FIELD_IMAGE = 'image';
export const FIELD_PORT = 'port';
export const FIELD_PROTOCOL = 'protocol';
export const FIELD_GROUP = 'group';
export const FIELD_SVC_MAP = 'ingressServiceMap';
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
export const FIELD_INT_EDGE_LATENCY = 'interEdgeLatency';
export const FIELD_INT_EDGE_LATENCY_VAR = 'interEdgeLatencyVariation';
export const FIELD_INT_EDGE_THROUGPUT = 'interEdgeThroughput';
export const FIELD_INT_EDGE_PKT_LOSS = 'interEdgePacketLoss';
export const FIELD_INT_FOG_LATENCY = 'interFogLatency';
export const FIELD_INT_FOG_LATENCY_VAR = 'interFogLatencyVariation';
export const FIELD_INT_FOG_THROUGPUT = 'interFogThroughput';
export const FIELD_INT_FOG_PKT_LOSS = 'interFogPacketLoss';
export const FIELD_EDGE_FOG_LATENCY = 'edgeFogLatency';
export const FIELD_EDGE_FOG_LATENCY_VAR = 'edgeFogLatencyVariation';
export const FIELD_EDGE_FOG_THROUGPUT = 'edgeFogThroughput';
export const FIELD_EDGE_FOG_PKT_LOSS = 'edgeFogPacketLoss';
export const FIELD_LINK_LATENCY = 'terminalLinkLatency';
export const FIELD_LINK_LATENCY_VAR = 'terminalLinkLatencyVariation';
export const FIELD_LINK_THROUGPUT = 'terminalLinkThroughput';
export const FIELD_LINK_PKT_LOSS = 'terminalLinkPacketLoss';

export const getElemFieldVal = (elem, field) => {
  return (elem[field]) ? elem[field].val : null;
};

export const setElemFieldVal = (elem, field, val) => {
  elem[field] = { val: val, err: null };
};

export const getElemFieldErr = (elem, field) => {
  return (elem[field]) ? elem[field].err : null;
};

export const setElemFieldErr = (elem, field, err) => {
  elem[field].err = err;
};

export const createElem = (name) => {
  var elem = {};
  setElemFieldVal(elem, FIELD_TYPE,                   '');
  setElemFieldVal(elem, FIELD_PARENT,                 '');
  setElemFieldVal(elem, FIELD_NAME,                   name);
  setElemFieldVal(elem, FIELD_IMAGE,                  '');
  setElemFieldVal(elem, FIELD_PORT,                   '');
  setElemFieldVal(elem, FIELD_PROTOCOL,               '');
  setElemFieldVal(elem, FIELD_GROUP,                  '');
  setElemFieldVal(elem, FIELD_SVC_MAP,                '');
  setElemFieldVal(elem, FIELD_ENV_VAR,                '');
  setElemFieldVal(elem, FIELD_CMD,                    '');
  setElemFieldVal(elem, FIELD_CMD_ARGS,               '');
  setElemFieldVal(elem, FIELD_EXT_PORT,               '');
  setElemFieldVal(elem, FIELD_IS_EXTERNAL,            false);
  setElemFieldVal(elem, FIELD_CHART_ENABLED,          false);
  setElemFieldVal(elem, FIELD_CHART_LOC,              '');
  setElemFieldVal(elem, FIELD_CHART_VAL,              '');
  setElemFieldVal(elem, FIELD_CHART_GROUP,            '');
  setElemFieldVal(elem, FIELD_INT_DOM_LATENCY,        DEFAULT_LATENCY_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_LATENCY_VAR,    DEFAULT_LATENCY_JITTER_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_THROUGPUT,      DEFAULT_THROUGHPUT_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_DOM_PKT_LOSS,       DEFAULT_PACKET_LOSS_INTER_DOMAIN);
  setElemFieldVal(elem, FIELD_INT_ZONE_LATENCY,       DEFAULT_LATENCY_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_ZONE_LATENCY_VAR,   DEFAULT_LATENCY_JITTER_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_ZONE_THROUGPUT,     DEFAULT_THROUGHPUT_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_ZONE_PKT_LOSS,      DEFAULT_PACKET_LOSS_INTER_ZONE);
  setElemFieldVal(elem, FIELD_INT_EDGE_LATENCY,       DEFAULT_LATENCY_INTER_EDGE);
  setElemFieldVal(elem, FIELD_INT_EDGE_LATENCY_VAR,   DEFAULT_LATENCY_JITTER_INTER_EDGE);
  setElemFieldVal(elem, FIELD_INT_EDGE_THROUGPUT,     DEFAULT_THROUGHPUT_INTER_EDGE);
  setElemFieldVal(elem, FIELD_INT_EDGE_PKT_LOSS,      DEFAULT_PACKET_LOSS_INTER_EDGE);
  setElemFieldVal(elem, FIELD_INT_FOG_LATENCY,        DEFAULT_LATENCY_INTER_FOG);
  setElemFieldVal(elem, FIELD_INT_FOG_LATENCY_VAR,    DEFAULT_LATENCY_JITTER_INTER_FOG);
  setElemFieldVal(elem, FIELD_INT_FOG_THROUGPUT,      DEFAULT_THROUGHPUT_INTER_FOG);
  setElemFieldVal(elem, FIELD_INT_FOG_PKT_LOSS,       DEFAULT_PACKET_LOSS_INTER_FOG);
  setElemFieldVal(elem, FIELD_EDGE_FOG_LATENCY,       DEFAULT_LATENCY_EDGE_FOG);
  setElemFieldVal(elem, FIELD_EDGE_FOG_LATENCY_VAR,   DEFAULT_LATENCY_JITTER_EDGE_FOG);
  setElemFieldVal(elem, FIELD_EDGE_FOG_THROUGPUT,     DEFAULT_THROUGHPUT_EDGE_FOG);
  setElemFieldVal(elem, FIELD_EDGE_FOG_PKT_LOSS,      DEFAULT_PACKET_LOSS_EDGE_FOG);
  setElemFieldVal(elem, FIELD_LINK_LATENCY,           DEFAULT_LATENCY_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_LINK_LATENCY_VAR,       DEFAULT_LATENCY_JITTER_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_LINK_THROUGPUT,         DEFAULT_THROUGHPUT_TERMINAL_LINK);
  setElemFieldVal(elem, FIELD_LINK_PKT_LOSS,          DEFAULT_PACKET_LOSS_TERMINAL_LINK);

  return elem;
};
