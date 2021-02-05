import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import {Project, UpsertProjectRequest, Job, Application } from 'waypoint-pb';

class ProjectModel {
  name: string;
  applicationsList: [];
  dataSource: {
    git: any,
    local:any,
    url: any,
  };
  remoteEnabled: boolean;
}

interface ProjectSettingsArgs {
  project: ProjectModel
}

export default class AppFormProjectSettings extends Component<ProjectSettingsArgs> {
  // normal class body definition here
  @service api!: ApiService;
  @tracked project: ProjectModel;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { project } = this.args;
    this.project = project;
  }

  get dataSource(){
    return this.project.dataSource;
  }


  get git() {
    return this.dataSource?.git || {};
  }

  set git(args) {
    this.project.dataSource.git = args
    debugger;
  }

  setGitData(prop: string, value: any) {
    if (!this.dataSource.git) {
      this.dataSource.git = {};
    }
    this.dataSource.git[prop] = value;
  }

  @action
  setGitPath(e: any) {
    this.setGitData('path', e.target.value)
  }

  @action
  setGitUrl(e: any) {
    this.setGitData('url', e.target.value)
  }

  @action
  setGitRef(e: any) {
    this.setGitData('ref', e.target.value)
  }

  @action
  async saveSettings() {
    let project = this.project;
    project.dataSource = this.dataSource;
    let ref = new Project();
    ref.setName(project.name);
    let dataSource = new Job.DataSource();
    let git = new Job.Git();
    git.setUrl(project.dataSource.git.url);
    git.setPath(project.dataSource.git.path);
    git.setRef(project.dataSource.git.ref);
    dataSource.setGit(git);
    ref.setDataSource(dataSource);
    ref.setRemoteEnabled(project.remoteEnabled);
    const applist = project.applicationsList.map((app: any) => {
      return new Application(app);
    });
    ref.setApplicationsList(applist);
    let req = new UpsertProjectRequest();
    req.setProject(ref);
    await this.api.client.upsertProject(req, this.api.WithMeta());
  }
}
