import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import CurrentWorkspaceService from 'waypoint/services/current-workspace';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

export default class Index extends Route {
  @service api!: ApiService;
  @service currentWorkspace!: CurrentWorkspaceService;

  async model() {
    let meta = {
      authorization:
        'bM152PWkXxfoy4vA51JFhR7LsKQez6x23oi2RDqYk8DPjnRdGjWtS6J3CTywJSaBPQX7wZAgV61bFMMLWoqvpjUfr1pL2sq9AcDGL',
    };
    let resp = await this.api.client.listProjects(new Empty(), meta);
    let projects = resp.getProjectsList().map((p) => p.toObject());

    return projects;
  }
}
