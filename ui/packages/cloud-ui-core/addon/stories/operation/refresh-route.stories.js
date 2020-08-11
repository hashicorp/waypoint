import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Operation::RefreshRoute',
  component: 'OperationRefreshRoute',
};

// add stories by adding more exported functions
export let OperationRefreshRoute = () => ({
  template: hbs`<Operation::RefreshRoute />`,
  context: {
    // add items to the component rendering context here
  }
});
