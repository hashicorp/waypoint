import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';

type Model = {
  releases: AppRouteModel['releases'];
  releaseDeploymentPairs: Record<number, number>;
};

export default class Releases extends Route {
  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let rdPairs = {};
    for (let release of app.releases) {
      let matchingDeployment = app.deployments.find((obj) => obj.id === release?.deploymentId);
      if (matchingDeployment) {
        rdPairs[release.sequence] = matchingDeployment.sequence;
      }
    }

    return {
      releases: app.releases,
      releaseDeploymentPairs: rdPairs,
    };
  }
}
