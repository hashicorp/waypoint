import type RouterService from '@ember/routing/router-service';
type Transition = ReturnType<RouterService['transitionTo']>;

declare module 'ember-simple-auth/services/session' {
  export default interface SessionService {
    authenticate(authenticator: string, params: unknown): Promise<void>;
    isAuthenticated: boolean;
    invalidate(): Promise<void>;
    attemptedTransition: Transition;
    data: SessionData;
  }

  interface SessionData {
    authenticated?: Record<string, unknown>;
  }
}
