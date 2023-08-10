/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { CompleteOIDCAuthRequest, Ref } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Route from '@ember/routing/route';
import SessionService from 'ember-simple-auth/services/session';
import { parseResponse } from 'waypoint/authenticators/oidc';
import { inject as service } from '@ember/service';

export default class AuthIndex extends Route {
  @service session!: SessionService;
  @service api!: ApiService;

  async model(): Promise<void> {
    let oidcParams = parseResponse(window.location.search);
    let completeAuthRequest = new CompleteOIDCAuthRequest();
    completeAuthRequest.setCode(oidcParams.code);
    let authMethodName = window.localStorage.getItem('waypointOIDCAuthMethod');
    let authMethodRef = new Ref.AuthMethod();
    if (authMethodName) {
      authMethodRef.setName(authMethodName);
    }
    completeAuthRequest.setAuthMethod(authMethodRef);
    completeAuthRequest.setRedirectUri(window.location.origin + window.location.pathname);
    let nonce = window.localStorage.getItem('waypointOIDCNonce');
    if (nonce) {
      completeAuthRequest.setNonce(nonce);
    }
    completeAuthRequest.setState(oidcParams.state);
    let resp = await this.api.client.completeOIDCAuth(completeAuthRequest, this.api.WithMeta());
    let respObject = resp.toObject();
    await this.session.authenticate('authenticator:oidc', respObject);
  }
}
