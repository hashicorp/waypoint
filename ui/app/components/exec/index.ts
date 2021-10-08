import Component from '@glimmer/component';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { ExecStreamRequest } from 'waypoint-pb';
import { ExecWebSocketAddon } from 'waypoint/utils/exec-websocket-addon';
import SessionService from 'ember-simple-auth/services/session';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface ExecComponentArgs {
  deploymentId: string;
}

export default class ExecComponent extends Component<ExecComponentArgs> {
  @service session!: SessionService;

  @tracked deploymentId: string;
  @tracked terminal!: Terminal;
  @tracked socket!: WebSocket;

  constructor(owner: unknown, args: ExecComponentArgs) {
    super(owner, args);
    let { deploymentId } = this.args;
    this.deploymentId = deploymentId;
    this.terminal = createTerminal({ inputDisabled: false });
    this.terminal.focus();
    this.startExecStream(deploymentId);
  }

  async startExecStream(deploymentId: string): Promise<void> {
    let token = this.session.data.authenticated?.token;
    let apiHost = window.location.hostname;
    let socket = new WebSocket(`wss://${apiHost}:9702/v1/exec?token=${token}`);
    socket.binaryType = 'arraybuffer';
    this.socket = socket;
    // The socket addon handles all terminal input/output
    let socketAddon = new ExecWebSocketAddon(socket, { bidirectional: true, deploymentId });
    this.terminal.loadAddon(socketAddon);
  }

  disconnect(): void {
    let execStreamRequest = new ExecStreamRequest();
    execStreamRequest.setInputEof(new Empty());
    this.socket.send(execStreamRequest.serializeBinary());
  }

  willDestroy(): void {
    this.disconnect();
    super.willDestroy();
  }
}
