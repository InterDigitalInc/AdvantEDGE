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

const EXEC_CHANGE_SCENARIO_LIST = 'EXEC_CHANGE_SCENARIO_LIST';
function execChangeScenarioList(scenarios) {
  return {
    type: EXEC_CHANGE_SCENARIO_LIST,
    payload: scenarios
  };
}

export {
  EXEC_CHANGE_SCENARIO_LIST,
  execChangeScenarioList
};

const initialState = {
  scenarios: []
};

export function execApiResultsReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SCENARIO_LIST:
    return updateObject(state, {scenarios: action.payload});
  default:
    return state;
  }
}
