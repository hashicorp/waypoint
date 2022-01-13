import { triggerEvent, triggerKeyEvent, focus } from '@ember/test-helpers';

/**
 * Returns the text from an xterm.js terminal element.
 *
 * Simulates select-all, then copy, then returns the resultant clipboard data.
 *
 * Assumes only one xterm instance is rendered.
 */
export async function getTerminalText(): Promise<string> {
  await focus('.xterm-helper-textarea');

  // Trigger cmd+a and ctrl+a to select all
  await triggerKeyEvent('.xterm-helper-textarea', 'keydown', 65, { metaKey: true });
  await triggerKeyEvent('.xterm-helper-textarea', 'keydown', 65, { ctrlKey: true });

  let clipboardData = new DataTransfer();
  await triggerEvent('.xterm', 'copy', { clipboardData });

  return clipboardData.getData('text').trim();
}
