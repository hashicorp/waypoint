import Component from '@glimmer/component';
import { assert } from '@ember/debug';
import { PADDING_SIZE_SCALE } from './consts';

/**
 *
 * `Box` is a container to control spacing within an element.
 *
 *
 * ```
 * <Box
 *   @padding="sm"
 * />
 * ```
 *
 * @class Box
 *
 */

export default class BoxComponent extends Component {
  /**
   * Controls the inner padding of the Box container and aligns to the PDS
   * size scale. A single value or an array of values are accepted values.
   * Options: '2xs', 'xs', 'sm', 'md', 'lg', 'xl', '2xl'
   * @argument padding
   * @default 'sm'
   * @type {string}
   */

  get padding() {
    return this.getSpacingArray();
  }

  getSpacingArray() {
    let { padding = '' } = this.args;
    let spacing;

    if (!padding) {
      spacing = ['sm', 'sm', 'sm', 'sm'];
    } else {
      spacing = padding.split(' ');
    }

    switch (spacing.length) {
      case 1: {
        spacing = [spacing[0], spacing[0], spacing[0], spacing[0]];
        break;
      }
      case 2: {
        spacing = [spacing[0], spacing[1], spacing[0], spacing[1]];
        break;
      }
      case 3: {
        spacing = [spacing[0], spacing[1], spacing[2], spacing[1]];
        break;
      }
    }

    for (let size of spacing) {
      assert(
        `@padding size for ${this.constructor.name} must be one of the following: ${PADDING_SIZE_SCALE.join(
          ', '
        )}, receieved: ${size}`,
        PADDING_SIZE_SCALE.includes(size)
      );
    }

    return spacing;
  }
}
