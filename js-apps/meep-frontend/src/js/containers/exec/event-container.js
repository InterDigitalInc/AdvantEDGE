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
import { Grid, GridCell } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
import { Button } from '@rmwc/button';

import {
  uiExecChangeEventCreationMode,
  uiExecChangeEventReplayMode
} from '../../state/ui';

import {
  EXEC_BTN_MANUAL_REPLAY,
  EXEC_BTN_AUTO_REPLAY,
  EXEC_BTN_SAVE_REPLAY
} from '../../meep-constants';

const styles = {
  button: {
    marginRight: 10
  }
};

class EventContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sourceNodeId: ''
    };
  }

  componentDidMount() { }

  componentWillUnmount() {
    clearInterval(this.dataTimer);
  }

  // CREATE EVENT PANE
  onCreateEvent() {
    this.props.changeEventCreationMode(true);
    this.props.changeEventReplayMode(false);
  }

  // SHOW REPLAY EVENT PANE
  onReplayEvent() {
    this.props.changeEventReplayMode(true);
    this.props.changeEventCreationMode(false);
  }

  render() {

    if (!this.props.eventCfgMode) {
      return null;
    }

    return (
      <>
        <Elevation
          z={2}
          className="component-style"
          style={{ padding: 10, marginBottom: 10 }}
        >
          <Grid>
            <GridCell span={6}>
              <div style={{ marginBottom: 10 }}>
                <span className="mdc-typography--headline6">
                  Event
                </span>
              </div>
            </GridCell>
            <GridCell span={6}>
              <div align={'right'}>
                <Button
                  outlined
                  style={styles.button}
                  onClick={this.props.onCloseEventCfg}
                >
                  Close
                </Button>
              </div>
            </GridCell>
          </Grid>

          <Grid style={{ marginBottom: 10 }}>
            <GridCell span={6}>
              <Button
                outlined
                style={styles.button}
                onClick={() => this.onCreateEvent()}
                data-cy={EXEC_BTN_MANUAL_REPLAY}
              >
                MANUAL
              </Button>
              <Button
                outlined
                style={styles.button}
                onClick={() => this.onReplayEvent()}
                data-cy={EXEC_BTN_AUTO_REPLAY}
              >
                AUTO-REPLAY
              </Button>
              <Button
                outlined
                style={styles.button}
                onClick={this.props.onSaveReplay}
                data-cy={EXEC_BTN_SAVE_REPLAY}
              >
                SAVE EVENTS AS ...
              </Button>
            </GridCell>
          </Grid>

          <div>{this.props.replayStatus ? this.props.replayStatus.status : 'NONE'}</div>
          
        </Elevation>
      </>
    );
  }
}

const mapStateToProps = state => {
  return {
    eventCreationMode: state.exec.eventCreationMode,
    eventReplayMode: state.exec.eventReplayMode,
    replayStatus: state.ui.replayStatus
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeEventCreationMode: mode => dispatch(uiExecChangeEventCreationMode(mode)),
    changeEventReplayMode: mode => dispatch(uiExecChangeEventReplayMode(mode))
  };
};

const ConnectedEventContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(EventContainer);

export default ConnectedEventContainer;
