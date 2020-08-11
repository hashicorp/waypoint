import { helper as buildHelper } from '@ember/component/helper';

const I18N_PREFIX = 'helpers.options-for-http-error';

export const DEFAULT_ERROR_CODE = 'default';
export const ERROR_CODE_MAPPING = {
  403: { label: `${I18N_PREFIX}.403.label`, message: `${I18N_PREFIX}.403.message`, iconType: 'disabled' },
  404: {
    label: `${I18N_PREFIX}.404.label`,
    message: `${I18N_PREFIX}.404.message`,
    iconType: 'help-circle-outline',
  },
  500: {
    label: `${I18N_PREFIX}.500.label`,
    message: `${I18N_PREFIX}.500.message`,
    iconType: 'alert-circle-outline',
  },
  [DEFAULT_ERROR_CODE]: {
    label: `${I18N_PREFIX}.default.label`,
    message: `${I18N_PREFIX}.default.message`,
    iconType: 'alert-circle-outline',
  },
};

export const ERROR_CODE_SCALE = Object.keys(ERROR_CODE_MAPPING);

export function optionForHttpError([code = 'default']) {
  let option = ERROR_CODE_MAPPING[code] || ERROR_CODE_MAPPING['default'];

  return {
    ...option,
  };
}

export default buildHelper(optionForHttpError);
