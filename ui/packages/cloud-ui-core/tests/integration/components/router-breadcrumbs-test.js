import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import Service from '@ember/service';
import { create } from 'ember-cli-page-object';
import { tracked } from '@glimmer/tracking';
import routerBreadcrumbsPageObject, {
  containerSelector,
  crumbSelector,
} from 'cloud-ui-core/test-support/pages/components/router-breadcrumbs';
import { makeRouteStub } from 'dummy/tests/helpers/make-route-stub';

const routerBreadcrumbs = create(routerBreadcrumbsPageObject);

class RouterStub extends Service {
  @tracked currentRouteName = 'cloud.orgs.detail.projects.detail.hvns.list';
}

module('Integration | Component | RouterBreadcrumbs', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.owner.unregister('router:main');
    this.owner.unregister('service:router');
    this.owner.register('service:router', RouterStub);
    this.router = this.owner.lookup('service:router');
  });

  test('it renders', async function(assert) {
    this.owner.register('route:test', makeRouteStub());
    this.owner.register('route:test.one', makeRouteStub());
    this.owner.register('route:test.one.two', makeRouteStub());
    this.owner.register('route:test.one.two.three', makeRouteStub());
    this.owner.register('route:test.one.two.three.four', makeRouteStub());
    this.router.currentRouteName = 'test.one.two.three.four';

    await render(hbs`<RouterBreadcrumbs />`);

    assert.dom(containerSelector).exists();
    assert.dom(crumbSelector).exists({ count: 5 });
    assert.equal(routerBreadcrumbs.crumbsSelector[0].text, 'Test');
    assert.equal(routerBreadcrumbs.crumbsSelector[1].text, 'One');
    assert.equal(routerBreadcrumbs.crumbsSelector[2].text, 'Two');
    assert.equal(routerBreadcrumbs.crumbsSelector[3].text, 'Three');
    assert.equal(routerBreadcrumbs.crumbsSelector[4].text, 'Four');
  });

  test('it renders: with @breadcrumb on routes', async function(assert) {
    this.owner.register(
      'route:test',
      makeRouteStub({
        breadcrumb: function() {
          return { title: 'Organizations' };
        },
      })
    );
    this.owner.register(
      'route:test.one',
      makeRouteStub({
        breadcrumb: function() {
          return { title: 'Organization Details' };
        },
      })
    );
    this.owner.register(
      'route:test.one.two',
      makeRouteStub({
        breadcrumb: function() {
          return { title: 'Projects' };
        },
      })
    );
    this.owner.register(
      'route:test.one.two.three',
      makeRouteStub({
        breadcrumb: function() {
          return { title: 'Project Detail' };
        },
      })
    );
    this.owner.register(
      'route:test.one.two.three.four',
      makeRouteStub({
        breadcrumb: function() {
          return { title: 'Networks' };
        },
      })
    );
    this.router.currentRouteName = 'test.one.two.three.four';

    await render(hbs`<RouterBreadcrumbs />`);

    assert.dom(containerSelector).exists();
    assert.dom(crumbSelector).exists({ count: 5 });
    assert.equal(routerBreadcrumbs.crumbsSelector[0].text, 'Organizations');
    assert.equal(routerBreadcrumbs.crumbsSelector[1].text, 'Organization Details');
    assert.equal(routerBreadcrumbs.crumbsSelector[2].text, 'Projects');
    assert.equal(routerBreadcrumbs.crumbsSelector[3].text, 'Project Detail');
    assert.equal(routerBreadcrumbs.crumbsSelector[4].text, 'Networks');
  });

  test('it renders: with @breadcrumb on routes: hide crumb', async function(assert) {
    this.owner.register(
      'route:test',
      makeRouteStub({
        breadcrumb: function() {
          return { title: 'Organizations' };
        },
      })
    );
    this.owner.register('route:test.one', makeRouteStub({ breadcrumb: null }));
    this.owner.register('route:test.one.two', makeRouteStub({ breadcrumb: null }));
    this.owner.register('route:test.one.two.three', makeRouteStub({ breadcrumb: null }));
    this.owner.register(
      'route:test.one.two.three.four',
      makeRouteStub({
        breadcrumb: function() {
          return { title: 'Networks' };
        },
      })
    );
    this.router.currentRouteName = 'test.one.two.three.four';

    await render(hbs`<RouterBreadcrumbs />`);

    assert.dom(containerSelector).exists();
    assert.dom(crumbSelector).exists({ count: 2 });
    assert.equal(routerBreadcrumbs.crumbsSelector[0].text, 'Organizations');
    assert.equal(routerBreadcrumbs.crumbsSelector[1].text, 'Networks');
  });
});
