import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Box',
  component: 'Box',
};

export let Box = () => ({
  template: hbs`<Box>some content</Box>`,
  context: {},
});

export let Box2XSmallPaddingSquare = () => ({
  template: hbs`<Box @padding="2xs">some content</Box>`,
  context: {},
});

export let BoxXSmallPaddingSquare = () => ({
  template: hbs`<Box @padding="xs">some content</Box>`,
  context: {},
});

export let BoxMediumPaddingNonSquare = () => ({
  template: hbs`<Box @padding="md">some content</Box>`,
  context: {},
});

export let BoxLargePaddingNonSquare = () => ({
  template: hbs`<Box @padding="lg">some content</Box>`,
  context: {},
});

export let BoxXLPaddingNonSquare = () => ({
  template: hbs`<Box @padding="xl">some content</Box>`,
  context: {},
});

export let Box2XLargePaddingNonSquare = () => ({
  template: hbs`<Box @padding="2xl">some content</Box>`,
  context: {},
});
