import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { ConfigSetRequest, ConfigVar, Project } from 'waypoint-pb';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';

interface VariableArgs {
  variable: ConfigVar.AsObject;
  isEditing: boolean;
  isCreating: boolean;
  saveVariableSettings: (
    variable: ConfigVar.AsObject,
    initialVariable: ConfigVar.AsObject
  ) => Promise<Project.AsObject>;
  deleteVariable: (variable: ConfigVar.AsObject) => Promise<void>;
}

export default class ProjectConfigVariablesListItemComponent extends Component<VariableArgs> {
  @service api!: ApiService;
  @service flashMessages;

  initialVariable: ConfigVar.AsObject;
  @tracked variable: ConfigVar.AsObject;
  @tracked isCreating: boolean;
  @tracked isEditing: boolean;

  constructor(owner: any, args: VariableArgs) {
    super(owner, args);
    let { variable, isEditing, isCreating } = args;
    this.variable = variable;
    this.isEditing = isEditing;
    this.isCreating = isCreating;
    this.storeInitialVariable();
  }

  get isEditable(): boolean {
    // Temporarily making Dynamic vars uneditable
    return !this.variable.dynamic;
  }

  storeInitialVariable() {
    this.initialVariable = JSON.parse(JSON.stringify(this.variable));
  }

  @action
  async deleteVariable(variable) {
    await this.args.deleteVariable(variable);
  }

  @action
  editVariable() {
    this.isEditing = true;
    this.storeInitialVariable();
  }

  @action
  async saveVariable(e) {
    e.preventDefault();
    if (this.variable.name === '' || this.variable.pb_static === '') {
      return this.flashMessages.error('Variable keys or values can not be empty');
    }
    let req = new ConfigSetRequest();
    let savedVars = await this.args.saveVariableSettings(this.variable, this.initialVariable);
    this.isCreating = false;
    this.isEditing = false;
  }

  @action
  cancelCreate() {
    this.isCreating = false;
    this.isEditing = false;
    this.deleteVariable(this.variable);
  }

  @action
  cancelEdit() {
    this.isCreating = false;
    this.isEditing = false;
    this.variable = this.initialVariable;
  }

}
