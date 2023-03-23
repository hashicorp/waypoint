/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

type Args = {
  expanded?: boolean;
  isExpandable?: boolean;
};

export default class extends Component<Args> {
  @tracked expanded = this.args.expanded ?? true;

  get isExpandable(): boolean {
    return this.args.isExpandable ?? true;
  }

  @action
  toggleExpanded(): void {
    this.expanded = !this.expanded;
  }
}
