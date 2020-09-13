import { helper } from '@ember/component/helper';

// iconForComponent
export function iconForComponent([component]: [string]): string {
  switch (component) {
    case 'kubernetes':
      return 'logo-kubernetes-color';
    case 'docker':
      return 'logo-docker-color';
    case 'google-cloud-run':
      return 'logo-gcp-color';
    default:
      return 'help-circle-outline';
  }
}

export default helper(iconForComponent);
