import {
  AuthMethod,
  GetAuthMethodRequest,
  GetOIDCAuthURLRequest,
  ListOIDCAuthMethodsResponse,
  Ref,
} from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';

interface OIDCAuthButtonsArgs {
  model: ListOIDCAuthMethodsResponse.AsObject;
}

export default class OIDCAuthButtonsComponent extends Component<OIDCAuthButtonsArgs> {
  @tracked model!: ListOIDCAuthMethodsResponse.AsObject;
  @service api!: ApiService;

  @task
  async initializeOIDCFlow(authMethodProvider: AuthMethod.AsObject): Promise<void> {
    let authMethodProviderName = authMethodProvider.name;
    let urlRequest = new GetOIDCAuthURLRequest();
    let authMethodrequest = new GetAuthMethodRequest();

    let authMethodRef = new Ref.AuthMethod();
    authMethodRef.setName(authMethodProviderName);
    authMethodrequest.setAuthMethod(authMethodRef);
    urlRequest.setAuthMethod(authMethodRef);
    let randomArray = new Uint32Array(10);
    window.crypto.getRandomValues(randomArray);
    let nonce = randomArray.join('').slice(0, 20);
    urlRequest.setNonce(nonce);
    window.localStorage.setItem('waypointOIDCNonce', nonce);
    window.localStorage.setItem('waypointOIDCAuthMethod', authMethodProviderName);
    let redirectUri = `${window.location.origin}/auth/oidc-callback`;
    urlRequest.setRedirectUri(redirectUri);
    let authUrl = await this.api.client.getOIDCAuthURL(urlRequest, this.api.WithMeta());
    await window.location.replace(authUrl.getUrl());
  }
}
