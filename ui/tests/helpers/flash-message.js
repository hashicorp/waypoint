/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import FlashObject from 'ember-cli-flash/flash/object';

// This prevents ember-cli-flash from starting timers and other behavior that
// gums up our tests.
//
// eslint-disable-next-line @typescript-eslint/no-empty-function
FlashObject.reopen({ init() {} });
