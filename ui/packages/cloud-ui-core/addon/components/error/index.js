import Component from '@glimmer/component';
import { DEFAULT_ICON_TYPE } from './consts';

/**
 *
 * `Error` displays relevant error messaging.
 *
 *
 * ```
 * <Error @iconType="cancel-square-outline">
 *   <:title>Not Found</:title>
 *   <:subtitle>Error 404</:subtitle>
 *   <:content>Some content message</:content>
 *   <:footer>
 *    <LinkTo @route="cloud">
 *      <Icon @type='chevron-left' @size='sm' aria-hidden='true' />
 *      Go back
 *    </LinkTo>
 *  </:footer>
 * </Error>
 * ```
 *
 * @class Error
 *
 */

export default class ErrorComponent extends Component {
  /**
   * The icon to render via an Icon component arg.
   * @argument iconType
   * @type {string}
   */

  get iconType() {
    let { iconType } = this.args;
    return iconType || DEFAULT_ICON_TYPE;
  }
}
