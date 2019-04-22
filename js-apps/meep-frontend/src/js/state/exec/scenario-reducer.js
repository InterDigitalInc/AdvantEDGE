/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
 import { updateObject } from '../../util/update';

// EXEC_CHANGE_SCENARIO
const EXEC_CHANGE_SCENARIO = 'EXEC_CHANGE_SCENARIO';
function execChangeScenario(scenario) {
  return {
    type: EXEC_CHANGE_SCENARIO,
    payload: scenario
  };
}

// EXEC_CHANGE_SCENARIO_NAME
const EXEC_CHANGE_SCENARIO_NAME = 'CFG_CHANGE_SCENARIO_NAME';
function execChangeScenarioName(name) {
  return {
    type: EXEC_CHANGE_SCENARIO_NAME,
    payload: name
  };
}

export {
  // Action creators
  execChangeScenario,
  execChangeScenarioName
};

const initialState = {
  name: 'none',
  deployment: {
    domains: [],
    interDomainLatency: 50,
    interDomainLatencyVariation: 10,
    interDomainPacketLoss: 0,
    interDomainThroughput: 1000000
  }
};

export function scenarioReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SCENARIO_NAME:
    return updateObject(state, {name: action.payload});
  case EXEC_CHANGE_SCENARIO:
    return action.payload;
  default:
    return state;
  }
}
