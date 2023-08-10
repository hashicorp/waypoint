/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';
import { formatDistance } from 'date-fns';

// dateFormatDistance
export function dateFormatDistance([date, baseDate]: [number, number]): string {
  if (!date || !baseDate) {
    return 'unknown';
  }
  let start = new Date(baseDate * 1000);
  let end = new Date(date * 1000);
  return formatDistance(end, start, { includeSeconds: true });
}

export default helper(dateFormatDistance);
