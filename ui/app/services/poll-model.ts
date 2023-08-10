/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import Service from '@ember/service';
import Route from '@ember/routing/route';
import { task, timeout } from 'ember-concurrency';
import { taskFor } from 'ember-concurrency-ts';

// Seconds for polling
const INTERVAL = 15000;

export default class PollModelService extends Service {
  route!: Route;

  setup(route: Route): void {
    this.route = route;

    // Start polling
    this.start();
  }

  willDestroy(): void {
    this.stop();
    super.willDestroy();
  }

  start(): void {
    if (taskFor(this.poll).isRunning) {
      return;
    }

    taskFor(this.poll).perform();
  }

  stop(): void {
    taskFor(this.poll).cancelAll();
  }

  @task({
    restartable: true,
    maxConcurrency: 1,
  })
  async poll(): Promise<void> {
    // eslint-disable-next-line no-constant-condition
    while (true) {
      if (Ember.testing) {
        return;
      }

      await timeout(INTERVAL);

      try {
        this.route.refresh();
      } catch (e) {
        console.log(e);
      }
    }
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    pollModel: PollModelService;
  }
}
