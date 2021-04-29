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
import { MEEP_DLG_SAVE_SCENARIO } from '../../meep-constants';

class IDSaveScenarioDialog extends Component {
  constructor(props) {
    super(props);
    this.state = {
      err: null,
      scenarioName: null
    };
  }

  changeScenarioName(name) {
    var err = null;

    if (name) {
      if (name.length > 20) {
        err = 'Maximum 20 characters';
      } else if (!name.match(/^(([a-z0-9][-a-z0-9.]*)?[a-z0-9])+$/)) {
        err = 'Lowercase alphanumeric or \'-\'';
      }
    } else {
      err = 'Please enter a scenario name';
    }

    this.setState({
      scenarioName: name,
      err: err
    });
  }

  saveScenario() {
    this.props.saveScenario(this.scenarioName());
  }

  scenarioName() {
    return this.state.scenarioName === null
      ? this.props.scenarioName
      : this.state.scenarioName;
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => this.saveScenario()}
        okDisabled={
          (!this.state.scenarioName && this.props.scenarioNameRequired) ||
          this.state.err
        }
        cydata={MEEP_DLG_SAVE_SCENARIO}
      >
        <span style={styles.text}>
            Store the scenario in the MEEP Controller<br/>
            NOTE: overwrites any existing scenario with the same name
        </span>

        <TextField
          style={{ width: '100%', marginTop: 20 }}
          label={'Scenario Name'}
          invalid={
            this.state.err ||
            (!this.state.scenarioName && this.props.scenarioNameRequired)
          }
          onChange={e => this.changeScenarioName(e.target.value)}
          value={this.scenarioName()}
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

export default IDSaveScenarioDialog;
