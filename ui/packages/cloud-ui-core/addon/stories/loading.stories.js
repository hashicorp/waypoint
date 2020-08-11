import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Loading',
  component: 'Loading',

  subcomponents: {
    'L.Elapsed': 'LoadingElapsed',
    'L.Header': 'LoadingHeader',
    'L.Message': 'LoadingMessage',
  },
};

// add stories by adding more exported functions
export let Loading = () => ({
  template: hbs`
    <Loading as |L|>
      <L.Elapsed></L.Elapsed>
      <L.Header>Some Header</L.Header>
      <L.Message>Some Message</L.Message>
    </Loading>
  `,
  context: {
    // add items to the component rendering context here
  },
});
