import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
interface StatusBadgeArgs {
  state: number;
  iconOnly: boolean;
}

export default class StatusBadge extends Component<StatusBadgeArgs> {
  @tracked iconOnly = false;
  @tracked state: number;

  constructor(owner: any, args: StatusBadgeArgs) {
    super(owner, args);
    this.state = args.state;
    this.iconOnly = args.iconOnly;
  }

  get statusClass() {
    let { state } = this.args;
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
      default:
        return 'unknown';
    }
  }

  get iconType() {
    let { state } = this.args;
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
      default:
        return 'help-circle-outline';
    }
  }
}
