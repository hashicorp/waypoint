import { Model, belongsTo } from 'miragejs';
import { Artifact } from 'waypoint-pb';

export default Model.extend({
  build: belongsTo(),

  toProtobuf(): Artifact {
    let result = new Artifact();

    result.setArtifactJson(JSON.stringify(this.artifact));

    return result;
  },
});
