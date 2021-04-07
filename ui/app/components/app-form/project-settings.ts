import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import RouterService from '@ember/routing/router-service';
import ApiService from 'waypoint/services/api';
import FlashMessagesService from 'waypoint/services/flash-messages';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { Project, UpsertProjectRequest, Job, Application } from 'waypoint-pb';
import parseUrl from 'parse-url';

const FORMAT = {
  HCL: 0,
  JSON: 1,
};
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
        password: string,
      },
      ssh: {
        user: string,
        password: string,
        privateKeyPem: string
      }
    },
    local:any,
  };
  dataSourcePoll: {
    enabled: boolean,
    interval: string,
  };
  remoteEnabled: boolean;
  waypointHcl: string;
  waypointHclFormat: number;
}

const DEFAULT_PROJECT_MODEL = {
  name: '',
  applicationsList: [],
  dataSource: {
    git: {
      url: '',
      path: '',
      ref: 'HEAD',
      basic: {
        username: '',
        password: ''
      },
      ssh: {
        user: '',
        password: '',
        privateKeyPem: '',
      }
    },
    local: null,
  },
  dataSourcePoll: {
    enabled: false,
    interval: '2m',
  },
  remoteEnabled: null,
  waypointHcl: '',
  waypointHclFormat: FORMAT.HCL,
};

interface ProjectSettingsArgs {
  project: ProjectModel;
}

export default class AppFormProjectSettings extends Component<ProjectSettingsArgs> {
  // normal class body definition here
  @service api!: ApiService;
  @service flashMessages!: FlashMessagesService;
  @service router!: RouterService;
  @tracked project: ProjectModel;
  @tracked authCase: number;
  @tracked serverHcl: boolean;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { project } = this.args;
    this.project = project;
    // restore git settings state from existing datasource
    if (this.project?.dataSource?.git && this.project?.dataSource?.git?.url) {
      // Set authCase if it exists
      if (project.dataSource?.git?.ssh?.privateKeyPem) {
        this.authCase = 5;
      } else {
        if (project.dataSource?.git?.basic?.username) {
          this.authCase = 4;
        } else {
          this.authCase = 0;
        }
        // Set empty default git auth if not defined
        if (!this.project?.dataSource?.git?.basic?.username || !this.project?.dataSource?.git?.ssh?.privateKeyPem) {
          this.project.dataSource.git.basic = DEFAULT_PROJECT_MODEL.dataSource.git.basic;
          this.project.dataSource.git.ssh = DEFAULT_PROJECT_MODEL.dataSource.git.ssh;
        }
      }
    } else {
      // set empty default dataSource data if non-existent
      this.project.dataSource = DEFAULT_PROJECT_MODEL.dataSource;
      this.authCase = 4;
    }

    if (!this.project?.dataSourcePoll) {
      this.project.dataSourcePoll = DEFAULT_PROJECT_MODEL.dataSourcePoll;
    }

    if (this.project?.waypointHcl) {
      this.serverHcl = true;
    } else {
      this.serverHcl = false;
      this.project.waypointHcl = DEFAULT_PROJECT_MODEL.waypointHcl;
      this.project.waypointHclFormat = DEFAULT_PROJECT_MODEL.waypointHclFormat;
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

  get authNotSet() {
    return this.authCase === 0;
  }

  get git() {
    return this.dataSource?.git;
  }


  validateGitUrl() {
    let gitUrl = parseUrl(this.project.dataSource.git.url);
    // If basic auth, match https url
    if (this.authCase == 4 || this.authCase == 0) {
      if (gitUrl.protocol != 'https') {
        this.flashMessages.error('Git url needs to use "https:" protocol');
        return false;
      }
    }
    // If ssh force users to use a git: url
    if (this.authCase == 5) {
      if (gitUrl.protocol != 'ssh') {
        this.flashMessages.error('Git url needs to use "git:" protocol');
        return false;
      }
    }
    return true;
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
      this.project.dataSource.git.ssh = {
        user: '',
        password: '',
        privateKeyPem: '',
      };
    }
    this.project.dataSource.git.ssh[path] = e.target.value;
  }

  @action
  async saveSettings(e: Event) {
    e.preventDefault();
    if (!this.validateGitUrl()) {
      return;
    }
    let project = this.project;
    project.dataSource = this.dataSource;
    let ref = new Project();
    ref.setName(project.name);
    let dataSource = new Job.DataSource();
    let dataSourcePoll = new Project.Poll();
    dataSourcePoll.setEnabled(project.dataSourcePoll.enabled);
    dataSourcePoll.setInterval(project.dataSourcePoll.interval);
    // Git settings
    let git = new Job.Git();
    git.setUrl(project.dataSource.git.url);
    git.setPath(project.dataSource.git.path);
    if (!project.dataSource.git.ref) {
      git.setRef('HEAD');
    } else {
      git.setRef(project.dataSource.git.ref);
    }

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
      let encoder = new window.TextEncoder();
      gitSSH.setPrivateKeyPem(encoder.encode(this.git.ssh.privateKeyPem));
      gitSSH.setUser(this.git.ssh.user);
      gitSSH.setPassword(this.git.ssh.password);
      git.setSsh(gitSSH);
      git.clearBasic();
    }

    if (this.authNotSet) {
      git.clearBasic();
      git.clearSsh();
    }

    dataSource.setGit(git);
    ref.setDataSource(dataSource);
    ref.setDataSourcePoll(dataSourcePoll);

    if (this.serverHcl && this.project.waypointHcl) {
      let hclEncoder = new window.TextEncoder();
      let waypointHcl = hclEncoder.encode(this.project.waypointHcl);
      // Hardcode hcl for now
      ref.setWaypointHclFormat(FORMAT.HCL);
      ref.setWaypointHcl(waypointHcl);
    }
    let applist = project.applicationsList.map((app: Application.AsObject) => {
      return new Application(app);
    });
    ref.setApplicationsList(applist);

    // Build and trigger request
    let req = new UpsertProjectRequest();
    req.setProject(ref);
    try {
      await this.api.client.upsertProject(req, this.api.WithMeta());
      this.flashMessages.success('Settings saved');
      this.router.transitionTo('workspace.projects.project', this.project.name);
    } catch (err) {
      this.flashMessages.error('Failed to save Settings', { content: err.message, sticky: true });
    }
  }
}
