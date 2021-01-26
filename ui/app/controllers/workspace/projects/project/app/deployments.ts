import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class WorkspaceProjectsProjectAppDeployments extends Controller {
  queryParams = ['destroyed']

  @tracked destroyed = false;

  get deployments() {
    if (this.destroyed) {
      return this.model
    } else {
      const deploys = this.model.filter(deployment => deployment.state != 4)
      return deploys;
    }
  }
}
