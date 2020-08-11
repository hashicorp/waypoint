import Route from '@ember/routing/route';

export function makeRouteStub({ breadcrumb = {} } = {}) {
  return class extends Route {
    breadcrumb = breadcrumb;
  };
}
