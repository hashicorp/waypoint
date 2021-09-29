import { module, test } from 'qunit';

import { createTerminal } from 'waypoint/utils/create-terminal';
import hbs from 'htmlbars-inline-precompile';
import { later } from '@ember/runloop';
import { render } from '@ember/test-helpers';
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
    }, 300);
    // Clean up terminal
    terminal.dispose();
  });
});
