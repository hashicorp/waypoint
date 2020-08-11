import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Loading/Header',
  component: 'LoadingHeader',
};

// add stories by adding more exported functions
export let LoadingHeader = () => ({
  template: hbs`<Loading::Header>Some Header</Loading::Header>`,
  context: {
    // add items to the component rendering context here
  },
});
