import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { Deployment } from 'waypoint-pb';
import { action } from '@ember/object';
export default class WorkspaceProjectsProjectAppDeployments extends Controller {
  queryParams = ['destroyed'];

  @tracked destroyed = false;

  get hasMoreDeployments(): boolean {
    return this.model.filter((deployment: Deployment.AsObject) => deployment.state == 4).length > 0;
  }

  get deployments(): Deployment.AsObject[] {
    if (this.destroyed) {
      return this.model;
    } else {
      let deploys = this.model.filter((deployment: Deployment.AsObject) => deployment.state != 4);
      return deploys;
    }
  }

  get deploymentsByGeneration(): GenerationGroup[] {
    let result: GenerationGroup[] = [];

    for (let deployment of this.deployments) {
      let id = deployment.generation?.id ?? deployment.id;
      let group = result.find((group) => group.generationID === id);

      if (!group) {
        group = new GenerationGroup(id);
        result.push(group);
      }

      group.deployments.push(deployment);
    }

    return result;
  }

  @action
  showDestroyed(): void {
    this.destroyed = true;
  }
}

class GenerationGroup {
  deployments: Deployment.AsObject[] = [];

  constructor(public generationID: string) {}
}
