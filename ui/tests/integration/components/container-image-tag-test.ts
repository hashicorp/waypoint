/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | container-image-tag', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    this.set('statusReport', {
      resourcesList: [
        {
          type: 'container',
          stateJson: '{"Config": {"Image": "docker:tag"}}',
        },
      ],
    });

    await render(hbs`<ContainerImageTag @statusReport={{this.statusReport}}/>`);

    assert.dom('[data-test-image-ref-uri]').hasText('docker');
    assert.dom('[data-test-image-ref-tag]').hasText('tag');
  });

  test('it renders n/a when a status report does not exist', async function (assert) {
    this.set('statusReport2', {});

    await render(hbs`<ContainerImageTag @statusReport={{this.statusReport2}}/>`);

    assert.dom(this.element).hasText('n/a');
  });

  test('it renders multiple tags', async function (assert) {
    this.set('multiStatus', {
      resourcesList: [
        {
          stateJson: '{ "Config": { "Image": "docker:tag" }, "Pod": { "Image": "kubernetes:latest" } }',
        },
      ],
    });

    await render(hbs`<ContainerImageTag @statusReport={{this.multiStatus}}/>`);
    assert.dom('[data-test-image-ref]').exists({ count: 2 });
  });

  test('it handles refs like "localhost:5000/image-name:latest"', async function (assert) {
    this.set('statusReport', {
      resourcesList: [
        {
          type: 'container',
          stateJson: '{"Config": {"Image": "localhost:5000/image-name:latest"}}',
        },
      ],
    });

    await render(hbs`<ContainerImageTag @statusReport={{this.statusReport}}/>`);

    assert.dom('[data-test-image-ref-uri]').hasText('localhost:5000/image-name');
    assert.dom('[data-test-image-ref-tag]').hasText('latest');
  });

  test('it handles refs like "localhost:5000/image-name"', async function (assert) {
    this.set('statusReport', {
      resourcesList: [
        {
          type: 'container',
          stateJson: '{"Config": {"Image": "localhost:5000/image-name"}}',
        },
      ],
    });

    await render(hbs`<ContainerImageTag @statusReport={{this.statusReport}}/>`);

    assert.dom('[data-test-image-ref-uri]').hasText('localhost:5000/image-name');
    assert.dom('[data-test-image-ref-tag]').doesNotExist();
  });

  test('it handles refs with digests', async function (assert) {
    this.set('statusReport', {
      resourcesList: [
        {
          type: 'container',
          stateJson:
            '{"Config": {"Image": "localhost:5000/image-name@sha256:aaaaf56b44807c64d294e6c8059b479f35350b454492398225034174808d1726"}}',
        },
      ],
    });

    await render(hbs`<ContainerImageTag @statusReport={{this.statusReport}}/>`);

    assert.dom('[data-test-image-ref-uri]').hasText('localhost:5000/image-name');
    assert.dom('[data-test-image-ref-tag]').hasText('sha256:aaaaf56');
  });
});
