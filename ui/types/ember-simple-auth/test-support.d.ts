/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import SessionService from 'ember-simple-auth/services/session';

declare module 'ember-simple-auth/test-support' {
  export function authenticateSession(sessionData?: Record<string, unknown>): Promise<void>;
  export function currentSession(): SessionService;
  export function invalidateSession(): Promise<void>;
}
