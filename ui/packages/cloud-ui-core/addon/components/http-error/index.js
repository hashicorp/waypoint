import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 *
 * `HttpError` wraps an Error component to help with rendering HTTP Errors.
 *
 *
 * ```
 * <HttpError
 *   @code="404"
 * />
 * ```
 *
 * @class HttpError
 *
 */

export default class HttpErrorComponent extends Component {

  @service config;
  /**
   * A route path for the Go back link.
   * @argument previousRoute
   * @type {?string}
   */

  /**
   * A string to override the default title for the error code.
   * @argument title
   * @type {?string}
   */

  /**
   * A string to override the default message for the error code.
   * @argument message
   * @type {?string}
   */

  /**
   * The HTTP Error code to use when mapping error messaging.
   * @argument code
   * @type {string}
   */

  get previousRoute() {
    return this.args.previousRoute || this.config?.app?.defaultRoute;
  }
}
