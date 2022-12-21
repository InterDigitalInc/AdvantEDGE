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

import React, { Component }  from 'react';
import { connect } from 'react-redux';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
import { Typography } from '@rmwc/typography';
import '../../../img/AdvantEDGE-logo_Blue-01.png';

import {
  MEEP_DOCS_URL,
  MEEP_ARCHITECTURE_URL,
  MEEP_USAGE_URL,
  MEEP_HELP_GUI_URL,
  MEEP_CONTRIBUTING_URL,
  MEEP_DISCUSSIONS_URL,
  MEEP_ISSUES_URL,
  MEEP_LICENSE_URL
} from '../../meep-constants';

class HomePageContainer extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Grid>

        <GridCell span={2}/>
        <GridCell span={8}>

          <GridInner align='center'>
            <GridCell style={{ marginBottom: 10 }} span={12}>
              <Elevation className='idcc-elevation' style={styles.banner} z={3}>
                <div style={styles.bannerImage}/>
              </Elevation>
            </GridCell>
          </GridInner>

          <GridInner>
            <GridCell span={12}>
              <Elevation className='idcc-elevation' z={3}>
                <div style={{ padding: 30 }}>
                  <div style={styles.title}>
                    <Typography theme='primary' use='headline4'>Welcome to AdvantEDGE</Typography>
                  </div>

                  <div>
                    <Typography className='text-color-main' use='body1'>
                      <p> AdvantEDGE is a Mobile Edge Emulation Platform (MEEP) that runs on Docker & Kubernetes.</p>
                      <p>
                        AdvantEDGE provides an emulation environment, enabling experimentation with Edge
                        Computing Technologies, Applications, and Services. The platform facilitates exploring
                        edge / fog deployment models and their impact on applications and services in short
                        and agile iterations.
                      </p>
                    </Typography>

                    <Typography theme='primary' use='headline6'>Motivation</Typography>

                    <Typography className='text-color-main' use='body1'>
                      <ul>
                        <li>Accelerate Mobile Edge Computing adoption</li>
                        <li>Help Discover new edge application use cases & services</li>
                        <li>Help answer questions such as:
                          <ul>
                            <li>Where should my application components be located in the edge network?</li>
                            <li>How do network characteristics (such as latency, jitter, and packet loss) impact my application or service?</li>
                            <li>How will my application behave when the user moves within and across access networks?</li>
                          </ul>
                        </li>
                      </ul>
                    </Typography>

                    <Typography theme='primary' use='headline6'>Intended Users</Typography>

                    <Typography className='text-color-main' use='body1'>
                      <ul>
                        <li>Edge Application Developers</li>
                        <li>Edge Network and Service Designers</li>
                        <li>Edge Researchers</li>
                        <li>Technologists that are simply interested learning how the Edge works</li>
                      </ul>
                    </Typography>

                    <Typography theme='primary' use='headline6'>Getting started</Typography>

                    <Typography className='text-color-main' use='body1'>
                      <p>If you made it here, AdvantEDGE was successfully installed. Go ahead and experiment with the platform!</p>
                      <p>Need some help?</p>
                      <ul>
                        <li>Check out the <a className='idcc-link' href={MEEP_DOCS_URL} target='_blank'>AdvantEDGE Docs</a></li>
                        <li>Learn more about platform <a className='idcc-link' href={MEEP_ARCHITECTURE_URL} target='_blank'>Architecture & Features</a></li>
                        <li>Follow the platform <a className='idcc-link' href={MEEP_USAGE_URL} target='_blank'>Usage Workflow</a></li>
                        <li>Have a look at the <a className='idcc-link' href={MEEP_HELP_GUI_URL} target='_blank'>GUI Help Page</a></li>
                      </ul>
                    </Typography>

                    <Typography theme='primary' use='headline6'>How to Contribute</Typography>

                    <Typography className='text-color-main' use='body1'>
                      <p>
                        If you like this project and would like to participate in its evolution, you can find information
                        on contributing <a className='idcc-link' href={MEEP_CONTRIBUTING_URL} target='_blank'>here</a>.
                      </p>
                      <p>
                        We welcome questions, feedback and improvement suggestions
                        via <a className='idcc-link' href={MEEP_DISCUSSIONS_URL} target='_blank'>Discussions</a> and bug reporting
                        via <a className='idcc-link' href={MEEP_ISSUES_URL} target='_blank'>Issues</a>.
                      </p>
                      <p>Hope to hear from you...</p>
                    </Typography>

                    <Typography theme='primary' use='headline6'>Licensing</Typography>

                    <Typography className='text-color-main' use='body1'>
                      <p>
                        AdvantEDGE is licensed under under
                        the <a className='idcc-link' href={MEEP_LICENSE_URL} target='_blank'>Apache License, Version 2.0</a>
                      </p>
                    </Typography>
                  </div>
                </div>
              </Elevation>
            </GridCell>
          </GridInner>

        </GridCell>
      </Grid>
    );
  }
}

const styles = {
  banner: {
    height: 150,
    padding: 10,
    backgroundColor: 'white'
  },
  bannerImage: {
    backgroundSize: 'contain',
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundImage: 'url(../../../img/AdvantEDGE-logo_Blue-01.png)',
    height: '100%',
    width: '100%'
  },
  text: {
    color: 'black',
    fontSize: '1.2rem'
  },
  title: {
    marginTop: 10,
    marginBottom: 20
  }
};

const mapStateToProps = () => {
  return {
  };
};

const ConnectedHomePageContainer = connect(
  mapStateToProps
)(HomePageContainer);

export default ConnectedHomePageContainer;
