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

import { connect } from 'react-redux';
import React, { Component }  from 'react';

import {
  MEEP_TAB_CFG,
  MEEP_TAB_EXEC,
  MEEP_TAB_MON,
  MEEP_TAB_SET,
  PAGE_CONFIGURE,
  PAGE_EXECUTE,
  PAGE_MONITOR,
  PAGE_SETTINGS
} from '../meep-constants';

import { uiChangeCurrentPage } from '../state/ui';

import {
  List,
  ListItem
} from '@rmwc/list';

import { Icon } from '@rmwc/icon';

class MeepDrawer extends Component {

  constructor(props) {
    super(props);
    this.state = {
      dismissibleOpen: true
    };
  }

  handleItemClick(page) {
    this.props.changeCurrentPage(page);
  }

  styleForPage(page) {
    var style = (this.props.currentPage === page) ? {backgroundColor: '#E0F0F9'} : null;
    return style;
  }

  render() {
    return (
      <div className="component-style" style={{overflow: 'hidden', position: 'relative'}}>
        <div
          style={containerStyle}
          open={this.props.open}>
          <div style={{marginTop: '-80px'}}>
            <div className="idcc-margin-top mdc-top-app-bar--fixed-adjust"></div>
            <List>
              <ListItem data-cy={MEEP_TAB_CFG} style={this.styleForPage(PAGE_CONFIGURE)} onClick = {() => {this.handleItemClick(PAGE_CONFIGURE);}}>
                <Icon icon="sort" iconOptions={{strategy: 'ligature'}} style={iconStyles}/>
                <span style={textStyles}>Configure</span>
              </ListItem>
              <ListItem data-cy={MEEP_TAB_EXEC} style={this.styleForPage(PAGE_EXECUTE)} onClick = {() => {this.handleItemClick(PAGE_EXECUTE);}}>
                <Icon icon="forward" iconOptions={{strategy: 'ligature'}} style={iconStyles}/>
                <span style={textStyles}>Execute</span>
              </ListItem>
              <ListItem data-cy={MEEP_TAB_MON} style={this.styleForPage(PAGE_MONITOR)} onClick = {() => {this.handleItemClick(PAGE_MONITOR);}}>
                <Icon icon="tv" iconOptions={{strategy: 'ligature'}} style={iconStyles}/>
                <span style={textStyles}>Monitor</span>
              </ListItem>
              <ListItem data-cy={MEEP_TAB_SET} style={this.styleForPage(PAGE_SETTINGS)} onClick = {() => {this.handleItemClick(PAGE_SETTINGS);}}>
                <Icon icon="settings" iconOptions={{strategy: 'ligature'}} style={iconStyles}/>
                <span style={textStyles}>Settings</span>
              </ListItem>
            </List>
          </div>
        </div>
      </div>
    );
  }
}

const containerStyle = {
  borderRight: '1px solid #e4e4e4',
  height: '100vh'
};

const textStyles = {
  marginLeft: '30px',
  marginRight: '90px',
  fontSize: 14
};

const iconStyles = {
  color: 'gray'
};

const mapDispatchToProps = dispatch => {
  return {
    changeCurrentPage: (page) => dispatch(uiChangeCurrentPage(page))
  };
};

const mapStateToProps = state => {
  return {
    currentPage: state.ui.page,
    mainDrawerOpen: state.ui.mainDrawerOpen
  };
};

const ConnectedMeepDrawer = connect(
  mapStateToProps,
  mapDispatchToProps
)(MeepDrawer);

export default ConnectedMeepDrawer;
