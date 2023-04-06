/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';
import * as Anser from 'anser';
import { htmlSafe } from '@ember/template';

// ansiToHtml
export function ansiToHtml([text]: [string]): string | ReturnType<typeof htmlSafe> {
  if (!text) return '';

  // Simple escaping
  text = text.replace(/</g, '&lt;').replace(/>/g, '&gt;');

  return htmlSafe(Anser.default.ansiToHtml(text));
}

export default helper(ansiToHtml);
