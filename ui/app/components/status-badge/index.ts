import Component from '@glimmer/component';

interface StatusBadgeArgs {
  state?: State;
  iconOnly?: boolean;
  message?: string;
}

type State = 'UNKNOWN' | 'ALIVE' | 'READY' | 'DOWN' | 'PARTIAL';

export default class StatusBadge extends Component<StatusBadgeArgs> {
  get statusClass(): string {
    switch (this.args.state) {
      case 'ALIVE':
        return 'alive';
      case 'READY':
        return 'ready';
      case 'DOWN':
        return 'down';
      case 'PARTIAL':
        return 'partial';
      case 'UNKNOWN':
      default:
        return 'unknown';
    }
  }

  get iconType(): string {
    switch (this.args.state) {
      case 'ALIVE':
        return 'run';
      case 'READY':
        return 'check-plain';
      case 'DOWN':
        return 'cancel-circle-fill';
      case 'PARTIAL':
        return 'alert-triangle';
      case 'UNKNOWN':
      default:
        return 'help-circle-outline';
    }
  }
}
