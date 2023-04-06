/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import type { Model as Project } from 'waypoint/routes/workspace/projects/project';

export default class extends Service {
  @tracked current?: Project;
}
