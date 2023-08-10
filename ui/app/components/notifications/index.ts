/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';

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
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;
}
