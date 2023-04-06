/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { FitAddon } from 'xterm-addon-fit';
import { Terminal } from 'xterm';
import { action } from '@ember/object';
import { assert } from '@ember/debug';
import { debounce } from '@ember/runloop';
import { tracked } from '@glimmer/tracking';

interface TerminalComponentArgs {
  terminal?: Terminal;
}

export default class LogTerminal extends Component<TerminalComponentArgs> {
  element!: HTMLElement;
  terminal?: Terminal;
  fitAddon!: FitAddon;
  @tracked isFollowingLogs: boolean;

  constructor(owner: unknown, args: TerminalComponentArgs) {
    super(owner, args);
    let { terminal } = args;
    assert('A terminal object must be passed to the component', !!terminal);
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
    this.isFollowingLogs = this.isScrolledToBottom();
  }

  viewPortdidScroll = (): void => {
    this.setIsFollowingLogs();
  };

  @action
  followLogs(): void {
    this.terminal?.scrollToBottom();
  }

  willDestroyNode = (): void => {
    this.terminal?.dispose();
    delete this.terminal;
    this.element.querySelector('.xterm-viewport')?.removeEventListener('scroll', this.viewPortdidScroll);
  };

  didResize = (): void => {
    debounce(this, this.fitTerminal, 150);
  };

  fitTerminal(): void {
    // Set terminal size to a minimum
    // before calling resize to avoid reflows when sizing down
    this.terminal?.resize(40, 1);
    this.fitAddon.fit();
  }

  isScrolledToBottom(): boolean {
    let viewport = this.element.querySelector('.xterm-viewport') as HTMLElement;
    if (viewport) {
      return viewport?.scrollTop >= viewport?.scrollHeight - viewport?.offsetHeight;
    } else {
      return false;
    }
  }
}
