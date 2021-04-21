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

import React, { Component } from 'react';
import IDDialog from './id-dialog';
import { Grid, GridCell } from '@rmwc/grid';
import { Button } from '@rmwc/button';
import GitHubIcon from '@/img/logo-github.svg';
import GitLabIcon from '@/img/logo-gitlab.svg';

import {
  OAUTH_PROVIDER_GITHUB,
  OAUTH_PROVIDER_GITLAB
} from '@/js/meep-constants';

class IDSignInOAuthDialog extends Component {
  constructor(props) {
    super(props);
    this.state = {
    };
  }

  render() {
    return (
      <IDDialog
        title={this.props.title}
        open={this.props.open}
        onClose={this.props.onClose}
      >
        <Grid style={{ marginBottom: 20 }}>
          <GridCell span={6}>
            <Button style={styles.button} outlined onClick={() => this.props.onSignIn(OAUTH_PROVIDER_GITHUB)}>
              <img style={styles.icon} src={GitHubIcon}/>
              GitHub
            </Button>
          </GridCell>
          <GridCell span={6}>
            <Button style={styles.button} outlined onClick={() => this.props.onSignIn(OAUTH_PROVIDER_GITLAB)}>
              <img style={styles.icon} src={GitLabIcon}/>
              GitLab (EOL Account)
            </Button>
          </GridCell>
        </Grid>
        <span style={styles.text}>
          Authenticating with an external provider will:
          <ul>
            <li>Redirect the browser to the provider login page</li>
            <li>Request authorization to read your public user name</li>
            <li>Create your session on successful login and authorization</li>
          </ul>
          NOTE: Login & authorization may be seamless if already performed.
        </span>
      </IDDialog>
    );
  }
}

const styles = {
  button: {
    width: '100%',
    height: '48px',
    whiteSpace: 'nowrap'
  },
  icon: {
    height: '75%',
    marginRight: 10
  },
  text: {
    color: 'gray'
  }
};

export default IDSignInOAuthDialog;
