/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import React, { Component }  from 'react';
import { TextField, TextFieldHelperText } from '@rmwc/textfield';

import IDDialog from './id-dialog';
import { MEEP_DLG_NEW_SCENARIO, MEEP_DLG_NEW_SCENARIO_NAME } from '../../meep-constants';

class IDNewScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={
      scenarioName: '',
      err: null
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
    }

    this.setState({scenarioName: name, err: err});
  }

  /**
     * Callback function to receive the result of the getScenario operation.
     * @callback module:api/ScenarioConfigurationApi~getScenarioCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Scenario} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */
  getScenarioNewCb(error/*, data, response*/) {

    if (error == null) {
      // TODO: consider showing an alert
      return;
    }

    // Validate scenario name
    if (this.state.scenarioName == '' || this.state.err != null) {
      // TODO: consider showing an alert
      return;
    }

    // Clear scenario
    this.props.createScenario(this.state.scenarioName);
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => {
          this.props.api.getScenario(this.state.scenarioName, (error, data, response) => {
            this.getScenarioNewCb(error, data, response);
          });
        }}
        cydata={MEEP_DLG_NEW_SCENARIO}
      >
        <TextField outlined style={{width: '100%'}}
          label={'Scenario Name'}
          onChange={(e) => {this.changeScenarioName(e.target.value);}}
          value={this.state.scenarioName}
          invalid={(this.state.err) ? true : false}
          data-cy={MEEP_DLG_NEW_SCENARIO_NAME}
        />
        <TextFieldHelperText validationMsg={true}>
          <span>
            {this.state.err}
          </span>
        </TextFieldHelperText>
      </IDDialog>
    );
  }
}

export default IDNewScenarioDialog;
