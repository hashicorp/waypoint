declare module 'ember-simple-auth/services/session' {
  interface SessionService {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    authenticate(authenticator: string, params: any): Promise<void>;
    isAuthenticated: boolean;
    invalidate(): Promise<void>;
    data: SessionData;
  }

  interface SessionData {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    authenticated?: any;
  }
}
