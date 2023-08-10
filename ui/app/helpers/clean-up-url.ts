/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

export function cleanUpUrl([str]: string[]): string {
  let cleanStr = str.replace('https://', '').replace('http://', '').replace('www.', '');
  return cleanStr;
}

export default helper(cleanUpUrl);
