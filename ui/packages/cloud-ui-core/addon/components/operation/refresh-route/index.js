import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';

/**
 *
 * `OperationRefreshRoute` is a renderless component that encapsulates route
 * refresh behavior. It uses the operation service to determine if the `model`
 * arg has any related operations that have recently changed. If there are
 * recently changed operations, then the `route` arg will be looked up and `refresh`
 * will be called on it, causing any model hooks on the route to re-run and load
 * any updated data.
 *
 * If `routeFallback` was passed and the `refresh` call fails, the component will
 * attempt to transition to the `routeFallback`.
 *
 *
 * ```
 * <Operation::RefreshRoute
 *   @route='cloud.consul.detail'
 *   @routeFallback='cloud.consul'
 *   @model={{@model.cluster}}
 * />
 * ```
 *
 * @class OperationRefreshRoute
 *
 */

export default class OperationRefreshRouteComponent extends Component {
  @service operation;
  @service router;

  /**
   * The `route` is the name of the route that will be refreshed if shouldRefresh is true.
   *
   * @argument route
   * @type {string}
   */

  /**
   * `routeFallback` is the name of a route to transition to if a refresh of `route`
   * fails.
   *
   * @argument routeFallback
   * @type {string}
   */

  /**
   * `model` is the object or list of objects to match against operation links
   *
   * @argument model
   * @type {array|object}
   */

  /*
   * `refreshRoute` looks up the route class for the named `this.args.route`.
   * If found, the route's class will have `refresh` called on it, triggering
   * all of the route hooks to be invalidated and re-fetched.
   *
   * This effectively will update the data from the route's model hook,
   * updating the UI if any of that data has changed.
   *
   */
  @action
  refreshRoute() {
    let { route, routeFallback } = this.args;
    let container = getOwner(this);
    let routeClass = container.lookup(`route:${route}`);

    routeClass.refresh().catch(e => {
      if (routeFallback) {
        return this.router.transitionTo(routeFallback);
      } else {
        throw e;
      }
    });
  }

  /*
   * `maybeRefresh` gets called by a modifier in the template that is triggered
   * anytime `this.operation.changedOperations` changes.
   *
   * If there are operations related to the passed `model` arg, then
   * `refreshRoute` will be called - otherwise the method returns early
   */
  @action
  maybeRefresh() {
    let modelIds = [this.args.model].flatMap(arr => arr).map(m => m.id);
    // this is the first filling of the operations service so we want
    // to skip it causing a refresh
    if (this.operation.firstFetch) {
      return;
    }
    // relatedOperations are any changeOperations from the operations service
    // where the link's uuid matches one of the current models' id
    let relatedOperations = this.operation.changedOperations.filter(op => {
      return modelIds.includes(op.link.uuid);
    });

    if (!relatedOperations.length) {
      return;
    }
    this.refreshRoute();
  }
}
