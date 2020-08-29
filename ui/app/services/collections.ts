import Service, { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import Route from '@ember/routing/route';
import { Ref } from 'waypoint-pb';
import { task } from 'ember-concurrency-decorators';
import { taskFor } from 'ember-concurrency-ts';
import { tracked } from '@glimmer/tracking';
import { timeout, waitForProperty } from 'ember-concurrency';

// Seconds for polling
const INTERVAL = 15000;

export default class CollectionsService extends Service {
  @service api!: ApiService;
  @tracked applicationObject?: Ref.Application.AsObject;
  @tracked workspaceObject?: Ref.Workspace.AsObject;
  @tracked collection?: any;

  // overidden
  async fetchData(): Promise<any> {}

  setup(workspace: Ref.Workspace.AsObject, application: Ref.Application.AsObject, route?: Route) {
    // Optionally configure an observer that refreshes a route
    // model hook
    if (route) {
      this.addObserver('collection', route, route.refresh);
    }

    this.workspaceObject = workspace;
    this.applicationObject = application;

    // Start the API interactions
    this.start();
  }

  willDestroy() {
    this.stop();
    super.willDestroy();
  }

  start() {
    // Make start() idempotent
    if (taskFor(this.poll).isRunning) {
      return;
    }

    taskFor(this.poll).perform();
  }

  stop() {
    taskFor(this.poll).cancelAll();
  }

  async waitForCollection() {
    await waitForProperty(this, 'collection', this.collection);
    return this.collection;
  }

  @task({
    restartable: true,
    maxConcurrency: 1,
  })
  async poll() {
    while (true) {
      try {
        await this.fetchData();
      } catch (e) {
        console.log(e);
      }

      await timeout(INTERVAL);
    }
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    collections: CollectionsService;
  }
}
