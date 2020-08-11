import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import errorPageObject from 'cloud-ui-core/test-support/pages/components/error';

const error = create(errorPageObject);

module('Integration | Component | error', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Error>
        <:title>Not Found</:title>
        <:subtitle>Error 404</:subtitle>
        <:content>Some content message</:content>
        <:footer>
          <LinkTo @route="cloud">
            <Icon @type='chevron-left' @size='sm' aria-hidden='true' />
            Go back
          </LinkTo>
        </:footer>
      </Error>
    `);

    assert.ok(error.showsContainer, 'shows container');
    assert.ok(error.showsIcon, 'shows icon');
    assert.ok(error.showsTitle, 'shows title');
    assert.equal(error.titleText, 'Not Found', 'title text displays');
    assert.ok(error.showsSubtitle, 'shows subtitle');
    assert.equal(error.subtitleText, 'Error 404', 'title text displays');
    assert.ok(error.showsContent, 'shows content');
    assert.equal(error.contentText, 'Some content message', 'title text displays');
    assert.ok(error.showsFooter, 'shows footer');
  });
});
