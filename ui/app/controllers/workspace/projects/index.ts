/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class ProjectsIndex extends Controller {
  queryParams = ['cli'];

  @tracked cli = null;
}
