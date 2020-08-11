import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import definitionListPageObject, {
  containerSelector,
  keySelector,
  valueSelector,
} from 'cloud-ui-core/test-support/pages/components/definition-list';

let definitionList = create(definitionListPageObject);

module('Integration | Component | definition-list', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <DefinitionList as |DL|>
        <DL.Key>key1</DL.Key>
        <DL.Value>value1</DL.Value>
        <DL.Key>key2</DL.Key>
        <DL.Value>value2</DL.Value>
      </DefinitionList>
    `);

    assert.dom(containerSelector).exists();
    assert.dom(keySelector).exists({ count: 2 });
    assert.dom(valueSelector).exists({ count: 2 });
    assert.equal(definitionList.keys[0].text, 'key1');
    assert.equal(definitionList.keys[1].text, 'key2');
    assert.equal(definitionList.values[0].text, 'value1');
    assert.equal(definitionList.values[1].text, 'value2');
  });
});
