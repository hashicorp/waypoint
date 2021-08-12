import { Server } from 'ember-cli-mirage';
import login from '../helpers/login';

export default function (server: Server): void {
  let project1 = server.create('project', { name: 'microchip' });
  let project2 = server.create('project', 'mutable-deployments', { name: 'microchip-mutable' });
  // maybe 2 applications
  // 1 with mutable deployments ?
  let application1 = server.create('application', { name: 'wp-bandwidth', project: project1 });
  let application2 = server.create('application', { name: 'wp-nginx', project: project1 });

  // builds
  server.createList('deployment', 3, 'docker', '5-minutes-old-success', { application: application1 });
  // releases (2 done, 1 ongoing) (or 1 success, 1 failed, 1 ongoing)
  login();
}
