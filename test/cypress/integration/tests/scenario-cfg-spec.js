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

// Import MEEP Contstants
import * as meep from '../../../../js-apps/meep-frontend/src/js/meep-constants';

// Import element utils
import {
  // Field Names
  FIELD_TYPE,
  FIELD_PARENT,
  FIELD_NAME,
  FIELD_IMAGE,
  FIELD_PORT,
  FIELD_PROTOCOL,
  FIELD_GROUP,
  FIELD_SVC_MAP,
  FIELD_GPU_COUNT,
  FIELD_GPU_TYPE,
  FIELD_ENV_VAR,
  FIELD_CMD,
  FIELD_CMD_ARGS,
  FIELD_EXT_PORT,
  FIELD_IS_EXTERNAL,
  FIELD_CHART_ENABLED,
  FIELD_CHART_LOC,
  FIELD_CHART_VAL,
  FIELD_CHART_GROUP,
  FIELD_INT_DOM_LATENCY,
  FIELD_INT_DOM_LATENCY_VAR,
  FIELD_INT_DOM_THROUGPUT,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGPUT,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INT_EDGE_LATENCY,
  FIELD_INT_EDGE_LATENCY_VAR,
  FIELD_INT_EDGE_THROUGPUT,
  FIELD_INT_EDGE_PKT_LOSS,
  FIELD_INT_FOG_LATENCY,
  FIELD_INT_FOG_LATENCY_VAR,
  FIELD_INT_FOG_THROUGPUT,
  FIELD_INT_FOG_PKT_LOSS,
  FIELD_EDGE_FOG_LATENCY,
  FIELD_EDGE_FOG_LATENCY_VAR,
  FIELD_EDGE_FOG_THROUGPUT,
  FIELD_EDGE_FOG_PKT_LOSS,
  FIELD_TERM_LINK_LATENCY,
  FIELD_TERM_LINK_LATENCY_VAR,
  FIELD_TERM_LINK_THROUGPUT,
  FIELD_TERM_LINK_PKT_LOSS,

  getElemFieldVal,
} from '../../../../js-apps/meep-frontend/src/js/util/elem-utils';

// Import Test utility functions
import { selector, click, type, select, verify, verifyEnabled, verifyForm } from '../util/util';

// Scenario Configuration Tests
describe('Scenario Configuration', function() {

  // Test Variables
  let defaultScenario = 'None';
  let dummyScenario = 'dummy-scenario';

  // Test Setup
  beforeEach(() => {
    var meepUrl = Cypress.env('meep_url');
    if (meepUrl == null) {
      meepUrl = 'http://127.0.0.1:30000';
    }

    cy.viewport(1920, 1080);
    cy.visit(meepUrl);
  });

  it('Create, Save, & Delete Scenario', function() {
    // Go to configuration page
    cy.log('Go to configuration page');
    click(meep.MEEP_TAB_CFG);
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
    verifyEnabled(meep.CFG_BTN_NEW_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_OPEN_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_SAVE_SCENARIO, false);
    verifyEnabled(meep.CFG_BTN_DEL_SCENARIO, false);
    verifyEnabled(meep.CFG_BTN_IMP_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_EXP_SCENARIO, false);

    // Make sure scenario does not exist in MEEP
    cy.log('Make sure scenario does not exist in MEEP');
    click(meep.CFG_BTN_OPEN_SCENARIO);
    cy.wait(50);
    verify(meep.MEEP_DLG_OPEN_SCENARIO, 'not.contain', dummyScenario);
    click(meep.MEEP_DLG_OPEN_SCENARIO, 'Cancel');

    // Create new scenario
    cy.log('Create new scenario: ' + dummyScenario);
    click(meep.CFG_BTN_NEW_SCENARIO);
    type(meep.MEEP_DLG_NEW_SCENARIO_NAME, dummyScenario);
    click(meep.MEEP_DLG_NEW_SCENARIO, 'Ok');
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', dummyScenario);
    verifyEnabled(meep.CFG_BTN_NEW_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_OPEN_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_SAVE_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_DEL_SCENARIO, false);
    verifyEnabled(meep.CFG_BTN_IMP_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_EXP_SCENARIO, true);

    // Make sure scenario is not yet stored in MEEP
    cy.log('Make sure scenario is not yet stored in MEEP');
    click(meep.CFG_BTN_OPEN_SCENARIO);
    cy.wait(50);
    verify(meep.MEEP_DLG_OPEN_SCENARIO, 'not.contain', dummyScenario);
    click(meep.MEEP_DLG_OPEN_SCENARIO, 'Cancel');

    // Save scenario
    cy.log('Save scenario: ' + dummyScenario);
    click(meep.CFG_BTN_SAVE_SCENARIO);
    click(meep.MEEP_DLG_SAVE_SCENARIO, 'Ok');
    verifyEnabled(meep.CFG_BTN_NEW_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_OPEN_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_SAVE_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_DEL_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_IMP_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_EXP_SCENARIO, true);

    // Make sure scenario is stored in MEEP
    cy.log('Make sure scenario is stored in MEEP');
    click(meep.CFG_BTN_OPEN_SCENARIO);
    cy.wait(50);
    verify(meep.MEEP_DLG_OPEN_SCENARIO, 'contain', dummyScenario);
    click(meep.MEEP_DLG_OPEN_SCENARIO, 'Cancel');

    // Delete scenario
    cy.log('Delete scenario: ' + dummyScenario);
    click(meep.CFG_BTN_DEL_SCENARIO);
    click(meep.MEEP_DLG_DEL_SCENARIO, 'Ok');
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', defaultScenario);
    verifyEnabled(meep.CFG_BTN_NEW_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_OPEN_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_SAVE_SCENARIO, false);
    verifyEnabled(meep.CFG_BTN_DEL_SCENARIO, false);
    verifyEnabled(meep.CFG_BTN_IMP_SCENARIO, true);
    verifyEnabled(meep.CFG_BTN_EXP_SCENARIO, false);
  });

  it('Create Full Scenario', function() {
    let operatorName = 'operator1';
    let zoneName = 'zone1';
    let edgeName = 'edge1';
    let edgeAppName = 'edge1-app1';
    let poaName = 'poa1';
    let fogName = 'fog1';
    let fogAppName = 'fog1-app1';
    let ueName = 'ue1';
    let ueAppName = 'ue1-app1';
    let ueAppExtName = 'ue1-app1-ext';
    let cloudName = 'cloud1';
    let cloudAppName = 'cloud1-app1';

    // Create new dummy scenario
    cy.log('Create & validate new scenario: ' + dummyScenario);
    createNewScenario(dummyScenario);
    cy.wait(50);
    validateScenario(dummyScenario);

    // Close new element creation
    cy.log('Close new element creation');
    click(meep.CFG_BTN_NEW_ELEM);
    verifyForm(meep.CFG_ELEM_TYPE, true);
    verifyForm(meep.CFG_ELEM_NAME, true);
    verifyEnabled(meep.MEEP_BTN_CANCEL, true);
    verifyEnabled(meep.MEEP_BTN_APPLY, true);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, false);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    click(meep.MEEP_BTN_CANCEL);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);

    // Domain
    cy.log('Add new domain and verify default & configured settings: ' + operatorName);
    addDomain(operatorName, dummyScenario);
    validateDomain(operatorName, dummyScenario);

    // Zone
    cy.log('Add new zone and verify default & configured settings: ' + zoneName);
    addZone(zoneName, operatorName);
    validateZone(zoneName, operatorName);

    // Edge
    cy.log('Add new edge and verify default & configured settings: ' + edgeName);
    addEdge(edgeName, zoneName);
    validateEdge(edgeName, zoneName);

    // Edge App
    cy.log('Add new edge app and verify default & configured settings: ' + edgeAppName);
    addEdgeApp(edgeAppName, edgeName);
    validateEdgeApp(edgeAppName, edgeName);

    // POA
    cy.log('Add new poa and verify default & configured settings: ' + poaName);
    addPoa(poaName, zoneName);
    validatePoa(poaName, zoneName);

    // Fog
    cy.log('Add new fog and verify default & configured settings: ' + fogName);
    addFog(fogName, poaName);
    validateFog(fogName, poaName);

    // Fog App
    cy.log('Add new fog app and verify default & configured settings: ' + fogAppName);
    addFogApp(fogAppName, fogName);
    validateFogApp(fogAppName, fogName);

    // UE
    cy.log('Add new UE and verify default & configured settings: ' + ueName);
    addUe(ueName, poaName);
    validateUe(ueName, poaName);

    // UE App
    cy.log('Add new UE app and verify default & configured settings: ' + ueAppName);
    addUeApp(ueAppName, ueName);
    validateUeApp(ueAppName, ueName);

    // UE App (External)
    cy.log('Add new UE app (external) and verify default & configured settings: ' + ueAppExtName);
    addUeAppExt(ueAppExtName, ueName);
    validateUeAppExt(ueAppExtName, ueName);

    // Cloud
    cy.log('Add new cloud and verify default & configured settings: ' + cloudName);
    addCloud(cloudName, dummyScenario);
    validateCloud(cloudName, dummyScenario);

    // Cloud App
    cy.log('Add new cloud app and verify default & configured settings: ' + cloudAppName);
    addCloudApp(cloudAppName, cloudName);
    validateCloudApp(cloudAppName, cloudName);

    // Save scenario
    cy.log('Save scenario: ' + dummyScenario);
    click(meep.CFG_BTN_SAVE_SCENARIO);
    click(meep.MEEP_DLG_SAVE_SCENARIO, 'Ok');

    // Open scenario stored in MEEP
    cy.log('Open scenario is stored in MEEP');
    click(meep.CFG_BTN_OPEN_SCENARIO);
    cy.wait(50);
    verify(meep.MEEP_DLG_OPEN_SCENARIO, 'contain', dummyScenario);
    select(meep.MEEP_DLG_OPEN_SCENARIO_SELECT, dummyScenario);
    click(meep.MEEP_DLG_OPEN_SCENARIO, 'Ok');
    cy.wait(50);

    // Validate Loaded scenario entries match saved scenario values
    cy.log('Validate Loaded scenario entries match saved scenario values');
    validateScenario(dummyScenario);
    validateDomain(operatorName, dummyScenario);
    validateZone(zoneName, operatorName);
    validateEdge(edgeName, zoneName);
    validateEdgeApp(edgeAppName, edgeName);
    validatePoa(poaName, zoneName);
    validateFog(fogName, poaName);
    validateFogApp(fogAppName, fogName);
    validateUe(ueName, poaName);
    validateUeApp(ueAppName, ueName);
    validateUeAppExt(ueAppExtName, ueName);
    validateCloud(cloudName, dummyScenario);
    validateCloudApp(cloudAppName, cloudName);

    // Delete scenario
    cy.log('Delete scenario: ' + dummyScenario);
    click(meep.CFG_BTN_DEL_SCENARIO);
    click(meep.MEEP_DLG_DEL_SCENARIO, 'Ok');
  });

  // Retrieve Element entry from Application table
  function getEntry(entries, name) {
    if (entries) {
      for (var i = 0; i < entries.length; i++) {
        if (getElemFieldVal(entries[i], FIELD_NAME) == name) {
          return entries[i];
        }
      }
    }
    return null;
  }

  // ==============================
  // SCENARIO
  // ==============================

  function createNewScenario(name) {
    // Go to configuration page
    click(meep.MEEP_TAB_CFG);

    // Create new scenario
    click(meep.CFG_BTN_NEW_SCENARIO);
    type(meep.MEEP_DLG_NEW_SCENARIO_NAME, name);
    click(meep.MEEP_DLG_NEW_SCENARIO, 'Ok');
    verify(meep.MEEP_LBL_SCENARIO_NAME, 'contain', name);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateScenario(name) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_SCENARIO);
    });
  }

  // ==============================
  // DOMAIN
  // ==============================

  let interZoneLatency = '12';
  let interZoneLatencyVar = '4';
  let interZonePktLoss = '2';
  let interZoneThroughput = '2000';

  function addDomain(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_OPERATOR);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_THROUGHPUT, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_INTER_ZONE));
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_LATENCY, interZoneLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, interZoneLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, interZonePktLoss);
    type(meep.CFG_ELEM_THROUGHPUT, interZoneThroughput);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateDomain(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_OPERATOR);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY), interZoneLatency);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY_VAR), interZoneLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_PKT_LOSS), interZonePktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGPUT), interZoneThroughput);
    });
  }

  // ==============================
  // ZONE
  // ==============================

  function addZone(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_ZONE);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateZone(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_ZONE);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
    });
  }

  // ==============================
  // EDGE
  // ==============================

  function addEdge(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_EDGE);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateEdge(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_EDGE);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
    });
  }

  // ==============================
  // EDGE APP
  // ==============================

  let edgeAppImg = 'nginx';
  let edgeAppPort = '1234';
  let edgeAppExtPort = '32323';
  let edgeAppProt = 'TCP';
  let edgeAppGroup = 'edge-svc';
  let edgeAppGpuCount = '1';
  let edgeAppGpuType = 'NVIDIA';
  let edgeAppEnv = 'ENV_VAR=my-env-var';
  let edgeAppCmd = '/bin/bash';
  let edgeAppArgs = '-c, export;';

  function addEdgeApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_EDGE_APP);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, edgeAppImg);
    type(meep.CFG_ELEM_PORT, edgeAppPort);
    type(meep.CFG_ELEM_EXT_PORT, edgeAppExtPort);
    select(meep.CFG_ELEM_PROT, edgeAppProt);
    type(meep.CFG_ELEM_GROUP, edgeAppGroup);
    type(meep.CFG_ELEM_GPU_COUNT, edgeAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, edgeAppGpuType);
    type(meep.CFG_ELEM_ENV, edgeAppEnv);
    type(meep.CFG_ELEM_CMD, edgeAppCmd);
    type(meep.CFG_ELEM_ARGS, edgeAppArgs);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateEdgeApp(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_EDGE_APP);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_IMAGE), edgeAppImg);
      assert.equal(getElemFieldVal(entry, FIELD_PORT), edgeAppPort);
      assert.equal(getElemFieldVal(entry, FIELD_EXT_PORT), edgeAppExtPort);
      assert.equal(getElemFieldVal(entry, FIELD_PROTOCOL), edgeAppProt);
      assert.equal(getElemFieldVal(entry, FIELD_GROUP), edgeAppGroup);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_COUNT), edgeAppGpuCount);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_TYPE), edgeAppGpuType);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), edgeAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), edgeAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), edgeAppArgs);
    });
  }

  // ==============================
  // POA
  // ==============================
    
  let linkLatency = '2';
  let linkLatencyVar = '3';
  let linkPktLoss = '4';
  let linkThroughput = '5';

  function addPoa(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_POA);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_TERMINAL_LINK));
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_LATENCY, linkLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, linkLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, linkPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT, linkThroughput);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validatePoa(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_POA);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY), linkLatency);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY_VAR), linkLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_PKT_LOSS), linkPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGPUT), linkThroughput);
    });
  }

  // ==============================
  // FOG
  // ==============================

  function addFog(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_FOG);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateFog(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_FOG);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
    });
  }

  // ==============================
  // FOG APP
  // ==============================

  let fogAppImg = 'nginx';
  let fogAppPort = '5678';
  let fogAppExtPort = '31313';
  let fogAppProt = 'UDP';
  let fogAppGroup = 'fog-svc';
  let fogAppGpuCount = '2';
  let fogAppGpuType = 'NVIDIA';
  let fogAppEnv = 'ENV_VAR=my-env-var';
  let fogAppCmd = '/bin/bash';
  let fogAppArgs = '-c, export;';

  function addFogApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_EDGE_APP);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, fogAppImg);
    type(meep.CFG_ELEM_PORT, fogAppPort);
    type(meep.CFG_ELEM_EXT_PORT, fogAppExtPort);
    select(meep.CFG_ELEM_PROT, fogAppProt);
    type(meep.CFG_ELEM_GROUP, fogAppGroup);
    type(meep.CFG_ELEM_GPU_COUNT, fogAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, fogAppGpuType);
    type(meep.CFG_ELEM_ENV, fogAppEnv);
    type(meep.CFG_ELEM_CMD, fogAppCmd);
    type(meep.CFG_ELEM_ARGS, fogAppArgs);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateFogApp(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_EDGE_APP);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_IMAGE), fogAppImg);
      assert.equal(getElemFieldVal(entry, FIELD_PORT), fogAppPort);
      assert.equal(getElemFieldVal(entry, FIELD_EXT_PORT), fogAppExtPort);
      assert.equal(getElemFieldVal(entry, FIELD_PROTOCOL), fogAppProt);
      assert.equal(getElemFieldVal(entry, FIELD_GROUP), fogAppGroup);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_COUNT), fogAppGpuCount);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_TYPE), fogAppGpuType);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), fogAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), fogAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), fogAppArgs);
    });
  }

  // ==============================
  // UE
  // ==============================

  function addUe(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_UE);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateUe(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_UE);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
    });
  }

  // ==============================
  // UE APP
  // ==============================

  let ueAppImg = 'nginx';
  let ueAppGpuCount = '3';
  let ueAppGpuType = 'NVIDIA';
  let ueAppEnv = 'ENV_VAR=my-env-var';
  let ueAppCmd = '/bin/bash';
  let ueAppArgs = '-c, export;';

  // Add new ue app element
  function addUeApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_UE_APP);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, ueAppImg);
    type(meep.CFG_ELEM_GPU_COUNT, ueAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, ueAppGpuType);
    type(meep.CFG_ELEM_ENV, ueAppEnv);
    type(meep.CFG_ELEM_CMD, ueAppCmd);
    type(meep.CFG_ELEM_ARGS, ueAppArgs);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateUeApp(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_UE_APP);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_IMAGE), ueAppImg);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_COUNT), ueAppGpuCount);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_TYPE), ueAppGpuType);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), ueAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), ueAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), ueAppArgs);
    });
  }

  // ==============================
  // EXTERNAL UE APP
  // ==============================

  function addUeAppExt(name, parent) {
    // TODO -- Add test
  }

  function validateUeAppExt(name, parent) {
    // cy.window().then((win) => {
    //     var entry = getEntry(win.meepStore.getState().cfg.table.entries, name)
    //     assert.isNotNull(entry);
    //     assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_UE_APP);
    //     assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
    //     assert.equal(getElemFieldVal(entry, FIELD_IMAGE), ueAppImg);
    //     assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), ueAppEnv);
    //     assert.equal(getElemFieldVal(entry, FIELD_CMD), ueAppCmd);
    //     assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), ueAppArgs);
    // })
  }

  // ==============================
  // DISTANT CLOUD
  // ==============================

  function addCloud(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_DC);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }

  function validateCloud(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_DC);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
    });
  }

  // ==============================
  // DISTANT CLOUD APP
  // ==============================

  let cloudAppImg = 'nginx';
  let cloudAppPort = '9101';
  let cloudAppExtPort = '30303';
  let cloudAppProt = 'TCP';
  let cloudAppGpuCount = '4';
  let cloudAppGpuType = 'NVIDIA';
  let cloudAppEnv = 'ENV_VAR=my-env-var';
  let cloudAppCmd = '/bin/bash';
  let cloudAppArgs = '-c, export;';

  function addCloudApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_CLOUD_APP);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, cloudAppImg);
    type(meep.CFG_ELEM_PORT, cloudAppPort);
    type(meep.CFG_ELEM_EXT_PORT, cloudAppExtPort);
    select(meep.CFG_ELEM_PROT, cloudAppProt);
    type(meep.CFG_ELEM_GPU_COUNT, cloudAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, cloudAppGpuType);
    type(meep.CFG_ELEM_ENV, cloudAppEnv);
    type(meep.CFG_ELEM_CMD, cloudAppCmd);
    type(meep.CFG_ELEM_ARGS, cloudAppArgs);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
  }
    
  function validateCloudApp(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_CLOUD_APP);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_IMAGE), cloudAppImg);
      assert.equal(getElemFieldVal(entry, FIELD_PORT), cloudAppPort);
      assert.equal(getElemFieldVal(entry, FIELD_EXT_PORT), cloudAppExtPort);
      assert.equal(getElemFieldVal(entry, FIELD_PROTOCOL), cloudAppProt);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_COUNT), cloudAppGpuCount);
      assert.equal(getElemFieldVal(entry, FIELD_GPU_TYPE), cloudAppGpuType);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), cloudAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), cloudAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), cloudAppArgs);
    });
  }

});


