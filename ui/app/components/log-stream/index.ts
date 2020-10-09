import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import { inject as service } from '@ember/service';
import { GetLogStreamRequest, LogBatch } from 'waypoint-pb';
import { formatRFC3339 } from 'date-fns';

interface LogStreamArgs {
  req: GetLogStreamRequest;
}

export default class LogStream extends Component<LogStreamArgs> {
  @service api!: ApiService;
  @tracked lines: string[];

  constructor(owner: any, args: any) {
    super(owner, args);
    this.lines = [];
    this.start();
  }

  addLine(line: string) {
    this.lines = [...this.lines, line];
  }

  async start() {
    const onData = (response: LogBatch) => {
      response.getLinesList().forEach((entry) => {
        const prefix = formatRFC3339(entry.getTimestamp()!.toDate());
        this.addLine(`${prefix}: ${entry.getLine()}`);
      });
    };

    const onStatus = (status: any) => {
      this.addLine(status.details);
    };

    var stream = this.api.client.getLogStream(this.args.req, this.api.WithMeta());

    stream.on('data', onData);
    stream.on('status', onStatus);
  }
}
