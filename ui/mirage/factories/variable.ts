import { Factory, trait } from 'ember-cli-mirage';
import faker from '../faker';

export default Factory.extend({
  'random-str': trait({
    name: () => faker.hacker.noun(),
    str: () => faker.hacker.adjective(),
  }),
  'random-hcl': trait({
    name: () => faker.hacker.noun(),
    hcl: () => faker.hacker.adjective(),
  }),
  'is-sensitive': trait({
    name: () => faker.hacker.noun(),
    str: () => faker.hacker.adjective(),
    sensitive: true,
  }),
});
