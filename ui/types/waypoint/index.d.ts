import Ember from 'ember';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  interface Array<T> extends Ember.ArrayPrototypeExtensions<T> {}
  // interface Function extends Ember.FunctionPrototypeExtensions {}
}

import 'ember-concurrency-async';
import 'ember-concurrency-ts/async';

export {};
