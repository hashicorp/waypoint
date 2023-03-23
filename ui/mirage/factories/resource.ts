/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory, trait } from 'ember-cli-mirage';
import { fakeId } from '../utils';

export default Factory.extend({
  state: () => ({}),

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
              image: 'marketing-public/wp-matrix@sha256:c47cbb1d0526ad29183fb14919ff6c757ec31173',
            },
            {
              image: 'localhost:5000/wp-matrix:a-very-long-but-still-human-readable-tag',
            },
            {
              image: 'marketing-public/wp-matrix:latest',
            },
            {
              image: 'quay.io/marketing-public/wp-matrix',
            },
          ],
        },
      },
    }),
  }),

  'random-service': trait({
    id: () => fakeId(),
    name: 'web',
    platform: 'kubernetes',
    type: 'service',
    categoryDisplayHint: 'ROUTER',
    createdTime: () => new Date(),
    state: () => ({
      ipAddress: '10.104.177.149',
      service: {
        metadata: {
          name: 'web',
          namespace: 'default',
          uid: '267ef8f3-2b41-4e3e-88b4-b28799dc87a0',
          resourceVersion: '405154',
        },
        spec: {
          ports: [{ protocol: 'TCP', port: 80, targetPort: 'http' }],
          selector: {
            name: 'web-v2',
            'waypoint.hashicorp.com/id': '01FGW6JS8XWG4G66YD8RH75BNM',
          },
          clusterIP: '10.104.177.149',
          clusterIPs: ['10.104.177.149'],
          type: 'ClusterIP',
          sessionAffinity: 'None',
          ipFamilies: ['IPv4'],
          ipFamilyPolicy: 'SingleStack',
        },
      },
    }),
    health: 'READY',
    healthMessage: 'Available: Deployment has minimum availability.',
  }),
});
