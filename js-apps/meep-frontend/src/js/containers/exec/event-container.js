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
import { Typography } from '@rmwc/typography';

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

const StatusTable = props => {
  var timeToNextEvent = props.timeToNextEvent ? props.timeToNextEvent / 1000 : 0;
  var timeRemaining = props.timeRemaining ? props.timeRemaining / 1000 : 0;

  return (
    <Grid style={{border: '1px solid #e4e4e4'}}>
      <GridCell span={2}>
        <Typography use="headline6" style={{ marginLeft: 10 }}>Status:</Typography>
      </GridCell>
      <GridCell align={'middle'} span={3}>
        <Typography use="subtitle2" style={{ marginRight: 10 }}>REPLAY FILE:</Typography>
        <Typography use="body2">{props.name}</Typography>
      </GridCell>
      <GridCell align={'middle'} span={3}>
        <Typography use="subtitle2" style={{ marginRight: 10 }}>EVENT COUNT:</Typography>
        <Typography use="body2">{props.index} / {props.maxIndex}</Typography>
      </GridCell>
      <GridCell align={'middle'} span={4}>
        <Typography use="subtitle2" style={{ marginRight: 10 }}>NEXT/LAST EVENT (S):</Typography>
        <Typography use="body2">{timeToNextEvent.toFixed(2)} / {timeRemaining.toFixed(2)}</Typography>
      </GridCell>
    </Grid>
  );
};

class EventContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      replayFileName: null,
      eventIndex: null,
      maxIndex: null,
      loopMode: null,
      timeToNextEvent: null,
      timeRemaining: null
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
    this.props.changeEventCreationMode(false);
    this.props.changeEventReplayMode(true);

    // Refresh 
    this.props.onShowReplay();
  }

  render() {
    if (!this.props.eventCfgMode) {
      return null;
    }

    const replayStatus = this.props.replayStatus;

    return (
      <>
        <Elevation
          z={2}
          className="component-style"
          style={{ padding: 10, marginBottom: 10 }}
        >
          <Grid style={{ marginBottom: 10 }}>
            <GridCell span={2}>
              <div style={{ marginBottom: 10 }}>
                <span className="mdc-typography--headline6">
                  Event
                </span>
              </div>
            </GridCell>
            <GridCell span={10}>
              <div align={'right'}>
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
                  SAVE EVENTS
                </Button>
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

          <Grid>
            <GridCell span={12}>
              {replayStatus &&
                <StatusTable
                  name={replayStatus.replayFileRunning}
                  index={replayStatus.index}
                  maxIndex={replayStatus.maxIndex}
                  loopMode={replayStatus.loopMode}
                  timeRemaining={replayStatus.timeRemaining}
                  timeToNextEvent={replayStatus.timeToNextEvent}
                />
              }
            </GridCell>
          </Grid>
        </Elevation>
      </>
    );
  }
}

const mapStateToProps = state => {
  return {
    eventCreationMode: state.exec.eventCreationMode,
    eventReplayMode: state.exec.eventReplayMode,
    replayStatus: state.exec.state.replayStatus
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
