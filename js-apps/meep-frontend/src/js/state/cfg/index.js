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
import { typeReducer } from './type-reducer';
import { stateReducer } from './state-reducer';
import { scenarioReducer } from './scenario-reducer';
import { cfgVisReducer } from './vis-reducer';
import { cfgTableReducer } from './table-reducer';
import { cfgElementConfigurationReducer } from './element-configuration';
import { cfgApiResultsReducer } from './api-results';

export * from './type-reducer';
export * from './state-reducer';
export * from './scenario-reducer';
export * from './vis-reducer';
export * from './table-reducer';
export * from './element-configuration';
export * from './api-results';

const cfgReducer = combineReducers({
  type: typeReducer,
  state: stateReducer,
  scenario: scenarioReducer,
  vis: cfgVisReducer,
  table: cfgTableReducer,
  elementConfiguration: cfgElementConfigurationReducer,
  apiResults: cfgApiResultsReducer
});

export default cfgReducer;
