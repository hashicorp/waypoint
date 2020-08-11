import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select } from '@storybook/addon-knobs';
import { DEFAULT_VARIANT, COMPONENT_SCALE, VARIANT_SCALE } from 'cloud-ui-core/addon/components/typography/consts';

export default {
  title: 'Typography',
  component: 'Typography',
  decorators: [withKnobs],
};

// add stories by adding more exported functions
export let Typography = () => ({
  template: hbs`
    <Typography
      @component={{component}}
      @variant={{variant}}
    >
      Testing
    </Typography>
  `,
  context: {
    component: select('Component', COMPONENT_SCALE),
    variant: select('Variant', VARIANT_SCALE, DEFAULT_VARIANT),
  },
});
