/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { triggerEvent, triggerKeyEvent, focus } from '@ember/test-helpers';

// At the time of writing, xterm.js only implements the Macintosh “select all”
// keyboard shortcut (cmd+a), not the Linux and Windows equivalent (ctrl+a).
// What follows is a dreadful hack to make `getTerminalText` work in Linux CI.
// Note: this needs to happen before xterm is required, so you’ll find an
// import in `tests/test-helper.js`.
Object.defineProperty(navigator, 'platform', {
  value: 'Macintosh',
  configurable: true,
  enumerable: true,
  writable: false,
});

/**
 * Returns the text from an xterm.js terminal element.
 *
 * Simulates select-all, then copy, then returns the resultant clipboard data.
 *
 * Assumes only one xterm instance is rendered.
 */
export async function getTerminalText(): Promise<string> {
  await focus('.xterm-helper-textarea');

  // Trigger cmd+a to select all
  await triggerKeyEvent('.xterm-helper-textarea', 'keydown', 65, { metaKey: true });

  let clipboardData = new DataTransfer();
  await triggerEvent('.xterm', 'copy', { clipboardData });

  return clipboardData.getData('text').trim();
}
