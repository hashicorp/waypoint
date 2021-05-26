import { Factory, trait } from 'ember-cli-mirage';
import faker from '../faker';

export default Factory.extend({
  simple: trait({
    name: 'simple-project',
  }),

  'with-random-name': trait({
    name: () => faker.hacker.noun(),
  }),

  // This is our primary demo trait for development mode
  'marketing-public': trait({
    name: 'marketing-public',
    afterCreate(project, server) {
      let application = server.create('application', 'with-random-name', { project });

      let builds = [
        server.create('build', 'random', 'seconds-old-success', { sequence: 4, application }),
        server.create('build', 'random', 'minutes-old-success', { sequence: 3, application }),
        server.create('build', 'random', 'hours-old-success', { sequence: 2, application }),
        server.create('build', 'random', 'days-old-success', { sequence: 1, application }),
      ];

      let deployments = [
        server.create('deployment', 'random', 'seconds-old-success', {
          sequence: 4,
          application,
          build: builds[0],
        }),
        server.create('deployment', 'random', 'minutes-old-success', {
          sequence: 3,
          application,
          build: builds[1],
        }),
        server.create('deployment', 'random', 'hours-old-success', {
          sequence: 2,
          application,
          build: builds[2],
        }),
        server.create('deployment', 'random', 'days-old-success', {
          sequence: 1,
          application,
          build: builds[3],
        }),
      ];

      server.create('release', 'random', 'minutes-old-success', {
        sequence: 3,
        application,
        deployment: deployments[2],
      });
      server.create('release', 'random', 'hours-old-success', {
        sequence: 2,
        application,
        deployment: deployments[1],
      });
      server.create('release', 'random', 'days-old-success', {
        sequence: 1,
        application,
        deployment: deployments[0],
      });
    },
  }),

  'with-remote-runners': trait({
    remoteEnabled: true,
  }),
});
