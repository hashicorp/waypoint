/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';

export function enforceProtocol(params: [string] /*, hash*/): string {
  let str = params[0];

  let isHttps = str.startsWith('https://');
  let isHttp = str.startsWith('http://');
  if (isHttps || isHttp) {
    return str;
  } else {
    return `https://${str}`;
  }
}

export default helper(enforceProtocol);
