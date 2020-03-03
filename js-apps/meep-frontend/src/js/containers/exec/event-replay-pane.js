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

import { uiExecChangeReplayFileSelected } from '../../state/ui';

import {
  EXEC_EVT_REPLAY_FILES,
  EXEC_BTN_REPLAY_START,
  EXEC_BTN_REPLAY_STOP
} from '../../meep-constants';

import {
  execChangeReplayFilesList
} from '../../state/exec';

import { PAGE_EXECUTE } from '../../meep-constants';

const ReplayFileSelect = props => {
  return (
    <Grid style={styles.field}>
      <GridCell span={12}>
        <Select
          style={styles.select}
          label="Replay file"
          fullwidth="true"
          outlined
          options={props.replayFiles}
          onChange={props.onChange}
          onClick={props.onClick}
          value={props.replayFileSelected}
          data-cy={EXEC_EVT_REPLAY_FILES}
        />
      </GridCell>
    </Grid>
  );
};

class EventReplayPane extends Component {
  constructor(props) {
    super(props);
    this.state = {
      description: null
    };
  }

  triggerReplay(play) {
    if (play) {
      if (this.props.replayLoop) {
        this.props.api.loopReplay(this.props.replayFileSelected, (error) => {
          if (error) {
            // TODO consider showing an alert
            // console.log(error);
          }
        });
      } else {
        this.props.api.playReplayFile(this.props.replayFileSelected, (error) => {
          if (error) {
            // TODO consider showing an alert
            // console.log(error);
          }
        });
      }
    } else { //stop
      this.props.api.stopReplayFile(this.props.replayFileSelected, (error) => {
        if (error) {
          // TODO consider showing an alert
          // console.log(error);
        }
      });
    }
  }

  changeLoop(checked) {
    this.props.onReplayLoopChanged(checked);
  }

  /**
   * Callback function to receive the result of the getReplayList operation.
   * @callback module:api/EventReplayApi~getReplayFileListCallback
   * @param {String} error Error message, if any.
   * @param {module:model/ReplayFileList} data The data returned by the service call.
   */
  getReplayFileListCb(error, data) {
    if (error !== null) {
      // TODO: consider showing an alert/toast
      return;
    }
    this.props.changeReplayFilesList(data.replayFiles);

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
    this.state.description = data.description;

  }

  updateReplayFileList() {
    this.props.api.getReplayFileList((error, data, response) => {
      this.getReplayFileListCb(error, data, response);
    });
  }

  getDescription(name) {
    this.props.api.getReplayFile(name, (error, data, response) => {
      this.getReplayFileCb(error, data, response);
    });
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
        <ReplayFileSelect
          replayFiles={this.props.replayFiles}
          replayFileSelected={this.props.replayFileSelected}
          onClick={() => this.updateReplayFileList()}
          onChange={event => {
            this.props.changeReplayFileSelected(event.target.value);
            this.getDescription(event.target.value);
          }}
        />
        <div style={styles.block}>
          <Typography use="subtitle2">{this.state.description}</Typography>
        </div>
        <Grid style={{ marginBottom: 10 }}>
          <GridCell span={2}>
            <Checkbox
              checked={this.props.replayLoop}
              onChange={e => this.changeLoop(e.target.checked)}
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
                onClick={() => this.triggerReplay(true)}
                data-cy={EXEC_BTN_REPLAY_START}
              >
                START
              </Button>
              <Button
                outlined
                style={styles.button}
                onClick={() => this.triggerReplay(false)}
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
    page: state.ui.page
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeReplayFileSelected: name => dispatch(uiExecChangeReplayFileSelected(name)),
    changeReplayFilesList: list => dispatch(execChangeReplayFilesList(list))
  };
};

const ConnectedEventReplayPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(EventReplayPane);

export default ConnectedEventReplayPane;
