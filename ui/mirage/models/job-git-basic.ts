import { Model, belongsTo } from 'miragejs';
import { Job } from 'waypoint-pb';

export default Model.extend({
  parent: belongsTo('job-git', { inverse: 'basic' }),

  toProtobuf(): Job.Git.Basic {
    let result = new Job.Git.Basic();

    result.setUsername(this.username);
    result.setPassword(this.password);

    return result;
  },
});
