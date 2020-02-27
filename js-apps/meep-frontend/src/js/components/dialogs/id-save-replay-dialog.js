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

import React, { Component } from 'react';
import { TextField, TextFieldHelperText } from '@rmwc/textfield';
import IDDialog from './id-dialog';
import { MEEP_DLG_SAVE_REPLAY } from '../../meep-constants';

class IDSaveReplayDialog extends Component {
  constructor(props) {
    super(props);
    this.state = {
      err: null,
      replayName: null,
      description: null
    };
  }

  changeReplayName(name) {
    var err = null;

    if (name) {
      if (name.length > 20) {
        err = 'Maximum 20 characters';
      } else if (!name.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
        err = 'Lowercase alphanumeric or \'-\'';
      }
    } else {
      err = 'Please enter a replay file name';
    }
    this.setState({
      replayName: name,
      err: err
    });
  }

  changeDescription(desc) {
    var err = null;

    if (desc) {
      if (desc.length > 30) {
        err = 'Maximum 30 characters';
      }
    }
    this.setState({
      description: desc,
      err: err
    });
  }

  saveReplay() {
    this.props.saveReplay(this.state);
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => this.saveReplay()}
        okDisabled={
          (!this.state.replayName && this.props.replayNameRequired) ||
          this.state.err
        }
        cydata={MEEP_DLG_SAVE_REPLAY}
      >
        <span style={styles.text}>
          {
            'Store the events as a replay file in the MEEP Controller for current deployed scenario (overwrites any existing replay file with the same name)'
          }
        </span>

        <TextField
          outlined
          style={{ width: '100%' }}
          label={'Replay Name'}
          invalid={
            this.state.err ||
            (!this.state.replayName && this.props.replayNameRequired)
          }
          onChange={e => this.changeReplayName(e.target.value)}
          value={this.replayName}
        />
       <TextField
          outlined
          style={{ width: '100%' }}
          label={'Replay Description'}
          invalid={
            this.state.err 
          }
          onChange={e => this.changeDescription(e.target.value)}
          value={this.description}
        />

        <TextFieldHelperText validationMsg={true}>
          <span>{this.state.err}</span>
        </TextFieldHelperText>
      </IDDialog>
    );
  }
}

const styles = {
  text: {
    color: 'gray'
  }
};

export default IDSaveReplayDialog;
