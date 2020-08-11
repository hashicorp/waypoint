import Component from '@glimmer/component';

/**
 *
 * `Crumb` component generates a LinkTo component.
 *
 *
 * * ## Example usage
 *
 * ```
 * <Crumb @route="some.route" @models={{(array 1 2)}}>
 *   Resource
 * </Crumb>
 * ```
 *
 * @class Crumb
 *
 */

export default class CrumbComponent extends Component {
  /**
   * The array of models to represent the route. This could be anything LinkTo can
   * accept as a valid route @models arg.
   * @argument models
   * @type {string[]}
   */

  /**
   * The string value to represent the route. This could be anything LinkTo can
   * accept as a valid route.
   * @argument route
   * @type {string}
   */

  get models() {
    return this.args.models || [];
  }
}
