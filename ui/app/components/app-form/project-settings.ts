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

const DEFAULT_PROJECT_MODEL = {
  name: '',
  applicationsList: [],
  dataSource: {
    git: {
      url: '',
      path: '',
      ref: '',
      basic: {
        username: '',
        password: ''
      },
      ssh: {
        privateKeyPem: ''
      }
    },
    local: null,
  },
  remoteEnabled: null
};

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
    this.project = Object.assign(DEFAULT_PROJECT_MODEL, project);
    if (this.project?.dataSource?.git) {
      this.authCase = project.dataSource?.git?.ssh?.privateKeyPem ? 5 : 4;
    } else {
      this.authCase = 4;
    }
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


  @action
  setAuthCase(val:any) {
    this.authCase = val;
  }

  @action
  setBasicAuth(path: string, e: any) {
    if (!this.project.dataSource?.git?.basic) {
      this.project.dataSource.git.basic = {
        username: '',
        password: ''
      };
    }
    this.project.dataSource.git.basic[path] = e.target.value;
  }

  @action
  setSshAuth(path: string, e: any) {
    if (!this.project.dataSource?.git?.ssh) {
      this.project.dataSource.git.basic = {
        privateKeyPem: '',
      };
    }
    this.project.dataSource.git.ssh[path] = e.target.value;
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
      git.clearSsh();
    }

    if (this.authSSH) {
      let gitSSH = new Job.Git.SSH();
      gitSSH.setPrivateKeyPem(this.git.ssh.privateKeyPem);
      git.setSsh(gitSSH);
      git.clearBasic();
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
