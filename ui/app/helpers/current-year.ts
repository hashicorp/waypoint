import { helper } from '@ember/component/helper';

// currentYear
export function currentYear([]): number {
  return new Date().getFullYear();
}

export default helper(currentYear);
