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

import React, { Component }  from 'react';
import { connect } from 'react-redux';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Elevation } from '@rmwc/elevation';
import '../../../img/AdvantEDGE-logo_Blue-01.png';

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
                <div style={{ padding: 20 }}>
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
