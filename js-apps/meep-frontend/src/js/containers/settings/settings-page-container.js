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
import React, { Component } from 'react';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { TextField } from '@rmwc/textfield';
import { Checkbox } from '@rmwc/checkbox';
import { Elevation } from '@rmwc/elevation';
import { Button } from '@rmwc/button';
import IDConfirmDialog from '../../components/dialogs/id-confirm-dialog';

import {
  meepSetDefaultState
} from '../../state/meep-reducer';

import {
  uiSetAutomaticRefresh,
  uiChangeRefreshInterval,
  uiChangeDevMode,
  uiChangeCurrentDialog
} from '../../state/ui';

import {
  PAGE_SETTINGS,
  SET_EXEC_REFRESH_CHECKBOX,
  SET_EXEC_REFRESH_INT,
  SET_VIS_CFG_CHECKBOX,
  SET_VIS_CFG_LABEL,
  SET_RESET_SETTINGS_BUTTON,
  IDC_DIALOG_CLEAR_UI_CACHE
} from '../../meep-constants';

/*global __VERSION__*/

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
      this.setState({ error: false });
    } else {
      this.props.stopRefresh();
      this.setState({ error: true });
    }
  }

  handleCheckboxChange(val) {
    this.props.setAutomaticRefresh(val);
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

  showDialog(id) {
    this.props.showDialog(id);
  }

  closeDialog() {
    this.showDialog(null);
  }

  render() {
    if (this.props.page !== PAGE_SETTINGS) {
      return null;
    }

    return (
      <div>
        <IDConfirmDialog
          title="Clear UI cache (reset default frontend state)"
          open={this.props.currentDialog === IDC_DIALOG_CLEAR_UI_CACHE}
          onClose={() => {
            this.closeDialog();
          }}
          onSubmit={() => this.props.meepSetDefaultState()}
        />
        <div style={{ width: '100%' }}>
          <Grid style={{ width: '100%' }}>
            <GridInner>
              <GridCell span={12}>
                <Elevation
                  className="component-style"
                  z={2}
                  style={styles.elevation}
                >
                  <div style={styles.headline}>
                    <span className="mdc-typography--headline6">
                      Execution:{' '}
                    </span>
                  </div>
                  <div style={styles.content}>
                    <Grid span={12}>
                      <GridCell span={2}>
                        <div style={{ marginTop: 20 }}>
                          <Checkbox
                            checked={this.props.automaticRefresh}
                            onChange={e =>
                              this.handleCheckboxChange(e.target.checked)
                            }
                            data-cy={SET_EXEC_REFRESH_CHECKBOX}
                          >
                            Automatic refresh:
                          </Checkbox>
                        </div>
                      </GridCell>
                      <GridCell span={10}>
                        <TextField
                          outlined
                          style={this.styles().interval}
                          label="Interval (ms)"
                          onChange={e => this.handleIntervalChange(e.target.value)}
                          value={this.props.refreshInterval}
                          disabled={!this.props.automaticRefresh}
                          data-cy={SET_EXEC_REFRESH_INT}
                        />
                      </GridCell>
                    </Grid>

                    <Grid>
                      <GridCell span={2}></GridCell>
                      <GridCell span={10} style={this.styles().errorGridCell}>
                        <p style={this.styles().errorText}>
                          500 &lt; value &lt; 60000
                        </p>
                      </GridCell>
                    </Grid>
                  </div>
                </Elevation>

                <Elevation
                  className="component-style"
                  z={2}
                  style={styles.elevation}
                >
                  <div style={styles.headline}>
                    <span className="mdc-typography--headline6">
                      Development:{' '}
                    </span>
                  </div>
                  <div style={styles.content}>
                    <CheckableSettingItem
                      stateItem={this.props.devMode}
                      changeStateItem={this.props.changeDevMode}
                      stateItemName={SET_VIS_CFG_LABEL}
                      cydata={SET_VIS_CFG_CHECKBOX}
                    />
                  </div>
                </Elevation>

                <Elevation
                  className="component-style"
                  z={2}
                  style={styles.elevation}
                >
                  <div style={styles.headline}>
                    <span className="mdc-typography--headline6">
                      Local Storage:{' '}
                    </span>
                  </div>
                  <div style={styles.content}>
                    <Button
                      raised
                      style={styles.button}
                      onClick={() => this.showDialog(IDC_DIALOG_CLEAR_UI_CACHE)}
                      cydata={SET_RESET_SETTINGS_BUTTON}>
                      CLEAR UI CACHE
                    </Button>
                  </div>
                </Elevation>

                <Elevation
                  className="component-style"
                  z={2}
                  style={styles.elevation}
                >
                  <div style={styles.headline}>
                    <span className="mdc-typography--headline6">
                      About:{' '}
                    </span>
                  </div>
                  <div style={styles.content}>
                    <Grid>
                      <GridCell span={2}>UI Version:</GridCell>
                      <GridCell span={10}>{__VERSION__}</GridCell>
                    </Grid>
                  </div>
                </Elevation>

              </GridCell>
            </GridInner>
          </Grid>
        </div>
      </div>
    );
  }
}

const CheckableSettingItem = ({
  stateItem,
  changeStateItem,
  stateItemName,
  cydata
}) => {
  return (
    <Grid span={12}>
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
  elevation: {
    padding: 15,
    marginBottom: 10
  },
  headline: {
    marginBottom: 10
  },
  content: {
    marginBottom: 10
  },
  button: {
    color: 'white'
  }
};

const mapStateToProps = state => {
  return {
    automaticRefresh: state.ui.automaticRefresh,
    refreshInterval: state.ui.refreshInterval,
    devMode: state.ui.devMode,
    page: state.ui.page,
    currentDialog: state.ui.currentDialog
  };
};

const mapDispatchToProps = dispatch => {
  return {
    setAutomaticRefresh: val => dispatch(uiSetAutomaticRefresh(val)),
    changeRefreshInterval: val => dispatch(uiChangeRefreshInterval(val)),
    changeDevMode: mode => dispatch(uiChangeDevMode(mode)),
    showDialog: type => dispatch(uiChangeCurrentDialog(type)),
    meepSetDefaultState: () => dispatch(meepSetDefaultState())
  };
};

const ConnectedSettingsPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(SettingsPageContainer);

export default ConnectedSettingsPageContainer;
