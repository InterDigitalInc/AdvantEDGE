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

import { connect } from 'react-redux';
import React, { Component } from 'react';
import { Button } from '@rmwc/button';

import {
  EXEC_STATE_DEPLOYED,
  EXEC_BTN_SAVE_SCENARIO,
  EXEC_BTN_DEPLOY,
  EXEC_BTN_TERMINATE,
  EXEC_BTN_REFRESH,
  EXEC_BTN_EVENT,
  EXEC_BTN_CONFIG
} from '../../meep-constants';

import {
  scenarioPodsPending,
  scenarioPodsTerminating,
  scenarioPodsTerminated
} from '../../state/exec';

class ExecPageScenarioButtons extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  canDeploy() {
    return (
      this.props.podsTerminated &&
      this.props.scenarioState.scenario !== EXEC_STATE_DEPLOYED
    );
  }

  canTerminate() {
    return (
      !this.props.podsTerminating &&
      this.props.scenarioState.scenario === EXEC_STATE_DEPLOYED &&
      this.props.okToTerminate
    );
  }

  canSaveScenario() {
    return (
      !this.props.podsPending &&
      !this.props.podsTerminating &&
      !this.props.podsTerminated
    );
  }

  canRefresh() {
    return (
      !this.props.podsPending &&
      !this.props.podsTerminating &&
      !this.props.podsTerminated
    );
  }

  canOpenDashCfg() {
    return (
      !this.props.podsPending &&
      !this.props.podsTerminating &&
      !this.props.podsTerminated
    );
  }

  canOpenEventCfg() {
    return (
      !this.props.podsPending &&
      !this.props.podsTerminating &&
      !this.props.podsTerminated
    );
  }

  render() {
    return (
      <div style={{ marginTop: 10 }}>
        <Button
          raised
          style={styles.section1}
          onClick={this.props.onDeploy}
          disabled={!this.canDeploy()}
          data-cy={EXEC_BTN_DEPLOY}
        >
          DEPLOY
        </Button>
        <Button
          raised
          style={styles.section1}
          onClick={() => this.props.onSaveScenario()}
          disabled={!this.canSaveScenario()}
          data-cy={EXEC_BTN_SAVE_SCENARIO}
        >
          SAVE
        </Button>
        <Button
          raised
          style={styles.section1}
          onClick={this.props.onTerminate}
          disabled={!this.canTerminate()}
          data-cy={EXEC_BTN_TERMINATE}
        >
          TERMINATE
        </Button>
        <Button
          raised
          style={styles.section1}
          onClick={this.props.onRefresh}
          disabled={!this.canRefresh()}
          data-cy={EXEC_BTN_REFRESH}
        >
          REFRESH
        </Button>
        <Button
          raised
          style={styles.section2}
          onClick={this.props.onOpenEventCfg}
          disabled={!this.canOpenEventCfg()}
          data-cy={EXEC_BTN_EVENT}
        >
          EVENT
        </Button>
        <Button
          raised
          style={styles.section2}
          onClick={this.props.onOpenDashCfg}
          disabled={!this.canOpenDashCfg()}
          data-cy={EXEC_BTN_CONFIG}
        >
          DASHBOARD
        </Button>
      </div>
    );
  }
}

const styles = {
  section1: {
    color: 'white',
    marginRight: 5
  },
  section2: {
    color: 'white',
    marginRight: 5,
    marginLeft: 10
  }
};

const mapStateToProps = state => {
  return {
    podsTerminated: scenarioPodsTerminated(state),
    podsTerminating: scenarioPodsTerminating(state),
    podsPending: scenarioPodsPending(state),
    eventCreationMode: state.ui.eventCreationMode,
    scenarioState: state.exec.state,
    okToTerminate: state.exec.state.okToTerminate
  };
};

const ConnectedExecPageScenarioButtons = connect(
  mapStateToProps,
  null
)(ExecPageScenarioButtons);

export default ConnectedExecPageScenarioButtons;
