/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

/**
  Determines if the given haystack contains the given needle. Delegates to
  `haystack.includes`, if such a method is present, otherwise returns false.

  JavaScript’s String, Array, and TypedArray all implement `includes`.

  For example:

  ```handlebars
  <input type="checkbox" checked={{includes @ingredients "cheese"}}>

  {{#if (includes @ingredients "cheese")}}
    With cheese
  {{/if}}
  ```

  Why do we have this helper when ember-composable-helpers already ships an
  `includes` helper?
  https://github.com/DockYard/ember-composable-helpers#includes

  Simply because the `includes` helper from ember-composable-helpers only works
  with Arrays and Ember arrays.

  @param {any} haystack - the string, array etc. to search within
  @param {any} needle - the thing to search for within haystack
  @returns {boolean} whether or not haystack includes needle
 */
export default helper(([haystack, needle]: [unknown, unknown]) => includes(haystack, needle));

export function includes(haystack: unknown, needle: unknown): boolean {
  if (isHaystack(haystack)) {
    return haystack.includes(needle);
  } else {
    return false;
  }
}

function isHaystack(haystack?: unknown): haystack is Haystack {
  return !!haystack && typeof (haystack as Haystack).includes === 'function';
}

interface Haystack {
  includes(needle: unknown): boolean;
}
