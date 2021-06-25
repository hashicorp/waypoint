import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import { inject as service } from '@ember/service';

import { GetJobStreamRequest, GetJobStreamResponse } from 'waypoint-pb';
import { ITerminalOptions, Terminal } from 'xterm';

const ANSI_UI_GRAY_400 = '\x1b[38;2;142;150;163m';
const ANSI_WHITE = '\x1b[0m';

interface LogTerminalArgs {
  inputDisabled: boolean;
  jobId: string;
}

export default class LogTerminal extends Component<LogTerminalArgs> {
  @service api!: ApiService;
  terminal: any;
  inputDisabled: boolean;
  jobId: string;

  constructor(owner: any, args: any) {
    super(owner, args);
    let { inputDisabled, jobId } = args;
    this.jobId = jobId;
    this.inputDisabled = inputDisabled;

    let terminalOptions: ITerminalOptions = {
      fontFamily: 'monospace',
      fontWeight: '400',
      logLevel: 'debug',
    };

    if (this.inputDisabled) {
      terminalOptions.disableStdin = true;
      terminalOptions.cursorBlink = false;
      terminalOptions.cursorStyle = 'bar';
      terminalOptions.cursorWidth = 1;
    }
    let terminal = new Terminal(terminalOptions);
    this.terminal = terminal;
    this.start();
  }

  didInsertNode(element) {
    this.terminal.open(element);
    this.terminal.write(ANSI_UI_GRAY_400);
    this.terminal.writeln('Welcome to Waypoint...');
  }

  willDestroyNode() {
    this.terminal.dispose();
  }

  writeTerminalOutput(response: GetJobStreamResponse) {
    let event = response.getEventCase();
    if (event == GetJobStreamResponse.EventCase.TERMINAL) {
      let terminal = response.getTerminal();
      if (!terminal) {
        this.terminal.writeln('status', { msg: 'Logs are no longer available for this operation' });
      } else {
        terminal.getEventsList().forEach((event) => {
          let line = event.getLine();
          let step = event.getStep();
          if (line && line.getMsg()) {
            this.terminal.writeln(line.getMsg());
          }

          if (step && step.getOutput()) {
            let newStep = step.toObject();

            if (step.getOutput_asU8().length > 0) {
              newStep.output = new TextDecoder().decode(step.getOutput_asU8());
            }

            this.terminal.writeUtf8(step.getOutput_asU8());
          }
        });
      }
    }
  }

  onData = (response: GetJobStreamResponse) => {
    this.writeTerminalOutput(response);
  }

  onStatus = (status: any) => {
    if (status.details) {
      this.terminal.writeln(status);
    }
  }

  async start() {
    let req = new GetJobStreamRequest();
    req.setJobId(this.jobId);
    let stream = this.api.client.getJobStream(req, this.api.WithMeta());

    stream.on('data', this.onData);
    stream.on('status', this.onStatus);
  }
}
