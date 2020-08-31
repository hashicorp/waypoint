import Ember from 'ember';

declare global {
  interface Array<T> extends Ember.ArrayPrototypeExtensions<T> {}
  // interface Function extends Ember.FunctionPrototypeExtensions {}
}

import 'ember-concurrency-async';
import 'ember-concurrency-ts/async';

export {};
