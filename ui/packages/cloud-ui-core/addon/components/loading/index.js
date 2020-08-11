import Component from '@glimmer/component';

/**
 *
 * `Loading` displays a loading icon by default but can also display a header
 * and a message for use when a page is "initializing" or loading any async
 * request.
 *
 *
 * ```
 * <Loading as |L|>
 *   <L.Elapsed />
 *   <L.Header>
 *     {{t 'components.page.hvns.detail.initializing.header'}}
 *   </L.Header>
 *   <L.Message>
 *     {{t 'components.page.hvns.detail.initializing.message'}}
 *   </L.Message>
 * </Loading>
 * ```
 *
 * @class Loading
 *
 */

export default class LoadingComponent extends Component {}
