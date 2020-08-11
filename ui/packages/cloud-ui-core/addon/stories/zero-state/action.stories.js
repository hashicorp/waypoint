import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'ZeroState/Action',
  component: 'ZeroStateAction',
};

export let ZeroStateAction = () => ({
  template: hbs`
    <ZeroState::Action>
      <button type='submit'>
        {{t 'components.page.hvns.create.title'}}
      </button>
    </ZeroState::Action>
  `,
  context: {},
});
