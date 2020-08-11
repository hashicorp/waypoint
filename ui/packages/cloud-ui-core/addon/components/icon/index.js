import Component from '@glimmer/component';
import { assert } from '@ember/debug';
import { SIZE_SCALE } from './consts';

/**
 *
 * `Icon` renders a `<Pds::Icon />` with some extra functionality.
 *
 * The SVG fill color can be changed by changing the color of the element.
 * The size can be overriden by using font sizes for the element.
 *
 *
 * ```
 * <Icon
 *   @type="chevron-up"
 *   @size="md"
 * />
 * ```
 *
 * @class Icon
 *
 */

/**
 *
 * The name of the SVG to render
 * @argument type
 * @type string
 *
 */

/**
 * A t-shirt size of the icon. Can be one of:
 * 'sm', 'md', 'lg', 'xl', '2xl'. Defaults to 'lg'.
 *
 *
 * @argument size
 * @default 'lg'
 *
 * @type string
 *
 */

export default class IconComponent extends Component {
  get size() {
    let { size } = this.args;
    if (size) {
      assert(
        `@size for ${this.toString()} must be one of the following: ${SIZE_SCALE.join(
          ', '
        )}, receieved: ${size}`,
        SIZE_SCALE.includes(size)
      );
    }
    return size || 'lg';
  }

  get iconClass() {
    return `icon--${this.size}`;
  }
}
