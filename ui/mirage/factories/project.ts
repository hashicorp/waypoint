import { Factory, trait, association } from 'ember-cli-mirage';
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

  'with-input-variables': trait({
    name: 'with-input-variables',

    afterCreate(project, server) {
      server.createList('variable', 2, 'random-str', { project });
      server.create('variable', 'random-hcl', { project });
      server.create('variable', 'is-sensitive', { project });
    },
  }),

  // This is our primary demo trait for development mode
  'marketing-public': trait({
    name: 'marketing-public',

    dataSource: association('marketing-public'),
    dataSourcePoll: association('every-2-minutes'),

    afterCreate(project, server) {
      let application = server.create('application', 'with-random-name', { project });

      let prodWorkspace = server.create('workspace', { name: 'production' });

      server.create('config-variable', 'random', { project, name: 'test' });

      let builds = [
        server.create('build', 'random', 'seconds-old-success', { sequence: 4, application }),
        server.create('build', 'random', 'minutes-old-success', { sequence: 3, application }),
        server.create('build', 'random', 'hours-old-success', { sequence: 2, application }),
        server.create('build', 'random', 'days-old-success', { sequence: 1, application }),
      ];

      let prodBuilds = [
        server.create('build', 'random', 'seconds-old-success', {
          sequence: 5,
          application,
          workspace: prodWorkspace,
        }),
      ];

      let deployments = [
        server.create('deployment', 'random', 'seconds-old-success', {
          sequence: 4,
          application,
          build: builds[0],
          statusReport: server.create('status-report', 'alive', 'with-deployment-resources', { application }),
        }),
        server.create('deployment', 'random', 'minutes-old-success', {
          sequence: 3,
          application,
          build: builds[1],
          statusReport: server.create('status-report', 'ready', 'with-deployment-resources', { application }),
        }),
        server.create('deployment', 'random', 'hours-old-success', {
          sequence: 2,
          application,
          build: builds[2],
          statusReport: server.create('status-report', 'partial', 'with-deployment-resources', {
            application,
          }),
        }),
        server.create('deployment', 'random', 'days-old-success', {
          sequence: 1,
          application,
          build: builds[3],
          statusReport: server.create('status-report', 'down', 'with-deployment-resources', { application }),
        }),
      ];

      server.create('deployment', 'random', 'seconds-old-success', {
        sequence: 5,
        application,
        workspace: prodWorkspace,
        build: prodBuilds[0],
        statusReport: server.create('status-report', 'alive', 'with-deployment-resources', { application }),
      });

      server.create('release', 'random', 'minutes-old-success', {
        sequence: 3,
        application,
        deployment: deployments[0],
        statusReport: server.create('status-report', 'ready', 'with-release-resources', { application }),
      });
      server.create('release', 'random', 'hours-old-success', {
        sequence: 2,
        application,
        deployment: deployments[1],
      });
      server.create('release', 'random', 'days-old-success', {
        sequence: 1,
        application,
        deployment: deployments[2],
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
          pushedArtifact: server.create('pushed-artifact', 'docker', 'days-old-success'),
        }),
        server.create('build', 'docker', 'days-old-success', {
          application,
          sequence: 2,
          pushedArtifact: server.create('pushed-artifact', 'docker', 'days-old-success'),
        }),
        server.create('build', 'docker', 'hours-old-success', {
          application,
          sequence: 3,
          pushedArtifact: server.create('pushed-artifact', 'docker', 'hours-old-success'),
        }),
        server.create('build', 'docker', 'hours-old-success', {
          application,
          sequence: 4,
          pushedArtifact: server.create('pushed-artifact', 'docker', 'hours-old-success'),
        }),
        server.create('build', 'docker', 'minutes-old-success', {
          application,
          sequence: 5,
          pushedArtifact: server.create('pushed-artifact', 'docker', 'minutes-old-success'),
        }),
        server.create('build', 'docker', 'minutes-old-success', {
          application,
          sequence: 6,
          pushedArtifact: server.create('pushed-artifact', 'docker', 'minutes-old-success'),
        }),
        server.create('build', 'docker', 'seconds-old-success', {
          application,
          sequence: 7,
          pushedArtifact: server.create('pushed-artifact', 'docker', 'seconds-old-success'),
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
          state: 'DESTROYED',
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
          build: builds[3],
          sequence: 4,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'minutes-old-success', {
          application,
          generation: generations[1],
          build: builds[4],
          sequence: 5,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'minutes-old-success', {
          application,
          generation: generations[1],
          build: builds[5],
          sequence: 6,
        }),
        server.create('deployment', 'random', 'nomad-jobspec', 'seconds-old-success', {
          application,
          generation: generations[1],
          build: builds[6],
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

  'example-nodejs': trait({
    name: 'example-nodejs',

    afterCreate(project, server) {
      let application = server.create('application', { name: 'example-nodejs', project });

      server.create('build', 'pack', 'minutes-old-success', { application, sequence: 1 });
      server.create('build', 'pack', 'seconds-old-success', { application, sequence: 2 });

      server.create('deployment', 'docker', 'minutes-old-success', {
        application,
        sequence: 1,
        statusReport: server.create('status-report', 'ready', { application }),
        deployUrl: `https://instantly-worthy-shrew--v1.waypoint.run`,
      });
      server.create('deployment', 'docker', 'seconds-old-success', {
        application,
        sequence: 2,
        statusReport: server.create('status-report', 'ready', { application }),
        deployUrl: `https://instantly-worthy-shrew--v2.waypoint.run`,
      });

      server.create('release', 'docker', 'minutes-old-success', {
        application,
        sequence: 1,
        statusReport: server.create('status-report', 'ready', { application }),
      });
      server.create('release', 'docker', 'seconds-old-success', {
        application,
        sequence: 2,
        statusReport: server.create('status-report', 'ready', { application }),
      });
    },
  }),
});
