import * as AnsiColors from 'ansi-colors';

import { GetJobStreamRequest, GetJobStreamResponse } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface OperationLogsArgs {
  jobId: string;
}

export default class OperationLogs extends Component<OperationLogsArgs> {
  @service api!: ApiService;

  @tracked terminal!: Terminal;

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
    this.startTerminalStream();
  }

  writeTerminalOutput(response: GetJobStreamResponse): void {
    let event = response.getEventCase();
    if (event == GetJobStreamResponse.EventCase.TERMINAL) {
      let terminalOutput = response.getTerminal();
      if (!terminalOutput) {
        this.terminal.writeln('Logs are no longer available for this operation');
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
      this.terminal.writeln(status.details);
    }
  };
}
