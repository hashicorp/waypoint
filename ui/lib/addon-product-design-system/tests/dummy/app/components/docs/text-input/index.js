import Component from '@glimmer/component';

/**
 * @class DocsTextInput
 */

export default class DocsTextInput extends Component {
  get logicallyInvalid () {
    if (this.args.disabled) {
      return false;
    }
    return this.args.required && !this.args.value;
  }
}
