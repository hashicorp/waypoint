import Component from '@glimmer/component';
import { action } from '@ember/object';

interface Args {
  json?: string;
}

export default class extends Component<Args> {
  get prettyJSON(): string {
    let uglyJSON = this.args.json ?? '{}';
    let data = JSON.parse(uglyJSON);
    let result = JSON.stringify(data, null, 2);

    return result;
  }

  get codeMirrorOptions(): Record<string, unknown> {
    return {
      mode: 'json',
      readOnly: true,
      viewportMargin: Infinity,
    };
  }

  @action
  onScroll(event: Event): void {
    event.stopImmediatePropagation();
  }
}
