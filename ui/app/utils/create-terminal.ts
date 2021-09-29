import { ITerminalOptions, Terminal } from 'xterm';

import terminalTheme from 'waypoint/utils/terminal-theme';

interface TerminalOptions {
  inputDisabled: boolean;
  domRendering?: boolean | undefined;
}

export function createTerminal(options: TerminalOptions): Terminal {
  let terminalOptions: ITerminalOptions = {
    fontFamily: 'ui-monospace,Menlo,monospace',
    fontWeight: '400',
    lineHeight: 1.4,
    fontSize: 12,
    fontWeightBold: '700',
    theme: terminalTheme.light,
  };

  // This should be used for DOM rendering in tests only
  if (options.domRendering === true) {
    terminalOptions.rendererType = 'dom';
  }

  // Switch to dark theme if enabled
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    terminalOptions.theme = terminalTheme.dark;
  }

  if (options.inputDisabled) {
    terminalOptions.disableStdin = true;
    terminalOptions.cursorBlink = false;
    terminalOptions.cursorStyle = 'bar';
    terminalOptions.cursorWidth = 1;
  }

  let terminal = new Terminal(terminalOptions);
  return terminal;
}
