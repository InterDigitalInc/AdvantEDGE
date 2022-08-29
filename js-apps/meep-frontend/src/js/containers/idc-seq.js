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


// %%{init: {'theme': 'base', 'themeVariables': { 'actorBkg': '#FF9800'}}}%%
const SEQ_DEFAULT = `
sequenceDiagram
edn2_ees1->>ecs1: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs1-->>edn2_ees1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
edn2_ees1->>ecs2: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs2-->>edn2_ees1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs1->>ecs2: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs2-->>ecs1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs2->>ecs3: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs3-->>ecs2: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs3->>ecs4: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs4-->>ecs3: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs4->>ecs5: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs5-->>ecs4: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
edn2_ees1->>ecs1: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs1-->>edn2_ees1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
edn2_ees1->>ecs2: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs2-->>edn2_ees1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs1->>ecs2: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs2-->>ecs1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs2->>ecs3: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs3-->>ecs2: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs3->>ecs4: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs4-->>ecs3: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs4->>ecs5: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs5-->>ecs4: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
edn2_ees1->>ecs1: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs1-->>edn2_ees1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
edn2_ees1->>ecs2: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs2-->>edn2_ees1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs1->>ecs2: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs2-->>ecs1: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs2->>ecs3: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs3-->>ecs2: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs3->>ecs4: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs4-->>ecs3: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
ecs4->>ecs5: [21:26:23.495] /eecs-eesregistration/v1/registrations
ecs5-->>ecs4: [21:26:23.516] /eecs-eesregistration/v1/registrations (22.09 ms)
`;

class IDCSeq extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.mermaidRef = createRef();

    this.props.changeExecSeqChart(SEQ_DEFAULT);

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

  // initMermaid () {
  //   const {
  //     code,
  //     history,
  //     match: { url }
  //   } = this.props
  //   try {
  //     mermaid.parse(code)
  //     // Replacing special characters '<' and '>' with encoded '&lt;' and '&gt;'
  //     let _code = code
  //     _code = _code.replace(/</g, '&lt;')
  //     _code = _code.replace(/>/g, '&gt;')
  //     // Overriding the innerHTML with the updated code string
  //     this.container.innerHTML = _code
  //     mermaid.init(undefined, this.container)
  //   } catch (e) {
  //     // {str, hash}
  //     const base64 = Base64.encodeURI(e.str || e.message)
  //     history.push(`${url}/error/${base64}`)
  //   }
  // }

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
