import Component from '@glimmer/component';

/**
 *
 * `Toolbar` renders a bar above a list that contains search or filtering elements
 *  as well as well as related buttons or links. A yielded Filters component is
 *  the container for any filter elements and will be displayed to the left of the
 *  bar, Actions is the container for buttons or links and is rendered to the right.
 *
 *
 * ```
 * <Toolbar as |T|>
 *   <T.Filters>
 *     <input type="search" />
 *   </T.Filters>
 *   <T.Actions>
 *     <a href="#">Link here</a>
 *   </T.Actions>
 * </Toolbar>
 * ```
 *
 * @class Toolbar
 * @yield {ToolbarFilters} Filters `Toolbar::Filters` component
 * @yield {ToolbarActions} Actions `Toolbar::Actions` component
 *
 */

export default class ToolbarComponent extends Component {}
