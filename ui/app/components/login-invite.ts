import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
import { action } from '@ember/object';
import SessionService from 'waypoint/services/session';
import RouterService from '@ember/routing/router-service';
import ApiService from 'waypoint/services/api';
import { ConvertInviteTokenRequest } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

interface InviteLoginFormArgs {
  inviteToken: string;
  cli: boolean;
}

export default class InviteLoginForm extends Component<InviteLoginFormArgs> {
  @service session!: SessionService;
  @service router!: RouterService;
  @service api!: ApiService;

  @tracked inviteToken = '';
  @tracked cli = null;

  constructor(owner: any, args: any) {
    super(owner, args);

    let { cli } = this.args;

    // If this is a CLI invite login, do it automatically when the component loads
    if (cli) {
      this.login();
    }
  }

  @action
  async login() {
    var req = new ConvertInviteTokenRequest();
    req.setToken(this.inviteToken);
    var resp = await this.api.client.convertInviteToken(req, this.api.WithMeta());
    await this.session.setToken(resp.getToken());

    // If this is an invite for a new user, take them to on-boarding, otherwise, take
    // them to the workspaces page with a query parameter to specify it came
    // from the CLI
    if (this.cli) {
      return this.router.transitionTo('workspaces', { queryParams: { cli: this.cli } });
    } else {
      return this.router.transitionTo('onboarding.install');
    }
  }
}
