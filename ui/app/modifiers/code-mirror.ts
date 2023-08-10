/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Modifier from 'ember-modifier';
import codemirror from 'codemirror';

import './utils/register-waypoint-hcl-mode';
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/selection/active-line';
import 'codemirror/mode/javascript/javascript';

const _PRESET_DEFAULTS: codemirror.EditorConfiguration = {
  mode: 'waypointHCL',
  theme: 'monokai',
  lineNumbers: true,
  cursorBlinkRate: 500,
  matchBrackets: true,
  autoCloseBrackets: true,
  styleActiveLine: true,
};
interface Args {
  positional: never;
  named: {
    value?: string;
    onInput?: (value: string) => void;
    options?: codemirror.EditorConfiguration;
  };
}
export default class CodeMirrorModifier extends Modifier<Args> {
  _editor!: codemirror.Editor;

  didInstall(): void {
    this._setup();
  }

  didUpdateArguments(): void {
    let value = this.args.named.value ?? '';
    let options = this.args.named.options;

    if (value !== this._editor.getValue()) {
      this._editor.setValue(value);
    }

    if (options) {
      eachEntry(options, (key, value) => {
        this._editor.setOption(key, value);
      });
    }
  }

  _onChange(editor: codemirror.Editor): void {
    let newVal = editor.getValue();

    if (typeof this.args.named.onInput === 'function') {
      this.args.named.onInput(newVal);
    }
  }

  _setup(): void {
    if (!this.element) {
      throw new Error('CodeMirror modifier has no element');
    }

    let editor = codemirror(this.element, {
      ..._PRESET_DEFAULTS,
      ...this.args.named.options,
      value: this.args.named.value ?? '',
    });

    editor.on('change', (editor) => {
      this._onChange(editor);
    });

    this._editor = editor;
  }
}

// Object.entries loses type information, so this is a workaround.
function eachEntry<T>(object: T, callback: (key: keyof T, value: T[keyof T]) => void): void {
  for (let key in object) {
    callback(key, object[key]);
  }
}
