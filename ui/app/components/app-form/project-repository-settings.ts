import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import RouterService from '@ember/routing/router-service';
import { isEmpty } from '@ember/utils';
import ApiService from 'waypoint/services/api';
import FlashMessagesService from 'waypoint/services/flash-messages';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { Project } from 'waypoint-pb';
import parseUrl from 'parse-url';

const FORMAT = {
  HCL: 0,
  JSON: 1,
};

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
  variablesList: [],
};

interface ProjectSettingsArgs {
  project: Project.AsObject;
}

export default class AppFormProjectRepositorySettings extends Component<ProjectSettingsArgs> {
  // normal class body definition here
  @service api!: ApiService;
  @service flashMessages!: FlashMessagesService;
  @service router!: RouterService;
  @tracked project: Project.AsObject;
  @tracked authCase: number;
  @tracked serverHcl: boolean;

  constructor(owner: any, args: any) {
    super(owner, args);
    this.project = JSON.parse(JSON.stringify(DEFAULT_PROJECT_MODEL)) as Project.AsObject; // to ensure we're doing a deep copy
    let { project } = this.args;
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

  get decodedWaypointHcl(): string {
    return atob((this.project.waypointHcl as string) || '');
  }

  get decodedPrivateKey(): string {
    return atob((this.git?.ssh?.privateKeyPem as string) || '');
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
        currentModel[key] = DEFAULT_PROJECT_MODEL[key];
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

    let value = e.target.value;
    // if private key, encode input to base 64
    if (path === 'privateKeyPem') {
      value = btoa(value);
    }

    this.project.dataSource.git.ssh[path] = value;
  }

  @action
  setWaypointHcl(e: any) {
    this.project.waypointHcl = btoa(e.target.value);
  }

  @action
  async saveSettings(e: Event) {
    e.preventDefault();
    if (!this.validateGitUrl()) {
      return;
    }

    if (!this.serverHcl) {
      this.project.waypointHcl = '';
    }

    try {
      await this.api.upsertProject(this.project, this.authCase);
      this.flashMessages.success('Settings saved');
      this.router.transitionTo('workspace.projects.project', this.project.name);
    } catch (err) {
      this.flashMessages.error('Failed to save Settings', { content: err.message, sticky: true });
    }
  }
}
