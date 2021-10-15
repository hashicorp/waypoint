import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  example: trait({
    username: 'example',
    password: 'example',
  }),
});
