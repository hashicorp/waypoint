/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { ExecStreamRequest } from 'waypoint-pb';
import { ExecWebSocketAddon } from 'waypoint/utils/exec-websocket-addon';
import SessionService from 'ember-simple-auth/services/session';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import config from 'waypoint/config/environment';

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
    let token = this.session.data.authenticated?.token as string;
    let url = new URL(config.apiAddress || window.location.toString());
    url.protocol = 'wss:';
    url.pathname = '/v1/exec';
    url.searchParams.append('token', token);
    let socket = new WebSocket(url);
    socket.binaryType = 'arraybuffer';
    this.socket = socket;
    // The socket addon handles all terminal input/output
    let socketAddon = new ExecWebSocketAddon(socket, { bidirectional: true, deploymentId });
    this.terminal.loadAddon(socketAddon);
  }

  disconnect(): void {
    // send EOF message then close socket connection
    let execStreamRequest = new ExecStreamRequest();
    execStreamRequest.setInputEof(new Empty());
    this.socket.send(execStreamRequest.serializeBinary());
    this.socket.close();
  }

  willDestroy(): void {
    this.disconnect();
    super.willDestroy();
  }
}
