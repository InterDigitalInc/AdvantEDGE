import React from 'react';

import {
  NO_SCENARIO_NAME,
  MEEP_LBL_SCENARIO_NAME
} from '../meep-constants';

const HeadlineBar = ({titleLabel, scenarioName}) => {
  var name = (scenarioName === NO_SCENARIO_NAME) ? 'None' : scenarioName;
  return (
    <div style={{marginTop: 14, marginBottom: 14}}>
      <span className="mdc-typography--headline6">{titleLabel}: </span>
      <span data-cy={MEEP_LBL_SCENARIO_NAME} className="idcc-margin-left mdc-typography--headline6 mdc-theme--primary">
        {name}
      </span>    
    </div>
  );
}; 

export default HeadlineBar;

