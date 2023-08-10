/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { later } from '@ember/runloop';

export default class CopyableCode extends Component {
  @tracked copySuccess = false;

  @action
  onSuccess(): void {
    this.copySuccess = true;

    later(() => {
      this.copySuccess = false;
    }, 2000);
  }
}
