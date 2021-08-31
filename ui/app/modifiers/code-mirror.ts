import { getOwner } from '@ember/application';
import { guidFor } from '@ember/object/internals';
import Modifier from 'ember-modifier';
import codemirror from 'codemirror';

import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/edit/closebrackets';
import 'codemirror/addon/selection/active-line';
import 'codemirror/addon/mode/simple';
import CodeMirror from 'waypoint/services/code-mirror';

interface Args {
  positional: never;
  named: {
    value: string;
    onInput: Function;
    options: Object;
  };
}

export default class CodeMirrorModifier extends Modifier<Args> {
  _editor!: codemirror.Editor;

  get cmService(): CodeMirror {
    return getOwner(this).lookup('service:code-mirror');
  }

  didInstall(): void {
    this._setup();
  }

  willRemove(): void {
    this._cleanup();
  }

  _onChange(editor): void {
    let newVal = editor.getValue();
    this.args.named.onInput(newVal);
    this.args.named.value = newVal;
  }

  _setup(): void {
    if (!this.element) {
      throw new Error('CodeMirror modifier has no element');
    }

    // Assign an ID to this element if there is none. This is to
    // ensure that there are unique IDs in the code-mirror service
    // registry.
    if (!this.element.id) {
      this.element.id = guidFor(this.element);
    }

    const _PRESET_DEFAULTS = {
      theme: 'monokai',
      lineNumbers: true,
      cursorBlinkRate: 500,
      matchBrackets: true,
      autoCloseBrackets: true,
      styleActiveLine: true,
    };

    let editor = codemirror(
      this.element,
      Object.assign(
        { value: this.args.named.value ? this.args.named.value : '' },
        this.args.named.options ? this.args.named.options : _PRESET_DEFAULTS
      )
    );

    editor.on('change', (editor) => {
      this._onChange(editor);
    });

    if (this.cmService) {
      this.cmService.registerInstance(this.element.id, editor);
    }

    this._editor = editor;
    this._editor.setOption('mode', 'waypointHCL');
  }

  _cleanup(): void {
    if (this.cmService) {
      this.cmService.unregisterInstance(this.element.id);
    }
  }
}

codemirror.defineSimpleMode('waypointHCL', {
  start: [
    { regex: /(\${)([^}]*)(})/, token: 'null' }, // TODO: formatting within string
    { regex: /"(?:[^\\]|\\.)*?(?:"|$)/, token: 'string' }, // strings
    { regex: /(\w+)(\s+)(=)/, token: ['keyword', 'null', 'null'] }, // assignment of variables
    {
      regex: /(build|deploy|release|hook|registry|type|runner|url)( )({)/,
      token: ['keyword', 'null', 'null'],
    }, // top level keywords
    { regex: /(variable)\b/, token: 'keyword' }, // input variable keyword
    { regex: /true|false|null|undefined/, token: 'atom' }, // bool keywords
    { regex: /0x[a-f\d]+|[-+]?(?:\.\d+|\d+\.?\d*)(?:e[-+]?\d+)?/i, token: 'number' }, // numbers
    { regex: /(#|\/\/)(\s*\S*)/, token: 'comment' }, // single line comments
    { regex: /(path)(.)/, token: ['variable-2', 'null'] }, // path variables
    { regex: /(workspace)(.)(\S*)/, token: ['string', 'null', 'string'] }, // workspace variables
    { regex: /\/\*/, token: 'comment', next: 'comment' }, // multi-line comment
    { regex: /[-+\/*=<>!]+/, token: 'operator' }, // operators
    { regex: /[\{\[\(]/, indent: true }, // for auto indent
    { regex: /[\}\]\)]/, dedent: true },
  ],
  // The multi-line comment state.
  comment: [
    { regex: /.*?\*\//, token: 'comment', next: 'start' },
    { regex: /.*/, token: 'comment' },
  ],
  meta: {
    dontIndentStates: ['comment'],
    lineComment: '#',
  },
});
