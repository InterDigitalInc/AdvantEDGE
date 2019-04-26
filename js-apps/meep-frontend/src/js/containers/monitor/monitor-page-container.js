/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
import { connect } from 'react-redux';
import React, { Component }  from 'react';
import { Grid, GridCell, GridInner } from '@rmwc/grid';
import { Select } from '@rmwc/select';
import { Elevation } from '@rmwc/elevation';
import { Button } from '@rmwc/button';
import Iframe from 'react-iframe';
import HeadlineBar from '../../components/headline-bar';

import {
  uiSetAutomaticRefresh
} from '../../state/ui';

import {
  changeDashboardUrl
} from '../../state/monitor';

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

const DashboardContainer = (props) => {
  if (!props.dashboardUrl) {
    return null;
  }
  return (
    <Grid style={{width: '100%', height: '100%'}} >
      <GridInner style={{width: '100%', height: '100%'}}>
        <GridCell span={12} style={styles.inner}>
          <Elevation className="component-style" z={2} style={{width: '100%', height: '100%', display: 'flex', flexDirection: 'column'}}>                
            <div style={{flex: 1, padding: 10}}>
              <div data-cy={MON_DASHBOARD_IFRAME} style={{height: '100%'}}>
                <Iframe
                  url={props.dashboardUrl}
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
  );
};

const MonitorPageHeadlineBar = (props) => {
  return(
    <div style={{width: '100%'}}>   
      <Grid style={styles.headlineGrid}>
        <GridCell span={12}>
          <Elevation className="component-style" z={2} style={styles.headline}>
            <GridInner>
              <GridCell align={'middle'} span={6}>
                <HeadlineBar
                  titleLabel="Deployed Scenario"
                  scenarioName={props.scenarioName}
                />
              </GridCell>
              <GridCell span={4}>
                <Select
                  style={{width: '100%'}}
                  label='Dashboard'
                  outlined
                  options={selectOptions}
                  onChange={(e) => props.onChangeDashboard(e)}
                  value={props.dashboardUrl}
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
  );
};

class MonitorPageContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  handleSelectionChange(e) {
    this.props.changeDashboardUrl(e.target.value);
  }

  render() {
    return (
      <div style={{width: '100%', height: '100%'}}>
        <MonitorPageHeadlineBar
          scenarioName={this.props.scenarioName}
          onChangeDashboard={(e) => this.handleSelectionChange(e)}
        />
        <DashboardContainer dashboardUrl={this.props.dashboardUrl}/>
      </div>
    );
  }
}

const styles = {
  headlineGrid: {
    marginBottom: 10
  },
  headline: {
    paddingTop: 10,
    paddingRight: 10,
    paddingBottom: 10,
    paddingLeft: 10,
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
    dashboardUrl: state.monitor.dashboardUrl
  };
};

const mapDispatchToProps = dispatch => {
  return {
    setAutomaticRefresh: (val) => dispatch(uiSetAutomaticRefresh(val)),
    changeDashboardUrl: (url) => dispatch(changeDashboardUrl(url))
  };
};

const ConnectedMonitorPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(MonitorPageContainer);

export default ConnectedMonitorPageContainer;
