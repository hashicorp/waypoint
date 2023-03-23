/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { VersionInfo, GetVersionInfoResponse } from 'waypoint-pb';
import { Response } from 'miragejs';
import { RouteHandler } from '../types';

function createVersionInfo(): VersionInfo {
  let versionInfo = new VersionInfo();
  let protocolVersion = new VersionInfo.ProtocolVersion();
  protocolVersion.setCurrent(1);
  versionInfo.setApi(protocolVersion);
  versionInfo.setEntrypoint(protocolVersion);
  versionInfo.setVersion('[Mirage]');
  return versionInfo;
}

export function get(this: RouteHandler): Response {
  let resp = new GetVersionInfoResponse();
  let versionInfo = createVersionInfo();
  resp.setInfo(versionInfo);
  return this.serialize(resp, 'application');
}
