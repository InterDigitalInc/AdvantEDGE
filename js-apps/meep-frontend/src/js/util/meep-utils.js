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

import {
  createNewScenario
} from './scenario-utils';

// MEEP Controller Frontend state information
export const createMeepState = ({ui}) => {
  return {
    ui: ui ? ui : {
      page: PAGE_CONFIGURE,
      mainDrawerOpen: true,
      eventCreationMode: false,
      execCurrentEvent: null,
      currentEventType: MOBILITY_EVENT, // Should be moved somewhere else
      devMode: false,
      currentDialog: '',
      automaticRefresh: false,
      refreshInterval: 1000,
      execShowApps: false,
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
};