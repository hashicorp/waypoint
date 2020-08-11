import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import menuPageObject, {
  MENU_SELECTOR,
  MENU_TRIGGER_SELECTOR,
  MENU_CONTENT_SELECTOR,
} from 'cloud-ui-core/test-support/pages/components/menu';
import sinon from 'sinon';

const menu = create(menuPageObject);

module('Integration | Component | Menu', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
    <Menu as |M|>
      <M.Trigger>
        Open this...
      </M.Trigger>
      <M.Content>
        Put some content here
      </M.Content>
    </Menu>
      `);
    assert.notOk(menu.isOpen, 'renders closed by default');
    assert.dom(MENU_TRIGGER_SELECTOR).containsText('Open this...', 'renders to the trigger element');
    assert.dom(MENU_CONTENT_SELECTOR).containsText('Put some content here', 'renders to the content element');
  });

  test('it renders: @isOpen', async function(assert) {
    await render(hbs`
    <Menu @isOpen={{true}} as |M|>
      <M.Trigger>
        Open this...
      </M.Trigger>
      <M.Content>
        Put some content here
      </M.Content>
    </Menu>
      `);
    assert.ok(menu.isOpen, 'renders open');
  });

  test('it changes open attribute on click, clickOutside', async function(assert) {
    await render(hbs`
    <Menu as |M|>
      <M.Trigger>
        Open this...
      </M.Trigger>
      <M.Content>
        Put some content here
      </M.Content>
    </Menu>
      `);

    await menu.click();
    assert.ok(menu.isOpen);

    await click(this.element);
    assert.notOk(menu.isOpen);
  });

  test('it calls onToggle', async function(assert) {
    let onToggle = sinon.spy();
    this.onToggle = onToggle;
    await render(hbs`
    <Menu @onToggle={{onToggle}} as |M|>
      <M.Trigger>
        Open this...
      </M.Trigger>
      <M.Content>
        Put some content here
      </M.Content>
    </Menu>
      `);

    await menu.click();
    assert.ok(onToggle.calledOnce);
    assert.ok(onToggle.calledWith(true));
  });

  test('it forwards attributes', async function(assert) {
    await render(hbs`
    <Menu class="menu-test" as |M|>
      <M.Trigger class="menu-test-trigger" aria-label="The real label">
        Open this...
      </M.Trigger>
      <M.Content class="menu-test-content">
        Put some content here
      </M.Content>
    </Menu>
      `);

    assert.dom(MENU_SELECTOR).hasClass('menu-test');
    assert.dom(MENU_TRIGGER_SELECTOR).hasClass('menu-test-trigger');
    assert.dom(MENU_TRIGGER_SELECTOR).hasAttribute('aria-label', 'The real label');
    assert.dom(MENU_CONTENT_SELECTOR).hasClass('menu-test-content');
  });
});
