/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { ContextObject, RunOptions } from 'axe-core';
import {
  setupGlobalA11yHooks,
  setEnableA11yAudit,
  setRunOptions,
  DEFAULT_A11Y_TEST_HELPER_NAMES,
} from 'ember-a11y-testing/test-support';

const A11Y_TEST_HELPER_NAMES: typeof DEFAULT_A11Y_TEST_HELPER_NAMES = [
  ...DEFAULT_A11Y_TEST_HELPER_NAMES,
  'render',
  'focus',
];

// ember-a11y-testing allows us to pass `include` and `exclude` context
// parameters as run options. This isn’t documented, and isn’t represented in
// the type defintions but you can see it’s covered by the test suite here:
// https://github.com/ember-a11y/ember-a11y-testing/blob/v4.0.7/tests/acceptance/a11y-audit-test.ts
//
// This reassures TypeScript that `include` and `exclude` are acceptable.
type OptionsWithContext = RunOptions & ContextObject;

// Selectors of elements to exclude from a11y auditing. See the following docs
// for more:
// https://github.com/dequelabs/axe-core/blob/develop/doc/API.md#include-exclude-object
const include = [['#ember-testing-container']];
const exclude = [
  ['.pds-logomark'],
  ['.pds-tabNav'],
  ['.card-header'],
  ['.x-toggle-btn'],
  ['.flight-sprite-container'],
];

export function setup(): void {
  setupGlobalA11yHooks(() => true, { helpers: A11Y_TEST_HELPER_NAMES });
  setEnableA11yAudit(true);
  setRunOptions({ include, exclude } as OptionsWithContext);
}
