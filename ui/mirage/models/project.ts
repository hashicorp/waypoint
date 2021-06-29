import { Model, hasMany } from 'ember-cli-mirage';
import { Project, Ref } from 'waypoint-pb';

export default Model.extend({
  applications: hasMany(),
  links: hasMany(),

  toProtobuf(): Project {
    let result = new Project();

    result.setApplicationsList(this.applications.models.map((a) => a.toProtobuf()));
    // TODO: result.setDataSource(...)
    // TODO: result.setDataSourcePoll(...)
    // TODO: result.setExtension(...)
    result.setFileChangeSignal(this.fileChangeSignal);
    result.setName(this.name);
    result.setRemoteEnabled(this.remoteEnabled);
    result.setWaypointHcl(this.waypointHcl);
    result.setWaypointHclFormat(Project.Format.HCL);
    result.setLinksList(this.links.models.map((l) => l.toProtobuf()));

    return result;
  },

  toProtobufRef(): Ref.Project {
    let result = new Ref.Project();

    // TODO: result.setExtension(...)
    result.setProject(this.name);

    return result;
  },
});
