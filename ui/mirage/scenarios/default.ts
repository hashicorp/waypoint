import { Server } from 'ember-cli-mirage';
import login from '../helpers/login';

export default function (server: Server): void {
  server.create('project', 'marketing-public');
  server.create('project', 'mutable-deployments');
  server.create('project', 'example-nodejs');
  login();
}
