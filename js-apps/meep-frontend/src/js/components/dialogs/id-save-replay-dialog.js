/*
 * Copyright (c) 2022  The AdvantEDGE Authors
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
import { 
  MEEP_DLG_SAVE_REPLAY,
  MEEP_DLG_SAVE_REPLAY_NAME,
  MEEP_DLG_SAVE_REPLAY_DESCRIPTION
} from '../../meep-constants';

class IDSaveReplayDialog extends Component {
  constructor(props) {
    super(props);
    this.state = {
      replayName: null,
      replayErr: null,
      description: null,
      descriptionErr: null
    };
  }

  changeReplayName(name) {
    var err = null;

    if (name) {
      if (name.length > 30) {
        err = 'Maximum 30 characters';
      } else if (!name.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
        err = 'Lowercase alphanumeric or \'-\'';
      }
    } else {
      err = 'Please enter a replay file name';
    }
    this.setState({
      replayName: name,
      replayErr: err
    });
  }

  changeDescription(desc) {
    var err = null;

    if (desc) {
      if (desc.length > 100) {
        err = 'Maximum 100 characters';
      }
    }
    this.setState({
      description: desc,
      descriptionErr: err
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
          this.state.replayErr || this.state.descriptionErr
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
            this.state.replayErr
          }
          onChange={e => this.changeReplayName(e.target.value)}
          value={this.replayName}
          data-cy={MEEP_DLG_SAVE_REPLAY_NAME}
        />
        <TextFieldHelperText validationMsg={true}>
          <span>{this.state.replayErr}</span>
        </TextFieldHelperText>

        <TextField
          outlined
          style={{ width: '100%' }}
          label={'Replay Description'}
          invalid={
            this.state.descriptionErr 
          }
          onChange={e => this.changeDescription(e.target.value)}
          value={this.description}
          data-cy={MEEP_DLG_SAVE_REPLAY_DESCRIPTION}
        />
        <TextFieldHelperText validationMsg={true}>
          <span>{this.state.descriptionErr}</span>
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
