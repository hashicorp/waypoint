import { helper } from '@ember/component/helper';

let BUSY_STATES_FOR_TYPE = {
  'hashicorp.network.hvn': ['UNSET', 'CREATING', 'DELETING'],
  'hashicorp.network.peering': ['UNSET', 'CREATING', 'PENDING_ACCEPTANCE', 'DELETING'],
  'hashicorp.consul.cluster': ['UNSET', 'PENDING', 'CREATING', 'UPDATING', 'RESTORING', 'DELETING'],
};

export function busyStatesForResource([type]) {
  return BUSY_STATES_FOR_TYPE[type].slice(0);
}

export default helper(busyStatesForResource);
