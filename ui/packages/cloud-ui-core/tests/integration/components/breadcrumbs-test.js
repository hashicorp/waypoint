import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import breadcrumbs from 'cloud-ui-core/test-support/pages/components/breadcrumbs';
import breadcrumbsCrumb from 'cloud-ui-core/test-support/pages/components/breadcrumbs/crumb';

module('Integration | Component | Breadcrumbs', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Breadcrumbs>
        <Breadcrumbs::Crumb
          @route="cloud.orgs.detail.index"
        >
          Organizations
        </Breadcrumbs::Crumb>
        <Breadcrumbs::Crumb
          @route="cloud.orgs.detail.projects.detail.index"
        >
          Project
        </Breadcrumbs::Crumb>
        <Breadcrumbs::Crumb
          @route="cloud.orgs.detail.projects.detail.hvns.list"
        >
          Resource
        </Breadcrumbs::Crumb>
      </Breadcrumbs>
    `);
    assert.dom(breadcrumbs.navSelector).exists('nav exists');
    assert.dom(breadcrumbsCrumb.crumbSelector).exists({ count: 3 });
    assert.dom(breadcrumbsCrumb.separatorSelector).exists({ count: 3 });
    assert.dom(breadcrumbsCrumb.separatorSelector).hasAttribute('aria-hidden', 'true');
    assert.dom(breadcrumbsCrumb.separatorSelector).hasText('/');
  });
});
