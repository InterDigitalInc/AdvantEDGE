/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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

const initialState = {
  metrics: [],
  participants: [],
  chart: ''
};

// EXEC_CHANGE_SEQ
const EXEC_CHANGE_SEQ = 'EXEC_CHANGE_SEQ';
export function execChangeSeq(seq) {
  return {
    type: EXEC_CHANGE_SEQ,
    payload: seq
  };
}

// EXEC_CHANGE_SEQ_METRICS
const EXEC_CHANGE_SEQ_METRICS = 'EXEC_CHANGE_SEQ_METRICS';
export function execChangeSeqMetrics(metrics) {
  return {
    type: EXEC_CHANGE_SEQ_METRICS,
    payload: metrics
  };
}

// EXEC_CHANGE_SEQ_PARTICIPANTS
const EXEC_CHANGE_SEQ_PARTICIPANTS = 'EXEC_CHANGE_SEQ_PARTICIPANTS';
export function execChangeSeqParticipants(participants) {
  return {
    type: EXEC_CHANGE_SEQ_PARTICIPANTS,
    payload: participants
  };
}

// EXEC_CHANGE_SEQ_CHART
const EXEC_CHANGE_SEQ_CHART = 'EXEC_CHANGE_SEQ_CHART';
export function execChangeSeqChart(chart) {
  return {
    type: EXEC_CHANGE_SEQ_CHART,
    payload: chart
  };
}

export function execSeqReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SEQ:
    return action.payload;
  case EXEC_CHANGE_SEQ_METRICS:
    return updateObject(state, { metrics: action.payload });
  case EXEC_CHANGE_SEQ_PARTICIPANTS:
    return updateObject(state, { participants: action.payload });
  case EXEC_CHANGE_SEQ_CHART:
    return updateObject(state, { chart: action.payload });
  default:
    return state;
  }
}
