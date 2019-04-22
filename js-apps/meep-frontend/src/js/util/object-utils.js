/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */
export function deepCopy(source) {
  var dest = JSON.parse(JSON.stringify(source));
  return dest;
}
