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

import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component } from 'react';
import { Select } from '@rmwc/select';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Button } from '@rmwc/button';
import { TextField, TextFieldHelperText } from '@rmwc/textfield';
import { Checkbox } from '@rmwc/checkbox';
import { Typography } from '@rmwc/typography';

import { updateObject } from '../../util/object-util';
import { createUniqueName } from '../../util/elem-utils';

import IDSelect from '../../components/helper-components/id-select';
import CancelApplyPair from '../../components/helper-components/cancel-apply-pair';
import NCGroup from '../../components/helper-components/nc-group';

import {
  // Network Element Fields
  FIELD_TYPE,
  FIELD_PARENT,
  FIELD_NAME,
  FIELD_IMAGE,
  FIELD_PORT,
  FIELD_PROTOCOL,
  FIELD_GROUP,
  FIELD_GPU_COUNT,
  FIELD_GPU_TYPE,
  FIELD_PLACEMENT_ID,
  FIELD_ENV_VAR,
  FIELD_CMD,
  FIELD_CMD_ARGS,
  FIELD_EXT_PORT,
  FIELD_IS_EXTERNAL,
  FIELD_MCC,
  FIELD_MNC,
  FIELD_DEFAULT_CELL_ID,
  FIELD_CELL_ID,
  FIELD_CHART_ENABLED,
  FIELD_CHART_LOC,
  FIELD_CHART_VAL,
  FIELD_CHART_GROUP,
  getElemFieldVal,
  setElemFieldVal,
  getElemFieldErr,
  setElemFieldErr
} from '../../util/elem-utils';

import {
  CFG_ELEM_MODE_NEW,
  CFG_ELEM_MODE_EDIT,
  CFG_ELEM_MODE_CLONE,
  cfgElemUpdate,
  cfgElemClone
} from '../../state/cfg';

import {
  TYPE_CFG,

  // Network element types
  ELEMENT_TYPE_SCENARIO,
  ELEMENT_TYPE_OPERATOR,
  ELEMENT_TYPE_OPERATOR_GENERIC,
  ELEMENT_TYPE_OPERATOR_CELL,
  ELEMENT_TYPE_ZONE,
  ELEMENT_TYPE_POA,
  ELEMENT_TYPE_POA_GENERIC,
  ELEMENT_TYPE_POA_CELL,
  ELEMENT_TYPE_DC,
  ELEMENT_TYPE_CN,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_UE,
  ELEMENT_TYPE_MECSVC,
  ELEMENT_TYPE_UE_APP,
  ELEMENT_TYPE_EDGE_APP,
  ELEMENT_TYPE_CLOUD_APP,

  // GPU types
  GPU_TYPE_NVIDIA,

  // NC Group Prefixes
  PREFIX_INT_DOM,
  PREFIX_INT_ZONE,
  PREFIX_INTRA_ZONE,
  PREFIX_TERM_LINK,
  PREFIX_LINK,
  PREFIX_APP,

  // Cypress test data
  CFG_ELEM_TYPE,
  CFG_ELEM_PARENT,
  CFG_ELEM_NAME,
  CFG_ELEM_IMG,
  CFG_ELEM_GROUP,
  CFG_ELEM_ENV,
  CFG_ELEM_PORT,
  CFG_ELEM_EXT_PORT,
  CFG_ELEM_PROT,
  CFG_ELEM_GPU_COUNT,
  CFG_ELEM_GPU_TYPE,
  CFG_ELEM_PLACEMENT_ID,
  CFG_ELEM_CMD,
  CFG_ELEM_ARGS,
  CFG_ELEM_EXTERNAL_CHECK,
  CFG_ELEM_MNC,
  CFG_ELEM_MCC,
  CFG_ELEM_DEFAULT_CELL_ID,
  CFG_ELEM_CELL_ID,
  CFG_ELEM_CHART_CHECK,
  CFG_ELEM_CHART_LOC,
  CFG_ELEM_CHART_GROUP,
  CFG_ELEM_CHART_ALT_VAL,
  CFG_ELEM_INGRESS_SVC_MAP,
  CFG_ELEM_EGRESS_SVC_MAP,
  CFG_BTN_NEW_ELEM,
  CFG_BTN_DEL_ELEM,
  CFG_BTN_CLONE_ELEM,

  // Layout type
  MEEP_COMPONENT_TABLE_LAYOUT
} from '../../meep-constants';

// ELEMENT VALIDATION

const SERVICE_PORT_MIN = 1;
const SERVICE_PORT_MAX = 65535;
const SERVICE_NODE_PORT_MIN = 30000;
const SERVICE_NODE_PORT_MAX = 32767;
const GPU_COUNT_MIN = 1;
const GPU_COUNT_MAX = 4;

const validateName = val => {
  if (val) {
    if (val.length > 30) {
      return 'Maximum 30 characters';
    } else if (!val.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
      return 'Lowercase alphanumeric or \'-\' or \'.\'';
    }
  }
  return null;
};

const validateFullName = val => {
  if (val) {
    if (val.length > 60) {
      return 'Maximum 60 characters';
    } else if (!val.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
      return 'Lowercase alphanumeric or \'-\' or \'.\'';
    }
  }
  return null;
};

const validateVariableName = val => {
  if (val) {
    if (val.length > 30) {
      return 'Maximum 30 characters';
    } else if (!val.match(/^(([_a-z0-9A-Z][_-a-z0-9.]*)?[_a-z0-9A-Z])+$/)) {
      return 'Alphanumeric or \'-\' or \'.\'';
    }
  }
  return null;
};

// const validateOptionalName = (val) => {
//   if (val ==='') {return null;}
//   return validateName(val);
// };

// const validateChars = (val) => {
//   /*eslint-disable */
//   if (val.match(/^.*?(?=[\^#%&$\*:<>\?/\{\|\} ]).*$/)) {
//   /*eslint-enable */
//     return 'Invalid characters';
//   }
//   return null;
// };

const notNull = val => val;
const validateNotNull = val => {
  if (!notNull(val)) {
    return 'Value is required';
  }
};

const validateNumber = val => {
  if (isNaN(val)) {
    return 'Must be a number';
  }
  return null;
};

const validateInt = val => {
  const numberError = validateNumber(val);
  if (numberError) {
    return numberError;
  }
  return val.indexOf('.') === -1 ? null : 'Must be an integer';
};

const validatePath = val => {
  /*eslint-disable */
  if (val.match(/^.*?(?=[\^#%&$\*<>\?\{\|\} ]).*$/)) {
    /*eslint-enable */
    return 'Invalid characters';
  }
  return null;
};

const validatePort = port => {
  if (port === '') {
    return null;
  }

  const notIntError = validateInt(port);
  if (notIntError) {
    return notIntError;
  }

  const p = Number(port);
  if (p !== '' && (p < SERVICE_PORT_MIN || p > SERVICE_PORT_MAX)) {
    return SERVICE_PORT_MIN + ' < port < ' + SERVICE_PORT_MAX;
  }
  return null;
};

const validateGpuCount = count => {
  if (count === '') {
    return null;
  }

  const notIntError = validateInt(count);
  if (notIntError) {
    return notIntError;
  }

  const p = Number(count);
  if (p !== '' && (p < GPU_COUNT_MIN || p > GPU_COUNT_MAX)) {
    return GPU_COUNT_MIN + ' < count < ' + GPU_COUNT_MAX;
  }
  return null;
};

const validateCellularMccMnc = val => {
  if (val) {
    if (val.length > 3) {
      return 'Maximum 3 numeric characters';
    } else if (!val.match(/^(([0-9][0-9]*)?[0-9])+$/)) {
      return 'Numeric characters only';
    }
  }
  return null;
};

const validateCellularCellId = val => {
  if (val) {
    if (val.length > 7) {
      return 'Maximum 7 characters';
    } else if (!val.match(/^(([_a-f0-9A-F][_-a-f0-9]*)?[_a-f0-9A-F])+$/)) {
      return 'Alphanumeric hex characters only';
    }
  }
  return null;
};

const validateExternalPort = port => {
  if (port === '') {
    return null;
  }

  const notIntError = validateInt(port);
  if (notIntError) {
    return notIntError;
  }

  const p = Number(port);
  if (p !== '' && (p < SERVICE_NODE_PORT_MIN || p > SERVICE_NODE_PORT_MAX)) {
    return SERVICE_NODE_PORT_MIN + ' < ext. port < ' + SERVICE_NODE_PORT_MAX;
  }
  return null;
};

const validateProtocol = protocol => {
  if (protocol === '') {
    return null;
  }

  if (protocol) {
    if (protocol !== '' && protocol !== 'TCP' && protocol !== 'UDP') {
      return 'Must be TCP or UDP';
    }
  }
  return null;
};

// Validates list of similar comma-separated entries
const validateEntries = validator => entries => {
  return _.chain(entries.split(','))
    .map(validator)
    .flatten()
    .value()
    .join(', \n');
};

const validateIngressServiceMappingEntry = entry => {
  if (entry === '') {
    return null;
  }

  const args = entry.split(':');
  if (args.length !== 4) {
    return ` ${'Ext Port:Svc Name:Port:Protocol[,Ext Port: Svc Name:Port:Protocol]'}`;
  }

  return [
    validateExternalPort(args[0]),
    validateFullName(args[1]),
    validatePort(args[2]),
    validateProtocol(args[3])
  ].filter(notNull);
};

const validateEgressServiceMappingEntry = entry => {
  if (entry === '') {
    return null;
  }

  const args = entry.split(':');
  if (args.length !== 5) {
    return ` ${'Svc Name:ME Svc Name:IP:Port:Protocol[,Svc Name:ME Svc Name:IP:Port:Protocol]'}`;
  }

  return [
    validateFullName(args[0]),
    validateFullName(args[1]),
    // validateIP(args[2]), <-- TODO
    validatePort(args[3]),
    validateProtocol(args[4])
  ].filter(notNull);
};

const validateEnvironmentVariableEntry = entry => {
  if (entry === '') {
    return null;
  }

  const parts = entry.split('=');
  if (parts.length !== 2) {
    return `${'VAR=value[,VAR=value]'}`;
  }

  return [validateVariableName(parts[0]), validateNotNull(parts[1])].filter(
    notNull
  );
};

const validateChartGroupEntry = entry => {
  if (entry === '') {
    return null;
  }

  const args = entry.split(':');
  if (args.length !== 4) {
    return ` ${'Svc instance:svc group name:port:protocol'}`;
  }

  return [
    validateFullName(args[0]),
    validateName(args[1]),
    validatePort(args[2]),
    validateProtocol(args[3])
  ]
    .filter(notNull)
    .join(',');
};

const validateIngressServiceMapping = entries =>
  validateEntries(validateIngressServiceMappingEntry)(entries);
const validateEgressServiceMapping = entries =>
  validateEntries(validateEgressServiceMappingEntry)(entries);
const validateEnvironmentVariables = entries =>
  validateEntries(validateEnvironmentVariableEntry)(entries);

const validateCommandArguments = () => null;

// COMPONENTS
const CfgTextField = props => {
  var err = props.element[props.fieldName]
    ? props.element[props.fieldName].err
    : null;
  return (
    <>
      <TextField
        outlined
        style={{ width: '100%' }}
        label={props.label}
        type={props.type}
        onChange={event => {
          var err = props.validate ? props.validate(event.target.value) : null;
          var val =
            event.target.value && props.isNumber && !err
              ? Number(event.target.value)
              : event.target.value;
          props.onUpdate(props.fieldName, val, err);
        }}
        invalid={err}
        value={
          props.element[props.fieldName]
            ? props.element[props.fieldName].val
            : ''
        }
        disabled={props.disabled}
        data-cy={props.cydata}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>{getElemFieldErr(props.element, props.fieldName)}</span>
      </TextFieldHelperText>
    </>
  );
};

const CfgTextFieldCell = props => {
  return (
    <GridCell span={props.span}>
      <CfgTextField {...props} />
    </GridCell>
  );
};

const PortProtocolGroup = ({ onUpdate, element }) => {
  return (
    <Grid>
      <CfgTextFieldCell
        span={4}
        onUpdate={onUpdate}
        element={element}
        validate={validatePort}
        isNumber={true}
        label="Port #"
        fieldName={FIELD_PORT}
        cydata={CFG_ELEM_PORT}
      />
      <CfgTextFieldCell
        span={4}
        onUpdate={onUpdate}
        element={element}
        validate={validateExternalPort}
        isNumber={true}
        label="External Port #"
        fieldName={FIELD_EXT_PORT}
        cydata={CFG_ELEM_EXT_PORT}
      />

      <GridCell span={4} style={{ paddingTop: 16 }}>
        <Select
          style={{ width: '100%' }}
          label="Protocol"
          outlined
          value={getElemFieldVal(element, FIELD_PROTOCOL)}
          options={['TCP', 'UDP']}
          onChange={event => onUpdate(FIELD_PROTOCOL, event.target.value, null)}
          data-cy={CFG_ELEM_PROT}
        />
      </GridCell>
    </Grid>
  );
};

const gpuTypes = [GPU_TYPE_NVIDIA];

const GpuGroup = ({ onUpdate, element }) => {
  var type = getElemFieldVal(element, FIELD_GPU_TYPE) || '';

  return (
    <Grid>
      <CfgTextFieldCell
        span={4}
        onUpdate={onUpdate}
        element={element}
        validate={validateGpuCount}
        isNumber={true}
        label="GPU Count"
        fieldName={FIELD_GPU_COUNT}
        cydata={CFG_ELEM_GPU_COUNT}
      />
      <GridCell span={8} style={{ paddingTop: 16 }}>
        <IDSelect
          label="GPU Type"
          span={8}
          options={gpuTypes}
          onChange={elem => onUpdate(FIELD_GPU_TYPE, elem.target.value, null)}
          value={type}
          disabled={false}
          cydata={CFG_ELEM_GPU_TYPE}
        />
      </GridCell>
    </Grid>
  );
};

const CommandGroup = ({ onUpdate, element }) => {
  return (
    <Grid>
      <CfgTextFieldCell
        span={4}
        onUpdate={onUpdate}
        element={element}
        label="Command"
        validate={validatePath}
        fieldName={FIELD_CMD}
        cydata={CFG_ELEM_CMD}
      />
      <CfgTextFieldCell
        span={8}
        onUpdate={onUpdate}
        element={element}
        label="Arguments"
        validate={validateCommandArguments}
        fieldName={FIELD_CMD_ARGS}
        cydata={CFG_ELEM_ARGS}
      />
    </Grid>
  );
};

const NCGroups = ({ prefixes, onUpdate, element }) => {
  return _.map(prefixes, p => {
    return (
      <NCGroup
        onUpdate={onUpdate}
        type={TYPE_CFG}
        element={element}
        prefix={p}
        key={p}
        layout={MEEP_COMPONENT_TABLE_LAYOUT}
      />
    );
  });
};

const ExternalFields = ({ element, onUpdate }) => {
  return (
    <>
      <CfgTextField
        onUpdate={onUpdate}
        element={element}
        label="IngressServiceMapping"
        validate={validateIngressServiceMapping}
        fieldName="ingressServiceMap"
        cydata={CFG_ELEM_INGRESS_SVC_MAP}
      />
      <CfgTextField
        onUpdate={onUpdate}
        element={element}
        label="EgressServiceMapping"
        validate={validateEgressServiceMapping}
        fieldName="egressServiceMap"
        cydata={CFG_ELEM_EGRESS_SVC_MAP}
      />
    </>
  );
};

const UserChartFields = ({ element, onUpdate }) => {
  return (
    <>
      <CfgTextField
        onUpdate={onUpdate}
        element={element}
        label="User Chart Location"
        validate={validatePath}
        fieldName={FIELD_CHART_LOC}
        cydata={CFG_ELEM_CHART_LOC}
      />
      <CfgTextField
        onUpdate={onUpdate}
        element={element}
        label="User Chart Group"
        validate={validateChartGroupEntry}
        fieldName={FIELD_CHART_GROUP}
        cydata={CFG_ELEM_CHART_GROUP}
      />
      <CfgTextField
        onUpdate={onUpdate}
        element={element}
        label="User Chart Alternate Values"
        validate={validatePath}
        fieldName={FIELD_CHART_VAL}
        cydata={CFG_ELEM_CHART_ALT_VAL}
      />
    </>
  );
};

// Display element-specific form fields
const TypeRelatedFormFields = ({ onUpdate, element }) => {
  var type = getElemFieldVal(element, FIELD_TYPE);
  var isExternal = getElemFieldVal(element, FIELD_IS_EXTERNAL);
  var chartEnabled = getElemFieldVal(element, FIELD_CHART_ENABLED);

  switch (type) {
  case ELEMENT_TYPE_SCENARIO:
    return (
      <NCGroups
        onUpdate={onUpdate}
        element={element}
        prefixes={[PREFIX_INT_DOM]}
      />
    );
  case ELEMENT_TYPE_OPERATOR:
    return (
      <>
        <NCGroups
          onUpdate={onUpdate}
          element={element}
          prefixes={[PREFIX_INT_ZONE]}
        />
      </>
    );
  case ELEMENT_TYPE_OPERATOR_CELL:
    return (
      <>
        <NCGroups
          onUpdate={onUpdate}
          element={element}
          prefixes={[PREFIX_INT_ZONE]}
        />
        <Grid>
          <CfgTextFieldCell
            span={3}
            onUpdate={onUpdate}
            element={element}
            validate={validateCellularMccMnc}
            label="MCC"
            fieldName={FIELD_MCC}
            cydata={CFG_ELEM_MCC}
          />
          <CfgTextFieldCell
            span={3}
            onUpdate={onUpdate}
            element={element}
            validate={validateCellularMccMnc}
            label="MNC"
            fieldName={FIELD_MNC}
            cydata={CFG_ELEM_MNC}
          />
          <CfgTextFieldCell
            span={6}
            onUpdate={onUpdate}
            element={element}
            validate={validateCellularCellId}
            label="Default cell Id"
            fieldName={FIELD_DEFAULT_CELL_ID}
            cydata={CFG_ELEM_DEFAULT_CELL_ID}
          />
        </Grid>
      </>
    );
  case ELEMENT_TYPE_ZONE:
    return (
      <NCGroups
        onUpdate={onUpdate}
        element={element}
        prefixes={[PREFIX_INTRA_ZONE]}
      />
    );
  case ELEMENT_TYPE_POA:
    return (
      <NCGroups
        onUpdate={onUpdate}
        element={element}
        prefixes={[PREFIX_TERM_LINK]}
      />
    );
  case ELEMENT_TYPE_POA_CELL:
    return (
      <>
        <NCGroups
          onUpdate={onUpdate}
          element={element}
          prefixes={[PREFIX_TERM_LINK]}
        />
        <CfgTextFieldCell
          onUpdate={onUpdate}
          element={element}
          validate={validateCellularCellId}
          label="Cell Id"
          fieldName={FIELD_CELL_ID}
          cydata={CFG_ELEM_CELL_ID}
        />
      </>
    );
  case ELEMENT_TYPE_UE:
  case ELEMENT_TYPE_DC:
  case ELEMENT_TYPE_EDGE:
  case ELEMENT_TYPE_FOG:
    return (
      <NCGroups
        onUpdate={onUpdate}
        element={element}
        prefixes={[PREFIX_LINK]}
      />
    );
  case ELEMENT_TYPE_UE_APP:
    return (
        <>
          <NCGroups
            onUpdate={onUpdate}
            element={element}
            prefixes={[PREFIX_APP]}
          />

          <Checkbox
            checked={isExternal}
            onChange={e => onUpdate(FIELD_IS_EXTERNAL, e.target.checked, null)}
            data-cy={CFG_ELEM_EXTERNAL_CHECK}
          >
            External App
          </Checkbox>

          {isExternal ? (
            <>
              <ExternalFields onUpdate={onUpdate} element={element} />
              <CfgTextField
                onUpdate={onUpdate}
                element={element}
                label="Placement Identifier"
                fieldName={FIELD_PLACEMENT_ID}
                cydata={CFG_ELEM_PLACEMENT_ID}
              />
            </>
          ) : (
              <>
                <Checkbox
                  checked={chartEnabled}
                  onChange={e =>
                    onUpdate(FIELD_CHART_ENABLED, e.target.checked, null)
                  }
                  data-cy={CFG_ELEM_CHART_CHECK}
                >
                  User-Defined Chart
                </Checkbox>

                {chartEnabled ? (
                  <UserChartFields onUpdate={onUpdate} element={element} />
                ) : (
                    <>
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Container Image Name"
                        validate={validatePath}
                        fieldName={FIELD_IMAGE}
                        cydata={CFG_ELEM_IMG}
                      />
                      <GpuGroup onUpdate={onUpdate} element={element} />
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Environment variables"
                        validate={validateEnvironmentVariables}
                        fieldName={FIELD_ENV_VAR}
                        cydata={CFG_ELEM_ENV}
                      />
                      <CommandGroup onUpdate={onUpdate} element={element} />
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Placement Identifier"
                        fieldName={FIELD_PLACEMENT_ID}
                        cydata={CFG_ELEM_PLACEMENT_ID}
                      />
                    </>
                )}
              </>
          )}
        </>
    );
  case ELEMENT_TYPE_CLOUD_APP:
  case ELEMENT_TYPE_MECSVC:
    return (
        <>
          <NCGroups
            onUpdate={onUpdate}
            element={element}
            prefixes={[PREFIX_APP]}
          />

          <Checkbox
            checked={isExternal}
            onChange={e => onUpdate(FIELD_IS_EXTERNAL, e.target.checked, null)}
            data-cy={CFG_ELEM_EXTERNAL_CHECK}
          >
            External App
          </Checkbox>

          {isExternal ? (
            <>
              <ExternalFields onUpdate={onUpdate} element={element} />
              <CfgTextField
                onUpdate={onUpdate}
                element={element}
                label="Placement Identifier"
                fieldName={FIELD_PLACEMENT_ID}
                cydata={CFG_ELEM_PLACEMENT_ID}
              />
            </>
          ) : (
              <>
                <Checkbox
                  checked={chartEnabled}
                  onChange={e =>
                    onUpdate(FIELD_CHART_ENABLED, e.target.checked, null)
                  }
                  data-cy={CFG_ELEM_CHART_CHECK}
                >
                  User-Defined Chart
                </Checkbox>

                {chartEnabled ? (
                  <UserChartFields onUpdate={onUpdate} element={element} />
                ) : (
                    <>
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Container Image Name"
                        validate={validatePath}
                        fieldName={FIELD_IMAGE}
                        cydata={CFG_ELEM_IMG}
                      />
                      <PortProtocolGroup onUpdate={onUpdate} element={element} />
                      <GpuGroup onUpdate={onUpdate} element={element} />
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Environment variables"
                        validate={validateEnvironmentVariables}
                        fieldName={FIELD_ENV_VAR}
                        cydata={CFG_ELEM_ENV}
                      />
                      <CommandGroup onUpdate={onUpdate} element={element} />
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Placement Identifier"
                        fieldName={FIELD_PLACEMENT_ID}
                        cydata={CFG_ELEM_PLACEMENT_ID}
                      />
                    </>
                )}
              </>
          )}
        </>
    );
  case ELEMENT_TYPE_EDGE_APP:
    return (
        <>
          <NCGroups
            onUpdate={onUpdate}
            element={element}
            prefixes={[PREFIX_APP]}
          />

          <Checkbox
            checked={isExternal}
            onChange={e => onUpdate(FIELD_IS_EXTERNAL, e.target.checked, null)}
            data-cy={CFG_ELEM_EXTERNAL_CHECK}
          >
            External App
          </Checkbox>

          {isExternal ? (
            <>
              <ExternalFields onUpdate={onUpdate} element={element} />
              <CfgTextField
                onUpdate={onUpdate}
                element={element}
                label="Placement Identifier"
                fieldName={FIELD_PLACEMENT_ID}
                cydata={CFG_ELEM_PLACEMENT_ID}
              />
            </>
          ) : (
              <>
                <Checkbox
                  checked={chartEnabled}
                  onChange={e =>
                    onUpdate(FIELD_CHART_ENABLED, e.target.checked, null)
                  }
                  data-cy={CFG_ELEM_CHART_CHECK}
                >
                  User-Defined Chart
                </Checkbox>

                {chartEnabled ? (
                  <UserChartFields onUpdate={onUpdate} element={element} />
                ) : (
                    <>
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Container Image Name"
                        validate={validatePath}
                        fieldName={FIELD_IMAGE}
                        cydata={CFG_ELEM_IMG}
                      />
                      <PortProtocolGroup onUpdate={onUpdate} element={element} />
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Group Service Name"
                        validate={validateName}
                        fieldName={FIELD_GROUP}
                        cydata={CFG_ELEM_GROUP}
                      />
                      <GpuGroup onUpdate={onUpdate} element={element} />
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Environment variables"
                        validate={validateEnvironmentVariables}
                        fieldName={FIELD_ENV_VAR}
                        cydata={CFG_ELEM_ENV}
                      />
                      <CommandGroup onUpdate={onUpdate} element={element} />
                      <CfgTextField
                        onUpdate={onUpdate}
                        element={element}
                        label="Placement Identifier"
                        fieldName={FIELD_PLACEMENT_ID}
                        cydata={CFG_ELEM_PLACEMENT_ID}
                      />
                    </>
                )}
              </>
          )}
        </>
    );

  default:
    return null;
  }
};

const elementTypes = [
  {
    label: 'Logical Domain',
    options: [ELEMENT_TYPE_OPERATOR_GENERIC, ELEMENT_TYPE_OPERATOR_CELL]
  },
  {
    label: 'Logical Zone',
    options: [ELEMENT_TYPE_ZONE]
  },
  {
    label: 'Network Location',
    options: [ELEMENT_TYPE_POA_GENERIC, ELEMENT_TYPE_POA_CELL]
  },
  {
    label: 'Physical Location',
    options: [
      ELEMENT_TYPE_UE,
      ELEMENT_TYPE_FOG,
      ELEMENT_TYPE_EDGE,
      ELEMENT_TYPE_DC
      // ELEMENT_TYPE_CN
    ]
  },
  {
    label: 'Process',
    options: [
      ELEMENT_TYPE_UE_APP,
      // ELEMENT_TYPE_MECSVC,
      ELEMENT_TYPE_EDGE_APP,
      ELEMENT_TYPE_CLOUD_APP
    ]
  }
];

var parentTypes = {};
parentTypes[ELEMENT_TYPE_SCENARIO] = null;
parentTypes[ELEMENT_TYPE_OPERATOR] = [ELEMENT_TYPE_SCENARIO];
parentTypes[ELEMENT_TYPE_OPERATOR_CELL] = [ELEMENT_TYPE_SCENARIO];
parentTypes[ELEMENT_TYPE_EDGE] = [ELEMENT_TYPE_ZONE];
parentTypes[ELEMENT_TYPE_ZONE] = [ELEMENT_TYPE_OPERATOR, ELEMENT_TYPE_OPERATOR_CELL];
parentTypes[ELEMENT_TYPE_POA] = [ELEMENT_TYPE_ZONE];
parentTypes[ELEMENT_TYPE_POA_CELL] = [ELEMENT_TYPE_ZONE];
parentTypes[ELEMENT_TYPE_CN] = [ELEMENT_TYPE_ZONE];
parentTypes[ELEMENT_TYPE_FOG] = [ELEMENT_TYPE_POA, ELEMENT_TYPE_POA_CELL];
parentTypes[ELEMENT_TYPE_UE] = [ELEMENT_TYPE_POA, ELEMENT_TYPE_POA_CELL];
parentTypes[ELEMENT_TYPE_DC] = [ELEMENT_TYPE_SCENARIO];
parentTypes[ELEMENT_TYPE_UE_APP] = [ELEMENT_TYPE_UE];
parentTypes[ELEMENT_TYPE_MECSVC] = [
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_CN
];
parentTypes[ELEMENT_TYPE_EDGE_APP] = [ELEMENT_TYPE_FOG, ELEMENT_TYPE_EDGE];
parentTypes[ELEMENT_TYPE_CLOUD_APP] = [ELEMENT_TYPE_DC];

const getParentTypes = type => {
  return parentTypes[type];
};

const buttonStyles = {
  marginRight: 10,
  marginBottom: 5
};

const ElementCfgButtons = ({
  configuredElement,
  configMode,
  onNewElement,
  onDeleteElement,
  onCloneElement
}) => {
  const canCreateNewElement = () => {
    return !configuredElement;
  };

  const canDeleteOrCloneElement = () => {
    return configuredElement && configMode === CFG_ELEM_MODE_EDIT;
  };

  return (
    <>
      <Button
        outlined
        data-cy={CFG_BTN_NEW_ELEM}
        style={buttonStyles}
        onClick={() => onNewElement()}
        disabled={!canCreateNewElement()}
      >
        NEW
      </Button>

      <Button
        outlined
        data-cy={CFG_BTN_DEL_ELEM}
        style={buttonStyles}
        onClick={() => onDeleteElement()}
        disabled={!canDeleteOrCloneElement()}
      >
        DELETE
      </Button>

      <Button
        outlined
        data-cy={CFG_BTN_CLONE_ELEM}
        style={buttonStyles}
        onClick={() => onCloneElement()}
        disabled={!canDeleteOrCloneElement()}
      >
        CLONE
      </Button>
    </>
  );
};

const getSuggestedName = ( type, elements ) => {
  var suggestedPrefix = '';
  switch(type) {
  case ELEMENT_TYPE_UE_APP:
    suggestedPrefix = 'ue-app';
    break;
  case ELEMENT_TYPE_EDGE_APP:
    suggestedPrefix = 'edge-app';
    break;
  case ELEMENT_TYPE_CLOUD_APP:
    suggestedPrefix = 'cloud-app';
    break;
  case ELEMENT_TYPE_DC:
    suggestedPrefix = 'cloud';
    break;
  case ELEMENT_TYPE_POA_CELL:
    suggestedPrefix = 'poa-cell';
    break;
  case ELEMENT_TYPE_OPERATOR_CELL:
    suggestedPrefix = 'operator-cell';
    break;
  default:
    suggestedPrefix = type.toLowerCase();
  }

  return createUniqueName(elements, suggestedPrefix);
};

const getElementTypeOverride = (type) => {
  var typeOverride = '';
  switch(type) {
  case ELEMENT_TYPE_POA:
    typeOverride = ELEMENT_TYPE_POA_GENERIC;
    break;
  case ELEMENT_TYPE_OPERATOR:
    typeOverride = ELEMENT_TYPE_OPERATOR_GENERIC;
    break;
  default:
    typeOverride = type;
  }
  return typeOverride;
};

const getElementTypeOverrideBack = (typeOverride) => {
  var type = '';
  switch(typeOverride) {
  case ELEMENT_TYPE_POA_GENERIC:
    type = ELEMENT_TYPE_POA;
    break;
  case ELEMENT_TYPE_OPERATOR_GENERIC:
    type = ELEMENT_TYPE_OPERATOR;
    break;
  default:
    type = typeOverride;
  }
  return type;
};

const HeaderGroup = ({ element, onTypeChange, onUpdate, typeDisabled, parentDisabled, nameDisabled }) => {
  var type = getElemFieldVal(element, FIELD_TYPE) || '';
  var parent = getElemFieldVal(element, FIELD_PARENT) || '';
  var parentElements = element.parentElements || [parent];

  var typeOverride = getElementTypeOverride(type);

  return (
    <>
      <Grid style={{ marginTop: 20 }}>
        {type !== 'SCENARIO' && (
          <IDSelect
            label="Element Type"
            span={6}
            options={elementTypes}
            onChange={elem => onTypeChange(elem.target.value)}
            value={typeOverride}
            disabled={typeDisabled}
            cydata={CFG_ELEM_TYPE}
          />
        )}
        {type && type !== 'SCENARIO' && (
          <IDSelect
            label="Parent Node"
            span={6}
            options={parentElements}
            onChange={elem => onUpdate(FIELD_PARENT, elem.target.value, null)}
            value={parent}
            disabled={parentDisabled}
            cydata={CFG_ELEM_PARENT}
          />
        )}
      </Grid>
      <Grid>
        <CfgTextFieldCell
          span={12}
          onUpdate={onUpdate}
          element={element}
          validate={validateName}
          label="Unique Element Name"
          fieldName={FIELD_NAME}
          disabled={nameDisabled}
          cydata={CFG_ELEM_NAME}
        />
      </Grid>
    </>
  );
};

export class CfgNetworkElementContainer extends Component {
  constructor(props) {
    super(props);
  }

  // Element update handler
  onUpdateElement(name, val, err) {
    var updatedElem = updateObject({}, this.props.configuredElement);
    setElemFieldVal(updatedElem, name, val);
    setElemFieldErr(updatedElem, name, err);

    this.props.cfgElemUpdate(updatedElem);
  }

  // Element clone handler
  onCloneElement(newName) {
    var clonedElem = updateObject({}, this.props.configuredElement);
    setElemFieldVal(clonedElem, FIELD_NAME, newName);
    setElemFieldVal(clonedElem, FIELD_PARENT, null);
    var elementType = getElemFieldVal(clonedElem, FIELD_TYPE);
    clonedElem.parentElements = this.elementsOfType(getParentTypes(elementType));

    this.props.cfgElemClone(clonedElem);
  }

  // Retrieve names of elements with matching type
  elementsOfType(types) {
    return _.chain(this.props.tableData)
      .filter(e => {
        var elemType = getElemFieldVal(e, FIELD_TYPE);
        return _.includes(types, elemType);
      })
      .map(e => {
        return getElemFieldVal(e, FIELD_NAME);
      })
      .value();
  }

  // Element configuration type change handler
  onElementTypeChange(elementType) {
    var elem = updateObject({}, this.props.configuredElement);

    //override the frontend terminology
    var elementTypeOverride = getElementTypeOverrideBack(elementType);
    setElemFieldVal(elem, FIELD_TYPE, elementTypeOverride);
    setElemFieldVal(elem, FIELD_PARENT, null);

    elem.parentElements = this.elementsOfType(getParentTypes(elementTypeOverride));

    if (this.props.configMode !== CFG_ELEM_MODE_CLONE) {
      setElemFieldVal(elem, FIELD_NAME, getSuggestedName(elementTypeOverride, this.props.tableData));
    }
    this.props.cfgElemUpdate(elem);
  }

  render() {
    const element = this.props.configuredElement;
    return (
      <div className="cfg-network-element-div" style={styles.outer}>
        <Grid>
          <GridCell span={12}>
            <div style={styles.block}>
              <Typography use="headline6">Element Configuration</Typography>
            </div>
          </GridCell>
          <GridCell span={12}>
            <GridInner align={'left'}>
              <GridCell span={12}>
                <ElementCfgButtons
                  configuredElement={element}
                  configMode={this.props.configMode}
                  onNewElement={this.props.onNewElement}
                  onDeleteElement={() => {
                    this.props.onDeleteElement(element);
                  }}
                  onCloneElement={() => {
                    this.onCloneElement(createUniqueName(this.props.tableData, getElemFieldVal(element, FIELD_NAME) + '-copy'));
                  }}
                />
              </GridCell>
            </GridInner>
          </GridCell>
        </Grid>

        {element && (
          <>
            <HeaderGroup
              element={element}
              onTypeChange={type => {
                this.onElementTypeChange(type);
              }}
              onUpdate={(name, val, err) => {
                this.onUpdateElement(name, val, err);
              }}
              typeDisabled={this.props.configMode === CFG_ELEM_MODE_CLONE || this.props.configMode === CFG_ELEM_MODE_EDIT}
              parentDisabled={this.props.configMode === CFG_ELEM_MODE_EDIT}
              nameDisabled={getElemFieldVal(element, FIELD_TYPE) === ELEMENT_TYPE_SCENARIO && this.props.configMode !== CFG_ELEM_MODE_NEW}
            />

            <TypeRelatedFormFields
              element={element}
              onUpdate={(name, val, err) => {
                this.onUpdateElement(name, val, err);
              }}
            />

            <div
              id="new-element-error-message"
              className="idcc-margin-top mdc-typography--body1"
            >
              {this.props.errorMessage}
            </div>

            <CancelApplyPair
              saveDisabled={(this.props.isModified === false) ? true : false}
              onCancel={this.props.onCancelElement}
              onApply={() => {
                (this.props.configMode === CFG_ELEM_MODE_CLONE) ? this.props.onApplyCloneElement(element) : this.props.onSaveElement(element);
              }}

            />
          </>
        )}
      </div>
    );
  }
}

const styles = {
  outer: {
    padding: 10,
    height: '100%'
  },
  block: {
    marginBottom: 10
  },
  field: {
    marginBottom: 10
  },
  button: {
    color: 'white'
  },
  select: {
    width: '100%'
  }
};

const mapStateToProps = state => {
  return {
    tableData: state.cfg.table.entries,
    configuredElement: state.cfg.elementConfiguration.configuredElement,
    configMode: state.cfg.elementConfiguration.configurationMode,
    isModified: state.cfg.elementConfiguration.isModified,
    errorMessage: state.cfg.elementConfiguration.errorMessage
  };
};

const mapDispatchToProps = dispatch => {
  return {
    cfgElemUpdate: element => dispatch(cfgElemUpdate(element)),
    cfgElemClone: element => dispatch(cfgElemClone(element))
  };
};

const ConnectedCfgNetworkElementContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(CfgNetworkElementContainer);

export default ConnectedCfgNetworkElementContainer;
