import Application from '@ember/application';

export function initialize(application: Application): void {
  let waypointToken = window.localStorage.waypointAuthToken;
  if (waypointToken) {
    window.localStorage.setItem(
      'ember_simple_auth-session',
      JSON.stringify({ authenticated: { authenticator: 'authenticator:token', token: waypointToken } })
    );
    window.localStorage.removeItem('waypointAuthToken');
  }
}

export default {
  initialize
};
