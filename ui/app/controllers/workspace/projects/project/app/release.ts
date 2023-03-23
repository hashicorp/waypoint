/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { Model } from 'waypoint/routes/workspace/projects/project/app/release';
import { tracked } from '@glimmer/tracking';

export default class extends Controller {
  @tracked model!: Model;

  get shouldShowURL(): boolean {
    return !!this.model.url && this.isLatest;
  }

  get isLatest(): boolean {
    return this.model.id === this.model.latestRelease?.id;
  }
}
