import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'ConfigDownloadButton',
  component: 'ConfigDownloadButton',
};

// add stories by adding more exported functions
export let ConfigDownloadButton = () => ({
  template: hbs`
  <ConfigDownloadButton @resource={{ resource }} @resourceType='consul'>
    foo bar
  </ConfigDownloadButton>
  `,
  context: {
    resource: {
      id: "aaa-bbb-cccc",
      location: {
        organizationId: "aaa",
        projectId: "bbb"
      }
    }
  }
});
