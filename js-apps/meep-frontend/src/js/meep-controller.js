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

// Import CSS
import 'material-design-icons/iconfont/material-icons.css';
import 'vis/dist/vis.min.css';
import '../css/meep-controller.scss';

// Import module dependencies
import 'material-design-icons';
import React from 'react';
import ReactDOM from 'react-dom';
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import { Provider } from 'react-redux';
import meepReducer from './state/meep-reducer';
import { execDisplayedScenarioMiddleware } from './middlewares/exec-displayed-scenario-middleware';

import { saveUIState, loadUIState } from './util/persist';

import { createMeepState } from './util/meep-utils';

// UI Components
import MeepContainer from './containers/meep-container';

// Initialize variables and listeners when document ready
var loadedUIState = loadUIState();

// Uncomment if logger middleware is needed
// var logger = store => () => action => {
//   console.log(`logger - action: ${action.type}. payload: `, action.payload);
//   console.log('state: ', store.getState());
// };

const meepState = createMeepState({ ui: loadedUIState });

const meepStore = createStore(
  meepReducer,
  meepState,
  applyMiddleware(
    thunk,
    execDisplayedScenarioMiddleware
  )
);
window.meepStore = meepStore;

meepStore.subscribe(() => {
  saveUIState(meepStore.getState().ui);
});

// Monitor Page
let meepContainerPlaceholder = document.getElementById('meep-container');
ReactDOM.render(
  <Provider store={meepStore}>
    <MeepContainer />
  </Provider>,
  meepContainerPlaceholder
);
