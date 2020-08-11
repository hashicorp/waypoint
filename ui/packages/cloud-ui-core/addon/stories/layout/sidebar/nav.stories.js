import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'LayoutSidebarNav',
  component: 'LayoutSidebarNav',
};

export let LayoutSidebarNav = () => ({
  template: hbs`
    <Layout::Sidebar as |S|>
      <S.Nav aria-labelledby='navLabel' as |SN|>
        <SN.Header id='projectNavLabel'>
          Project
        </SN.Header>
        <LinkTo @route='cloud.orgs.detail.projects.detail.index'>
          Overview
        </LinkTo>
        <SN.Section>
          <SN.Subheader id='navLabel'>
            Menu Subheader
          </SN.Subheader>
          <LinkTo @route='cloud.orgs.detail.projects'>
            Projects
          </LinkTo>
          <LinkTo @route='cloud.orgs.detail.projects.detail.hvns.list'>
            Link
          </LinkTo>
        </SN.Section>
      </S.Nav>
    </Layout::Sidebar>
  `,
});
