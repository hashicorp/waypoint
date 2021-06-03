import { Server } from 'ember-cli-mirage';

export default function (server: Server): void {
  server.create('project', 'marketing-public');
  server.create('project', 'mutable-deployments');
}
