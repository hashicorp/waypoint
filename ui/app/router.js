import EmberRouter from '@ember/routing/router';
import config from './config/environment';

export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  this.route('workspace', { path: '/:workspace_id' }, function () {
    this.route('project', { path: '/:project_id' }, function () {
      this.route('app', { path: '/app/:app_id' }, function () {
        this.route('builds');
      });
    });
  });
});
