import Component from '@glimmer/component';
import KEYS from 'waypoint/utils/keys';
import SessionService from 'waypoint/services/session';
import { Terminal } from 'xterm';
import config from 'waypoint/config/environment';
import { createTerminal } from 'waypoint/utils/create-terminal';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

interface ExecComponentArgs {
  deploymentId: string;
}

const BACKSPACE_ONE_CHARACTER = '\x08 \x08';

// eslint-disable-next-line no-control-regex
const UNPRINTABLE_CHARACTERS_REGEX = /[\x00-\x1F]/g;

export default class ExecComponent extends Component<ExecComponentArgs> {
  @service session!: SessionService;

  @tracked deploymentId: string;
  @tracked terminal!: Terminal;
  @tracked command!: string;

  constructor(owner: unknown, args: ExecComponentArgs) {
    super(owner, args);
    let { deploymentId } = this.args;
    this.deploymentId = deploymentId;

    this.terminal = createTerminal({ inputDisabled: false });
    this.startExecStream(deploymentId);
  }

  async startExecStream(deploymentId: string): void {
    let protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let socket = new WebSocket(`ws://localhost:9702/v1/exec`, );
    socket.addEventListener('open', (event) => {
      socket.send(JSON.stringify({ version: 1, auth_token: this.token || '' }));
      console.log(event);
      socket.send('Hello Server!');
    });
    socket.addEventListener('message', (event) => {
      console.log(event);
    });
    socket.addEventListener('close', (event) => {
      console.log(event);
    });
    socket.addEventListener('error', (event) => {
      console.log(event);
    });
    // this.dataListener = this.terminal.onData((data) => {
    //   this.handleDataEvent(data);
    // });
  }

  async openSocketStream() {
    // Todo: handle different hosts/ports
    let protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let socket = new WebSocket(`ws://localhost:9702/v1/exec`, []);
    return socket;
  }

  handleDataEvent = (data) => {
    console.log(data);
    if (
      data === KEYS.LEFT_ARROW ||
      data === KEYS.UP_ARROW ||
      data === KEYS.RIGHT_ARROW ||
      data === KEYS.DOWN_ARROW
    ) {
      // Ignore arrow keys
    } else if (data === KEYS.DELETE) {
      if (this.command.length > 0) {
        this.terminal.write(BACKSPACE_ONE_CHARACTER);
        this.command = this.command.slice(0, -1);
      }
    } else if (data.length > 0) {
      let strippedData = data.replace(UNPRINTABLE_CHARACTERS_REGEX, '');
      this.terminal.write(strippedData);
      this.command = `${this.command}${strippedData}`;
    }
  }
}
