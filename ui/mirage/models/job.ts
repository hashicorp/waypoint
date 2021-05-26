import { Model, belongsTo } from 'miragejs';

export default Model.extend({
  application: belongsTo(),
  workspace: belongsTo(),
});
