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
import {
  MEEP_DLG_NEW_SANDBOX,
  MEEP_DLG_NEW_SANDBOX_NAME
} from '../../meep-constants';

class IDNewSandboxDialog extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sandboxName: '',
      err: null
    };
  }

  changeSandboxName(name) {
    var err = null;

    if (name) {
      if (name.length > 20) {
        err = 'Maximum 20 characters';
      } else if (!name.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
        err = 'Lowercase alphanumeric or \'-\'';
      }
    }

    this.setState({ sandboxName: name, err: err });
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => {this.props.createSandbox(this.state.sandboxName);}}
        cydata={MEEP_DLG_NEW_SANDBOX}
      >
        <TextField
          outlined
          style={{ width: '100%' }}
          label={'Sandbox Name'}
          onChange={e => {
            this.changeSandboxName(e.target.value);
          }}
          value={this.state.sandboxName}
          invalid={this.state.err ? true : false}
          data-cy={MEEP_DLG_NEW_SANDBOX_NAME}
        />
        <TextFieldHelperText validationMsg={true}>
          <span>{this.state.err}</span>
        </TextFieldHelperText>
      </IDDialog>
    );
  }
}

export default IDNewSandboxDialog;
