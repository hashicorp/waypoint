import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select } from '@storybook/addon-knobs';
import { SOURCE_SCALE, VARIANT_SCALE } from 'cloud-ui-core/addon/helpers/option-for-icon-badge';

export default {
  title: 'IconBadge',
  component: 'IconBadge',
  decorators: [withKnobs],
};

export let IconBadge = () => ({
  template: hbs`
    <IconBadge
      @highlightLabel={{highlightLabel}}
      @source={{source}}
      @variant={{variant}}
    />
  `,
  context: {
    source: select('Source', [...SOURCE_SCALE], SOURCE_SCALE[0]),
    variant: select('Variant', [...VARIANT_SCALE], VARIANT_SCALE[0]),
    highlightLabel: select('Highlight Label', [true, false], false),
  },
});

export let IconBadgeWithLabel = () => ({
  template: hbs`
    <IconBadge
      @label="My Label"
      @source={{source}}
      @variant={{variant}}
    />
  `,
  context: {
    source: select('Source', [...SOURCE_SCALE], SOURCE_SCALE[0]),
    variant: select('Variant', [...VARIANT_SCALE], VARIANT_SCALE[0]),
  },
});
