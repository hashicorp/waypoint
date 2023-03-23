/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Modifier | code-mirror', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    let valueToChange = 'initial value in editor';
    this.set('value', valueToChange);
    this.set('onInput', () => null);

    await render(
      hbs`
      <div
        {{code-mirror
          value=this.value
          onInput=this.onInput
          options=(hash screenReaderLabel="test")
        }}
      />
    `
    );

    assert.dom('.CodeMirror').exists();
    assert.dom('.CodeMirror').containsText(valueToChange);
  });

  test('it renders even with undefined value/onInput', async function (assert) {
    this.set('value', undefined);
    this.set('onInput', undefined);

    await render(hbs`
      <div
        {{code-mirror
          value=this.value
          onInput=this.onInput
          options=(hash screenReaderLabel="test")
        }}
      />
    `);
    assert.dom('.CodeMirror').exists();
  });

  test('it calls onInput when new text added', async function (assert) {
    this.set('value', '');
    this.set('onInput', (value: string) => this.set('value', value));
    this.set('options', {
      lineNumbers: false,
      screenReaderLabel: 'test',
    }); // otherwise the #s appear when comparing text

    await render(hbs`
      <div
        {{code-mirror
          value=this.value
          onInput=this.onInput
          options=this.options
        }}
      />
    `);

    let textArea = this.element.querySelector('textarea') as HTMLElement;
    // if set as value on initial render, it won't get deleted on the second fillIn call
    let firstValue = 'first value';
    await fillIn(textArea, firstValue);

    assert.dom('.CodeMirror-code').matchesText(firstValue);

    let newValue = 'second value';
    await fillIn(textArea, newValue);

    /* eslint-disable ember/no-get */
    assert.equal(this.get('value'), newValue);
    assert.dom('.CodeMirror-code').doesNotContainText(firstValue);
    assert.dom('.CodeMirror-code').matchesText(newValue);
  });

  test('it renders user-specified options', async function (assert) {
    this.set('value', '');
    this.set('onInput', () => null);
    this.set('options', {
      lineNumbers: false,
      theme: 'default',
      screenReaderLabel: 'test',
    });

    await render(hbs`
      <div
        {{code-mirror
          value=this.value
          onInput=this.onInput
          options=(hash
            screenReaderLabel="test"
          )
        }}
      />
    `); // without options
    // default options
    assert.dom('.cm-s-monokai').exists();
    assert.dom('.CodeMirror-code').containsText('1');

    await render(hbs`
      <div
        {{code-mirror
          value=this.value
          onInput=this.onInput
          options=this.options
        }}
      />
    `); // with options
    // codemirror's real default theme should be set and linenumbers gone
    assert.dom('.cm-s-monokai').doesNotExist();
    assert.dom('.CodeMirror-code').doesNotContainText('1');
  });
});
