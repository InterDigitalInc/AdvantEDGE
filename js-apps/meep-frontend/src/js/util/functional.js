/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

export const pipe = (...fns) => val => fns.reduce((acc, f) => f(acc), val);
export const filter = fn => array => array.filter(fn);

export const idlog = label => val => {
  /*eslint-disable */
  console.log(`${label}: `, val);
  /*eslint-enable */
  return val;
};
