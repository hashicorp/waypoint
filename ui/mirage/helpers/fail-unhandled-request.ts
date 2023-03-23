/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import * as QUnit from 'qunit';

/**
 * Reports an unhandled request to QUnit, to surface missing handlers more
 * clearly.
 */
export default function failUnhandledRequest(verb: string, path: string): void {
  let result = false;
  let message = `There is no Mirage handler for ${verb} ${path}. Please define one in ui/mirage/config.ts.`;
  // Technically it is possible to get the real stack using `new
  // Error().stack` but honestly it’s pretty opaque and not terribly useful
  // for debugging. Easier to bring folks here so they can add a breakpoint
  // and dig around.
  let source = 'ui/mirage/helpers/fail-unhandled-request.ts:16';

  QUnit.config.current.assert.pushResult({ result, message, source });
}
