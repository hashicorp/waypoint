/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { Model as ProjectRouteModel } from 'waypoint/routes/workspace/projects/project';

export default class extends Controller {
  @tracked project?: ProjectRouteModel;
}
