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
import * as YAML from 'yamljs';

import IDDialog from './id-dialog';

class IDExportScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={
      filename:null
    };
    this.exportScenarioTextFile= '';
  }

  makeTextFile(text) {
    var data = new Blob([text], {type: 'text/plain'});
    // If we are replacing a previously generated file we need to
    // manually revoke the object URL to avoid memory leaks.
    if (this.state.exportScenarioTextFile !== null) {
      window.URL.revokeObjectURL(this.state.exportScenarioTextFile);
    }

    this.exportScenarioTextFile = window.URL.createObjectURL(data);

    return this.exportScenarioTextFile;
  }

  exportScenario() {

    if (this.state.filename == '') {
      // console.log('Invalid file name provided');
      // TODO: consider showing an alert
      return;
    }

    var filename = (this.state.filename === null) ? this.props.scenarioName + '.yaml' : this.state.filename;
    var link = document.getElementById('export-scenario-link');
    link.href = this.makeTextFile(YAML.stringify(this.props.scenario, 20, 4));
    // link.href = makeTextFile(JSON.stringify(meep.cfg.scenario, null, "\t"));
    link.download = filename;
    link.click();
  }

  render() {

    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => {
          this.exportScenario();
          this.setState({filename: null, err: null});
        }}
        okDisabled={this.state.err}
      >
        <TextField outlined style={{width: '100%'}}
          label={'Export filename'}
          invalid={this.state.err}
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
          value={this.state.filename === null ?  this.props.scenarioName + '.yaml' : this.state.filename}
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

export default IDExportScenarioDialog;
