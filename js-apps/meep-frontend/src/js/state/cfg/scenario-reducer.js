/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
 import { updateObject } from '../../util/update';

// CHANGE_SCENARIO
const CFG_CHANGE_SCENARIO = 'CFG_CHANGE_SCENARIO';
function cfgChangeScenario(scenario) {
  return {
    type: CFG_CHANGE_SCENARIO,
    payload: scenario
  };
}

// CHANGE_SCENARIO_NAME
const CFG_CHANGE_SCENARIO_NAME = 'CFG_CHANGE_SCENARIO_NAME';
function cfgChangeScenarioName(name) {
  return {
    type: CFG_CHANGE_SCENARIO_NAME,
    payload: name
  };
}

export { cfgChangeScenario, cfgChangeScenarioName };

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
  case CFG_CHANGE_SCENARIO_NAME:
    return updateObject(state, {name: action.payload});
  case CFG_CHANGE_SCENARIO:
    return action.payload;
  default:
    return state;
  }
}
