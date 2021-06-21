import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { Variable } from 'waypoint-pb';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';

interface VariableArgs {
  variable: Variable.AsObject;
  isEditing: boolean;
  isCreating: boolean;
  saveVariableSettings: (variable: Variable.AsObject) => Promise<void>;
  deleteVariable: (variable: Variable.AsObject) => Promise<void>;
}

export default class ProjectInputVariablesListComponent extends Component<VariableArgs> {
  initialVariable?: Variable.AsObject;
  @service api!: ApiService;
  @tracked variable: Variable.AsObject;
  @tracked isCreating: boolean;
  @tracked isEditing: boolean;
  @tracked isHcl: boolean;

  constructor(owner: any, args: VariableArgs) {
    super(owner, args);
    this.isCreating = false;
    let { variable, isEditing, isCreating } = args;
    this.variable = variable;
    this.isEditing = isEditing;
    this.isCreating = isCreating;
    this.isHcl = false;
    if (variable.hcl) {
      this.isHcl = true;
    }
    if (this.isEditing) {
      this.initialVariable = JSON.parse(JSON.stringify(this.variable));
    }
  }

  @action
  async deleteVariable(variable) {
    await this.args.deleteVariable(variable);
  }

  @action
  editVariable() {
    this.isEditing = true;
    this.initialVariable = JSON.parse(JSON.stringify(this.variable));
  }

  @action
  async saveVariable(e) {
    e.preventDefault();
    await this.args.saveVariableSettings(this.variable);
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

  @action
  toggleHcl(variable) {
    if (this.isHcl) {
      this.isHcl = false;
      this.variable.str = variable.hcl;
      this.variable.hcl = '';
    } else {
      this.isHcl = true;
      this.variable.hcl = variable.str;
      this.variable.str = '';
    }
  }
}
