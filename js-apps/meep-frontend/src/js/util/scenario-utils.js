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
import * as visdata from 'vis-data';
import { updateObject, deepCopy } from './object-util';
import uuid from 'uuid';

import {
  // Element Fields
  FIELD_TYPE,
  FIELD_PARENT,
  FIELD_NAME,
  FIELD_IMAGE,
  FIELD_PORT,
  FIELD_PROTOCOL,
  FIELD_GROUP,
  FIELD_GPU_COUNT,
  FIELD_GPU_TYPE,
  FIELD_PLACEMENT_ID,
  FIELD_INGRESS_SVC_MAP,
  FIELD_EGRESS_SVC_MAP,
  FIELD_ENV_VAR,
  FIELD_CMD,
  FIELD_CMD_ARGS,
  FIELD_EXT_PORT,
  FIELD_IS_EXTERNAL,
  FIELD_MCC,
  FIELD_MNC,
  FIELD_DEFAULT_CELL_ID,
  FIELD_CELL_ID,
  FIELD_NR_CELL_ID,
  FIELD_MAC_ID,
  FIELD_UE_MAC_ID,
  FIELD_GEO_LOCATION,
  FIELD_GEO_RADIUS,
  FIELD_GEO_PATH,
  FIELD_GEO_EOP_MODE,
  FIELD_GEO_VELOCITY,
  FIELD_CHART_ENABLED,
  FIELD_CHART_LOC,
  FIELD_CHART_VAL,
  FIELD_CHART_GROUP,
  FIELD_CONNECTED,
  FIELD_CONNECTIVITY_MODEL,
  FIELD_DN_NAME,
  FIELD_DN_LADN,
  FIELD_DN_ECSP,
  FIELD_WIRELESS,
  FIELD_WIRELESS_TYPE,
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
  FIELD_META_DISPLAY_MAP_COLOR,
  FIELD_META_DISPLAY_MAP_ICON,
  createElem,
  getElemFieldVal,
  setElemFieldVal,
  createUniqueName,
  FIELD_CPU_MIN,
  FIELD_CPU_MAX,
  FIELD_MEMORY_MIN,
  FIELD_MEMORY_MAX
} from './elem-utils';

import {
  ELEMENT_TYPE_SCENARIO,
  ELEMENT_TYPE_OPERATOR,
  ELEMENT_TYPE_OPERATOR_CELL,
  ELEMENT_TYPE_ZONE,
  ELEMENT_TYPE_POA,
  ELEMENT_TYPE_POA_4G,
  ELEMENT_TYPE_POA_5G,
  ELEMENT_TYPE_POA_WIFI,
  ELEMENT_TYPE_DC,
  ELEMENT_TYPE_CN,
  ELEMENT_TYPE_EDGE,
  ELEMENT_TYPE_FOG,
  ELEMENT_TYPE_UE,
  ELEMENT_TYPE_MECSVC,
  ELEMENT_TYPE_UE_APP,
  ELEMENT_TYPE_EDGE_APP,
  ELEMENT_TYPE_CLOUD_APP,
  DEFAULT_LATENCY_INTER_DOMAIN,
  DEFAULT_LATENCY_JITTER_INTER_DOMAIN,
  DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN,
  DEFAULT_THROUGHPUT_DL_INTER_DOMAIN,
  DEFAULT_THROUGHPUT_UL_INTER_DOMAIN,
  DEFAULT_PACKET_LOSS_INTER_DOMAIN,
  DEFAULT_LATENCY_INTER_ZONE,
  DEFAULT_LATENCY_JITTER_INTER_ZONE,
  DEFAULT_THROUGHPUT_DL_INTER_ZONE,
  DEFAULT_THROUGHPUT_UL_INTER_ZONE,
  DEFAULT_PACKET_LOSS_INTER_ZONE,
  DEFAULT_LATENCY_INTRA_ZONE,
  DEFAULT_LATENCY_JITTER_INTRA_ZONE,
  DEFAULT_THROUGHPUT_DL_INTRA_ZONE,
  DEFAULT_THROUGHPUT_UL_INTRA_ZONE,
  DEFAULT_PACKET_LOSS_INTRA_ZONE,
  DEFAULT_LATENCY_TERMINAL_LINK,
  DEFAULT_LATENCY_JITTER_TERMINAL_LINK,
  DEFAULT_THROUGHPUT_DL_TERMINAL_LINK,
  DEFAULT_THROUGHPUT_UL_TERMINAL_LINK,
  DEFAULT_PACKET_LOSS_TERMINAL_LINK,
  DEFAULT_LATENCY_LINK,
  DEFAULT_LATENCY_JITTER_LINK,
  DEFAULT_THROUGHPUT_DL_LINK,
  DEFAULT_THROUGHPUT_UL_LINK,
  DEFAULT_PACKET_LOSS_LINK,
  DEFAULT_LATENCY_APP,
  DEFAULT_LATENCY_JITTER_APP,
  DEFAULT_THROUGHPUT_DL_APP,
  DEFAULT_THROUGHPUT_UL_APP,
  DEFAULT_PACKET_LOSS_APP,
  // DEFAULT_LATENCY_DC,
  DOMAIN_TYPE_STR,
  DOMAIN_CELL_TYPE_STR,
  PUBLIC_DOMAIN_TYPE_STR,
  ZONE_TYPE_STR,
  COMMON_ZONE_TYPE_STR,
  POA_TYPE_STR,
  POA_4G_TYPE_STR,
  POA_5G_TYPE_STR,
  POA_WIFI_TYPE_STR,
  DEFAULT_NL_TYPE_STR,
  UE_TYPE_STR,
  FOG_TYPE_STR,
  EDGE_TYPE_STR,
  CN_TYPE_STR,
  DC_TYPE_STR,
  MEC_SVC_TYPE_STR,
  UE_APP_TYPE_STR,
  EDGE_APP_TYPE_STR,
  CLOUD_APP_TYPE_STR,

  DEFAULT_CONNECTIVITY_MODEL,

  // Logical Scenario types
  TYPE_SCENARIO,
  TYPE_DOMAIN,
  TYPE_ZONE,
  TYPE_NET_LOC,
  TYPE_PHY_LOC,
  TYPE_PROCESS
} from '../meep-constants';

import {
  META_DISPLAY_MAP_COLOR,
  META_DISPLAY_MAP_ICON
} from './meta-keys';

// Import images used in JS
import * as poaImage from '../../img/tower-02-idcc.svg';
import * as edgeImage from '../../img/edge-idcc.svg';
import * as fogImage from '../../img/fog-idcc.svg';
import * as screenImage from '../../img/Screen-02-idcc.svg';
import * as cloudServerBlackBlue from '../../img/cloud-server-black-blue.svg';
import * as cloudBlack from '../../img/cloud-black.svg';
import * as cloudOutlineBlack from '../../img/cloud-outline-black.svg';
import * as switchBlue from '../../img/switch-blue.svg';
import * as droneBlue from '../../img/drone-blue.svg';
import * as droneBlack from '../../img/drone-black.svg';

// Parse scenario to populate visualization & table
export function parseScenario(scenario) {
  if (!scenario) {
    return null;
  }

  var nodes = new Array();
  var edges = new Array();
  var ueList = new Array();
  var poaList = new Array();
  var computeList = new Array();

  // Add scenario to graph and table
  addScenarioNode(scenario, nodes);

  // Domains
  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];

    // Add domain to graph and table (ignore public domain)
    if (domain.name !== PUBLIC_DOMAIN_TYPE_STR) {
      addDomainNode(domain, scenario, nodes, edges);
    }

    // Zones
    for (var j in domain.zones) {
      var zone = domain.zones[j];

      // Add zone to graph and table (ignore common zone)
      if (zone.name.indexOf(COMMON_ZONE_TYPE_STR) === -1) {
        const parent =
          domain.type === PUBLIC_DOMAIN_TYPE_STR ? scenario : domain;
        addZoneNode(zone, parent, nodes, edges);
      }

      // Network Locations
      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];

        // Add Network Location to graph and table (ignore default network location)
        if (nl.name.indexOf(DEFAULT_NL_TYPE_STR) === -1) {
          const parent = domain.type === PUBLIC_DOMAIN_TYPE_STR ? scenario :
            zone.type === COMMON_ZONE_TYPE_STR ? domain : zone;
          addNlNode(nl, parent, nodes, edges);

          // Add NL with geodata to map
          if (nl.geoData && nl.geoData.location) {
            var nlGeoDataAsset = updateObject({assetName: nl.name, assetType: 'POA', subType: nl.type}, nl.geoData);
            poaList.push(nlGeoDataAsset);
          }
        }

        // Physical Locations
        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];

          // Add Physical Location to graph and table
          const parent = domain.type === PUBLIC_DOMAIN_TYPE_STR ? scenario :
            zone.type === COMMON_ZONE_TYPE_STR ? domain :
              nl.type === DEFAULT_NL_TYPE_STR ? zone : nl;
          addPlNode(pl, parent, nodes, edges);

          // Add PL with geodata to map
          if (pl.geoData && pl.geoData.location) {
            var plGeoDataAsset = updateObject({assetName: pl.name, subType: pl.type}, pl.geoData);
            if (pl.type === UE_TYPE_STR) {
              plGeoDataAsset.assetType = 'UE';
              ueList.push(plGeoDataAsset);
            } else {
              plGeoDataAsset.assetType = 'COMPUTE';
              computeList.push(plGeoDataAsset);
            }
          }

          // Processes
          for (var m in pl.processes) {
            var proc = pl.processes[m];

            // Add Process to graph and table (for configuration only)
            addProcessNode(proc, pl, nodes, edges);
          }
        }
      }
    }
  }

  // Update table
  var table = {};
  table.data = { edges: edges, nodes: nodes };
  table.entries = _.reduce(table.data.nodes, (nodeMap, node) => {
    nodeMap[node.name] = updateObject(node, getElementFromScenario(scenario, node.id));
    return nodeMap;
  }, {});

  // Update visualization data
  var visData = {};
  visData.nodes = new visdata.DataSet(nodes);
  visData.edges = new visdata.DataSet(edges);

  // Update map data
  var mapData = {};
  mapData.ueList = _.sortBy(ueList, ['assetName']);
  mapData.poaList = _.sortBy(poaList, ['assetName']);
  mapData.computeList = _.sortBy(computeList, ['assetName']);

  return { table: table, visData: visData, mapData: mapData };
}

function findIdInScenario(scenario, uniqueId) {

  // Domains
  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];

    // Add domain to graph and table (ignore public domain)
    if (domain.id === uniqueId) {
      return true;
    }

    // Zones
    for (var j in domain.zones) {
      var zone = domain.zones[j];

      if (zone.id === uniqueId) {
        return true;
      }

      // Network Locations
      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];

        if (nl.id === uniqueId) {
          return true;
        }

        // Physical Locations
        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];

          if (pl.id === uniqueId) {
            return true;
          }

          // Processes
          for (var m in pl.processes) {
            var proc = pl.processes[m];

            if (proc.id === uniqueId) {
              return true;
            }
          }
        }
      }
    }
  }
  return false;
}

export const getUniqueId = (scenario) => {
  var uniqueId = uuid();
  var isUniqueId = false;
  while(!isUniqueId) {
    isUniqueId = true;
    if (findIdInScenario(scenario, uniqueId)) {
      uniqueId = uuid();
      isUniqueId = false;
    }
  }
  return uniqueId;
};

// Add network element to scenario
export function addElementToScenario(scenario, element) {
  var scenarioElement;
  var type = getElemFieldVal(element, FIELD_TYPE);
  var name = getElemFieldVal(element, FIELD_NAME);
  var parent = getElemFieldVal(element, FIELD_PARENT);
  var uniqueId = getUniqueId(scenario);

  // Prepare network element to be added to scenario
  switch (type) {
  case ELEMENT_TYPE_OPERATOR:
    scenarioElement = createDomain(uniqueId, name, element);
    break;
  case ELEMENT_TYPE_OPERATOR_CELL:
    scenarioElement = createDomainCell(uniqueId, name, element);
    break;
  case ELEMENT_TYPE_ZONE:
    scenarioElement = createZone(uniqueId, name, element);
    break;
  case ELEMENT_TYPE_POA_4G:
    scenarioElement = createPoa4G(uniqueId, name, element);
    break;
  case ELEMENT_TYPE_POA_5G:
    scenarioElement = createPoa5G(uniqueId, name, element);
    break;
  case ELEMENT_TYPE_POA_WIFI:
    scenarioElement = createPoaWIFI(uniqueId, name, element);
    break;
  case ELEMENT_TYPE_POA:
    scenarioElement = createPoa(uniqueId, name, element);
    break;
  case ELEMENT_TYPE_DC:
    setElemFieldVal(element, FIELD_PARENT,
      PUBLIC_DOMAIN_TYPE_STR + '-' + COMMON_ZONE_TYPE_STR + '-' + DEFAULT_NL_TYPE_STR);
    scenarioElement = createPL(uniqueId, name, DC_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_CN:
    setElemFieldVal(element, FIELD_PARENT,
      (parent += '-' + COMMON_ZONE_TYPE_STR + '-' + DEFAULT_NL_TYPE_STR)
    );
    scenarioElement = createPL(uniqueId, name, CN_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_EDGE:
    setElemFieldVal(element, FIELD_PARENT, (parent += '-' + DEFAULT_NL_TYPE_STR));
    scenarioElement = createPL(uniqueId, name, EDGE_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_FOG:
    scenarioElement = createPL(uniqueId, name, FOG_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_UE:
    scenarioElement = createPL(uniqueId, name, UE_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_MECSVC:
    scenarioElement = createProcess(uniqueId, name, MEC_SVC_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_UE_APP:
    scenarioElement = createProcess(uniqueId, name, UE_APP_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_EDGE_APP:
    scenarioElement = createProcess(uniqueId, name, EDGE_APP_TYPE_STR, element);
    break;
  case ELEMENT_TYPE_CLOUD_APP:
    scenarioElement = createProcess(uniqueId, name, CLOUD_APP_TYPE_STR, element);
    break;
  default:
    break;
  }

  // Find parent node
  parent = getElemFieldVal(element, FIELD_PARENT);
  if (scenario.name === parent) {
    scenario.deployment.domains.push(scenarioElement);
    return;
  }

  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];
    if (domain.name === parent) {
      if (domain.zones === undefined) {
        domain.zones = [];
      }
      domain.zones.push(scenarioElement);
      return;
    }

    for (var j in domain.zones) {
      var zone = domain.zones[j];
      if (zone.name === parent) {
        if (zone.networkLocations === undefined) {
          zone.networkLocations = [];
        }
        zone.networkLocations.push(scenarioElement);
        return;
      }

      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];
        if (nl.name === parent) {
          if (nl.physicalLocations === undefined) {
            nl.physicalLocations = [];
          }
          nl.physicalLocations.push(scenarioElement);
          return;
        }

        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];
          if (pl.name === parent) {
            if (pl.processes === undefined) {
              pl.processes = [];
            }
            pl.processes.push(scenarioElement);
            return;
          }
        }
      }
    }
  }
}

// Update network element in scenario
export function updateElementInScenario(scenario, element) {
  var name = getElemFieldVal(element, FIELD_NAME);
  var id = element.id;

  // Find element in scenario
  if (scenario.name === name) {
    if (!scenario.deployment.netChar) {
      scenario.deployment.netChar = {};
    }
    scenario.deployment.netChar.latency = getElemFieldVal(element, FIELD_INT_DOM_LATENCY);
    scenario.deployment.netChar.latencyVariation = getElemFieldVal(element, FIELD_INT_DOM_LATENCY_VAR);
    scenario.deployment.netChar.latencyDistribution = getElemFieldVal(element, FIELD_INT_DOM_LATENCY_DIST);
    scenario.deployment.netChar.throughputDl = getElemFieldVal(element, FIELD_INT_DOM_THROUGHPUT_DL);
    scenario.deployment.netChar.throughputUl = getElemFieldVal(element, FIELD_INT_DOM_THROUGHPUT_UL);
    scenario.deployment.netChar.packetLoss = getElemFieldVal(element, FIELD_INT_DOM_PKT_LOSS);

    if (!scenario.deployment.connectivity) {
      scenario.deployment.connectivity = {};
    }
    scenario.deployment.connectivity.model = getElemFieldVal(element, FIELD_CONNECTIVITY_MODEL);
    return;
  }

  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];
    if (domain.id === id) {

      if (!domain.netChar) {
        domain.netChar = {};
      }
      domain.netChar.latency = getElemFieldVal(element, FIELD_INT_ZONE_LATENCY);
      domain.netChar.latencyVariation = getElemFieldVal(element, FIELD_INT_ZONE_LATENCY_VAR);
      domain.netChar.throughputDl = getElemFieldVal(element, FIELD_INT_ZONE_THROUGHPUT_DL);
      domain.netChar.throughputUl = getElemFieldVal(element, FIELD_INT_ZONE_THROUGHPUT_UL);
      domain.netChar.packetLoss = getElemFieldVal(element, FIELD_INT_ZONE_PKT_LOSS);

      if (domain.type === DOMAIN_CELL_TYPE_STR) {
        var cellularDomainConfig = {
          mcc: getElemFieldVal(element, FIELD_MCC),
          mnc: getElemFieldVal(element, FIELD_MNC),
          defaultCellId: getElemFieldVal(element, FIELD_DEFAULT_CELL_ID)
        };
        domain.cellularDomainConfig = cellularDomainConfig;
      }

      //if domain name changed, other elements created based on that name must also be updated (default ones)
      for (var i2 in domain.zones) {
        var zoneCommon = domain.zones[i2];
        if (zoneCommon.id === domain.name + '-' + COMMON_ZONE_TYPE_STR) {
          for (var i3 in zoneCommon.networkLocations) {
            var nlDomainCommon = zoneCommon.networkLocations[i3];
            if (nlDomainCommon.id === zoneCommon.name + '-' + DEFAULT_NL_TYPE_STR) {
              nlDomainCommon.id = name + '-' + COMMON_ZONE_TYPE_STR + '-' + DEFAULT_NL_TYPE_STR;
              nlDomainCommon.name = name + '-' + COMMON_ZONE_TYPE_STR + '-' + DEFAULT_NL_TYPE_STR;
              break;
            }
          }
          zoneCommon.id = name + '-' + COMMON_ZONE_TYPE_STR;
          zoneCommon.name = name + '-' + COMMON_ZONE_TYPE_STR;
          break;
        }
      }
      domain.label = name;
      domain.name = name;
      return;
    }

    for (var j in domain.zones) {
      var zone = domain.zones[j];
      if (zone.id === id) {

        if (!zone.netChar) {
          zone.netChar = {};
        }
        zone.netChar.latency = getElemFieldVal(element, FIELD_INTRA_ZONE_LATENCY);
        zone.netChar.latencyVariation = getElemFieldVal(element, FIELD_INTRA_ZONE_LATENCY_VAR);
        zone.netChar.throughputDl = getElemFieldVal(element, FIELD_INTRA_ZONE_THROUGHPUT_DL);
        zone.netChar.throughputUl = getElemFieldVal(element, FIELD_INTRA_ZONE_THROUGHPUT_UL);
        zone.netChar.packetLoss = getElemFieldVal(element, FIELD_INTRA_ZONE_PKT_LOSS);

        //if zone name changed, other elements created based on that name must also be updated (default ones)
        for (var j2 in zone.networkLocations) {
          var nlZoneCommon = zone.networkLocations[j2];
          if (nlZoneCommon.id === zone.name + '-' + DEFAULT_NL_TYPE_STR) {
            nlZoneCommon.id = name + '-' + DEFAULT_NL_TYPE_STR;
            nlZoneCommon.name = name + '-' + DEFAULT_NL_TYPE_STR;
          }
        }

        if (!zone.meta) {
          zone.meta = {};
        }
        zone.meta[META_DISPLAY_MAP_COLOR] = getElemFieldVal(element, FIELD_META_DISPLAY_MAP_COLOR);

        zone.label = name;
        zone.name = name;
        return;
      }

      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];
        if (nl.id === id) {

          if (!nl.netChar) {
            nl.netChar = {};
          }
          nl.netChar.latency = getElemFieldVal(element, FIELD_TERM_LINK_LATENCY);
          nl.netChar.latencyVariation = getElemFieldVal(element, FIELD_TERM_LINK_LATENCY_VAR);
          nl.netChar.throughputDl = getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_DL);
          nl.netChar.throughputUl = getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_UL);
          nl.netChar.packetLoss = getElemFieldVal(element, FIELD_TERM_LINK_PKT_LOSS);

          if (nl.type === POA_4G_TYPE_STR) {
            var poa4GConfig = {
              cellId: getElemFieldVal(element, FIELD_CELL_ID)
            };
            nl.poa4GConfig = poa4GConfig;
          }
          if (nl.type === POA_5G_TYPE_STR) {
            var poa5GConfig = {
              cellId: getElemFieldVal(element, FIELD_NR_CELL_ID)
            };
            nl.poa5GConfig = poa5GConfig;
          }
          if (nl.type === POA_WIFI_TYPE_STR) {
            var poaWifiConfig = {
              macId: getElemFieldVal(element, FIELD_MAC_ID)
            };
            nl.poaWifiConfig = poaWifiConfig;
          }

          if (!nl.geoData) {
            nl.geoData = {};
          }
          var nlLocation = getElemFieldVal(element, FIELD_GEO_LOCATION);
          nl.geoData.location = !nlLocation ? null : {
            type: 'Point',
            coordinates: JSON.parse(nlLocation)
          };
          var radius = getElemFieldVal(element, FIELD_GEO_RADIUS);
          nl.geoData.radius = (radius === '') ? null : radius;

          nl.label = name;
          nl.name = name;
          return;
        }

        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];
          if (pl.id === id) {

            if (!pl.netChar) {
              pl.netChar = {};
            }
            pl.netChar.latency = getElemFieldVal(element, FIELD_LINK_LATENCY);
            pl.netChar.latencyVariation = getElemFieldVal(element, FIELD_LINK_LATENCY_VAR);
            pl.netChar.throughputDl = getElemFieldVal(element, FIELD_LINK_THROUGHPUT_DL);
            pl.netChar.throughputUl = getElemFieldVal(element, FIELD_LINK_THROUGHPUT_UL);
            pl.netChar.packetLoss = getElemFieldVal(element, FIELD_LINK_PKT_LOSS);

            pl.connected = getElemFieldVal(element, FIELD_CONNECTED);
            var wireless = getElemFieldVal(element, FIELD_WIRELESS);
            pl.wireless = wireless;
            pl.wirelessType = wireless ? getElemFieldVal(element, FIELD_WIRELESS_TYPE) : '';

            if (!pl.dataNetwork) {
              pl.dataNetwork = {};
            }
            pl.dataNetwork.dnn = getElemFieldVal(element, FIELD_DN_NAME);
            pl.dataNetwork.ladn = getElemFieldVal(element, FIELD_DN_LADN);
            pl.dataNetwork.ecsp = getElemFieldVal(element, FIELD_DN_ECSP);

            if (!pl.geoData) {
              pl.geoData = {};
            }
            var plLocation = getElemFieldVal(element, FIELD_GEO_LOCATION);
            pl.geoData.location = !plLocation ? null : {
              type: 'Point',
              coordinates: JSON.parse(plLocation)
            };
            var path = getElemFieldVal(element, FIELD_GEO_PATH);
            pl.geoData.path = !path ? null : {
              type: 'LineString',
              coordinates: JSON.parse(path)
            };
            pl.geoData.eopMode = getElemFieldVal(element, FIELD_GEO_EOP_MODE);
            const velocity = getElemFieldVal(element, FIELD_GEO_VELOCITY);
            pl.geoData.velocity = velocity ? velocity : null;

            pl.macId = getElemFieldVal(element, FIELD_UE_MAC_ID);

            pl.label = name;
            pl.name = name;
            return;
          }

          for (var m in pl.processes) {
            var process = pl.processes[m];
            if (process.id === id) {
              pl.processes[m] = createProcess(
                process.id,
                name,
                process.type,
                element
              );
              return;
            }
          }
        }
      }
    }
  }
}

// Clone network element in scenario
export function cloneElementInScenario(scenario, element, table) {
  var inDomainCloneBranch = false, inZoneCloneBranch = false, inNlCloneBranch = false, inPlCloneBranch = false;
  var newZoneRootParentName = '';
  var newNlRootParentName = '';
  var newPlRootParentName = '';
  var newProcessRootParentName = '';
  var elementFromScenario;
  var parent = getElemFieldVal(element, FIELD_PARENT);

  // Domains
  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];

    // Add domain to graph and table (ignore public domain)
    if (domain.id === element.id) {
      newZoneRootParentName = cloneElement(scenario, element, parent, true, table);
      inDomainCloneBranch = true;
    } else {
      inDomainCloneBranch = false;
    }

    // Zones
    for (var j in domain.zones) {
      var zone = domain.zones[j];

      if (inDomainCloneBranch) {
        if (zone.name.indexOf(COMMON_ZONE_TYPE_STR) !== -1) {
          newNlRootParentName = newZoneRootParentName + COMMON_ZONE_TYPE_STR;
        } else {
          elementFromScenario = getElementFromScenario(scenario, zone.id);
          newNlRootParentName = cloneElement(scenario, elementFromScenario, newZoneRootParentName, false, table);
        }
      } else {
        if (zone.id === element.id) {
          newNlRootParentName = cloneElement(scenario, element, parent, true, table);
          inZoneCloneBranch = true;
        } else {
          inZoneCloneBranch = false;
        }
      }

      // Network Locations
      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];

        if (inDomainCloneBranch || inZoneCloneBranch) {
          if (nl.name.indexOf(DEFAULT_NL_TYPE_STR) !== -1) {
            newPlRootParentName = newNlRootParentName;
          } else {
            elementFromScenario = getElementFromScenario(scenario, nl.id);
            newPlRootParentName = cloneElement(scenario, elementFromScenario, newNlRootParentName, false, table);
          }
        } else {
          if (nl.id === element.id) {
            newPlRootParentName = cloneElement(scenario, element, parent, true, table);
            inNlCloneBranch = true;
          } else {
            inNlCloneBranch = false;
          }
        }

        // Physical Locations
        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];

          if (inDomainCloneBranch || inZoneCloneBranch || inNlCloneBranch) {
            elementFromScenario = getElementFromScenario(scenario, pl.id);
            newProcessRootParentName = cloneElement(scenario, elementFromScenario, newPlRootParentName, false, table);
          } else {
            if (pl.id === element.id) {
              newProcessRootParentName = cloneElement(scenario, element, parent, true, table);
              inPlCloneBranch = true;
            } else {
              inPlCloneBranch = false;
            }
          }

          // Processes
          for (var m in pl.processes) {
            var proc = pl.processes[m];

            if (inDomainCloneBranch || inZoneCloneBranch || inNlCloneBranch || inPlCloneBranch) {
              elementFromScenario = getElementFromScenario(scenario, proc.id);
              cloneElement(scenario, elementFromScenario, newProcessRootParentName, false, table);
            } else {
              if (proc.id === element.id) {
                cloneElement(scenario, element, parent, true, table);
              }
            }
          }
        }
      }
    }

    if(inDomainCloneBranch || inZoneCloneBranch || inNlCloneBranch || inPlCloneBranch) {
      break;
    }
  }
}

// CLONE ELEMENT, return new element name
function cloneElement(scenario, element, newParentName, isRoot, table) {
  let newElement = deepCopy(element);

  var name = getElemFieldVal(element, FIELD_NAME);
  if (isRoot === false) {
    name = createUniqueName(table.entries, name + '-copy');
    setElemFieldVal(newElement, FIELD_NAME, name);
  }
  setElemFieldVal(newElement, FIELD_PARENT, newParentName);

  // The following element fields cause issues when duplicated in the scenario
  // For now, set these values to null when cloning
  // TODO -- Improve frontend cloning or move scenario configuration to brackend
  if (getElemFieldVal(element, FIELD_EXT_PORT)) {
    setElemFieldVal(newElement, FIELD_EXT_PORT, null);
  }
  if (getElemFieldVal(element, FIELD_INGRESS_SVC_MAP)) {
    setElemFieldVal(newElement, FIELD_INGRESS_SVC_MAP, null);
  }
  if (getElemFieldVal(element, FIELD_EGRESS_SVC_MAP)) {
    setElemFieldVal(newElement, FIELD_EGRESS_SVC_MAP, null);
  }

  // add new element to scenario
  // new id and label will be created as part of the addNewElementToScenario called by newScenarioElem
  addElementToScenario(scenario, newElement);
  return name;
}

// Remove the specific element and its children from the scenario
export function removeElementFromScenario(scenario, element) {
  var name = getElemFieldVal(element, FIELD_NAME);

  // Loop through scenario until element is found
  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];
    if (domain.name === name) {
      scenario.deployment.domains.splice(i, 1);
      return;
    }

    for (var j in domain.zones) {
      var zone = domain.zones[j];
      if (zone.name === name) {
        domain.zones.splice(j, 1);
        return;
      }

      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];
        if (nl.name === name) {
          zone.networkLocations.splice(k, 1);
          return;
        }

        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];
          if (pl.name === name) {
            nl.physicalLocations.splice(l, 1);
            return;
          }

          for (var m in pl.processes) {
            var process = pl.processes[m];
            if (process.name === name) {
              pl.processes.splice(m, 1);
              return;
            }
          }
        }
      }
    }
  }
}

// Create a new scenario with given name
export function createNewScenario(name) {
  var scenario = {
    name: name,
    deployment: {
      netChar: {
        latency: parseInt(DEFAULT_LATENCY_INTER_DOMAIN),
        latencyVariation: parseInt(DEFAULT_LATENCY_JITTER_INTER_DOMAIN),
        latencyDistribution: DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN,
        throughputDl: parseInt(DEFAULT_THROUGHPUT_DL_INTER_DOMAIN),
        throughputUl: parseInt(DEFAULT_THROUGHPUT_UL_INTER_DOMAIN),
        interDomainPacketLoss: parseInt(DEFAULT_PACKET_LOSS_INTER_DOMAIN)
      },
      connectivity: {
        model: DEFAULT_CONNECTIVITY_MODEL
      },
      domains: name === 'None' ? [] : [createDefaultDomain()]
    }
  };
  return scenario;
}

export function createProcess(uniqueId, name, type, element) {
  var isExternal = getElemFieldVal(element, FIELD_IS_EXTERNAL);
  var port = getElemFieldVal(element, FIELD_PORT);
  var gpuCount = getElemFieldVal(element, FIELD_GPU_COUNT);
  var cpuMin = getElemFieldVal(element, FIELD_CPU_MIN);
  var cpuMax = getElemFieldVal(element, FIELD_CPU_MAX);
  var memoryMin = getElemFieldVal(element, FIELD_MEMORY_MIN);
  var memoryMax = getElemFieldVal(element, FIELD_MEMORY_MAX);
  var process = {
    id: uniqueId,
    name: name,
    type: type,
    isExternal: isExternal,
    userChartLocation: null,
    userChartAlternateValues: null,
    userChartGroup: null,
    image: null,
    environment: null,
    commandArguments: null,
    commandExe: null,
    serviceConfig: null,
    gpuConfig: null,
    cpuConfig: null,
    memoryConfig: null,
    externalConfig: null,
    netChar: {
      latency: parseInt(DEFAULT_LATENCY_APP),
      latencyVariation: parseInt(DEFAULT_LATENCY_JITTER_APP),
      throughputDl: parseInt(DEFAULT_THROUGHPUT_DL_APP),
      throughputUl: parseInt(DEFAULT_THROUGHPUT_UL_APP),
      packetLoss: parseInt(DEFAULT_PACKET_LOSS_APP)
    },
    placementId: null
  };

  if (isExternal) {
    process.externalConfig = {
      ingressServiceMap: getIngressServiceMapArray(
        getElemFieldVal(element, FIELD_INGRESS_SVC_MAP)),
      egressServiceMap: getEgressServiceMapArray(
        getElemFieldVal(element, FIELD_EGRESS_SVC_MAP))
    };
    process.placementId = getElemFieldVal(element, FIELD_PLACEMENT_ID);
  } else if (getElemFieldVal(element, FIELD_CHART_ENABLED)) {
    process.userChartLocation = getElemFieldVal(element, FIELD_CHART_LOC);
    process.userChartAlternateValues = getElemFieldVal(element, FIELD_CHART_VAL);
    process.userChartGroup = getElemFieldVal(element, FIELD_CHART_GROUP);
  } else {
    process.image = getElemFieldVal(element, FIELD_IMAGE);
    process.environment = getElemFieldVal(element, FIELD_ENV_VAR);
    process.commandArguments = getElemFieldVal(element, FIELD_CMD_ARGS);
    process.commandExe = getElemFieldVal(element, FIELD_CMD);
    process.serviceConfig = isNaN(port) || !port ? null : {
      name: name,
      meSvcName: getElemFieldVal(element, FIELD_GROUP),
      // TODO -- Add frontend support for multiple ports
      ports: [{
        protocol: getElemFieldVal(element, FIELD_PROTOCOL) === '' ? null :
          getElemFieldVal(element, FIELD_PROTOCOL).toUpperCase(),
        port: getElemFieldVal(element, FIELD_PORT) === '' ? null :
          getElemFieldVal(element, FIELD_PORT),
        externalPort: getElemFieldVal(element, FIELD_EXT_PORT) === '' ? null :
          getElemFieldVal(element, FIELD_EXT_PORT)
      }]
    };
    process.gpuConfig = isNaN(gpuCount) || !gpuCount ? null : {
      type: getElemFieldVal(element, FIELD_GPU_TYPE) === '' ? null :
        getElemFieldVal(element, FIELD_GPU_TYPE).toUpperCase(),
      count: gpuCount
    };
    process.cpuConfig = (cpuMin && !isNaN(cpuMin)) || (cpuMax && !isNaN(cpuMax)) ? {
      min: cpuMin && !isNaN(cpuMin) ? parseFloat(cpuMin) : null,
      max: cpuMax && !isNaN(cpuMax) ? parseFloat(cpuMax): null
    } : null;
    process.memoryConfig = (memoryMin && !isNaN(memoryMin)) || (memoryMax && !isNaN(memoryMax)) ? {
      min: memoryMin && !isNaN(memoryMin) ? parseInt(memoryMin) : null,
      max: memoryMax && !isNaN(memoryMax) ? parseInt(memoryMax) : null
    } : null;
    process.placementId = getElemFieldVal(element, FIELD_PLACEMENT_ID);
  }
  if (process.netChar) {
    process.netChar.latency = getElemFieldVal(element, FIELD_APP_LATENCY);
    process.netChar.latencyVariation = getElemFieldVal(element, FIELD_APP_LATENCY_VAR);
    process.netChar.throughputDl = getElemFieldVal(element, FIELD_APP_THROUGHPUT_DL);
    process.netChar.throughputUl = getElemFieldVal(element, FIELD_APP_THROUGHPUT_UL);
    process.netChar.packetLoss = getElemFieldVal(element, FIELD_APP_PKT_LOSS);
  }

  return process;
}

export function getIngressServiceMapStr(ingressServiceMapArray) {
  var ingressServiceMapStr = '';

  // Loop through service map array
  for (var i = 0; i < ingressServiceMapArray.length; i++) {
    var svcMap = ingressServiceMapArray[i];
    ingressServiceMapStr +=
      (i === 0 ? '' : ',') +
      svcMap.externalPort +
      ':' +
      svcMap.name +
      ':' +
      svcMap.port +
      ':' +
      svcMap.protocol;
  }
  return ingressServiceMapStr;
}

export function getIngressServiceMapArray(ingressServiceMapStr) {
  var ingressServiceMapArray = [];

  // Add service map entries, if any
  if (ingressServiceMapStr) {
    var scpMapList = ingressServiceMapStr.split(',');
    // Loop through service map list
    for (var i = 0; i < scpMapList.length; i++) {
      var svcMap = scpMapList[i].split(':');
      if (svcMap.length !== 4) {
        continue;
      }

      // Add service map to ingressServiceMap Array
      ingressServiceMapArray.push({
        externalPort: parseInt(svcMap[0]),
        name: svcMap[1],
        port: parseInt(svcMap[2]),
        protocol: svcMap[3].toUpperCase()
      });
    }
  }
  return ingressServiceMapArray;
}

export function getEgressServiceMapStr(egressServiceMapArray) {
  var egressServiceMapStr = '';

  // Loop through service map array
  for (var i = 0; i < egressServiceMapArray.length; i++) {
    var svcMap = egressServiceMapArray[i];
    egressServiceMapStr +=
      (i === 0 ? '' : ',') +
      svcMap.name +
      ':' +
      (svcMap.meSvcName ? svcMap.meSvcName : '') +
      ':' +
      svcMap.ip +
      ':' +
      svcMap.port +
      ':' +
      svcMap.protocol;
  }
  return egressServiceMapStr;
}

export function getEgressServiceMapArray(egressServiceMapStr) {
  var egressServiceMapArray = [];

  // Add service map entries, if any
  if (egressServiceMapStr) {
    var scpMapList = egressServiceMapStr.split(',');
    // Loop through service map list
    for (var i = 0; i < scpMapList.length; i++) {
      var svcMap = scpMapList[i].split(':');
      if (svcMap.length !== 5) {
        continue;
      }

      // Add service map to egressServiceMap Array
      egressServiceMapArray.push({
        name: svcMap[0],
        meSvcName: svcMap[1],
        ip: svcMap[2],
        port: parseInt(svcMap[3]),
        protocol: svcMap[4].toUpperCase()
      });
    }
  }
  return egressServiceMapArray;
}

export function createDomain(uniqueId, name, element) {
  var domain = {
    id: uniqueId,
    name: name,
    type: DOMAIN_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_INT_ZONE_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_INT_ZONE_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_INT_ZONE_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_INT_ZONE_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_INT_ZONE_PKT_LOSS)
    },
    zones: [createDefaultZone(name)]
  };
  return domain;
}

export function createDomainCell(uniqueId, name, element) {
  var domain = {
    id: uniqueId,
    name: name,
    type: DOMAIN_CELL_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_INT_ZONE_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_INT_ZONE_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_INT_ZONE_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_INT_ZONE_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_INT_ZONE_PKT_LOSS)
    },
    zones: [createDefaultZone(name)],
    cellularDomainConfig: {
      mcc: getElemFieldVal(element, FIELD_MCC),
      mnc: getElemFieldVal(element, FIELD_MNC),
      defaultCellId: getElemFieldVal(element, FIELD_DEFAULT_CELL_ID)
    }
  };
  return domain;
}

export function createDefaultDomain() {
  var domain = {
    id: PUBLIC_DOMAIN_TYPE_STR,
    name: PUBLIC_DOMAIN_TYPE_STR,
    type: PUBLIC_DOMAIN_TYPE_STR,
    netChar: {
      latency: parseInt(DEFAULT_LATENCY_INTER_ZONE),
      latencyVariation: parseInt(DEFAULT_LATENCY_JITTER_INTER_ZONE),
      throughputDl: parseInt(DEFAULT_THROUGHPUT_DL_INTER_ZONE),
      throughputUl: parseInt(DEFAULT_THROUGHPUT_UL_INTER_ZONE),
      packetLoss: parseInt(DEFAULT_PACKET_LOSS_INTER_ZONE)
    },
    zones: [createDefaultZone(PUBLIC_DOMAIN_TYPE_STR)]
  };
  return domain;
}

export function createPoa(uniqueId, name, element) {
  var location = getElemFieldVal(element, FIELD_GEO_LOCATION);
  var radius = getElemFieldVal(element, FIELD_GEO_RADIUS);
  var nl = {
    id: uniqueId,
    name: name,
    type: POA_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_TERM_LINK_PKT_LOSS)
    },
    geoData: !location ? null : {
      location: {
        type: 'Point',
        coordinates: JSON.parse(location)
      },
      radius: (radius === '') ? null : radius
    },
    physicalLocations: []
  };

  return nl;
}

export function createPoa4G(uniqueId, name, element) {
  var location = getElemFieldVal(element, FIELD_GEO_LOCATION);
  var radius = getElemFieldVal(element, FIELD_GEO_RADIUS);
  var nl = {
    id: uniqueId,
    name: name,
    type: POA_4G_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_TERM_LINK_PKT_LOSS)
    },
    physicalLocations: [],
    poa4GConfig: {
      cellId: getElemFieldVal(element, FIELD_CELL_ID)
    },
    geoData: !location ? null : {
      location: {
        type: 'Point',
        coordinates: JSON.parse(location)
      },
      radius: (radius === '') ? null : radius
    }
  };

  return nl;
}

export function createPoa5G(uniqueId, name, element) {
  var location = getElemFieldVal(element, FIELD_GEO_LOCATION);
  var radius = getElemFieldVal(element, FIELD_GEO_RADIUS);
  var nl = {
    id: uniqueId,
    name: name,
    type: POA_5G_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_TERM_LINK_PKT_LOSS)
    },
    physicalLocations: [],
    poa5GConfig: {
      cellId: getElemFieldVal(element, FIELD_NR_CELL_ID)
    },
    geoData: !location ? null : {
      location: {
        type: 'Point',
        coordinates: JSON.parse(location)
      },
      radius: (radius === '') ? null : radius
    }
  };

  return nl;
}

export function createPoaWIFI(uniqueId, name, element) {
  var location = getElemFieldVal(element, FIELD_GEO_LOCATION);
  var radius = getElemFieldVal(element, FIELD_GEO_RADIUS);
  var nl = {
    id: uniqueId,
    name: name,
    type: POA_WIFI_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_TERM_LINK_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_TERM_LINK_PKT_LOSS)
    },
    physicalLocations: [],
    poaWifiConfig: {
      macId: getElemFieldVal(element, FIELD_MAC_ID)
    },
    geoData: !location ? null : {
      location: {
        type: 'Point',
        coordinates: JSON.parse(location)
      },
      radius: (radius === '') ? null : radius
    }
  };

  return nl;
}

export function createDefaultNL(zoneName) {
  var nlName = zoneName + '-' + DEFAULT_NL_TYPE_STR;
  var nl = {
    id: nlName,
    name: nlName,
    type: DEFAULT_NL_TYPE_STR,
    netChar: {
      latency: parseInt(DEFAULT_LATENCY_TERMINAL_LINK),
      latencyVariation: parseInt(DEFAULT_LATENCY_JITTER_TERMINAL_LINK),
      throughputDl: parseInt(DEFAULT_THROUGHPUT_DL_TERMINAL_LINK),
      throughputUl: parseInt(DEFAULT_THROUGHPUT_UL_TERMINAL_LINK),
      packetLoss: parseInt(DEFAULT_PACKET_LOSS_TERMINAL_LINK)
    },
    geoData: null,
    physicalLocations: []
  };
  return nl;
}

export function createPL(uniqueId, name, type, element) {
  var location = getElemFieldVal(element, FIELD_GEO_LOCATION);
  var wireless = getElemFieldVal(element, FIELD_WIRELESS);
  var pl = {
    id: uniqueId,
    name: name,
    type: type,
    isExternal: getElemFieldVal(element, FIELD_IS_EXTERNAL),
    connected: getElemFieldVal(element, FIELD_CONNECTED),
    wireless: wireless,
    wirelessType: wireless ? getElemFieldVal(element, FIELD_WIRELESS_TYPE) : '',
    netChar: {
      latency: getElemFieldVal(element, FIELD_LINK_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_LINK_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_LINK_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_LINK_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_LINK_PKT_LOSS)
    },
    dataNetwork: {
      dnn: getElemFieldVal(element, FIELD_DN_NAME),
      ladn: getElemFieldVal(element, FIELD_DN_LADN),
      ecsp: getElemFieldVal(element, FIELD_DN_ECSP)
    },
    geoData: !location ? null : {
      location: {
        type: 'Point',
        coordinates: JSON.parse(location)
      }
    },
    macId: getElemFieldVal(element, FIELD_UE_MAC_ID),
    processes: []
  };

  var path = getElemFieldVal(element, FIELD_GEO_PATH);
  if (path && pl.geoData) {
    pl.geoData.path = {
      type: 'LineString',
      coordinates: JSON.parse(path)
    };
    pl.geoData.eopMode = getElemFieldVal(element, FIELD_GEO_EOP_MODE);
    const velocity = getElemFieldVal(element, FIELD_GEO_VELOCITY);
    pl.geoData.velocity = velocity ? velocity : null;
  }

  return pl;
}

export function createZone(uniqueId, name, element) {
  var zone = {
    id: uniqueId,
    name: name,
    type: ZONE_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_INTRA_ZONE_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_INTRA_ZONE_LATENCY_VAR),
      throughputDl: getElemFieldVal(element, FIELD_INTRA_ZONE_THROUGHPUT_DL),
      throughputUl: getElemFieldVal(element, FIELD_INTRA_ZONE_THROUGHPUT_UL),
      packetLoss: getElemFieldVal(element, FIELD_INTRA_ZONE_PKT_LOSS)
    },
    networkLocations: [createDefaultNL(name)],
    meta: {}
  };
  zone.meta[META_DISPLAY_MAP_COLOR] = getElemFieldVal(element, FIELD_META_DISPLAY_MAP_COLOR);
  return zone;
}

export function createDefaultZone(domainName) {
  var zoneName = domainName + '-' + COMMON_ZONE_TYPE_STR;
  var zone = {
    id: zoneName,
    name: zoneName,
    type: COMMON_ZONE_TYPE_STR,
    netChar: {
      latency: parseInt(DEFAULT_LATENCY_INTRA_ZONE),
      latencyVariation: parseInt(DEFAULT_LATENCY_JITTER_INTRA_ZONE),
      throughputDl: parseInt(DEFAULT_THROUGHPUT_DL_INTRA_ZONE),
      throughputUl: parseInt(DEFAULT_THROUGHPUT_UL_INTRA_ZONE),
      packetLoss: parseInt(DEFAULT_PACKET_LOSS_INTRA_ZONE)
    },
    networkLocations: [createDefaultNL(zoneName)]
  };
  return zone;
}

// Find the provided element in the scenario
export function getElementFromScenario(scenario, elementId) {
  // Create new element to be populated with scenario data
  var elem = createElem(elementId);

  // Check if scenario deployment is being requested
  if (scenario.name === elementId) {
    setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_SCENARIO);
    if (scenario.deployment.netChar) {
      setElemFieldVal(elem, FIELD_INT_DOM_LATENCY, scenario.deployment.netChar.latency || 0);
      setElemFieldVal(elem, FIELD_INT_DOM_LATENCY_VAR, scenario.deployment.netChar.latencyVariation || 0);
      setElemFieldVal(elem, FIELD_INT_DOM_LATENCY_DIST, scenario.deployment.netChar.latencyDistribution || DEFAULT_LATENCY_DISTRIBUTION_INTER_DOMAIN);
      setElemFieldVal(elem, FIELD_INT_DOM_THROUGHPUT_DL, scenario.deployment.netChar.throughputDl || 0);
      setElemFieldVal(elem, FIELD_INT_DOM_THROUGHPUT_UL, scenario.deployment.netChar.throughputUl || 0);
      setElemFieldVal(elem, FIELD_INT_DOM_PKT_LOSS, scenario.deployment.netChar.packetLoss || 0);
    }
    if (scenario.deployment.connectivity) {
      setElemFieldVal(elem, FIELD_CONNECTIVITY_MODEL, scenario.deployment.connectivity.model || DEFAULT_CONNECTIVITY_MODEL);
    }
    return elem;
  }

  // Loop through scenario until element is found
  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];
    if (domain.id === elementId) {

      switch (domain.type) {
      case DOMAIN_TYPE_STR:
        setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_OPERATOR);
        break;
      case DOMAIN_CELL_TYPE_STR:
        setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_OPERATOR_CELL);
        break;
      default:
        break;
      }

      setElemFieldVal(elem, FIELD_NAME, domain.name);
      setElemFieldVal(elem, FIELD_PARENT, scenario.name);
      if (domain.netChar) {
        setElemFieldVal(elem, FIELD_INT_ZONE_LATENCY, domain.netChar.latency || 0);
        setElemFieldVal(elem, FIELD_INT_ZONE_LATENCY_VAR, domain.netChar.latencyVariation || 0);
        setElemFieldVal(elem, FIELD_INT_ZONE_THROUGHPUT_DL, domain.netChar.throughputDl || 0);
        setElemFieldVal(elem, FIELD_INT_ZONE_THROUGHPUT_UL, domain.netChar.throughputUl || 0);
        setElemFieldVal(elem, FIELD_INT_ZONE_PKT_LOSS, domain.netChar.packetLoss || 0);
      }
      //only valid for OPERATOR_CELL
      if (domain.cellularDomainConfig) {
        setElemFieldVal(elem, FIELD_MCC, domain.cellularDomainConfig.mcc);
        setElemFieldVal(elem, FIELD_MNC, domain.cellularDomainConfig.mnc);
        setElemFieldVal(elem, FIELD_DEFAULT_CELL_ID, domain.cellularDomainConfig.defaultCellId);
      }

      return elem;
    }

    for (var j in domain.zones) {
      var zone = domain.zones[j];
      if (zone.id === elementId) {
        setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_ZONE);
        setElemFieldVal(elem, FIELD_NAME, zone.name);
        setElemFieldVal(elem, FIELD_PARENT, domain.type === PUBLIC_DOMAIN_TYPE_STR ? scenario.name : domain.name);

        if (zone.netChar) {
          setElemFieldVal(elem, FIELD_INTRA_ZONE_LATENCY, zone.netChar.latency || 0);
          setElemFieldVal(elem, FIELD_INTRA_ZONE_LATENCY_VAR, zone.netChar.latencyVariation || 0);
          setElemFieldVal(elem, FIELD_INTRA_ZONE_THROUGHPUT_DL, zone.netChar.throughputDl || DEFAULT_THROUGHPUT_DL_INTRA_ZONE);
          setElemFieldVal(elem, FIELD_INTRA_ZONE_THROUGHPUT_UL, zone.netChar.throughputUl || DEFAULT_THROUGHPUT_UL_INTRA_ZONE);
          setElemFieldVal(elem, FIELD_INTRA_ZONE_PKT_LOSS, zone.netChar.packetLoss || 0);
        }
        if (zone.meta) {
          setElemFieldVal(elem, FIELD_META_DISPLAY_MAP_COLOR, zone.meta[META_DISPLAY_MAP_COLOR]);
        }
        return elem;
      }

      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];
        if (nl.id === elementId) {
          switch (nl.type) {
          case POA_TYPE_STR:
            setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_POA);
            break;
          case POA_4G_TYPE_STR:
            setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_POA_4G);
            break;
          case POA_5G_TYPE_STR:
            setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_POA_5G);
            break;
          case POA_WIFI_TYPE_STR:
            setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_POA_WIFI);
            break;
          default:
            break;
          }

          setElemFieldVal(elem, FIELD_NAME, nl.name);
          setElemFieldVal(elem, FIELD_PARENT, domain.type === PUBLIC_DOMAIN_TYPE_STR ?
            scenario.name : zone.type === COMMON_ZONE_TYPE_STR ? domain.name : zone.name);
          if (nl.netChar) {
            setElemFieldVal(elem, FIELD_TERM_LINK_LATENCY, nl.netChar.latency || 0);
            setElemFieldVal(elem, FIELD_TERM_LINK_LATENCY_VAR, nl.netChar.latencyVariation || 0);
            setElemFieldVal(elem, FIELD_TERM_LINK_THROUGHPUT_DL, nl.netChar.throughputDl || 0);
            setElemFieldVal(elem, FIELD_TERM_LINK_THROUGHPUT_UL, nl.netChar.throughputUl || 0);
            setElemFieldVal(elem, FIELD_TERM_LINK_PKT_LOSS, nl.netChar.packetLoss || 0);
          }
          //only valid for specific POAs
          if (nl.poa4GConfig) {
            setElemFieldVal(elem, FIELD_CELL_ID, nl.poa4GConfig.cellId || '');
          }
          if (nl.poa5GConfig) {
            setElemFieldVal(elem, FIELD_NR_CELL_ID, nl.poa5GConfig.cellId || '');
          }
          if (nl.poaWifiConfig) {
            setElemFieldVal(elem, FIELD_MAC_ID, nl.poaWifiConfig.macId || '');
          }

          if (nl.geoData) {
            if (nl.geoData.location) {
              setElemFieldVal(elem, FIELD_GEO_LOCATION, JSON.stringify(nl.geoData.location.coordinates) || '');
            }
            setElemFieldVal(elem, FIELD_GEO_RADIUS, nl.geoData.radius || '');
          }
          if (nl.meta) {
            setElemFieldVal(elem, FIELD_META_DISPLAY_MAP_ICON, nl.meta[META_DISPLAY_MAP_ICON]);
          }
          return elem;
        }

        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];
          if (pl.id === elementId) {
            switch (pl.type) {
            case UE_TYPE_STR:
              setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_UE);
              break;
            case FOG_TYPE_STR:
              setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_FOG);
              break;
            case EDGE_TYPE_STR:
              setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_EDGE);
              break;
            case CN_TYPE_STR:
              setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_CN);
              break;
            case DC_TYPE_STR:
              setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_DC);
              break;
            default:
              break;
            }

            setElemFieldVal(elem, FIELD_NAME, pl.name);
            setElemFieldVal(elem, FIELD_PARENT, domain.type === PUBLIC_DOMAIN_TYPE_STR ? scenario.name :
              zone.type === COMMON_ZONE_TYPE_STR ? domain.name :
                nl.type === DEFAULT_NL_TYPE_STR ? zone.name : nl.name);
            if (pl.netChar) {
              setElemFieldVal(elem, FIELD_LINK_LATENCY, pl.netChar.latency || DEFAULT_LATENCY_LINK);
              setElemFieldVal(elem, FIELD_LINK_LATENCY_VAR, pl.netChar.latencyVariation || DEFAULT_LATENCY_JITTER_LINK);
              setElemFieldVal(elem, FIELD_LINK_THROUGHPUT_DL, pl.netChar.throughputDl || DEFAULT_THROUGHPUT_DL_LINK);
              setElemFieldVal(elem, FIELD_LINK_THROUGHPUT_UL, pl.netChar.throughputUl || DEFAULT_THROUGHPUT_UL_LINK);
              setElemFieldVal(elem, FIELD_LINK_PKT_LOSS, pl.netChar.packetLoss || DEFAULT_PACKET_LOSS_LINK);
            }
            setElemFieldVal(elem, FIELD_IS_EXTERNAL, pl.isExternal || false);
            setElemFieldVal(elem, FIELD_CONNECTED, pl.connected || false);
            setElemFieldVal(elem, FIELD_WIRELESS, pl.wireless || false);
            setElemFieldVal(elem, FIELD_WIRELESS_TYPE, pl.wirelessType || '');
            setElemFieldVal(elem, FIELD_UE_MAC_ID, pl.macId || '');

            if (pl.dataNetwork) {
              setElemFieldVal(elem, FIELD_DN_NAME, pl.dataNetwork.dnn || '');
              setElemFieldVal(elem, FIELD_DN_LADN, pl.dataNetwork.ladn || false);
              setElemFieldVal(elem, FIELD_DN_ECSP, pl.dataNetwork.ecsp || '');
            }

            if (pl.geoData) {
              if (pl.geoData.location) {
                setElemFieldVal(elem, FIELD_GEO_LOCATION, JSON.stringify(pl.geoData.location.coordinates) || '');
              }
              if (pl.geoData.path) {
                setElemFieldVal(elem, FIELD_GEO_PATH, JSON.stringify(pl.geoData.path.coordinates) || '');
              }
              setElemFieldVal(elem, FIELD_GEO_EOP_MODE, pl.geoData.eopMode || '');
              setElemFieldVal(elem, FIELD_GEO_VELOCITY, pl.geoData.velocity || '');
            }
            if (pl.meta) {
              setElemFieldVal(elem, FIELD_META_DISPLAY_MAP_ICON, pl.meta[META_DISPLAY_MAP_ICON]);
            }
            return elem;
          }

          for (var m in pl.processes) {
            var process = pl.processes[m];
            if (process.id === elementId) {
              switch (process.type) {
              case MEC_SVC_TYPE_STR:
                setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_MECSVC);
                break;
              case UE_APP_TYPE_STR:
                setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_UE_APP);
                break;
              case EDGE_APP_TYPE_STR:
                setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_EDGE_APP);
                break;
              case CLOUD_APP_TYPE_STR:
                setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_CLOUD_APP);
                break;
              default:
                break;
              }
              setElemFieldVal(elem, FIELD_PARENT, pl.name);
              setElemFieldVal(elem, FIELD_NAME, process.name);

              if (process.netChar) {
                setElemFieldVal(elem, FIELD_APP_LATENCY, process.netChar.latency || 0);
                setElemFieldVal(elem, FIELD_APP_LATENCY_VAR, process.netChar.latencyVariation || 0);
                setElemFieldVal(elem, FIELD_APP_THROUGHPUT_DL, process.netChar.throughputDl || DEFAULT_THROUGHPUT_DL_APP);
                setElemFieldVal(elem, FIELD_APP_THROUGHPUT_UL, process.netChar.throughputUl || DEFAULT_THROUGHPUT_UL_APP);
                setElemFieldVal(elem, FIELD_APP_PKT_LOSS, process.netChar.packetLoss || 0);
              }
              if (process.userChartLocation) {
                setElemFieldVal(elem, FIELD_CHART_ENABLED, true);
                setElemFieldVal(elem, FIELD_CHART_LOC, process.userChartLocation || '');
                setElemFieldVal(elem, FIELD_CHART_VAL, process.userChartAlternateValues || '');
                setElemFieldVal(elem, FIELD_CHART_GROUP, process.userChartGroup || '');
              } else {
                setElemFieldVal(elem, FIELD_IMAGE, process.image || '');
                setElemFieldVal(elem, FIELD_ENV_VAR, process.environment || '');
                setElemFieldVal(elem, FIELD_CMD, process.commandExe || '');
                setElemFieldVal(elem, FIELD_CMD_ARGS, process.commandArguments || '');
                setElemFieldVal(elem, FIELD_IS_EXTERNAL, process.isExternal || false);
                setElemFieldVal(elem, FIELD_PLACEMENT_ID, process.placementId || '');

                if (process.serviceConfig) {
                  setElemFieldVal(elem, FIELD_PORT, process.serviceConfig.ports[0].port || '');
                  setElemFieldVal(elem, FIELD_PROTOCOL, process.serviceConfig.ports[0].protocol || '');
                  setElemFieldVal(elem, FIELD_GROUP, process.serviceConfig.meSvcName || '');
                  setElemFieldVal(elem, FIELD_EXT_PORT, process.serviceConfig.ports[0].externalPort || '');
                }

                if (process.gpuConfig) {
                  setElemFieldVal(elem, FIELD_GPU_COUNT, process.gpuConfig.count || '');
                  setElemFieldVal(elem, FIELD_GPU_TYPE, process.gpuConfig.type || '');
                }

                if (process.cpuConfig) {
                  setElemFieldVal(elem, FIELD_CPU_MIN, process.cpuConfig.min || '');
                  setElemFieldVal(elem, FIELD_CPU_MAX, process.cpuConfig.max || '');
                }

                if (process.memoryConfig) {
                  setElemFieldVal(elem, FIELD_MEMORY_MIN, process.memoryConfig.min || '');
                  setElemFieldVal(elem, FIELD_MEMORY_MAX, process.memoryConfig.max || '');
                }
              }

              if (process.externalConfig) {
                if (process.externalConfig.ingressServiceMap) {
                  setElemFieldVal(elem, FIELD_INGRESS_SVC_MAP,
                    getIngressServiceMapStr(process.externalConfig.ingressServiceMap));
                }
                if (process.externalConfig.egressServiceMap) {
                  setElemFieldVal(elem, FIELD_EGRESS_SVC_MAP,
                    getEgressServiceMapStr(process.externalConfig.egressServiceMap));
                }
                setElemFieldVal(elem, FIELD_PLACEMENT_ID, process.placementId || '');
              }
              return elem;
            }
          }
        }
      }
    }
  }
}

function createTooltip(title) {
  var tooltip = document.createElement('div');
  tooltip.style.padding = '10px';
  if (title) {
    tooltip.innerHTML = '<b>' + title + '</b>';
  }
  return tooltip;
}
function addTitle(tooltip, title) {
  tooltip.innerHTML += '<b>' + title + '</b>';
}
function addName(tooltip, name) {
  tooltip.innerHTML += '<br>id: ' + name;
}
function addType(tooltip, type) {
  tooltip.innerHTML += '<br>type: ' + type;
}
function addNetChar(tooltip, netChar) {
  tooltip.innerHTML += '<br>latency: ' + (netChar.latency || 0) + ' ms';
  tooltip.innerHTML += '<br>jitter: ' + (netChar.latencyVariation || '0')  + ' ms';
  tooltip.innerHTML += '<br>packet loss: ' + (netChar.packetLoss || '0') + ' %';
  tooltip.innerHTML += '<br>UL throughput: ' + (netChar.throughputUl || 0) + ' mb/s';
  tooltip.innerHTML += '<br>DL throughput: ' + (netChar.throughputDl || 0) + ' mb/s';
}
function addConnectivityModel(tooltip, model) {
  tooltip.innerHTML += '<br>connectivity model: ' + (model || '');
}
function addCellId(tooltip, cellId) {
  tooltip.innerHTML += '<br>cell id: ' + (cellId || '');
}
function addMacAddress(tooltip, mac) {
  tooltip.innerHTML += '<br>mac: ' + (mac || '');
}
function addConnectionState(tooltip, state) {
  tooltip.innerHTML += '<br>state: ' + ((state) ? 'CONNECTED' : 'DISCONNECTED');
}
function addWirelessTypes(tooltip, type) {
  if (type) {
    tooltip.innerHTML += '<br>wireless types: ' + type;
  }
}

// Add scenario node
export function addScenarioNode(scenario, nodes) {
  var nodeTooltip = createTooltip('Node Configuration');
  addName(nodeTooltip, 'Internet');
  addType(nodeTooltip, 'scenario');
  addConnectivityModel(nodeTooltip, scenario.deployment.connectivity.model);

  var n = {
    id: scenario.name,
    name: scenario.name,
    title: nodeTooltip,
    label: 'Internet',
    level: 0
  };

  var image = getScenarioSpecificImage(n.label, TYPE_SCENARIO, null);
  if (image) {
    n['size'] = 40;
    n['image'] = image;
    n['group'] = 'scenarioImg';
  } else {
    n['group'] = 'scenario';
  }

  nodes.push(n);
}

// Add domain node
export function addDomainNode(domain, parent, nodes, edges) {
  var nodeTooltip = createTooltip('Node Configuration');
  addName(nodeTooltip, domain.name);
  addType(nodeTooltip, domain.type);

  var n = {
    id: domain.id,
    name: domain.name,
    title: nodeTooltip,
    label: domain.name,
    level: 1
  };

  var edgeTooltip = createTooltip('Link Configuration');
  addType(edgeTooltip, 'inter-domain');
  addNetChar(edgeTooltip, parent.deployment.netChar);

  var e = {
    title: edgeTooltip,
    label: (parent.deployment.netChar ? parent.deployment.netChar.latency || 0 : 0) + ' ms',
    from: parent.name,
    to: domain.id
  };

  var image = getScenarioSpecificImage(n.label, TYPE_DOMAIN, null);
  if (image) {
    n['size'] = 30;
    n['image'] = image;
    n['group'] = 'domainImg';
  } else {
    n['group'] = 'domain';
  }

  nodes.push(n);
  edges.push(e);
}

// Add zone node
export function addZoneNode(zone, parent, nodes, edges) {
  var nodeTooltip = createTooltip('Node Configuration');
  addName(nodeTooltip, zone.name);
  addType(nodeTooltip, zone.type);

  var n = {
    id: zone.id,
    name: zone.name,
    title: nodeTooltip,
    label: zone.name,
    level: 2
  };

  var edgeTooltip = createTooltip('Link Configuration');
  addType(edgeTooltip, 'inter-zone');
  addNetChar(edgeTooltip, parent.netChar);

  var e = {
    title: edgeTooltip,
    label: (parent.netChar ? parent.netChar.latency || 0 : 0) + ' ms',
    from: parent.id,
    to: zone.id,
    color: {
      color: '#606060',
      highlight: '#606060',
      hover: '#606060'
    }
  };

  var image = getScenarioSpecificImage(n.label, TYPE_ZONE, null);
  if (image) {
    n['size'] = 30;
    n['image'] = image;
    n['group'] = 'zoneImg';
  } else {
    n['group'] = 'zone';
  }

  nodes.push(n);
  edges.push(e);
}

// Add network location node
export function addNlNode(nl, parent, nodes, edges) {
  var nodeTooltip = createTooltip('Node Configuration');
  addName(nodeTooltip, nl.name);
  addType(nodeTooltip, nl.type);

  var n = {
    id: nl.id,
    name: nl.name,
    title: nodeTooltip,
    label: nl.name,
    level: 3
  };

  var edgeTooltip = createTooltip('Link Configuration');
  addType(edgeTooltip, 'intra-zone');
  addNetChar(edgeTooltip, parent.netChar);

  var e = {
    title: edgeTooltip,
    label: (parent.netChar ? parent.netChar.latency || 0 : 0) + ' ms',
    from: parent.id,
    to: nl.id
  };

  var image = getScenarioSpecificImage(n.label, TYPE_NET_LOC, null);
  if (image) {
    n['size'] = 30;
    n['image'] = image;
    n['group'] = 'nLocPoaImg';
  } else {
    n['group'] = 'nLocPoa';
  }

  // Set level and group based on PL type
  switch (nl.type) {
  case POA_4G_TYPE_STR:
    addCellId(nodeTooltip, nl.poa4GConfig.cellId);
    break;
  case POA_5G_TYPE_STR:
    addCellId(nodeTooltip, nl.poa5GConfig.cellId);
    break;
  case POA_WIFI_TYPE_STR:
    addMacAddress(nodeTooltip, nl.poaWifiConfig.macId);
    break;
  default:
    break;
  }

  nodes.push(n);
  edges.push(e);
}

// Add physical location node
export function addPlNode(pl, parent, nodes, edges) {
  var nodeTooltip = createTooltip('Node Configuration');
  addName(nodeTooltip, pl.name);
  addType(nodeTooltip, (pl.type === UE_TYPE_STR) ? 'TERMINAL' : pl.type);

  var n = {
    id: pl.id,
    name: pl.name,
    title: nodeTooltip,
    label: pl.name
  };

  var edgeTooltip = null;
  
  var e = {
    from: parent.id,
    to: pl.id
  };

  //the parent of a distant cloud is the scenario, which has no id, only a name
  if(pl.type === DC_TYPE_STR) {
    e.from = parent.name;
  }

  var latency = (parent.netChar) ? parent.netChar.latency || 0 : 0;
  var lineColor = (pl.connected) ? '#606060' : '#FF0000';
  e['color'] = {
    color: lineColor,
    highlight: lineColor,
    hover: lineColor
  };
  e['dashes'] = pl.wireless || false;

  // Set level and group based on PL type
  switch (pl.type) {
  case FOG_TYPE_STR: {
    latency = 0;
    n['level'] = 4;
    
    if (pl.isExternal) {
      n['group'] = 'pLocExtFog';
    } else {
      const image = getScenarioSpecificImage(n.label, TYPE_PHY_LOC, pl.type);
      if (image) {
        n['image'] = image;
        n['group'] = 'pLocIntFogImg';
      } else {
        n['group'] = 'pLocIntFog';
      }
    }
    break;
  }
  case EDGE_TYPE_STR: {
    edgeTooltip = createTooltip('Link Configuration');
    addType(edgeTooltip, 'intra-zone');
    addConnectionState(edgeTooltip, pl.connected);
    addNetChar(edgeTooltip, parent.netChar);

    n['level'] = 3;

    if (pl.isExternal) {
      n['group'] = 'pLocExtEdge';
    } else {
      const image = getScenarioSpecificImage(n.label, TYPE_PHY_LOC, pl.type);
      if (image) {
        n['image'] = image;
        n['group'] = 'pLocIntEdgeImg';
      } else {
        n['group'] = 'pLocIntEdge';
      }
    }
    break;
  }

  case UE_TYPE_STR: {
    addWirelessTypes(nodeTooltip, pl.wirelessType);

    edgeTooltip = createTooltip('Link Configuration');
    addType(edgeTooltip, 'terminal-link');
    addConnectionState(edgeTooltip, pl.connected);
    addNetChar(edgeTooltip, parent.netChar);

    n['level'] = 4;
     
    if (pl.isExternal) {
      const image = getScenarioSpecificImage(
        n.label + '-ext',
        TYPE_PHY_LOC,
        pl.type
      );
      if (image) {
        n['image'] = image;
        n['group'] = 'pLocExtUEImg';
      } else {
        n['group'] = 'pLocExtUE';
      }
    } else {
      const image = getScenarioSpecificImage(n.label, TYPE_PHY_LOC, pl.type);
      if (image) {
        n['image'] = image;
        n['group'] = 'pLocIntUEImg';
      } else {
        n['group'] = 'pLocIntUE';
      }
    }
    break;
  }

  case CN_TYPE_STR: {
    n['level'] = 2;
    n['group'] = pl.isExternal ? 'pLocExtCN' : 'pLocIntCN';
    break;
  }

  case DC_TYPE_STR: {
    edgeTooltip = createTooltip('Link Configuration');
    addType(edgeTooltip, 'inter-domain');
    addConnectionState(edgeTooltip, pl.connected);
    addNetChar(edgeTooltip, parent.deployment.netChar);
    latency = (parent.deployment.netChar) ? parent.deployment.netChar.latency || 0 : 0;
    
    n['level'] = -1;

    if (pl.isExternal) {
      n['group'] = 'pLocExtDC';
    } else {
      const image = getScenarioSpecificImage(n.label, TYPE_PHY_LOC, pl.type);
      if (image) {
        n['size'] = 40;
        n['image'] = image;
        n['group'] = 'pLocIntDCImg';
      } else {
        n['group'] = 'pLocIntDC';
      }
    }
    break;
  }

  default:
    break;
  }

  // Set tooltip
  if (edgeTooltip) {
    e['title'] = edgeTooltip;
  }

  // Set latency label
  e['label'] = latency + ' ms';

  nodes.push(n);
  edges.push(e);
}

// Add process node
export function addProcessNode(proc, parent, nodes, edges) {
  var nodeTooltip = createTooltip('Node Configuration');
  addName(nodeTooltip, proc.name);
  addType(nodeTooltip, proc.type);
  addTitle(nodeTooltip, '<br><br>Link Configuration');
  addType(nodeTooltip, 'application');
  addNetChar(nodeTooltip, proc.netChar);

  var n = {
    id: proc.id,
    name: proc.name,
    title: nodeTooltip,
    label: proc.name
  };

  var edgeTooltip = createTooltip('Link Configuration');
  addType(edgeTooltip, 'node');
  addNetChar(edgeTooltip, parent.netChar);

  var e = {
    title: edgeTooltip,
    label: (parent.netChar ? parent.netChar.latency || 0 : 0) + ' ms',
    from: parent.id,
    to: proc.id,
    color: {
      color: '#C0C0C0',
      highlight: '#C0C0C0',
      hover: '#C0C0C0'
    }
  };

  if (proc.type === 'EDGE-APP') {
    n['level'] = parent.type === 'EDGE' ? 4 : 5;
    n['group'] = proc.isExternal ? 'procExtEdgeApp' : 'procIntEdgeApp';
  } else if (proc.type === 'MEC-SVC') {
    n['level'] = parent.type === 'EDGE' ? 4 : 5;
    n['group'] = proc.isExternal ? 'procExtMecSvc' : 'procIntMecSvc';
  } else if (proc.type === 'UE-APP') {
    n['level'] = 5;
    n['group'] = proc.isExternal ? 'procExtUEApp' : 'procIntUEApp';
  } else if (proc.type === 'CLOUD-APP') {
    n['level'] = -2;
    n['group'] = proc.isExternal ? 'procExtCloudApp' : 'procIntCloudApp';
  }

  nodes.push(n);
  edges.push(e);
}

// Retrieve scenario-specific images for visualization
export function getScenarioSpecificImage(label, nodeType, plType) {
  var image = null;

  switch (nodeType) {
  case TYPE_SCENARIO:
    image = cloudBlack;
    break;
  case TYPE_DOMAIN:
    image = cloudOutlineBlack;
    break;
  case TYPE_ZONE:
    if (label === 'zone1') {
      // image = edgeImage;
      image = switchBlue;
    } else {
      image = switchBlue;
    }
    break;
  case TYPE_NET_LOC:
    image = poaImage;
    break;
  case TYPE_PHY_LOC:
    if (plType === UE_TYPE_STR) {
      var tmp = label.toLowerCase();
      if (tmp.includes('display')) {
        image = screenImage;
      } else if (tmp.includes('uav')) {
        image = tmp.includes('-ext') ? droneBlue : droneBlack;
      } else if (tmp.includes('emu-')) {
        image = tmp.includes('-ext') ? droneBlue : droneBlack;
      }
    } else if (plType === DC_TYPE_STR) {
      image = cloudServerBlackBlue;
    } else if (plType === EDGE_TYPE_STR) {
      image = edgeImage;
    } else if (plType === FOG_TYPE_STR) {
      image = fogImage;
    }
    break;
  case TYPE_PROCESS:
    // No supported images yet
    break;
  default:
    break;
  }

  return image;
}

export const getScenarioNodeChildren = node => {
  if (node.collapsed) {
    return null;
  }
  return (
    node.domains ||
    node.zones ||
    node.networkLocations ||
    node.physicalLocations ||
    node.processes
  );
};

export const isApp = node => {
  return (
    node.data.type &&
    (node.data.type === 'EDGE-APP' ||
      node.data.type === 'UE-APP' ||
      node.data.type === 'CLOUD-APP')
  );
};

export const getElementNames = (neType, scenario) => {
  var elementNames = [];
  for (var dInd in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[dInd];
    for (var zInd in domain.zones) {
      var zone = domain.zones[zInd];
      for (var nInd in zone.networkLocations) {
        var nl = zone.networkLocations[nInd];
        for (var plInd in nl.physicalLocations) {
          var pl = nl.physicalLocations[plInd];
          for (var prInd in pl.processes) {
            var pr = pl.processes[prInd];
            if (pr.type === neType) {
              elementNames.push(pr.name);
            }
          }
        }
      }
    }
  }

  return elementNames;
};