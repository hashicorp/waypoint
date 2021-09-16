import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { ConfigVar, Project } from 'waypoint-pb';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';

interface VariableArgs {
  variable: ConfigVar.AsObject;
  isEditing: boolean;
  isCreating: boolean;
  saveVariableSettings: (variable: ConfigVar.AsObject, deleteVariable?: boolean) => Promise<Project.AsObject>;
  deleteVariable: (variable: ConfigVar.AsObject) => Promise<void>;
  cancelCreate: () => void;
}

export default class ProjectConfigVariablesListItemComponent extends Component<VariableArgs> {
  @service api!: ApiService;
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;

  initialVariable!: ConfigVar.AsObject;
  @tracked variable: ConfigVar.AsObject;
  @tracked isCreating: boolean;
  @tracked isEditing: boolean;

  constructor(owner: unknown, args: VariableArgs) {
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

  storeInitialVariable(): void {
    this.initialVariable = JSON.parse(JSON.stringify(this.variable));
  }

  @action
  async deleteVariable(variable: ConfigVar.AsObject): Promise<void> {
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
    if (this.variable.name === '' || this.variable.pb_static === '') {
      this.flashMessages.error('Variable keys or values can not be empty');
      return;
    }
    await this.args.saveVariableSettings(this.variable, false);
    this.isCreating = false;
    this.isEditing = false;
  }

  @action
  cancelCreate(): void {
    this.isCreating = false;
    this.isEditing = false;
    this.args.cancelCreate();
  }

  @action
  cancelEdit(): void {
    this.isCreating = false;
    this.isEditing = false;
    this.variable = this.initialVariable;
  }
}
