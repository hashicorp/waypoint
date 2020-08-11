import Component from '@glimmer/component';

/**
 *
 * `ZeroState` is displayed when a page has nothing to display. You may render
 * a header, a message, and an action.
 *
 *
 * ```
 * <ZeroState as |ZS|>
 *   <ZS.Header>
 *     {{t 'components.page.hvns.list.empty.header'}}
 *   </ZS.Header>
 *   <ZS.Message>
 *     {{t 'components.page.hvns.list.empty.message'}}
 *   </ZS.Message>
 *   <ZS.Action>
 *     <button type='submit'>
 *       {{t 'components.page.hvns.create.title'}}
 *     </button>
 *   </ZS.Action>
 * </ZeroState>
 * ```
 *
 * @class ZeroState
 *
 */

export default class ZeroStateComponent extends Component {}
