import Component from '@glimmer/component';
import { assert } from '@ember/debug';
import { VARIANT_SCALE, DEFAULT_VARIANT, DEFAULT_VARIANT_MAPPING } from './consts';

/**
 *
 * `Button` renders a button that is styled based on variant and variant.
 *
 *
 * ```
 * <Button
 *   @variant="primary"
 *   @compact={{true}}
 * /></Button>
 * ```
 *
 * @class Button
 *
 */

export default class ButtonComponent extends Component {
  /**
   * Changes the variant of the Button.
   * @argument variant
   * @type {?string}
   */

  /**
   * Changes the padding of the Button.
   * @argument compact
   * @type {?boolean}
   */

  /**
   * Gets the variant arg or falls back to the default.
   * @method Button#variant
   * @return {string}
   */
  get variant() {
    let { variant = DEFAULT_VARIANT } = this.args;

    if (variant) {
      assert(
        `@variant for ${this.toString()} must be one of the following: ${VARIANT_SCALE.join(
          ', '
        )}, receieved: ${variant}`,
        VARIANT_SCALE.includes(variant)
      );
    }

    return variant;
  }

  /**
   * Get a class to apply to the button based on the variant argument.
   * @method Button#variantClass
   * @return {string} The css class to apply to the Button.
   */
  get variantClass() {
    return DEFAULT_VARIANT_MAPPING[this.variant] || DEFAULT_VARIANT_MAPPING[DEFAULT_VARIANT];
  }
}
