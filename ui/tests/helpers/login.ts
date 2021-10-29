import { authenticateSession } from 'ember-simple-auth/test-support';

export default async function login(waypointToken = 'default-test-token-value'): Promise<void> {
  return await authenticateSession({ token: waypointToken });
}
