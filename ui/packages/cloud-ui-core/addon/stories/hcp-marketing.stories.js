import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'HcpMarketing',
  component: 'HcpMarketing',
};

// add stories by adding more exported functions
export let HcpMarketing = () => ({
  template: hbs`<HcpMarketing />`,
  context: {
    // add items to the component rendering context here
  },
});
