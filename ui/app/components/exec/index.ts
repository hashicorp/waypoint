import Component from '@glimmer/component';
import { Terminal } from 'xterm';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { tracked } from '@glimmer/tracking';

interface ExecComponentArgs {
  deploymentId: string;
}

export default class ExecComponent extends Component<ExecComponentArgs> {
  @tracked deploymentId: string;
  @tracked terminal!: Terminal;

  constructor(owner: unknown, args: ExecComponentArgs) {
    super(owner, args);
    let { deploymentId } = this.args;
    this.deploymentId = deploymentId;

    this.terminal = createTerminal({ inputDisabled: false });
    this.startExecStream();
  }

  startExecStream(): void {
    // Todo
  }
}
