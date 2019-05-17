/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import React from 'react';
import { Grid, GridCell  } from '@rmwc/grid';
import { TextField, TextFieldHelperText } from '@rmwc/textfield';

import {
  firstLetterUpper
} from '../../util/stringManipulation';

import {
  // Field Names
  FIELD_INT_DOM_LATENCY,
  FIELD_INT_DOM_LATENCY_VAR,
  FIELD_INT_DOM_THROUGPUT,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGPUT,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INT_EDGE_LATENCY,
  FIELD_INT_EDGE_LATENCY_VAR,
  FIELD_INT_EDGE_THROUGPUT,
  FIELD_INT_EDGE_PKT_LOSS,
  FIELD_INT_FOG_LATENCY,
  FIELD_INT_FOG_LATENCY_VAR,
  FIELD_INT_FOG_THROUGPUT,
  FIELD_INT_FOG_PKT_LOSS,
  FIELD_EDGE_FOG_LATENCY,
  FIELD_EDGE_FOG_LATENCY_VAR,
  FIELD_EDGE_FOG_THROUGPUT,
  FIELD_EDGE_FOG_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGPUT,
  FIELD_LINK_PKT_LOSS,

  getElemFieldVal,
  getElemFieldErr
} from '../../util/elem-utils';

import {
  CFG_ELEM_LATENCY,
  CFG_ELEM_LATENCY_VAR,
  CFG_ELEM_PKT_LOSS,
  CFG_ELEM_THROUGHPUT,

  // NC Group Prefixes
  PREFIX_INT_DOM,
  PREFIX_INT_ZONE,
  PREFIX_INT_EDGE,
  PREFIX_INT_FOG,
  PREFIX_EDGE_FOG,
  PREFIX_TERM_LINK

} from '../../meep-constants';

const MIN_LATENCY_VALUE = 0;
const MAX_LATENCY_VALUE = 250000;

const MIN_LATENCY_VARIATION_VALUE = 0;
const MAX_LATENCY_VARIATION_VALUE = 250000;

const MIN_THROUGHPUT_VALUE = 1;
const MAX_THROUGHPUT_VALUE = 1000000;

const MIN_PACKET_LOSS_VALUE = 0;
const MAX_PACKET_LOSS_VALUE = 100;

const isInt = (val) => {
  return Number(val) %1 === 0 && val[val.length-1] !== '.';
};

const validateLatency = (val) => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val) {
    if (isNaN(val)) {return 'Latency value should be a number';}
    if ((val !== '') && val < MIN_LATENCY_VALUE || val > MAX_LATENCY_VALUE) {
      return `Out of range (${MIN_LATENCY_VALUE}-${MAX_LATENCY_VALUE})`;
    }
    if (!isInt(val)) {return 'Latency value should be an integer';}
    
  }
  return null;
};

const validateLatencyVariation = (val) => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val) {
    if (isNaN(val)) {return 'Latency variation should be a number';}
    if ((val !== '') && val < MIN_LATENCY_VARIATION_VALUE || val > MAX_LATENCY_VARIATION_VALUE) {
      return `Out of range (${MIN_LATENCY_VARIATION_VALUE}-${MAX_LATENCY_VARIATION_VALUE})`;
    }
    if (!isInt(val)) {
      return 'Latency variation value should be an integer';
    }
  }
  return null;
};

const validatePacketLoss = (val) => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val && val !== '0') {
    if (isNaN(val)) {return 'Packet loss value should be a number';}
    if ((val !== '') && val < MIN_PACKET_LOSS_VALUE || val > MAX_PACKET_LOSS_VALUE) {
      return `Out of range (${MIN_PACKET_LOSS_VALUE}-${MAX_PACKET_LOSS_VALUE})`;
    }
    if (!Number(val) || val[val.length-1] === '.') {
      return 'Must be a number with at most 7 decimal places';
    }
    if (val.length > 9) {
      return 'Too many decimal places';
    }
  }
  return null;
};

const validateThroughput = (val) => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val) {
    if (isNaN(val)) {return 'Throughput value should be a number';}
    if ((val !== '') && val < MIN_THROUGHPUT_VALUE || val > MAX_THROUGHPUT_VALUE) {
      return `Out of range (${MIN_THROUGHPUT_VALUE}-${MAX_THROUGHPUT_VALUE})`;
    }
  }
  return null;
};

const NCGroup = ({prefix, onUpdate, element}) => {
  const formLabel = (valueName) => {
    const space = prefix ? ' ' : '';
    return firstLetterUpper(prefix) + space + valueName;
  };

  const handleEvent = (event, fieldName, validate) => {
    var err = (validate) ? validate(event.target.value) : null;
    var val = (event.target.value && !err) ? Number(event.target.value) : event.target.value;
    onUpdate(fieldName, val, err);
    event.preventDefault();
  };

  // Retrieve field names
  var latencyFieldName = null;
  var latencyVarFieldName = null;
  var throughputFieldName = null;
  var packetLossFieldName = null;
  switch (prefix) {
  case PREFIX_INT_DOM:
    latencyFieldName = FIELD_INT_DOM_LATENCY;
    latencyVarFieldName = FIELD_INT_DOM_LATENCY_VAR;
    throughputFieldName = FIELD_INT_DOM_THROUGPUT;
    packetLossFieldName = FIELD_INT_DOM_PKT_LOSS;
    break;
  case PREFIX_INT_ZONE:
    latencyFieldName = FIELD_INT_ZONE_LATENCY;
    latencyVarFieldName = FIELD_INT_ZONE_LATENCY_VAR;
    throughputFieldName = FIELD_INT_ZONE_THROUGPUT;
    packetLossFieldName = FIELD_INT_ZONE_PKT_LOSS;
    break;
  case PREFIX_INT_EDGE:
    latencyFieldName = FIELD_INT_EDGE_LATENCY;
    latencyVarFieldName = FIELD_INT_EDGE_LATENCY_VAR;
    throughputFieldName = FIELD_INT_EDGE_THROUGPUT;
    packetLossFieldName = FIELD_INT_EDGE_PKT_LOSS;
    break;
  case PREFIX_INT_FOG:
    latencyFieldName = FIELD_INT_FOG_LATENCY;
    latencyVarFieldName = FIELD_INT_FOG_LATENCY_VAR;
    throughputFieldName = FIELD_INT_FOG_THROUGPUT;
    packetLossFieldName = FIELD_INT_FOG_PKT_LOSS;
    break;
  case PREFIX_EDGE_FOG:
    latencyFieldName = FIELD_EDGE_FOG_LATENCY;
    latencyVarFieldName = FIELD_EDGE_FOG_LATENCY_VAR;
    throughputFieldName = FIELD_EDGE_FOG_THROUGPUT;
    packetLossFieldName = FIELD_EDGE_FOG_PKT_LOSS;
    break;
  case PREFIX_TERM_LINK:
    latencyFieldName = FIELD_LINK_LATENCY;
    latencyVarFieldName = FIELD_LINK_LATENCY_VAR;
    throughputFieldName = FIELD_LINK_THROUGPUT;
    packetLossFieldName = FIELD_LINK_PKT_LOSS;
    break;
  default:
    return null;
  }

  return (
    <div>
      <Grid>
        <GridCell span="6">
          <TextField outlined style={{width: '100%'}}
            label={formLabel('Latency') + ' (ms)'}
            onChange={(e) => handleEvent(e, latencyFieldName, validateLatency)}
            value={getElemFieldVal(element, latencyFieldName)}
            invalid={getElemFieldErr(element, latencyFieldName) ? true : false}
            data-cy={CFG_ELEM_LATENCY}
          />
          <TextFieldHelperText validationMsg={true}>
            <span>
              {getElemFieldErr(element, latencyFieldName)}
            </span>
          </TextFieldHelperText>
        </GridCell>
        <GridCell span="6">
          <TextField outlined style={{width: '100%'}}
            label={formLabel('Latency Variation') + ' (ms)'}
            onChange={(e) => handleEvent(e, latencyVarFieldName, validateLatencyVariation)}
            value={getElemFieldVal(element, latencyVarFieldName)}
            invalid={getElemFieldErr(element, latencyVarFieldName) ? true : false}
            data-cy={CFG_ELEM_LATENCY_VAR}
          />
          <TextFieldHelperText validationMsg={true}>
            <span>
              {getElemFieldErr(element, latencyVarFieldName)}
            </span>
          </TextFieldHelperText>
        </GridCell>
      </Grid>

      <Grid style={{marginBottom: 10}}>
        <GridCell span="6">
          <TextField outlined style={{width: '100%'}}
            label={formLabel('Packet Loss') + ' (%)'}
            onChange={(e) =>  handleEvent(e, packetLossFieldName, validatePacketLoss)}
            value={getElemFieldVal(element, packetLossFieldName)}
            invalid={getElemFieldErr(element, packetLossFieldName) ? true : false}
            data-cy={CFG_ELEM_PKT_LOSS}
          />
          <TextFieldHelperText validationMsg={true}>
            <span>
              {getElemFieldErr(element, packetLossFieldName)}
            </span>
          </TextFieldHelperText>

        </GridCell>
        <GridCell span="6">
          <TextField outlined style={{width: '100%'}}
            label={formLabel('Throughput') + ' Mbps'}
            onChange={(e) =>  handleEvent(e, throughputFieldName, validateThroughput)}
            value={getElemFieldVal(element, throughputFieldName)}
            invalid={getElemFieldErr(element, throughputFieldName) ? true : false}
            data-cy={CFG_ELEM_THROUGHPUT}
          />
          <TextFieldHelperText validationMsg={true}>
            <span>
              {getElemFieldErr(element, throughputFieldName)}
            </span>
          </TextFieldHelperText>
        </GridCell>
      </Grid>
    </div>
  );
};

export default NCGroup;
