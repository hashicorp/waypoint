import { helper } from '@ember/component/helper';

// iconForComponent
export function iconForComponent([component]: [string]): string {
  switch (component) {
    case 'aws-ec2':
      return 'logo-aws-color';
    case 'aws-ecs':
      return 'logo-aws-color';
    case 'azure-container-instances':
      return 'logo-azure-color';
    case 'docker':
      return 'logo-docker-color';
    case 'google-cloud-run':
      return 'logo-gcp-color';
    case 'kubernetes':
      return 'logo-kubernetes-color-alt';
    case 'nomad':
    case 'nomad-jobspec':
      return 'logo-nomad-color';
    case 'pack':
      return 'logo-pack-color';
    default:
      return 'more-horizontal';
  }
}

export default helper(iconForComponent);
