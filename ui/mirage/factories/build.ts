import { Factory, trait, association } from 'ember-cli-mirage';
import { fakeId } from '../utils';

export default Factory.extend({
  id: () => fakeId(),
  sequence: (i) => i + 1,

  afterCreate(build, server) {
    if (!build.workspace) {
      let workspace =
        server.schema.workspaces.findBy({ name: 'default' }) || server.create('workspace', 'default');
      build.update('workspace', workspace);
    }

    build.pushedArtifact?.update('application', build.application);
    build.pushedArtifact?.update('workspace', build.workspace);
  },

  random: trait({
    labels: () => ({
      'common/vcs-ref': '0d56a9f8456b088dd0e4a7b689b842876fd47352',
      'common/vcs-ref-path': 'https://github.com/hashicorp/waypoint/commit/',
    }),
    component: association('builder', 'with-random-name'),
    status: association('random'),
    pushedArtifact: association('random'),
  }),

  docker: trait({
    component: association('builder', 'docker'),
  }),

  pack: trait({
    component: association('builder', 'pack'),
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
