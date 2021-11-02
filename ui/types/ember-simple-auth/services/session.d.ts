import type RouterService from '@ember/routing/router-service';
type Transition = ReturnType<RouterService['transitionTo']>;

declare module 'ember-simple-auth/services/session' {
  interface SessionService {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    authenticate(authenticator: string, params: any): Promise<void>;
    isAuthenticated: boolean;
    invalidate(): Promise<void>;
    attemptedTransition: Transition;
    data: SessionData;
  }

  interface SessionData {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    authenticated?: any;
  }
}
