import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import ApiService from 'waypoint/services/api';
import { UpsertProjectRequest, UpsertProjectResponse, Project, Hcl } from 'waypoint-pb';
import { WaypointClient } from 'waypoint-client';

module('Unit | Service | api', function (hooks) {
  setupTest(hooks);

  module('upsertProject', function () {
    test('sets remote_enabled to true when a data source is present', async function (assert) {
      let api: ApiService = this.owner.lookup('service:api');
      let result = setupMockUpsertProject(api.client);
      let project = projectObject({
        dataSource: {
          git: {
            url: 'https://github.com/hashicorp/waypoint-examples',
            ref: 'head',
            path: '',
            ignoreChangesOutsidePath: true,
            recurseSubmodules: 1,
          },
        },
        remoteEnabled: false,
      });

      api.upsertProject(project);

      assert.true(result.request?.getProject()?.getRemoteEnabled());
    });

    test('leaves remote_enabled alone when data source is absent', async function (assert) {
      let api: ApiService = this.owner.lookup('service:api');
      let result = setupMockUpsertProject(api.client);
      let project = projectObject({
        dataSource: undefined,
        remoteEnabled: true,
      });

      api.upsertProject(project);

      assert.true(result.request?.getProject()?.getRemoteEnabled());
    });
  });
});

/**
 * Mocks out the `upsertProject` method on a WaypointClient instance.
 *
 * Yes we could use a library like Sinon, but things get complicated making all
 * that work with TypeScript.
 *
 * @param client {WaypointClient}
 * @returns {MockResult} object that records interactions with upsertProject
 */
function setupMockUpsertProject(client: WaypointClient): MockResult {
  let result: MockResult = {};
  let upsertProject = (request: UpsertProjectRequest, _meta: never) => {
    result.request = request;
    return Promise.resolve(new UpsertProjectResponse());
  };

  client.upsertProject = upsertProject as WaypointClient['upsertProject'];

  return result;
}

interface MockResult {
  request?: UpsertProjectRequest;
}

/**
 * Provides defaults for a Project.AsObject.
 * @param attrs {Partial<Project.AsObject>} attributes to override
 * @returns {Project.AsObject} a complete Project.AsObject
 */
function projectObject(attrs: Partial<Project.AsObject>): Project.AsObject {
  return {
    applicationsList: [],
    dataSource: undefined,
    fileChangeSignal: 'HUP',
    name: 'test-project',
    remoteEnabled: false,
    variablesList: [],
    waypointHcl: '',
    waypointHclFormat: Hcl.Format.HCL,
    ...attrs,
  };
}
