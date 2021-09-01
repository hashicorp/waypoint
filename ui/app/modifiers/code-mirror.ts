import Modifier from 'ember-modifier';
import codemirror from 'codemirror';

import waypointHclMode from './utils/waypointHclMode';
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/selection/active-line';
import 'codemirror/addon/mode/simple';

const _PRESET_DEFAULTS: codemirror.EditorConfiguration = {
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
    value: string;
    onInput: (value: string) => void;
    options: Record<string, unknown>;
  };
}

// call to create the CodeMirror waypoint HCL mode for syntax highlighting
waypointHclMode();
export default class CodeMirrorModifier extends Modifier<Args> {
  _editor!: codemirror.Editor;

  didInstall(): void {
    this._setup();
  }

  _onChange(editor: codemirror.Editor): void {
    let newVal = editor.getValue();
    this.args.named.onInput(newVal);
  }

  _setup(): void {
    if (!this.element) {
      throw new Error('CodeMirror modifier has no element');
    }

    let editor = codemirror(this.element, {
      ..._PRESET_DEFAULTS,
      ...this.args.named.options,
      value: this.args.named.value,
    });

    editor.on('change', (editor) => {
      this._onChange(editor);
    });

    this._editor = editor;
    this._editor.setOption('mode', 'waypointHCL');
  }
}
