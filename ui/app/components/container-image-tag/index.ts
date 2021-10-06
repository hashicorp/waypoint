import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { StatusReport } from 'waypoint-pb';
import { ImageRef, findImageRefs } from 'waypoint/utils/image-refs';

interface Args {
  statusReport: StatusReport.AsObject;
}

export default class extends Component<Args> {
  @service api!: ApiService;

  get states(): unknown {
    return this.args.statusReport.resourcesList
      ? this.args.statusReport.resourcesList.map((r) => JSON.parse(r.stateJson ?? '{}'))
      : [];
  }

  get imageRefs(): ImageRef[] {
    return findImageRefs(this.states);
  }
}
