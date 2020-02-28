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

const greyColor = 'grey';

const styles = {
  button: {
    marginRight: 0
  },
  slider: {
    container: {
      marginTop: 10,
      marginBottom: 10,
      color: greyColor
    },
    boundaryValues: {
      marginTop: 15
    },
    title: {
      marginBottom: 0
    }
  },
  section1: {
    color: 'white',
    marginRight: 5
  },
  section2: {
    color: 'white',
    marginRight: 5,
    marginLeft: 10
  }
};

const ConfigurationView = props => {
  return (
    <>
      <Grid style={{ marginBottom: 10 }}>
        <GridCell span={6}>
          <Button
            raised
            style={styles.section1}
            onClick={props.onCreateEvent}
            data-cy={EXEC_BTN_MANUAL_REPLAY}
          >
            MANUAL
          </Button>
          <Button
            raised
            style={styles.section1}
            onClick={props.onReplayEvent}
            data-cy={EXEC_BTN_AUTO_REPLAY}
          >
            AUTO-REPLAY
          </Button>
          <Button
            raised
            style={styles.section1}
            onClick={props.onSaveReplay}
            data-cy={EXEC_BTN_SAVE_REPLAY}
          >
            SAVE EVENTS AS ...
          </Button>
        </GridCell>
      </Grid>
    </>
  );
};

const EventConfiguration = props => {
  if (!props.eventCfgMode) {
    return null;
  }

  let configurationView = null;

  configurationView = (
    <ConfigurationView
      onCreateEvent={props.onCreateEvent}
      onReplayEvent={props.onReplayEvent}
      onSaveReplay={props.onSaveReplay}
    />
  );
  return (
    <Elevation
      z={2}
      className="component-style"
      style={{ padding: 10, marginBottom: 10 }}
    >
      <Grid>
        <GridCell span={11}>
          <div style={{ marginBottom: 10 }}>
            <span className="mdc-typography--headline6">
              Event
            </span>
          </div>
        </GridCell>
        <GridCell span={1}>
          <Button
            outlined
            style={styles.button}
            onClick={() => props.onCloseEventCfg()}
          >
            Close
          </Button>
        </GridCell>
      </Grid>
      {configurationView}
    </Elevation>
  );
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

  changeReplayLoop(checked) {
    this.props.onReplayLoopChanged(checked);
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

    return (
      <>
        <EventConfiguration
          eventCfgMode={this.props.eventCfgMode}
          onCloseEventCfg={this.props.onCloseEventCfg}
          onCreateEvent={() => this.onCreateEvent()}
          onReplayEvent={() => this.onReplayEvent()}
          onSaveReplay={this.props.onSaveReplay}
          changeReplayLoop={checked => this.changeReplayLoop(checked)}
          replayLoop={this.props.replayLoop}
        />
      </>
    );
  }
}

const mapStateToProps = state => {
  return {
    eventCreationMode: state.exec.eventCreationMode,
    eventReplayMode: state.exec.eventReplayMode
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
