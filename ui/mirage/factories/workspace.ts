import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  default: trait({
    name: 'default',
  }),
});
