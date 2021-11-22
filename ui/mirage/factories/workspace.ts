import { Factory } from 'miragejs';
import { trait } from 'ember-cli-mirage';

export default Factory.extend({
  default: trait({
    name: 'default',
  }),
});
