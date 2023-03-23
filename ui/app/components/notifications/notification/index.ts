/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import FlashObject from 'ember-cli-flash/flash/object';

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

interface ExtendedFlashObject extends FlashObject {
  content?: string;
  actionText?: string;
  onAction?: (close: Event) => void;
}

type Args = {
  flash: ExtendedFlashObject;
};

export default class NotificationsNotificationComponent extends Component<Args> {
  @action
  onAction(close: Event): void {
    if (typeof this.args.flash.onAction === 'function') {
      this.args.flash.onAction(close);
    }
  }
}
