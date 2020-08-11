import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import waitForError from 'dummy/tests/helpers/wait-for-error';
import { create } from 'ember-cli-page-object';
import iconBadgePageObject from 'cloud-ui-core/test-support/pages/components/icon-badge';
import { SOURCE_SCALE, VARIANT_SCALE } from 'dummy/helpers/option-for-icon-badge';

const iconBadge = create(iconBadgePageObject);

module('Integration | Component | IconBadge', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders: all @variant sizes', async function(assert) {
    assert.expect(VARIANT_SCALE.length * SOURCE_SCALE.length * 3);
    for (let source of SOURCE_SCALE) {
      this.set('source', source);
      for (let variant of VARIANT_SCALE) {
        this.set('variant', variant);
        await render(hbs`<IconBadge @source={{source}} @variant={{variant}} />`);
        assert.ok(iconBadge.showsContainer, `the ${source} ${variant} container renders`);
        assert.ok(iconBadge.showsIcon, `the ${source} ${variant} icon renders`);
        assert.ok(iconBadge.showsLabel, `the ${source} ${variant} label renders`);
      }
    }
  });

  test('it renders: throws invalid source', async function(assert) {
    let promise = waitForError();
    render(hbs`<IconBadge @source="MISSING" @variant="PENDING" />`);
    let err = await promise;
    assert.ok(err.message.includes('@source must '), "errors when passed a source that's not allowed");
  });
});
