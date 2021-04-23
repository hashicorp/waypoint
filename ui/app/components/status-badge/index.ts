import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
interface StatusBadgeArgs {
  model: Object
}

export default class StatusBadge extends Component<StatusBadgeArgs> {
  @tracked model = {};

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
      case 3:
        return 'error';
        break;
      case 2:
        return 'success';
        break;
      case 1:
        return 'running';
        break;
      case 0:
        return 'unknown';
        break;
    }
  }
}
