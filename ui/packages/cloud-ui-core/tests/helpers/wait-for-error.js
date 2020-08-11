import { waitUntil } from '@ember/test-helpers';
import Ember from 'ember';

export default function waitForError(opts) {
  let orig = Ember.onerror;

  let error = null;
  Ember.onerror = err => {
    error = err;
  };

  return waitUntil(() => error, opts).finally(() => {
    Ember.onerror = orig;
  });
}
