/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { VersionInfo as info } from 'waypoint-pb';

type Args = Record<string, never>;

export default class VersionInfo extends Component<Args> {
  @service api!: ApiService;
  @tracked versionInfo: info.AsObject | undefined;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.getVersionInfo();
  }

  async getVersionInfo(): Promise<void> {
    let resp = await this.api.client.getVersionInfo(new Empty(), this.api.WithMeta());
    let versionInfo = resp?.getInfo();

    this.versionInfo = versionInfo?.toObject();
  }
}
