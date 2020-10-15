export default function login(token?: string) {
  window.localStorage.waypointAuthToken = token || 'default-test-token-value';
}
