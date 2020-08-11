import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'ZeroState/Header',
  component: 'ZeroStateHeader',
};

export let ZeroStateHeader = () => ({
  template: hbs`
    <ZeroState::Header>
      {{t 'components.page.hvns.list.empty.header'}}
    </ZeroState::Header>
  `,
  context: {},
});
