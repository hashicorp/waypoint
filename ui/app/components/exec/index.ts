import * as protobuf from 'google-protobuf';

import { ExecStreamRequest, ExecStreamResponse } from 'waypoint-pb';

import Component from '@glimmer/component';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import KEYS from 'waypoint/utils/keys';
import { Message } from 'google-protobuf';
import SessionService from 'waypoint/services/session';
import { Terminal } from 'xterm';
import config from 'waypoint/config/environment';
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
    this.startExecStream(deploymentId);
  }

  async startExecStream(deploymentId: string): void {
    let protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let socket = new WebSocket(`wss://localhost:9702/v1/exec?token=${this.session.token}`);
    socket.binaryType = 'arraybuffer';
    this.socket = socket;
    socket.addEventListener('open', (event) => {
      // socket.send(JSON.stringify({ version: 1, auth_token: this.token || '' }));
      console.log(event);
      this.sendStart(deploymentId);
      this.setupWinch();
    });
    socket.addEventListener('message', (event) => {
      console.log(event);
      let reader = new FileReader();
      let output = event.data;
      let resp = ExecStreamResponse.deserializeBinary(output);
      console.log(resp.getEventCase());
      if (resp.hasOpen()) {
        console.log(resp.getOpen()?.toString());
        this.terminal.writeln('Connected to instance...');
      }
      if (resp.getEventCase() === 3) {
        console.log(resp.getOutput()?.getChannel());
        console.log(resp.getOutput()?.getData_asU8());
        console.log(resp.getOutput()?.toObject());
      }

      if (resp.getEventCase() === 1) {
        console.log(resp.getOutput()?.getChannel());
        console.log(resp.getOutput()?.getData_asU8());
        console.log(resp.getOutput()?.toObject());
        this.terminal.writeUtf8(resp.getOutput()?.getData_asU8());
      }
      console.log(resp.toObject());
    });
    socket.addEventListener('close', (event) => {
      let reader = new FileReader();
      let output = event.data;
      let resp = ExecStreamResponse.deserializeBinary(output);
      console.log(resp.getEventCase());
      console.log(resp.toObject());
    });
    socket.addEventListener('error', (event) => {
      let reader = new FileReader();
      let output = event.data;
      let resp = ExecStreamResponse.deserializeBinary(output);
      console.log(resp.getEventCase());
      console.log(resp.toObject());
    });
    this.terminal.onData((data) => {
      this.handleDataEvent(data);
    });
  }

  async openSocketStream() {
    // Todo: handle different hosts/ports
    let protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let socket = new WebSocket(`ws://localhost:9702/v1/exec`, ['']);
    return socket;
  }

  sendCommand() {
    let execStreamRequest = new ExecStreamRequest();
    let input = new ExecStreamRequest.Input();
    let encoder = new TextEncoder();
    input.setData(encoder.encode(`${this.command}\r`));
    execStreamRequest.setInput(input);
    this.socket.send(execStreamRequest.serializeBinary());
    this.command = '';
  }

  sendStart(deploymentId) {
    let execStreamStartRequest = new ExecStreamRequest();
    let start = new ExecStreamRequest.Start();
    // Important: ArgsList can't be empty
    start.setArgsList(['/bin/bash']);
    let streaminput = new ExecStreamRequest.Input();
    execStreamStartRequest.setInput(streaminput);
    start.setDeploymentId(deploymentId);
    let pty = new ExecStreamRequest.PTY();
    pty.setTerm('bash');
    let windowSize = new ExecStreamRequest.WindowSize();
    windowSize.setCols(this.terminal.cols);
    windowSize.setRows(this.terminal.rows);
    windowSize.setHeight(this.terminal.element?.offsetHeight);
    windowSize.setWidth(this.terminal.element?.offsetWidth);
    pty.setWindowSize(windowSize);
    pty.setEnable(true);
    start.setPty(pty);
    execStreamStartRequest.setStart(start);
    // Send start message
    this.socket.send(execStreamStartRequest.serializeBinary());
  }

  setupWinch() {
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

  disconnect() {
    let execStreamRequest = new ExecStreamRequest();
    execStreamRequest.setInputEof(new Empty());
    this.socket.send(execStreamRequest.serializeBinary());
  }

  willDestroy(): void {
    this.disconnect();
    super.willDestroy();
  }

  handleDataEvent = (data) => {
    console.log(data);
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
      this.terminal.writeln('');
      this.sendCommand();
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
  };
}
