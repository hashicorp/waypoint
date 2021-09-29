import * as AnsiColors from 'ansi-colors';

import { GetJobStreamRequest, GetJobStreamResponse } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { Status } from 'grpc-web';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface OperationLogsArgs {
  jobId: string;
}

// Mappings for message styles
// https://github.com/hashicorp/waypoint-plugin-sdk/blob/baf566811af680c5df138f9915d756f67d271b1a/terminal/ui.go#L126-L135
const STYLE_TO_ANSI: Record<string, (msg: string) => string> = {
  header: AnsiColors.bold,
  error: AnsiColors.red,
  'error-bold': AnsiColors.red.bold,
  warning: AnsiColors.yellow,
  'warning-bold': AnsiColors.yellow.bold,
  info: AnsiColors.cyan,
  success: AnsiColors.green,
  'success-bold': AnsiColors.green.bold,
  '': AnsiColors.bold,
  default: (msg) => msg,
};
export default class OperationLogs extends Component<OperationLogsArgs> {
  @service api!: ApiService;

  @tracked terminal!: Terminal;

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
    let formattedMsg = (STYLE_TO_ANSI[line.toObject().style] || STYLE_TO_ANSI.default)(msg);
    this.terminal.writeln(formattedMsg);
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

  onStatus = (status: Status): void => {
    if (status.details) {
      this.terminal.writeln(status.details);
    }
  };
}
