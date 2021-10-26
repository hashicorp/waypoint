import { GetAuthMethodRequest, GetOIDCAuthURLRequest, Ref } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Route from '@ember/routing/route';
import SessionService from 'waypoint/services/session';
import { inject as service } from '@ember/service';

export default class AuthOIDCRedirect extends Route {
  @service session!: SessionService;
  @service api!: ApiService;

  async model() {
    let params = this.paramsFor('auth.oidc');
    let urlRequest = new GetOIDCAuthURLRequest();
    let authMethodrequest = new GetAuthMethodRequest();

    let authMethodRef = new Ref.AuthMethod();
    authMethodRef.setName(params.provider_name);
    authMethodrequest.setAuthMethod(authMethodRef);
    let authMethod = await this.api.client.getAuthMethod(authMethodrequest, this.api.WithMeta());
    urlRequest.setAuthMethod(authMethodRef);
    let randomArray = new Uint32Array(10);
    window.crypto.getRandomValues(randomArray);
    let nonce = randomArray.join('').slice(0, 20);
    urlRequest.setNonce(nonce);
    window.localStorage.setItem('waypointOIDCNonce', nonce);
    let redirectUri = `${window.location.origin}/auth/${params.provider_name}/oidc-redirect`;
    urlRequest.setRedirectUri(redirectUri);
    let authUrl = await this.api.client.getOIDCAuthURL(urlRequest, this.api.WithMeta());
    window.location.replace(authUrl.getUrl());
  }
}
