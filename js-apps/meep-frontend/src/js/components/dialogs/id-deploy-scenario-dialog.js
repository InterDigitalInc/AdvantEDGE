/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import React, { Component }  from 'react';

import IDDialog from './id-dialog';
import IDSelect from '../helper-components/id-select';
import { MEEP_DLG_DEPLOY_SCENARIO, MEEP_DLG_DEPLOY_SCENARIO_SELECT } from '../../meep-constants';

class IDDeployScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={
      selectedScenario: null
    };
  }

  onDeployScenario() {
    if (this.state.selectedScenario === '') {
      // console.log('Invalid scenario name');
      // TODO: consider showing an alert
      return;
    }
    this.props.api.activateScenario(this.state.selectedScenario, this.props.activateScenarioCb);
  }

  render() {
    return (
      <IDDialog
        title='Deploy Scenario'
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => {this.onDeployScenario();}}
        cydata={MEEP_DLG_DEPLOY_SCENARIO}
      >
        <IDSelect
          label={this.props.label || 'Scenario'}
          value={this.props.value}
          options={this.props.options}
          onChange={
            (e) => {
              this.setState({selectedScenario: e.target.value});
            }
          }
          cydata={MEEP_DLG_DEPLOY_SCENARIO_SELECT}
        />
      </IDDialog>
    );
  }
}

export default IDDeployScenarioDialog;

//
