/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

// currentYear
export function currentYear(): number {
  return new Date().getFullYear();
}

export default helper(currentYear);
