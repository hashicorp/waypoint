import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import externalLinkPageObject, {
  EXTERNAL_LINK_CONTAINER_SELECTOR,
} from 'cloud-ui-core/test-support/pages/components/external-link';

const externalLink = create(externalLinkPageObject);

module('Integration | Component | ExternalLink', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <ExternalLink href="https://www.google.com">
        template block text
      </ExternalLink>
    `);

    assert.equal(externalLink.text, 'template block text');
    assert.dom(EXTERNAL_LINK_CONTAINER_SELECTOR).hasAttribute('href', 'https://www.google.com');
    assert.dom(EXTERNAL_LINK_CONTAINER_SELECTOR).hasAttribute('target', '_blank');
    assert.dom(EXTERNAL_LINK_CONTAINER_SELECTOR).hasAttribute('rel', 'noopener noreferrer');
  });
});
