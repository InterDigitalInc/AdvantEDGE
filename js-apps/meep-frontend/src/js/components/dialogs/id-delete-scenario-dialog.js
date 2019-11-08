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
import { MEEP_DLG_DEL_SCENARIO } from '../../meep-constants';

class IDDeleteScenarioDialog extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
        onSubmit={this.props.deleteScenario}
        cydata={MEEP_DLG_DEL_SCENARIO}
      >
        <span style={styles.text}>
          {
            'Are you sure you want to delete the current scenario from the MEEP Controller?'
          }
        </span>
      </IDDialog>
    );
  }
}

const styles = {
  text: {
    color: 'gray'
  }
};

export default IDDeleteScenarioDialog;
