import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, focus } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | status-badge', function (hooks) {
  setupRenderingTest(hooks);

  test('renders the correct class for @state', async function (assert) {
    await render(hbs`
      <StatusBadge @state={{this.state}} />
    `);

    assert.dom('.badge-status--unknown').exists();

    this.set('state', 'UNKNOWN');
    assert.dom('.badge-status--unknown').exists();

    this.set('state', 'ALIVE');
    assert.dom('.badge-status--alive').exists();

    this.set('state', 'READY');
    assert.dom('.badge-status--ready').exists();

    this.set('state', 'DOWN');
    assert.dom('.badge-status--down').exists();

    this.set('state', 'PARTIAL');
    assert.dom('.badge-status--partial').exists();
  });

  test('renders the correct icon for @state', async function (assert) {
    await render(hbs`
      <StatusBadge @state={{this.state}} />
    `);

    assert.dom('[data-test-icon-type="help-circle-outline"]').exists();

    this.set('state', 'UNKNOWN');
    assert.dom('[data-test-icon-type="help-circle-outline"]').exists();

    this.set('state', 'ALIVE');
    assert.dom('[data-test-icon-type="run"]').exists();

    this.set('state', 'READY');
    assert.dom('[data-test-icon-type="check-plain"]').exists();

    this.set('state', 'DOWN');
    assert.dom('[data-test-icon-type="cancel-circle-fill"]').exists();

    this.set('state', 'PARTIAL');
    assert.dom('[data-test-icon-type="alert-triangle"]').exists();
  });

  test('renders the correct text for @state', async function (assert) {
    await render(hbs`
      <StatusBadge @state={{this.state}} />
    `);

    assert.dom('[data-test-status-badge]').includesText('Unknown');

    this.set('state', 'UNKNOWN');
    assert.dom('[data-test-status-badge]').includesText('Unknown');

    this.set('state', 'ALIVE');
    assert.dom('[data-test-status-badge]').includesText('Starting…');

    this.set('state', 'READY');
    assert.dom('[data-test-status-badge]').includesText('Up');

    this.set('state', 'DOWN');
    assert.dom('[data-test-status-badge]').includesText('Down');

    this.set('state', 'PARTIAL');
    assert.dom('[data-test-status-badge]').includesText('Partial');
  });

  test('does not render text if @iconOnly={{true}}', async function (assert) {
    await render(hbs`
      <StatusBadge @state="READY" @iconOnly={{true}} />
    `);

    assert.dom('[data-test-status-badge]').doesNotIncludeText('Up');
  });

  test('renders a default tooltip for @state', async function (assert) {
    await render(hbs`
      <StatusBadge @state="READY" />
    `);

    await focus('[data-test-status-badge]');

    assert.dom('.ember-tooltip').includesText('Application is ready');
  });

  test('renders @message as a tooltip', async function (assert) {
    await render(hbs`
      <StatusBadge
        @state={{this.state}}
        @message="Test message"
      />
    `);

    await focus('[data-test-status-badge]');

    assert.dom('.ember-tooltip').includesText('Test message');
  });
});
