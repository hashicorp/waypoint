import Component from '@glimmer/component';
import RouterService from '@ember/routing/router-service';
import { TaskGenerator, task } from 'ember-concurrency';
import { taskFor } from 'ember-concurrency-ts';
import { inject as service } from '@ember/service';
import { Workspace } from 'waypoint-pb';
import ApiService from 'waypoint/services/api';

type RouteInfo = RouterService['currentRoute'];

export default class extends Component {
  @service router!: RouterService;
  @service api!: ApiService;

  constructor(...args: ConstructorParameters<typeof Component>) {
    super(...args);

    taskFor(this.loadWorkspaces).perform();
  }

  @task({ restartable: true })
  *loadWorkspaces(): TaskGenerator<Workspace.AsObject[]> {
    return yield this.api.listWorkspaces();
  }

  get workspaces(): Workspace.AsObject[] {
    return taskFor(this.loadWorkspaces).lastSuccessful?.value ?? [];
  }

  get workspaceRoute(): RouteInfo | undefined {
    return this.router.currentRoute.find((route) => route.name === 'workspace');
  }

  get workspace(): string | undefined {
    return this.workspaceRoute?.params?.workspace_id;
  }

  get projectRoute(): RouteInfo | undefined {
    return this.router.currentRoute.find((route) => route.name === 'workspace.projects.project');
  }

  get project(): string | undefined {
    return this.projectRoute?.params?.project_id;
  }

  get appRoute(): RouteInfo | undefined {
    return this.router.currentRoute.find((route) => route.name === 'workspace.projects.project.app');
  }

  get app(): string | undefined {
    return this.appRoute?.params?.app_id;
  }
}
