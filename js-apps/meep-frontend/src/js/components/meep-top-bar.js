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
import React, { Component }  from 'react';
import ReactTooltip from 'react-tooltip';
import { connect } from 'react-redux';
import '@/css/meep-controller.scss';

import {
  Toolbar,
  ToolbarRow,
  ToolbarSection
  // ToolbarTitle
} from '@rmwc/toolbar';

import { Elevation } from '@rmwc/elevation';

import { Grid, GridCell } from '@rmwc/grid';

import { TabBar, Tab } from '@rmwc/tabs';

import { uiChangeCurrentPage } from '@/js/state/ui';

import {
  MEEP_TAB_EXEC,
  MEEP_TAB_MON,
  MEEP_TAB_SET,
  MEEP_TAB_CFG,
  PAGE_EXECUTE,
  PAGE_MONITOR,
  PAGE_SETTINGS,
  PAGE_CONFIGURE
} from '@/js/meep-constants';

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

class MeepTopBar extends Component {
  constructor(props) {
    super(props);
    /*eslint-disable */
    this.logo = require('../../img/ID-Icon-01-idcc.svg');
    this.advantEdge = require('../../img/AdvantEDGE-logo-NoTagline_White_RGB.png');
    /* eslint-enable */
  }

  handleItemClick(page) {
    this.props.currentPage !== page ? this.props.changeCurrentPage(page) : '';
  }

  render() {
    return (
      <div>
        <Toolbar>
          <Elevation z={4}>
            <ToolbarRow>
              <ToolbarSection alignStart>
                <img
                  id='idcc-logo'
                  className='idcc-toolbar-menu mdc-top-app-bar__navigation-icon'
                  src={this.logo}
                  alt=''
                  onClick={this.props.toggleMainDrawer}
                />
                <img id='AdvantEdgeLogo' height={50} src={this.advantEdge} alt='' />
                {/* <ToolbarTitle>
                  <span style={styles.title}>{this.props.title}</span>
                </ToolbarTitle> */}
                <Grid>
                  <GridCell span="12">
                    <TabBar>
                      <GridCell span="2">
                        <Tab
                          data-cy={MEEP_TAB_CFG}
                          style={styles.mdcTab}
                          label="Configure"
                          onClick={() => { this.handleItemClick(PAGE_CONFIGURE); }}
                        />
                      </GridCell>
                      <GridCell span="2">
                        <Tab
                          data-cy={MEEP_TAB_EXEC}
                          style={styles.mdcTab}
                          label="Execute"
                          onClick={() => { this.handleItemClick(PAGE_EXECUTE); }}
                        />
                      </GridCell>
                      <GridCell span="2">
                        <Tab
                          data-cy={MEEP_TAB_MON}
                          style={styles.mdcTab}
                          label="Monitor"
                          onClick={() => { this.handleItemClick(PAGE_MONITOR); }}
                        />
                      </GridCell>
                      <GridCell span="2">
                        <Tab
                          data-cy={MEEP_TAB_SET}
                          style={styles.mdcTab}
                          label="Settings"
                          onClick={() => { this.handleItemClick(PAGE_SETTINGS); }}
                        />
                      </GridCell>
                    </TabBar>
                  </GridCell>
                </Grid>
              </ToolbarSection>

              <ToolbarSection alignEnd>
                <CorePodsLed
                  corePodsRunning={this.props.corePodsRunning}
                  corePodsErrors={this.props.corePodsErrors}
                />
              </ToolbarSection>
            </ToolbarRow>
          </Elevation>
        </Toolbar>
      </div>
    );
  }
}

const styles = {
  title: {
    color: 'white',
    fontFamily: 'Gill Sans, Gill Sans MT, Calibri, Trebuchet MS, sans-serif',
    fontSize: 22
  },
  mdcTab: {
    fontSize: 17,
    fontFamily: 'Roboto'
  }
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentPage: page => dispatch(uiChangeCurrentPage(page))
  };
};

const mapStateToProps = state => {
  return {
    currentPage: state.ui.page
  };
};

const ConnectedMeepTopBar = connect(
  mapStateToProps,
  mapDispatchToProps
)(MeepTopBar);

export default ConnectedMeepTopBar;
