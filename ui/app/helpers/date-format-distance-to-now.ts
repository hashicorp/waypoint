/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';
import { formatDistanceToNow } from 'date-fns';

// dateFormatDistanceToNow
export function dateFormatDistanceToNow([date]: [number]): string {
  if (!date) {
    return 'unknown';
  }

  return formatDistanceToNow(date * 1000, { includeSeconds: true, addSuffix: true });
}

export default helper(dateFormatDistanceToNow);
