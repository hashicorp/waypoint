import Service from '@ember/service';
import codemirror from 'codemirror';

// This service chiefly exists now for testing purposes.
export default class CodeMirror extends Service {
  _instances = Object.create(null);

  instanceFor(id: string): codemirror.Editor {
    return this._instances[id];
  }

  registerInstance(id: string, instance: codemirror.Editor): codemirror.Editor {
    this._instances[id] = instance;

    return instance;
  }

  unregisterInstance(id: string): void {
    delete this._instances[id];
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    'code-mirror': CodeMirror;
  }
}
