import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 *
 * `Notifications` utilizes the FlashMessage queue to render Notification
 * child elements.
 *
 *
 * ```
 * <Notifications />
 * ```
 *
 * @class Notifications
 *
 */

export default class NotificationsComponent extends Component {
  @service flashMessages!: any;
}
