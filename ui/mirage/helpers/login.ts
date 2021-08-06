export default function login(token?: string): void {
  window.localStorage.waypointAuthToken = token || 'default-test-token-value';
}
