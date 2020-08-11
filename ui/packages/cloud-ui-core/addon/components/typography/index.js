import Component from '@glimmer/component';
import { assert } from '@ember/debug';
import { DEFAULT_VARIANT, DEFAULT_VARIANT_MAPPING, DEFAULT_CLASS_MAPPING } from './consts';

/**
 *
 * `Typography` applies styles based on align, color, and variant properties
 * for raw text or as a wrapper element.
 *
 *
 * ```
 * <Typography
 *   @align="center"
 *   @color="primary"
 *   @variant="h1"
 * />
 * ```
 *
 * @class Typography
 *
 */

export default class TypographyComponent extends Component {
  /**
   * The typography alignment to apply.
   * @argument align
   * @type {?string}
   */

  /**
   * The typography color to apply.
   * @argument color
   * @type {?string}
   */

  /**
   * The typography component tag to apply. This will overwrite variantMapping.
   * @argument component
   * @type {?string}
   */

  /**
   * The typography variant to render.
   * @argument Typography#variant
   * @type {?string}
   */

  /**
   * The typography variant to render but also returns the default if nothing
   * is set..
   * @method Typography#variant
   * @type {?string}
   */
  get variant() {
    let { variant = DEFAULT_VARIANT } = this.args;
    let variantScale = Object.keys(this.variantMapping);

    if (variant) {
      assert(
        `@variant for ${this.toString()} must be one of the following: ${variantScale.join(
          ', '
        )}, receieved: ${variant}`,
        variantScale.includes(variant)
      );
    }

    return variant;
  }

  /**
   * Get the variant mapping.
   * @method Typography#variantMapping
   * @return {string} The map to find a tag based on a variant.
   */
  get variantMapping() {
    return this.args.variantMapping || DEFAULT_VARIANT_MAPPING;
  }

  /**
   * Get a tag to render based on the argument `variantMapping` or, if that's
   * not passed in, the default mapping will be used.
   * @method Typography#getComponentTag
   * @return {string} The html tag to use in a dynamic element render.
   */
  get componentTag() {
    return this.args.component || this.variantMapping[this.variant] || 'span';
  }

  /**
   * Get a class string to set based on the class mapping.
   * @method Typography#componentClasses
   * @return {string} The classes to apply to the Typography element.
   */
  get componentClasses() {
    let classes = DEFAULT_CLASS_MAPPING[this.variant] || [];
    classes.push(this.variantMapping[this.variant]);
    return classes.join(' ');
  }
}
