import ApiService from 'waypoint/services/api';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { ListOIDCAuthMethodsResponse } from 'waypoint-pb';
import Route from '@ember/routing/route';
import { SessionService } from 'ember-simple-auth/services/session';
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

  async model(): Promise<ListOIDCAuthMethodsResponse.AsObject | undefined> {
    if (this.authMethods.getAuthMethodsList().length) {
      let providers = this.authMethods.toObject();
      return providers;
    } else {
      return;
    }
  }
}
