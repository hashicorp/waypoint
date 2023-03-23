/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Project } from 'waypoint-pb';
import { helper } from '@ember/component/helper';

export function projectIsRemoteUppable(params: Array<Project.AsObject> /*, hash*/): boolean {
  let project = params[0] as Project.AsObject;
  // We only want to display the Up button only in this case:
  // if a project has a git datasource, and the dataSourcePoll is not enabled
  return !!project?.dataSource?.git?.url && !project?.dataSourcePoll?.enabled;
}

export default helper(projectIsRemoteUppable);
