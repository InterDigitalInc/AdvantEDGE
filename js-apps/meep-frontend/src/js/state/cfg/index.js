/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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
