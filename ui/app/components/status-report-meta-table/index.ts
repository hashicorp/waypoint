/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';
import type ProjectService from 'waypoint/services/project';
import {
  Ref,
  ExpediteStatusReportRequest,
  GetJobStreamRequest,
  GetJobStreamResponse,
  Job,
} from 'waypoint-pb';
import { DeploymentExtended } from 'waypoint/services/api';

interface Args {
  model: DeploymentExtended;
  artifactType: string;
}

export default class StatusReportMetaTable extends Component<Args> {
  @service api!: ApiService;
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;
  @service declare project: ProjectService;
  @tracked isRefreshRunning = false;

  get artifactType(): Args['artifactType'] {
    return this.args.artifactType;
  }

  get model(): Args['model'] {
    return this.args.model;
  }

  get statusReport(): Args['model']['statusReport'] {
    return this.model.statusReport;
  }

  get projectHasDataSource(): boolean {
    let dataSource = this.project.current?.dataSource;
    return Boolean(dataSource && !dataSource.local);
  }

  @action
  async refreshHealthCheck(e: Event): Promise<void> {
    e.preventDefault();

    let ref = new Ref.Operation();
    ref.setId(this.args.model.id);

    let workspace = new Ref.Workspace();
    let wkspName = this.args.model.workspace?.workspace || 'default';
    workspace.setWorkspace(wkspName);

    let req = new ExpediteStatusReportRequest();
    req.setWorkspace(workspace);

    if (this.artifactType === 'Deployment') {
      req.setDeployment(ref);
    } else if (this.artifactType === 'Release') {
      req.setRelease(ref);
    }

    let resp = await this.api.client.expediteStatusReport(req, this.api.WithMeta()).catch((error) => {
      this.flashMessages.error(error.message);
    });

    if (resp && resp?.getJobId()) {
      this.isRefreshRunning = true;

      let streamReq = new GetJobStreamRequest();
      streamReq.setJobId(resp.getJobId());
      let jobStream = await this.api.client.getJobStream(streamReq, this.api.WithMeta());

      // handler for job stream when receiving data
      let onData = async (response: GetJobStreamResponse) => {
        let event = response.getEventCase();
        if (event === GetJobStreamResponse.EventCase.STATE) {
          let state = response.getState()?.getCurrent() as Job.State;
          if (state === 5) {
            this.isRefreshRunning = false;
          }
        }
      };

      jobStream.on('data', onData);
    }
  }
}
