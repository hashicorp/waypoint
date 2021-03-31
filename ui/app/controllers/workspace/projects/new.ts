import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import { Project, UpsertProjectRequest } from 'waypoint-pb';


export default class WorkspaceProjectsNew extends Controller {
  @service api!: ApiService;
  @tracked createGit = false;

  @action
  async saveProject() {
    let project = this.model;
    let ref = new Project();
    ref.setName(project.name);
    let req = new UpsertProjectRequest();
    req.setProject(ref);
    let newProject = await this.api.client.upsertProject(req, this.api.WithMeta());
    if (this.createGit) {
      this.transitionToRoute('workspace.projects.project.settings', newProject.toObject().project?.name);
    } else {
      this.transitionToRoute('workspace.projects.project', newProject.toObject().project?.name);
    }
  }

}
