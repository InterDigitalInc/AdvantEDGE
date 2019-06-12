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

const data = {  
  'children':[  
    {  
      'name':'boss1',
      'children':[  
        {  
          'name':'mister_a',
          'colname':'level3'
        },
        {  
          'name':'mister_b',
          'colname':'level3'
        },
        {  
          'name':'mister_c',
          'colname':'level3'
        },
        {  
          'name':'mister_d',
          'colname':'level3'
        }
      ],
      'colname':'level2'
    },
    {  
      'name':'boss2',
      'children':[  
        {  
          'name':'mister_e',
          'colname':'level3'
        },
        {  
          'name':'mister_f',
          'colname':'level3'
        },
        {  
          'name':'mister_g',
          'colname':'level3'
        },
        {  
          'name':'mister_h',
          'colname':'level3'
        }
      ],
      'colname':'level2'
    }
  ],
  'name':'CEO'
};

const computeDisplayedScenario = scenario => {
  // TODO: replaced with real computed scenario
  let root = scenario.deployment;
  root.iconName = 'cloud-black.svg';

  visitScenario(setImageForChildren)(root);
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
      type: 'COMMON'
    };
  }

  return null;
};

const getImageForType = type => {
  switch(type) {
  case 'ZONE':
    return 'tower-02-idcc.svg';
  case 'EDGE':
    return 'edge-idcc.svg';
  case 'PUBLIC':
    return '';
  case 'OPERATOR':
    return 'fog-idcc.svg';
  case 'COMMON':
    return 'Gear-01-idcc.svg';
  default:
    return 'Gear-01-idcc.svg';
  }
};

const setImageForChildren = node => {
  const info = getChildrenInfo(node);
  if (!info) {
    return;
  }

  const {fieldName, type} = info;
  if (fieldName && type) {
    _.each(node[fieldName], c => {
      c.iconName = getImageForType(type);
      // console.log(`iconName for type ${type} is ${c.iconName}`);
    });
  }
};

const visitScenario = f => node => {
  // console.log(`Visiting scenario for node with type ${node.type}`);
  const ff = f;
  ff(node);
  const info = getChildrenInfo(node);

  if (!info) {
    return;
  }
  
  const {fieldName, type} = info;
  if (fieldName && type) {
    _.each(node[fieldName], c => {
      visitScenario(ff)(c);
    });
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