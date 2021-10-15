import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

type Args = {
  expanded?: boolean;
  isExpandable?: boolean;
};

export default class extends Component<Args> {
  @tracked expanded = this.args.expanded ?? true;
  @tracked isExpandable = this.args.isExpandable ?? true;

  @action
  toggleExpanded(): void {
    this.expanded = !this.expanded;
  }
}
