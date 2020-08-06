import EmberRouter from '@ember/routing/router';
import config from './config/environment';

export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  this.route('auth', { path: '/auth' });
  this.route('workspaces', { path: '/' });
  this.route('workspace', { path: '/:workspace_id' }, function () {
    this.route('projects', { path: '/projects' }, function () {
      this.route('project', { path: '/:project_id' }, function () {
        this.route('apps');
        this.route('app', { path: '/app/:app_id' }, function () {
          this.route('builds');
        });
      });
    });
  });
});
