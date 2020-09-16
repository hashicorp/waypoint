import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import { inject as service } from '@ember/service';
import { GetJobStreamRequest, GetJobStreamResponse } from 'waypoint-pb';

interface OperationLogsArgs {
  jobId: string;
}

export default class OperationLogs extends Component<OperationLogsArgs> {
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
    const onData = (response: GetJobStreamResponse) => {
      let terminal = response.getTerminal();

      if (!terminal) {
        this.addLine('Logs are no longer available for this operation');
      } else {
        terminal.getEventsList().forEach((event) => {
          event.getLine();
        });
      }
    };

    const onStatus = (status: any) => {
      this.addLine(status.details);
    };

    let req = new GetJobStreamRequest();
    req.setJobId(this.args.jobId);
    let stream = this.api.client.getJobStream(req, this.api.WithMeta());

    stream.on('data', onData);
    stream.on('status', onStatus);
  }
}
