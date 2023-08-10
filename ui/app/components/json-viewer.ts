/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  json?: string;
  label?: string;
}

export default class extends Component<Args> {
  get parseResult(): { data: unknown } | { error: Error } {
    if (!this.args.json) {
      return { error: new Error('No source JSON provided') };
    }

    try {
      return { data: JSON.parse(this.args.json) };
    } catch (error) {
      return { error };
    }
  }

  get formattedJSON(): string | undefined {
    if ('data' in this.parseResult) {
      return JSON.stringify(this.parseResult.data, null, 2);
    } else {
      return;
    }
  }

  get screenReaderLabel(): string {
    return this.args.label ?? 'JSON';
  }

  get codeMirrorOptions(): Record<string, unknown> {
    return {
      mode: { name: 'javascript', json: true },
      readOnly: true,
      viewportMargin: Infinity,
      screenReaderLabel: this.screenReaderLabel,
    };
  }
}
