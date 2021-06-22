import { Model, hasMany } from 'ember-cli-mirage';
import { Project, Ref, Variable } from 'waypoint-pb';

export default Model.extend({
  applications: hasMany(),
  variables: hasMany(),
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
    // Somehow adding the toProtoBuf method to the variable wasn't working
    // (probably because of the embedded relationship), so we're converting here for now.
    let varProtosList = this.variables.models.map((a) => {
      let variable = new Variable();
      variable.setName(a.name);
      variable.setServer();
      if (a.hcl) {
        variable.setStr('');
        variable.setHcl(a.hcl);
      } else {
        if (a.str) {
          variable.setHcl('');
          variable.setStr(a.str);
        }
      }
      return variable;
    });
    result.setVariablesList(varProtosList);
    return result;
  },

  toProtobufRef(): Ref.Project {
    let result = new Ref.Project();

    // TODO: result.setExtension(...)
    result.setProject(this.name);

    return result;
  },
});
