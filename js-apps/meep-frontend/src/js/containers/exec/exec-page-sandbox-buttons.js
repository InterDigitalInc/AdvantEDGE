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

import {
  EXEC_BTN_NEW_SANDBOX,
  EXEC_BTN_DELETE_SANDBOX
} from '../../meep-constants';

class ExecPageSandboxButtons extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  canDelete() {
    return this.props.sandbox;
  }

  render() {
    return (
      <div>
        <Button
          raised
          style={styles.buttonWithMargin}
          onClick={this.props.onNewSandbox}
          data-cy={EXEC_BTN_NEW_SANDBOX}
        >
          NEW
        </Button>
        <Button
          raised
          style={styles.button}
          onClick={this.props.onDeleteSandbox}
          disabled={!this.canDelete()}
          data-cy={EXEC_BTN_DELETE_SANDBOX}
        >
          DELETE
        </Button>
      </div>
    );
  }
}

const styles = {
  button: {
    color: 'white',
    marginRight: 5
  },
  buttonWithMargin: {
    color: 'white',
    marginRight: 5,
    marginLeft: 10
  }
};

const mapStateToProps = state => {
  return {
    scenarioState: state.exec.state
  };
};

const ConnectedExecPageSandboxButtons = connect(
  mapStateToProps,
  null
)(ExecPageSandboxButtons);

export default ConnectedExecPageSandboxButtons;
