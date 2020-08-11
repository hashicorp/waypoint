import Helper from '@ember/component/helper';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { busyStatesForResource } from './busy-states-for-resource';

/*
 * A helper that determines if a busy state should be shown for a resource.
 * Returns `true` if a busy indicator should be shown, or `false` if not.
 *
 * @example
 * {{resource-is-busy @model.network type='hashicorp.network.hvn'}}
 *
 *
 */

export default class ResourceIsBusy extends Helper {
  @service operation;

  @action
  compute([resource], { type }) {
    // are there any related operations that are not 'DONE' that are linked to
    // the resource?
    let matchedWaitingOp = this.operation.operations.find(op => {
      return op.state !== 'DONE' && op.link.uuid === resource.id;
    });

    // is the resource in any of the states where it would be considered "busy" ?
    let busyStates = busyStatesForResource([type]);

    return !!matchedWaitingOp || busyStates.includes(resource.state);
  }
}
