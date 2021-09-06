import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import ApiService from 'waypoint/services/api';
import { inject as service } from '@ember/service';
import { GetLogStreamRequest, LogBatch } from 'waypoint-pb';
import { Status } from 'grpc-web';
import { formatRFC3339 } from 'date-fns';

interface LogStreamArgs {
  req: GetLogStreamRequest;
}

export default class LogStream extends Component<LogStreamArgs> {
  @service api!: ApiService;

  @tracked lines: string[];
  @tracked isFollowingLogs = true;
  @tracked badgeCount = 0;

  constructor(owner: unknown, args: LogStreamArgs) {
    super(owner, args);
    this.lines = [];
    this.start();
  }

  addLine(line: string): void {
    this.lines = [...this.lines, line];
    if (this.isFollowingLogs === false) {
      this.badgeCount = this.badgeCount + 1;
    }
  }

  @action
  followLogs(element: HTMLElement | Event): void {
    if (element instanceof Event) {
      if (element.target instanceof HTMLElement) {
        element = element.target;
      } else {
        return;
      }
    }

    let scrollableElement = element.closest('.output-scroll-y');

    if (!scrollableElement) {
      return;
    }

    scrollableElement.scroll(0, scrollableElement.scrollHeight);
    this.badgeCount = 0;
  }

  @action
  updateScroll(element: HTMLElement): void {
    if (this.isFollowingLogs === true) {
      element.scrollIntoView(false);
      this.badgeCount = 0;
    }
  }

  async start(): Promise<void> {
    let onData = (response: LogBatch) => {
      response.getLinesList().forEach((entry) => {
        let ts = entry.getTimestamp();
        if (!ts) {
          return;
        }
        let prefix = formatRFC3339(ts.toDate());
        this.addLine(`${prefix}: ${entry.getLine()}`);
      });
    };

    let onStatus = (status: Status) => {
      this.addLine(status.details);
    };

    let stream = this.api.client.getLogStream(this.args.req, this.api.WithMeta());

    stream.on('data', onData);
    stream.on('status', onStatus);
  }
}
