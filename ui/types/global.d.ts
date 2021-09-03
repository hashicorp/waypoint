declare module 'ember-cli-mirage/test-support' {
  export function setupMirage(hooks: NestedHooks): void;
}

declare module 'ember-test-helpers' {
  interface TestContext {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    server: any;
  }
}

declare module 'ember-a11y-testing/test-support/audit' {
  export default function a11yAudit(
    target?: string | Element,
    axeOptions?: Record<string, unknown>
  ): Promise<void>;
}

declare module '@ember/routing/router-service' {
  import Route from '@ember/routing/route';

  type Transition = ReturnType<Route['transitionTo']>;

  export default class RouterService {
    // This method comes from ember-router-service-refresh-polyfill,
    // which does not provide its own type declarations.
    refresh(pivotRouteName?: string): Transition;
  }
}
