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

import React from 'react';
import { shallow } from 'enzyme';
import toJson from 'enzyme-to-json';

// Component to be tested
import IDSelect from '../../../../src/js/components/helper-components/id-select';


describe('<IDSelect />', () => {

  describe('render()', () => {
    test('renders the component without props', () => {
      const wrapper = shallow(<IDSelect />);
      expect(wrapper.find('GridCell').get(0).props.span).toBeUndefined();
      expect(wrapper.find('Select').get(0).props.label).toBeUndefined();
      expect(wrapper.find('Select').get(0).props.options).toBeUndefined();
      expect(wrapper.find('Select').get(0).props.onChange).toBeUndefined();
      expect(wrapper.find('Select').get(0).props.disabled).toBeUndefined();
      expect(wrapper.find('Select').get(0).props.value).toBeUndefined();
    });

    test('renders the component with props', () => {
      const props = {
        span: 100,
        label: 'myLabel',
        options: 'myOptions',
        onChange: function myFunction(){},
        disabled: true,
        value: 'myValue'
      };
      const wrapper = shallow(<IDSelect {...props} />);
      expect(wrapper.find('GridCell').get(0).props.span).toEqual(props.span);
      expect(wrapper.find('Select').get(0).props.label).toEqual(props.label);
      expect(wrapper.find('Select').get(0).props.options).toEqual(props.options);
      expect(wrapper.find('Select').get(0).props.onChange).toEqual(props.onChange);
      expect(wrapper.find('Select').get(0).props.disabled).toEqual(props.disabled);
      expect(wrapper.find('Select').get(0).props.value).toEqual(props.value);
    });
  });

});
