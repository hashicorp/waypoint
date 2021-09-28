import Component from '@glimmer/component';
import { FitAddon } from 'xterm-addon-fit';
import { Terminal } from 'xterm';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

interface TerminalComponentArgs {
  terminal: Terminal;
}

export default class LogTerminal extends Component<TerminalComponentArgs> {
  element!: HTMLElement;
  terminal: Terminal;
  fitAddon!: FitAddon;
  @tracked isFollowingLogs: boolean;

  constructor(owner: unknown, args: TerminalComponentArgs) {
    super(owner, args);
    let { terminal } = args;
    this.terminal = terminal;
    this.isFollowingLogs = true;
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
      this.terminal.onScroll(() => this.viewPortdidScroll());
      this.element.querySelector('.xterm-viewport')?.addEventListener('scroll', this.viewPortdidScroll);
    }
  }

  setIsFollowingLogs(): void {
    let viewport = this.element.querySelector('.xterm-viewport');
    if (viewport) {
      this.isFollowingLogs = viewport?.scrollTop >= viewport?.scrollHeight - viewport?.offsetHeight;
    }
  }

  viewPortdidScroll = (): void => {
    this.setIsFollowingLogs();
  };

  @action
  followLogs(): void {
    this.terminal.scrollToBottom();
  }

  willDestroyNode = (): void => {
    this.terminal.dispose();
    this.element.querySelector('.xterm-viewport')?.removeEventListener('scroll', this.viewPortdidScroll);
  };

  didResize = (): void => {
    this.fitAddon.fit();
  };
}
