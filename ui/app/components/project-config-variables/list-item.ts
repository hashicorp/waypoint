import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { ConfigVar, Project } from 'waypoint-pb';

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

  storeInitialVariable() {
    this.initialVariable = JSON.parse(JSON.stringify(this.variable));
  }

  @action
  async deleteVariable(variable) {
    // await this.args.deleteVariable(variable);
  }

  @action
  editVariable() {
    this.isEditing = true;
    this.storeInitialVariable();
  }

  @action
  async saveVariable(e) {
    e.preventDefault();
    // todo
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
