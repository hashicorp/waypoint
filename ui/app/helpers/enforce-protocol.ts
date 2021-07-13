import { helper } from '@ember/component/helper';

export function enforceProtocol(params/*, hash*/) {
  let str = params[0];

  let isHttps = str.indexOf('https://') !== -1;
  let isHttp = str.indexOf('http://') !== -1;
  if (isHttps || isHttp) {
    return str;
  } else {
    return `https://${str}`;
  }
}

export default helper(enforceProtocol);
