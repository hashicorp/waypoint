import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { Deployment, UI } from 'waypoint-pb';
import { action } from '@ember/object';
export default class WorkspaceProjectsProjectAppDeployments extends Controller {
  queryParams = [
    {
      isShowingDestroyed: {
        as: 'destroyed',
      },
    },
  ];

  @tracked isShowingDestroyed = false;

  get hasMoreDeployments(): boolean {
    return (
      this.model.filter((bundle: UI.DeploymentBundle.AsObject) => bundle.deployment?.state == 4).length > 0
    );
  }

  get deployments(): UI.DeploymentBundle.AsObject[] {
    return this.model;
  }

  get deploymentsByGeneration(): GenerationGroup[] {
    let result: GenerationGroup[] = [];

    for (let bundle of this.deployments) {
      let { deployment } = bundle;

      if (!deployment) {
        continue;
      }
      let id = deployment.generation?.id ?? deployment.id;
      let group = result.find((group) => group.generationID === id);

      if (!group) {
        group = new GenerationGroup(id);
        result.push(group);
      }

      group.bundles.push(bundle);
    }

    return result;
  }

  @action
  showDestroyed(): void {
    this.isShowingDestroyed = true;
  }
}

class GenerationGroup {
  bundles: UI.DeploymentBundle.AsObject[] = [];

  constructor(public generationID: string) {}
}
