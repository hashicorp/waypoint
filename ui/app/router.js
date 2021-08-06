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
      this.route('project');
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
          this.route('build-id', { path: '/build/:build_id' });
          this.route('build', { path: '/build/seq/:sequence' });
          this.route('deployments');
          this.route('deployment-id', { path: '/deployment/:deployment_id' });
          this.route('deployment', { path: '/deployment/seq/:sequence' });
          this.route('releases');
          this.route('release-id', { path: '/release/:release_id' });
          this.route('release', { path: '/release/seq/:sequence' });
          this.route('logs');
          this.route('exec');
        });
        this.route('settings', function () {
          this.route('repository', { path: '/' });
          this.route('variables');
        });
      });
      this.route('new');
    });
  });
});
