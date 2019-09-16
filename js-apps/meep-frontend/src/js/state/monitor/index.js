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

const initialState = {
  currentDashboardUrl: '',
  dashboardOptions: [
    {
      label:  'Latency Dashboard',
      value: 'http://' + location.hostname + ':32003/app/kibana#/dashboard/6745bb30-c29c-11e8-95a0-933bd4e05896?embed=true&_g=(refreshInterval%3A(pause%3A!f%2Cvalue%3A5000)%2Ctime%3A(from%3Anow-60s%2Cmode%3Arelative%2Cto%3Anow))'
    },
    {
      label:  'Demo Service Internal UE (ue1)',
      value: 'http://' + location.hostname + ':32003/app/kibana#/dashboard/434d37b0-1b6d-11e9-b72d-e70da2a5e139?embed=true&_g=(refreshInterval%3A(pause%3A!f%2Cvalue%3A5000)%2Ctime%3A(from%3Anow-15m%2Cmode%3Arelative%2Cto%3Anow))'
    },
    {
      label:  'Demo Service External UE (ue2-ext)',
      value: 'http://' + location.hostname + ':32003/app/kibana#/dashboard/788a4f70-1b73-11e9-b72d-e70da2a5e139?embed=true&_g=(refreshInterval%3A(pause%3A!f%2Cvalue%3A5000)%2Ctime%3A(from%3Anow-15m%2Cmode%3Arelative%2Cto%3Anow))'
    }
  ],
  editedDashboardOptions: null
};

const ADD_DASHBOARD_OPTION = 'ADD_DASHBOARD_OPTION';
export const addDashboardOption = (option) => {
  return {
    type: ADD_DASHBOARD_OPTION,
    payload: option
  };
};
  
const CHANGE_DASHBOARD_URL = 'CHANGE_DASHBOARD_URL';
export function changeDashboardUrl(url) {
  return {
    type: CHANGE_DASHBOARD_URL,
    payload: url
  };
}

const CHANGE_EDITED_DASHBOARD_OPTIONS = 'CHANGE_EDITED_DASHBOARD_OPTIONS';
export function changeEditedDashboardOptions(mode) {
  return {
    type: CHANGE_EDITED_DASHBOARD_OPTIONS,
    payload: mode
  };
}

const CHANGE_DASHBOARD_OPTIONS = 'CHANGE_DASHBOARD_OPTIONS';
export function changeDashboardOptions(mode) {
  return {
    type: CHANGE_DASHBOARD_OPTIONS,
    payload: mode
  };
}
  
export default function settingsReducer(state = initialState, action) {
  switch (action.type) {
  case CHANGE_DASHBOARD_URL:
    return {...state, currentDashboardUrl: action.payload};
  case ADD_DASHBOARD_OPTION:
    return {...state, dashboardOptions: [...state.dashboardOptions, action.payload]};
  case CHANGE_EDITED_DASHBOARD_OPTIONS:
    return {...state, editedDashboardOptions: action.payload};
  case CHANGE_DASHBOARD_OPTIONS:
    return {...state, dashboardOptions: action.payload};
  default:
    return state;
  }
}
  