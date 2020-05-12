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

import _ from 'lodash';
import { createSelector } from 'reselect';
import { updateObject } from '../../util/object-util';
import { EXEC_STATE_IDLE } from '../../meep-constants';

// Pod phases
const CORE_PODS_PHASE_RUNNING = 'Running';
// const CORE_PODS_PHASE_PENDING = 'Pending';
// const SCENARIO_PODS_PHASE_RUNNING = 'Running';
const SCENARIO_PODS_PHASE_TERMINATING = 'Terminating';
const SCENARIO_PODS_PHASE_PENDING = 'Pending';

// Base selectors
const serviceMapsSelector = state => state.exec.state.serviceMaps;
const scenarioPodsPhasesSelector = state => state.exec.state.scenarioPodsPhases;
const corePodsPhasesSelector = state => state.exec.state.corePodsPhases;
const sandboxSelector = state => state.ui.sandbox;


// returns true if all core pods are in the RUNNING phase
export const corePodsRunning = createSelector(
  corePodsPhasesSelector,
  sandboxSelector,
  (phases, sandbox) => {
    return _.reduce(
      phases,
      (status, podStatus) => {
        return status && ((podStatus.sandbox === 'default' || podStatus.sandbox === sandbox) ?
          podStatus.logicalState === CORE_PODS_PHASE_RUNNING : true);
      },
      true
    );
  }
);

export const corePodsErrors = createSelector(
  corePodsPhasesSelector,
  sandboxSelector,
  (phases, sandbox) => {
    var statii = _.chain(phases)
      .filter(item => {
        return (item.sandbox === 'default' || item.sandbox === sandbox) &&
          item.logicalState !== CORE_PODS_PHASE_RUNNING;
      })
      .map(p => {
        return { name: p.name, status: p.logicalState };
      })
      .value();
    return statii;
  }
);

// Returns whether there scenario pods terminating
export const scenarioPodsPending = createSelector(
  scenarioPodsPhasesSelector,
  pods => {
    var phasePending = false;
    _.each(pods, pod => {
      phasePending =
        phasePending || pod.logicalState === SCENARIO_PODS_PHASE_PENDING;
    });

    return phasePending === true;
  }
);

export const scenarioPodsTerminating = createSelector(
  scenarioPodsPhasesSelector,
  pods => {
    var phaseTerminating = false;
    _.each(pods, pod => {
      phaseTerminating =
        phaseTerminating ||
        pod.logicalState === SCENARIO_PODS_PHASE_TERMINATING;
    });

    return phaseTerminating === true;
  }
);

export const scenarioPodsTerminated = createSelector(
  scenarioPodsPhasesSelector,
  pods => {
    return pods === null || pods === undefined || !pods.length;
  }
);

// Returns a list of scenario posds info and adds serviceMaps to external ones
export const podsWithServiceMaps = createSelector(
  [serviceMapsSelector, scenarioPodsPhasesSelector],
  (sms, spps) => {
    var podsWithInfo = _.map(spps, spp => {
      // If has info, add it
      var newSpp = null;
      _.each(sms, sm => {
        if (spp.name === sm.node) {
          if (!newSpp) {
            newSpp = updateObject({}, spp);
            newSpp.ingressServiceMap = [];
            newSpp.egressServiceMap = [];
          }
          _.each(sm.ingressServiceMap, entry => {
            newSpp.ingressServiceMap.push(entry);
          });
          _.each(sm.egressServiceMap, entry => {
            newSpp.egressServiceMap.push(entry);
          });
        }
      });

      if (newSpp) {
        return newSpp;
      }
      return spp;
    });

    return podsWithInfo;
  }
);

const EXEC_CHANGE_SCENARIO_STATE = 'EXEC_CHANGE_SCENARIO_STATE';
export function execChangeScenarioState(state) {
  return {
    type: EXEC_CHANGE_SCENARIO_STATE,
    payload: state
  };
}

const EXEC_CHANGE_CORE_PODS_PHASES = 'EXEC_CHANGE_CORE_PODS_PHASES';
export function execChangeCorePodsPhases(pods) {
  return {
    type: EXEC_CHANGE_CORE_PODS_PHASES,
    payload: pods
  };
}

const EXEC_CHANGE_SCENARIO_PODS_PHASES = 'EXEC_CHANGE_SCENARIO_PODS_PHASES';
export function execChangeScenarioPodsPhases(phases) {
  return {
    type: EXEC_CHANGE_SCENARIO_PODS_PHASES,
    payload: phases
  };
}

const EXEC_CHANGE_SERVICE_MAPS = 'EXEC_CHANGE_SERVICE_MAPS';
export function execChangeServiceMaps(serviceMaps) {
  return {
    type: EXEC_CHANGE_SERVICE_MAPS,
    payload: serviceMaps
  };
}

const CHANGE_OK_TO_TERMINATE = 'CHANGE_OK_TO_TERMINATE';
export function execChangeOkToTerminate(ok) {
  return {
    type: CHANGE_OK_TO_TERMINATE,
    payload: ok
  };
}

const EXEC_CHANGE_REPLAY_STATUS = 'EXEC_CHANGE_REPLAY_STATUS';
export function execChangeReplayStatus(status) {
  return {
    type: EXEC_CHANGE_REPLAY_STATUS,
    payload: status
  };
}

// Initial state
const initialState = {
  scenario: EXEC_STATE_IDLE,
  terminateButtonEnabled: false,
  corePodsPhases: [],
  scenarioPodsPhases: [],
  serviceMaps: [],
  okToTerminate: false,
  replayStatus: null
};

export function stateReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SCENARIO_STATE:
    return updateObject(state, { scenario: action.payload });
  case EXEC_CHANGE_CORE_PODS_PHASES:
    return updateObject(state, { corePodsPhases: action.payload });
  case EXEC_CHANGE_SCENARIO_PODS_PHASES:
    return updateObject(state, { scenarioPodsPhases: action.payload });
  case EXEC_CHANGE_SERVICE_MAPS:
    return updateObject(state, { serviceMaps: action.payload });
  case CHANGE_OK_TO_TERMINATE:
    return updateObject(state, { okToTerminate: action.payload });
  case EXEC_CHANGE_REPLAY_STATUS:
    return updateObject(state, { replayStatus: action.payload });
  default:
    return state;
  }
}
