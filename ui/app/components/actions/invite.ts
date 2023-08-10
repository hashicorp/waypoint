/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { InviteTokenRequest } from 'waypoint-pb';
import SessionService from 'ember-simple-auth/services/session';
import { action } from '@ember/object';
import { later } from '@ember/runloop';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class ActionsInvite extends Component {
  @service api!: ApiService;
  @service session!: SessionService;

  @tracked token = '';
  @tracked hintIsVisible = false;
  @tracked copySuccess = false;

  selectContents(element: HTMLInputElement): void {
    element.focus();
    element.select();
  }

  @action
  onSuccess(): void {
    this.copySuccess = true;

    later(() => {
      this.copySuccess = false;
    }, 2000);
  }

  @action
  async createToken(): Promise<void> {
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
  async toggleHint(): Promise<boolean> {
    // Create a token if one doesn't exist
    if (this.token == '') await this.createToken();

    if (this.hintIsVisible === true) {
      return (this.hintIsVisible = false);
    } else {
      return (this.hintIsVisible = true);
    }
  }
}
