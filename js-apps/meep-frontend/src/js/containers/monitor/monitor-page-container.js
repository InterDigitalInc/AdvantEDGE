import { connect } from 'react-redux';
import React, { Component }  from 'react';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Select } from '@rmwc/select';
import { Elevation } from '@rmwc/elevation';
import { Button } from '@rmwc/button';
import Iframe from 'react-iframe';
import HeadlineBar from '../../components/headline-bar';

import {
  uiSetAutomaticRefresh,
  PAGE_MONITOR
} from '../../state/ui';

import { 
  MON_DASHBOARD_SELECT,
  MON_DASHBOARD_IFRAME
} from '../../meep-constants';

const selectOptions = [
  {
    label:  'Latency Dashboard',
    value: 'http://' + location.hostname + ':32003/app/kibana#/dashboard/6745bb30-c29c-11e8-95a0-933bd4e05896?embed=true&_g=(refreshInterval%3A(pause%3A!f%2Cvalue%3A5000)%2Ctime%3A(from%3Anow-60s%2Cmode%3Arelative%2Cto%3Anow))',
  },
  {
    label:  'Demo Service Internal UE (ue1)',
    value: 'http://' + location.hostname + ':32003/app/kibana#/dashboard/434d37b0-1b6d-11e9-b72d-e70da2a5e139?embed=true&_g=(refreshInterval%3A(pause%3A!f%2Cvalue%3A5000)%2Ctime%3A(from%3Anow-15m%2Cmode%3Arelative%2Cto%3Anow))',
  },
  {
    label:  'Demo Service External UE (ue2-ext)',
    value: 'http://' + location.hostname + ':32003/app/kibana#/dashboard/788a4f70-1b73-11e9-b72d-e70da2a5e139?embed=true&_g=(refreshInterval%3A(pause%3A!f%2Cvalue%3A5000)%2Ctime%3A(from%3Anow-15m%2Cmode%3Arelative%2Cto%3Anow))',
  }
];

const kibanaDashboardUrl = 'http://' + location.hostname + ':32003/app/kibana#/dashboard';

class MonitorPageContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      iFrameUrl: '',
      dashboardSelected: null
    };
  }

  handleSelectionChange(e) {
    this.setState({
      iFrameUrl: e.target.value,
      dashboardSelected: true
    });
  }

  render() {

    if (this.props.page != PAGE_MONITOR) {
      return null;
    }
        
    return (
      <div style={{width: '100%', height: '100%'}}>
        <div style={{width: '100%'}}>   
          <Grid style={styles.headlineGrid}>
            <GridCell span={12}>
              <Elevation className="component-style" z={2} style={styles.headline}>
                <GridInner>
                  <GridCell align={'middle'} span={6}>
                    <HeadlineBar
                      titleLabel="Deployed Scenario"
                      scenarioName={this.props.scenarioName}
                    />
                  </GridCell>
                  <GridCell span={4}>
                    <Select
                      style={{width: '100%'}}
                      label='Dashboard'
                      outlined
                      options={selectOptions}
                      onChange={(e) => this.handleSelectionChange(e)}
                      data-cy={MON_DASHBOARD_SELECT}
                    />
                  </GridCell>
                  <GridCell span={2} style={{paddingTop: 8}}>
                    <Button raised
                      style={styles.button}
                      onClick={() => window.open(kibanaDashboardUrl, '_blank')}
                    >
                        OPEN KIBANA
                    </Button>
                  </GridCell>
                </GridInner>
              </Elevation>
            </GridCell>
          </Grid>
        </div>
                
        <Grid style={{width: '100%', height: '100%'}} hidden={!this.state.dashboardSelected}>
          <GridInner style={{width: '100%', height: '100%'}}>
            <GridCell span={12} style={styles.inner}>
              <Elevation className="component-style" z={2} style={{width: '100%', height: '100%', display: 'flex', flexDirection: 'column'}}>
                                    
                <div style={{flex: 1, padding: 10}}>
                  <div data-cy={MON_DASHBOARD_IFRAME} style={{height: '100%'}}>
                    <Iframe
                      url={this.state.iFrameUrl}
                      id="myId"
                      display="initial"
                      position="relative"
                      allowFullScreen
                      styles={{width: '100%', height: '100%'}}
                    />
                  </div>
                </div>
              </Elevation>
            </GridCell>
          </GridInner>
        </Grid>
                   
      </div>
    );
  }
}

const styles = {
  headlineGrid: {
    marginBottom: 10
  },
  headline: {
    padding: 10,
    marginBotton: 25
  },
  inner: {
    height: '100%',
  },
  page: {
    height: 1500,
    marginBottom: 10
  },
  cfgTable: {
    marginTop: 20,
    padding: 10
  },
  button: {
    color: 'white',
    marginRight: 5
  }
};

const mapStateToProps = state => {
  return {
    automaticRefresh: state.ui.automaticRefresh, 
    devMode: state.ui.devMode,
    page: state.ui.page,
    scenarioName: state.exec.scenario.name,
  };
};

const mapDispatchToProps = dispatch => {
  return {
    setAutomaticRefresh: (val) => dispatch(uiSetAutomaticRefresh(val))
  };
};

const ConnectedMonitorPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(MonitorPageContainer);

export default ConnectedMonitorPageContainer;
