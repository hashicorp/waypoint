import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Breadcrumbs',
  component: 'Breadcrumbs',
};

// add stories by adding more exported functions
export let Breadcrumbs = () => ({
  template: hbs`
    <Breadcrumbs>
      <Breadcrumbs::Crumb
        @route="cloud.orgs.detail"
      >
        Organizations
      </Breadcrumbs::Crumb>
      <Breadcrumbs::Crumb
        @route="cloud.orgs.detail.projects.detail"
      >
        Project
      </Breadcrumbs::Crumb>
      <Breadcrumbs::Crumb
        @route="cloud.orgs.detail.projects.detail.hvns.detail"
      >
        Resource
      </Breadcrumbs::Crumb>
    </Breadcrumbs>
    `,
  context: {},
});
