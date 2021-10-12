import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class AuthController extends Controller {
  queryParams = ['state', 'code', 'scope', 'authuser', 'prompt'];

  @tracked state = null;
  @tracked code = null;
  @tracked scope = null;
}
