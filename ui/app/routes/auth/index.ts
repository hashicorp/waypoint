import ApiService from 'waypoint/services/api';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { ListOIDCAuthMethodsResponse } from 'waypoint-pb';
import Route from '@ember/routing/route';
import SessionService from 'waypoint/services/session';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class AuthIndex extends Route {
  @service session!: SessionService;
  @service api!: ApiService;
  @tracked authMethods!: ListOIDCAuthMethodsResponse;

  async beforeModel(): Promise<void> {
    let authMethods = await this.api.client.listOIDCAuthMethods(new Empty(), this.api.WithMeta());
    this.authMethods = authMethods;
    return;
  }
  // todo: move this to beforeModel then use model hook to get the primary OIDC provider and render it
  async model(): Promise<ListOIDCAuthMethodsResponse.AsObject | undefined> {
    if (this.authMethods.getAuthMethodsList().length) {
      let providers = this.authMethods.toObject();
      return providers;
    } else {
      return;
    }
  }
}
