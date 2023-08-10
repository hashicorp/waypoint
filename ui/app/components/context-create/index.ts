/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { tracked } from '@glimmer/tracking';
import { LoginTokenRequest } from 'waypoint-pb';

type Args = Record<string, never>;

export default class ContextCreate extends Component<Args> {
  @service api!: ApiService;
  @tracked token = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.createToken();
  }

  async createToken(): Promise<void> {
    let resp = await this.api.client.generateLoginToken(new LoginTokenRequest(), this.api.WithMeta());
    this.token = resp.getToken();
  }

  get hostname(): string {
    return `${window.location.hostname}:9701`;
  }

  get contextName(): string {
    return `${window.location.hostname}-ui`;
  }
}
