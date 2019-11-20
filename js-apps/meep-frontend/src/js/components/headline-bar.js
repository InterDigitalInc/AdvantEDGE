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

import React from 'react';

import { NO_SCENARIO_NAME, MEEP_LBL_SCENARIO_NAME } from '../meep-constants';

const HeadlineBar = ({ titleLabel, scenarioName }) => {
  var name = scenarioName === NO_SCENARIO_NAME ? 'None' : scenarioName;
  return (
    <div style={{ marginTop: 14, marginBottom: 14 }}>
      <span className="mdc-typography--headline6">{titleLabel}: </span>
      <span
        data-cy={MEEP_LBL_SCENARIO_NAME}
        className="idcc-margin-left mdc-typography--headline6 mdc-theme--primary"
      >
        {name}
      </span>
    </div>
  );
};

export default HeadlineBar;
