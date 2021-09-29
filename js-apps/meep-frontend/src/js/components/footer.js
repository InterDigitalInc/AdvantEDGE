/*
 * Copyright (c) 2021  InterDigital Communications, Inc
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
import { Grid, GridCell } from '@rmwc/grid';

import {
  MEEP_DOCS_URL,
  MEEP_GITHUB_URL,
  MEEP_DISCUSSIONS_URL,
  MEEP_LICENSE_URL,
  MEEP_CONTRIBUTING_URL
} from '../meep-constants';

const Footer = () => {

  return (
    <Grid style={styles.footer}>
      <GridCell span={2}/>
      <GridCell span={3} align={'middle'} style={styles.copyright}>
        <Typography use="body1">© 2021 InterDigital, Inc</Typography>
      </GridCell>
      <GridCell span={5} align={'middle'}>
        <div align={'right'}>
          <a href={MEEP_DOCS_URL} target="_blank" style={styles.link}>
            <Typography use="body1">Documentation</Typography>
          </a>
          <span style={styles.separator}>|</span>
          <a href={MEEP_GITHUB_URL} target="_blank" style={styles.link}>
            <Typography use="body1">GitHub</Typography>
          </a>
          <span style={styles.separator}>|</span>
          <a href={MEEP_DISCUSSIONS_URL} target="_blank" style={styles.link}>
            <Typography use="body1">Discussions</Typography>
          </a>
          <span style={styles.separator}>|</span>
          <a href={MEEP_LICENSE_URL} target="_blank" style={styles.link}>
            <Typography use="body1">License</Typography>
          </a>
          <span style={styles.separator}>|</span>
          <a href={MEEP_CONTRIBUTING_URL} target="_blank" style={styles.link}>
            <Typography use="body1">Contributing</Typography>
          </a>
        </div>
      </GridCell>
    </Grid>
  );
};

const styles = {
  footer: {
    backgroundColor: '#379DD8',
    padding: 5
  },
  copyright: {
    color: 'white',
    padding: 10
  },
  link: {
    color: 'white',
    textDecoration: 'none',
    marginLeft: 5,
    marginRight: 10
  },
  separator: {
    color: '#FFFFFF',
    margin: 3
  }
};

export default Footer;
