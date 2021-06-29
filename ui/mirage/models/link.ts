import { Model, belongsTo } from 'miragejs';
import { Project } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),

  toProtobuf(): Project.Link {
    let result = new Project.Link();

    result.setUrl(this.url);
    result.setLabel(this.label);

    return result;
  },
});
