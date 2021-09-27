import { GetLogStreamRequest, LogBatch } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/terminal';
import { formatRFC3339 } from 'date-fns';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface LogStreamArgs {
  req: GetLogStreamRequest;
}

export default class LogStream extends Component<LogStreamArgs> {
  @service api!: ApiService;

  @tracked terminal!: Terminal;


  constructor(owner: unknown, args: LogStreamArgs) {
    super(owner, args);
    this.terminal = createTerminal({ inputDisabled: true });
    this.startTerminalStream();
  }

  async startTerminalStream(): Promise<void> {
    let stream = this.api.client.getLogStream(this.args.req, this.api.WithMeta());

    stream.on('data', this.onData);
    stream.on('status', this.onStatus);
  }

  onData = (response: LogBatch): void => {
    response.getLinesList().forEach((entry) => {
      let ts = entry.getTimestamp();
      if (!ts) {
        return;
      }
      let prefix = formatRFC3339(ts.toDate());
      this.terminal.writeln(`${prefix}: ${entry.getLine()}`);
    });
  };

  onStatus = (status: any): void => {
    if (status.details) {
      this.terminal.writeln(status.details);
    }
  };
}
