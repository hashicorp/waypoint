import { VersionInfo, GetVersionInfoResponse } from 'waypoint-pb';
import { Response, RouteHandler } from 'miragejs';

function createVersionInfo(): VersionInfo {
  let versionInfo = new VersionInfo();
  let protocolVersion = new VersionInfo.ProtocolVersion();
  protocolVersion.setCurrent(1);
  versionInfo.setApi(protocolVersion);
  versionInfo.setEntrypoint(protocolVersion);
  versionInfo.setVersion('0.4.2');
  return versionInfo;
}

export function get(this: RouteHandler): Response {
  let resp = new GetVersionInfoResponse();
  let versionInfo = createVersionInfo();
  resp.setInfo(versionInfo);
  return this.serialize(resp, 'application');
}
