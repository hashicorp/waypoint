import { helper } from '@ember/component/helper';

export default helper(([str]: string[]) => str?.substr(0, 7));
