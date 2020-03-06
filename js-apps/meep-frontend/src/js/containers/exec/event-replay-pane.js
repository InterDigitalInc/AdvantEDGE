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
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import { Select } from '@rmwc/select';
import { Grid, GridInner, GridCell } from '@rmwc/grid';
import { Typography } from '@rmwc/typography';

import {
  uiExecChangeReplayFileSelected,
  uiExecChangeReplayFileDesc,
  uiExecChangeReplayLoop
} from '../../state/ui';

import {
  EXEC_EVT_REPLAY_FILES,
  EXEC_BTN_REPLAY_START,
  EXEC_BTN_REPLAY_STOP
} from '../../meep-constants';

import { PAGE_EXECUTE } from '../../meep-constants';

class EventReplayPane extends Component {
  constructor(props) {
    super(props);
    this.state = {
      description: null
    };
  }

  playReplay(name, loop) {
    if (name !== 'None') {
      if (loop) {
        this.props.api.loopReplay(name, (error) => {
          if (error) {
            // TODO consider showing an alert
            // console.log(error);
          }
        });
      } else {
        this.props.api.playReplayFile(name, (error) => {
          if (error) {
            // TODO consider showing an alert
            // console.log(error);
          }
        });
      }
    }
  }

  stopReplay(name) {
    this.props.api.stopReplayFile(name, (error) => {
      if (error) {
        // TODO consider showing an alert
        // console.log(error);
      }
    });
  }

  /**
   * Callback function to receive the result of the getReplayFile operation.
   * @callback module:api/EventReplayApi~getReplayFileCallback
   * @param {String} error Error message, if any.
   * @param {module:model/Replay} data The data returned by the service call.
   */
  getReplayFileCb(error, data) {
    if (error !== null) {
      // TODO: consider showing an alert/toast
      return;
    }
    this.props.changeReplayFileDesc(data.description);
  }

  getDescription(name) {
    this.props.api.getReplayFile(name, (error, data, response) => {
      this.getReplayFileCb(error, data, response);
    });
  }

  replayRunning() {
    return this.props.replayStatus ? true : false;
  }

  canPlay() {
    return !this.replayRunning() &&
      this.props.replayFileSelected &&
      this.props.replayFileSelected !== 'None';
  }

  render() {
    if (this.props.page !== PAGE_EXECUTE || this.props.hide) {
      return null;
    }

    return (
      <div style={{ padding: 10 }}>
        <div style={styles.block}>
          <Typography use="headline6">Replay Events</Typography>
        </div>
        <Grid style={styles.field}>
          <GridCell span={12}>
            <Select
              style={styles.select}
              label="Replay file"
              fullwidth="true"
              outlined
              options={this.props.replayFiles}
              onChange={event => {
                this.props.changeReplayFileSelected(event.target.value);
                this.getDescription(event.target.value);
              }}
              value={this.props.replayFileSelected}
              disabled={this.replayRunning()}
              data-cy={EXEC_EVT_REPLAY_FILES}
            />
          </GridCell>
          <GridCell span={12}>
            <Typography use="subtitle2">{this.props.replayFileDesc}</Typography>
          </GridCell>
        </Grid>

        <Grid style={{ marginBottom: 10 }}>
          <GridCell span={2}>
            <Checkbox
              checked={this.props.replayLoop}
              onChange={e => this.props.changeReplayLoop(e.target.checked)}
              disabled={this.replayRunning()}
            >
              Loop
            </Checkbox>
          </GridCell>
        </Grid>
        <Grid style={{ marginTop: 10 }}>
          <GridInner>
            <GridCell span={12}>
              <Button
                outlined
                style={styles.button}
                onClick={() => this.playReplay(this.props.replayFileSelected, this.props.replayLoop)}
                disabled={!this.canPlay()}
                data-cy={EXEC_BTN_REPLAY_START}
              >
                PLAY
              </Button>
              <Button
                outlined
                style={styles.button}
                onClick={() => this.stopReplay(this.props.replayFileSelected)}
                disabled={!this.replayRunning()}
                data-cy={EXEC_BTN_REPLAY_STOP}
              >
                STOP
              </Button>
              <Button
                outlined
                style={styles.button}
                onClick={this.props.onClose}
              >
                Close
              </Button>
            </GridCell>
          </GridInner>
        </Grid>
      </div>
    );
  }
}

const styles = {
  button: {
    marginRight: 10
  },
  block: {
    marginBottom: 20
  },
  field: {
    marginBottom: 10
  },
  select: {
    width: '100%'
  }
};

const mapStateToProps = state => {
  return {
    page: state.ui.page,
    replayStatus: state.exec.state.replayStatus,
    replayFiles: state.ui.replayFiles,
    replayFileSelected: state.ui.replayFileSelected,
    replayFileDesc: state.ui.replayFileDesc,
    replayLoop: state.ui.eventReplayLoop
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeReplayFileSelected: name => dispatch(uiExecChangeReplayFileSelected(name)),
    changeReplayFileDesc: name => dispatch(uiExecChangeReplayFileDesc(name)),
    changeReplayLoop: val => dispatch(uiExecChangeReplayLoop(val))
  };
};

const ConnectedEventReplayPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(EventReplayPane);

export default ConnectedEventReplayPane;
