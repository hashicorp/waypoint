import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class AuthController extends Controller {
  queryParams = ['token', 'cli'];

  @tracked token = null;
  @tracked cli = null;
}
