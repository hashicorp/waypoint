import { VersionInfo, GetVersionInfoResponse } from 'waypoint-pb';

function createVersionInfo(): VersionInfo {
  let versionInfo = new VersionInfo();
  let protocolVersion = new VersionInfo.ProtocolVersion();
  protocolVersion.setCurrent(1);
  versionInfo.setApi(protocolVersion);
  versionInfo.setEntrypoint(protocolVersion);
  versionInfo.setVersion('0.3.12');
  return versionInfo;
}

export function get(schema: any, { params, requestHeaders }) {
  let resp = new GetVersionInfoResponse();
  let versionInfo = createVersionInfo();
  resp.setInfo(versionInfo);
  return this.serialize(resp, 'application');
}
