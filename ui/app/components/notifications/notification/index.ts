import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 *
 * `NotificationsNotification` is the individual notification which renders
 * an Pds::Popup component which, in-turn, is wrapped by a FlashMessage
 * component and utilizes the same api. The first argument in your action
 * will be used as the string title of the notification. The second argument is
 * an object which can take `content`, `actionText`, and `onAction` properties.
 *
 * `content` is a string that is used below the title.
 * `actionText` is a string that is used in a Button component.
 * `onAction` is a function that will be called when the user clicks the button
 *     and gets passed the close function for the alert.
 *
 * Available methods are 'success', 'info', 'warning', and 'error'.
 *
 * this.flashMessages.success('Success!');
 * this.flashMessages.warning('Warning!');
 * this.flashMessages.info('Info!');
 * this.flashMessages.error('Danger!');
 *
 * this.flashMessages.success('Successfully saved!', {
 *   content: 'Success content, baby!',
 *   actionText: 'Celebrate',
 *   onAction: function(close) {
 *     console.log('celebrate');
 *     return close();
 *   },
 * });
 *
 *
 * ```
 * <Notifications::Notification
 *   @flash={{flash}}
 * />
 * ```
 *
 * @class NotificationsNotification
 *
 */

export default class NotificationsNotificationComponent extends Component {
  @action
  onAction(close: Event) {
    let { onAction } = this.args.flash;

    if (onAction && typeof onAction == 'function') {
      this.args.flash.onAction(close);
    }
  }
}
