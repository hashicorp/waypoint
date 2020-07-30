import Service, { inject as service } from '@ember/service';
import { Project, Ref } from 'waypoint-pb';

export default class CurrentProjectService extends Service {
  project?: Project;
  ref?: Ref.Project;

  setProject(project: Project) {
    this.project = project;
  }

  setRef(ref: Ref.Project) {
    this.ref = ref;
  }
}

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    currentProject: CurrentProjectService;
  }
}
