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
    git: {
      url: string,
      path: string,
      ref: string,
      basic: {
        username: string,
        password: string
      },
      ssh: {
        privateKeyPem: string
      }
    },
    local:any,
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
  @tracked authCase: number;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { project } = this.args;
    this.project = Object.assign(new ProjectModel(), project);
    this.authCase = 4;
  }

  get dataSource(){
    return this.project.dataSource;
  }

  get authSSH() {
    return this.authCase === 5;
  }

  get authBasic() {
    return this.authCase === 4;
  }

  get git() {
    return this.dataSource?.git;
  }

  set git(args: any) {
    this.project.dataSource.git = args;
  }

  setGitData(prop: string, value: any) {
    if (!this.dataSource) {
      this.project.dataSource = {};
    }
    if (!this.dataSource.git) {
      this.project.dataSource.git = {};
    }
    this.project.dataSource.git[prop] = value;
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
  setAuthCase(val:any) {
    this.authCase = val;
  }

  @action
  setGitSSH(e: any) {
    if (!this.git.ssh) {
      this.git.ssh = {};
    }
    this.git.ssh['privateKeyPem'] = e.target.value;
  }

  @action
  setGitPassword(e: any) {
    if (!this.git.basic) {
      this.git.basic = {};
    }
    this.git.basic['password'] = e.target.value;
  }

  @action
  setGitUsername(e: any) {
    if (!this.git.basic) {
      this.git.basic = {};
    }
    this.git.basic['username'] = e.target.value
  }

  @action
  async saveSettings() {
    let project = this.project;
    project.dataSource = this.dataSource;
    let ref = new Project();
    ref.setName(project.name);
    let dataSource = new Job.DataSource();
    // Git settings
    let git = new Job.Git();
    git.setUrl(project.dataSource.git.url);
    git.setPath(project.dataSource.git.path);
    git.setRef(project.dataSource.git.ref);

    // Git Authentication settings
    if (this.authBasic) {
      let gitBasic = new Job.Git.Basic();
      gitBasic.setUsername(this.git.basic.username);
      gitBasic.setPassword(this.git.basic.password);
      git.setBasic(gitBasic);
    }

    if (this.authSSH) {
      let gitSSH = new Job.Git.SSH();
      gitSSH.setPrivateKeyPem(this.git.ssh.privateKeyPem);
      git.setSsh(gitSSH);
    }

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
