import { Model } from 'miragejs';
import { OIDCAuthMethod } from 'waypoint-pb';

export default Model.extend({
  toProtobuf(): OIDCAuthMethod {
    let result = new OIDCAuthMethod();

    result.setDisplayName(this.displayName);
    result.setKind(0);
    result.setName(this.name);
    return result;
  },
});
