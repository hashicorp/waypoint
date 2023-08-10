/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { Breadcrumb } from 'waypoint/services/breadcrumbs';

export default class WorkspaceProjectsProjectSettingsInputVariables extends Route {
  breadcrumbs(): Breadcrumb[] {
    return [
      {
        label: 'Input Variables',
        route: 'workspace.projects.project.settings.config-variables',
      },
    ];
  }
}
