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
import React from 'react';
import ReactTooltip from 'react-tooltip';
import {
  Toolbar,
  ToolbarRow,
  ToolbarSection,
  ToolbarTitle
} from '@rmwc/toolbar';

import { Elevation } from '@rmwc/elevation';

const CorePodsLed = props => {
  /*eslint-disable */
  const greenLed = require('../../img/green-led.png');
  const redLed = require('../../img/red-led.png');
  /* eslint-enable */
  const tooltipType = props.corePodsRunning ? 'success' : 'error';
  const marginLeft = { marginLeft: -35 };
  return (
    <>
      <a data-tip data-for='led'>
        <img
          src={props.corePodsRunning ? greenLed : redLed}
          height={30}
          width={30}
          style={{ marginRight: 15 }}
        />
      </a>
      <ReactTooltip
        id='led'
        aria-haspopup='true'
        role='example'
        place='left'
        type={tooltipType}
      >
        <ul style={{ listStyle: 'none' }}>
          {props.corePodsErrors.length ? (
            _.map(props.corePodsErrors, e => {
              return (
                <li key={e.name} style={marginLeft}>
                  {`${e.name}: ${e.status}`}
                </li>
              );
            })
          ) : (
            <span style={marginLeft}>All systems GO!</span>
          )}
        </ul>
      </ReactTooltip>
    </>
  );
};

const MeepTopBar = props => {
  /*eslint-disable */
  const logo = require('../../img/ID-Icon-01-idcc.svg');
  const advantEdge = require('../../img/AdvantEDGE-logo-NoTagline_White_RGB.png');
  /* eslint-enable */
  return (
    <Toolbar>
      <Elevation z={4}>
        <ToolbarRow>
          <ToolbarSection alignStart>
            <img
              id='idcc-logo'
              className='idcc-toolbar-menu mdc-top-app-bar__navigation-icon'
              src={logo}
              alt=''
              onClick={() => {
                props.toggleMainDrawer();
              }}
            />
            <img id='AdvantEdgeLogo' height={50} src={advantEdge} alt='' />
            <ToolbarTitle>
              <span style={titleStyle}>{props.title}</span>
            </ToolbarTitle>
          </ToolbarSection>
          <ToolbarSection alignEnd>
            <CorePodsLed
              corePodsRunning={props.corePodsRunning}
              corePodsErrors={props.corePodsErrors}
            />
          </ToolbarSection>
        </ToolbarRow>
      </Elevation>
    </Toolbar>
  );
};

const titleStyle = {
  color: 'white',
  fontFamily: 'Gill Sans, Gill Sans MT, Calibri, Trebuchet MS, sans-serif',
  fontSize: 22
};

export default MeepTopBar;
