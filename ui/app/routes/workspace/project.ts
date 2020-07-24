import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import ApiService from 'waypoint/services/api';
import { Ref } from 'waypoint-pb';

interface ProjectModelParams {
  project_id: string;
}

export default class Project extends Route {
  @service api!: ApiService;

  async model(params: ProjectModelParams) {
    let proj = new Ref.Project();

    // Project based on id
    proj.setProject(params.project_id);

    // todo(pearkes): actually list applications once that api is ready
    // let resp = await this.api.client.listApplications(new Empty(), {})
    // let apps = resp.getProjectsList().map(p => p.toObject());
    var app = new Ref.Application();
    app.setProject(proj.getProject());
    app.setApplication('wp-gcr-deno-test');
    let apps = [app.toObject()];

    return {
      ref: proj as Ref.Project,
      project: proj.toObject(),
      applications: apps,
    };
  }
}
