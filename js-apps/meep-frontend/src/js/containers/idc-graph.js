/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component }  from 'react';
import { Graph } from 'react-d3-graph';
import ReactDOM from 'react-dom';
import { Button } from '@rmwc/button';
import * as d3 from 'd3';

import {
  getScenarioSpecificImage
} from '../util/scenario-utils';

import { updateObject } from '../util/object-util';
import {
  execChangeTable,
  execChangeVis,
  execVisFilteredData,
  execChangeDisplayedScenario
} from '../state/exec';
import { cfgChangeTable, cfgChangeVis, cfgElemEdit } from '../state/cfg';

import {
  TYPE_CFG,
  TYPE_EXEC
} from '../meep-constants';

import {
  FIELD_NAME,
  getElemFieldVal
} from '../util/elem-utils';

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

const translate = d => {
  return `translate(${d.X}, ${d.Y})`;
};

const curveGenerator = d => {
  return `M${d.X},${d.Y} C${d.X},${d.parent.Y + 150} ${d.parent.X},${d.parent.Y + 50} ${d.parent.X},${d.parent.Y}`;
};

const lineGenerator = d => {
  return `M${d.X},${d.Y} L${d.parent.X},${d.parent.Y}`;
};

class IDCNode extends Component {
  constructor(props) {
    super(props);

    this.state = {
      mouseDown: false,
      dragging: false,
      d: this.props.d
    };
  }

  render() {
    const d = this.props.d;

    const fill = this.highlighted ? 'red' : '#69b3a2';
    const radius = this.highlighted ? 14 : 12;
    const size=30;
    
    return (<g
      key={this.props.key}
      transform={translate(d)}
    >
      <image xlinkHref={`../img/${d.data.iconName}`} height={size} width={size} x={-size/2} y={-size/2}
        r={radius}
        style={{fill: fill}}
        stroke={'black'}
        strokeWidth={3}
        onMouseDown={ (e) => {
          console.log(`mouseDown event ${e}`);
          this.dragging = true;
          this.highlighted = true;

          this.mouseCoords={
            x: e.clientX - e.target.farthestViewportElement.parentNode.offsetLeft,
            y: e.clientY - e.target.farthestViewportElement.parentNode.offsetTop
          };

          this.props.updateParent();
        }}
        onMouseUp={ () => {
          this.dragging = false;
          this.highlighted = false;
        }}
        onClick={ () => {
          d.collapsed = !d.data.collapsed;
          this.props.updateParent();
        }}
        onMouseMove={ (e) => {
          if (!this.dragging) {
            return;
          }
          e.preventDefault();

          const newX = e.clientX - e.target.farthestViewportElement.parentNode.offsetLeft;
          const newY = e.clientY - e.target.farthestViewportElement.parentNode.offsetTop;

          const dx = newX - this.mouseCoords.x;
          const dy = newY - this.mouseCoords.y;

          this.mouseCoords.x = newX;
          this.mouseCoords.y = newY;

          const targetXY = e.currentTarget.parentNode.getAttribute('transform').substr(10).slice(0, -1).split(', ');
          const targetX = Number(targetXY[0]);
          const targetY = Number(targetXY[1]);

          // console.log(`(${d.x}, ${d.y}) -> (${X}, ${Y})`);
          d.X = targetX + dx;
          d.Y = targetY + dy;
        
          this.props.updateParent();
        }}
        onMouseOver={() => {
          this.highlighted = true;
          d.data.dR = 4;
          console.log('onMouseOver');
          this.props.updateParent();
        }}
        onMouseOut={() => {
          d.data.dR = 0;
          this.dragging = false;
          this.highlighted = false;
          this.props.updateParent();
          console.log('onMouseOut');
        }}
      />
      <text x={-size/2} y="35" className="tiny">{d.data.id}</text>
    </g>);
  }
}

const visitNodes = f => node => {
  f(node);
  if (node.children) {
    _.each(node.children, (c) => {
      visitNodes(f)(c);
    });
  }
};

const copyPositions = nodeSrc => nodeDest => {
  if (nodeSrc.X && nodeSrc.Y) {
    nodeDest.X = nodeSrc.X;
    nodeDest.Y = nodeSrc.Y;
    nodeDest.dR = nodeSrc.dR;
    nodeDest.iconName = nodeSrc.iconName;
    nodeDest.collapsed = nodeSrc.collapsed;
    nodeDest.name = nodeSrc.name;
  }
  
  if (nodeSrc.children && nodeDest.children && (nodeSrc.children.length === nodeDest.children.length)) {
    for (let i=0; i<nodeSrc.children.length; i++) {
      const cSrc = nodeSrc.children[i];
      const cDest = nodeDest.children[i];
      copyPositions(cSrc)(cDest);
    }
  }
};

const updateXY = node => {
  if (!node.X) {
    node.X = node.x;
  }
  if (!node.Y) {
    node.Y = node.y;
  }
};

const flipXY = node => {
  const margin = 30;
  const x = node.x;
  const y = node.y;
  node.x = y + margin; 
  node.y = x;
};

const clusterFlip = cluster => root => {
  cluster(root);
  visitNodes(flipXY)(root);
  visitNodes(updateXY)(root);
};

const data = {  
  'children':[  
    {  
      'name':'boss1',
      'children':[  
        {  
          'name':'mister_a',
          'colname':'level3'
        },
        {  
          'name':'mister_b',
          'colname':'level3'
        },
        {  
          'name':'mister_c',
          'colname':'level3'
        },
        {  
          'name':'mister_d',
          'colname':'level3'
        }
      ],
      'colname':'level2'
    },
    {  
      'name':'boss2',
      'children':[  
        {  
          'name':'mister_e',
          'colname':'level3'
        },
        {  
          'name':'mister_f',
          'colname':'level3'
        },
        {  
          'name':'mister_g',
          'colname':'level3'
        },
        {  
          'name':'mister_h',
          'colname':'level3'
        }
      ],
      'colname':'level2'
    }
  ],
  'name':'CEO'
};

const transformData = data => {
  return data;
};

const getChildren = (node) => {
  if (node.collapsed) {
    return null;
  }
  return node.domains || node.zones || node.networkLocations || node.physicalLocations || node.processes;
};

class IDCGraph extends Component {

  constructor(props) {
    super(props);
    this.state = {
      root: null
    };
  }

  componentDidMount() {

    // this.root.siblingIndex = 0;
    // const setSiblingIndices = node => {
    //   if (node.children) {
    //     _.each(node.children, (c, i) => {
    //       c.siblingIndex = i;
    //       setSiblingIndices(c);
    //     });
    //   }
    // };

    // setSiblingIndices(this.root);

    // this.mounted = true;

    // this.root = data;
    // this.setState({data: null});

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
    switch(this.props.type) {
    case TYPE_CFG:
      return this.props.cfgTable;
    case TYPE_EXEC:
      return this.props.execTable;
    default:
      return null;
    }
  }

  changeTable(table) {
    switch(this.props.type) {
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

  update() {
    // this.props.execChangeDisplayedScenario(this.root);
    this.setState({root: this.root});
  }

  positionNodes () {
    const data = this.root; // || this.props.displayedScenario;
    this.root = d3.hierarchy(data, getChildren);

    copyPositions(data)(this.root);

    // Create the cluster layout:
    let cluster = d3.cluster().size([this.props.height, this.props.width - 100]);
    clusterFlip(cluster)(this.root);
  }

  render() {

    if (!this.state.root) {
      this.root = this.props.displayedScenario;
      this.positionNodes();
    }
    
    if (!this.root) {
      return  null;
    }

    const lines = this.root.descendants().slice(1)
      .map((d,i) => <path
        key={'path' + i}
        d={lineGenerator(d)}
        style={{fill: 'none', 'strokeWidth': 2}}
        stroke={'#aaa'}
        className='line'
      />);

    const nodes = this.root.descendants()
      .map((d, i) =>
        <IDCNode
          key={`path${i}`}
          d={d}
          updateParent={() => this.update()}
        />
      );
         
    return (
      <svg
        height={this.props.height}
        width={this.props.width}
      >
      <>
        {lines}
        {nodes}
      </>
      </svg>
    );
  }
}


const mapStateToProps = state => {
  return {
    cfgTable: state.cfg.table,
    execTable: state.exec.table,
    execVis: state.exec.vis,
    cfgVis: state.cfg.vis,
    devMode: state.ui.devMode,
    execVisData: execVisFilteredData(state),
    displayedScenario: state.exec.displayedScenario
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeExecTable: (table) => {dispatch(execChangeTable(table));},
    changeCfgTable: (table) => {dispatch(cfgChangeTable(table));},
    changeExecVis: (vis) => dispatch(execChangeVis(vis)),
    changeCfgVis: (vis) => dispatch(cfgChangeVis(vis)),
    changeCfgElement: (element) => dispatch(cfgElemEdit(element)),
    execChangeDisplayedScenario: (scenario) => dispatch(execChangeDisplayedScenario(scenario))
  };
};

const ConnectedIDCGraph = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCGraph);

export default ConnectedIDCGraph;
