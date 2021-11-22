import { Model, hasMany, belongsTo } from 'miragejs';
import { Project, Ref } from 'waypoint-pb';

export default Model.extend({
  applications: hasMany(),
  variables: hasMany(),
  dataSource: belongsTo('job-data-source'),
  dataSourcePoll: belongsTo('project-poll'),

  name: undefined as string | undefined,

  toProtobuf(): Project {
    let result = new Project();

    result.setApplicationsList(this.applications.models.map((a) => a.toProtobuf()));
    result.setDataSource(this.dataSource?.toProtobuf());
    result.setDataSourcePoll(this.dataSourcePoll?.toProtobuf());
    result.setFileChangeSignal(this.fileChangeSignal);
    result.setName(this.name);
    result.setRemoteEnabled(this.remoteEnabled);
    result.setWaypointHcl(this.waypointHcl);
    result.setWaypointHclFormat(Project.Format.HCL);
    let varProtosList = this.variables.models.map((a) => a.toProtobuf());
    result.setVariablesList(varProtosList);
    return result;
  },

  toProtobufRef(): Ref.Project {
    let result = new Ref.Project();

    result.setProject(this.name);

    return result;
  },
});
