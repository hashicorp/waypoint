/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import codemirror from 'codemirror';
import 'codemirror/addon/mode/simple';

codemirror.defineSimpleMode('waypointHCL', {
  start: [
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
    { regex: /[-+/*=<>!]+/, token: 'operator' }, // operators
    { regex: /[{[(]/, indent: true }, // for auto indent
    { regex: /[}\])]/, dedent: true },
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
