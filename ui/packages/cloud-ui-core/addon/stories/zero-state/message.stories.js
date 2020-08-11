import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'ZeroState/Message',
  component: 'ZeroStateMessage',
};

export let ZeroStateMessage = () => ({
  template: hbs`
    <ZeroState::Message>
      {{t 'components.page.hvns.list.empty.message'}}
    </ZeroState::Message>
  `,
  context: {},
});
