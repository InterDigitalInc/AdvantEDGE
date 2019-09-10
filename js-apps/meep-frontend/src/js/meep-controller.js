/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
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
import { fixMetricsValuesMiddleware } from './middlewares/fix-metrics-values-middleware';


// To uncomment when save state is fixed
import {
  saveState
  // loadState
} from './util/persist';

// UI Components
import MeepContainer from './containers/meep-container';

// Constants
import {
  PAGE_CONFIGURE
} from './state/ui';

import {
  TYPE_CFG,
  TYPE_EXEC,
  CFG_STATE_IDLE,
  EXEC_STATE_IDLE,
  NO_SCENARIO_NAME
} from './meep-constants';

import {
  createNewScenario
} from './util/scenario-utils';

// MEEP Controller Frontend state information
const meep = {
  ui: {
    page: PAGE_CONFIGURE
  },
  cfg: {
    type: TYPE_CFG,
    state: CFG_STATE_IDLE,
    scenario: createNewScenario(NO_SCENARIO_NAME),
    vis: {
      type: TYPE_CFG,
      network: {},
      options: {},
      data: {
        nodes: [],
        edges: []
      },
      showConfig: false
    },
    table: {
      data: [],
      selected: [],
      order: 'asc',
      orderBy: 'name',
      rowsPerPage: 10,
      page: 0,
      refresh: () => {}
    }
  },
  exec: {
    type: TYPE_EXEC,
    state: EXEC_STATE_IDLE,
    scenario: {
      name: NO_SCENARIO_NAME
    },
    vis: {
      type: TYPE_EXEC,
      network: {},
      options: {},
      data: {
        nodes: [],
        edges: []
      },
      showConfig: false
    },
    table: {
      data: [],
      selected: [],
      order: 'asc',
      orderBy: 'name',
      rowsPerPage: 10,
      page: 0,
      refresh: () => {}
    }
  },
  settings: {
    debug: false
  }
};

// Attempting to load the meep (state) from localStorage first.
// If not found, initialize the Redux store with the above meep object.
// Will mode that code out when references to DOM elements are factored out (VIS stuff)
function createState(meepObject) {
  var state = JSON.parse(JSON.stringify(meepObject));
  delete state.cfg.vis.reportContainer;
  delete state.cfg.table.refresh;
  delete state.exec.vis.containers;
  delete state.exec.table.refresh;
  delete state.exec.vis.reportContainer;

  state.exec.state = {
    scenario: meepObject.exec.state,
    corePodsPhases: [],
    scenarioPodsPhases: [],
    serviceMaps: []
  };

  state.cfg.table.selected = [];

  state.cfg.elementConfiguration = {
    configuredElement: null,
    configurationMode: null
  };

  state.ui =  {
    devMode: false,
    currentDialog: '',
    execShowApps: false,
    mainDrawerOpen: true
  };

  return state;
}

// Initialize variables and listeners when document ready
var loadedState = null; //loadState();
let meepState = loadedState ? loadedState : createState(meep);

// Uncomment if logger middleware is needed
// var logger = store => () => action => {
//   console.log(`logger - action: ${action.type}. payload: `, action.payload);
//   console.log('state: ', store.getState());
// };

const meepStore = createStore(meepReducer, meepState, applyMiddleware(thunk, execDisplayedScenarioMiddleware, fixMetricsValuesMiddleware));
window.meepStore = meepStore;

// TODO: fix circularity in store
meepStore.subscribe(() => {
  saveState(meepStore.getState());
});

// Monitor Page
let meepContainerPlaceholder = document.getElementById('meep-container');
ReactDOM.render(
  <Provider store={meepStore}>
    <MeepContainer />
  </Provider>,
  meepContainerPlaceholder
);
