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

const initialState = {
  data: []
};

// ECEC_UPDATE_APP
const EXEC_CHANGE_APP_INSTANCE_TABLE = 'EXEC_CHANGE_APP_INSTANCE_TABLE';
export function execChangeAppInstanceTable(data) {
  return {
    type: EXEC_CHANGE_APP_INSTANCE_TABLE,
    payload: data
  };
}

export function appInstanceTableReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_APP_INSTANCE_TABLE:
    return updateObject(state, { data: action.payload });
  default:
    return state;
  }
}
