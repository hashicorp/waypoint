import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Application, ExecStreamRequest, ExecStreamResponse, FindExecInstanceRequest } from 'waypoint-pb';
import { setupConsoleLogger } from 'ember-a11y-testing/test-support';

export default class Exec extends Route {
  @service api!: ApiService;

  async model() {
    let application = this.modelFor('workspace.projects.project.app') as Application.AsObject;
    let latestDeploymentId = application.releases[0]?.deploymentId;


    let execInstanceRequest = new FindExecInstanceRequest();
    execInstanceRequest.setDeploymentId(latestDeploymentId);
    let instanceResponse = await this.api.client.findExecInstance(execInstanceRequest, this.api.WithMeta());
    let instance = instanceResponse.toObject()?.instance;
    let req = new ExecStreamRequest();

    let start = await new ExecStreamRequest.Start();
    start.setDeploymentId(latestDeploymentId);
    start.setArgsList(['']);
    start.setInstanceId(instance?.id);
    req.setStart(start);

    let stream = await this.api.client.startExecStream(req, this.api.WithMeta());


    stream.on('data', function (response) {
      console.log(response.toObject());
    })
    stream.on('open', function (response) {
      console.log(response)
    })
    stream.on('status', function (response) {
      console.log(response)
    })

    setTimeout( ()=> {
      req.setInput(new ExecStreamRequest.Input());
    }, 3000)
    // todo(pearkes): construct GetExecStreamRequest
  }
}
