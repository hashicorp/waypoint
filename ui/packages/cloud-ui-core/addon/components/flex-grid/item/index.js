import Component from '@glimmer/component';

/**
 *
 * `FlexGridItem` applies a class on a container to control the width of an item
 * within a `Grid`.
 *
 *
 * ```
 * <FlexGrid::Item
 *   @xs="6"
 *   @sm="6"
 *   @md="6"
 *   @lg="6"
 * >
 *   content
 * </FlexGrid::Item>
 * ```
 *
 * @class FlexGridItem
 *
 */

export default class FlexGridItemComponent extends Component {
  /**
   * `xs` applies the extra small class for this item.
   * @argument xs
   * @type {string}
   */
  /**
   * `xsOffset` applies the extra small offset class for this item.
   * @argument xsOffset
   * @type {string}
   */
  /**
   * `sm` applies the extra small class for this item.
   * @argument sm
   * @type {string}
   */
  /**
   * `smOffset` applies the small offset class for this item.
   * @argument smOffset
   * @type {string}
   */
  /**
   * `md` applies the extra small class for this item.
   * @argument md
   * @type {string}
   */
  /**
   * `mdOffset` applies the medium offset class for this item.
   * @argument mdOffset
   * @type {string}
   */
  /**
   * `lg` applies the extra small class for this item.
   * @argument lg
   * @type {string}
   */
  /**
   * `lgOffset` applies the large offset class for this item.
   * @argument lgOffset
   * @type {string}
   */
}
