import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { later } from '@ember/runloop';

import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { InviteTokenRequest } from 'waypoint-pb';
import SessionService from 'waypoint/services/session';

export default class ActionsInvite extends Component {
  @service api!: ApiService;
  @service session!: SessionService;

  @tracked token = '';
  @tracked hintIsVisible = false;
  @tracked copySuccess = false;

  constructor(owner: any, args: any) {
    super(owner, args);
  }

  selectContents(element: any) {
    element.focus();
    element.select();
  }

  @action
  onSuccess() {
    this.copySuccess = true;

    later(() => {
      this.copySuccess = false;
    }, 2000);
  }

  @action
  async createToken() {
    let req = new InviteTokenRequest();
    req.setDuration('12h');
    let resp = await this.api.client.generateInviteToken(req, this.api.WithMeta());
    this.token = resp.getToken();
  }

  get hostname(): string {
    // There's currently no way for us to retrieve this address from the API
    // so we assume this same URL the user is utilizing is also available to others
    return `${window.location.protocol}//${window.location.host}`;
  }

  @action
  async toggleHint() {
    // Create a token if one doesn't exist
    if (this.token == '') await this.createToken();

    if (this.hintIsVisible === true) {
      return (this.hintIsVisible = false);
    } else {
      return (this.hintIsVisible = true);
    }
  }
}
