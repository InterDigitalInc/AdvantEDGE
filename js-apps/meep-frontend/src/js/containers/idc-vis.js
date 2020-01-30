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
import { connect } from 'react-redux';
import React, { Component, createRef } from 'react';
import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import * as vis from 'vis';
import { updateObject } from '../util/object-util';
import {
  execChangeTable,
  execChangeVis,
  execVisFilteredData
} from '../state/exec';
import { cfgChangeTable, cfgChangeVis, cfgElemEdit } from '../state/cfg';

import { TYPE_CFG, TYPE_EXEC } from '../meep-constants';

import { FIELD_NAME, getElemFieldVal } from '../util/elem-utils';

function createBoxGroup(groups, name, bgColor) {
  groups[name] = {
    borderWidth: 2,
    font: {
      color: '#ffffff',
      size: 18,
      face: 'verdana'
    },
    shape: 'box',
    size: 21,
    color: {
      background: bgColor,
      border: '#000000'
    },
    shadow: {
      enabled: true,
      color: 'rgba(0,0,0,0.5)',
      x: 6,
      y: 6
    }
  };
}

// Create image group for setting visualization
function createImageGroup(groups, name) {
  groups[name] = {
    borderWidth: 0,
    font: {
      color: '#000000',
      size: 18,
      bold: true,
      face: 'verdana'
    },
    color: {
      background: '#FFFFFF'
    },
    shape: 'image',
    shapeProperties: {
      useBorderWithImage: true
    }
  };
}

// Set fixed to true for provided group
function setFixedGroup(group) {
  group['fixed'] = {
    x: true,
    y: true
  };
}

const visFilters = [
  'nodes',
  'edges',
  'layout',
  'interaction',
  'manipulation',
  'physics'
];

class IDCVis extends Component {
  constructor(props) {
    super(props);
    this.state = {};

    this.thisRef = createRef();
    this.configRef = createRef();
  }

  initializeVisualizationOptions(vis, container) {
    vis.options = {
      //clickToUse:true,
      configure: {
        enabled: false,
        filter: '',
        container: container,
        showButton: true
      },
      layout: {
        hierarchical: {
          enabled: true,
          levelSeparation: 320,
          nodeSpacing: 130,
          treeSpacing: 10,
          parentCentralization: false,
          direction: 'LR',
          sortMethod: 'directed'
        }
      },
      interaction: {
        hideEdgesOnDrag: true,
        hover: true,
        multiselect: true,
        navigationButtons: true
      },
      physics: {
        enabled: false,
        hierarchicalRepulsion: {
          centralGravity: 0
        },
        minVelocity: 0.75,
        solver: 'hierarchicalRepulsion'
      },
      nodes: {
        borderWidth: 2,
        font: {
          color: '#ffffff',
          size: 18,
          face: 'verdana'
        },
        shape: 'box',
        size: 21
      },
      edges: {
        width: 2,
        color: {
          color: '#000000'
        },
        font: {
          color: '#000000',
          background: '#FFFFFF',
          size: 18,
          face: 'verdana'
        }
      },
      groups: {}
    };

    var groups = vis.options.groups;

    // Scenario
    createBoxGroup(groups, 'scenario', '#ffffff');
    groups.scenario.borderWidth = 4;
    groups.scenario.font.color = '#000000';
    groups.scenario.font.size = 24;
    createImageGroup(groups, 'scenarioImg');

    // Domains
    createBoxGroup(groups, 'domain', '#ff3620');
    createImageGroup(groups, 'domainImg');

    // Zones
    createBoxGroup(groups, 'zone', '#ffa032');
    createImageGroup(groups, 'zoneImg');

    // Network locations
    createBoxGroup(groups, 'nLocPoa', '#7f7f7f');
    setFixedGroup(groups.nLocPoa);
    createImageGroup(groups, 'nLocPoaImg');
    setFixedGroup(groups.nLocPoaImg);

    // Physical locations
    createBoxGroup(groups, 'pLocIntUE', '#1d8a5c');
    createImageGroup(groups, 'pLocIntUEImg');
    createBoxGroup(groups, 'pLocExtUE', '#373c42');
    createImageGroup(groups, 'pLocExtUEImg');
    createBoxGroup(groups, 'pLocIntFog', '#3ebcfb');
    createImageGroup(groups, 'pLocIntFogImg');
    createBoxGroup(groups, 'pLocExtFog', '#006af8');

    createBoxGroup(groups, 'pLocIntEdge', '#3ebcfb');
    createImageGroup(groups, 'pLocIntEdgeImg');
    createBoxGroup(groups, 'pLocExtEdge', '#006af8');
    createBoxGroup(groups, 'pLocIntCN', '#ff69f9');
    createBoxGroup(groups, 'pLocExtCN', '#ff00f8');
    createBoxGroup(groups, 'pLocIntDC', '#939393');
    createImageGroup(groups, 'pLocIntDCImg');
    createBoxGroup(groups, 'pLocExtDC', '#252525');

    // Processes
    createBoxGroup(groups, 'procIntUEApp', '#d76a2f');
    createBoxGroup(groups, 'procExtUEApp', '#8e4721');
    createImageGroup(groups, 'procExtUEAppImg');
    createBoxGroup(groups, 'procIntEdgeApp', '#d76a2f');
    createBoxGroup(groups, 'procExtEdgeApp', '#8e4721');
    createBoxGroup(groups, 'procIntMecSvc', '#3ebcfb');
    createBoxGroup(groups, 'procExtMecSvc', '#006af8');
    createBoxGroup(groups, 'procIntCloudApp', '#000000');
    createBoxGroup(groups, 'procExtCloudApp', '#000000');
  }

  componentDidMount() {
    const newVis = updateObject({}, this.getVis());
    if (newVis.data.nodes.length < 1) {
      newVis.data.nodes = [
        {
          id: 'waiting',
          label: 'Waiting for Scenario...',
          title: 'Multi-Access Emulation Platform --- The MEEP AdvantEDGE!',
          level: 0,
          group: 'scenario'
        }
      ];
      newVis.data.edges = [];
    }

    this.initializeVisualizationOptions(newVis, this.configRef.current);

    var domNode = ReactDOM.findDOMNode(this);
    newVis.network = new vis.Network(
      domNode,
      this.props.type === TYPE_CFG ? newVis.data : this.props.execVisData,
      newVis.options
    );

    this.table = updateObject(this.getTable(), { data: newVis.data });
    this.changeVis(newVis);
    this.changeTable(this.table);

    // Configuration Visualization handlers
    if (this.props.type === TYPE_CFG) {
      newVis.network.on('click', obj => {
        if (!this.props.cfgVis.data.nodes.get) {
          return;
        }
        // meep.cfg.vis.reportContainer.innerHTML = "x:" + obj.pointer.canvas.x + ", y:" + obj.pointer.canvas.y;

        var clickedNodes = this.props.cfgVis.data.nodes.get(obj.nodes);

        // Highlight selected nodes in table
        const table = updateObject({}, this.props.cfgTable);
        table.selected = [];
        for (var i = 0; i < clickedNodes.length; i++) {
          var clickedNode = clickedNodes[i];
          table.selected.push(clickedNode.id);
        }

        this.changeTable(table);

        // Open first selected element in element configuration pane
        if (this.props.type === TYPE_CFG) {
          this.props.onEditElement(
            table.selected.length
              ? this.getElementByName(table.entries, table.selected[0])
              : null
          );
        }
      });
    }
  }

  getElementByName(entries, name) {
    for (var i = 0; i < entries.length; i++) {
      if (getElemFieldVal(entries[i], FIELD_NAME) === name) {
        return entries[i];
      }
    }
    return null;
  }

  getTable() {
    switch (this.props.type) {
    case TYPE_CFG:
      return this.props.cfgTable;
    case TYPE_EXEC:
      return this.props.execTable;
    default:
      return null;
    }
  }

  changeTable(table) {
    switch (this.props.type) {
    case TYPE_CFG:
      this.props.changeCfgTable(table);
      break;
    case TYPE_EXEC:
      this.props.changeExecTable(table);
      break;
    default:
      break;
    }
  }

  changeVis(vis) {
    switch (this.props.type) {
    case TYPE_CFG:
      this.props.changeCfgVis(vis);
      break;
    case TYPE_EXEC:
      this.props.changeExecVis(vis);
      break;
    default:
      break;
    }
  }

  getVis() {
    switch (this.props.type) {
    case TYPE_CFG:
      return this.props.cfgVis;
    case TYPE_EXEC:
      return this.props.execVis;
    default:
      return null;
    }
  }

  // Toggle visualization controls
  toggleConfig(filterStr) {
    var vis = this.getVis();
    if (
      vis.showConfig === false ||
      (vis.showConfig === true && vis.options.configure.filter === filterStr)
    ) {
      vis.showConfig = !vis.showConfig;
    }
    vis.options.configure.enabled = vis.showConfig;

    var subOptions;
    switch (vis.options.configure.filter) {
    case 'physics':
      subOptions = vis.network.getOptionsFromConfigurator();
      vis.options.physics = subOptions.physics;
      break;
    case 'manipulation':
      subOptions = vis.network.getOptionsFromConfigurator();
      vis.options.manipulation = subOptions.manipulation;
      break;
    case 'interaction':
      subOptions = vis.network.getOptionsFromConfigurator();
      vis.options.interaction = subOptions.interaction;
      break;
    case 'nodes':
      subOptions = vis.network.getOptionsFromConfigurator();
      vis.options.nodes = subOptions.nodes;
      break;
    case 'edges':
      subOptions = vis.network.getOptionsFromConfigurator();
      vis.options.edges = subOptions.edges;
      break;
    case 'layout':
      subOptions = vis.network.getOptionsFromConfigurator();
      vis.options.layout = subOptions.layout;
      break;
    default:
    }

    vis.options.configure.filter = filterStr;
    vis.network.setOptions(vis.options);
  }

  render() {
    return (
      <>
        <div
          className="vis-network-div"
          ref={this.thisRef}
          data-cy={this.props.cydata}
        >
          Vis Component
        </div>
        <div className="idcc-margin-top">
          {this.props.devMode
            ? _.map(visFilters, filter => {
              return (
                <Button
                  raised
                  style={buttonStyles}
                  key={filter}
                  onClick={() => {
                    this.toggleConfig(filter);
                  }}
                >
                  {filter}
                </Button>
              );
            })
            : null}
          <div className="idcc-margin-top" ref={this.configRef} />
        </div>
      </>
    );
  }
}

const buttonStyles = {
  color: 'white',
  marginRight: 5
};

const mapStateToProps = state => {
  return {
    cfgTable: state.cfg.table,
    execTable: state.exec.table,
    execVis: state.exec.vis,
    cfgVis: state.cfg.vis,
    devMode: state.ui.devMode,
    execVisData: execVisFilteredData(state)
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeExecTable: table => {
      dispatch(execChangeTable(table));
    },
    changeCfgTable: table => {
      dispatch(cfgChangeTable(table));
    },
    changeExecVis: vis => dispatch(execChangeVis(vis)),
    changeCfgVis: vis => dispatch(cfgChangeVis(vis)),
    changeCfgElement: element => dispatch(cfgElemEdit(element))
  };
};

const ConnectedIDCVis = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCVis);

export default ConnectedIDCVis;
