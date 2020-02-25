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

import { combineReducers } from 'redux';

import uiReducer from './ui';
import cfgReducer from './cfg';
import execReducer from './exec';
import settingsReducer from './settings';
import monitorReducer from './monitor';

const appReducer = combineReducers({
  ui: uiReducer,
  cfg: cfgReducer,
  exec: execReducer,
  monitor: monitorReducer,
  settings: settingsReducer
});

const UI_SET_DEFAULT_STATE = 'UI_SET_DEFAULT_STATE';
export function meepSetDefaultState() {
  return {
    type: UI_SET_DEFAULT_STATE,
    payload: ''
  };
}

const meepReducer = (state, action) => {
  switch (action.type) {
  case UI_SET_DEFAULT_STATE:
    return appReducer(undefined, action);
  default:
    return appReducer(state, action);
  }
};

export default meepReducer;