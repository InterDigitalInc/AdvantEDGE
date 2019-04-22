/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import _ from 'lodash';
import {createSelector } from 'reselect';
import { updateObject } from '../../util/update';
import { EXEC_STATE_IDLE } from '../../meep-constants';

// Pod phases
const CORE_PODS_PHASE_RUNNING = 'Running';
const CORE_PODS_PHASE_PENDING = 'Pending';
const SCENARIO_PODS_PHASE_RUNNING = 'Running';
const SCENARIO_PODS_PHASE_TERMINATING = 'Terminating';
const SCENARIO_PODS_PHASE_PENDING = 'Pending';

// Base selectors
const serviceMapsSelector = state => state.exec.state.serviceMaps;
const scenarioPodsPhasesSelector = state => state.exec.state.scenarioPodsPhases;
const corePodsPhasesSelector = state => state.exec.state.corePodsPhases;

// returns true if all core pods are in the RUNNING phase
const corePodsRunning = createSelector(corePodsPhasesSelector, (phases) => {
  return _.reduce(phases, (status, pod) => {
    return status && pod.logicalState === CORE_PODS_PHASE_RUNNING;
  }, true);
});

const corePodsErrors = createSelector(corePodsPhasesSelector, (phases) => {
  var statii = _.chain(phases).map(p => {
    return {name: p.name, status: p.logicalState};
  }).filter((item) => {
    return item.status != CORE_PODS_PHASE_RUNNING;
  }).value();
  return statii;
});

// Returns whether there scenario pods terminating
const scenarioPodsPending = createSelector(scenarioPodsPhasesSelector, (pods) => {
  var phasePending = false;
  _.each(pods, (pod) => {
    phasePending |= (pod.logicalState === SCENARIO_PODS_PHASE_PENDING);
  });

  return phasePending == true;
});

const scenarioPodsTerminating = createSelector(scenarioPodsPhasesSelector, (pods) => {
  var phaseTerminating = false;
  _.each(pods, (pod) => {
    phaseTerminating |= (pod.logicalState === SCENARIO_PODS_PHASE_TERMINATING);
  });

  return phaseTerminating == true;
});

const scenarioPodsTerminated = createSelector(scenarioPodsPhasesSelector, (pods) => {
  return pods == null || !pods.length;
});

// Returns a list of scenario posds info and adds serviceMaps to external ones
const podsWithServiceMaps = createSelector([serviceMapsSelector, scenarioPodsPhasesSelector], (sms, spps) => {
  var podsWithInfo = _.map(spps, (spp) => {
    // If has info, add it
    var newSpp = null;
    _.each(sms, (sm) => {
      if (spp.name === sm.client) {
        if (!newSpp) {
          newSpp = updateObject({}, spp);
          newSpp.serviceMaps = [];
        }
        _.each(sm.serviceMap, (entry) => {
          newSpp.serviceMaps.push(entry);
        });
      }
    });

    if (newSpp) {
      return newSpp;
    }
    return spp;
  });

  return podsWithInfo;
});

const EXEC_CHANGE_SCENARIO_STATE = 'EXEC_CHANGE_SCENARIO_STATE';
function execChangeScenarioState(state) {
  return {
    type: EXEC_CHANGE_SCENARIO_STATE,
    payload: state
  };
}

const EXEC_CHANGE_CORE_PODS_PHASES = 'EXEC_CHANGE_CORE_PODS_PHASES';
function execChangeCorePodsPhases(pods) {
  return {
    type: EXEC_CHANGE_CORE_PODS_PHASES,
    payload: pods
  };
}

const EXEC_CHANGE_SCENARIO_PODS_PHASES = 'EXEC_CHANGE_SCENARIO_PODS_PHASES';
function execChangeScenarioPodsPhases(phases) {
  return {
    type: EXEC_CHANGE_SCENARIO_PODS_PHASES,
    payload: phases
  };
}

const EXEC_CHANGE_SERVICE_MAPS = 'EXEC_CHANGE_SERVICE_MAPS';
function execChangeServiceMaps(serviceMaps) {
  return {
    type: EXEC_CHANGE_SERVICE_MAPS,
    payload: serviceMaps
  };
}

const CHANGE_OK_TO_TERMINATE = 'CHANGE_OK_TO_TERMINATE';
function execChangeOkToTerminate(ok) {
  return {
    type: CHANGE_OK_TO_TERMINATE,
    payload: ok
  };
}

export {
  // Actions
  execChangeScenarioState,
  execChangeScenarioPodsPhases,
  execChangeCorePodsPhases,
  execChangeServiceMaps,
  execChangeOkToTerminate,

  // Selectors
  corePodsRunning,
  corePodsErrors,
  podsWithServiceMaps,
  scenarioPodsPending,
  scenarioPodsTerminating,
  scenarioPodsTerminated,

  // Action types
  EXEC_CHANGE_CORE_PODS_PHASES,
  EXEC_CHANGE_SCENARIO_STATE,
  EXEC_CHANGE_SCENARIO_PODS_PHASES,
  EXEC_CHANGE_SERVICE_MAPS,

  // Core pods phases
  CORE_PODS_PHASE_RUNNING,
  CORE_PODS_PHASE_PENDING,

  // Scenario pods phases
  SCENARIO_PODS_PHASE_RUNNING,
  SCENARIO_PODS_PHASE_PENDING,
  SCENARIO_PODS_PHASE_TERMINATING
};

// Initial state
const initialState = {
  scenario: EXEC_STATE_IDLE,
  terminateButtonEnabled: false,
  corePodsPhases: [],
  scenarioPodsPhases: [],
  serviceMaps: [],
  okToTerminate: false
};

export function stateReducer(state = initialState, action) {
  switch (action.type) {
  case EXEC_CHANGE_SCENARIO_STATE:
    return updateObject(state, {scenario: action.payload});
  case EXEC_CHANGE_CORE_PODS_PHASES:
    return updateObject(state, {corePodsPhases: action.payload});
  case EXEC_CHANGE_SCENARIO_PODS_PHASES:
    return updateObject(state, {scenarioPodsPhases: action.payload});
  case EXEC_CHANGE_SERVICE_MAPS:
    return updateObject(state, {serviceMaps: action.payload});
  case CHANGE_OK_TO_TERMINATE:
    return updateObject(state, {okToTerminate: action.payload});
  default:
    return state;
  }
}
