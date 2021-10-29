// This initializer is used to migrate the session token in localStorage from the old format to the new one used by ember-simple-auth.
// It should be safe to remove after the feature has been released.

export function initialize(/*application: Application*/): void {
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
  initialize,
};
