import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 *
 * The `OperationFailureForResource` looks for failed operations
 * from the operations service using the passed resource. If the
 * resource's `state` is `FAILED` and a failed operation with a Link
 * to the resource is found, the operation's error message will be
 * rendered in an `AlertBanner` component.
 *
 *
 * ```
 * <Operation::FailureForResource @resource={{@model.cluster}} />
 * ```
 *
 * @class OperationFailureForResource
 *
 */

export default class OperationFailureForResourceComponent extends Component {
  /**
   * The `resource` is the object that will be used to find a related operation
   * failure. It needs to have `state` and `id` attributes.
   * @argument resource
   * @type {object}
   */

  @service operation;

  get failedOperationError() {
    if (!this.operation.operations) {
      return null;
    }

    let relatedOperations = this.operation.operations.reduce(
      (related, op) => {
        // it's not a related operation return the accumulator
        if (op.link.uuid !== this.args.resource.id) {
          return related;
        }
        if (op.state === 'RUNNING') {
          related.running.push(op);
        }

        // if the operation has an error and we've not already marked a related
        if (op.error && !related.errored) {
          related.errored = op;
        }
        return related;
      },
      {
        running: [],
        errored: null,
      }
    );

    // if there are running related operation, don't show any errors,
    // otherwise return the error from the errored operation
    return relatedOperations.running.length ? null : relatedOperations.errored?.error;
  }
}
