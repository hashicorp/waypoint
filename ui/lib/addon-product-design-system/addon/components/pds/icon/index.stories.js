import hbs from 'htmlbars-inline-precompile';
import types from '@hashicorp/structure-icons/dist/index';
import { withKnobs, select } from '@storybook/addon-knobs';
import DocsPage, { TITLE } from './docs.mdx';

const CONFIG = {
  title: TITLE,
  component: 'PdsIcon',
  decorators: [ withKnobs ],
  parameters: { docs: { page: DocsPage } },
};

const Index = () => ({
  template: hbs`
    <Pds::Icon
      @type={{type}}
    />
  `,
  context: {
    type: select('@type', types, 'bolt'),
  },
});

export {
  CONFIG as default,

  Index,
};
