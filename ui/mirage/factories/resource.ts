import { Factory, trait } from 'ember-cli-mirage';
import { fakeId } from '../utils';

export default Factory.extend({
  'random-deployment': trait({
    id: () => fakeId(),
    name: 'web-01ffr30gszyz43x7jxkkyt7zjk',
    platform: 'kubernetes',
    type: 'deployment',
    categoryDisplayHint: 'INSTANCE_MANAGER',
    createdTime: () => new Date(),
    state: () => ({
      deployment: {
        metadata: {
          name: 'wp-matrix-01ffr30gszyz43x7jxkkyt7zjk',
          namespace: 'default',
          uid: 'afe78df9-7032-425b-a330-4538cdfee81d',
          resourceVersion: '26926',
        },
      },
    }),
    health: 'READY',
    healthMessage: 'Available: Deployment has minimum availability.',

    afterCreate(resource, server) {
      server.create('resource', 'random-pod', { statusReport: resource.statusReport, parent: resource });
    },
  }),

  'random-pod': trait({
    id: () => fakeId(),
    name: 'web-01ffr30gszyz43x7jxkkyt7zjk-7c8cf49b76-f75s7',
    createdTime: () => new Date(),
    platform: 'kubernetes',
    type: 'pod',
    categoryDisplayHint: 'INSTANCE',
    health: 'READY',
    healthMessage: 'ready',
    state: () => ({
      hostIP: '192.168.65.4',
      ipAddress: '10.1.0.38',
      pod: {
        metadata: {
          name: 'wp-matrix-01ffr30gszyz43x7jxkkyt7zjk-7c8cf49b76-f75s7',
          generateName: 'wp-matrix-01ffr30gszyz43x7jxkkyt7zjk-7c8cf49b76-',
          namespace: 'default',
          labels: {
            app: 'wp-matrix-v1',
            name: 'wp-matrix-v1',
            'pod-template-hash': 'cf9f9fb8f',
            version: '01FGBX52AM9FTDR5ZTSSBCHF6Q',
            'waypoint.hashicorp.com/id': '01FGBX52AM9FTDR5ZTSSBCHF6Q',
          },
        },
        spec: {
          containers: [
            {
              image: 'marketing-public/wp-matrix:1',
            },
          ],
        },
      },
    }),
  }),
});
