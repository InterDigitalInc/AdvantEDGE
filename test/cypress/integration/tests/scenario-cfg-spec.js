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
  FIELD_PLACEMENT_ID,
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
  FIELD_INT_DOM_LATENCY_DIST,
  FIELD_INT_DOM_THROUGHPUT_DL,
  FIELD_INT_DOM_THROUGHPUT_UL,
  FIELD_INT_DOM_PKT_LOSS,
  FIELD_INT_ZONE_LATENCY,
  FIELD_INT_ZONE_LATENCY_VAR,
  FIELD_INT_ZONE_THROUGHPUT_DL,
  FIELD_INT_ZONE_THROUGHPUT_UL,
  FIELD_INT_ZONE_PKT_LOSS,
  FIELD_INTRA_ZONE_LATENCY,
  FIELD_INTRA_ZONE_LATENCY_VAR,
  FIELD_INTRA_ZONE_THROUGHPUT_DL,
  FIELD_INTRA_ZONE_THROUGHPUT_UL,
  FIELD_INTRA_ZONE_PKT_LOSS,
  FIELD_TERM_LINK_LATENCY,
  FIELD_TERM_LINK_LATENCY_VAR,
  FIELD_TERM_LINK_THROUGHPUT_DL,
  FIELD_TERM_LINK_THROUGHPUT_UL,
  FIELD_TERM_LINK_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGHPUT_DL,
  FIELD_LINK_THROUGHPUT_UL,
  FIELD_LINK_PKT_LOSS,
  FIELD_APP_LATENCY,
  FIELD_APP_LATENCY_VAR,
  FIELD_APP_THROUGHPUT_DL,
  FIELD_APP_THROUGHPUT_UL,
  FIELD_APP_PKT_LOSS,
  FIELD_MNC,
  FIELD_MCC,
  FIELD_MAC_ID,
  FIELD_UE_MAC_ID,
  FIELD_DEFAULT_CELL_ID,
  FIELD_CELL_ID,
  FIELD_NR_CELL_ID,

  getElemFieldVal,
  FIELD_CPU_MIN,
  FIELD_CPU_MAX,
  FIELD_MEMORY_MIN,
  FIELD_MEMORY_MAX,

  FIELD_META_DISPLAY_MAP_COLOR,
  FIELD_GEO_LOCATION,
  FIELD_GEO_RADIUS,
  FIELD_GEO_PATH,
  FIELD_GEO_VELOCITY,
  FIELD_GEO_EOP_MODE,
  FIELD_CONNECTED,
  FIELD_WIRELESS_TYPE,

} from '../../../../js-apps/meep-frontend/src/js/util/elem-utils';

// Import Test utility functions
import { selector, click, type, select, verify, verifyEnabled, verifyForm } from '../util/util';

// Scenario Configuration Tests
describe('Scenario Configuration', function () {

  // Test Variables
  let defaultScenario = 'None';
  let dummyScenario = 'dummy-scenario21';

  // Test Setup
  beforeEach(() => {
    var meepUrl = Cypress.env('meep_url');
    if (meepUrl == null) {
      meepUrl = 'http://127.0.0.1';
    }

    cy.viewport(1920, 1080);
    cy.visit(meepUrl);
  });

  it('Create, Save, & Delete Scenario', function () {
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

  it('Create Full Scenario', function () {
    let operatorName = 'operator1';
    let operatorCellName = 'operator-cell1';
    let zoneName = 'zone1';
    let edgeName = 'edge1';
    let edgeAppName = 'edge1-app1';
    let poaName = 'poa1';
    let poa4GName = 'poa-4g1';
    let poa5GName = 'poa-5g1';
    let poaWifiName = 'poa-wifi1';
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
    verifyEnabled(meep.MEEP_BTN_APPLY, false);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, false);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
    click(meep.MEEP_BTN_CANCEL);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);

    // Domain
    cy.log('Add new domain and verify default & configured settings: ' + operatorName);
    addDomain(operatorName, dummyScenario);
    validateDomain(operatorName, dummyScenario);

    // Domain Cell
    cy.log('Add new domain cell and verify default & configured settings: ' + operatorCellName);
    addDomainCell(operatorCellName, dummyScenario);
    validateDomainCell(operatorCellName, dummyScenario);

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

    // POA 4G
    cy.log('Add new poa 4G and verify default & configured settings: ' + poa4GName);
    addPoa4G(poa4GName, zoneName);
    validatePoa4G(poa4GName, zoneName);

    // POA 5G
    cy.log('Add new poa 5G and verify default & configured settings: ' + poa5GName);
    addPoa5G(poa5GName, zoneName);
    validatePoa5G(poa5GName, zoneName);

    // POA WIFI
    cy.log('Add new poa wifi and verify default & configured settings: ' + poaWifiName);
    addPoaWifi(poaWifiName, zoneName);
    validatePoaWifi(poaWifiName, zoneName);

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
    validateDomainCell(operatorCellName, dummyScenario);
    validateZone(zoneName, operatorName);
    validateEdge(edgeName, zoneName);
    validateEdgeApp(edgeAppName, edgeName);
    validatePoa(poaName, zoneName);
    validatePoa4G(poa4GName, zoneName);
    validatePoa5G(poa5GName, zoneName);
    validatePoaWifi(poaWifiName, zoneName);
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
      return entries[name] ? entries[name] : null;
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
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
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
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_OPERATOR_GENERIC);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_INTER_ZONE));

    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_LATENCY, interZoneLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, interZoneLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, interZonePktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, interZoneThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, interZoneThroughput-1);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
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
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGHPUT_DL), interZoneThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGHPUT_UL), interZoneThroughput-1);
    });
  }

  // ==============================
  // DOMAIN CELL
  // ==============================

  let interZoneLatency2 = '13';
  let interZoneLatencyVar2 = '4';
  let interZonePktLoss2 = '3';
  let interZoneThroughput2 = '2001';
  let mcc = '002';
  let mnc = '001';
  let defaultCellId = 'ABCDEF1';

  function addDomainCell(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_OPERATOR_CELL);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_INTER_ZONE));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_INTER_ZONE));
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_LATENCY, interZoneLatency2);
    type(meep.CFG_ELEM_LATENCY_VAR, interZoneLatencyVar2);
    type(meep.CFG_ELEM_PKT_LOSS, interZonePktLoss2);
    type(meep.CFG_ELEM_THROUGHPUT_DL, interZoneThroughput2);
    type(meep.CFG_ELEM_THROUGHPUT_UL, interZoneThroughput2-1);
    type(meep.CFG_ELEM_MCC, mcc);
    type(meep.CFG_ELEM_MNC, mnc);
    type(meep.CFG_ELEM_DEFAULT_CELL_ID, defaultCellId);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validateDomainCell(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_OPERATOR_CELL);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY), interZoneLatency2);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_LATENCY_VAR), interZoneLatencyVar2);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_PKT_LOSS), interZonePktLoss2);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGHPUT_DL), interZoneThroughput2);
      assert.equal(getElemFieldVal(entry, FIELD_INT_ZONE_THROUGHPUT_UL), interZoneThroughput2-1);
      assert.equal(getElemFieldVal(entry, FIELD_MCC), mcc);
      assert.equal(getElemFieldVal(entry, FIELD_MNC), mnc);
      assert.equal(getElemFieldVal(entry, FIELD_DEFAULT_CELL_ID), defaultCellId);
    });
  }


  // ==============================
  // ZONE
  // ==============================
  let intraZoneLatency = '2';
  let intraZoneLatencyVar = '3';
  let intraZonePktLoss = '4';
  let intraZoneThroughput = '5';
  let zoneColor = '#123DEF'

  function addZone(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_ZONE);
    select(meep.CFG_ELEM_PARENT, parent);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_INTRA_ZONE));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_INTRA_ZONE));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_INTRA_ZONE));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_INTRA_ZONE));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_INTRA_ZONE));
    //error value section
    type(meep.CFG_ELEM_META_DISPLAY_MAP_COLOR, 'red');
    click(meep.MEEP_BTN_APPLY);
    cy.contains('1 fields in error')
    //valid value section
    type(meep.CFG_ELEM_LATENCY, intraZoneLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, intraZoneLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, intraZonePktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, intraZoneThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, intraZoneThroughput-1);
    type(meep.CFG_ELEM_META_DISPLAY_MAP_COLOR, zoneColor);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validateZone(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_ZONE);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_LATENCY), intraZoneLatency);
      assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_LATENCY_VAR), intraZoneLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_PKT_LOSS), intraZonePktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_THROUGHPUT_DL), intraZoneThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_INTRA_ZONE_THROUGHPUT_UL), intraZoneThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_META_DISPLAY_MAP_COLOR), zoneColor);

    });
  }

  // ==============================
  // EDGE
  // ==============================

  //valid for every physical locations in the other test cases too
  let linkLatency = '2';
  let linkLatencyVar = '3';
  let linkPktLoss = '4';
  let linkThroughput = '5';
  let linkLocationCoordinates = '[7.419344,43.72764]';

  function addEdge(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_EDGE);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_LINK));
    verifyForm(meep.CFG_ELEM_CONNECTED, true, 'have.value', String(meep.OPT_CONNECTED.value));
    verifyForm(meep.CFG_ELEM_WIRELESS, true, 'have.value', String(meep.OPT_WIRED.value));
    type(meep.CFG_ELEM_LATENCY, linkLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, linkLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, linkPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, linkThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, linkThroughput-1);
    type(meep.CFG_ELEM_GEO_LOCATION, linkLocationCoordinates);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validateEdge(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_EDGE);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY), linkLatency);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY_VAR), linkLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_PKT_LOSS), linkPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_DL), linkThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_UL), linkThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), linkLocationCoordinates);
    });
  }

  // ==============================
  // EDGE APP
  // ==============================

  //valid for every app in the other test cases too
  let appLatency = '2';
  let appLatencyVar = '4';
  let appPktLoss = '5';
  let appThroughput = '6';

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
  let edgeAppPlacementId = 'node1';
  let edgeAppCpuMin = '0.5';
  let edgeAppCpuMax = '1';
  let edgeAppMemoryMin = '100';
  let edgeAppMemoryMax = '200';

  function addEdgeApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_EDGE_APP);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_APP));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_APP));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_APP));
    type(meep.CFG_ELEM_LATENCY, appLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, appLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, appPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, appThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, appThroughput-1);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, edgeAppImg);
    type(meep.CFG_ELEM_PORT, edgeAppPort);
    type(meep.CFG_ELEM_EXT_PORT, edgeAppExtPort);
    select(meep.CFG_ELEM_PROT, edgeAppProt);
    type(meep.CFG_ELEM_GROUP, edgeAppGroup);
    type(meep.CFG_ELEM_GPU_COUNT, edgeAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, edgeAppGpuType);
    type(meep.CFG_ELEM_CPU_MIN, edgeAppCpuMin);
    type(meep.CFG_ELEM_CPU_MAX, edgeAppCpuMax);
    type(meep.CFG_ELEM_MEMORY_MIN, edgeAppMemoryMin);
    type(meep.CFG_ELEM_MEMORY_MAX, edgeAppMemoryMax);
    type(meep.CFG_ELEM_ENV, edgeAppEnv);
    type(meep.CFG_ELEM_CMD, edgeAppCmd);
    type(meep.CFG_ELEM_ARGS, edgeAppArgs);
    type(meep.CFG_ELEM_PLACEMENT_ID, edgeAppPlacementId);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
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
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MIN), edgeAppCpuMin);
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MAX), edgeAppCpuMax);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MIN), edgeAppMemoryMin);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MAX), edgeAppMemoryMax);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), edgeAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), edgeAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), edgeAppArgs);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY), appLatency);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY_VAR), appLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_APP_PKT_LOSS), appPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_DL), appThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_UL), appThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_PLACEMENT_ID), edgeAppPlacementId);
    });
  }

  // ==============================
  // POA
  // ==============================

  let termLinkLatency = '2';
  let termLinkLatencyVar = '3';
  let termLinkPktLoss = '4';
  let termLinkThroughput = '5';
  let termLinkLocationCoordinates = '[7.419344,43.72764]';
  let termLinkRadius = '10';

  function addPoa(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_POA_GENERIC);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_TERMINAL_LINK));
    type(meep.CFG_ELEM_LATENCY, termLinkLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, termLinkLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, termLinkPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, termLinkThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, termLinkThroughput-1);
    type(meep.CFG_ELEM_GEO_LOCATION, termLinkLocationCoordinates);
    type(meep.CFG_ELEM_GEO_RADIUS, termLinkRadius);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validatePoa(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_POA);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY), termLinkLatency);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY_VAR), termLinkLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_PKT_LOSS), termLinkPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_DL), termLinkThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_UL), termLinkThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), termLinkLocationCoordinates);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_RADIUS), termLinkRadius);
    });
  }

  // ==============================
  // POA-4G
  // ==============================

  let termLinkLatency2 = '2';
  let termLinkLatencyVar2 = '3';
  let termLinkPktLoss2 = '4';
  let termLinkThroughput2 = '5';
  let termLinkLocationCoordinates2 = '[7.419344,43.72764]';
  let termLinkRadius2 = '10';
  let cellId = '1234567';


  function addPoa4G(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_POA_4G);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_TERMINAL_LINK));
    type(meep.CFG_ELEM_LATENCY, termLinkLatency2);
    type(meep.CFG_ELEM_LATENCY_VAR, termLinkLatencyVar2);
    type(meep.CFG_ELEM_PKT_LOSS, termLinkPktLoss2);
    type(meep.CFG_ELEM_THROUGHPUT_DL, termLinkThroughput2);
    type(meep.CFG_ELEM_THROUGHPUT_UL, termLinkThroughput2-1);
    type(meep.CFG_ELEM_GEO_LOCATION, termLinkLocationCoordinates2);
    type(meep.CFG_ELEM_GEO_RADIUS, termLinkRadius2);
    type(meep.CFG_ELEM_CELL_ID, cellId);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validatePoa4G(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_POA_4G);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY), termLinkLatency2);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY_VAR), termLinkLatencyVar2);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_PKT_LOSS), termLinkPktLoss2);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_DL), termLinkThroughput2);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_UL), termLinkThroughput2-1);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), termLinkLocationCoordinates2);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_RADIUS), termLinkRadius2);
      assert.equal(getElemFieldVal(entry, FIELD_CELL_ID), cellId);
    });
  }

  // ==============================
  // POA-5G
  // ==============================

  let termLinkLatency3 = '2';
  let termLinkLatencyVar3 = '3';
  let termLinkPktLoss3 = '4';
  let termLinkThroughput3 = '5';
  let termLinkLocationCoordinates3 = '[7.419344,43.72764]';
  let termLinkRadius3 = '10';
  let nrCellId = '3456789';

  function addPoa5G(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_POA_5G);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_TERMINAL_LINK));
    type(meep.CFG_ELEM_LATENCY, termLinkLatency3);
    type(meep.CFG_ELEM_LATENCY_VAR, termLinkLatencyVar3);
    type(meep.CFG_ELEM_PKT_LOSS, termLinkPktLoss3);
    type(meep.CFG_ELEM_THROUGHPUT_DL, termLinkThroughput3);
    type(meep.CFG_ELEM_THROUGHPUT_UL, termLinkThroughput3-1);
    type(meep.CFG_ELEM_GEO_LOCATION, termLinkLocationCoordinates3);
    type(meep.CFG_ELEM_GEO_RADIUS, termLinkRadius3);
    type(meep.CFG_ELEM_NR_CELL_ID, nrCellId);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validatePoa5G(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_POA_5G);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY), termLinkLatency3);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY_VAR), termLinkLatencyVar3);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_PKT_LOSS), termLinkPktLoss3);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_DL), termLinkThroughput3);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_UL), termLinkThroughput3-1);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), termLinkLocationCoordinates3);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_RADIUS), termLinkRadius3);
      assert.equal(getElemFieldVal(entry, FIELD_NR_CELL_ID), nrCellId);
    });
  }

  // ==============================
  // POA-WIFI
  // ==============================

  let termLinkLatency4 = '2';
  let termLinkLatencyVar4 = '3';
  let termLinkPktLoss4 = '4';
  let termLinkThroughput4 = '5';
  let termLinkLocationCoordinates4 = '[7.419344,43.72764]';
  let termLinkRadius4 = '10';
  let macId = '112233445566';

  function addPoaWifi(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_POA_WIFI);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_TERMINAL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_TERMINAL_LINK));
    type(meep.CFG_ELEM_LATENCY, termLinkLatency4);
    type(meep.CFG_ELEM_LATENCY_VAR, termLinkLatencyVar4);
    type(meep.CFG_ELEM_PKT_LOSS, termLinkPktLoss4);
    type(meep.CFG_ELEM_THROUGHPUT_DL, termLinkThroughput4);
    type(meep.CFG_ELEM_THROUGHPUT_UL, termLinkThroughput4-1);
    type(meep.CFG_ELEM_GEO_LOCATION, termLinkLocationCoordinates4);
    type(meep.CFG_ELEM_GEO_RADIUS, termLinkRadius4);
    type(meep.CFG_ELEM_MAC_ID, macId);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validatePoaWifi(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_POA_WIFI);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY), termLinkLatency4);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_LATENCY_VAR), termLinkLatencyVar4);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_PKT_LOSS), termLinkPktLoss4);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_DL), termLinkThroughput4);
      assert.equal(getElemFieldVal(entry, FIELD_TERM_LINK_THROUGHPUT_UL), termLinkThroughput4-1);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), termLinkLocationCoordinates4);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_RADIUS), termLinkRadius4);
      assert.equal(getElemFieldVal(entry, FIELD_MAC_ID), macId);
    });
  }

  // ==============================
  // FOG
  // ==============================

  function addFog(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_FOG);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_LINK));
    verifyForm(meep.CFG_ELEM_CONNECTED, true, 'have.value', String(meep.OPT_CONNECTED.value));
    verifyForm(meep.CFG_ELEM_WIRELESS, true, 'have.value', String(meep.OPT_WIRED.value));
    type(meep.CFG_ELEM_LATENCY, linkLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, linkLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, linkPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, linkThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, linkThroughput-1);
    type(meep.CFG_ELEM_GEO_LOCATION, linkLocationCoordinates);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validateFog(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_FOG);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY), linkLatency);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY_VAR), linkLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_PKT_LOSS), linkPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_DL), linkThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_UL), linkThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), linkLocationCoordinates);
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
  let fogAppPlacementId = 'node2';
  let fogAppCpuMin = '0.5';
  let fogAppCpuMax = '1';
  let fogAppMemoryMin = '100';
  let fogAppMemoryMax = '200';

  function addFogApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_EDGE_APP);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_APP));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_APP));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_APP));
    type(meep.CFG_ELEM_LATENCY, appLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, appLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, appPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, appThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, appThroughput-1);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, fogAppImg);
    type(meep.CFG_ELEM_PORT, fogAppPort);
    type(meep.CFG_ELEM_EXT_PORT, fogAppExtPort);
    select(meep.CFG_ELEM_PROT, fogAppProt);
    type(meep.CFG_ELEM_GROUP, fogAppGroup);
    type(meep.CFG_ELEM_GPU_COUNT, fogAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, fogAppGpuType);
    type(meep.CFG_ELEM_CPU_MIN, fogAppCpuMin);
    type(meep.CFG_ELEM_CPU_MAX, fogAppCpuMax);
    type(meep.CFG_ELEM_MEMORY_MIN, fogAppMemoryMin);
    type(meep.CFG_ELEM_MEMORY_MAX, fogAppMemoryMax);
    type(meep.CFG_ELEM_ENV, fogAppEnv);
    type(meep.CFG_ELEM_CMD, fogAppCmd);
    type(meep.CFG_ELEM_ARGS, fogAppArgs);
    type(meep.CFG_ELEM_PLACEMENT_ID, fogAppPlacementId);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
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
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MIN), fogAppCpuMin);
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MAX), fogAppCpuMax);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MIN), fogAppMemoryMin);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MAX), fogAppMemoryMax);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), fogAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), fogAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), fogAppArgs);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY), appLatency);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY_VAR), appLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_APP_PKT_LOSS), appPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_DL), appThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_UL), appThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_PLACEMENT_ID), fogAppPlacementId);
    });
  }

  // ==============================
  // UE
  // ==============================
  let linkWirelessType = 'wifi,4g'
  let linkPath = '[[7.419344,43.72764],[8.419344,43.72764]]';
  let linkPathMode = meep.GEO_EOP_MODE_REVERSE;
  let linkVelocity = '9';
  let ueMacId = '123456123456';

  function addUe(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_UE);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_LINK));
    verifyForm(meep.CFG_ELEM_CONNECTED, true, 'have.value', String(meep.OPT_CONNECTED.value));
    verifyForm(meep.CFG_ELEM_WIRELESS, true, 'have.value', String(meep.OPT_WIRELESS.value));
    type(meep.CFG_ELEM_LATENCY, linkLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, linkLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, linkPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, linkThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, linkThroughput-1);
    select(meep.CFG_ELEM_CONNECTED, meep.OPT_DISCONNECTED.label);
    type(meep.CFG_ELEM_WIRELESS_TYPE, linkWirelessType);
    type(meep.CFG_ELEM_GEO_LOCATION, linkLocationCoordinates);
    type(meep.CFG_ELEM_GEO_PATH, linkPath);
    type(meep.CFG_ELEM_GEO_VELOCITY, linkVelocity);
    select(meep.CFG_ELEM_GEO_EOP_MODE, linkPathMode);

    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_UE_MAC_ID, ueMacId);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validateUe(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_UE);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY), linkLatency);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY_VAR), linkLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_PKT_LOSS), linkPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_DL), linkThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_UL), linkThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_CONNECTED), false);
      assert.equal(getElemFieldVal(entry, FIELD_WIRELESS_TYPE), linkWirelessType);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), linkLocationCoordinates);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_PATH), linkPath);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_EOP_MODE), linkPathMode);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_VELOCITY), linkVelocity);
      assert.equal(getElemFieldVal(entry, FIELD_UE_MAC_ID), ueMacId);
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
  let ueAppPlacementId = 'node3';
  let ueAppCpuMin = '0.5';
  let ueAppCpuMax = '1';
  let ueAppMemoryMin = '100';
  let ueAppMemoryMax = '200';

  // Add new ue app element
  function addUeApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_UE_APP);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_APP));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_APP));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_APP));
    type(meep.CFG_ELEM_LATENCY, appLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, appLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, appPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, appThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, appThroughput-1);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, ueAppImg);
    type(meep.CFG_ELEM_GPU_COUNT, ueAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, ueAppGpuType);
    type(meep.CFG_ELEM_CPU_MIN, ueAppCpuMin);
    type(meep.CFG_ELEM_CPU_MAX, ueAppCpuMax);
    type(meep.CFG_ELEM_MEMORY_MIN, ueAppMemoryMin);
    type(meep.CFG_ELEM_MEMORY_MAX, ueAppMemoryMax);
    type(meep.CFG_ELEM_ENV, ueAppEnv);
    type(meep.CFG_ELEM_CMD, ueAppCmd);
    type(meep.CFG_ELEM_ARGS, ueAppArgs);
    type(meep.CFG_ELEM_PLACEMENT_ID, ueAppPlacementId);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
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
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MIN), ueAppCpuMin);
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MAX), ueAppCpuMax);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MIN), ueAppMemoryMin);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MAX), ueAppMemoryMax);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), ueAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), ueAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), ueAppArgs);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY), appLatency);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY_VAR), appLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_APP_PKT_LOSS), appPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_DL), appThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_UL), appThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_PLACEMENT_ID), ueAppPlacementId);
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
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_LINK));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_LINK));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_LINK));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_LINK));
    verifyForm(meep.CFG_ELEM_CONNECTED, true, 'have.value', String(meep.OPT_CONNECTED.value));
    verifyForm(meep.CFG_ELEM_WIRELESS, true, 'have.value', String(meep.OPT_WIRED.value));
    type(meep.CFG_ELEM_LATENCY, linkLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, linkLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, linkPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, linkThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, linkThroughput-1);
    type(meep.CFG_ELEM_GEO_LOCATION, linkLocationCoordinates);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
  }

  function validateCloud(name, parent) {
    cy.window().then((win) => {
      var entry = getEntry(win.meepStore.getState().cfg.table.entries, name);
      assert.isNotNull(entry);
      assert.equal(getElemFieldVal(entry, FIELD_TYPE), meep.ELEMENT_TYPE_DC);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY), linkLatency);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_LATENCY_VAR), linkLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_PKT_LOSS), linkPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_DL), linkThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_LINK_THROUGHPUT_UL), linkThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_GEO_LOCATION), linkLocationCoordinates);
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
  let cloudAppPlacementId = '';
  let cloudAppCpuMin = '0.5';
  let cloudAppCpuMax = '1';
  let cloudAppMemoryMin = '100';
  let cloudAppMemoryMax = '200';

  function addCloudApp(name, parent) {
    click(meep.CFG_BTN_NEW_ELEM);
    select(meep.CFG_ELEM_TYPE, meep.ELEMENT_TYPE_CLOUD_APP);
    verifyForm(meep.CFG_ELEM_LATENCY, true, 'have.value', String(meep.DEFAULT_LATENCY_APP));
    verifyForm(meep.CFG_ELEM_LATENCY_VAR, true, 'have.value', String(meep.DEFAULT_LATENCY_JITTER_APP));
    verifyForm(meep.CFG_ELEM_PKT_LOSS, true, 'have.value', String(meep.DEFAULT_PACKET_LOSS_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_DL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_DL_APP));
    verifyForm(meep.CFG_ELEM_THROUGHPUT_UL, true, 'have.value', String(meep.DEFAULT_THROUGHPUT_UL_APP));
    type(meep.CFG_ELEM_LATENCY, appLatency);
    type(meep.CFG_ELEM_LATENCY_VAR, appLatencyVar);
    type(meep.CFG_ELEM_PKT_LOSS, appPktLoss);
    type(meep.CFG_ELEM_THROUGHPUT_DL, appThroughput);
    type(meep.CFG_ELEM_THROUGHPUT_UL, appThroughput-1);
    select(meep.CFG_ELEM_PARENT, parent);
    type(meep.CFG_ELEM_NAME, name);
    type(meep.CFG_ELEM_IMG, cloudAppImg);
    type(meep.CFG_ELEM_PORT, cloudAppPort);
    type(meep.CFG_ELEM_EXT_PORT, cloudAppExtPort);
    select(meep.CFG_ELEM_PROT, cloudAppProt);
    type(meep.CFG_ELEM_GPU_COUNT, cloudAppGpuCount);
    select(meep.CFG_ELEM_GPU_TYPE, cloudAppGpuType);
    type(meep.CFG_ELEM_CPU_MIN, cloudAppCpuMin);
    type(meep.CFG_ELEM_CPU_MAX, cloudAppCpuMax);
    type(meep.CFG_ELEM_MEMORY_MIN, cloudAppMemoryMin);
    type(meep.CFG_ELEM_MEMORY_MAX, cloudAppMemoryMax);
    type(meep.CFG_ELEM_ENV, cloudAppEnv);
    type(meep.CFG_ELEM_CMD, cloudAppCmd);
    type(meep.CFG_ELEM_ARGS, cloudAppArgs);
    type(meep.CFG_ELEM_PLACEMENT_ID, cloudAppPlacementId);
    click(meep.MEEP_BTN_APPLY);
    verifyEnabled(meep.CFG_BTN_NEW_ELEM, true);
    verifyEnabled(meep.CFG_BTN_DEL_ELEM, false);
    verifyEnabled(meep.CFG_BTN_CLONE_ELEM, false);
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
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MIN), cloudAppCpuMin);
      assert.equal(getElemFieldVal(entry, FIELD_CPU_MAX), cloudAppCpuMax);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MIN), cloudAppMemoryMin);
      assert.equal(getElemFieldVal(entry, FIELD_MEMORY_MAX), cloudAppMemoryMax);
      assert.equal(getElemFieldVal(entry, FIELD_ENV_VAR), cloudAppEnv);
      assert.equal(getElemFieldVal(entry, FIELD_CMD), cloudAppCmd);
      assert.equal(getElemFieldVal(entry, FIELD_CMD_ARGS), cloudAppArgs);
      assert.equal(getElemFieldVal(entry, FIELD_PARENT), parent);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY), appLatency);
      assert.equal(getElemFieldVal(entry, FIELD_APP_LATENCY_VAR), appLatencyVar);
      assert.equal(getElemFieldVal(entry, FIELD_APP_PKT_LOSS), appPktLoss);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_DL), appThroughput);
      assert.equal(getElemFieldVal(entry, FIELD_APP_THROUGHPUT_UL), appThroughput-1);
      assert.equal(getElemFieldVal(entry, FIELD_PLACEMENT_ID), cloudAppPlacementId);
    });
  }

});


