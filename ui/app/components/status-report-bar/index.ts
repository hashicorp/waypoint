import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import FlashMessagesService from 'waypoint/services/flash-messages';
import {
  Deployment,
  Ref,
  ExpediteStatusReportRequest,
  StatusReport,
  GetJobStreamRequest,
  GetJobStreamResponse,
  Job,
} from 'waypoint-pb';

interface StatusReportBarArgs {
  model: Deployment.AsObject & WithStatusReport;
  artifactType: string;
}

interface WithStatusReport {
  statusReport?: StatusReport.AsObject;
}

export default class StatusReportBar extends Component<StatusReportBarArgs> {
  @service api!: ApiService;
  @service flashMessages!: FlashMessagesService;
  @tracked isRefreshRunning = false;

  constructor(owner: any, args: any) {
    super(owner, args);
  }

  get artifactType() {
    return this.args.artifactType;
  }

  get model() {
    return this.args.model;
  }

  get statusReport() {
    return this.model.statusReport;
  }

  @action
  async refreshHealthCheck(e: Event) {
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

    if (resp?.getJobId()) {
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
