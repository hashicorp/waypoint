// Types for compiled templates
declare module '@hashicorp/waypoint-ui-toolkit/templates/*' {
  import { TemplateFactory } from 'htmlbars-inline-precompile';
  const tmpl: TemplateFactory;
  export default tmpl;
}
