import { Model, belongsTo, hasMany } from 'ember-cli-mirage';
import { Application, Ref } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),
  builds: hasMany(),
  deployments: hasMany(),
  statusReports: hasMany(),

  toProtobuf(): Application {
    let result = new Application();

    // TODO: result.setFileChangeSignal(...)
    result.setName(this.name);
    result.setProject(this.project.toProtobufRef());

    return result;
  },

  toProtobufRef(): Ref.Application {
    let result = new Ref.Application();

    result.setApplication(this.name);
    result.setProject(this.project?.name);

    return result;
  },
});
