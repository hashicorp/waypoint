import { helper } from '@ember/component/helper';

export function enforceProtocol(params/*, hash*/) {
  let str = params[0];

  let isHttps = str.startsWith('https://');
  let isHttp = str.startsWith('http://');
  if (isHttps || isHttp) {
    return str;
  } else {
    return `https://${str}`;
  }
}

export default helper(enforceProtocol);
