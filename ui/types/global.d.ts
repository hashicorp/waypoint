/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

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
