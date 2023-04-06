/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import RouterService from '@ember/routing/router-service';
import { isEmpty } from '@ember/utils';
import ApiService from 'waypoint/services/api';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { Project, Job } from 'waypoint-pb';
import parseUrl from 'parse-url';
import { later } from '@ember/runloop';

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
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;
  @service router!: RouterService;
  @tracked project: Project.AsObject;
  @tracked authCase: number;
  @tracked serverHcl: boolean;
  defaultProject: Project.AsObject;

  constructor(owner: unknown, args: ProjectSettingsArgs) {
    super(owner, args);
    this.defaultProject = JSON.parse(JSON.stringify(DEFAULT_PROJECT_MODEL)) as Project.AsObject; // to ensure we're doing a deep copy
    this.project = this.defaultProject;
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

  get dataSource(): Job.DataSource.AsObject | undefined {
    return this.project.dataSource;
  }

  get authSSH(): boolean {
    return this.authCase === 5;
  }

  get authBasic(): boolean {
    return this.authCase === 4;
  }

  get authNotSet(): boolean {
    return this.authCase === 0;
  }

  get git(): Job.Git.AsObject | undefined {
    return this.dataSource?.git;
  }

  get decodedWaypointHcl(): string {
    return atob((this.project.waypointHcl as string) || '');
  }

  get decodedPrivateKey(): string {
    return atob((this.git?.ssh?.privateKeyPem as string) || '');
  }

  validateGitUrl(): boolean {
    let gitUrl = parseUrl(this.project.dataSource?.git?.url);
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
        this.flashMessages.error('Git url needs to use "ssh:" protocol');
        return false;
      }
    }
    return true;
  }

  populateExistingFields(projectFromArgs: Project.AsObject, currentModel: Project.AsObject): void {
    for (let [key, value] of Object.entries(projectFromArgs)) {
      // Guard against prototype pollution
      if (!Object.prototype.hasOwnProperty.call(currentModel, key)) {
        continue;
      }

      if (isEmpty(value)) {
        currentModel[key] = this.defaultProject[key];
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
  setAuthCase(val: number): void {
    this.authCase = val;
  }

  @action
  setBasicAuth(path: string, e: Event): void {
    if (!this.git) {
      return;
    }

    if (!isFormField(e.target)) {
      return;
    }

    if (!this.git.basic) {
      this.git.basic = {
        username: '',
        password: '',
      };
    }

    this.git.basic[path] = e.target.value;
  }

  @action
  setSshAuth(path: string, e: Event): void {
    if (!this.git) {
      return;
    }

    if (!isFormField(e.target)) {
      return;
    }

    if (!this.git.ssh) {
      this.git.ssh = {
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

    this.git.ssh[path] = value;
  }

  @action
  setWaypointHcl(value: string): void {
    this.project.waypointHcl = btoa(value);
  }

  @action
  async saveSettings(e: Event): Promise<void> {
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

      // Refresh project route to get the latest state of the InitOp (if any)
      this.router.refresh('workspace.projects.project');

      if (!Ember.testing) {
        // Optimistically refresh again a few seconds later, by which time
        // the InitOp is likely to have completed
        later(this.router, 'refresh', 'workspace.projects.project', 3000);
      }

      this.router.transitionTo('workspace.projects.project', this.project.name);
    } catch (err) {
      this.flashMessages.error('Failed to save Settings', { content: err.message, sticky: true });
    }
  }
}

interface FormField {
  value: string;
}

function isFormField(element: unknown): element is FormField {
  return typeof (element as FormField).value === 'string';
}
