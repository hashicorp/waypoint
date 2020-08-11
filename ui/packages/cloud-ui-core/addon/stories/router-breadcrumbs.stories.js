import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'RouterBreadcrumbs',
  component: 'RouterBreadcrumbs',
};

// add stories by adding more exported functions
export let RouterBreadcrumbs = () => ({
  template: hbs`<RouterBreadcrumbs />`,
  context: {
    // add items to the component rendering context here
  },
});
