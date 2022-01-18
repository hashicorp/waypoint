import ApplicationInstance from '@ember/application/instance';

export function initialize(appInstance: ApplicationInstance): void {
  let session = appInstance.lookup('service:session');

  session.on('invalidationSucceeded', () => {
    // Clear workspace on logout
    session.set('data.workspace', undefined);
  });
}

export default {
  name: 'session-workspace-storage',
  initialize,
};
