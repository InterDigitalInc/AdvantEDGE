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
import { NO_SCENARIO_NAME } from '../../meep-constants';

// EXEC_CHANGE_SCENARIO
export const EXEC_CHANGE_SCENARIO = 'EXEC_CHANGE_SCENARIO';
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
  name: NO_SCENARIO_NAME
};

export function scenarioReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SCENARIO_NAME:
    return updateObject(state, { name: action.payload });
  case EXEC_CHANGE_SCENARIO:
    return action.payload;
  default:
    return state;
  }
}
