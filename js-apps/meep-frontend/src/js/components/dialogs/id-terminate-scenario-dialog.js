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

import React, { Component }  from 'react';

import IDDialog from './id-dialog';
import { MEEP_DLG_TERMINATE_SCENARIO } from '../../meep-constants';

class IDTerminateScenarioDialog extends Component {

  constructor(props) {
    super(props);
    this.state={};
  }

  render() {
    return (
      <IDDialog
        title='Terminate Scenario'
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={() => this.props.onSubmit()}
        cydata={MEEP_DLG_TERMINATE_SCENARIO}
      >
        <span style={styles.text}>{'Are you sure you want to terminate the deployed scenario?'}</span>
      </IDDialog>
    );
  }
}

const styles = {
  text: {
    color: 'gray'
  }
};

export default IDTerminateScenarioDialog;
