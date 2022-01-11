import { DeploymentExtended, ReleaseExtended } from 'waypoint/services/api';

import { Model as AppRouteModel } from '../app';
import Route from '@ember/routing/route';
import { StatusReport } from 'waypoint-pb';

type Model = AppRouteModel['deployments'];

interface ResourceMap {
  resource: StatusReport.Resource.AsObject;
  type: string;
  source: DeploymentExtended | ReleaseExtended;
}
export default class Resources extends Route {
  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;

    let deployments = app.deployments;
    let releases = app.releases;

    let resources: ResourceMap[] = [];

    deployments.forEach((dep) => {
      dep.statusReport?.resourcesList.forEach((resource) => {
        resources.push({
          resource,
          type: 'deployment',
          source: dep,
        } as ResourceMap);
      });
    });
    releases.forEach((rel) => {
      rel.statusReport?.resourcesList.forEach((resource) => {
        resources.push({
          resource,
          type: 'release',
          source: rel,
        } as ResourceMap);
      });
    });

    return resources;
  }
}
