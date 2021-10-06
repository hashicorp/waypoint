export class ImageRef {
  ref: string;

  constructor(ref: string) {
    this.ref = ref;
  }

  get label(): string {
    return this.split[0];
  }

  get tag(): string {
    return this.split[1];
  }

  private get split(): string[] {
    return this.ref.split(':');
  }
}

/**
 * Returns a flat map of values for `image` properties within the given object.
 *
 * @param {object|array} obj search space
 * @param {ImageRef[]} [result=[]] starting result array (used internally, usually no need to pass this)
 * @returns {ImageRef[]} an array of found ImageRefs
 */
export function findImageRefs(obj: unknown, result: ImageRef[] = []): ImageRef[] {
  if (typeof obj !== 'object') {
    return result;
  }

  if (obj === null) {
    return result;
  }

  for (let [key, value] of Object.entries(obj)) {
    if (key.toLowerCase() === 'image' && typeof value === 'string') {
      if (!result.some((image) => image.ref === value)) {
        result.push(new ImageRef(value));
      }
    } else {
      findImageRefs(value, result);
    }
  }

  return result;
}
