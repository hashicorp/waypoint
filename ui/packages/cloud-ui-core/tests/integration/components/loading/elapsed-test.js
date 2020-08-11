import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import sinon from 'sinon';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import loadingElapsed from 'cloud-ui-core/test-support/pages/components/loading/elapsed';

module('Integration | Component | Loading/Elapsed', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.clock = sinon.useFakeTimers({
      now: Date.now(),
      shouldAdvanceTime: true,
    });
  });

  hooks.afterEach(function() {
    this.clock.restore();
  });

  test('it renders', async function(assert) {
    assert.expect(4);
    await render(hbs`<Loading::Elapsed />`);

    assert.dom(loadingElapsed.containerSelector).exists();
    assert.dom(loadingElapsed.containerSelector).containsText('Time Elapsed: --:--');
    assert.dom(loadingElapsed.containerSelector).hasClass('loadingState__elapsed');
    this.clock.tick(60 * 1000);
    assert.dom(loadingElapsed.containerSelector).containsText('Time Elapsed: 01:00');
  });

  test('it renders: @startTime', async function(assert) {
    let thirtyMinutesAgo = new Date().getTime() - 1000 * 60 * 30;
    this.set('startTime', thirtyMinutesAgo);
    await render(hbs`<Loading::Elapsed @startTime={{startTime}} />`);
    this.clock.tick(1000);
    assert.dom(loadingElapsed.containerSelector).containsText('Time Elapsed: 30:01');
  });
});
