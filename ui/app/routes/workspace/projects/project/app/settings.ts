import Route from '@ember/routing/route';

export default class WorkspaceProjectsProjectAppSettings extends Route {
  breadcrumbs(model: AppRouteModel) {
    if (!model) return [];
    return [
      {
        label: model.application.application,
        icon: 'git-repository',
        args: ['workspace.projects.project.app'],
      },
      {
        label: 'Settings',
        icon: 'settings',
        args: ['workspace.projects.project.app.settings'],
      },
    ];
  }
}
