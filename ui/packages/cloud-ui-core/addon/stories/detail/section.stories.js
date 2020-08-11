import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Detail/Section',
  component: 'DetailSection',
};

export let DetailSection = () => ({
  template: hbs`
    <Detail::Section @title="Title">
    </Detail::Section>
  `,
  context: {},
});

export let DetailSectionWithZeroState = () => ({
  template: hbs`
    <Detail::Section @title="Title" as |DS|>
      <DS.ZeroState>Some Empty Message</DS.ZeroState>
    </Detail::Section>
  `,
  context: {},
});
