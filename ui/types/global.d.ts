// Types for compiled templates
declare module 'waypoint/templates/*' {
  import { TemplateFactory } from 'htmlbars-inline-precompile';
  const tmpl: TemplateFactory;
  export default tmpl;
}

declare module 'ember-cli-mirage/test-support' {
  export function setupMirage(hooks: NestedHooks): void;
}