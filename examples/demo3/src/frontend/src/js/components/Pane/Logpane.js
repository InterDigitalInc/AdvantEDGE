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

import React from 'react';
import { Typography } from '@rmwc/typography';
import { Elevation } from '@rmwc/elevation';
import Title from '@/js/components/Title';
import '@material/typography/dist/mdc.typography.css';
export default function LogPane({ data }) {
  // Generate data
  function renderActivity(data) {
    if (data) {
      return data.map((element) => {
        return (
          <Typography use="caption" style={{ display: 'block' }}>
            {element}
          </Typography>
        );
      });
    }
  }

  return (
    <div>
      <Elevation
        z={2}
        className="component-style "
        style={{ padding: 10, marginBottom: 10 }}
      >
        <Title>Activity Logs</Title>
        <div style={{ height: '45vh', overflowX: 'auto', overflowY: 'auto', whiteSpace: 'nowrap'}}>
          {renderActivity(data)}
        </div>
      </Elevation>
    </div>
  );
}
