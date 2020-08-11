import Component from '@glimmer/component';

/**
 *
 * `Select` renders a `<select>` with associated `<options>` from an array that is passed to the component.
 *
 *
 * ```
 * <Select
 *   @options={{array '1' '2' '3'}}
 *   @value={{this.value}}
 *   {{on 'change' (set this 'value' (get _ 'target.selected))}}
 * />
 * ```
 *
 * @class Select
 *
 */

export default class SelectComponent extends Component {
  /**
   * The list of options to render in the select element. This can be an array of strings objects. If the members of options are objects, you will need to pass `valuePath` and `labelPath` in addition to `options`.
   * @argument options
   * @type {Array}
   * @required
   */
  /**
   * The current value of the `<select>`. `value` will be used to determine if an option should be marked as `selected`.
   * @argument value
   * @type {string}
   * @required
   */
  /**
   * When options is a list of objects, `valuePath` tells us the path to lookup a value for the purpose of setting `value` and `selected` on the associated `<option>` element.
   * @argument valuePath
   * @type {string}
   */
  /**
   * When options is a list of objects, `labelPath` tells us the path to lookup a vaule for the purpose of setting `value` and `selected` on the associated `<option>` element.
   * @argument labelPath
   * @type {string}
   */
}
