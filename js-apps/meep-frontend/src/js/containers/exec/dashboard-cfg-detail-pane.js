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

import _ from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import autoBind from 'react-autobind';
import { Grid, GridCell } from '@rmwc/grid';
import { Button } from '@rmwc/button';
import { Checkbox } from '@rmwc/checkbox';
import { TextField, TextFieldIcon, TextFieldHelperText } from '@rmwc/textfield';
import { updateObject } from '@/js/util/object-util';

import IDSelect from '../../components/helper-components/id-select';


import {
  NET_TOPOLOGY_VIEW,
  SEQ_DIAGRAM_VIEW,
  DATAFLOW_DIAGRAM_VIEW,
  NET_METRICS_PTP_VIEW,
  NET_METRICS_AGG_VIEW,
  WIRELESS_METRICS_PTP_VIEW,
  WIRELESS_METRICS_AGG_VIEW,
  CFG_BTN_EXP_MERMAID_SEQ,
  CFG_BTN_EXP_SDORG_SEQ,
  CFG_BTN_EXP_MERMAID_DF
} from '@/js/meep-constants';

import {
  uiExecChangeShowApps,
  uiExecChangePauseSeq,
  uiExecChangePauseDataflow
} from '@/js/state/ui';

import {
  DASH_CFG_VIEW_TYPE,
  DASH_CFG_SOURCE_NODE_SELECTED,
  DASH_CFG_DEST_NODE_SELECTED,
  DASH_CFG_PARTICIPANTS,
  DASH_CFG_MAX_MSG_COUNT,
  getDashCfgFieldVal,
  setDashCfgField,
  DASH_CFG_START_TIME
} from '@/js/util/dashboard-utils';

// COMPONENTS
const DashCfgTextField = props => {
  var err = props.viewCfg[props.fieldName] ? props.viewCfg[props.fieldName].err: null;

  return (
    <>
      <TextField
        outlined
        style={{ width: '100%', marginBottom: 0 }}
        label={props.label}
        withLeadingIcon={!props.icon ? null :
          <TextFieldIcon
            tabIndex="0"
            icon={props.icon}
            onClick={props.onIconClick}
          />
        }
        type={props.type}
        onChange={event => {
          var err = props.validate ? props.validate(event.target.value) : null;
          var val =
            event.target.value && props.isNumber && !err
              ? Number(event.target.value)
              : event.target.value;
          props.onUpdate(props.fieldName, val, err);
        }}
        invalid={err}
        value={
          props.viewCfg[props.fieldName]
            ? props.viewCfg[props.fieldName].val
            : ''
        }
        disabled={props.disabled}
        data-cy={props.cydata}
      />
      <TextFieldHelperText validationMsg={true}>
        <span>
          {props.viewCfg[props.fieldName] ? props.viewCfg[props.fieldName].err : ''}
        </span>
      </TextFieldHelperText>
    </>
  );
};

class DashCfgDetailPane extends Component {
  constructor(props) {
    super(props);
    autoBind(this);
  }

  onClose(e) {
    e.preventDefault();
    this.props.onClose(e);
  }

  onUpdateDashCfg(name, val, err) {
    var newViewCfg = updateObject({}, this.props.viewCfg);
    setDashCfgField(newViewCfg, name, val, err);
    this.props.changeViewCfg(this.props.currentView, newViewCfg);
  }

  onUpdateViewType(event) {
    this.onUpdateDashCfg(DASH_CFG_VIEW_TYPE, event.target.value, null);
  }
  onUpdateSourceNodeSelected(event) {
    this.onUpdateDashCfg(DASH_CFG_SOURCE_NODE_SELECTED, event.target.value, null);
  }
  onUpdateDestNodeSelected(event) {
    this.onUpdateDashCfg(DASH_CFG_DEST_NODE_SELECTED, event.target.value, null);
  }
  onChangeShowApps(event) {
    this.props.changeShowApps(event.target.checked);
  }
  onChangePauseSeq(event) {
    this.props.changePauseSeq(event.target.checked);
  }
  onChangePauseDataflow(event) {
    this.props.changePauseDataflow(event.target.checked);
  }
  onClickClear() {
    // Obtain last metric timestamp
    const metrics = this.props.execSeqMetrics;
    var lastMetricTime = (metrics && metrics.length > 0) ? metrics[metrics.length-1].time : '';

    // Set start time equal to last metric timestamp
    this.onUpdateDashCfg(DASH_CFG_START_TIME, lastMetricTime, null);
  }
  onClickFetchAll() {
    // Reset start time
    this.onUpdateDashCfg(DASH_CFG_START_TIME, '', null);
  }

  // Callback function for postSeqQuery
  postSeqQueryCb(error, data, format) {
    if (error !== null) {
      return;
    }
    var seqMetrics = (data && data.seqMetricString) ? data.seqMetricString : '';

    // Nothing to do if no new metrics
    if (seqMetrics.length === 0) {
      return;
    }
    var seqChart = '';
    // Add scenario-configured participants
    var participants = getDashCfgFieldVal(this.viewCfg, DASH_CFG_PARTICIPANTS);
    participants = _.split(participants, ',');
    _.forEach(participants, participant => {
      seqChart += ('participant ' + participant + '\n');
    });
    var scenarioName = this.props.scenarioName;
    if (format === 'mermaid'){
      seqChart = 'sequenceDiagram\n' + seqChart + seqMetrics;
      var filename = scenarioName + '-mermaid-seq.txt';
      var id = 'export-mermaid-link';
    }
    else{
      seqChart = seqChart + seqMetrics;
      filename = scenarioName + '-sdorg-seq.txt';
      id = 'export-sdorg-link';
    }
    var link = document.getElementById(id);
    link.href = this.makeTextFile(seqChart);
    link.download = filename;
    link.click();
  }

  // Export Mermaid Sequence Diagram
  exportMermaidSeq() {
    var seqMetricsQuery = {
      fields: ['mermaid'],
      scope: {
        limit: 0
      },
      responseType: 'stronly'
    };
    var format = 'mermaid';
    // Query sequence diagram
    this.props.metricsApi.postSeqQuery(seqMetricsQuery, (error, data) =>
      this.postSeqQueryCb(error, data, format)
    );
  }
  // Export Sdorg Sequence Diagram
  exportSdorgSeq() {
    var seqMetricsQuery = {
      fields: ['sdorg'],
      scope: {
        limit: 0
      },
      responseType: 'stronly'
    };
    var format = 'sdorg';
    // Query sequence diagram
    this.props.metricsApi.postSeqQuery(seqMetricsQuery, (error, data) =>
      this.postSeqQueryCb(error, data, format)
    );
  }

  // Callback function for postDataflowQuery
  postDataflowQueryCb(error, data) {
    if (error !== null) {
      return;
    }
    var dataflowChart = (data) ? data.dataflowMetricString : '';
    dataflowChart = 'stateDiagram\n' + dataflowChart;
    var filename = this.props.scenarioName + '-mermaid-df.txt';
    var id = 'export-dataflow-link';
    var link = document.getElementById(id);
    link.href = this.makeTextFile(dataflowChart);
    link.download = filename;
    link.click();
  }

  // Export Mermaid Dataflow Diagram
  exportMermaidDataflow() {
    var dataflowMetricsQuery = {
      fields: ['mermaid'],
      responseType: 'stronly'
    };
    this.props.metricsApi.postDataflowQuery(dataflowMetricsQuery, (error, data) =>
      this.postDataflowQueryCb(error, data)
    );
  }

  makeTextFile(text) {
    var data = new Blob([text], { type: 'text/plain'});
    var exportTextFile = window.URL.createObjectURL(data);
    return exportTextFile;
  }
  render() {
    this.viewCfg = this.props.viewCfg;
    this.netSrcNodeIds = this.props.appIds;
    this.netDstNodeIds = this.props.appIds;
    this.wirelessSrcNodeIds = this.props.ueIds;
    this.wirelessDstNodeIds = this.props.poaIds;
    this.viewType = getDashCfgFieldVal(this.viewCfg, DASH_CFG_VIEW_TYPE);
    this.sourceNodeSelected = getDashCfgFieldVal(this.viewCfg, DASH_CFG_SOURCE_NODE_SELECTED);
    this.destNodeSelected = getDashCfgFieldVal(this.viewCfg, DASH_CFG_DEST_NODE_SELECTED);

    return (
      <div>
        <Grid>
          <GridCell span={12} style={{marginBottom: 10}}>
            <IDSelect
              label={'View Type'}
              outlined
              options={this.props.dashboardViewsList}
              onChange={this.onUpdateViewType}
              value={this.viewType}
            />
          </GridCell>

          { this.viewType === NET_TOPOLOGY_VIEW ?
            <>
              <GridCell span={6}>
                <Checkbox
                  checked={this.props.showApps}
                  onChange={this.onChangeShowApps}
                >
                  Show Apps
                </Checkbox>
              </GridCell>
            </> : null
          }

          { this.viewType === SEQ_DIAGRAM_VIEW ?
            <>
              <GridCell span={12}>
                <DashCfgTextField
                  onUpdate={this.onUpdateDashCfg}
                  viewCfg={this.viewCfg}
                  // validate={validateParticipants}
                  isNumber={false}
                  label={'Participants'}
                  fieldName={DASH_CFG_PARTICIPANTS}
                />
              </GridCell>
              <GridCell span={12}>
                <DashCfgTextField
                  onUpdate={this.onUpdateDashCfg}
                  viewCfg={this.viewCfg}
                  // validate={validateParticipants}
                  isNumber={true}
                  label={'Max message count'}
                  fieldName={DASH_CFG_MAX_MSG_COUNT}
                />
              </GridCell>
              <GridCell span={12} style={{marginBottom: 10}}>
                <Button
                  outlined
                  onClick={this.onClickClear}
                >
                  Clear
                </Button>
                <Button
                  outlined
                  style={{marginLeft: 10}}
                  onClick={this.onClickFetchAll}
                >
                  Fetch all
                </Button>
                <Button
                  outlined
                  style={{marginLeft: 10}}
                  onClick={this.exportMermaidSeq}
                  data-cy={CFG_BTN_EXP_MERMAID_SEQ}
                >
                  Export mermaid
                </Button>
                <a id='export-mermaid-link' download='mermaid.txt' hidden></a>
                <Button
                  outlined
                  style={{marginLeft: 10}}
                  onClick={this.exportSdorgSeq}
                  data-cy={CFG_BTN_EXP_SDORG_SEQ}
                >
                  Export sdorg
                </Button>
                <a id='export-sdorg-link' download='sdorg.txt' hidden></a>
              </GridCell>
              <GridCell span={12}>
                <Checkbox
                  // style={{marginTop: 20}}
                  checked={this.props.pauseSeq}
                  onChange={this.onChangePauseSeq}
                >
                  Pause
                </Checkbox>
              </GridCell>
            </> : null
          }

          { this.viewType === DATAFLOW_DIAGRAM_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <Button
                  outlined
                  style={{marginLeft: 10}}
                  onClick={this.exportMermaidDataflow}
                  data-cy={CFG_BTN_EXP_MERMAID_DF}
                >
                  Export mermaid
                </Button>
                <a id='export-dataflow-link' download='df.txt' hidden></a>
              </GridCell>
              <GridCell span={12}>
                <Checkbox
                  checked={this.props.pauseDataflow}
                  onChange={this.onChangePauseDataflow}
                >
                  Pause
                </Checkbox>
              </GridCell>
            </> : null
          }

          { this.viewType === NET_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={this.netSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.netSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'Destination Node'}
                  outlined
                  options={this.netDstNodeIds}
                  onChange={this.onUpdateDestNodeSelected}
                  value={this.netDstNodeIds.includes(this.destNodeSelected) ? this.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.viewType === NET_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'Source Node'}
                  outlined
                  options={this.netSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.netSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.viewType === WIRELESS_METRICS_PTP_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'UE'}
                  outlined
                  options={this.wirelessSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.wirelessSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'POA'}
                  outlined
                  options={this.wirelessDstNodeIds}
                  onChange={this.onUpdateDestNodeSelected}
                  value={this.wirelessDstNodeIds.includes(this.destNodeSelected) ? this.destNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

          { this.viewType === WIRELESS_METRICS_AGG_VIEW ?
            <>
              <GridCell span={12} style={{marginBottom: 10}}>
                <IDSelect
                  label={'UE'}
                  outlined
                  options={this.wirelessSrcNodeIds}
                  onChange={this.onUpdateSourceNodeSelected}
                  value={this.wirelessSrcNodeIds.includes(this.sourceNodeSelected) ? this.sourceNodeSelected : 'None'}
                />
              </GridCell>
            </> : null
          }

        </Grid>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    execSeqMetrics: state.exec.seq.metrics,
    showApps: state.ui.execShowApps,
    pauseSeq: state.ui.execPauseSeq,
    pauseDataflow: state.ui.execPauseDataflow,
    scenarioName: state.exec.scenario.name
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeShowApps: show => dispatch(uiExecChangeShowApps(show)),
    changePauseSeq: pause => dispatch(uiExecChangePauseSeq(pause)),
    changePauseDataflow: pause => dispatch(uiExecChangePauseDataflow(pause))
  };
};

const ConnectedDashCfgDetailPane = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashCfgDetailPane);

export default ConnectedDashCfgDetailPane;
