import Component from '@glimmer/component';
import { FitAddon } from 'xterm-addon-fit';
import { Terminal } from 'xterm';
import { action } from '@ember/object';

interface TerminalComponentArgs {
  terminal: Terminal;
}

export default class LogTerminal extends Component<TerminalComponentArgs> {
  element!: HTMLElement;
  terminal: Terminal;
  fitAddon!: FitAddon;

  constructor(owner: unknown, args: TerminalComponentArgs) {
    super(owner, args);
    let { terminal } = args;
    this.terminal = terminal;
  }

  @action
  didInsertNode(element: HTMLElement): void {
    this.element = element;
    this._setup();
  }

  _setup(): void {
    if (this.terminal) {
      this.terminal.open(this.element);
      // Initial fit to component size
      // and resize addon setup
      let fitAddon = new FitAddon();
      this.fitAddon = fitAddon;
      this.terminal.loadAddon(fitAddon);
      this.fitAddon.fit();
    }
  }

  willDestroyNode = (): void => {
    this.terminal.dispose();
  };

  didResize = (): void => {
    this.fitAddon.fit();
  };
}
