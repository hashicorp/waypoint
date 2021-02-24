import Component from '@glimmer/component';
import { action } from '@ember/object';
import {Project, Application } from 'waypoint-pb';


interface AppFormAppSettingsArgs {
  app: Application,
  project: Project
}

export default class AppFormAppSettings extends Component<AppFormAppSettingsArgs> {
  @action
  saveSettings() {

  }
}
