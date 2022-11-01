import { Server } from 'ember-cli-mirage';

export default function (server: Server): void {
  server.create('project', 'marketing-public');
  server.create('project', 'mutable-deployments');
  server.create('project', 'example-nodejs');
  server.create('project', {
    name: 'init-test',
    dataSource: server.create('job-data-source', 'marketing-public'),
  });
}
