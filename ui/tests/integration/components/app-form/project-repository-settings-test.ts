/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, clearRender, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { Project } from 'waypoint-pb';

module('Integration | Component | app-form/project-repository-settings', function (hooks) {
  setupRenderingTest(hooks);

  test('second new project does not have previous input', async function (assert) {
    this.set('project', {});
    await render(hbs`<AppForm::ProjectRepositorySettings @project={{this.project}} />`);
    await fillIn('#git-source-url', 'https://github.projectone.git');
    await fillIn('#git-source-username', 'admin');
    await fillIn('#git-source-password', 'password');
    await clearRender();

    this.set('project2', new Project().toObject());
    await render(hbs`<AppForm::ProjectRepositorySettings @project={{this.project2}} />`);

    assert.dom('#git-source-url').hasValue('');
    assert.dom('#git-source-username').hasValue('');
    assert.dom('#git-source-password').hasValue('');
  });

  test('populated applications list does not break render', async function (assert) {
    this.set('project', {
      applicationsList: [
        {
          name: 'app-1',
          project: {
            project: 'project',
          },
        },
      ],
      dataSource: {
        git: {
          url: 'https://github.com',
        },
      },
    });
    await render(hbs`<AppForm::ProjectRepositorySettings @project={{this.project}} />`);

    assert.dom('#git-source-url').hasValue('https://github.com');
    assert.dom('#git-auth-not-set').isChecked();
  });

  test('cli generated project does not break render', async function (assert) {
    this.set('project', {
      applicationsList: [
        {
          name: 'app-1',
          project: {
            project: 'project',
          },
        },
      ],
      dataSource: undefined,
      dataSourcePoll: undefined,
    });
    await render(hbs`<AppForm::ProjectRepositorySettings @project={{this.project}} />`);

    assert.dom('#git-source-url').hasValue('');
    assert.dom('#git-auth-basic').isChecked();
    assert.dom('#git-source-username').hasValue('');
    assert.dom('#git-source-password').hasValue('');
  });

  test('basic auth project loads properly', async function (assert) {
    this.set('project', {
      dataSource: {
        git: {
          url: 'https://github.com',
          basic: {
            username: 'user',
            password: 'password',
          },
        },
      },
    });
    await render(hbs`<AppForm::ProjectRepositorySettings @project={{this.project}} />`);

    assert.dom('#git-auth-basic').isChecked();
    assert.dom('#git-source-username').hasValue('user');
    assert.dom('#git-source-password').hasValue('password');
  });

  test('ssh project loads properly', async function (assert) {
    this.set('project', {
      dataSource: {
        git: {
          url: 'https://github.com',
          ssh: {
            user: 'user',
            password: 'password',
            privateKeyPem: 'private key',
          },
        },
      },
    });
    await render(hbs`<AppForm::ProjectRepositorySettings @project={{this.project}} />`);

    assert.dom('#git-auth-ssh').isChecked();
    assert.dom('#git-source-ssh-user').hasValue('user');
    assert.dom('#git-source-ssh-password').hasValue('password');
    assert.dom('#git-source-ssh-key').hasValue(atob('private key'));
  });

  test('no auth loads properly', async function (assert) {
    this.set('project', {
      dataSource: {
        git: {
          url: 'https://github.com',
        },
      },
    });
    await render(hbs`<AppForm::ProjectRepositorySettings @project={{this.project}} />`);

    assert.dom('#git-auth-not-set').isChecked();
  });
});
