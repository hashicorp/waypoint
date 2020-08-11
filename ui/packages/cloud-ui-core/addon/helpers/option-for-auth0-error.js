import { helper as buildHelper } from '@ember/component/helper';

const I18N_PREFIX = 'helpers.options-for-auth0-error';

// AUth0 Error preview docs: https://auth0.com/docs/libraries/error-messages
export const DEFAULT_ERROR_CODE = 'default';
export const ERROR_CODE_MAPPING = {
  access_denied: {
    message: `${I18N_PREFIX}.access_denied.message`,
  },
  invalid_user_password: {
    message: `${I18N_PREFIX}.invalid_user_password.message`,
  },
  mfa_invalid_code: {
    message: `${I18N_PREFIX}.mfa_invalid_code.message`,
  },
  mfa_registration_required: {
    message: `${I18N_PREFIX}.mfa_registration_required.message`,
  },
  mfa_required: {
    message: `${I18N_PREFIX}.mfa_required.message`,
  },
  password_leaked: {
    message: `${I18N_PREFIX}.password_leaked.message`,
  },
  PasswordHistoryError: {
    message: `${I18N_PREFIX}.PasswordHistoryError.message`,
  },
  PasswordStrengthError: {
    message: `${I18N_PREFIX}.PasswordStrengthError.message`,
  },
  too_many_attempts: {
    message: `${I18N_PREFIX}.too_many_attempts.message`,
  },
  unauthorized: {
    message: `${I18N_PREFIX}.unauthorized.message`,
  },
  [DEFAULT_ERROR_CODE]: {
    message: `${I18N_PREFIX}.default.message`,
  },
};

export const ERROR_CODE_SCALE = Object.keys(ERROR_CODE_MAPPING);

export function optionForAuth0Error([code = 'default']) {
  let option = ERROR_CODE_MAPPING[code] || ERROR_CODE_MAPPING['default'];

  return {
    ...option,
  };
}

export default buildHelper(optionForAuth0Error);
