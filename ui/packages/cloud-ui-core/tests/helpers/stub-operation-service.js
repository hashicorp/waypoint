import Service from '@ember/service';
import { A } from '@ember/array';
import { tracked } from '@glimmer/tracking';

export default class OperationServiceStub extends Service {
  @tracked changedOperations = A([]);
  @tracked _operations = A([]);
  @tracked firstFetch = false;

  set operations(val) {
    val = A(val);
    if (this._operations.length === 0) {
      this._operations = this.changedOperations = val;
      this.firstFetch = true;
      return;
    }
    this.firstFetch = false;
    this.changedOperations = val.reduce((changed, op) => {
      let oldOp = this._operations.findBy('id', op.id);
      // if it's new push to changedOps
      if (!oldOp) {
        changed.push(op);
      } else if (
        // if it's old, and the states don't match
        oldOp &&
        op.state !== oldOp.state
      ) {
        changed.push(op);
      }
      return changed;
    }, []);

    // finally set the new operations
    this._operations = val;
  }

  get operations() {
    return this._operations;
  }
}

