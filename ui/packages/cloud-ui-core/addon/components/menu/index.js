import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 *
 * `Menu` renders a `<details>` element styled as a drop-down menu. It yields `Trigger` and `Content` contextual components. `Trigger` is a `<summary>` element and `Content` is a `div` that is revealed by `<details>`.
 *
 *
 * ## Example usage
 *
 * ```
 * <Menu as |M| >
 *   <M.Trigger>
 *     Open up!
 *   </M.Trigger>
 *   <M.Content>
 *     <a href="">Sign out</a>
 *   </M.Content>
 * </Menu>
 * ```
 *
 * @class Menu
 * @yield {MenuTrigger} Trigger `Menu::Trigger` component
 * @yield {MenuContent} Content `Menu::Content` component
 *
 */

export default class MenuComponent extends Component {
  /**
   * `isOpen` is a binding that controls if the current `<details>` element's `open` attribute.
   * @argument isOpen
   * @type {boolean}
   */

  /**
   * `onToggle` is a funciton that will be called with the details "toggle" event fires. It will be passed the value of the details `open` attribute.
   * @argument onToggle
   * @type {function}
   */
  @tracked _isOpen = false;

  get isOpen() {
    return this.args.isOpen || this._isOpen;
  }

  set isOpen(val) {
    if (this.args.onToggle && typeof this.args.onToggle === 'function') {
      this.args.onToggle(val);
    }
    this._isOpen = val;
  }

  @action
  trackOpen(evt) {
    this.isOpen = evt.target.open;
  }

  @action
  setupClickOutside(el) {
    document.addEventListener('click', evt => this.clickOutside(evt, el), true);
  }

  @action
  removeClickOutside(el) {
    document.removeEventListener('click', evt => this.clickOutside(evt, el), true);
  }

  @action
  clickOutside(evt, el) {
    if (!el.contains(evt.target)) {
      this.removeClickOutside(el);
      this.isOpen = false;
    }
  }
}
