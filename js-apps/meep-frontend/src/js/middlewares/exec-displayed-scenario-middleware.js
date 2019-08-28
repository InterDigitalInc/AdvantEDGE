import _ from 'lodash';

import {
  EXEC_CHANGE_SCENARIO,
  EXEC_CHANGE_DISPLAYED_SCENARIO,
  typeReducer
} from '../state/exec';

import {
  TYPE_SCENARIO,
  TYPE_DOMAIN,
  TYPE_ZONE,
  TYPE_NET_LOC,
  TYPE_PHY_LOC,
  TYPE_PROCESS 
} from '../meep-constants';
import { getScenarioSpecificImage } from '../util/scenario-utils';

const softwareIcon = 'software-icon.svg';

const computeDisplayedScenario = scenario => {
  // TODO: replaced with real computed scenario
  let root = scenario.deployment;
  root.iconName = 'cloud-black.svg';

  visitNodes(setImage)(root);
  return root;
};

const getChildrenInfo = node => {
  if (node.domains) {
    return {
      fieldName: 'domains',
      type: 'DOMAIN'
    };
  }
  if (node.zones) {
    return {
      fieldName: 'zones',
      type: 'ZONE'
    };
  }
  if (node.networkLocations) {
    return {
      fieldName: 'networkLocations',
      type: 'EDGE'
    };
  }
  if (node.physicalLocations) {
    return {
      fieldName: 'physicalLocations',
      type: 'EDGE'
    };
  }
  if (node.processes) {
    return {
      fieldName: 'processes',
      type: node.processes[0].type
    };
  }

  return null;
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

const getChildrenType = node => {
  let fieldName = getChildrenFieldName(node);

  if (fieldName && node[fieldName].length) {
    return node[fieldName][0].type;
  }

  return '';
};

const getImageForType = type => {
  switch(type) {
  case undefined:
    return 'cloud-black.svg';
  case 'DEFAULT':
    return 'tower-02-idcc.svg';
  case 'ZONE':
    return 'tower-02-idcc.svg';
  case 'EDGE':
    return 'edge-idcc.svg';
  case 'PUBLIC':
    return 'cloud-outline-black.svg';
  case 'OPERATOR':
    return 'fog-idcc.svg';
  case 'COMMON':
    return 'tower-02-idcc.svg';
  case 'UE-APP':
    return 'drone-blue.svg';
  case 'UE':
    return 'phone.svg';
  case 'EDGE-APP':
    return 'drone-black.svg';
  case 'CLOUD-APP':
    return 'drone-blue.svg';
  case 'POA':
    return 'switch-blue.svg';
  case 'FOG':
    return 'fog-idcc.svg';
  case 'DC':
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

export {
  execDisplayedScenarioMiddleware
};