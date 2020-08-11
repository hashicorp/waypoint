import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select } from '@storybook/addon-knobs';
import { DEFAULT_VARIANT, VARIANT_SCALE } from 'cloud-ui/components/layout/page-header/consts';

export default {
  title: 'Layout::PageHeader',
  component: 'LayoutPageHeader',
  decorators: [withKnobs],
};

// add stories by adding more exported functions
export let PageHeaderCreateOrEdit = () => ({
  template: hbs`
    <Layout::PageHeader @variant={{variant}} as |P|>
      <P.Breadcrumbs>
        <RouterBreadcrumbs />
      </P.Breadcrumbs>
      <P.Title>
        <Typography @variant='h1'>
        {{t 'components.page.hvns.create.title'}}
        </Typography>
        <Typography @variant='body1'>
          {{t 'components.page.hvns.create.description'}}
        </Typography>
      </P.Title>
    </Layout::PageHeader>
  `,
  context: {
    variant: select('Variant', VARIANT_SCALE, DEFAULT_VARIANT),
  },
});
export let PageHeaderList = () => ({
  template: hbs`
    <Layout::PageHeader @variant={{variant}} as |P|>
      <P.Breadcrumbs>
        <RouterBreadcrumbs />
      </P.Breadcrumbs>
      <P.Title>
        <Typography @variant='h1'>
          {{t 'components.page.hvns.list.title'}}
        </Typography>
      </P.Title>
      <P.Actions>
        <LinkTo
          class='button button--primary'
          @route='cloud.orgs.detail.projects.detail.hvns.create'
          data-test-network-create
        >
          <Icon @type='plus-plain' @size='md' />
          {{t 'components.page.hvns.create.action'}}
        </LinkTo>
      </P.Actions>
    </Layout::PageHeader>
  `,
  context: {
    variant: select('Variant', VARIANT_SCALE, VARIANT_SCALE[1]),
  },
});

export let PageHeaderDetail = () => ({
  template: hbs`
    <Layout::PageHeader @variant={{variant}} as |P|>
      <P.Breadcrumbs>
        <RouterBreadcrumbs />
      </P.Breadcrumbs>
      <P.Title>
        <Typography @variant='h1'>
          some-hvnnetwork-name
        </Typography>
      </P.Title>
      <P.Actions>
        <button type='button' data-test-delete-button>
          {{#if this.areYouSure}}
            Really delete this?
          {{else}}
            Delete this?
          {{/if}}
        </button>
      </P.Actions>
      <P.Tabs>
        <Tabs>
          <LinkTo @route='cloud.orgs.detail.projects.detail.hvns.detail'>
            {{t 'components.page.hvns.detail.tabs.overview'}}
          </LinkTo>
          <a href='/peering-connections'>
            {{t 'components.page.hvns.detail.tabs.peering-connections'}}
          </a>
        </Tabs>
      </P.Tabs>
    </Layout::PageHeader>
  `,
  context: {
    variant: select('Variant', VARIANT_SCALE, VARIANT_SCALE[1]),
  },
});
