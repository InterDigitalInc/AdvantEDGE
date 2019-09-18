/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import _ from 'lodash';
import React, { Component }  from 'react';

import {
  plusGenerator,
  minusGenerator,
  visitNodes,
  blue
} from './graph-utils';

const translate = d => {
  return `translate(${d.X}, ${d.Y})`;
};

const hideNode = node => {
  node.hidden = true;
};

const showNode = node => {
  node.hidden = false;
};

const hideChildren = node => {
  _.each(node.children, c => {
    visitNodes(hideNode)(c);
  });
};

const showChildren = node => {
  _.each(node.children, c => {
    visitNodes(showNode)(c);
  });
};

const Plus = props => {
  const d = props.d;

  const plusMinus = props.collapsible
    ? (d.collapsed ? plusGenerator : minusGenerator)
    : () => '';
    
  return (
    <path
      width={20}
      height={20}
      d={plusMinus()}
      style={{fill: blue, 'strokeWidth': 2}}
      stroke={blue}
      className='plus'
      onClick={() => {
        d.collapsed = !d.collapsed;
        if (d.collapsed) {
          hideChildren(d);
        } else {
          showChildren(d);
        }
        props.updateParent();
      }}
    />
  );
};

export default class IDCNode extends Component {
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

    const fill = this.highlighted ? '#69b3a2' : '#69b3a2';
    const radius = this.highlighted ? 14 : 12;
    const size=30;

    return (<g
      transform={translate(d)}
    >
      <Plus width={10} height={10} d={d} updateParent={this.props.updateParent}/>
      <circle xlinkHref={`../img/${d.data.iconName}`} height={size} width={size} cx={-size/2 + 15} cy={-size/2 + 15} /*filter={d.selected ? 'url(#filter)' : '' }*/
        r={radius}
        style={{fill: fill}}
        stroke={'black'}
        strokeWidth={3}
        onMouseDown={ (e) => {
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
        onClick={() => {
          d.selected = !d.selected;
          // this.props.updateParent();
          this.props.onClick({node: d});
        }}
        onMouseOver={() => {
          this.highlighted = true;
          d.highlighted = true;
          d.data.dR = 4;
          this.props.updateParent();
        }}
        onMouseOut={() => {
          d.data.dR = 0;
          this.dragging = false;
          this.highlighted = false;
          d.highlighted = false;
          this.props.updateParent();
        }}
      />
      <text x={-size/2} y="35" className="tiny" stroke={this.props.stroke} fontWeight={this.highlighted ? 'bold' : 'normal'}>{d.data.id}</text>
    </g>);
  }
}