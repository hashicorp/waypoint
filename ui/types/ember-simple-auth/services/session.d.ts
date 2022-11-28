import type RouterService from '@ember/routing/router-service';
type Transition = ReturnType<RouterService['transitionTo']>;

declare module 'ember-simple-auth/services/session' {
  type SessionEvent = 'authenticationSucceeded' | 'invalidationSucceeded';

  export default class SessionService {
    authenticate(authenticator: string, params: unknown): Promise<void>;
    isAuthenticated: boolean;
    invalidate(): Promise<void>;
    attemptedTransition: Transition;
    data: SessionData;

    on(event: SessionEvent, callback: () => void): void;
    set(key: string, value: unknown): void;
    setup(): Promise<void>;
    handleInvalidation(routeAfterInvalidation: string): void;
  }

  interface SessionData {
    authenticated?: Record<string, unknown>;
    workspace?: string;
  }
}
