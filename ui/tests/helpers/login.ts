/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { authenticateSession } from 'ember-simple-auth/test-support';

export default async function login(waypointToken = 'default-test-token-value'): Promise<void> {
  return await authenticateSession({ token: waypointToken });
}

export function setupSession(hooks: NestedHooks): void {
  hooks.beforeEach(() => login());
}
