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

import {
  TYPE_CFG,
  TYPE_EXEC,
  CFG_STATE_IDLE,
  EXEC_STATE_IDLE,
  NO_SCENARIO_NAME,
  PAGE_CONFIGURE,
  VIEW_NAME_NONE,
  VIS_VIEW,
  MOBILITY_EVENT
} from '../meep-constants';

import { createNewScenario } from './scenario-utils';

// MEEP Controller Frontend state information
export const createMeepState = ({ ui }) => {
  return {
    ui: ui
      ? ui
      : {
        page: PAGE_CONFIGURE,
        mainDrawerOpen: true,
        eventCreationMode: false,
        eventReplayMode: false,
        execCurrentEvent: null,
        execReplayFileSelected: '',
        currentEventType: MOBILITY_EVENT, // Should be moved somewhere else
        devMode: false,
        currentDialog: '',
        automaticRefresh: true,
        refreshInterval: 1000,
        execShowApps: false,
        eventReplayLoop: false,
        showDashboardConfig: false,
        dashboardConfigExpanded: false,
        dashboardView1: VIS_VIEW,
        dashboardView2: VIEW_NAME_NONE
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
      state: {
        scenario: EXEC_STATE_IDLE,
        terminateButtonEnabled: false,
        corePodsPhases: [],
        scenarioPodsPhases: [],
        serviceMaps: [],
        okToTerminate: false
      },
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
};
