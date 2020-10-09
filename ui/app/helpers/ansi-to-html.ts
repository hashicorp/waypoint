import { helper } from '@ember/component/helper';
import * as Anser from 'anser';

// ansiToHtml
export function ansiToHtml([text]: [string]): string {
  if (!text) return '';

  // Simple escaping
  text = text.replace(/</g, '&lt;').replace(/>/g, '&gt;');

  return Anser.default.ansiToHtml(text);
}

export default helper(ansiToHtml);
