import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Layout',
  component: 'Layout',
};

// add stories by adding more exported functions
export let Layout = () => ({
  template: hbs`
    <Layout>
      <Layout::Header>
        Header
      </Layout::Header>
      <Layout::Sidebar>
        Sidebar
      </Layout::Sidebar>
      <Layout::Content>
        Content
      </Layout::Content>
      <Layout::Drawer>
        Drawer
      </Layout::Drawer>
      <Layout::Footer>
        Footer
      </Layout::Footer>
    </Layout>
  `,
  context: {
    // add items to the component rendering context here
  },
});

export let LayoutWithoutDrawer = () => ({
  template: hbs`
    <Layout>
      <Layout::Header>
        Header
      </Layout::Header>
      <Layout::Sidebar>
        Sidebar
      </Layout::Sidebar>
      <Layout::Content>
        Content
      </Layout::Content>
      <Layout::Footer>
        Footer
      </Layout::Footer>
    </Layout>
  `,
  context: {
    // add items to the component rendering context here
  },
});

export let LayoutWithoutSidebar = () => ({
  template: hbs`
    <Layout>
      <Layout::Header>
        Header
      </Layout::Header>
      <Layout::Content>
        Content
      </Layout::Content>
      <Layout::Drawer>
        Drawer
      </Layout::Drawer>
      <Layout::Footer>
        Footer
      </Layout::Footer>
    </Layout>
  `,
  context: {
    // add items to the component rendering context here
  },
});

export let LayoutWithoutSidebarAndDrawer = () => ({
  template: hbs`
    <Layout>
      <Layout::Header>
        Header
      </Layout::Header>
      <Layout::Content>
        Content
      </Layout::Content>
      <Layout::Footer>
        Footer
      </Layout::Footer>
    </Layout>
  `,
  context: {
    // add items to the component rendering context here
  },
});
