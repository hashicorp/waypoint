import { Factory, trait } from 'ember-cli-mirage';

export default Factory.extend({
  'every-2-minutes': trait({
    enabled: true,
    interval: '2m',
  }),
});
