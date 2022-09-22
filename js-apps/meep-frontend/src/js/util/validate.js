/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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
import L from 'leaflet';

import {
  // Connectivity model
  CONNECTIVITY_MODEL_OPEN,
  CONNECTIVITY_MODEL_PDU,

  // Layout type
  GEO_EOP_MODE_LOOP,
  GEO_EOP_MODE_REVERSE
} from '@/js/meep-constants';

export const SERVICE_PORT_MIN = 1;
export const SERVICE_PORT_MAX = 65535;
export const SERVICE_NODE_PORT_MIN = 30000;
export const SERVICE_NODE_PORT_MAX = 32767;
export const GPU_COUNT_MIN = 1;
export const GPU_COUNT_MAX = 4;
export const CONNECTIVITY_MODELS = [CONNECTIVITY_MODEL_OPEN, CONNECTIVITY_MODEL_PDU];
export const EOP_MODES = [GEO_EOP_MODE_LOOP, GEO_EOP_MODE_REVERSE];
export const D2D_RADIUS_MIN = 1;
export const D2D_RADIUS_MAX = 10000;


export const validateName = val => {
  if (val) {
    if (val.length > 30) {
      return 'Maximum 30 characters';
    } else if (!val.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
      return 'Lowercase alphanumeric or \'-\' or \'.\'';
    }
  }
  return null;
};

export const validateFullName = val => {
  if (val) {
    if (val.length > 60) {
      return 'Maximum 60 characters';
    } else if (!val.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
      return 'Lowercase alphanumeric or \'-\' or \'.\'';
    }
  }
  return null;
};

export const validateVariableName = val => {
  if (val) {
    if (val.length > 30) {
      return 'Maximum 30 characters';
    } else if (!val.match(/^(([_a-z0-9A-Z][_-a-z0-9A-Z.]*)?[_a-z0-9A-Z])+$/)) {
      return 'Alphanumeric or \'_\' or \'.\'';
    }
  }
  return null;
};

export const validateDnn = val => {
  if (val) {
    if (val.length > 50) {
      return 'Maximum 50 characters';
    } else if (!val.match(/^(([a-z0-9A-Z][-a-z0-9A-Z.]*)?[a-z0-9A-Z])+$/)) {
      return 'Alphanumeric or \'-\' or \'.\'';
    }
  }
  return null;
};

export const validateEcsp = val => {
  if (val) {
    if (val.length > 50) {
      return 'Maximum 50 characters';
    } else if (!val.match(/^(([a-z0-9A-Z][ a-z0-9A-Z]*)?[a-z0-9A-Z])+$/)) {
      return 'Alphanumeric or \' \'';
    }
  }
  return null;
};

// export const validateOptionalName = (val) => {
//   if (val ==='') {return null;}
//   return validateName(val);
// };

// export const validateChars = (val) => {
//   /*eslint-disable */
//   if (val.match(/^.*?(?=[\^#%&$\*:<>\?/\{\|\} ]).*$/)) {
//   /*eslint-enable */
//     return 'Invalid characters';
//   }
//   return null;
// };

export const notNull = val => val;
export const validateNotNull = val => {
  if (!notNull(val)) {
    return 'Value is required';
  }
};

export const validateNumber = val => {
  if (isNaN(val)) {
    return 'Must be a number';
  }
  return null;
};

export const validateInt = val => {
  const numberError = validateNumber(val);
  if (numberError) {
    return numberError;
  }
  return val.indexOf('.') === -1 ? null : 'Must be an integer';
};

export const validatePositiveInt = val => {
  const intError = validateInt(val);
  if (intError) {
    return intError;
  }
  return val >= 0 ? null : 'Must be a positive integer';
};

export const validateD2DRadius = val => {
  const intError = validateInt(val);
  if (intError) {
    return intError;
  }
  if (val >= 0) {
    const p = Number(val);
    if (p !== '' && (p < D2D_RADIUS_MIN || p > D2D_RADIUS_MAX)) {
      return 'Out of range(' + D2D_RADIUS_MIN + '-' + D2D_RADIUS_MAX + ')';
    }
  }else {
    return 'Must be a positive integer';
  }
};

export const validatePath = val => {
  /*eslint-disable */
  if (val.match(/^.*?(?=[\^#%&$\*<>\?\{\|\} ]).*$/)) {
    /*eslint-enable */
    return 'Invalid characters';
  }
  return null;
};

export const validatePositiveFloat = val => {
  const floatError = validateNumber(val);
  if (floatError) {
    return floatError;
  }
  return val >= 0 ? null : 'Must be a positive float';
};

export const validateCpuValue = count => {
  if (count === '') {
    return null;
  }

  const notPosFloatError = validatePositiveFloat(count);
  if (notPosFloatError) {
    return notPosFloatError;
  }

  const p = Number(count);
  if (p !== '' && p === 0) {
    return 'Must be a float greater than 0';
  }
  return null;
};

export const validateMemoryValue = count => {
  if (count === '') {
    return null;
  }

  const notPosIntError = validatePositiveInt(count);
  if (notPosIntError) {
    return notPosIntError;
  }

  const p = Number(count);
  if (p !== '' && p === 0) {
    return 'Must be an integer greater than 0';
  }
  return null;
};

export const validatePort = port => {
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

export const validateGpuCount = count => {
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

export const validateWirelessType = val => {
  if (val) {
    if (!val.match(/^((,\s*)?(wifi|5g|4g|other))+$/)) {
      return 'Comma-separated values: wifi|5g|4g|other';
    }
  }
  return null;
};

export const validateCellularMccMnc = val => {
  if (val) {
    if (val.length > 3) {
      return 'Maximum 3 numeric characters';
    } else if (!val.match(/^(([0-9][0-9]*)?[0-9])+$/)) {
      return 'Numeric characters only';
    }
  }
  return null;
};

export const validateCellularCellId = val => {
  if (val) {
    if (val.length > 7) {
      return 'Maximum 7 characters';
    } else if (!val.match(/^(([_a-f0-9A-F][_-a-f0-9]*)?[_a-f0-9A-F])+$/)) {
      return 'Alphanumeric hex characters only';
    }
  }
  return null;
};

export const validateCellularNrCellId = val => {
  if (val) {
    if (val.length > 9) {
      return 'Maximum 9 characters';
    } else if (!val.match(/^(([_a-f0-9A-F][_-a-f0-9]*)?[_a-f0-9A-F])+$/)) {
      return 'Alphanumeric hex characters only';
    }
  }
  return null;
};

export const validateMacAddress = val => {
  if (val) {
    if (val.length > 12) {
      return 'Maximum 12 characters';
    } else if (!val.match(/^(([_a-f0-9A-F][_-a-f0-9]*)?[_a-f0-9A-F])+$/)) {
      return 'Alphanumeric hex characters only';
    }
  }
  return null;
};

export const validateLocation = val => {
  if (val) {
    try {
      L.GeoJSON.coordsToLatLng(JSON.parse(val));
    } catch(e) {
      return '[longitude,latitude]';
    }
  }
  return null;
};

export const validateGeoPath = val => {
  if (val) {
    // TODO -- Validate location format
    try {
      L.GeoJSON.coordsToLatLngs(JSON.parse(val),0);
    } catch(e) {
      return '[[longitude,latitude],...]';
    }
  }
  return null;
};

export const validateExternalPort = port => {
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

export const validateProtocol = protocol => {
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

export const validateColor = val => {
  if (val === '') {
    return null;
  }
  if (!val.match(/^#[0-9A-Fa-f]{6}$/)) {
    return 'Invalid hex format';
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

export const validateIngressServiceMappingEntry = entry => {
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

export const validateEgressServiceMappingEntry = entry => {
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

export const validateEnvironmentVariableEntry = entry => {
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

export const validateChartGroupEntry = entry => {
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

export const validateIngressServiceMapping = entries =>
  validateEntries(validateIngressServiceMappingEntry)(entries);

export const validateEgressServiceMapping = entries =>
  validateEntries(validateEgressServiceMappingEntry)(entries);

export const validateEnvironmentVariables = entries =>
  validateEntries(validateEnvironmentVariableEntry)(entries);

export const validateCommandArguments = () => null;
