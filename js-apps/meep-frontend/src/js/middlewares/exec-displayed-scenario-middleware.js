/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

import {
  EXEC_CHANGE_SCENARIO,
  EXEC_CHANGE_DISPLAYED_SCENARIO
} from '../state/exec';

import {
  DOMAIN_TYPE_STR,
  DOMAIN_CELL_TYPE_STR,
  PUBLIC_DOMAIN_TYPE_STR,
  ZONE_TYPE_STR,
  COMMON_ZONE_TYPE_STR,
  // NL_TYPE_STR,
  POA_TYPE_STR,
  // POA_CELL_TYPE_STR,
  DEFAULT_NL_TYPE_STR,
  UE_TYPE_STR,
  FOG_TYPE_STR,
  EDGE_TYPE_STR,
  // CN_TYPE_STR,
  DC_TYPE_STR,
  // MEC_SVC_TYPE_STR,
  UE_APP_TYPE_STR,
  EDGE_APP_TYPE_STR,
  CLOUD_APP_TYPE_STR
} from '../meep-constants';

const computeDisplayedScenario = scenario => {
  // TODO: replaced with real computed scenario
  let root = scenario.deployment;
  root.iconName = 'cloud-black.svg';

  visitNodes(setImage)(root);
  return root;
};

const getChildrenFieldName = node => {
  let fieldName = null;
  if (node.domains) {
    fieldName = 'domains';
  }
  if (node.zones) {
    fieldName = 'zones';
  }
  if (node.networkLocations) {
    fieldName = 'networkLocations';
  }
  if (node.physicalLocations) {
    fieldName = 'physicalLocations';
  }
  if (node.processes) {
    fieldName = 'processes';
  }

  return fieldName;
};

const getImageForType = type => {
  switch (type) {
  case undefined:
    return 'cloud-black.svg';
  case DEFAULT_NL_TYPE_STR:
    return 'tower-02-idcc.svg';
  case ZONE_TYPE_STR:
    return 'tower-02-idcc.svg';
  case EDGE_TYPE_STR:
    return 'edge-idcc.svg';
  case PUBLIC_DOMAIN_TYPE_STR:
    return 'cloud-outline-black.svg';
  case DOMAIN_TYPE_STR:
    return 'fog-idcc.svg';
  case DOMAIN_CELL_TYPE_STR:
    return 'fog-idcc.svg';
  case COMMON_ZONE_TYPE_STR:
    return 'tower-02-idcc.svg';
  case UE_APP_TYPE_STR:
    return 'drone-blue.svg';
  case UE_TYPE_STR:
    return 'phone.svg';
  case EDGE_APP_TYPE_STR:
    return 'drone-black.svg';
  case CLOUD_APP_TYPE_STR:
    return 'drone-blue.svg';
  case POA_TYPE_STR:
    return 'switch-blue.svg';
  case FOG_TYPE_STR:
    return 'fog-idcc.svg';
  case DC_TYPE_STR:
    return 'cloud-outline-black.svg';
  default:
    return 'Gear-01-idcc.svg';
  }
};

const setImage = node => {
  const iconName = getImageForType(node.type);
  node.iconName = iconName;
};

const visitNodes = f => node => {
  f(node);
  // console.log('visitingNode ' + node.name + ' of type: ' + node.type);
  const childrenFieldName = getChildrenFieldName(node);
  if (node[childrenFieldName]) {
    _.each(node[childrenFieldName], c => visitNodes(f)(c));
  }
};

const execDisplayedScenarioMiddleware = store => next => action => {
  if (action.type === EXEC_CHANGE_SCENARIO) {
    const displayedScenario = computeDisplayedScenario(action.payload);
    store.dispatch({
      type: EXEC_CHANGE_DISPLAYED_SCENARIO,
      payload: displayedScenario
    });
  }

  next(action);
};

export { execDisplayedScenarioMiddleware };
