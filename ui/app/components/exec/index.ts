import * as AnsiColors from 'ansi-colors';

import { ExecStreamRequest, ExecStreamResponse } from 'waypoint-pb';

import Component from '@glimmer/component';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import KEYS from 'waypoint/utils/keys';
import SessionService from 'ember-simple-auth/services/session';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface ExecComponentArgs {
  deploymentId: string;
}

const BACKSPACE_ONE_CHARACTER = '\x08 \x08';

// eslint-disable-next-line no-control-regex
const UNPRINTABLE_CHARACTERS_REGEX = /[\x00-\x1F]/g;

export default class ExecComponent extends Component<ExecComponentArgs> {
  @service session!: SessionService;

  @tracked deploymentId: string;
  @tracked terminal!: Terminal;
  @tracked command!: string;
  @tracked socket!: WebSocket;

  constructor(owner: unknown, args: ExecComponentArgs) {
    super(owner, args);
    let { deploymentId } = this.args;
    this.deploymentId = deploymentId;

    this.command = '';

    this.terminal = createTerminal({ inputDisabled: false });
    this.startExecStream();
  }

  get hasDeployment(): boolean {
    return !!this.deploymentId;
  }

  async startExecStream(): Promise<void> {
    let token = this.session.data.authenticated?.token;
    let protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let socket = await new WebSocket(`wss://localhost:9702/v1/exec?token=${token}`);
    socket.binaryType = 'arraybuffer';
    this.socket = socket;
    this.terminal.writeln(AnsiColors.bold.cyan('Connecting...'));
    if (!this.hasDeployment) {
      this.terminal.writeln(AnsiColors.bold.red('No deployment available'));
    }
    this.setUpSocketListeners();
    this.terminal.onData((data) => {
      console.log(this.socket.readyState);
      this.handleDataEvent(data);
    });
    this.terminal.onBinary((data) => {
      console.log(this.socket.readyState);
      // this.sendData();
      this.sendBinary(data);
    });
  }

  setUpSocketListeners(): void {
    this.socket.addEventListener('open', () => {
      this.sendStart(this.deploymentId);
      this.setupWinch();
      this.terminal.focus();
    });
    this.socket.addEventListener('message', (event) => {
      let output = event.data;
      let resp = ExecStreamResponse.deserializeBinary(output);
      if (resp.hasOpen()) {
        this.terminal.clear();
        this.terminal.writeln(AnsiColors.bold.cyan(`Connecting to deployment: ${this.deploymentId}`));
      }
      // Open
      if (resp.getEventCase() === 3) {
        this.terminal.writeln(AnsiColors.bold.cyan('Connection opened'));
      }
      // Exit
      if (resp.getEventCase() === 2) {
        let exitCode = resp.getExit() || 'unknown';
        this.terminal.writeln(AnsiColors.yellow(`Exit code: ${exitCode}`));
        this.terminal.writeln(AnsiColors.yellow('Connection closed...'));
      }
      // Output
      if (resp.getEventCase() === 1) {
        this.terminal.writeUtf8(resp.getOutput()?.getData_asU8());
      }
      // Event not set
      if (resp.getEventCase() === 0) {
        this.terminal.writeln(AnsiColors.yellow('Connection closed...'));
      }
    });
    this.socket.addEventListener('close', (event) => {
      let output = event.data;
      if (output) {
        let resp = ExecStreamResponse.deserializeBinary(output);
        this.terminal.writeUtf8(resp.getOutput()?.getData_asU8());
        this.terminal.writeln(AnsiColors.yellow('Connection closed...'));
      }
    });
    this.socket.addEventListener('error', (event) => {
      let output = event.data;
      if (output) {
        let resp = ExecStreamResponse.deserializeBinary(output);
        this.terminal.writeUtf8(AnsiColors.red(resp.getOutput()?.getData_asU8()));
      }
    });
  }

  sendData(data: string): void {
    let execStreamRequest = new ExecStreamRequest();
    let input = new ExecStreamRequest.Input();
    let encoder = new TextEncoder();
    input.setData(encoder.encode(data));
    execStreamRequest.setInput(input);
    this.socket.send(execStreamRequest.serializeBinary());
  }

  sendBinary(data) {
    let buffer = new Uint8Array(data.length);
    for (let i = 0; i < data.length; ++i) {
      buffer[i] = data.charCodeAt(i) & 255;
    }
    let execStreamRequest = new ExecStreamRequest();
    let input = new ExecStreamRequest.Input();
    input.setData(buffer);
    execStreamRequest.setInput(input);
    this.socket.send(execStreamRequest.serializeBinary());
  }

  sendStart(deploymentId: string): void {
    let execStreamStartRequest = new ExecStreamRequest();
    let start = new ExecStreamRequest.Start();
    // Important: ArgsList can't be empty
    start.setArgsList(['/bin/bash']);
    let streaminput = new ExecStreamRequest.Input();
    execStreamStartRequest.setInput(streaminput);
    start.setDeploymentId(deploymentId);
    let pty = new ExecStreamRequest.PTY();
    pty.setTerm('bash');
    if (this.terminal.element) {
      let windowSize = new ExecStreamRequest.WindowSize();
      windowSize.setCols(this.terminal.cols);
      windowSize.setRows(this.terminal.rows);
      windowSize.setHeight(this.terminal.element.offsetHeight);
      windowSize.setWidth(this.terminal.element.offsetWidth);
      pty.setWindowSize(windowSize);
    }
    pty.setEnable(true);
    start.setPty(pty);
    execStreamStartRequest.setStart(start);
    // Send start message
    this.socket.send(execStreamStartRequest.serializeBinary());
  }

  setupWinch(): void {
    if (this.terminal.element) {
      // setup winch
      let execStreamWinchRequest = new ExecStreamRequest();
      let windowSize = new ExecStreamRequest.WindowSize();
      windowSize.setCols(this.terminal.cols);
      windowSize.setRows(this.terminal.rows);
      windowSize.setHeight(this.terminal.element?.offsetHeight);
      windowSize.setWidth(this.terminal.element?.offsetWidth);
      execStreamWinchRequest.setWinch(windowSize);
      this.socket.send(execStreamWinchRequest.serializeBinary());
    }
  }

  disconnect(): void {
    let execStreamRequest = new ExecStreamRequest();
    execStreamRequest.setInputEof(new Empty());
    this.socket.send(execStreamRequest.serializeBinary());
    this.socket.close();
  }

  willDestroy(): void {
    this.disconnect();
    super.willDestroy();
  }

  handleDataEvent = (data: string): void => {
    console.log(this.command);
    console.log(this.terminal.buffer);
    // if (
    //   data === KEYS.LEFT_ARROW ||
    //   data === KEYS.UP_ARROW ||
    //   data === KEYS.RIGHT_ARROW ||
    //   data === KEYS.DOWN_ARROW
    // ) {
    //   // Ignore arrow keys
    // } else
    if (data === KEYS.CONTROL_U) {
      this.terminal.write(BACKSPACE_ONE_CHARACTER.repeat(this.command.length));
      this.command = '';
    } else if (data === KEYS.ENTER) {
      this.sendData(data);
      this.command = '';
    } else if (data === KEYS.DELETE) {
      if (this.command.length > 0) {
        this.terminal.write(BACKSPACE_ONE_CHARACTER);
        this.command = this.command.slice(0, -1);
      }
    } else if (data.length > 0) {
      // let strippedData = data.replace(UNPRINTABLE_CHARACTERS_REGEX, '');
      this.command = `${this.command}${data}`;
      this.sendData(data);
    }
    console.log(this.command);
  };
}
