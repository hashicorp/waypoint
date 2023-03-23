/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';
import { fakeComponentForKind } from '../utils';
import { Component } from 'waypoint-pb';

// eslint-disable-next-line ember/require-tagless-components
export default Factory.extend({
  unknown: trait({
    type: 'UNKNOWN',
  }),

  builder: trait({
    type: 'BUILDER',
  }),

  registry: trait({
    type: 'REGISTRY',
  }),

  platform: trait({
    type: 'PLATFORM',
  }),

  'release-manager': trait({
    type: 'RELEASEMANAGER',
  }),

  docker: trait({
    name: 'docker',
  }),

  pack: trait({
    name: 'pack',
  }),

  nomad: trait({
    name: 'nomad',
  }),

  'nomad-jobspec': trait({
    name: 'nomad-jobspec',
  }),

  'nomad-jobspec-canary': trait({
    name: 'nomad-jobspec-canary',
  }),

  kubernetes: trait({
    name: 'kubernetes',
  }),

  'kubernetes-apply': trait({
    name: 'kubernetes-apply',
  }),

  'aws-ecr': trait({
    name: 'aws-ecr',
  }),

  'with-random-name': trait({
    afterCreate(component) {
      component.update('name', randomNameForType(component.type));
    },
  }),
});

function randomNameForType(type): string {
  let kind = Component.Type[type as keyof typeof Component.Type];
  return fakeComponentForKind(kind);
}
