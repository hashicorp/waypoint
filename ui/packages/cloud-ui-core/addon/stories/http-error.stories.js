import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select, text } from '@storybook/addon-knobs';
import { DEFAULT_ERROR_CODE, ERROR_CODE_SCALE } from 'cloud-ui-core/addon/helpers/option-for-http-error';

export default {
  title: 'HttpError',
  component: 'HttpError',
  decorators: [withKnobs],
};

export let HttpError = () => ({
  template: hbs`<HttpError @code={{code}} @title={{title}} @message={{message}} @previousRoute={{previousRoute}} />`,
  context: {
    code: select('Error Code', ERROR_CODE_SCALE, DEFAULT_ERROR_CODE),
    message: text('Message', null),
    previousRoute: select('Previous Route', [null, 'cloud', 'cloud.orgs'], null),
    title: text('Title', null),
  },
});
