import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class WorkspaceController extends Controller {
  queryParams = ['cli'];

  @tracked cli = null;
}
