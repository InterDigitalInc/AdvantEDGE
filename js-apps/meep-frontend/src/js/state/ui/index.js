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

export const PAGE_CONFIGURE = 'page-configure-link';
export const PAGE_EXECUTE = 'page-execute-link';
export const PAGE_MONITOR = 'page-monitor-link';
export const PAGE_SETTINGS = 'page-settings-link';

const UE_MOBILITY_EVENT = 'UE-MOBILITY';
const NETWORK_CHARACTERISTICS_EVENT = 'NETWORK-CHARACTERISTICS-UPDATE';

// Change the current page
const CHANGE_CURRENT_PAGE = 'CHANGE_CURRENT_PAGE';
function uiChangeCurrentPage(page) {
  return {
    type: CHANGE_CURRENT_PAGE,
    payload: page
  };
}

const TOGGLE_MAIN_DRAWER = 'TOGGLE_MAIN_DRAWER';
function uiToggleMainDrawer() {
  return {
    type: TOGGLE_MAIN_DRAWER,
    payload: null
  };
}

const EXEC_CHANGE_CURRENT_EVENT = 'EXEC_CHANGE_CURRENT_EVENT';
function uiExecChangeCurrentEvent(event) {
  return {
    type: EXEC_CHANGE_CURRENT_EVENT,
    payload: event
  };
}

const EXEC_CHANGE_EVENT_CREATION_MODE = 'EXEC_CHANGE_EVENT_CREATION_MODE';
function uiExecChangeEventCreationMode(val) {
  return {
    type: EXEC_CHANGE_EVENT_CREATION_MODE,
    payload: val
  };
}

const UI_CHANGE_DEV_MODE = 'UI_CHANGE_DEV_MODE';
function uiChangeDevMode(mode) {
  return {
    type: UI_CHANGE_DEV_MODE,
    payload: mode
  };
}

// Dialog Types
// CFG
const IDC_DIALOG_OPEN_SCENARIO = 'IDC_DIALOG_OPEN_SCENARIO';
const IDC_DIALOG_NEW_SCENARIO = 'IDC_DIALOG_NEW_SCENARIO';
const IDC_DIALOG_SAVE_SCENARIO = 'IDC_DIALOG_SAVE_SCENARIO';
const IDC_DIALOG_DELETE_SCENARIO = 'IDC_DIALOG_DELETE_SCENARIO';
const IDC_DIALOG_EXPORT_SCENARIO = 'IDC_DIALOG_EXPORT_SCENARIO';
const IDC_DIALOG_TERMINATE_SCENARIO = 'IDC_DIALOG_TERMINATE_SCENARIO';
const IDC_DIALOG_CONFIRM = 'IDC_DIALOG_CONFIRM';


// EXEC
const IDC_DIALOG_DEPLOY_SCENARIO = 'IDC_DIALOG_DEPLOY_SCENARIO';


const UI_CHANGE_CURRENT_DIALOG = 'UI_CHANGE_CURRENT_DIALOG';
const uiChangeCurrentDialog = (type) => {
  return {
    type: UI_CHANGE_CURRENT_DIALOG,
    payload: type
  };
};

const UI_SET_AUTOMATIC_REFRESH = 'UI_SET_AUTOMATIC_REFRESH';
const uiSetAutomaticRefresh = (val) => {
  return {
    type: UI_SET_AUTOMATIC_REFRESH,
    payload: val
  };
};

const UI_CHANGE_REFRESH_INTERVAL = 'UI_CHANGE_REFRESH_INTERVAL';
const uiChangeRefreshInterval = (val) => {
  return {
    type: UI_CHANGE_REFRESH_INTERVAL,
    payload: val
  };
};

const UI_EXEC_CHANGE_SHOW_APPS = 'UI_EXEC_CHANGE_SHOW_APPS';
const uiExecChangeShowApps = (show) => {
  return {
    type: UI_EXEC_CHANGE_SHOW_APPS,
    payload: show
  };
};

export {
  // Event types
  UE_MOBILITY_EVENT,
  NETWORK_CHARACTERISTICS_EVENT,

  // Action types
  EXEC_CHANGE_CURRENT_EVENT,

  // Dialogs types
  IDC_DIALOG_OPEN_SCENARIO,
  IDC_DIALOG_NEW_SCENARIO,
  IDC_DIALOG_SAVE_SCENARIO,
  IDC_DIALOG_DELETE_SCENARIO,
  IDC_DIALOG_EXPORT_SCENARIO,
  IDC_DIALOG_DEPLOY_SCENARIO,
  IDC_DIALOG_TERMINATE_SCENARIO,
  IDC_DIALOG_CONFIRM,

  // Action creators
  uiChangeCurrentPage,
  uiToggleMainDrawer,
  uiExecChangeEventCreationMode,
  uiExecChangeCurrentEvent,
  uiChangeDevMode,
  uiChangeCurrentDialog,
  uiSetAutomaticRefresh,
  uiChangeRefreshInterval,
  uiExecChangeShowApps
};

const initialState = {
  page: PAGE_CONFIGURE,
  mainDrawerOpen: true,
  eventCreationMode: false,
  execCurrentEvent: null,
  currentEventType: UE_MOBILITY_EVENT,
  devMode: false,
  currentDialog: '',
  automaticRefresh: false,
  refreshInterval: 1000,
  execShowApps: false
};

export default function uiReducer(state = initialState, action) {
  switch (action.type) {
  case CHANGE_CURRENT_PAGE:
    return updateObject(state, {page: action.payload});
  case TOGGLE_MAIN_DRAWER:
    return updateObject(state, {mainDrawerOpen: !state.mainDrawerOpen});
  case EXEC_CHANGE_CURRENT_EVENT:
    return updateObject(state, {execCurrentEvent: action.payload});
  case UI_CHANGE_DEV_MODE:
    return updateObject(state, {devMode: action.payload || false});
  case UI_CHANGE_CURRENT_DIALOG:
    return updateObject(state, {currentDialog: action.payload});
  case EXEC_CHANGE_EVENT_CREATION_MODE:
    return updateObject(state, {eventCreationMode: action.payload});
  case UI_SET_AUTOMATIC_REFRESH:
    return updateObject(state, {automaticRefresh: action.payload});
  case UI_CHANGE_REFRESH_INTERVAL:
    return updateObject(state, {refreshInterval: action.payload});
  case UI_EXEC_CHANGE_SHOW_APPS:
    return updateObject(state, {execShowApps: action.payload});
  default:
    return state;
  }
}
