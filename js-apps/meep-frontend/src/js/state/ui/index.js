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

import { updateObject } from '../../util/object-util';
import {
  PAGE_CONFIGURE,
  VIEW_NAME_NONE,
  NET_TOPOLOGY_VIEW,
  MOBILITY_EVENT
} from '../../meep-constants';

const initialState = {
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
  dashCfgMode: false,
  eventCfgMode: false,
  dashboardView1: NET_TOPOLOGY_VIEW,
  dashboardView2: VIEW_NAME_NONE,
  sourceNodeSelected: '',
  destNodeSelected: '',
  eventReplayMode: false,
  eventReplayLoop: false,
  replayFiles: [],
  replayFileSelected: '',
  replayFileDesc: ''
};

// Change the current page
const CHANGE_CURRENT_PAGE = 'CHANGE_CURRENT_PAGE';
export function uiChangeCurrentPage(page) {
  return {
    type: CHANGE_CURRENT_PAGE,
    payload: page
  };
}

const TOGGLE_MAIN_DRAWER = 'TOGGLE_MAIN_DRAWER';
export function uiToggleMainDrawer() {
  return {
    type: TOGGLE_MAIN_DRAWER,
    payload: null
  };
}

const UI_EXEC_CHANGE_SANDBOX = 'UI_EXEC_CHANGE_SANDBOX';
export function uiExecChangeSandbox(name) {
  return {
    type: UI_EXEC_CHANGE_SANDBOX,
    payload: name
  };
}

const UI_EXEC_CHANGE_SANDBOX_LIST = 'UI_EXEC_CHANGE_SANDBOX_LIST';
export function uiExecChangeSandboxList(list) {
  return {
    type: UI_EXEC_CHANGE_SANDBOX_LIST,
    payload: list
  };
}

const UI_EXEC_CHANGE_CURRENT_EVENT = 'UI_EXEC_CHANGE_CURRENT_EVENT';
export function uiExecChangeCurrentEvent(event) {
  return {
    type: UI_EXEC_CHANGE_CURRENT_EVENT,
    payload: event
  };
}

const UI_EXEC_CHANGE_EVENT_CREATION_MODE = 'UI_EXEC_CHANGE_EVENT_CREATION_MODE';
export function uiExecChangeEventCreationMode(val) {
  return {
    type: UI_EXEC_CHANGE_EVENT_CREATION_MODE,
    payload: val
  };
}

const UI_EXEC_CHANGE_EVENT_REPLAY_MODE = 'UI_EXEC_CHANGE_EVENT_REPLAY_MODE';
export function uiExecChangeEventReplayMode(val) {
  return {
    type: UI_EXEC_CHANGE_EVENT_REPLAY_MODE,
    payload: val
  };
}

const UI_EXEC_CHANGE_DASH_CFG_MODE = 'UI_EXEC_CHANGE_DASH_CFG_MODE';
export function uiExecChangeDashCfgMode(val) {
  return {
    type: UI_EXEC_CHANGE_DASH_CFG_MODE,
    payload: val
  };
}

const UI_EXEC_CHANGE_EVENT_CFG_MODE = 'UI_EXEC_CHANGE_EVENT_CFG_MODE';
export function uiExecChangeEventCfgMode(val) {
  return {
    type: UI_EXEC_CHANGE_EVENT_CFG_MODE,
    payload: val
  };
}

const UI_CHANGE_DEV_MODE = 'UI_CHANGE_DEV_MODE';
export function uiChangeDevMode(mode) {
  return {
    type: UI_CHANGE_DEV_MODE,
    payload: mode
  };
}

const UI_CHANGE_CURRENT_DIALOG = 'UI_CHANGE_CURRENT_DIALOG';
export function uiChangeCurrentDialog(type) {
  return {
    type: UI_CHANGE_CURRENT_DIALOG,
    payload: type
  };
}

const UI_SET_AUTOMATIC_REFRESH = 'UI_SET_AUTOMATIC_REFRESH';
export function uiSetAutomaticRefresh(val) {
  return {
    type: UI_SET_AUTOMATIC_REFRESH,
    payload: val
  };
}

const UI_CHANGE_REFRESH_INTERVAL = 'UI_CHANGE_REFRESH_INTERVAL';
export function uiChangeRefreshInterval(val) {
  return {
    type: UI_CHANGE_REFRESH_INTERVAL,
    payload: val
  };
}

const UI_EXEC_CHANGE_SHOW_APPS = 'UI_EXEC_CHANGE_SHOW_APPS';
export function uiExecChangeShowApps(show) {
  return {
    type: UI_EXEC_CHANGE_SHOW_APPS,
    payload: show
  };
}

const UI_EXEC_CHANGE_DASHBOARD_VIEW1 = 'UI_EXEC_CHANGE_DASHBOARD_VIEW1';
export function uiExecChangeDashboardView1(name) {
  return {
    type: UI_EXEC_CHANGE_DASHBOARD_VIEW1,
    payload: name
  };
}

const UI_EXEC_CHANGE_DASHBOARD_VIEW2 = 'UI_EXEC_CHANGE_DASHBOARD_VIEW2';
export function uiExecChangeDashboardView2(name) {
  return {
    type: UI_EXEC_CHANGE_DASHBOARD_VIEW2,
    payload: name
  };
}

const UI_EXEC_CHANGE_SOURCE_NODE_SELECTED = 'UI_EXEC_CHANGE_SOURCE_NODE_SELECTED';
export function uiExecChangeSourceNodeSelected(node) {
  return {
    type: UI_EXEC_CHANGE_SOURCE_NODE_SELECTED,
    payload: node
  };
}

const UI_EXEC_CHANGE_DEST_NODE_SELECTED = 'UI_EXEC_CHANGE_DEST_NODE_SELECTED';
export function uiExecChangeDestNodeSelected(node) {
  return {
    type: UI_EXEC_CHANGE_DEST_NODE_SELECTED,
    payload: node
  };
}

const UI_EXEC_CHANGE_REPLAY_FILES_LIST = 'UI_EXEC_CHANGE_REPLAY_FILES_LIST';
export function uiExecChangeReplayFilesList(replayFiles) {
  return {
    type: UI_EXEC_CHANGE_REPLAY_FILES_LIST,
    payload: replayFiles
  };
}

const UI_EXEC_CHANGE_REPLAY_FILE_SELECTED = 'UI_EXEC_CHANGE_REPLAY_FILE_SELECTED';
export function uiExecChangeReplayFileSelected(name) {
  return {
    type: UI_EXEC_CHANGE_REPLAY_FILE_SELECTED,
    payload: name
  };
}

const UI_EXEC_CHANGE_REPLAY_FILE_DESC = 'UI_EXEC_CHANGE_REPLAY_FILE_DESC';
export function uiExecChangeReplayFileDesc(desc) {
  return {
    type: UI_EXEC_CHANGE_REPLAY_FILE_DESC,
    payload: desc
  };
}

const UI_EXEC_CHANGE_REPLAY_LOOP = 'UI_EXEC_CHANGE_REPLAY_LOOP';
export const uiExecChangeReplayLoop = val => {
  return {
    type: UI_EXEC_CHANGE_REPLAY_LOOP,
    payload: val
  };
};

export default function uiReducer(state = initialState, action) {
  switch (action.type) {
  case CHANGE_CURRENT_PAGE:
    return updateObject(state, { page: action.payload });
  case TOGGLE_MAIN_DRAWER:
    return updateObject(state, { mainDrawerOpen: !state.mainDrawerOpen });
  case UI_EXEC_CHANGE_SANDBOX:
    return updateObject(state, { sandbox: action.payload });
  case UI_EXEC_CHANGE_SANDBOX_LIST:
    return updateObject(state, { sandboxes: action.payload });
  case UI_EXEC_CHANGE_CURRENT_EVENT:
    return updateObject(state, { execCurrentEvent: action.payload });
  case UI_CHANGE_DEV_MODE:
    return updateObject(state, { devMode: action.payload || false });
  case UI_CHANGE_CURRENT_DIALOG:
    return updateObject(state, { currentDialog: action.payload });
  case UI_EXEC_CHANGE_EVENT_CREATION_MODE:
    return updateObject(state, { eventCreationMode: action.payload });
  case UI_EXEC_CHANGE_EVENT_REPLAY_MODE:
    return updateObject(state, { eventReplayMode: action.payload });
  case UI_EXEC_CHANGE_DASH_CFG_MODE:
    return updateObject(state, { dashCfgMode: action.payload });
  case UI_EXEC_CHANGE_EVENT_CFG_MODE:
    return updateObject(state, { eventCfgMode: action.payload });
  case UI_SET_AUTOMATIC_REFRESH:
    return updateObject(state, { automaticRefresh: action.payload });
  case UI_CHANGE_REFRESH_INTERVAL:
    return updateObject(state, { refreshInterval: action.payload });
  case UI_EXEC_CHANGE_SHOW_APPS:
    return updateObject(state, { execShowApps: action.payload });
  case UI_EXEC_CHANGE_DASHBOARD_VIEW1:
    return updateObject(state, { dashboardView1: action.payload });
  case UI_EXEC_CHANGE_DASHBOARD_VIEW2:
    return updateObject(state, { dashboardView2: action.payload });
  case UI_EXEC_CHANGE_SOURCE_NODE_SELECTED:
    return updateObject(state, { sourceNodeSelected: action.payload });
  case UI_EXEC_CHANGE_DEST_NODE_SELECTED:
    return updateObject(state, { destNodeSelected: action.payload });
  case UI_EXEC_CHANGE_REPLAY_FILES_LIST:
    return updateObject(state, { replayFiles: action.payload });
  case UI_EXEC_CHANGE_REPLAY_FILE_SELECTED:
    return updateObject(state, { replayFileSelected: action.payload });
  case UI_EXEC_CHANGE_REPLAY_FILE_DESC:
    return updateObject(state, { replayFileDesc: action.payload });
  case UI_EXEC_CHANGE_REPLAY_LOOP:
    return updateObject(state, { eventReplayLoop: action.payload });
  default:
    return state;
  }
}
