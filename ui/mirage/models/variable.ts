import { Model, belongsTo } from 'ember-cli-mirage';
import { Variable } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),

  toProtobuf(): Variable {
    let result = new Variable();

    result.setServer();
    result.setName(this.name);
    if (this.hcl) {
      result.setStr('');
      result.setHcl(this.hcl);
    } else {
      if (this.str) {
        result.setHcl('');
        result.setStr(this.str);
      }
    }

    return result;
  },
});
