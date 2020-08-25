import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { action } from '@ember/object';
import SessionService from 'waypoint/services/session';
import RouterService from '@ember/routing/router-service';
import ApiService from 'waypoint/services/api';
import { ConvertInviteTokenRequest } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

export default class InviteLoginForm extends Component {
  @service session!: SessionService;
  @service router!: RouterService;
  @service api!: ApiService;

  @tracked inviteToken = '';

  @action
  async login() {
    var req = new ConvertInviteTokenRequest();
    req.setToken(this.inviteToken);
    var resp = await this.api.client.convertInviteToken(req, this.api.WithMeta());
    await this.session.setToken(resp.getToken());
    return this.router.transitionTo('workspaces');
  }
}
