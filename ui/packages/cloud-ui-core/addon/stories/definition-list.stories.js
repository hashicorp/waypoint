import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'DefinitionList',
  component: 'DefinitionList',
};

export let DefinitionList = () => ({
  template: hbs`
    <DefinitionList as |DL|>
      <DL.Key>Term1</DL.Key>
      <DL.Value>Value1</DL.Value>
      <DL.Key>Term2</DL.Key>
      <DL.Value>Value2</DL.Value>
      <DL.Key>Term3</DL.Key>
      <DL.Value>Term3</DL.Value>
    </DefinitionList>
  `,
  context: {},
});
