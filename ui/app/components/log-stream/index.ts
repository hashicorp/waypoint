import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

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
  @tracked isFollowingLogs = true;
  @tracked badgeCount = 0;

  constructor(owner: any, args: any) {
    super(owner, args);
    this.lines = [];
    this.start();
  }

  addLine(line: string) {
    this.lines = [...this.lines, line];
    if (this.isFollowingLogs === false) {
      this.badgeCount = this.badgeCount + 1;
    }
  }

  @action
  followLogs(element: any) {
    let scrollableElement = element.target
      ? element.target.closest('.output-scroll-y')
      : element.closest('.output-scroll-y');

    scrollableElement.scroll(0, scrollableElement.scrollHeight);
    this.badgeCount = 0;
  }

  @action
  updateScroll(element: any) {
    if (this.isFollowingLogs === true) {
      element.scrollIntoView(false);
      this.badgeCount = 0;
    }
  }

  async start() {
    let onData = (response: LogBatch) => {
      response.getLinesList().forEach((entry) => {
        let prefix = formatRFC3339(entry.getTimestamp()!.toDate());
        this.addLine(`${prefix}: ${entry.getLine()}`);
      });
    };

    let onStatus = (status: any) => {
      this.addLine(status.details);
    };

    let stream = this.api.client.getLogStream(this.args.req, this.api.WithMeta());

    stream.on('data', onData);
    stream.on('status', onStatus);
  }
}
