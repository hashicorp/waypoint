import Controller from '@ember/controller';
import RouterService from '@ember/routing/router-service';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class extends Controller {
  @tracked isSwitchingWorkspace = false;
  @service router!: RouterService;
}
