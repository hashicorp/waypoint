import Component from '@glimmer/component';
import { action } from '@ember/object';
import { Project, Application } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';
interface AppFormAppSettingsArgs {
  app: Application;
  project: Project;
}

export default class AppFormAppSettings extends Component<AppFormAppSettingsArgs> {
  @tracked project;
  @tracked app;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { project, app } = this.args;
    this.project = project;
    this.app = app;
  }
  // disabled git settings if project-level git config is not set up
  get gitDisabled() {
    return !(this.project.dataSource?.git?.basic || this.project.dataSource?.git?.ssh)
  }

  @action
  saveSettings() {

  }
}
