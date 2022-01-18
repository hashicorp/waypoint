import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import SessionService from 'ember-simple-auth/services/session';

export default class Workspaces extends Route {
  @service api!: ApiService;
  @service session!: SessionService;

  async redirect(): Promise<void> {
    let workspaces = await this.api.listWorkspaces();
    let workspaceNames = workspaces.map((w) => w.name).sort();
    let storedWorkspaceName = this.session.data.workspace;

    if (storedWorkspaceName && workspaceNames.includes(storedWorkspaceName)) {
      this.transitionTo('workspace', storedWorkspaceName);
      return;
    }

    if (workspaceNames.includes('default')) {
      this.transitionTo('workspace', 'default');
      return;
    }

    if (workspaceNames.length) {
      this.transitionTo('workspace', workspaceNames[0]);
      return;
    }

    this.transitionTo('workspace', 'default');
  }
}
