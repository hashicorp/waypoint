/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

const wordBreak = /(?:^|\s|-|\/)\S/g;
const replace = /_+|-+/g;

// componentName
export function componentName([component]: [string]): string {
  if (!component) {
    return 'Unknown';
  }

  // Split into words based on common characters and uppercase
  // the first letter of each word
  let result = component.toLowerCase().replace(wordBreak, function (m) {
    return m.toUpperCase();
  });

  // Replace any separators that are not human readable
  result = result.replace(replace, ' ');

  // Replace brand initialisms
  result = result.replace('Aws', 'AWS');
  result = result.replace('Ecr', 'ECR');

  return result;
}

export default helper(componentName);
