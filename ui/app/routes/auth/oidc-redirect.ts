import ApiService from 'waypoint/services/api';
import OIDCAuthenticator from 'waypoint/authenticators/oidc';
import Route from '@ember/routing/route';
import SessionService from 'ember-simple-auth/services/session';
import { inject as service } from '@ember/service';

export default class AuthIndex extends Route {
  @service session!: SessionService;
  @service api!: ApiService;

  beforeModel() {
    let oidcParams = OIDCAuthenticator.parseResponse(window.location.hash);
    console.log(oidcParams);
  }
}
