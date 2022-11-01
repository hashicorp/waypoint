import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { Project, Variable } from 'waypoint-pb';
import { inject as service } from '@ember/service';

interface VariableArgs {
  variable: Variable.AsObject;
  isEditing: boolean;
  isCreating: boolean;
  saveVariableSettings: (
    variable: Variable.AsObject,
    initialVariable: Variable.AsObject
  ) => Promise<Project.AsObject>;
  deleteVariable: (variable: Variable.AsObject) => Promise<void>;
}

export default class ProjectInputVariablesListItemComponent extends Component<VariableArgs> {
  @service('pdsFlashMessages') flashMessages;

  initialVariable: Variable.AsObject;
  @tracked variable: Variable.AsObject;
  @tracked isCreating: boolean;
  @tracked isEditing: boolean;
  @tracked writeOnly: boolean;

  constructor(owner: unknown, args: VariableArgs) {
    super(owner, args);
    let { variable, isEditing, isCreating } = args;
    this.variable = variable;
    this.isEditing = isEditing;
    this.isCreating = isCreating;
    this.writeOnly = false;
    this.initialVariable = JSON.parse(JSON.stringify(this.variable));
  }

  get isHcl(): boolean {
    return !!this.variable.hcl;
  }

  get isSensitive(): boolean {
    return !!this.variable.sensitive;
  }

  storeInitialVariable(): void {
    this.initialVariable = JSON.parse(JSON.stringify(this.variable));
    this.writeOnly = this.variable.sensitive ? true : false;
  }

  @action
  async deleteVariable(variable: Variable.AsObject): Promise<void> {
    await this.args.deleteVariable(variable);
  }

  @action
  editVariable(): void {
    this.isEditing = true;
    this.storeInitialVariable();
  }

  @action
  async saveVariable(e: Event): Promise<void> {
    e.preventDefault();
    // Validate non-empty var name & value
    if (
      this.variable.name === '' ||
      (this.isHcl && this.variable.hcl === '') ||
      (!this.isHcl && this.variable.str === '')
    ) {
      return this.flashMessages.error('Variable keys or values can not be empty');
    }
    let savedProject = await this.args.saveVariableSettings(this.variable, this.initialVariable);
    if (savedProject) {
      let newVariable = savedProject.variablesList.find((v) => v.name === this.variable.name);
      if (newVariable) {
        this.variable = newVariable;
      }
    }
    this.isCreating = false;
    this.isEditing = false;
  }

  @action
  cancelCreate(): void {
    this.isCreating = false;
    this.isEditing = false;
    this.deleteVariable(this.variable);
  }

  @action
  cancelEdit(): void {
    this.isCreating = false;
    this.isEditing = false;
    this.variable = this.initialVariable;
  }

  @action
  toggleHcl(variable: Variable.AsObject): void {
    if (this.isHcl) {
      this.variable.str = variable.hcl;
      this.variable.hcl = '';
    } else {
      this.variable.hcl = variable.str;
      this.variable.str = '';
    }
  }

  @action
  toggleSensitive(variable: Variable.AsObject): void {
    this.variable.sensitive = !variable.sensitive;
  }
}
