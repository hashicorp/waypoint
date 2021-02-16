import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import ApiService from 'waypoint/services/api';
import { Project, UpsertProjectRequest } from 'waypoint-pb';


export default class WorkspaceProjectsNew extends Controller {
  @service api!: ApiService;

  @action
  async saveProject() {
    let project = this.model;
    let ref = new Project();
    ref.setName(project.name);
    let req = new UpsertProjectRequest();
    req.setProject(ref);
    await this.api.client.upsertProject(req, this.api.WithMeta());
    // route to new project settings
  }

}
