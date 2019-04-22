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
import { MEEP_DLG_OPEN_SCENARIO, MEEP_DLG_OPEN_SCENARIO_SELECT } from '../../meep-constants';

class IDOpenScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={
      selectedScenario: null
    };
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => {
          this.props.api.getScenario(this.state.selectedScenario, this.props.getScenarioLoadCb);
        }}
        cydata={MEEP_DLG_OPEN_SCENARIO}
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
          cydata={MEEP_DLG_OPEN_SCENARIO_SELECT}
        />
      </IDDialog>
    );
  }
}

export default IDOpenScenarioDialog;
