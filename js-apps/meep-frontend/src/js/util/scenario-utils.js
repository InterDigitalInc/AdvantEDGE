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
import * as vis from 'vis';
import { updateObject } from './object-util';

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
  FIELD_INTRA_ZONE_LATENCY,
  FIELD_INTRA_ZONE_LATENCY_VAR,
  FIELD_INTRA_ZONE_THROUGPUT,
  FIELD_INTRA_ZONE_PKT_LOSS,
  FIELD_TERM_LINK_LATENCY,
  FIELD_TERM_LINK_LATENCY_VAR,
  FIELD_TERM_LINK_THROUGPUT,
  FIELD_TERM_LINK_PKT_LOSS,
  FIELD_LINK_LATENCY,
  FIELD_LINK_LATENCY_VAR,
  FIELD_LINK_THROUGPUT,
  FIELD_LINK_PKT_LOSS,
  FIELD_APP_LATENCY,
  FIELD_APP_LATENCY_VAR,
  FIELD_APP_THROUGPUT,
  FIELD_APP_PKT_LOSS,
  createElem,
  getElemFieldVal,
  setElemFieldVal
} from './elem-utils';

import {
  ELEMENT_TYPE_SCENARIO,
  ELEMENT_TYPE_OPERATOR,
  ELEMENT_TYPE_ZONE,
  ELEMENT_TYPE_POA,
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
  DEFAULT_THROUGHPUT_INTER_DOMAIN,
  DEFAULT_PACKET_LOSS_INTER_DOMAIN,
  DEFAULT_LATENCY_INTER_ZONE,
  DEFAULT_LATENCY_JITTER_INTER_ZONE,
  DEFAULT_THROUGHPUT_INTER_ZONE,
  DEFAULT_PACKET_LOSS_INTER_ZONE,
  DEFAULT_LATENCY_INTRA_ZONE,
  DEFAULT_LATENCY_JITTER_INTRA_ZONE,
  DEFAULT_THROUGHPUT_INTRA_ZONE,
  DEFAULT_PACKET_LOSS_INTRA_ZONE,
  DEFAULT_LATENCY_TERMINAL_LINK,
  DEFAULT_LATENCY_JITTER_TERMINAL_LINK,
  DEFAULT_THROUGHPUT_TERMINAL_LINK,
  DEFAULT_PACKET_LOSS_TERMINAL_LINK,
  DEFAULT_LATENCY_LINK,
  DEFAULT_LATENCY_JITTER_LINK,
  DEFAULT_THROUGHPUT_LINK,
  DEFAULT_PACKET_LOSS_LINK,
  DEFAULT_LATENCY_APP,
  DEFAULT_LATENCY_JITTER_APP,
  DEFAULT_THROUGHPUT_APP,
  DEFAULT_PACKET_LOSS_APP,
  // DEFAULT_LATENCY_DC,
  DOMAIN_TYPE_STR,
  PUBLIC_DOMAIN_TYPE_STR,
  ZONE_TYPE_STR,
  COMMON_ZONE_TYPE_STR,
  NL_TYPE_STR,
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

  // Logical Scenario types
  TYPE_SCENARIO,
  TYPE_DOMAIN,
  TYPE_ZONE,
  TYPE_NET_LOC,
  TYPE_PHY_LOC,
  TYPE_PROCESS
} from '../meep-constants';

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
          const parent =
            domain.type === PUBLIC_DOMAIN_TYPE_STR
              ? scenario
              : zone.type === COMMON_ZONE_TYPE_STR
                ? domain
                : zone;
          addNlNode(nl, parent, nodes, edges);
        }

        // Physical Locations
        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];

          // Add Physical Location to graph and table
          const parent =
            domain.type === PUBLIC_DOMAIN_TYPE_STR
              ? scenario
              : zone.type === COMMON_ZONE_TYPE_STR
                ? domain
                : nl.type === DEFAULT_NL_TYPE_STR
                  ? zone
                  : nl;
          addPlNode(pl, parent, nodes, edges);

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
  table.entries = _.map(table.data.nodes, node => {
    var elemFromScenario = getElementFromScenario(scenario, node.id);
    return updateObject(node, elemFromScenario);
  });

  // Update visualization data
  var visData = {};
  visData.nodes = new vis.DataSet(nodes);
  visData.edges = new vis.DataSet(edges);

  return { table: table, visData: visData };
}

// Add network element to scenario
export function addElementToScenario(scenario, element) {
  var scenarioElement;
  var type = getElemFieldVal(element, FIELD_TYPE);
  var name = getElemFieldVal(element, FIELD_NAME);
  var parent = getElemFieldVal(element, FIELD_PARENT);

  // Prepare network element to be added to scenario
  switch (type) {
  case ELEMENT_TYPE_OPERATOR: {
    scenarioElement = createDomain(name, element);
    break;
  }
  case ELEMENT_TYPE_ZONE: {
    scenarioElement = createZone(name, element);
    break;
  }
  case ELEMENT_TYPE_POA: {
    scenarioElement = createNL(name, element);
    break;
  }
  case ELEMENT_TYPE_DC: {
    setElemFieldVal(
      element,
      FIELD_PARENT,
      PUBLIC_DOMAIN_TYPE_STR +
        '-' +
        COMMON_ZONE_TYPE_STR +
        '-' +
        DEFAULT_NL_TYPE_STR
    );
    scenarioElement = createPL(name, DC_TYPE_STR, element);
    break;
  }
  case ELEMENT_TYPE_CN: {
    setElemFieldVal(
      element,
      FIELD_PARENT,
      (parent += '-' + COMMON_ZONE_TYPE_STR + '-' + DEFAULT_NL_TYPE_STR)
    );
    scenarioElement = createPL(name, CN_TYPE_STR, element);
    break;
  }
  case ELEMENT_TYPE_EDGE: {
    setElemFieldVal(
      element,
      FIELD_PARENT,
      (parent += '-' + DEFAULT_NL_TYPE_STR)
    );
    scenarioElement = createPL(name, EDGE_TYPE_STR, element);
    break;
  }
  case ELEMENT_TYPE_FOG: {
    scenarioElement = createPL(name, FOG_TYPE_STR, element);
    break;
  }
  case ELEMENT_TYPE_UE: {
    scenarioElement = createPL(name, UE_TYPE_STR, element);
    break;
  }

  case ELEMENT_TYPE_MECSVC: {
    scenarioElement = createProcess(name, MEC_SVC_TYPE_STR, element);
    break;
  }

  case ELEMENT_TYPE_UE_APP: {
    scenarioElement = createProcess(name, UE_APP_TYPE_STR, element);
    break;
  }

  case ELEMENT_TYPE_EDGE_APP: {
    scenarioElement = createProcess(name, EDGE_APP_TYPE_STR, element);
    break;
  }

  case ELEMENT_TYPE_CLOUD_APP: {
    scenarioElement = createProcess(name, CLOUD_APP_TYPE_STR, element);
    break;
  }
  default: {
    break;
  }
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

  // Find element in scenario
  if (scenario.name === name) {
    scenario.deployment.interDomainLatency = getElemFieldVal(
      element,
      FIELD_INT_DOM_LATENCY
    );
    scenario.deployment.interDomainLatencyVariation = getElemFieldVal(
      element,
      FIELD_INT_DOM_LATENCY_VAR
    );
    scenario.deployment.interDomainThroughput = getElemFieldVal(
      element,
      FIELD_INT_DOM_THROUGPUT
    );
    scenario.deployment.interDomainPacketLoss = getElemFieldVal(
      element,
      FIELD_INT_DOM_PKT_LOSS
    );
    return;
  }

  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];
    if (domain.name === name) {
      domain.interZoneLatency = getElemFieldVal(
        element,
        FIELD_INT_ZONE_LATENCY
      );
      domain.interZoneLatencyVariation = getElemFieldVal(
        element,
        FIELD_INT_ZONE_LATENCY_VAR
      );
      domain.interZoneThroughput = getElemFieldVal(
        element,
        FIELD_INT_ZONE_THROUGPUT
      );
      domain.interZonePacketLoss = getElemFieldVal(
        element,
        FIELD_INT_ZONE_PKT_LOSS
      );
      return;
    }

    for (var j in domain.zones) {
      var zone = domain.zones[j];
      if (zone.name === name) {
        if (zone.netChar) {
          zone.netChar.latency = getElemFieldVal(element, FIELD_INTRA_ZONE_LATENCY);
          zone.netChar.latencyVariation = getElemFieldVal(
            element,
            FIELD_INTRA_ZONE_LATENCY_VAR
          );
          zone.netChar.throughput = getElemFieldVal(
            element,
            FIELD_INTRA_ZONE_THROUGPUT
          );
          zone.netChar.packetLoss = getElemFieldVal(
            element,
            FIELD_INTRA_ZONE_PKT_LOSS
          );
        }
        return;
      }

      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];
        if (nl.name === name) {
          nl.terminalLinkLatency = getElemFieldVal(
            element,
            FIELD_TERM_LINK_LATENCY
          );
          nl.terminalLinkLatencyVariation = getElemFieldVal(
            element,
            FIELD_TERM_LINK_LATENCY_VAR
          );
          nl.terminalLinkThroughput = getElemFieldVal(
            element,
            FIELD_TERM_LINK_THROUGPUT
          );
          nl.terminalLinkPacketLoss = getElemFieldVal(
            element,
            FIELD_TERM_LINK_PKT_LOSS
          );
          return;
        }

        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];
          if (pl.name === name) {
            pl.linkLatency = getElemFieldVal(element, FIELD_LINK_LATENCY);
            pl.linkLatencyVariation = getElemFieldVal(
              element,
              FIELD_LINK_LATENCY_VAR
            );
            pl.linkThroughput = getElemFieldVal(element, FIELD_LINK_THROUGPUT);
            pl.linkPacketLoss = getElemFieldVal(element, FIELD_LINK_PKT_LOSS);
            return;
          }

          for (var m in pl.processes) {
            var process = pl.processes[m];
            if (process.name === name) {
              pl.processes[m] = createProcess(
                process.name,
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
      interDomainLatency: parseInt(DEFAULT_LATENCY_INTER_DOMAIN),
      interDomainLatencyVariation: parseInt(
        DEFAULT_LATENCY_JITTER_INTER_DOMAIN
      ),
      interDomainThroughput: parseInt(DEFAULT_THROUGHPUT_INTER_DOMAIN),
      interDomainPacketLoss: parseInt(DEFAULT_PACKET_LOSS_INTER_DOMAIN),
      domains: name === 'None' ? [] : [createDefaultDomain()]
    }
  };
  return scenario;
}

export function createProcess(name, type, element) {
  var isExternal = getElemFieldVal(element, FIELD_IS_EXTERNAL);
  var port = getElemFieldVal(element, FIELD_PORT);
  var gpuCount = getElemFieldVal(element, FIELD_GPU_COUNT);
  var process = {
    id: name,
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
    externalConfig: null,
    appLatency: parseInt(DEFAULT_LATENCY_APP),
    appLatencyVariation: parseInt(DEFAULT_LATENCY_JITTER_APP),
    appThroughput: parseInt(DEFAULT_THROUGHPUT_APP),
    appPacketLoss: parseInt(DEFAULT_PACKET_LOSS_APP),
    placementId: null
  };

  if (isExternal) {
    process.externalConfig = {
      ingressServiceMap: getIngressServiceMapArray(
        getElemFieldVal(element, FIELD_INGRESS_SVC_MAP)
      ),
      egressServiceMap: getEgressServiceMapArray(
        getElemFieldVal(element, FIELD_EGRESS_SVC_MAP)
      )
    };
    process.placementId = getElemFieldVal(element, FIELD_PLACEMENT_ID);
  } else if (getElemFieldVal(element, FIELD_CHART_ENABLED)) {
    process.userChartLocation = getElemFieldVal(element, FIELD_CHART_LOC);
    process.userChartAlternateValues = getElemFieldVal(
      element,
      FIELD_CHART_VAL
    );
    process.userChartGroup = getElemFieldVal(element, FIELD_CHART_GROUP);
  } else {
    process.image = getElemFieldVal(element, FIELD_IMAGE);
    process.environment = getElemFieldVal(element, FIELD_ENV_VAR);
    process.commandArguments = getElemFieldVal(element, FIELD_CMD_ARGS);
    process.commandExe = getElemFieldVal(element, FIELD_CMD);
    process.serviceConfig =
      isNaN(port) || !port
        ? null
        : {
          name: name,
          meSvcName: getElemFieldVal(element, FIELD_GROUP),
          // TODO -- Add frontend support for multiple ports
          ports: [
            {
              protocol:
                getElemFieldVal(element, FIELD_PROTOCOL) === ''
                  ? null
                  : getElemFieldVal(element, FIELD_PROTOCOL).toUpperCase(),
              port:
                getElemFieldVal(element, FIELD_PORT) === ''
                  ? null
                  : getElemFieldVal(element, FIELD_PORT),
              externalPort:
                getElemFieldVal(element, FIELD_EXT_PORT) === ''
                  ? null
                  : getElemFieldVal(element, FIELD_EXT_PORT)
            }
          ]
        };
    process.gpuConfig =
      isNaN(gpuCount) || !gpuCount
        ? null
        : {
          type:
            getElemFieldVal(element, FIELD_GPU_TYPE) === ''
              ? null
              : getElemFieldVal(element, FIELD_GPU_TYPE).toUpperCase(),
          count: gpuCount
        };
    process.placementId = getElemFieldVal(element, FIELD_PLACEMENT_ID);
  }

  process.appLatency = getElemFieldVal(element, FIELD_APP_LATENCY);
  process.appLatencyVariation = getElemFieldVal(element, FIELD_APP_LATENCY_VAR);
  process.appThroughput = getElemFieldVal(element, FIELD_APP_THROUGPUT);
  process.appPacketLoss = getElemFieldVal(element, FIELD_APP_PKT_LOSS);

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

export function createDomain(name, element) {
  var domain = {
    id: name,
    name: name,
    type: DOMAIN_TYPE_STR,
    interZoneLatency: getElemFieldVal(element, FIELD_INT_ZONE_LATENCY),
    interZoneLatencyVariation: getElemFieldVal(
      element,
      FIELD_INT_ZONE_LATENCY_VAR
    ),
    interZoneThroughput: getElemFieldVal(element, FIELD_INT_ZONE_THROUGPUT),
    interZonePacketLoss: getElemFieldVal(element, FIELD_INT_ZONE_PKT_LOSS),
    zones: [createDefaultZone(name)]
  };
  return domain;
}

export function createDefaultDomain() {
  var domain = {
    id: PUBLIC_DOMAIN_TYPE_STR,
    name: PUBLIC_DOMAIN_TYPE_STR,
    type: PUBLIC_DOMAIN_TYPE_STR,
    interZoneLatency: parseInt(DEFAULT_LATENCY_INTER_ZONE),
    interZoneLatencyVariation: parseInt(DEFAULT_LATENCY_JITTER_INTER_ZONE),
    interZoneThroughput: parseInt(DEFAULT_THROUGHPUT_INTER_ZONE),
    interZonePacketLoss: parseInt(DEFAULT_PACKET_LOSS_INTER_ZONE),
    zones: [createDefaultZone(PUBLIC_DOMAIN_TYPE_STR)]
  };
  return domain;
}

export function createNL(name, element) {
  var nl = {
    id: name,
    name: name,
    type: NL_TYPE_STR,
    terminalLinkLatency: getElemFieldVal(element, FIELD_TERM_LINK_LATENCY),
    terminalLinkLatencyVariation: getElemFieldVal(
      element,
      FIELD_TERM_LINK_LATENCY_VAR
    ),
    terminalLinkThroughput: getElemFieldVal(element, FIELD_TERM_LINK_THROUGPUT),
    terminalLinkPacketLoss: getElemFieldVal(element, FIELD_TERM_LINK_PKT_LOSS),
    physicalLocations: []
  };
  return nl;
}

export function createDefaultNL(zoneName) {
  var nlName = zoneName + '-' + DEFAULT_NL_TYPE_STR;
  var nl = {
    id: nlName,
    name: nlName,
    type: DEFAULT_NL_TYPE_STR,
    terminalLinkLatency: parseInt(DEFAULT_LATENCY_TERMINAL_LINK),
    terminalLinkLatencyVariation: parseInt(
      DEFAULT_LATENCY_JITTER_TERMINAL_LINK
    ),
    terminalLinkThroughput: parseInt(DEFAULT_THROUGHPUT_TERMINAL_LINK),
    terminalLinkPacketLoss: parseInt(DEFAULT_PACKET_LOSS_TERMINAL_LINK),
    physicalLocations: []
  };
  return nl;
}

export function createPL(name, type, element) {
  var pl = {
    id: name,
    name: name,
    type: type,
    isExternal: getElemFieldVal(element, FIELD_IS_EXTERNAL),
    linkLatency: parseInt(DEFAULT_LATENCY_LINK),
    linkLatencyVariation: parseInt(DEFAULT_LATENCY_JITTER_LINK),
    linkThroughput: parseInt(DEFAULT_THROUGHPUT_LINK),
    linkPacketLoss: parseInt(DEFAULT_PACKET_LOSS_LINK),
    processes: []
  };
  pl.linkLatency = getElemFieldVal(element, FIELD_LINK_LATENCY);
  pl.linkLatencyVariation = getElemFieldVal(element, FIELD_LINK_LATENCY_VAR);
  pl.linkThroughput = getElemFieldVal(element, FIELD_LINK_THROUGPUT);
  pl.linkPacketLoss = getElemFieldVal(element, FIELD_LINK_PKT_LOSS);

  return pl;
}

export function createZone(name, element) {
  var zone = {
    id: name,
    name: name,
    type: ZONE_TYPE_STR,
    netChar: {
      latency: getElemFieldVal(element, FIELD_INTRA_ZONE_LATENCY),
      latencyVariation: getElemFieldVal(element, FIELD_INTRA_ZONE_LATENCY_VAR),
      throughput: getElemFieldVal(element, FIELD_INTRA_ZONE_THROUGPUT),
      packetLoss: getElemFieldVal(element, FIELD_INTRA_ZONE_PKT_LOSS)
    },
    networkLocations: [createDefaultNL(name)]
  };
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
      throughput: parseInt(DEFAULT_THROUGHPUT_INTRA_ZONE),
      packetLoss: parseInt(DEFAULT_PACKET_LOSS_INTRA_ZONE)
    },
    networkLocations: [createDefaultNL(zoneName)]
  };
  return zone;
}

// Find the provided element in the scenario
export function getElementFromScenario(scenario, elementName) {
  // Create new element to be populated with scenario data
  var elem = createElem(elementName);

  // Check if scenario deployment is being requested
  if (scenario.name === elementName) {
    setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_SCENARIO);
    setElemFieldVal(
      elem,
      FIELD_INT_DOM_LATENCY,
      scenario.deployment.interDomainLatency || 0
    );
    setElemFieldVal(
      elem,
      FIELD_INT_DOM_LATENCY_VAR,
      scenario.deployment.interDomainLatencyVariation || 0
    );
    setElemFieldVal(
      elem,
      FIELD_INT_DOM_THROUGPUT,
      scenario.deployment.interDomainThroughput || 0
    );
    setElemFieldVal(
      elem,
      FIELD_INT_DOM_PKT_LOSS,
      scenario.deployment.interDomainPacketLoss || 0
    );
    return elem;
  }

  // Loop through scenario until element is found
  for (var i in scenario.deployment.domains) {
    var domain = scenario.deployment.domains[i];
    if (domain.name === elementName) {
      setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_OPERATOR);
      setElemFieldVal(elem, FIELD_PARENT, scenario.name);
      setElemFieldVal(
        elem,
        FIELD_INT_ZONE_LATENCY,
        domain.interZoneLatency || 0
      );
      setElemFieldVal(
        elem,
        FIELD_INT_ZONE_LATENCY_VAR,
        domain.interZoneLatencyVariation || 0
      );
      setElemFieldVal(
        elem,
        FIELD_INT_ZONE_THROUGPUT,
        domain.interZoneThroughput || 0
      );
      setElemFieldVal(
        elem,
        FIELD_INT_ZONE_PKT_LOSS,
        domain.interZonePacketLoss || 0
      );
      return elem;
    }

    for (var j in domain.zones) {
      var zone = domain.zones[j];
      if (zone.name === elementName) {
        setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_ZONE);
        setElemFieldVal(
          elem,
          FIELD_PARENT,
          domain.type === PUBLIC_DOMAIN_TYPE_STR ? scenario.name : domain.name
        );

        if (zone.netChar) {
          setElemFieldVal(elem, FIELD_INTRA_ZONE_LATENCY, zone.netChar.latency || 0);
          setElemFieldVal(
            elem,
            FIELD_INTRA_ZONE_LATENCY_VAR,
            zone.netChar.latencyVariation || 0
          );
          setElemFieldVal(
            elem,
            FIELD_INTRA_ZONE_THROUGPUT,
            zone.netChar.throughput || 0
          );
          setElemFieldVal(
            elem,
            FIELD_INTRA_ZONE_PKT_LOSS,
            zone.netChar.packetLoss || 0
          );
        }
        return elem;
      }

      for (var k in zone.networkLocations) {
        var nl = zone.networkLocations[k];
        if (nl.name === elementName) {
          setElemFieldVal(elem, FIELD_TYPE, ELEMENT_TYPE_POA);
          setElemFieldVal(
            elem,
            FIELD_PARENT,
            domain.type === PUBLIC_DOMAIN_TYPE_STR
              ? scenario.name
              : zone.type === COMMON_ZONE_TYPE_STR
                ? domain.name
                : zone.name
          );
          setElemFieldVal(
            elem,
            FIELD_TERM_LINK_LATENCY,
            nl.terminalLinkLatency || 0
          );
          setElemFieldVal(
            elem,
            FIELD_TERM_LINK_LATENCY_VAR,
            nl.terminalLinkLatencyVariation || 0
          );
          setElemFieldVal(
            elem,
            FIELD_TERM_LINK_THROUGPUT,
            nl.terminalLinkThroughput || 0
          );
          setElemFieldVal(
            elem,
            FIELD_TERM_LINK_PKT_LOSS,
            nl.terminalLinkPacketLoss || 0
          );
          return elem;
        }

        for (var l in nl.physicalLocations) {
          var pl = nl.physicalLocations[l];
          if (pl.name === elementName) {
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
            setElemFieldVal(
              elem,
              FIELD_PARENT,
              domain.type === PUBLIC_DOMAIN_TYPE_STR
                ? scenario.name
                : zone.type === COMMON_ZONE_TYPE_STR
                  ? domain.name
                  : nl.type === DEFAULT_NL_TYPE_STR
                    ? zone.name
                    : nl.name
            );
            setElemFieldVal(elem, FIELD_LINK_LATENCY, pl.linkLatency || 0);
            setElemFieldVal(
              elem,
              FIELD_LINK_LATENCY_VAR,
              pl.linkLatencyVariation || 0
            );
            setElemFieldVal(
              elem,
              FIELD_LINK_THROUGPUT,
              pl.linkThroughput || DEFAULT_THROUGHPUT_LINK
            );
            setElemFieldVal(elem, FIELD_LINK_PKT_LOSS, pl.linkPacketLoss || 0);
            setElemFieldVal(elem, FIELD_IS_EXTERNAL, pl.isExternal || false);
            return elem;
          }

          for (var m in pl.processes) {
            var process = pl.processes[m];
            if (process.name === elementName) {
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

              setElemFieldVal(elem, FIELD_APP_LATENCY, process.appLatency || 0);
              setElemFieldVal(
                elem,
                FIELD_APP_LATENCY_VAR,
                process.appLatencyVariation || 0
              );
              setElemFieldVal(
                elem,
                FIELD_APP_THROUGPUT,
                process.appThroughput || DEFAULT_THROUGHPUT_APP
              );
              setElemFieldVal(
                elem,
                FIELD_APP_PKT_LOSS,
                process.appPacketLoss || 0
              );

              if (process.userChartLocation) {
                setElemFieldVal(elem, FIELD_CHART_ENABLED, true);
                setElemFieldVal(
                  elem,
                  FIELD_CHART_LOC,
                  process.userChartLocation || ''
                );
                setElemFieldVal(
                  elem,
                  FIELD_CHART_VAL,
                  process.userChartAlternateValues || ''
                );
                setElemFieldVal(
                  elem,
                  FIELD_CHART_GROUP,
                  process.userChartGroup || ''
                );
              } else {
                setElemFieldVal(elem, FIELD_IMAGE, process.image || '');
                setElemFieldVal(elem, FIELD_ENV_VAR, process.environment || '');
                setElemFieldVal(elem, FIELD_CMD, process.commandExe || '');
                setElemFieldVal(
                  elem,
                  FIELD_CMD_ARGS,
                  process.commandArguments || ''
                );
                setElemFieldVal(
                  elem,
                  FIELD_IS_EXTERNAL,
                  process.isExternal || false
                );
                setElemFieldVal(
                  elem,
                  FIELD_PLACEMENT_ID,
                  process.placementId || ''
                );

                if (process.serviceConfig) {
                  setElemFieldVal(
                    elem,
                    FIELD_PORT,
                    process.serviceConfig.ports[0].port || ''
                  );
                  setElemFieldVal(
                    elem,
                    FIELD_PROTOCOL,
                    process.serviceConfig.ports[0].protocol || ''
                  );
                  setElemFieldVal(
                    elem,
                    FIELD_GROUP,
                    process.serviceConfig.meSvcName || ''
                  );
                  setElemFieldVal(
                    elem,
                    FIELD_EXT_PORT,
                    process.serviceConfig.ports[0].externalPort || ''
                  );
                }

                if (process.gpuConfig) {
                  setElemFieldVal(
                    elem,
                    FIELD_GPU_COUNT,
                    process.gpuConfig.count || ''
                  );
                  setElemFieldVal(
                    elem,
                    FIELD_GPU_TYPE,
                    process.gpuConfig.type || ''
                  );
                }
              }

              if (process.externalConfig) {
                if (process.externalConfig.ingressServiceMap) {
                  setElemFieldVal(
                    elem,
                    FIELD_INGRESS_SVC_MAP,
                    getIngressServiceMapStr(
                      process.externalConfig.ingressServiceMap
                    )
                  );
                }
                if (process.externalConfig.egressServiceMap) {
                  setElemFieldVal(
                    elem,
                    FIELD_EGRESS_SVC_MAP,
                    getEgressServiceMapStr(
                      process.externalConfig.egressServiceMap
                    )
                  );
                }
                setElemFieldVal(
                  elem,
                  FIELD_PLACEMENT_ID,
                  process.placementId || ''
                );
              }
              return elem;
            }
          }
        }
      }
    }
  }
}

// Add scenario node
export function addScenarioNode(scenario, nodes) {
  var n = {
    id: scenario.name,
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
  var n = {
    id: domain.id,
    label: domain.name,
    level: 1
  };

  var e = {
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

  var latency = parseInt(parent.deployment.interDomainLatency);
  if (!isNaN(latency)) {
    e['label'] = String(latency / 2) + ' ms';
  }

  nodes.push(n);
  edges.push(e);
}

// Add zone node
export function addZoneNode(zone, parent, nodes, edges) {
  var n = {
    id: zone.id,
    label: zone.name,
    level: 2
  };

  var e = {
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

  // var latency = "0";
  // if (latency) {
  //     e["label"] = latency + " ms";
  // }

  nodes.push(n);
  edges.push(e);
}

// Add network location node
export function addNlNode(nl, parent, nodes, edges) {
  var n = {
    id: nl.id,
    label: nl.name,
    level: 3
  };

  var e = {
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

  var latency = (parent.netChar) ? parent.netChar.latency : 0;
  if (latency) {
    e['label'] = latency + ' ms';
  }

  nodes.push(n);
  edges.push(e);
}

// Add physical location node
export function addPlNode(pl, parent, nodes, edges) {
  var n = {
    id: pl.id,
    label: pl.name
  };

  var e = {
    from: parent.name,
    to: pl.id
  };

  var latency = null;

  // Set level and group based on PL type
  switch (pl.type) {
  case FOG_TYPE_STR: {
    // latency = "0";
    e['color'] = {
      color: '#606060',
      highlight: '#606060',
      hover: '#606060'
    };
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
    // latency = "0";
    e['color'] = {
      color: '#606060',
      highlight: '#606060',
      hover: '#606060'
    };
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
    latency = parent.terminalLinkLatency;
    e['dashes'] = true;
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
    // latency = "0";
    n['level'] = 2;
    n['group'] = pl.isExternal ? 'pLocExtCN' : 'pLocIntCN';
    break;
  }

  case DC_TYPE_STR: {
    var interDomainLatency = parseInt(parent.deployment.interDomainLatency);
    if (!isNaN(interDomainLatency)) {
      latency = String(interDomainLatency / 2);
    }
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

  // Set latency label
  if (latency) {
    e['label'] = latency + ' ms';
  }

  nodes.push(n);
  edges.push(e);
}

// Add process node
export function addProcessNode(proc, parent, nodes, edges) {
  var n = {
    id: proc.id,
    label: proc.name
  };

  var e = {
    from: parent.name,
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
