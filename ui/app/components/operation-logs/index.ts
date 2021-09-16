import * as AnsiColors from 'ansi-colors';

import { GetJobStreamRequest, GetJobStreamResponse } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { Status } from 'grpc-web';
import { Terminal } from 'xterm';
import { action } from '@ember/object';
import { createTerminal } from 'waypoint/utils/terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface OperationLogsArgs {
  jobId: string;
}

type LogLine = Record<string, unknown>;

export default class OperationLogs extends Component<OperationLogsArgs> {
  @service api!: ApiService;

  @tracked terminal!: Terminal;
  @tracked logLines: LogLine[];
  @tracked isFollowingLogs = true;
  @tracked badgeCount = 0;

  // https://github.com/hashicorp/waypoint-plugin-sdk/blob/baf566811af680c5df138f9915d756f67d271b1a/terminal/ui.go#L126-L135
  headerStyle = 'header';
  errorStyle = 'error';
  errorBoldStyle = 'error-bold';
  warningStyle = 'warning';
  warningBoldStyle = 'warning-bold';
  infoStyle = 'info';
  successStyle = 'success';
  successBoldStyle = 'success-bold';

  typeLine = 'line';
  typeStep = 'step';
  typeStepGroup = 'step-group';
  typeStatus = 'status';

  constructor(owner: unknown, args: OperationLogsArgs) {
    super(owner, args);
    this.terminal = createTerminal({ inputDisabled: true });
    // this.logLines = [];
    this.startTerminalStream();
  }

  addLogLine(t: string, logLine: LogLine): void {
    this.logLines = [...this.logLines, { type: t, logLine: logLine }];
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
  }

  @action
  updateScroll(element: HTMLElement): void {
    if (this.isFollowingLogs === true) {
      element.scrollIntoView(false);
      this.badgeCount = 0;
    }
  }

  writeTerminalOutput(response: GetJobStreamResponse): void {
    let event = response.getEventCase();
    if (event == GetJobStreamResponse.EventCase.TERMINAL) {
      let terminalOutput = response.getTerminal();
      if (!terminalOutput) {
        this.terminal.writeln('status', { msg: 'Logs are no longer available for this operation' });
      } else {
        terminalOutput.getEventsList().forEach((event) => {
          let line = event.getLine();
          let step = event.getStep();
          if (line && line.getMsg()) {
            this.writeLine(line);
          }

          if (step && step.getOutput()) {
            let newStep = step.toObject();

            if (step.getOutput_asU8().length > 0) {
              newStep.output = new TextDecoder().decode(step.getOutput_asU8());
            }

            this.terminal.write(step.getOutput_asU8());
          }
        });
      }
    }
  }

  writeLine(line: GetJobStreamResponse.Terminal.Event.Line): void {
    let msg = line.getMsg();
    switch (line.toObject().style) {
      case 'header':
        msg = AnsiColors.bold(msg);
        break;
      case 'error':
        msg = AnsiColors.red(msg);
        break;
      case 'error-bold':
        msg = AnsiColors.bold.red(msg);
        break;
      case 'warning':
        msg = AnsiColors.yellow(msg);
        break;
      case 'warning-bold':
        msg = AnsiColors.yellow.bold(msg);
        break;
      case 'info':
        msg = AnsiColors.cyan(msg);
        break;
      case 'success':
        msg = AnsiColors.green(msg);
        break;
      case 'success-bold':
        msg = AnsiColors.green.bold(msg);
        break;
      case '':
        msg = AnsiColors.bold(msg);
        break;
      default:
        break;
    }
    this.terminal.writeln(msg);
  }

  async startTerminalStream(): Promise<void> {
    let req = new GetJobStreamRequest();
    req.setJobId(this.args.jobId);
    let stream = this.api.client.getJobStream(req, this.api.WithMeta());

    stream.on('data', this.onData);
    stream.on('status', this.onStatus);
  }

  onData = (response: GetJobStreamResponse): void => {
    this.writeTerminalOutput(response);
  };

  onStatus = (status: any): void => {
    if (status.details) {
      this.terminal.writeln(status);
    }
  };
}
