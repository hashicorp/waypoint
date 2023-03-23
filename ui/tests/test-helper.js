/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Application from 'waypoint/app';
import config from 'waypoint/config/environment';
import * as QUnit from 'qunit';
import { setApplication } from '@ember/test-helpers';
import { setup } from 'qunit-dom';
import { start } from 'ember-qunit';
import { setup as setupA11yTesting } from './helpers/a11y';
import './helpers/flash-message';
import './helpers/xterm';

setApplication(Application.create(config.APP));

setupA11yTesting();

setup(QUnit.assert);

start();
