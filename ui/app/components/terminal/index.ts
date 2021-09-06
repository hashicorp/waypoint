import Component from '@glimmer/component';
import { ITerminalOptions, Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import terminalTheme from 'waypoint/utils/terminal-theme';

export default class LogTerminal extends Component {
  terminal: Terminal;
  inputDisabled: boolean;
  fitAddon: FitAddon;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { inputDisabled } = args;
    this.inputDisabled = inputDisabled;

    let terminalOptions: ITerminalOptions = {
      fontFamily: 'ui-monospace,Menlo,monospace',
      fontWeight: '400',
      logLevel: 'debug',
      lineHeight: 1.4,
      fontSize: 12,
      fontWeightBold: '700',
      theme: terminalTheme.light,
    };

    if (this.inputDisabled) {
      terminalOptions.disableStdin = true;
      terminalOptions.cursorBlink = false;
      terminalOptions.cursorStyle = 'bar';
      terminalOptions.cursorWidth = 1;
    }

    let terminal = new Terminal(terminalOptions);
    this.terminal = terminal;
    // Setup resize Addon
    let fitAddon = new FitAddon();
    this.fitAddon = fitAddon;
    this.terminal.loadAddon(fitAddon);
  }

  didInsertNode(element: any): void {
    this.terminal.open(element);
    // Initial fit to component size
    this.fitAddon.fit();

    this.terminal.writeln('Welcome to Waypoint...');
  }

  willDestroyNode(): void {
    this.terminal.dispose();
  }

  didResize(e: Event): void {
    this.fitAddon.fit();
    if (this.terminal.resized) {
      this.terminal.resized(e);
    }
  }
}
