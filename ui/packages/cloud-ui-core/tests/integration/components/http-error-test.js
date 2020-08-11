import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupIntl, t } from 'ember-intl/test-support';
import { create } from 'ember-cli-page-object';
import errorPageObject, { ERROR_ICON_SELECTOR } from 'cloud-ui-core/test-support/pages/components/error';
import httpErrorPageObject from 'cloud-ui-core/test-support/pages/components/http-error';
import { ERROR_CODE_MAPPING, ERROR_CODE_SCALE } from 'dummy/helpers/option-for-http-error';

const error = create(errorPageObject);
const httpError = create(httpErrorPageObject);

module('Integration | Component | http-error', function(hooks) {
  setupRenderingTest(hooks);
  setupIntl(hooks);

  test('it renders', async function(assert) {
    await render(hbs`<HttpError />`);
    assert.ok(httpError.showsContainer, `the container renders`);
  });

  test('it renders: all @code options', async function(assert) {
    assert.expect(ERROR_CODE_SCALE.length * 4);

    for (let code of ERROR_CODE_SCALE) {
      this.set('code', code);
      await render(hbs`<HttpError @code={{code}} />`);
      assert.dom(ERROR_ICON_SELECTOR).hasClass('icon');
      assert.equal(error.titleText, t(ERROR_CODE_MAPPING[code].label), 'title text displays');
      assert.equal(error.subtitleText, `Error ${code}`, 'subtitle text displays');
      assert.equal(error.contentText, t(ERROR_CODE_MAPPING[code].message), 'content text displays');
    }
  });

  test('it renders: with override args', async function(assert) {
    this.code = 404;
    this.title = 'Some Title';
    this.message = 'Some Message';
    this.previousRoute = 'cloud.orgs';

    await render(
      hbs`<HttpError @code={{code}} @title={{title}} @message={{message}} @previousRoute={{previousRoute}} />`
    );

    assert.equal(error.titleText, this.title, 'title text displays');
    assert.equal(error.subtitleText, `Error ${this.code}`, 'subtitle text displays');
    assert.equal(error.contentText, this.message, 'content text displays');
  });
});
