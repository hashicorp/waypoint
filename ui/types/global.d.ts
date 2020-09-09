declare module 'ember-cli-mirage/test-support' {
  export function setupMirage(hooks: NestedHooks): void;
}

declare module 'ember-a11y-testing/test-support/audit' {
  export default function a11yAudit(target?: string | Element, axeOptions?: object): Promise<void>;
}
