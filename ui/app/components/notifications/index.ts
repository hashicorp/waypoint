import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import FlashMessagesService from 'waypoint/services/flash-messages';

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
  @service flashMessages!: FlashMessagesService;
}
