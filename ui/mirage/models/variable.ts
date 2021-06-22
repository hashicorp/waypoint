import { Model, belongsTo, hasMany } from 'ember-cli-mirage';
import { Variable, Ref } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),

  toProtobuf(): Variable {
    let result = new Variable();

    result.setServer();
    result.setName(this.name);
    if (this.hcl) {
      result.setHcl(this.hcl);
    } else {
      result.setStr(this.str);
    }

    return result;
  },
});
