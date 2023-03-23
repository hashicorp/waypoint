/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';
import faker from '../faker';
import { Status } from 'waypoint-pb';

export default Factory.extend({
  afterCreate(status) {
    let minutes = faker.random.number({ min: 1, max: 10 });
    let startTime = new Date(status.completeTime.valueOf() - minutes * 60 * 1000);

    status.update('startTime', startTime);
  },

  random: trait({
    state: () => randomStateName(),
    details: 'Example status details',
    completeTime: () => faker.date.recent(),
  }),

  success: trait({
    state: 'SUCCESS',
  }),

  error: trait({
    state: 'ERROR',
  }),

  'seconds-old': trait({
    completeTime: () => new Date(),
  }),

  'minutes-old': trait({
    completeTime: () => new Date(new Date().valueOf() - faker.random.number({ min: 1, max: 15 }) * 60 * 1000),
  }),

  'hours-old': trait({
    completeTime: () =>
      new Date(new Date().valueOf() - faker.random.number({ min: 1, max: 5 }) * 60 * 60 * 1000),
  }),

  'days-old': trait({
    completeTime: () =>
      new Date(new Date().valueOf() - faker.random.number({ min: 1, max: 5 }) * 24 * 60 * 60 * 1000),
  }),
});

type StateName = keyof typeof Status.State;
function randomStateName(): StateName {
  return sample(Object.keys(Status.State)) as StateName;
}

function sample<T>(array: T[]): T {
  return array[Math.floor(Math.random() * array.length)];
}
