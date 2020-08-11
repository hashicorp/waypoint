import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'TopNav',
  component: 'TopNav',
};

// add stories by adding more exported functions
export let TopNav = () => ({
  template: hbs`<TopNav />`,
  context: {
    // add items to the component rendering context here
  },
});
