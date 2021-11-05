import { CompleteOIDCAuthRequest, Ref } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Route from '@ember/routing/route';
import SessionService from 'ember-simple-auth/services/session';
import { parseResponse } from 'waypoint/authenticators/oidc';
import { inject as service } from '@ember/service';

export default class AuthIndex extends Route {
  @service session!: SessionService;
  @service api!: ApiService;

  async beforeModel() {
    let oidcParams = parseResponse(window.location.search);

    console.log(oidcParams);
    let oidcModel = this.modelFor('auth.oidc');
    console.log(oidcModel);
    let completeAuthRequest = new CompleteOIDCAuthRequest();
    completeAuthRequest.setCode(oidcParams.code);
    let authMethodName = window.localStorage.getItem('waypointOIDCAuthMethod');
    let authMethodRef = new Ref.AuthMethod();
    authMethodRef.setName(authMethodName);
    completeAuthRequest.setAuthMethod(authMethodRef);
    completeAuthRequest.setRedirectUri(window.location.origin + window.location.pathname);
    let nonce = window.localStorage.getItem('waypointOIDCNonce');
    console.log(nonce);
    completeAuthRequest.setNonce(nonce);
    completeAuthRequest.setState(oidcParams.state);
    let resp = await this.api.client.completeOIDCAuth(completeAuthRequest, this.api.WithMeta());
    let respObject = resp.toObject();
    this.session.authenticate('authenticator:oidc', respObject);
  }
}
