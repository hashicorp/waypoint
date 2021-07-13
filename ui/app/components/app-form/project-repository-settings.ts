import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import RouterService from '@ember/routing/router-service';
import { isEmpty } from '@ember/utils';
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
      url: string;
      path: string;
      ref: string;
      basic: {
        username: string;
        password: string;
      };
      ssh: {
        user: string;
        password: string;
        privateKeyPem: string;
      };
    };
    local: any;
  };
  dataSourcePoll: {
    enabled: boolean;
    interval: string;
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
        password: '',
      },
      ssh: {
        user: '',
        password: '',
        privateKeyPem: '',
      },
    },
    local: null,
  },
  dataSourcePoll: {
    enabled: false,
    interval: '2m',
  },
  remoteEnabled: false,
  waypointHcl: '',
  waypointHclFormat: FORMAT.HCL,
};

interface ProjectSettingsArgs {
  project: ProjectModel;
}

export default class AppFormProjectRepositorySettings extends Component<ProjectSettingsArgs> {
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
    this.project = JSON.parse(JSON.stringify(DEFAULT_PROJECT_MODEL)); // to ensure we're doing a deep copy
    this.populateExistingFields(project, this.project);
    this.authCase = 4;
    this.serverHcl = !!this.project?.waypointHcl;

    // restore git settings state if editing existing project
    if (this.git?.url) {
      if (project.dataSource?.git?.ssh?.privateKeyPem) {
        this.setAuthCase(5);
        return;
      }
      if (!project.dataSource?.git?.basic?.username) {
        this.setAuthCase(0);
        return;
      }
    }
  }

  get dataSource() {
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
    if (this.authBasic || this.authNotSet) {
      if (gitUrl.protocol !== 'https') {
        this.flashMessages.error('Git url needs to use "https:" protocol');
        return false;
      }
    }
    // If ssh force users to use a git: url
    if (this.authSSH) {
      if (gitUrl.protocol !== 'ssh') {
        this.flashMessages.error('Git url needs to use "git:" protocol');
        return false;
      }
    }
    return true;
  }

  populateExistingFields(projectFromArgs, currentModel) {
    for (let [key, value] of Object.entries(projectFromArgs)) {
      if (isEmpty(value)) {
        continue;
      }

      // if the value is a nested object, recursively call this function
      if (typeof value === 'object') {
        this.populateExistingFields(
          projectFromArgs[key],
          !isEmpty(currentModel[key]) ? currentModel[key] : {}
        );
      }

      if (value !== currentModel[key]) {
        currentModel[key] = value;
      }
    }
  }

  @action
  setAuthCase(val: any) {
    this.authCase = val;
  }

  @action
  setBasicAuth(path: string, e: any) {
    if (!this.git?.basic) {
      this.git.basic = {
        username: '',
        password: '',
      };
    }
    this.git.basic[path] = e.target.value;
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
