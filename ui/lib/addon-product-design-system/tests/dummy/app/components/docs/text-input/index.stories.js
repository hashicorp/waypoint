import hbs from 'htmlbars-inline-precompile';
import { withKnobs, boolean } from '@storybook/addon-knobs';
import DocsPage, { TITLE } from './docs.mdx';

export default {
  title: TITLE,
  component: 'DocsTextInput',
  decorators: [ withKnobs ],
  parameters: { docs: { page: DocsPage } },
};

export const Index = () => ({
  template: hbs`
    <Docs::TextInput
      @dirty={{dirty}}
      @disabled={{disabled}}
      @focused={{focused}}
      @hovered={{hovered}}
      @invalid={{invalid}}
      @required={{required}}
    />
  `,
  context: {
    required: boolean(':required', false),
    disabled: boolean(':disabled', false),
    focused:  boolean(':focus', false),
    hovered:  boolean(':hover', false),
    invalid:  boolean('.pds-invalid', false),
    dirty: boolean('.pds-dirty', false),
  }
});
