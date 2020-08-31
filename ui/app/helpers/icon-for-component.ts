import { helper } from '@ember/component/helper';

// iconForComponent
export function iconForComponent([component]: [string]): string {
  switch (component) {
    case 'kubernetes':
      return 'logo-kubernetes-color';
    case 'docker':
      return 'box-outline';
    case 'google-cloud-run':
      return 'logo-gcp-color';
    default:
      return 'source-file';
  }
}

export default helper(iconForComponent);
