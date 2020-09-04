import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class ActionsInvite extends Component {
  @tracked hintIsVisible = false;

  selectContents(element: any) {
    element.focus();
    element.select();
  }

  @action
  toggleHint() {
    if (this.hintIsVisible === true) {
      return this.hintIsVisible = false;
    } else {
      return this.hintIsVisible = true;
    };
  }
}
