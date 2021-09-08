import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { UI } from 'waypoint-pb';
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

  get deploymentBundles(): UI.DeploymentBundle.AsObject[] {
    return this.model;
  }

  get hasMoreDeployments(): boolean {
    return this.deploymentBundles.filter((bundle) => bundle.deployment?.state == 4).length > 0;
  }

  get deploymentBundlesByGeneration(): GenerationGroup[] {
    let result: GenerationGroup[] = [];

    for (let bundle of this.deploymentBundles) {
      let { deployment } = bundle;
      let id = deployment?.generation?.id ?? deployment?.id ?? 'unknown';
      let group = result.find((group) => group.generationID === id);

      if (!group) {
        group = new GenerationGroup(id);
        result.push(group);
      }

      group.deploymentBundles.push(bundle);
    }

    return result;
  }

  @action
  showDestroyed(): void {
    this.isShowingDestroyed = true;
  }
}

class GenerationGroup {
  deploymentBundles: UI.DeploymentBundle.AsObject[] = [];

  constructor(public generationID: string) {}
}
