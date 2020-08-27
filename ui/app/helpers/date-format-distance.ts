import { helper } from '@ember/component/helper';
import { formatDistance } from 'date-fns';

export function dateFormatDistance([date, baseDate]: [number, number]): string {
  debugger;
  return formatDistance(date * 1000, baseDate * 1000, { includeSeconds: true });
}

export default helper(dateFormatDistance);
