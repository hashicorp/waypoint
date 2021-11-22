import { Model, belongsTo, hasMany } from 'miragejs';
import { Application, Ref } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),
  builds: hasMany(),
  deployments: hasMany(),
  releases: hasMany(),
  statusReports: hasMany(),

  name: undefined as string | undefined,
  projectId: undefined as string | undefined,

  toProtobuf(): Application {
    let result = new Application();

    result.setFileChangeSignal(this.fileChangeSignal);
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
