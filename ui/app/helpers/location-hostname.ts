/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

// locationHostname returns a guess of the gRPC address based on the window hostname
export function locationHostname(): string {
  return `${window.location.hostname}:9701`;
}

export default helper(locationHostname);
