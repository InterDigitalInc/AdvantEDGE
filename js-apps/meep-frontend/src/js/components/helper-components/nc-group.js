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

import React from 'react';
import { Grid, GridCell } from '@rmwc/grid';
import { TextField, TextFieldHelperText } from '@rmwc/textfield';
import { Select } from '@rmwc/select';

import {
  // Field Names
  FIELD_INT_DOM_LATENCY,
  FIELD_INT_DOM_LATENCY_VAR,
  FIELD_INT_DOM_LATENCY_DIST,
  FIELD_INT_DOM_THROUGHPUT_DL,
  FIELD_INT_DOM_THROUGHPUT_UL,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGHPUT_DL,
  FIELD_INT_ZONE_THROUGHPUT_UL,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INTRA_ZONE_LATENCY,
  FIELD_INTRA_ZONE_LATENCY_VAR,
  FIELD_INTRA_ZONE_THROUGHPUT_DL,
  FIELD_INTRA_ZONE_THROUGHPUT_UL,
  FIELD_INTRA_ZONE_PKT_LOSS,
  FIELD_TERM_LINK_LATENCY,
  FIELD_TERM_LINK_LATENCY_VAR,
  FIELD_TERM_LINK_THROUGHPUT_DL,
  FIELD_TERM_LINK_THROUGHPUT_UL,
  FIELD_TERM_LINK_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGHPUT_DL,
  FIELD_LINK_THROUGHPUT_UL,
  FIELD_LINK_PKT_LOSS,
  FIELD_APP_LATENCY,
  FIELD_APP_LATENCY_VAR,
  FIELD_APP_THROUGHPUT_DL,
  FIELD_APP_THROUGHPUT_UL,
  FIELD_APP_PKT_LOSS,
  getElemFieldVal,
  getElemFieldErr
} from '../../util/elem-utils';

import {
  CFG_ELEM_LATENCY,
  CFG_ELEM_LATENCY_VAR,
  CFG_ELEM_LATENCY_DIST,
  CFG_ELEM_PKT_LOSS,
  CFG_ELEM_THROUGHPUT_DL,
  CFG_ELEM_THROUGHPUT_UL,

  // NC Group Prefixes
  PREFIX_INT_DOM,
  PREFIX_INT_ZONE,
  PREFIX_INTRA_ZONE,
  PREFIX_TERM_LINK,
  PREFIX_LINK,
  PREFIX_APP
} from '../../meep-constants';

const MIN_LATENCY_VALUE = 0;
const MAX_LATENCY_VALUE = 250000;

const MIN_LATENCY_VARIATION_VALUE = 0;
const MAX_LATENCY_VARIATION_VALUE = 250000;

const MIN_THROUGHPUT_VALUE = 1;
const MAX_THROUGHPUT_VALUE = 1000000;

const MIN_PACKET_LOSS_VALUE = 0;
const MAX_PACKET_LOSS_VALUE = 100;

const isInt = val => {
  return Number(val) % 1 === 0 && val[val.length - 1] !== '.';
};

const validateLatency = val => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val) {
    if (isNaN(val)) {
      return 'Latency value should be a number';
    }
    if ((val !== '' && val < MIN_LATENCY_VALUE) || val > MAX_LATENCY_VALUE) {
      return `Out of range (${MIN_LATENCY_VALUE}-${MAX_LATENCY_VALUE})`;
    }
    if (!isInt(val)) {
      return 'Latency value should be an integer';
    }
  }
  return null;
};

const validateLatencyVariation = val => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val) {
    if (isNaN(val)) {
      return 'Latency variation should be a number';
    }
    if (
      (val !== '' && val < MIN_LATENCY_VARIATION_VALUE) ||
      val > MAX_LATENCY_VARIATION_VALUE
    ) {
      return `Out of range (${MIN_LATENCY_VARIATION_VALUE}-${MAX_LATENCY_VARIATION_VALUE})`;
    }
    if (!isInt(val)) {
      return 'Latency variation value should be an integer';
    }
  }
  return null;
};

const validatePacketLoss = val => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val && val !== '0') {
    if (isNaN(val)) {
      return 'Packet loss value should be a number';
    }
    if (
      (val !== '' && val < MIN_PACKET_LOSS_VALUE) ||
      val > MAX_PACKET_LOSS_VALUE
    ) {
      return `Out of range (${MIN_PACKET_LOSS_VALUE}-${MAX_PACKET_LOSS_VALUE})`;
    }
    if (!Number(val) || val[val.length - 1] === '.') {
      return 'Must be a number with at most 7 decimal places';
    }
    if (val.length > 9) {
      return 'Too many decimal places';
    }
  }
  return null;
};

const validateThroughput = val => {
  if (val === '' || val === undefined) {
    return 'Value is required';
  }
  if (val) {
    if (isNaN(val)) {
      return 'Throughput value should be a number';
    }
    if (
      (val !== '' && val < MIN_THROUGHPUT_VALUE) ||
      val > MAX_THROUGHPUT_VALUE
    ) {
      return `Out of range (${MIN_THROUGHPUT_VALUE}-${MAX_THROUGHPUT_VALUE})`;
    }
  }
  return null;
};

const NCLayout = props => {
  return (
    <>
      <Grid>
        <GridCell span="4">{props.latencyComponent}</GridCell>
        <GridCell span="4">{props.latencyVariationComponent}</GridCell>
        <GridCell span="4">{props.packetLossComponent}</GridCell>
      </Grid>

      {props.latencyDistributionComponent === null ? null :
        <Grid style={{ marginTop: 10 }}>
          <GridCell span="12">{props.latencyDistributionComponent}</GridCell>
        </Grid>
      }

      <Grid>
        <GridCell span="6">{props.throughputDlComponent}</GridCell>
        <GridCell span="6">{props.throughputUlComponent}</GridCell>
      </Grid>
    </>
  );
};

const NCGroup = ({ prefix, onUpdate, element }) => {

  const handleEvent = (event, fieldName, validate) => {
    var err = validate ? validate(event.target.value) : null;
    var val =
      event.target.value && !err
        ? Number(event.target.value)
        : event.target.value;
    onUpdate(fieldName, val, err);
    event.preventDefault();
  };

  // Retrieve field names
  var latencyFieldName = null;
  var latencyVarFieldName = null;
  var latencyDistFieldName = null;
  var throughputDlFieldName = null;
  var throughputUlFieldName = null;
  var packetLossFieldName = null;
  switch (prefix) {
  case PREFIX_INT_DOM:
    latencyFieldName = FIELD_INT_DOM_LATENCY;
    latencyVarFieldName = FIELD_INT_DOM_LATENCY_VAR;
    latencyDistFieldName = FIELD_INT_DOM_LATENCY_DIST;
    throughputDlFieldName = FIELD_INT_DOM_THROUGHPUT_DL;
    throughputUlFieldName = FIELD_INT_DOM_THROUGHPUT_UL;
    packetLossFieldName = FIELD_INT_DOM_PKT_LOSS;
    break;
  case PREFIX_INT_ZONE:
    latencyFieldName = FIELD_INT_ZONE_LATENCY;
    latencyVarFieldName = FIELD_INT_ZONE_LATENCY_VAR;
    throughputDlFieldName = FIELD_INT_ZONE_THROUGHPUT_DL;
    throughputUlFieldName = FIELD_INT_ZONE_THROUGHPUT_UL;
    packetLossFieldName = FIELD_INT_ZONE_PKT_LOSS;
    break;
  case PREFIX_INTRA_ZONE:
    latencyFieldName = FIELD_INTRA_ZONE_LATENCY;
    latencyVarFieldName = FIELD_INTRA_ZONE_LATENCY_VAR;
    throughputDlFieldName = FIELD_INTRA_ZONE_THROUGHPUT_DL;
    throughputUlFieldName = FIELD_INTRA_ZONE_THROUGHPUT_UL;
    packetLossFieldName = FIELD_INTRA_ZONE_PKT_LOSS;
    break;
  case PREFIX_TERM_LINK:
    latencyFieldName = FIELD_TERM_LINK_LATENCY;
    latencyVarFieldName = FIELD_TERM_LINK_LATENCY_VAR;
    throughputDlFieldName = FIELD_TERM_LINK_THROUGHPUT_DL;
    throughputUlFieldName = FIELD_TERM_LINK_THROUGHPUT_UL;
    packetLossFieldName = FIELD_TERM_LINK_PKT_LOSS;
    break;
  case PREFIX_LINK:
    latencyFieldName = FIELD_LINK_LATENCY;
    latencyVarFieldName = FIELD_LINK_LATENCY_VAR;
    throughputDlFieldName = FIELD_LINK_THROUGHPUT_DL;
    throughputUlFieldName = FIELD_LINK_THROUGHPUT_UL;
    packetLossFieldName = FIELD_LINK_PKT_LOSS;
    break;
  case PREFIX_APP:
    latencyFieldName = FIELD_APP_LATENCY;
    latencyVarFieldName = FIELD_APP_LATENCY_VAR;
    throughputDlFieldName = FIELD_APP_THROUGHPUT_DL;
    throughputUlFieldName = FIELD_APP_THROUGHPUT_UL;
    packetLossFieldName = FIELD_APP_PKT_LOSS;
    break;
  default:
    return null;
  }

  const latencyComponent = (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={'Latency (ms)'}
        onChange={e => handleEvent(e, latencyFieldName, validateLatency)}
        value={getElemFieldVal(element, latencyFieldName)}
        invalid={getElemFieldErr(element, latencyFieldName) ? true : false}
        data-cy={CFG_ELEM_LATENCY}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(element, latencyFieldName)}</span>
      </TextFieldHelperText>
    </>
  );

  const latencyVariationComponent = (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={'Jitter (ms)'}
        onChange={e =>
          handleEvent(e, latencyVarFieldName, validateLatencyVariation)
        }
        value={getElemFieldVal(element, latencyVarFieldName)}
        invalid={getElemFieldErr(element, latencyVarFieldName) ? true : false}
        data-cy={CFG_ELEM_LATENCY_VAR}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(element, latencyVarFieldName)}</span>
      </TextFieldHelperText>
    </>
  );

  const latencyDistributionComponent = (
    <>
      <Select
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label="Latency Distribution"
        value={getElemFieldVal(element, latencyDistFieldName)}
        options={['Normal', 'Pareto', 'ParetoNormal', 'Uniform']}
        onChange={e => onUpdate(latencyDistFieldName, e.target.value, null)}
        data-cy={CFG_ELEM_LATENCY_DIST}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(element, latencyDistFieldName)}</span>
      </TextFieldHelperText>
    </>
  );

  const packetLossComponent = (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={'Packet Loss (%)'}
        onChange={e => handleEvent(e, packetLossFieldName, validatePacketLoss)}
        value={getElemFieldVal(element, packetLossFieldName)}
        invalid={getElemFieldErr(element, packetLossFieldName) ? true : false}
        data-cy={CFG_ELEM_PKT_LOSS}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(element, packetLossFieldName)}</span>
      </TextFieldHelperText>
    </>
  );

  const throughputDlComponent = (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={'DL Throughput (Mbps)'}
        onChange={e => handleEvent(e, throughputDlFieldName, validateThroughput)}
        value={getElemFieldVal(element, throughputDlFieldName)}
        invalid={getElemFieldErr(element, throughputDlFieldName) ? true : false}
        data-cy={CFG_ELEM_THROUGHPUT_DL}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(element, throughputDlFieldName)}</span>
      </TextFieldHelperText>
    </>
  );

  const throughputUlComponent = (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={'UL Throughput (Mbps)'}
        onChange={e => handleEvent(e, throughputUlFieldName, validateThroughput)}
        value={getElemFieldVal(element, throughputUlFieldName)}
        invalid={getElemFieldErr(element, throughputUlFieldName) ? true : false}
        data-cy={CFG_ELEM_THROUGHPUT_UL}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(element, throughputUlFieldName)}</span>
      </TextFieldHelperText>
    </>
  );

  return (
    <NCLayout
      latencyComponent={latencyComponent}
      latencyVariationComponent={latencyVariationComponent}
      latencyDistributionComponent={(prefix === PREFIX_INT_DOM) ? latencyDistributionComponent : null}
      packetLossComponent={packetLossComponent}
      throughputDlComponent={throughputDlComponent}
      throughputUlComponent={throughputUlComponent}
    />
  );
};

export default NCGroup;
