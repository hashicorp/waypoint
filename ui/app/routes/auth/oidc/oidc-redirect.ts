import { CompleteOIDCAuthRequest, Ref } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Route from '@ember/routing/route';
import SessionService from 'waypoint/services/session';
import { inject as service } from '@ember/service';

export default class AuthOIDCRedirect extends Route {
  @service session!: SessionService;
  @service api!: ApiService;


  async model(params) {
    let oidcParams = this.paramsFor('auth.oidc');
    console.log(oidcParams);
    let oidcModel = this.modelFor('auth.oidc');
    console.log(oidcModel);
    let completeAuthRequest = new CompleteOIDCAuthRequest();
    completeAuthRequest.setCode(params.code);
    let authMethodRef = new Ref.AuthMethod();
    authMethodRef.setName(oidcParams.provider_name);
    completeAuthRequest.setAuthMethod(authMethodRef);
    completeAuthRequest.setRedirectUri(window.location.origin + window.location.pathname);
    let randomArray = new Uint32Array(10);
    window.crypto.getRandomValues(randomArray);
    let nonce = randomArray.join('').slice(0, 20);
    console.log(nonce);
    completeAuthRequest.setNonce(nonce);
    completeAuthRequest.setState(params.state);
    let resp = await this.api.client.completeOIDCAuth(completeAuthRequest, this.api.WithMeta());
    console.log(resp.toObject());
  }

  // example query params:
  // state=st_gyOv0ngwidW7YcMIR0qQ&code=4/0AX4XfWiuIlC4c5s0_p-lKyxXYPLNfrmrDmajR_8p648XiYPvjHWIArLOAgXMhkeQWKRidw&scope=openid&authuser=0&prompt=consent
}
