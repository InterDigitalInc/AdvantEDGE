
import React from 'react';
import { shallow } from 'enzyme';
import toJson from 'enzyme-to-json';

// Component to be tested
import { CfgNetworkElementContainer } from '../../../../src/js/containers/cfg/cfg-network-element-container';


describe('<CfgNetworkElementContainer />', () => {

  describe('render()', () => {
    test('renders the component without a configured element', () => {
      const wrapper = shallow(<CfgNetworkElementContainer />);
      expect(wrapper.getElement()).toBe(null);
    });

    // test('renders the component with a configured element', () => {
    //     const wrapper = shallow(<CfgNetworkElementContainer ??? />);
    // });
  });

});

