import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class ActionsRelease extends Component {
  @tracked hintIsVisible = false;

  @action
  showHint() {
    this.hintIsVisible = true;
  }

  @action
  hideHint() {
    this.hintIsVisible = false;
  }
}
