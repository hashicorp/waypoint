import hbs from 'htmlbars-inline-precompile';
import types from '@hashicorp/structure-icons/dist/index';
import { withKnobs, select } from '@storybook/addon-knobs';
import { SIZE_SCALE } from 'cloud-ui-core/addon/components/icon/consts';
export default {
  title: 'Icon',
  component: 'Icon',
  decorators: [withKnobs],
};

export let basic = () => ({
  template: hbs`<Icon @size={{size}} @type={{type}} />`,
  context: {
    type: select('Type', types, 'bolt'),
    size: select('Size', SIZE_SCALE, 'lg'),
  },
});
