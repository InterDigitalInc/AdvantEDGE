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

import { Toolbar, ToolbarRow, ToolbarSection } from '@rmwc/toolbar';
import { Elevation } from '@rmwc/elevation';
import { Grid, GridCell } from '@rmwc/grid';
import { TabBar, Tab } from '@rmwc/tabs';
import { Typography } from '@rmwc/typography';
import { IconButton } from '@rmwc/icon-button';
import { Menu, MenuItem, MenuSurfaceAnchor } from '@rmwc/menu';

import {
  uiChangeCurrentPage,
  uiChangeUserMenuDisplay,
  uiChangeCurrentTab
} from '@/js/state/ui';

import {
  MEEP_TAB_EXEC,
  MEEP_TAB_MON,
  MEEP_TAB_SET,
  MEEP_TAB_CFG,
  MEEP_TAB_HOME,
  PAGE_EXECUTE,
  PAGE_MONITOR,
  PAGE_SETTINGS,
  PAGE_CONFIGURE,
  PAGE_LOGIN,
  STATUS_SIGNED_IN,
  STATUS_SIGNIN_NOT_SUPPORTED,
  PAGE_LOGIN_INDEX,
  PAGE_CONFIGURE_INDEX,
  PAGE_EXECUTE_INDEX,
  PAGE_MONITOR_INDEX,
  PAGE_SETTINGS_INDEX
} from '@/js/meep-constants';

const CorePodsLed = props => {
  /*eslint-disable */
  const greenLed = require('@/img/green-led.png');
  const redLed = require('@/img/red-led.png');
  /* eslint-enable */
  const tooltipType = props.corePodsRunning ? 'success' : 'error';
  const marginLeft = { marginLeft: -35 };
  return (
    <div>
      { props.signInStatus === STATUS_SIGNED_IN || props.signInStatus === STATUS_SIGNIN_NOT_SUPPORTED ?
        <>
          <a data-tip data-for='led'>
            <img
              src={props.corePodsRunning ? greenLed : redLed}
              height={30}
              width={30}
              style={{ marginRight: 15, marginTop: 7 }}
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
        : null
      }
    </div>
  );
};

class MeepTopBar extends Component {
  constructor(props) {
    super(props);
    /*eslint-disable */
    this.logo = require('@/img/ID-Icon-01-idcc.svg');
    this.advantEdge = require('@/img/AdvantEDGE-logo-NoTagline_White_RGB.png');
    /* eslint-enable */
  }

  componentWillMount() {
    this.props.changeUserMenuDisplay(false);
  }

  handleItemClick(page, tabIndex) {
    this.props.changeUserMenuDisplay(false);
    if (this.props.currentPage !== page) {
      this.props.changeCurrentPage(page);
      this.props.changeTabIndex(tabIndex);
    }
  }

  render() {
    let hideTabs = !(this.props.signInStatus === STATUS_SIGNED_IN ||
      this.props.signInStatus === STATUS_SIGNIN_NOT_SUPPORTED);
    return (
      <div>
        <Toolbar>
          <Elevation z={4}>
            <ToolbarRow>
              <ToolbarSection alignStart style={{display:'contents'}}>
                <img
                  id='idcc-logo'
                  className='idcc-toolbar-menu mdc-top-app-bar__navigation-icon'
                  src={this.logo}
                  alt=''
                  onClick={() => this.handleItemClick(PAGE_LOGIN, PAGE_LOGIN_INDEX)}
                />
                <img id='AdvantEdgeLogo' height={50} src={this.advantEdge} alt='' />
                <Grid>
                  <GridCell span="12">
                    <TabBar
                      className='menu-tabs'
                      activeTabIndex={this.props.activeTabIndex}
                      onActivate={evt => this.props.changeTabIndex(evt.detail.index)}
                    >
                      <GridCell span="2">
                        <Tab
                          data-cy={MEEP_TAB_HOME}
                          style={styles.mdcTab}
                          label="Home"
                          onClick={() => { this.handleItemClick(PAGE_LOGIN, PAGE_LOGIN_INDEX); }}
                        />
                      </GridCell>
                      <GridCell span="2" style={{visibility:hideTabs?'hidden':null}}>
                        <Tab
                          hidden={hideTabs}
                          data-cy={MEEP_TAB_CFG}
                          style={styles.mdcTab}
                          label="Configure"
                          onClick={() => { this.handleItemClick(PAGE_CONFIGURE, PAGE_CONFIGURE_INDEX); }}
                        />
                      </GridCell>
                      <GridCell span="2" style={{visibility:hideTabs?'hidden':null}}>
                        <Tab
                          data-cy={MEEP_TAB_EXEC}
                          style={styles.mdcTab}
                          label="Execute"
                          onClick={() => { this.handleItemClick(PAGE_EXECUTE, PAGE_EXECUTE_INDEX); }}
                        />
                      </GridCell>
                      <GridCell span="2" style={{visibility:hideTabs?'hidden':null}}>
                        <Tab
                          data-cy={MEEP_TAB_MON}
                          style={styles.mdcTab}
                          label="Monitor"
                          onClick={() => { this.handleItemClick(PAGE_MONITOR, PAGE_MONITOR_INDEX); }}
                        />
                      </GridCell>
                      <GridCell span="2" style={{visibility:hideTabs?'hidden':null}}>
                        <Tab
                          data-cy={MEEP_TAB_SET}
                          style={styles.mdcTab}
                          label="Settings"
                          onClick={() => { this.handleItemClick(PAGE_SETTINGS, PAGE_SETTINGS_INDEX); }}
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
                  signInStatus={this.props.signInStatus}
                />
                <GridCell align={'middle'}>
                  { this.props.signInStatus === STATUS_SIGNED_IN ?
                    <MenuSurfaceAnchor style={{ height: 48 }}>
                      <Menu
                        open={this.props.userMenuDisplay}
                        onSelect={() => {}}
                        onClose={() => this.props.changeUserMenuDisplay(false)}
                        anchorCorner={'bottomLeft'}
                        align={'left'}
                        style={{ whiteSpace: 'nowrap', marginTop: 5 }}
                      >
                        <MenuItem>
                          <Typography use="body1">Signed in as <b>{this.props.signInUsername}</b></Typography>
                        </MenuItem>
                        <div style={{ width: '100%', borderTop: '1px solid #e4e4e4'}} />
                        <MenuItem onClick={() => {
                          this.props.onClickSignIn();
                          this.props.changeUserMenuDisplay(false);
                        }}>
                          <Typography use="body1">Sign out</Typography>
                        </MenuItem>
                      </Menu>
                      <IconButton
                        icon="account_circle"
                        className='user-icon'
                        style={styles.icon}
                        onClick={() => this.props.changeUserMenuDisplay(true)}
                      />
                    </MenuSurfaceAnchor>
                    : null
                  }
                </GridCell>
              </ToolbarSection>
            </ToolbarRow>
          </Elevation>
        </Toolbar>
      </div>
    );
  }
}

const styles = {
  mdcTab: {
    fontSize: 15,
    fontFamily: 'Roboto'
  },
  icon: {
    color: '#ffffff',
    padding: 5,
    marginRight: 10
  }
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentPage: page => dispatch(uiChangeCurrentPage(page)),
    changeTabIndex: index => dispatch(uiChangeCurrentTab(index)),
    changeUserMenuDisplay: val => dispatch(uiChangeUserMenuDisplay(val))
  };
};

const mapStateToProps = state => {
  return {
    currentPage: state.ui.page,
    userMenuDisplay: state.ui.userMenuDisplay,
    signInStatus: state.ui.signInStatus,
    signInUsername: state.ui.signInUsername,
    activeTabIndex: state.ui.activeTabIndex
  };
};

const ConnectedMeepTopBar = connect(
  mapStateToProps,
  mapDispatchToProps
)(MeepTopBar);

export default ConnectedMeepTopBar;
