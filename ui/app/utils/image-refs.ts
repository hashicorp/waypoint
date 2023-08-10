/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import parse, { Ref } from 'docker-parse-image';

/**
 * Returns a flat map of values for `image` properties within the given object.
 *
 * @param {object|array} obj search space
 * @param {Ref[]} [result=[]] starting result array (used internally, usually no need to pass this)
 * @returns {Ref[]} an array of found ImageRefs
 */
export function findImageRefs(obj: unknown, result: Ref[] = []): Ref[] {
  if (typeof obj !== 'object') {
    return result;
  }

  if (obj === null) {
    return result;
  }

  for (let [key, value] of Object.entries(obj)) {
    if (key.toLowerCase() === 'image' && typeof value === 'string') {
      if (result[value]) {
        // We’ve already seen this ref, continue.
        continue;
      }

      let ref = parse(value);

      result.push(ref);

      // The result array also acts as a map of ref strings to the resultant Ref
      // objects. This little trick is purely internal. Think of it as a “seen”
      // list.
      result[value] = ref;
    } else {
      findImageRefs(value, result);
    }
  }

  return result;
}
