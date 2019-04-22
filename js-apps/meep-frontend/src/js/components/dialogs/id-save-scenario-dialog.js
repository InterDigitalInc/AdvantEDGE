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
import { MEEP_DLG_SAVE_SCENARIO } from '../../meep-constants';

class IDSaveScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={
      err: null,
      filename: null
    };
  }

  saveScenario() {
    this.props.saveScenario(this.scenarioName());
  }

  scenarioName() {
    return this.state.filename === null ?  this.props.scenarioName : this.state.filename;
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => this.saveScenario()}
        okDisabled={(!this.state.filename && this.props.scenarioNameRequired)|| this.state.err}
        cydata={MEEP_DLG_SAVE_SCENARIO}
      >
        <span style={styles.text}>{'Store the scenario in the MEEP Controller (overwrites any existing scenario with the same name)'}</span>

        <TextField outlined style={{width: '100%'}}
          label={'Scenario Name'}
          invalid={this.state.err || (!this.state.filename && this.props.scenarioNameRequired)}
          onChange={
            (e) => {
              const val = e.target.value;
              const err = (!val && val !=null)
                ? 'Please enter a filename'
                : '';
              this.setState({
                filename: val,
                err: err
              });
            }
          }
          value={this.scenarioName()}
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

const styles = {
  text: {
    color: 'gray'
  }
};

export default IDSaveScenarioDialog;
