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

import { CFG_STATE_IDLE } from '../../meep-constants';

const initialState = CFG_STATE_IDLE;

// CHANGE_STATE
const CFG_CHANGE_STATE = 'CFG_CHANGE_STATE';
export function cfgChangeState(state) {
  return {
    type: CFG_CHANGE_STATE,
    payload: state
  };
}

export function stateReducer(state = initialState, action) {
  switch (action.type) {
  case CFG_CHANGE_STATE:
    return action.payload;
  default:
    return state;
  }
}
