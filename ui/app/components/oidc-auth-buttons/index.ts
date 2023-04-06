/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import {
  AuthMethod,
  GetAuthMethodRequest,
  GetOIDCAuthURLRequest,
  ListOIDCAuthMethodsResponse,
  Ref,
} from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import Component from '@glimmer/component';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

interface OIDCAuthButtonsArgs {
  model: ListOIDCAuthMethodsResponse.AsObject;
}

export default class OIDCAuthButtonsComponent extends Component<OIDCAuthButtonsArgs> {
  model!: ListOIDCAuthMethodsResponse.AsObject;
  @service api!: ApiService;
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;

  @task
  async initializeOIDCFlow(authMethodProvider: AuthMethod.AsObject): Promise<void> {
    let authMethodProviderName = authMethodProvider.name;
    let urlRequest = new GetOIDCAuthURLRequest();
    let authMethodrequest = new GetAuthMethodRequest();

    let authMethodRef = new Ref.AuthMethod();
    authMethodRef.setName(authMethodProviderName);
    authMethodrequest.setAuthMethod(authMethodRef);
    urlRequest.setAuthMethod(authMethodRef);

    let nonce = this._generateNonce();
    urlRequest.setNonce(nonce);

    this._storeOIDCAuthData(nonce, authMethodProviderName);
    let redirectUri = `${window.location.origin}/auth/oidc-callback`;
    urlRequest.setRedirectUri(redirectUri);
    try {
      let authUrl = await this.api.client.getOIDCAuthURL(urlRequest, this.api.WithMeta());
      await window.location.replace(authUrl.getUrl());
    } catch (err) {
      this.flashMessages.error('Authentication failed', { content: err.message, sticky: true });
    }
  }

  // Generate a 20-char nonce, using window.crypto to
  // create a sufficiently-large output then trimming
  _generateNonce(): string {
    let randomArray = new Uint32Array(10);
    window.crypto.getRandomValues(randomArray);
    return randomArray.join('').slice(0, 20);
  }

  // Store OIDC Data in LocalStorage, this gets cleaned up on authentication
  _storeOIDCAuthData(nonce: string, authMethod: string): void {
    window.localStorage.setItem('waypointOIDCNonce', nonce);
    window.localStorage.setItem('waypointOIDCAuthMethod', authMethod);
  }
}
