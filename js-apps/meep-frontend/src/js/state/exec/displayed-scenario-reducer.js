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

// EXEC_CHANGE_SCENARIO
export const EXEC_CHANGE_DISPLAYED_SCENARIO = 'EXEC_CHANGE_DISPLAYED_SCENARIO';
function execChangeDisplayedScenario(scenario) {
  return {
    type: EXEC_CHANGE_DISPLAYED_SCENARIO,
    payload: scenario
  };
}

export {
  // Action creators
  execChangeDisplayedScenario
};

const initialState = {};

export function displayedScenarioReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_DISPLAYED_SCENARIO:
    return action.payload;
  default:
    return state;
  }
}
