import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select, text } from '@storybook/addon-knobs';

export default {
  title: 'AlertBanner',
  component: 'AlertBanner',
  decorators: [withKnobs],
};

let VARIANT = ['error', 'info', 'success', 'warning'];

/**
 *
 * ⚠️Yieldable Named Blocks are currently not working in Ember Storybook but
 * will be fixed with the Octane 3.20 Release.
 *  In the meantime, don't be alarmed when you will see double yields😿😿.
 *
 *
 */
export let AlertBanner = () => ({
  template: hbs`
    <AlertBanner @variant={{variant}}>
      <:title>{{title}}</:title>
      <:content>{{content}}</:content>
      <:action>
        <Button @variant='link' @compact={{true}}>
          {{actionText}}
        </Button>
      </:action>
    </AlertBanner>
  `,
  context: {
    variant: select('@variant', VARIANT, 'info'),
    title: text('title', 'Alert Title'),
    content: text('content', 'Sweet caramels ice cream cupcake carrot cake chocolate cake.'),
    actionText: text('action', 'Request more sweets'),
  },
});
