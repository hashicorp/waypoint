import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class extends Controller {
  @tracked expanded: string[] = [];

  @action
  toggleExpanded(name: string): void {
    if (this.expanded.includes(name)) {
      this.expanded = this.expanded.filter((s) => s !== name);
    } else {
      this.expanded = [...this.expanded, name];
    }
  }
}
