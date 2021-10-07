import ApiService from 'waypoint/services/api';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import Route from '@ember/routing/route';
import SessionService from 'waypoint/services/session';
import { inject as service } from '@ember/service';

export default class Application extends Route {
  @service session!: SessionService;
  @service api!: ApiService;
  // todo: move this to beforeModel then use model hook to get the primary OIDC provider and render it
  async model() {
    let model = await this.api.client.listOIDCAuthMethods(new Empty(), this.api.WithMeta());
    return model.toObject().authMethodsList;
  }
}
