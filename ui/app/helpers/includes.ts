import { helper } from '@ember/component/helper';

interface Haystack {
  includes(needle: unknown): boolean;
}

export default helper(([haystack, needle]: [unknown, unknown]): boolean => {
  if (isHaystack(haystack)) {
    return haystack.includes(needle);
  } else {
    return false;
  }
});

function isHaystack(haystack: unknown): haystack is Haystack {
  return typeof (haystack as Haystack).includes === 'function';
}
