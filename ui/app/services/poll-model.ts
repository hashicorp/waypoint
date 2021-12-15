import Ember from 'ember';
import Service from '@ember/service';
import Route from '@ember/routing/route';
import Transition from '@ember/routing/-private/transition';
import { task, timeout } from 'ember-concurrency';
import { taskFor } from 'ember-concurrency-ts';

// Seconds for polling
const INTERVAL = 15000;

export default class PollModelService extends Service {
  route!: Route;
  transition!: Transition;

  setup(route: Route, transition: Transition): void {
    this.route = route;
    this.transition = transition;

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
        let updatedRouteModel = await this.route.model(
          this.route.paramsFor(this.route.routeName),
          this.transition
        );
        this.route.controllerFor(this.route.routeName).model = updatedRouteModel;
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
