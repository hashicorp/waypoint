import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Loading/Elapsed',
  component: 'LoadingElapsed',
};

// add stories by adding more exported functions
export let LoadingElapsed = () => ({
  template: hbs`<Loading::Elapsed />`,
  context: {
    // add items to the component rendering context here
  },
});
