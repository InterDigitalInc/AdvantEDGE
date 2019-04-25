/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
export const firstLetterUpper = (str) => {
  if (!str) return '';
  return str.charAt(0).toUpperCase() + str.slice(1);
};

export const camelCasePrefix = (prefix) => {
  if(!prefix) return '';
  const array = prefix.split('-');
  var f = array[0];
  array[0] = f.charAt(0).toLowerCase() + f.slice(1);
  return array.join('');
};
