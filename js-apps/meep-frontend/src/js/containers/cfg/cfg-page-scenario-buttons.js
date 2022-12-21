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

import { connect } from 'react-redux';
import React, { Component } from 'react';
import { Button } from '@rmwc/button';
import { TextField } from '@rmwc/textfield';

import {
  CFG_STATE_IDLE,
  CFG_STATE_NEW,
  CFG_STATE_LOADED,
  CFG_BTN_NEW_SCENARIO,
  CFG_BTN_OPEN_SCENARIO,
  CFG_BTN_SAVE_SCENARIO,
  CFG_BTN_DEL_SCENARIO,
  CFG_BTN_IMP_SCENARIO,
  CFG_BTN_EXP_SCENARIO
} from '../../meep-constants';

class CfgPageScenarioButtons extends Component {
  constructor(props) {
    super(props);
  }

  canCreateNewScenario() {
    const cfgState = this.props.cfgState;
    return (
      this.props.cfgState === CFG_STATE_IDLE ||
      cfgState === CFG_STATE_NEW ||
      cfgState === CFG_STATE_LOADED
    );
  }

  canOpenScenario() {
    const cfgState = this.props.cfgState;
    return (
      cfgState === CFG_STATE_IDLE ||
      cfgState === CFG_STATE_NEW ||
      cfgState === CFG_STATE_LOADED
    );
  }

  canSaveScenario() {
    const cfgState = this.props.cfgState;
    return cfgState === CFG_STATE_NEW || cfgState === CFG_STATE_LOADED;
  }

  canDeleteScenario() {
    const cfgState = this.props.cfgState;
    return cfgState === CFG_STATE_LOADED;
  }

  canImportScenario() {
    const cfgState = this.props.cfgState;
    return (
      cfgState === CFG_STATE_IDLE ||
      cfgState === CFG_STATE_NEW ||
      cfgState === CFG_STATE_LOADED
    );
  }

  canExportScenario() {
    const cfgState = this.props.cfgState;
    return cfgState === CFG_STATE_NEW || cfgState === CFG_STATE_LOADED;
  }

  render() {
    const input = (
      <TextField
        type="file"
        hidden
        ref={input => (this.inputElement = input)}
        onClick={() => this.props.onInputScenario(this.inputElement.input_)}
        onChange={() => this.props.onInputScenario(this.inputElement.input_)}
        style={{ height: '0%', width: '0%', marginTop: -20, paddingTop: -20 }}
      />
    );
    return (
      <>
        <Button
          raised
          style={styles.button}
          onClick={() => this.props.onNewScenario()}
          disabled={!this.canCreateNewScenario()}
          data-cy={CFG_BTN_NEW_SCENARIO}
        >
          NEW
        </Button>
        <Button
          raised
          style={styles.button}
          onClick={() => this.props.onOpenScenario()}
          disabled={!this.canOpenScenario()}
          data-cy={CFG_BTN_OPEN_SCENARIO}
        >
          OPEN
        </Button>
        <Button
          raised
          style={styles.button}
          onClick={() => this.props.onSaveScenario()}
          disabled={!this.canSaveScenario()}
          data-cy={CFG_BTN_SAVE_SCENARIO}
        >
          SAVE
        </Button>
        <Button
          raised
          style={styles.button}
          onClick={() => this.props.onDeleteScenario()}
          disabled={!this.canDeleteScenario()}
          data-cy={CFG_BTN_DEL_SCENARIO}
        >
          DELETE
        </Button>

        {input}

        <Button
          raised
          style={{ ...styles.button, marginLeft: 10 }}
          onClick={() => {
            this.inputElement.input_.click();
          }}
          disabled={!this.canImportScenario()}
          data-cy={CFG_BTN_IMP_SCENARIO}
        >
          IMPORT
        </Button>

        <Button
          raised
          style={styles.button}
          onClick={() => this.props.onExportScenario()}
          disabled={!this.canExportScenario()}
          data-cy={CFG_BTN_EXP_SCENARIO}
        >
          EXPORT
        </Button>

        <a id="export-scenario-link" download="config.yaml" hidden></a>
      </>
    );
  }
}

const styles = {
  button: {
    color: 'white',
    marginRight: 5
  },
  icon: {
    color: 'white'
  }
};

const mapStateToProps = state => {
  return {
    cfgTable: state.cfg.table,
    execVis: state.exec.vis,
    cfgVis: state.cfg.vis,
    devMode: state.ui.devMode,
    cfgState: state.cfg.state,
    scenarioName: state.cfg.scenario.name
  };
};

const ConnectedCfgPageScenarioButtons = connect(
  mapStateToProps,
  null
)(CfgPageScenarioButtons);

export default ConnectedCfgPageScenarioButtons;
