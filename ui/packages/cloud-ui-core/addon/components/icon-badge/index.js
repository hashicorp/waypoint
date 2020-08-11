import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 *
 * `IconBadge` renders an icon and label based on a source arg and variant arg.
 *
 *
 * ```
 * <IconBadge
 *   @source="region"
 *   @variant="aws"
 * />
 * ```
 *
 * @class IconBadge
 *
 */

export default class IconBadgeComponent extends Component {
  @service intl;

  /**
   * Adds the text color class to the label as well as the icon.
   * @argument highlightLabel
   * @type {?string}
   */

  /**
   * Overrides the label of the IconBadge.
   * @argument label
   * @type {?string}
   */

  /**
   * Changes the source of the IconBadge variants.
   * @argument source
   * @type {string}
   */

  /**
   * Changes the text and icon of the IconBadge.
   * @argument variant
   * @type {string}
   */
}
