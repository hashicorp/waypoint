import { helper } from '@ember/component/helper';
import { formatDistance } from 'date-fns';

// dateFormatDistance
export function dateFormatDistance([date, baseDate]: [number, number]): string {
  return formatDistance(date * 1000, baseDate * 1000, { includeSeconds: true });
}

export default helper(dateFormatDistance);
