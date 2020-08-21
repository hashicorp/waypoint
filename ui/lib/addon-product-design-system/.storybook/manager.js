import { create } from '@storybook/theming/create';
import { addons } from '@storybook/addons';

addons.setConfig({
  showPanel: true,
  panelPosition: 'bottom',
  theme: create({
    base: 'light',
    brandTitle: 'Product Design System',
    brandImage: '',
    brandUrl: '#',
  })
});
