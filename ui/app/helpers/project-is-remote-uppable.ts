import { helper } from '@ember/component/helper';

export function projectIsRemoteUppable(params/*, hash*/) {
  let project = params[0] as Project.AsObject;
  // We only want to display the Up button only in this case:
  // if a project has a git datasource, and the dataSourcePoll is not enabled
  return !!project?.dataSource?.git?.url && !project?.dataSourcePoll?.enabled;
}

export default helper(projectIsRemoteUppable);
