import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'CodeBlock',
  component: 'CodeBlock',
};

export let CodeBlock = () => ({
  template: hbs`<CodeBlock>some code</CodeBlock>`,
});
