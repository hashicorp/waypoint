/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper } from '@ember/component/helper';

// iconForComponent
export function iconForComponent([component]: [string]): string {
  switch (component) {
    case 'aws-ec2':
    case 'aws-ecs':
    case 'aws-ecr':
      return 'aws-color';
    case 'azure-container-instances':
      return 'azure-color';
    case 'docker':
      return 'docker-color';
    case 'google-cloud-run':
      return 'gcp-color';
    case 'kubernetes':
    case 'kubernetes-apply':
      return 'kubernetes-color';
    case 'nomad':
    case 'nomad-jobspec':
      return 'nomad-color';
    case 'nomad-jobspec-canary':
      return 'nomad-color';
    case 'pack':
      return 'pack-color';
    default:
      return 'more-horizontal';
  }
}

export default helper(iconForComponent);
