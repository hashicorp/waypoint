import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { OnDemandRunnerConfig } from 'waypoint-pb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

export type Model = OnDemandRunnerConfig.AsObject[];

export default class extends Route {
  @service api!: ApiService;

  async model(): Promise<Model> {
    let empty = new Empty();
    let meta = this.api.WithMeta();
    let response = await this.api.client.listOnDemandRunnerConfigs(empty, meta);

    return response.toObject().configsList;
  }
}
