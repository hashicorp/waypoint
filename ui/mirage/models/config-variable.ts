import { Model, belongsTo } from 'ember-cli-mirage';
import { ConfigVar } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),

  toProtobuf(): ConfigVar {
    let result = new ConfigVar();

    // TODO: result.setProject();
    result.setName(this.name);
    result.setStatic(this.pb_static);

    return result;
  },
});
