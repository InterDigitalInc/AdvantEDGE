import { combineReducers } from 'redux';

import uiReducer from './ui';
import cfgReducer from './cfg';
import execReducer from './exec';
import settingsReducer from './settings';

const meepReducer = combineReducers({
  ui: uiReducer,
  cfg: cfgReducer,
  exec: execReducer,
  settings: settingsReducer
});

export default meepReducer;
