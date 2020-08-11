import Component from '@glimmer/component';

/**
 *
 * `AlertBanner` component is used to display messaging, there are 4 types: error, info, success and warning.
 *
 *
 * ```
 * <AlertBanner>
 *   <:title>Page Title</:title>
 *   <:content>Sweet roll ice cream cupcake carrot cake chocolate cake.</:content>
 *   <:action><Button @variant='link' @compact={{true}}>Click here</Button></:action>
 * </AlertBanner>
 * ```
 *
 * @class AlertBanner
 *
 *
 */

export default class AlertBannerComponent extends Component {
  /**
   * Determines which alert-banner theme to apply.  Possibly values are `error`, `info`, `success` and `warning`.
   * @argument variant
   * @default 'info'
   * @type {string}
   */
  get variant() {
    return this.args.variant || 'info';
  }
  get iconType() {
    if (this.args.variant === 'error') {
      return 'cancel-square-fill';
    }

    if (this.args.variant === 'info') {
      return 'info-circle-fill';
    }

    if (this.args.variant === 'success') {
      return 'check-circle-fill';
    }

    if (this.args.variant === 'warning') {
      return 'alert-triangle';
    }

    return 'info-circle-fill';
  }
}
