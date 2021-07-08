import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { Project, UpsertProjectRequest, Variable } from 'waypoint-pb';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import RouterService from '@ember/routing/router-service';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

interface ProjectSettingsArgs {
  project: Project.AsObject;
}

export default class ProjectInputVariablesListComponent extends Component<ProjectSettingsArgs> {
  @service api!: ApiService;
  @service router!: RouterService;
  @service flashMessages;
  @tracked project;
  @tracked variablesList: Array<Variable.AsObject>;
  @tracked isCreating: boolean;
  @tracked activeVariable;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { project } = args;
    this.project = project;
    this.variablesList = this.project.variablesList;
    this.activeVariable = null;
    this.isCreating = false;
  }

  @action
  addVariable() {
    this.isCreating = true;
    let newVar = new Variable();
    // Seems like setServer (with empty arguments?) is required to make it a server variable
    newVar.setServer();
    let newVarObj = newVar.toObject();
    this.variablesList = [newVarObj, ...this.variablesList];
    this.activeVariable = newVarObj;
  }

  @action
  async deleteVariable(variable) {
    this.variablesList = this.variablesList.filter((v) => v.name !== variable.name);
    await this.saveVariableSettings();
  }

  @action
  cancelCreate() {
    this.activeVariable = null;
    this.isCreating = false;
  }

  @action
  async saveVariableSettings(
    variable?: Variable.AsObject,
    initialVariable?: Variable.AsObject
  ): Promise<Project.AsObject | void> {
    let project = this.project;
    let ref = new Project();
    ref.setName(project.name);
    if (variable && initialVariable) {
      let existingVarIndex = this.variablesList.findIndex((v) => v.name === initialVariable.name);
      if (existingVarIndex !== -1) {
        this.variablesList.splice(existingVarIndex, 1, variable);
        this.variablesList = [...this.variablesList];
      }
    }
    let varProtosList = this.variablesList.map((v: Variable.AsObject) => {
      let variable = new Variable();
      variable.setName(v.name);
      variable.setServer(new Empty());
      if (v.hcl) {
        variable.setHcl(v.hcl);
      } else {
        variable.setStr(v.str);
      }
      return variable;
    });
    ref.setVariablesList(varProtosList);
    // Build and trigger request
    let req = new UpsertProjectRequest();
    req.setProject(ref);
    try {
      let resp = await this.api.client.upsertProject(req, this.api.WithMeta());
      let respProject = resp.toObject().project;
      this.project = respProject;
      this.flashMessages.success('Settings saved');
      this.activeVariable = null;
      this.isCreating = false;
      this.router.refresh('workspace.projects.project');
      return respProject;
    } catch (err) {
      this.flashMessages.error('Failed to save Settings', { content: err.message, sticky: true });
    }
  }
}
