import hbs from 'htmlbars-inline-precompile';
import { withKnobs, text, boolean } from '@storybook/addon-knobs';
import DocsPage, { TITLE } from './docs.mdx';

export default {
  title: TITLE,
  component: 'DocsTextField',
  decorators: [ withKnobs ],
  parameters: { docs: { page: DocsPage } },
};

export const Index = () => ({
  template: hbs`<Docs::TextField
    @errorMessage={{errorMessage}}
    @helpText={{helpText}}
    @isDirty={{isDirty}}
    @isDisabled={{isDisabled}}
    @isInvalid={{isInvalid}}
    @isRequired={{isRequired}}
    @name={{name}}
    @value={{value}}
  />`,
  context: {
    // Field
    name: text('Label', 'Username', 'Field'),
    helpText: text('Help Text', 'Enter a unique username.', 'Field'),
    errorMessage: text('Error Message', 'Username is unavailable.', 'Field'),
    // Input
    isDisabled: boolean(':disabled', false, 'Input'),
    isRequired: boolean(':required', false, 'Input'),
    value: text('Value', 'c@pta1nMarv3l', 'Input'),
    isDirty: boolean('.pds-dirty', false, 'Input'),
    isInvalid: boolean('.pds-invalid', true, 'Input'),
  }
});
