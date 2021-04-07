import EmberRouter from '@ember/routing/router';
import config from 'waypoint/config/environment';

export default class Router extends EmberRouter {
  location = config.locationType;
  rootURL = config.rootURL;
}

Router.map(function () {
  this.route('auth', function () {
    this.route('invite');
    this.route('token');
  });
  this.route('onboarding', function () {
    this.route('install', function () {
      this.route('manual');
      this.route('homebrew');
      this.route('chocolatey');
      this.route('linux', function () {
        this.route('ubuntu');
        this.route('centos');
        this.route('fedora');
        this.route('amazon');
      });
    });
    this.route('connect');
    this.route('start');
  });
  this.route('workspaces', { path: '/' }, function () {
    this.route('projects', function () {
      this.route('project', function () {});
    });
  });
  this.route('workspace', { path: '/:workspace_id' }, function () {
    this.route('projects', { path: '/' }, function () {
      this.route('project', { path: '/:project_id' }, function () {
        this.route('apps', function () {
          this.route('new');
        });
        this.route('app', { path: '/app/:app_id' }, function () {
          this.route('builds');
          this.route('build', { path: '/build/:build_id' });
          this.route('deployments');
          this.route('deployment', { path: '/deployment/:deployment_id' });
          this.route('releases');
          this.route('release', { path: '/release/:release_id' });
          this.route('logs');
          this.route('exec');
        });
        this.route('settings');
      });
      this.route('new');
    });
  });
});
