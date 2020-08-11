import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { capitalize } from '@ember/string';
import { getOwner } from '@ember/application';

import { BREADCRUMBS_EXCLUDED_ROUTES } from './consts';

/**
 *
 * `RouterBreadcrumbs` automatically generates breadcrumbs based on the router.
 *
 *
 * ### API
 * This component looks up the current route's class and can utilize
 * breadcrumb-specific properties for generating custom breadcrumb titles or for
 * overriding an entire breadcrumb path.
 *
 *
 * #### breadcrumb (@type {Object})
 * `breadcrumb`: This can be an object or a getter method that returns an object.
 *
 *
 * `breadcrumb.title`: A string to override the label of the crumb.
 *
 *
 * #### breadcrumbs (@type {Function})
 * `breadcrumbs`: A function that returns an array of crumbs. Each crumb should
 *     have at least a route key with an optional models key.
 *
 *
 * `breadcrumbs[].route`: The route path to send to a LinkTo @route arg.
 *
 *
 * `breadcrumbs[].models`: An array of full models which contain an `id` key and
 *     any other model data to be used in generating a label and a link. Each
 *     model eventually gets mapped down to its `id` to be sent to LinkTo @models
 *     arg.
 *
 *
 * ```
 * <RouterBreadcrumbs
 *   @hide={{true}}
 * />
 * ```
 *
 * @class RouterBreadcrumbs
 *
 */

export default class RouterBreadcrumbsComponent extends Component {
  @service router;

  /**
   * `hide` allows you to block the rendering of breadcrumbs.
   * @argument hide
   * @type {?boolean}
   */

  get crumbs() {
    let route = getOwner(this).lookup(`route:${this.router.currentRouteName}`);
    if (!route) {
      return [];
    }
    let breadcrumbs = route.breadcrumbs ? route.breadcrumbs : null;

    return this.args.hide
      ? []
      : this.hydrateCrumbs(breadcrumbs ? breadcrumbs : this.buildCrumbs(this.router.currentRouteName));
  }

  buildCrumbs(routeName = '') {
    return routeName.split('.').reduce((crumbs, part, index, parts) => {
      // Don't add a breadcrumb for excluded routes or if there's no 'part'.
      if (!part || BREADCRUMBS_EXCLUDED_ROUTES.includes(part)) {
        return crumbs;
      }

      let path = parts.slice(0, index + 1).join('.');
      let route = getOwner(this).lookup(`route:${path}`);

      // Explicitly setting the crumb to null in the route drops it from the list
      if (route.breadcrumb === null) {
        return crumbs;
      }

      crumbs.push({
        route: path,
        models: [],
      });

      return crumbs;
    }, []);
  }

  hydrateCrumbs(crumbs = []) {
    return crumbs.map(crumb => {
      let path = crumb.route;
      let route = getOwner(this).lookup(`route:${path}`);
      let parts = path.split('.');
      let title = parts[parts.length - 1]
        .split('.')
        .map(capitalize)
        .join(' ');

      return {
        ...crumb,
        title,
        models:
          crumb.models && crumb.models.length
            ? crumb.models.map(function(c) {
                return c ? c.id : c;
              })
            : [],
        ...(route.breadcrumb && typeof route.breadcrumb === 'function' ? route.breadcrumb(crumb.models) : {}),
      };
    });
  }
}
