import { Factory, trait } from 'ember-cli-mirage';
import faker from '../faker';

export default Factory.extend({
  simple: trait({
    name: 'simple-project',
  }),

  'with-random-name': trait({
    name: () => faker.hacker.noun(),
  }),

  'with-remote-runners': trait({
    remoteEnabled: true,
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

  // For demoing and working against mutable deployments
  'mutable-deployments': trait({
    name: 'mutable-project',

    afterCreate(project, server) {
      let application = server.create('application', { name: 'mutable-application', project });

      let builds = [
        server.create('build', 'docker', 'days-old-success', {
          application,
          sequence: 1,
        }),
        server.create('build', 'docker', 'days-old-success', {
          application,
          sequence: 2,
        }),
        server.create('build', 'docker', 'hours-old-success', {
          application,
          sequence: 3,
        }),
        server.create('build', 'docker', 'hours-old-success', {
          application,
          sequence: 4,
        }),
        server.create('build', 'docker', 'minutes-old-success', {
          application,
          sequence: 5,
        }),
        server.create('build', 'docker', 'minutes-old-success', {
          application,
          sequence: 6,
        }),
        server.create('build', 'docker', 'seconds-old-success', {
          application,
          sequence: 7,
        }),
      ];

      let generations = [
        server.create('generation', {
          id: 'job-v1',
          initialSequence: 1,
        }),
        server.create('generation', {
          id: 'job-v2',
          initialSequence: 4,
        }),
      ];

      let deployments = [
        server.create('deployment', 'random', 'nomad-jobspec', 'days-old-success', {
          application,
          generation: generations[0],
          build: builds[0],
          sequence: 1,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'days-old-success', {
          application,
          generation: generations[0],
          build: builds[1],
          sequence: 2,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'hours-old-success', {
          application,
          generation: generations[0],
          build: builds[2],
          sequence: 3,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'hours-old-success', {
          application,
          generation: generations[1],
          build: builds[2],
          sequence: 4,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'minutes-old-success', {
          application,
          generation: generations[1],
          build: builds[2],
          sequence: 5,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'minutes-old-success', {
          application,
          generation: generations[1],
          build: builds[2],
          sequence: 6,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'seconds-old-success', {
          application,
          generation: generations[1],
          build: builds[2],
          sequence: 7,
        }),
      ];

      server.create('release', 'random', 'nomad-jobspec', 'seconds-old-success', {
        sequence: 1,
        deployment: deployments[0],
        application,
      });
    },
  }),
});
