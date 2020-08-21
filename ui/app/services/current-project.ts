import Service from '@ember/service';
import { Project, Ref } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

export default class CurrentProjectService extends Service {
  @tracked project?: Project;
  @tracked ref?: Ref.Project;

  get name() {
    return this.project?.getName() || this.ref?.getProject();
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    currentProject: CurrentProjectService;
  }
}
