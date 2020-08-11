import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import ctaLink from 'cloud-ui-core/test-support/pages/components/cta-link';

module('Integration | Component | CtaLink', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<CtaLink @route='cloud.orgs'>Create Network</CtaLink>`);
    assert.ok(ctaLink.isPresent);
  });
});
