import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { ConfigVar } from 'waypoint-pb';

interface ProjectConfigArgs {
  variablesList: ConfigVar.AsObject[];
}

export default class ProjectConfigVariablesListComponent extends Component<ProjectConfigArgs> {
  @tracked variablesList: Array<ConfigVar.AsObject>;
  @tracked isCreating: boolean;
  @tracked activeVariable;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { variablesList } = args;
    this.variablesList = variablesList;
    this.activeVariable = null;
    this.isCreating = false;
  }

  @action
  addVariable() {
    this.isCreating = true;
    let newVar = new ConfigVar();
    let newVarObj = newVar.toObject();
    this.variablesList = [newVarObj, ...this.variablesList];
    this.activeVariable = newVarObj;
  }

}
