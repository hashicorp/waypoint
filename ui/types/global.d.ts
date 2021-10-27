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

declare module 'miragejs' {
  /**
   * RouteHandler is the context (the `this`) in which our Mirage route handlers
   * are executed. Mirage itself does not export this type declaration in any
   * form so we’re merging it into the module declaration.
   */
  export interface RouteHandler {
    serialize(response: unknown, serializerType: string): Response;
  }
}
