import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class WorkspaceProjectsProjectAppExec extends Controller {
  queryParams = ['hasExec'];

  @tracked hasExec = null;
}
