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

import { connect } from 'react-redux';
import React, { Component }  from 'react';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { TextField } from '@rmwc/textfield';
import { Checkbox } from '@rmwc/checkbox';
import { Elevation } from '@rmwc/elevation';

import {
  uiSetAutomaticRefresh,
  uiChangeRefreshInterval,
  uiChangeDevMode,
  uiExecChangeShowDashboardConfig
} from '../../state/ui';

import {
  PAGE_SETTINGS,
  SET_EXEC_REFRESH_CHECKBOX,
  SET_EXEC_REFRESH_INT,
  SET_VIS_CFG_CHECKBOX,
  SET_VIS_CFG_LABEL,
  SET_DASHBOARD_CFG_CHECKBOX,
  SET_DASHBOARD_CFG_LABEL
} from '../../meep-constants';

class SettingsPageContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      error: false
    };
  }

  validateInterval(val) {
    if (isNaN(val) || val < 500 || 60000 < val) {
      return false;
    }
    return true;
  }

  handleIntervalChange(val) {
    this.props.changeRefreshInterval(val);
    if (this.validateInterval(val)) {
      this.props.startRefresh();
      this.setState({error: false});
    } else {
      this.props.stopRefresh();
      this.setState({error: true});
    }
  }

  handleCheckboxChange(val) {
    this.props.setAutomaticRefresh(val)
    if (val && this.validateInterval(this.props.refreshInterval)) {
      this.props.startRefresh();
    } else {
      this.props.stopRefresh();
    }
  }

  styles() {
    var styles = {
      interval: {

      },
      errorText: {
        display: 'none'
      },
      errorGridCell: {
        marginTop: -15,
        marginLeft: 25,
        paddingBottom: 10
      }
    };

    if (this.state.error) {
      delete styles.errorText.display;
      styles.errorText.fontSize = 14;
      styles.errorText.color = 'rgb(176, 0, 32)';
    }

    return styles;
  }

  render() {

    if (this.props.page !== PAGE_SETTINGS) {
      return null;
    }

    return (
      <div style={{width: '100%'}}>
        <Grid style={{width: '100%'}}>
          <GridInner>
            <GridCell span={12} style={styles.inner}>

              <Elevation className="component-style" z={2} style={{paddingBottom: 10, marginBottom: 10}}>
                <Grid>
                  <GridCell span={12} style={{paddingLeft: 10, paddingTop: 10}}>
                    <div>
                      <span className="mdc-typography--headline6">Execution: </span>
                    </div>
                  </GridCell>
                </Grid>
                <Grid span={12}>
                  <GridCell span={2}>
                    <div style={{marginTop: 20}}>
                      <Checkbox
                        checked={this.props.automaticRefresh}
                        onChange={e => this.handleCheckboxChange(e.target.checked)}
                        data-cy={SET_EXEC_REFRESH_CHECKBOX}>
                          Automatic refresh:
                      </Checkbox>
                    </div>
                  </GridCell>
                  <GridCell span={2}>
                    <TextField outlined style={this.styles().interval}
                      label="Interval (ms)"
                      onChange={(e) => this.handleIntervalChange(e.target.value)}
                      value={this.props.refreshInterval}
                      disabled={!this.props.automaticRefresh}
                      data-cy={SET_EXEC_REFRESH_INT}
                    />
                  </GridCell>
                  <GridCell span={8}>
                  </GridCell>
                </Grid>

                <Grid>
                  <GridCell span={2}>
                  </GridCell>
                  <GridCell span={2} style={this.styles().errorGridCell}>
                    <p style={this.styles().errorText}>
                      500 &lt; value &lt; 60000
                    </p>
                  </GridCell>
                  <GridCell span={8}>
                  </GridCell>
                </Grid>
              </Elevation>

              <Elevation className="component-style" z={2} style={{paddingBottom: 10, marginBottom: 10}}>
                <Grid>
                  <GridCell span={12} style={{paddingLeft: 10, paddingTop: 10}}>
                    <div>
                      <span className="mdc-typography--headline6">Development: </span>
                    </div>
                  </GridCell>
                </Grid>
                <CheckableSettingItem
                  stateItem={this.props.devMode}
                  changeStateItem={this.props.changeDevMode}
                  stateItemName={SET_VIS_CFG_LABEL}
                  cydata={SET_VIS_CFG_CHECKBOX}
                />
                <CheckableSettingItem
                  stateItem={this.props.showDashboardConfig}
                  changeStateItem={this.props.changeShowDashboardConfig}
                  stateItemName={SET_DASHBOARD_CFG_LABEL}
                  cydata={SET_DASHBOARD_CFG_CHECKBOX}
                />
              </Elevation>

            </GridCell>
          </GridInner>
        </Grid>
      </div>
    );
  }
}

const CheckableSettingItem = ({stateItem, changeStateItem, stateItemName, cydata}) => {
  return (
    <Grid span={12} style={{marginTop: 10}}>
      <GridCell span={12}>
        <div>
          <Checkbox
            checked={stateItem}
            onChange={e => changeStateItem(e.target.checked)}
            data-cy={cydata}
          >
            {stateItemName}
          </Checkbox>
        </div>
      </GridCell>
    </Grid>
  );
};

const styles = {
  headlineGrid: {
    marginBottom: 10
  },
  headline: {
    padding: 10,
    marginBotton: 25
  },
  inner: {
    height: '100%'
  },
  page: {
    height: 1500,
    marginBottom: 10
  },
  cfgTable: {
    marginTop: 20,
    padding: 10
  }
};

const mapStateToProps = state => {
  return {
    automaticRefresh: state.ui.automaticRefresh,
    refreshInterval: state.ui.refreshInterval,
    devMode: state.ui.devMode,
    page: state.ui.page,
    showDashboardConfig: state.ui.showDashboardConfig
  };
};

const mapDispatchToProps = dispatch => {
  return {
    setAutomaticRefresh: (val) => dispatch(uiSetAutomaticRefresh(val)),
    changeRefreshInterval: (val) => dispatch(uiChangeRefreshInterval(val)),
    changeDevMode: (mode) => dispatch(uiChangeDevMode(mode)),
    changeShowDashboardConfig: (show) => dispatch(uiExecChangeShowDashboardConfig(show))
  };
};

const ConnectedSettingsPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(SettingsPageContainer);

export default ConnectedSettingsPageContainer;
