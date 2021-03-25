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
import { connect } from 'react-redux';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Button } from '@rmwc/button';
import { Elevation } from '@rmwc/elevation';

import GitHubIcon from '@/img/logo-github.svg';
import AdvantEdgeLogo from '@/img/AdvantEDGE-logo_Blue-01.png';
import BackgroundWP from '@/img/network.png';

import {
  OAUTH_PROVIDER_GITHUB,
  STATUS_SIGNED_OUT
} from '@/js/meep-constants';

class LoginPageContainer extends Component {
  constructor(props) {
    super(props);
  }

  componentDidUpdate(prevProps) {
    prevProps.signInStatus !== this.props.signInStatus ? this.updateLoginBox() : null;
  }

  componentDidMount() {
    this.props.signInStatus !== STATUS_SIGNED_OUT ? this.updateLoginBox() : null;
  }

  updateLoginBox() {
    let elevation = document.getElementById('elevationTag');
    if (this.props.signInStatus !== STATUS_SIGNED_OUT) {
      elevation.style.width = '53%';
      elevation.style.marginLeft = '25%';
    } else {
      elevation.style.width = '80%';
      elevation.style.marginLeft = '10%';
    }
  }

  render() {
    let signedOut = this.props.signInStatus === STATUS_SIGNED_OUT;
    return (
      <div style={{ overflow: 'hidden', position: 'relative' }}>
        <img
          src={BackgroundWP}
          style={styles.background}
        />  
        <div style={{ position: 'relative' }}>
          <img
            src={AdvantEdgeLogo}
            style={styles.logo}
          />
          <Elevation z={3} style={styles.elevation} id='elevationTag'>
            <Grid style={ signedOut ? styles.gridLine : null}>
              <GridCell span={signedOut ? 8 : 12}>
                <span style={styles.text}>
                  <p> AdvantEDGE is a Mobile Edge Emulation Platform (MEEP) that runs on Docker & Kubernetes.</p>
                  <p>
                    AdvantEDGE provides an emulation environment, enabling experimentation with Edge
                    Computing Technologies, Applications, and Services. The platform facilitates exploring
                    edge / fog deployment models and their impact on applications and services in short
                    and agile iterations.
                  </p>
                  <h3>Motivation</h3>
                  <ul>
                    <li>Accelerate Mobile Edge Computing adoption</li>
                    <li>Help Discover new edge application use cases & services</li>
                    <li>
                      Help answer questions such as:
                      <ul>
                        <li>Where should my application components be located in the edge network?</li>
                        <li>How do network characteristics (such as latency, jitter, and packet loss) impact my application or service?</li>
                        <li>How will my application behave when the user moves within and across access networks?</li>
                      </ul>
                    </li>
                  </ul>
                </span>
              </GridCell>
              {signedOut ?
                <GridCell span={4}>
                  <div>
                    <span style={styles.text}>
                      <p>Authenticating with an external provider will:</p>
                      <ul>
                        <li>Redirect the browser to the provider login page</li>
                        <li>Request authorization to read your public user name</li>
                        <li>Allows you to use AdvantEDGE on successful login and authorization</li>
                      </ul>
                      NOTE: Login & authorization may be seamless if already performed.
                    </span>
                    <Button
                      style={styles.button}
                      outlined
                      onClick={() => this.props.onSignIn(OAUTH_PROVIDER_GITHUB)}
                    >
                      <img style={styles.icon} src={GitHubIcon}/>
                      Sign in with GitHub
                    </Button>
                  </div>
                </GridCell>
                : null
              }
            </Grid>
          </Elevation>
          <div style={styles.footer}>
            <hr/>
            <Grid>
              <GridInner>
                <GridCell span="2"></GridCell>
                <GridCell span="10">
                  <GridInner>
                    <GridCell span="2">
                      <a href="https://github.com/InterDigitalInc/AdvantEDGE/wiki" style={styles.headerText}>
                        Wiki
                      </a>
                    </GridCell>
                    <GridCell span="2">
                      <a href="https://github.com/InterDigitalInc/AdvantEDGE" style={styles.headerText}>
                        Github
                      </a>
                    </GridCell>
                    <GridCell span="2">
                      <a href="https://github.com/InterDigitalInc/AdvantEDGE/discussions" style={styles.headerText}>
                        Discussions
                      </a>
                    </GridCell>
                    <GridCell span="2">
                      <a href="https://github.com/InterDigitalInc/AdvantEDGE/blob/master/LICENSE" style={styles.headerText}>
                        License
                      </a>
                    </GridCell>
                    <GridCell span="2">
                      <a href="https://github.com/InterDigitalInc/AdvantEDGE/blob/master/CONTRIBUTING.md" style={styles.headerText}>
                        Contributing
                      </a>
                    </GridCell>
                  </GridInner>
                </GridCell>
              </GridInner>
            </Grid>
          </div>
        </div>
      </div>
    );
  }
}

const styles = {
  button: {
    width: '100%',
    height: '50px',
    whiteSpace: 'nowrap',
    marginTop: 30
  },
  icon: {
    height: '75%',
    marginRight: 10
  },
  text: {
    color: 'black',
    fontSize: '1.2rem'
  },
  elevation: {
    padding: '30px',
    width: '80%',
    marginLeft: '10%',
    marginTop: '3%',
    background: 'white'
  },
  logo: {
    height: 120,
    width: 500,
    marginLeft: '40%'
  },
  headerText: {
    fontFamily: 'sans-serif',
    marginLeft: '23%',
    color: 'black',
    fontSize: '1.3rem'
  },
  footer: {
    marginTop: '9%',
    marginLeft: '15%',
    width: '75%'
  },
  gridLine: {
    background: 'linear-gradient(#9d9d9d,#9d9d9d) center/2px 100% no-repeat',
    backgroundPosition: '66% 0'
  },
  background: {
    opacity: '6%',
    position: 'absolute',
    left: 0,
    top: 0,
    width: '100%',
    height: 'auto'
  }
};

const mapStateToProps = state => {
  return {
    signInStatus: state.ui.signInStatus
  };
};

const ConnectedLoginPageContainer = connect(
  mapStateToProps
)(LoginPageContainer);

export default ConnectedLoginPageContainer;
