import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select } from '@storybook/addon-knobs';

export default {
  title: 'Paper',
  component: 'Paper',
  decorators: [withKnobs],
};

export let Paper = () => ({
  template: hbs`<Paper @variant={{variant}} @square={{square}}>some content</Paper>`,
  context: {
    square: select('Square', [true, false], false),
    variant: select('Variant', [null, 'outlined'], null),
  },
});
