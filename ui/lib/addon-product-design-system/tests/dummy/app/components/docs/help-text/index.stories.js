import hbs from 'htmlbars-inline-precompile';
import DocsPage, { TITLE } from './docs.mdx';

export default {
  title: TITLE,
  component: 'DocsHelpText',
  parameters: { docs: { page: DocsPage } },
};

export const Index = () => ({
  template: hbs`<Docs::HelpText />`,
});
