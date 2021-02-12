import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class ActionsRelease extends Component {
  @tracked hintIsVisible = false;

  @action
  toggleHint() {
    if (this.hintIsVisible === true) {
      return this.hintIsVisible = false;
    } else {
      return this.hintIsVisible = true;
    };
  }
}
