import { Factory, trait } from 'ember-cli-mirage';
import faker from '../faker';

export default Factory.extend({
  random: trait({
    name: () => faker.hacker.noun(),
    pb_static: () => faker.hacker.adjective(),
  }),

  dynamic: trait({
    name: () => faker.hacker.noun(),
    dynamic: {
      from: () => 'my-config-map',
      configMap: () => [['my-config-map', 'port']],
    }
  }),
});
