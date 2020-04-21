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

import IDDialog from './id-dialog';
import IDSelect from '../helper-components/id-select';
import {
  SANDBOX_NAME,
  MEEP_DLG_DEPLOY_SCENARIO,
  MEEP_DLG_DEPLOY_SCENARIO_SELECT
} from '../../meep-constants';

class IDDeployScenarioDialog extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedScenario: null
    };
  }

  onDeployScenario() {
    if (this.state.selectedScenario === '') {
      // console.log('Invalid scenario name');
      // TODO: consider showing an alert
      return;
    }

    // this.props.api.activateScenario(
    //   this.state.selectedScenario,
    //   null,
    //   this.props.activateScenarioCb
    // );

    var sandboxConfig = {
      scenarioName: this.state.selectedScenario
    };
    this.props.api.createSandboxWithName(
      SANDBOX_NAME,
      sandboxConfig,
      this.props.createSandboxWithNameCb
    );
  }

  render() {
    return (
      <IDDialog
        title="Deploy Scenario"
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => {
          this.onDeployScenario();
        }}
        cydata={MEEP_DLG_DEPLOY_SCENARIO}
      >
        <IDSelect
          label={this.props.label || 'Scenario'}
          value={this.props.value}
          options={this.props.options}
          onChange={e => {
            this.setState({ selectedScenario: e.target.value });
          }}
          cydata={MEEP_DLG_DEPLOY_SCENARIO_SELECT}
        />
      </IDDialog>
    );
  }
}

export default IDDeployScenarioDialog;

//
