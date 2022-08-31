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

// import _ from 'lodash';
import { connect } from 'react-redux';
import React, { Component, createRef } from 'react';
import mermaid from 'mermaid';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';
import {
  execChangeSeqChart
} from '../state/exec';


class IDCSeq extends Component {
  constructor(props) {
    super(props);
    this.state = {};
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
    mermaid.init(undefined, '.seq-mermaid');
  }

  render() {
    return (
      <TransformWrapper>
        <TransformComponent
          wrapperStyle={{width: '100%', height: '100%'}}
        >
          <div
            ref={this.mermaidRef}
            className='seq-mermaid'
            data-cy={this.props.cydata}
          >
            {this.props.execSeqChart}
          </div>
        </TransformComponent>
      </TransformWrapper>
    );
  }
}

const mapStateToProps = state => {
  return {
    execSeqChart: state.exec.seq.chart
  };
};

const mapDispatchToProps = dispatch => {
  return {
    changeExecSeqChart: chart => dispatch(execChangeSeqChart(chart))
  };
};

const ConnectedIDCSeq = connect(
  mapStateToProps,
  mapDispatchToProps
)(IDCSeq);

export default ConnectedIDCSeq;
