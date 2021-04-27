import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
interface StatusBadgeArgs {
  model: Object
}

export default class StatusBadge extends Component<StatusBadgeArgs> {
  @tracked model = {};
  @tracked iconOnly = false;

  constructor(owner: any, args: StatusBadgeArgs) {
    super(owner, args);
    this.model = args.model;
  }

  get state() {
    return this.model?.status?.state;
  }

  get statusClass() {
    let state = this.state;
    switch (state) {
      case 4:
        return 'partial';
        break;
      case 3:
        return 'down';
        break;
      case 2:
        return 'ready';
        break;
      case 1:
        return 'alive';
        break;
      case 0:
        return 'unknown';
        break;
    }
  }

  get iconType() {
    let state = this.state;
    switch (state) {
      case 4:
        return 'alert-triangle';
        break;
      case 3:
        return 'cancel-circle-fill';
        break;
      case 2:
        return 'check-plain';
        break;
      case 1:
        return 'run';
        break;
      case 0:
        return 'help-circle-outline';
        break;
    }
  }
}
