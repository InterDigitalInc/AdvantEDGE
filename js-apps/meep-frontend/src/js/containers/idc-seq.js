/*
 * Copyright (c) 2022  InterDigital Communications, Inc
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
import { connect } from 'react-redux';
import React, { Component, createRef } from 'react';
import mermaid from 'mermaid';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';

class IDCSeq extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.mermaidChart = '';
    this.mermaidRef = createRef();

    mermaid.initialize({
      startOnLoad: true,
      useMaxWidth: true
    });
  }

  componentDidMount() {
    mermaid.init(undefined, '.seq-mermaid');
  }

  componentDidUpdate() {
    // Remove data-processed attribute to allow diagram refresh
    this.mermaidRef.current.removeAttribute('data-processed');
    this.mermaidRef.current.innerHTML = this.mermaidChart;
    mermaid.init(undefined, this.mermaidRef.current);
  }

  formatSequenceChart() {
    var seqChart = '';

    // Return default diagram if no metrics available yet
    if (this.props.execSeqMetrics.length === 0) {
      // Default diagram
      seqChart = 'flowchart LR\nid1(Sequence diagram waiting for metrics to display...)';
    } else {
      seqChart = 'sequenceDiagram\n';

      // Add dashboard-configured participants 
      if (this.props.participants) {
        var dashParticipants = _.split(this.props.participants, ',');
        _.forEach(dashParticipants, participant => {
          if (participant) {
            seqChart += ('participant ' + participant + '\n');
          }
        });
      }

      // Add scenario-configured participants
      _.forEach(this.props.execSeqParticipants, participant => {
        seqChart += ('participant ' + participant + '\n');
      });

      // Add metrics
      _.forEach(this.props.execSeqMetrics, metric => {
        seqChart += (metric.mermaid + '\n');
      });
    }

    this.mermaidChart = seqChart;
  }

  render() {
    // Format sequence diagram
    this.formatSequenceChart();

    return (
      <TransformWrapper>
        <TransformComponent
          wrapperStyle={{width: '100%', height: '100%', textAlign: 'bottom'}}
        >
          <div
            ref={this.mermaidRef}
            className='seq-mermaid'
            data-cy={this.props.cydata}
          >
            {this.mermaidChart}
          </div>
        </TransformComponent>
      </TransformWrapper>
    );
  }
}

const mapStateToProps = state => {
  return {
    execSeqMetrics: state.exec.seq.metrics,
    execSeqParticipants: state.exec.seq.participants,
    execSeqChart: state.exec.seq.chart
  };
};

const mapDispatchToProps = () => {
  return {
  };
};

const ConnectedIDCSeq = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCSeq);

export default ConnectedIDCSeq;
