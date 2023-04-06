/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { ConfigGetRequest, ConfigSetRequest, ConfigVar, Project, Ref } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import { BufferedChangeset } from 'validated-changeset';
import Component from '@glimmer/component';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface ProjectConfigArgs {
  variablesList: ConfigVar.AsObject[];
  project: Project.AsObject;
}

type ConfigVarChangeset = Partial<BufferedChangeset> & ConfigVar.AsObject;

export default class ProjectConfigVariablesListComponent extends Component<ProjectConfigArgs> {
  @service api!: ApiService;
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;

  @tracked variablesList: Array<ConfigVar.AsObject>;
  @tracked project: Project.AsObject;
  @tracked isCreating: boolean;
  @tracked activeVariable;

  constructor(owner: unknown, args: ProjectConfigArgs) {
    super(owner, args);
    let { variablesList, project } = args;
    this.variablesList = variablesList;
    this.project = project;
    this.activeVariable = null;
    this.isCreating = false;
  }

  @action
  cancelCreate(): void {
    this.variablesList.splice(0, 1);
    this.variablesList = [...this.variablesList];
    this.isCreating = false;
  }

  @action
  async deleteVariable(variable: ConfigVar.AsObject): Promise<void> {
    await this.saveVariableSettings(variable, true);
  }

  @action
  addVariable(): void {
    this.isCreating = true;
    let newVar = new ConfigVar();
    let newVarObj = newVar.toObject();
    this.variablesList = [newVarObj, ...this.variablesList];
    this.activeVariable = newVarObj;
  }

  @action
  async saveVariableSettings(variable: ConfigVarChangeset, deleteVariable?: boolean): Promise<void> {
    let req = new ConfigSetRequest();

    let projectRef = new Ref.Project();
    projectRef.setProject(this.project.name);

    let newVar = new ConfigVar();

    newVar.setProject(projectRef);

    newVar.setName(variable.name);

    if (variable?.pb_static) {
      newVar.setStatic(variable.pb_static);
    }

    if (variable.internal) {
      newVar.setInternal(variable.internal);
    }

    if (variable.nameIsPath) {
      newVar.setNameIsPath(variable.nameIsPath);
    }

    if (deleteVariable) {
      newVar.setUnset(new Empty());
    }

    req.setVariablesList([newVar]);
    try {
      await this.api.client.setConfig(req, this.api.WithMeta());
      // If Config Var was saved correctly, refetch Config
      let configRef = new Ref.Project();
      configRef.setProject(this.project.name);
      let configReq = new ConfigGetRequest();
      configReq.setProject(configRef);

      let config = await this.api.client.getConfig(configReq, this.api.WithMeta());

      this.variablesList = config?.toObject().variablesList;
      this.activeVariable = null;
      this.isCreating = false;
    } catch (err) {
      this.flashMessages.error('Failed to save Variable', { content: err.message, sticky: true });
    }
  }
}
