import { ITerminalOptions, Terminal } from 'xterm';

import terminalTheme from 'waypoint/utils/terminal-theme';

interface TerminalOptions {
  inputDisabled: boolean;
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

  if (options.inputDisabled) {
    terminalOptions.disableStdin = true;
    terminalOptions.cursorBlink = false;
    terminalOptions.cursorStyle = 'bar';
    terminalOptions.cursorWidth = 1;
  }

  let terminal = new Terminal(terminalOptions);
  return terminal;
}
