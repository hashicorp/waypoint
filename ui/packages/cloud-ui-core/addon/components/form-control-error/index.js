import Component from '@glimmer/component';

/**
 *
 * `FormControlError` displays an `<Icon @type='cancel-square-fill' @size='md' aria-hidden='true' />`
 *  along with an error message.<br />This component is intended to be used for form field error messaging.
 *
 *
 * ```
 * <FormControlError >
    <:message>{{this.errors.field_violations.name}}</:message>
 * </FormControlError>
 * ```
 *
 * @class FormControlError
 *
 */

export default class FormControlErrorComponent extends Component {}
