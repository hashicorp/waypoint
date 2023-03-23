/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';

class RouterStub extends Service {
  currentRouteName = 'workspace.projects.project.app.release';
}

module('Integration | Component | timeline', function (hooks) {
  setupRenderingTest(hooks);

  test('happy path', async function (assert) {
    this.owner.register('service:router', RouterStub);
    this.set('timeline', {
      build: {
        sequence: 3,
        status: {
          state: 2,
          completeTime: {
            seconds: 1639166870,
          },
        },
      },
      deployment: {
        sequence: 2,
        status: {
          state: 4,
          completeTime: {
            seconds: 1639166879,
          },
        },
      },
      release: {
        sequence: 1,
        status: {
          state: 4,
          completeTime: {
            seconds: 1639166880,
          },
        },
      },
    });
    await render(hbs`<Timeline @model={{this.timeline}}/>`);

    let listItems = this.element.querySelectorAll('li');
    assert.equal(listItems.length, 3);
    assert.dom(listItems[0]).containsText('Build');
    assert.dom(listItems[1]).containsText('Deployment');
    assert.dom(listItems[2]).containsText('Release');
    assert.dom(listItems[2]).containsText('You are here');
  });

  test('it does not render timestamps if status unavailable', async function (assert) {
    this.owner.register('service:router', RouterStub);
    this.set('timeline', {
      build: { sequence: 4 },
      deployment: { sequence: 3 },
      release: { sequence: 2 },
    });
    await render(hbs`<Timeline @model={{this.timeline}}/>`);

    let listItems = this.element.querySelectorAll('li');
    assert.equal(listItems.length, 3);

    for (let item of listItems) {
      assert.dom('.timeline__li__timestamp', item).doesNotExist();
    }
  });

  test('it does not render you are here badge if no router.currentRouteName value', async function (assert) {
    this.owner.register('service:router', RouterStub);
    this.set('timeline', {
      build: { sequence: 4 },
      deployment: { sequence: 3 },
    });
    await render(hbs`<Timeline @model={{this.timeline}}/>`);

    let listItems = this.element.querySelectorAll('li');
    assert.equal(listItems.length, 2);
    assert.dom('.timeline__current-marker', this.element).doesNotExist();
  });

  test('it renders an empty ul with no model', async function (assert) {
    await render(hbs`<Timeline />`);
    assert.dom(this.element).hasNoText();
    assert.dom('ul').exists();
    assert.dom('li').doesNotExist();
  });
});
