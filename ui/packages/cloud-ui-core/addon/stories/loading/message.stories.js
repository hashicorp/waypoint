import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Loading/Message',
  component: 'LoadingMessage',
};

// add stories by adding more exported functions
export let LoadingMessage = () => ({
  template: hbs`<Loading::Message>Some Message</Loading::Message>`,
  context: {
    // add items to the component rendering context here
  },
});
