import Route from '@ember/routing/route';

export default class extends Route {
  redirect(): void {
    this.transitionTo('workspace.projects.project.app.deployment.operation-logs');
  }
}
