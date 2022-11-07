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

import { updateObject } from '../../util/object-util';
import { createNewScenario } from '../../util/scenario-utils';
import {
  NO_SCENARIO_NAME
} from '../../meep-constants';

const initialState = createNewScenario(NO_SCENARIO_NAME);

// CHANGE_SCENARIO
const CFG_CHANGE_SCENARIO = 'CFG_CHANGE_SCENARIO';
export function cfgChangeScenario(scenario) {
  return {
    type: CFG_CHANGE_SCENARIO,
    payload: scenario
  };
}

// CHANGE_SCENARIO_NAME
const CFG_CHANGE_SCENARIO_NAME = 'CFG_CHANGE_SCENARIO_NAME';
export function cfgChangeScenarioName(name) {
  return {
    type: CFG_CHANGE_SCENARIO_NAME,
    payload: name
  };
}

export function scenarioReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_CHANGE_SCENARIO_NAME:
    return updateObject(state, { name: action.payload });
  case CFG_CHANGE_SCENARIO:
    return action.payload;
  default:
    return state;
  }
}
