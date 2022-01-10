import Route from '@ember/routing/route';
import { Model as AppRouteModel } from '../app';

type Model = {
  builds: AppRouteModel['builds'];
  buildDeploymentPairs: Record<number, number>;
};

export default class Builds extends Route {
  async model(): Promise<Model> {
    let app = this.modelFor('workspace.projects.project.app') as AppRouteModel;
    let bdPairs = {};
    for (let build of app.builds) {
      let matchingDeployment = app.deployments.find((obj) => {
        if (obj.preload && obj.preload.artifact) {
          return obj.preload.artifact.id === build.pushedArtifact?.id;
        }
        return obj.artifactId === build.pushedArtifact?.id;
      });

      if (matchingDeployment) {
        bdPairs[build.sequence] = matchingDeployment.sequence;
      }
    }

    return {
      builds: app.builds,
      buildDeploymentPairs: bdPairs,
    };
  }
}
