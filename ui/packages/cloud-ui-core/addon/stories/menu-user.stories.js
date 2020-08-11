import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'MenuUser',
  component: 'MenuUser',
};

// add stories by adding more exported functions
export let basic = () => ({
  template: hbs`<MenuUser />`,
  context: {
    // add items to the component rendering context here
  },
});
