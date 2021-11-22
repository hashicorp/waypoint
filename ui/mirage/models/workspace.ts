import { Model } from 'miragejs';
import { Ref, Workspace } from 'waypoint-pb';

export default Model.extend({
  name: undefined as string | undefined,

  toProtobuf(): Workspace {
    let result = new Workspace();

    // TODO: result.setActiveTime
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
