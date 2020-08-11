import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'ZeroState',
  component: 'ZeroState',

  subcomponents: {
    'ZS.Header': 'ZeroStateHeader',
    'ZS.Message': 'ZeroStateMessage',
  },
};

export let ZeroState = () => ({
  template: hbs`
    <ZeroState as |ZS|>
      <ZS.Header>
        {{t 'components.page.hvns.list.empty.header'}}
      </ZS.Header>
      <ZS.Message>
        {{t 'components.page.hvns.list.empty.message'}}
      </ZS.Message>
      <ZS.Action>
        <button type='submit'>
          {{t 'components.page.hvns.create.title'}}
        </button>
      </ZS.Action>
    </ZeroState>
  `,
  context: {},
});
