import { Model, hasMany } from 'ember-cli-mirage';
import { Ref, Workspace } from 'waypoint-pb';

export default Model.extend({
  builds: hasMany(),
  statusReports: hasMany(),

  toProtobuf(): Workspace {
    let result = new Workspace();

    // TODO: result.setActiveTime
    // TODO: result.setExtension
    result.setName(this.name);
    // TODO: result.setProjectsList

    return result;
  },

  toProtobufRef(): Ref.Workspace {
    let result = new Ref.Workspace();

    result.setWorkspace(this.name);

    return result;
  },
});
