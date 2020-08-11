import hbs from 'htmlbars-inline-precompile';
import { withKnobs } from '@storybook/addon-knobs';
import DocsPage, { TITLE } from './docs.mdx';

const CONFIG = {
  title: TITLE,
  component: '<%= jsClass %>',
  decorators: [ withKnobs ],
  parameters: { docs: { page: DocsPage } },
};

const Index = () => ({
  template: hbs`<<%= tagName %> />`,
});

export {
  CONFIG as default,

  Index,
};
