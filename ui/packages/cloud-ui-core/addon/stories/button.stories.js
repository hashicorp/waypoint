import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select } from '@storybook/addon-knobs';
import { DEFAULT_VARIANT, DEFAULT_COMPACT, VARIANT_SCALE } from 'cloud-ui-core/addon/components/button/consts';

export default {
  title: 'Button',
  component: 'Button',
  decorators: [withKnobs],
};

// add stories by adding more exported functions
export let Button = () => ({
  template: hbs`
    <Button
      @variant={{variant}}
      @compact={{compact}}
    >
      Create Network
    </Button>
  `,
  context: {
    variant: select('Variant', [null, ...VARIANT_SCALE], DEFAULT_VARIANT),
    compact: select('Compact', [true, false], true),
  },
});

export let ButtonAllVariants = () => ({
  template: hbs`
    <Button @variant="primary" @compact={{compact}}>Primary</Button>
    <Button @variant="secondary" @compact={{compact}}>Secondary</Button>
    <Button @variant="warning" @compact={{compact}}>Warning</Button>
    <Button @variant="ghost" @compact={{compact}}>Ghost</Button>
    <Button @variant="ghost-background" @compact={{compact}}>Ghost</Button>
    <Button @variant="link" @compact={{compact}}>Link</Button>
  `,
  context: {
    compact: select('Compact', [true, false], DEFAULT_COMPACT),
  },
});