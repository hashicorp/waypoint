import { helper } from '@ember/component/helper';
import { formatDistanceToNow } from 'date-fns';

// dateFormatDistanceToNow
export function dateFormatDistanceToNow([date]: [number]): string {
  return formatDistanceToNow(date * 1000, { includeSeconds: true, addSuffix: true });
}

export default helper(dateFormatDistanceToNow);
