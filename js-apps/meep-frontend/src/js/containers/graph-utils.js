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

import _ from 'lodash';

export const blue = '#5DBCD2';

export const lineGeneratorNodes = n1 => n2 => {
  if (!n1 || !n2) {
    return '';
  }
  return `M${n1.X},${n1.Y} L${n2.X},${n2.Y}`;
};

export const plusGenerator = () => {
  const s = 2;
  return `M25 -20 h${s} v${2 * s} h${2 * s} v${s} h-${2 * s} v${2 *
    s} h-${s} v-${2 * s} h-${2 * s} v-${s} h${2 * s} z`;
};

export const minusGenerator = () => {
  const s = 4;
  return `M25 -20 h${3 * s} v${s} h-${3 * s} z`;
};

export const curveGeneratorNodes = n1 => n2 => {
  if (!n1 || !n2) {
    return '';
  }
  return `M${n1.X},${n1.Y} C${n1.X},${n2.Y + 150} ${n1.X},${n2.Y + 50} ${
    n2.X
  },${n2.Y}`;
};

export const visitNodes = f => node => {
  f(node);
  if (node.children) {
    _.each(node.children, c => {
      visitNodes(f)(c);
    });
  }
};

export const isNodeSelected = n => n.selected;
export const isNodeHighlighted = n => n.highlighted;
