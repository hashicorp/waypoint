import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { getUnixTime, subMinutes } from 'date-fns';
import { a11yAudit } from 'ember-a11y-testing/test-support';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';

module('Integration | Component | app-item/build', function (hooks) {
  setupRenderingTest(hooks);

  test('with a build and a push', async function (assert) {
    this.set('build', {
      sequence: 3,
      status: {
        state: 2,
        startTime: minutesAgo(3),
        completeTime: minutesAgo(2),
      },
      component: {
        type: 1,
        name: 'docker',
      },
      pushedArtifact: {
        component: {
          type: 2,
          name: 'docker',
        },
        status: {
          state: 2,
          startTime: minutesAgo(2),
          completeTime: minutesAgo(1),
        },
      },
    });

    await render(hbs`
      <ul>
        <AppItem::Build @build={{this.build}} />
      </ul>
    `);
    await a11yAudit();

    assert.dom('[data-test-app-item-build]').includesText('v3');
    assert.dom('[data-test-icon-type="logo-docker-color"]').exists();
    assert.dom('[data-test-app-item-build]').includesText('Pushed to Docker');
    assert.dom('[data-test-operation-status-indicator="success"]').exists();
    assert.dom('[data-test-app-item-build]').includesText('1 minute ago');
    assert.dom('[data-test-app-item-build]').includesText('Built in 1 minute');
  });

  test('with no push', async function (assert) {
    this.set('build', {
      sequence: 3,
      status: {
        state: 2,
        startTime: minutesAgo(3),
        completeTime: minutesAgo(2),
      },
      component: {
        type: 1,
        name: 'docker',
      },
      pushedArtifact: null,
    });

    await render(hbs`
      <ul>
        <AppItem::Build @build={{this.build}} />
      </ul>
    `);
    await a11yAudit();

    assert.dom('[data-test-app-item-build]').includesText('v3');
    assert.dom('[data-test-icon-type="logo-docker-color"]').exists();
    assert.dom('[data-test-app-item-build]').includesText('Built with Docker');
    assert.dom('[data-test-operation-status-indicator="success"]').exists();
    assert.dom('[data-test-app-item-build]').includesText('2 minutes ago');
  });
});

function minutesAgo(n: number): Timestamp.AsObject {
  let now = new Date();
  let date = subMinutes(now, n);
  let result = {
    seconds: getUnixTime(date),
    nanos: 0,
  };

  return result;
}
