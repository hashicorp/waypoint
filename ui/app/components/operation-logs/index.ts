/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import * as AnsiColors from 'ansi-colors';

import { GetJobStreamRequest, GetJobStreamResponse } from 'waypoint-pb';
import { WaypointClient } from 'waypoint-client';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { Status } from 'grpc-web';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

interface OperationLogsArgs {
  jobId: string;
}

type JobStream = ReturnType<WaypointClient['getJobStream']>;

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
  @tracked hasLogs: boolean;
  @tracked terminal!: Terminal;

  stream?: JobStream;
  currentJobId?: string;

  constructor(owner: unknown, args: OperationLogsArgs) {
    super(owner, args);
    this.hasLogs = false;
    this.terminal = createTerminal({ inputDisabled: true });

    this.startTerminalStream();
  }

  @action
  changeJob(): void {
    if (this.args.jobId === this.currentJobId) {
      return;
    }
    this.cleanUpStream();
    this.terminal.clear();
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
    // Push job errors to build logs
    if (event == GetJobStreamResponse.EventCase.STATE) {
      if (response?.getState()?.toObject()?.job?.error) {
        let err = response?.getState()?.toObject()?.job?.error?.message || '';
        this.terminal.writeln(STYLE_TO_ANSI['error-bold'](err));
      }
    }
  }

  writeLine(line: GetJobStreamResponse.Terminal.Event.Line): void {
    let msg = line.getMsg();
    let formattedMsg = (STYLE_TO_ANSI[line.toObject().style] || STYLE_TO_ANSI.default)(msg);
    this.terminal.writeln(formattedMsg);
  }

  setHasLogs(): void {
    if (!this.hasLogs) {
      this.hasLogs = true;
    }
  }

  async startTerminalStream(): Promise<void> {
    this.currentJobId = this.args.jobId;

    let req = new GetJobStreamRequest();
    req.setJobId(this.args.jobId);

    this.stream = this.api.client.getJobStream(req, this.api.WithMeta());

    this.stream.on('data', this.onData);
    this.stream.on('status', this.onStatus);
  }

  onData = (response: GetJobStreamResponse): void => {
    this.setHasLogs();
    this.writeTerminalOutput(response);
  };

  onStatus = (status: Status): void => {
    this.setHasLogs();
    if (status.details) {
      this.terminal.writeln(status.details);
    }
  };

  cleanUpStream(): void {
    if (this.stream) {
      this.stream.cancel();
      this.stream = undefined;
      this.currentJobId = undefined;
    }
  }

  willDestroy(): void {
    super.willDestroy();

    this.cleanUpStream();
  }
}
