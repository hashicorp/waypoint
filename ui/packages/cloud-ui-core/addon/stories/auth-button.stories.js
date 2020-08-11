import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'AuthButton',
  component: 'AuthButton',
};

// add stories by adding more exported functions
export let AuthButton = () => ({
  template: hbs`<AuthButton />`,
  context: {
    // add items to the component rendering context here
  },
});
