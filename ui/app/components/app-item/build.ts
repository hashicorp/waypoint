import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class AppItemBuild extends Component {
  @tracked hintIsVisible = false;

  @action
  toggleHint(): boolean {
    if (this.hintIsVisible === true) {
      return (this.hintIsVisible = false);
    } else {
      return (this.hintIsVisible = true);
    }
  }
}
