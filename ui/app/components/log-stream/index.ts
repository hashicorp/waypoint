import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import { inject as service } from '@ember/service';
import { GetLogStreamRequest, LogBatch } from 'waypoint-pb';

interface LogStreamArgs {
  req: GetLogStreamRequest;
}

export default class LogStream extends Component<LogStreamArgs> {
  @service api!: ApiService;
  @tracked lines: string[];

  constructor(owner: any, args: any) {
    super(owner, args);
    this.start();
    this.lines = [];
  }

  addLine(line: string) {
    this.lines = [...this.lines, line];
  }

  async start() {
    const onData = (response: LogBatch) => {
      response.getLinesList().forEach((entry) => {
        this.addLine(entry.getLine());
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
