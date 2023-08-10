/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import * as AnsiColors from 'ansi-colors';

import { ExecStreamRequest, ExecStreamResponse } from 'waypoint-pb';
import { IDisposable, Terminal } from 'xterm';

import { AttachAddon } from 'xterm-addon-attach';
import KEYS from 'waypoint/utils/keys';
import { tracked } from '@glimmer/tracking';

const BACKSPACE_ONE_CHARACTER = '\x08 \x08';

// eslint-disable-next-line no-control-regex
const UNPRINTABLE_CHARACTERS_REGEX = /[\x00-\x1F]/g;
interface IAttachOptions {
  bidirectional?: boolean;
  deploymentId?: string;
}

type WaypointCloseEvent = CloseEvent & MessageEvent;

export class ExecWebSocketAddon extends AttachAddon {
  private _socket: WebSocket;
  private _deploymentId?: string;
  private _bidirectional: boolean;
  private _disposables: IDisposable[] = [];
  encoder: TextEncoder;
  @tracked terminal!: Terminal;
  // Command is used to store the current command until it gets sent over WS
  @tracked command!: string;

  constructor(socket: WebSocket, options?: IAttachOptions) {
    super(socket, options);
    this._socket = socket;
    this._deploymentId = options?.deploymentId;
    this.command = '';
    // always set binary type to arraybuffer, we do not handle blobs
    this._socket.binaryType = 'arraybuffer';
    this._bidirectional = !(options && options.bidirectional === false);
    this.encoder = new TextEncoder();
  }

  activate(terminal: Terminal): void {
    this.terminal = terminal;
    if (!this._deploymentId) {
      // If there is no deployment ID to connect to, warn the user.
      terminal.clear();
      terminal.writeln(AnsiColors.yellow('No deployment available'));
      return;
    }

    this._disposables.push(
      addSocketListener(this._socket, 'open', () => {
        terminal.writeln(AnsiColors.bold.cyan('Connecting...'));
        if (this._deploymentId) {
          this.sendStart(this._deploymentId, terminal);
          this.setupWinch(terminal);
        }
      }),
      addSocketListener(this._socket, 'message', (event: MessageEvent) => {
        let output = event.data;
        let resp = ExecStreamResponse.deserializeBinary(output);
        if (resp.hasOpen()) {
          // remove "Connecting..." message
          terminal.clear();
          terminal.writeln(AnsiColors.bold.cyan(`Connected to deployment: ${this._deploymentId}`));
        }
        if (resp.getEventCase() === 2 && resp.getExit()) {
          let exitCode = resp.getExit() || 'unknown';
          terminal.writeln(AnsiColors.yellow(`Exit code: ${exitCode}`));
          terminal.writeln(AnsiColors.yellow('Connection closed...'));
        }
        // Output
        if (resp.getEventCase() === 1 && resp.getOutput()) {
          let uint8data = resp.getOutput()?.getData_asU8() as Uint8Array;
          terminal.write(uint8data);
        }
        // Event not set
        if (resp.getEventCase() === 0) {
          terminal.writeln(AnsiColors.yellow('Connection closed...'));
        }
      }),
      addSocketListener(this._socket, 'close', (event: WaypointCloseEvent) => {
        let output = event.data;
        if (output) {
          let resp = ExecStreamResponse.deserializeBinary(output);
          let uint8data = resp.getOutput()?.getData_asU8() as Uint8Array;
          terminal.write(uint8data);
          terminal.writeln(AnsiColors.yellow('Connection closed...'));
        }
      }),
      addSocketListener(this._socket, 'error', (event: MessageEvent) => {
        let output = event.data;
        if (output) {
          let resp = ExecStreamResponse.deserializeBinary(output);
          let uint8data = resp.getOutput()?.getData_asU8() as Uint8Array;
          terminal.write(uint8data);
        }
      })
    );

    if (this._bidirectional) {
      this._disposables.push(terminal.onData((data) => this.handleOnData(data)));
      this._disposables.push(terminal.onBinary((data) => this.handleOnBinary(data)));
    }
  }

  handleOnData(data: string): void {
    if (
      data === KEYS.LEFT_ARROW ||
      data === KEYS.UP_ARROW ||
      data === KEYS.RIGHT_ARROW ||
      data === KEYS.DOWN_ARROW
    ) {
      // Ignore arrow keys
    } else if (data === KEYS.CONTROL_U) {
      this.terminal.write(BACKSPACE_ONE_CHARACTER.repeat(this.command.length));
      this.command = '';
    } else if (data === KEYS.ENTER) {
      // We remove the characters here because since the response already contains those
      this.terminal.write(BACKSPACE_ONE_CHARACTER.repeat(this.command.length));
      this.command += KEYS.ENTER;
      this._sendData(this.command);
      this.command = '';
    } else if (data === KEYS.DELETE) {
      if (this.command.length > 0) {
        this.terminal.write(BACKSPACE_ONE_CHARACTER);
        this.command = this.command.slice(0, -1);
      }
    } else if (data.length > 0) {
      let strippedData = data.replace(UNPRINTABLE_CHARACTERS_REGEX, '');
      this.terminal.write(strippedData);
      this.command = `${this.command}${strippedData}`;
    }
  }

  handleOnBinary(data: string): void {
    this._sendBinary(data);
  }

  // dispose is called by the RenderTerminal component when it gets removed
  dispose(): void {
    this._disposables.forEach((d) => d.dispose());
    this._disposables.length = 0;
  }

  private _sendData(data: string): void {
    // TODO: do something better than just swallowing
    // the data if the socket is not in a working condition
    if (this._socket.readyState !== 1) {
      return;
    }

    let execStreamRequest = new ExecStreamRequest();
    let input = new ExecStreamRequest.Input();
    input.setData(this.encoder.encode(data));
    execStreamRequest.setInput(input);
    this._socket.send(execStreamRequest.serializeBinary());
  }

  private _sendBinary(data: string): void {
    if (this._socket.readyState !== 1) {
      return;
    }
    let buffer = new Uint8Array(data.length);
    for (let i = 0; i < data.length; ++i) {
      buffer[i] = data.charCodeAt(i) & 255;
    }
    let execStreamRequest = new ExecStreamRequest();
    let input = new ExecStreamRequest.Input();
    input.setData(buffer);
    execStreamRequest.setInput(input);
    this._socket.send(execStreamRequest.serializeBinary());
  }

  sendStart(deploymentId: string, terminal: Terminal): void {
    let execStreamStartRequest = new ExecStreamRequest();
    let start = new ExecStreamRequest.Start();
    // Important: ArgsList can't be empty
    start.setArgsList(['/bin/bash']);
    let streaminput = new ExecStreamRequest.Input();
    execStreamStartRequest.setInput(streaminput);
    start.setDeploymentId(deploymentId);
    let pty = new ExecStreamRequest.PTY();
    pty.setTerm('bash');
    if (terminal.element) {
      let windowSize = new ExecStreamRequest.WindowSize();
      windowSize.setCols(terminal.cols);
      windowSize.setRows(terminal.rows);
      windowSize.setHeight(terminal.element.offsetHeight);
      windowSize.setWidth(terminal.element.offsetWidth);
      pty.setWindowSize(windowSize);
    }
    pty.setEnable(true);
    start.setPty(pty);
    execStreamStartRequest.setStart(start);
    // Send start message
    this._socket.send(execStreamStartRequest.serializeBinary());
  }

  setupWinch(terminal: Terminal): void {
    if (terminal.element) {
      // setup winch
      let execStreamWinchRequest = new ExecStreamRequest();
      let windowSize = new ExecStreamRequest.WindowSize();
      windowSize.setCols(terminal.cols);
      windowSize.setRows(terminal.rows);
      windowSize.setHeight(terminal.element?.offsetHeight);
      windowSize.setWidth(terminal.element?.offsetWidth);
      execStreamWinchRequest.setWinch(windowSize);
      this._socket.send(execStreamWinchRequest.serializeBinary());
    }
  }
}

function addSocketListener<K extends keyof WebSocketEventMap>(
  socket: WebSocket,
  type: K,
  handler: (this: WebSocket, ev: WebSocketEventMap[K]) => unknown
): IDisposable {
  socket.addEventListener(type, handler);
  return {
    dispose: () => {
      if (!handler) {
        // Already disposed
        return;
      }
      socket.removeEventListener(type, handler);
    },
  };
}
