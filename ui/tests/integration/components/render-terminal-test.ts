/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-disable qunit/require-expect */
import { clearRender, render } from '@ember/test-helpers';
import { module, test } from 'qunit';

import { createTerminal } from 'waypoint/utils/create-terminal';
import hbs from 'htmlbars-inline-precompile';
import { later } from '@ember/runloop';
import settled from '@ember/test-helpers/settled';
import { setupRenderingTest } from 'ember-qunit';

module('Integration | Component | render-terminal', function (hooks) {
  setupRenderingTest(hooks);

  test('basic rendering', async function (assert) {
    // Setup terminal externally, as expected by the component
    let terminal = createTerminal({ inputDisabled: true });
    this.set('terminal', terminal);
    // pass terminal and render
    await render(hbs`<RenderTerminal @terminal={{this.terminal}}/>`);
    // write line and see if it renders
    terminal.writeln('Welcome to Waypoint!');
    // We have to use the runloop as writeln isn't async
    // Note that even the xterm.js rendering tests use polling to evaluate rendering
    later(() => {
      assert.equal(terminal?.buffer?.active?.getLine(0)?.translateToString(true), 'Welcome to Waypoint!');
    }, 10);
    assert.dom('[data-test-xterm-pane]').exists('the xterm pane renders');
    assert.dom('.xterm').exists('the xterm terminal renders');
    // Clean up terminal
    await settled();
    await clearRender();
    later(() => {
      assert.dom('.xterm').doesNotExist('the xterm terminal is destroyed');
      assert.dom(terminal.element?.assignedSlot).doesNotExist('the xterm terminal is destroyed');
    }, 10);
    await settled();
  });
});
