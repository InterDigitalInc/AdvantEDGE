/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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

// Dashhboard configuration fields
export const DASH_CFG_VIEW_TYPE = 'viewType';
export const DASH_CFG_SOURCE_NODE_SELECTED = 'sourceNodeSelected';
export const DASH_CFG_DEST_NODE_SELECTED = 'destNodeSelected';
export const DASH_CFG_PARTICIPANTS = 'participants';
export const DASH_CFG_MAX_MSG_COUNT = 'maxMsgCount';
export const DASH_CFG_START_TIME = 'startTime';


export const getDashCfgFieldVal = (cfg, field) => {
  return (cfg && cfg[field]) ? cfg[field].val : null;
};

export const getDashCfgFieldErr = (cfg, field) => {
  return (cfg && cfg[field]) ? cfg[field].err : null;
};

export const setDashCfgField = (cfg, field, val, err) => {
  if (cfg) {
    cfg[field] = {
      val: val,
      err: err
    };
  }
};
