import { Factory, trait, association } from 'ember-cli-mirage';
import { fakeId } from '../utils';

export default Factory.extend({
  id: () => fakeId(),
  sequence: (i) => i + 1,

  afterCreate(deployment, server) {
    if (!deployment.workspace) {
      let workspace =
        server.schema.workspaces.findBy({ name: 'default' }) || server.create('workspace', 'default');
      deployment.update('workspace', workspace);
    }
  },

  random: trait({
    component: association('platform', 'with-random-name'),
    status: association('random'),
    state: 'CREATED',
    labels: () => ({
      'common/vcs-ref': '0d56a9f8456b088dd0e4a7b689b842876fd47352',
      'common/vcs-ref-path': 'https://github.com/hashicorp/waypoint/commit/',
    }),

    afterCreate(deployment) {
      let url = `https://wildly-intent-honeybee--v${deployment.sequence}.waypoint.run`;
      deployment.update('deployUrl', url);
    },
  }),

  docker: trait({
    component: association('platform', 'docker'),
  }),

  nomad: trait({
    component: association('platform', 'nomad'),
  }),

  'nomad-jobspec': trait({
    component: association('platform', 'nomad-jobspec'),
  }),

  kubernetes: trait({
    component: association('platform', 'kubernetes'),
  }),

  'kubernetes-apply': trait({
    component: association('platform', 'kubernetes-apply'),
  }),

  'seconds-old-success': trait({
    status: association('random', 'success', 'seconds-old'),
  }),

  'minutes-old-success': trait({
    status: association('random', 'success', 'minutes-old'),
  }),

  'hours-old-success': trait({
    status: association('random', 'success', 'hours-old'),
  }),

  'days-old-success': trait({
    status: association('random', 'success', 'days-old'),
  }),
});
