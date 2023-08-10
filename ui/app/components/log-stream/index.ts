/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { GetLogStreamRequest, LogBatch } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { Status } from 'grpc-web';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { formatRFC3339 } from 'date-fns';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface LogStreamArgs {
  req: GetLogStreamRequest;
}

export default class LogStream extends Component<LogStreamArgs> {
  @service api!: ApiService;

  @tracked terminal!: Terminal;
  @tracked hasLogs: boolean;

  constructor(owner: unknown, args: LogStreamArgs) {
    super(owner, args);
    this.hasLogs = false;
    this.terminal = createTerminal({ inputDisabled: true });
    this.startTerminalStream();
  }

  async startTerminalStream(): Promise<void> {
    let stream = this.api.client.getLogStream(this.args.req, this.api.WithMeta());

    stream.on('data', this.onData);
    stream.on('status', this.onStatus);
  }

  setHasLogs(): void {
    if (!this.hasLogs) {
      this.hasLogs = true;
    }
  }

  onData = (response: LogBatch): void => {
    this.setHasLogs();
    response.getLinesList().forEach((entry) => {
      let ts = entry.getTimestamp();
      if (!ts) {
        return;
      }
      let prefix = formatRFC3339(ts.toDate());
      this.terminal.writeln(`${prefix}: ${entry.getLine()}`);
    });
  };

  onStatus = (status?: Status): void => {
    this.setHasLogs();
    if (status?.details) {
      this.terminal.writeln(status.details);
    }
  };
}
