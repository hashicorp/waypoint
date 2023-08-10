/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import RouterService from '@ember/routing/router-service';
import { StatusReport } from 'waypoint-pb';
import { findImageRefs } from 'waypoint/utils/image-refs';

interface Args {
  resource?: StatusReport.Resource.AsObject;
}

type LabelMap = Record<string, unknown>;

interface MaybePodState {
  pod?: {
    metadata?: {
      labels?: LabelMap;
    };
    spec?: {
      containers?: { image?: string }[];
    };
  };
}

export default class extends Component<Args> {
  @service router!: RouterService;

  get state(): MaybePodState | undefined {
    try {
      return JSON.parse(this.args.resource?.stateJson ?? '{}');
    } catch (error) {
      console.error(error);
      return;
    }
  }

  get labels(): LabelMap | undefined {
    return this.state?.pod?.metadata?.labels;
  }

  get imageRefs(): ReturnType<typeof findImageRefs> {
    return findImageRefs(this.state);
  }

  get hasLabels(): boolean {
    return !!this.labels && Object.keys(this.labels).length !== 0;
  }
}
