/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import ApiService from 'waypoint/services/api';
import { Project, UpsertProjectRequest } from 'waypoint-pb';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';
export default class WorkspaceProjectsCreate extends Controller {
  @service api!: ApiService;
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;

  @tracked createGit = false;

  @action
  async saveProject(e: Event): Promise<void> {
    e.preventDefault();
    let project = this.model;
    let ref = new Project();
    ref.setName(project.name);
    let req = new UpsertProjectRequest();
    req.setProject(ref);
    try {
      let newProject = await this.api.client.upsertProject(req, this.api.WithMeta());
      this.flashMessages.success(`Project "${project.name}" created`);
      if (this.createGit) {
        this.transitionToRoute('workspace.projects.project.settings', newProject.toObject().project?.name);
      } else {
        this.transitionToRoute('workspace.projects.project', newProject.toObject().project?.name);
      }
    } catch (err) {
      this.flashMessages.error('Failed to create project', { content: err.message, sticky: true });
    }
  }
}
