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
import { createElem } from '../../util/elem-utils';

const EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT =
  'EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT';

// CFG_SET_EDITED_ELEMENT
function execChangeSelectedScenarioElement(element) {
  return {
    type: EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT,
    payload: element
  };
}

export { execChangeSelectedScenarioElement };

const initialState = createElem('dummy');

export function execSelectedScenarioElement(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SELECTED_SCENARIO_ELEMENT:
    return updateObject({}, action.payload);
  default:
    return state;
  }
}
