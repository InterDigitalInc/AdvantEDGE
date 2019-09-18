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

import { updateObject } from '../../util/object-util';

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
