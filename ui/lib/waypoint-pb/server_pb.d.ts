import * as jspb from 'google-protobuf'

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_rpc_status_pb from 'api-common-protos/google/rpc/status_pb';
import * as any_pb from 'opaqueany/any_pb';


export class UI extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UI.AsObject;
  static toObject(includeInstance: boolean, msg: UI): UI.AsObject;
  static serializeBinaryToWriter(message: UI, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UI;
  static deserializeBinaryFromReader(message: UI, reader: jspb.BinaryReader): UI;
}

export namespace UI {
  export type AsObject = {
  }

  export class GetProjectRequest extends jspb.Message {
    getProject(): Ref.Project | undefined;
    setProject(value?: Ref.Project): GetProjectRequest;
    hasProject(): boolean;
    clearProject(): GetProjectRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetProjectRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetProjectRequest): GetProjectRequest.AsObject;
    static serializeBinaryToWriter(message: GetProjectRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetProjectRequest;
    static deserializeBinaryFromReader(message: GetProjectRequest, reader: jspb.BinaryReader): GetProjectRequest;
  }

  export namespace GetProjectRequest {
    export type AsObject = {
      project?: Ref.Project.AsObject,
    }
  }


  export class GetProjectResponse extends jspb.Message {
    getProject(): Project | undefined;
    setProject(value?: Project): GetProjectResponse;
    hasProject(): boolean;
    clearProject(): GetProjectResponse;

    getLatestInitJob(): Job | undefined;
    setLatestInitJob(value?: Job): GetProjectResponse;
    hasLatestInitJob(): boolean;
    clearLatestInitJob(): GetProjectResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetProjectResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetProjectResponse): GetProjectResponse.AsObject;
    static serializeBinaryToWriter(message: GetProjectResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetProjectResponse;
    static deserializeBinaryFromReader(message: GetProjectResponse, reader: jspb.BinaryReader): GetProjectResponse;
  }

  export namespace GetProjectResponse {
    export type AsObject = {
      project?: Project.AsObject,
      latestInitJob?: Job.AsObject,
    }
  }


  export class ListDeploymentsRequest extends jspb.Message {
    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): ListDeploymentsRequest;
    hasApplication(): boolean;
    clearApplication(): ListDeploymentsRequest;

    getWorkspace(): Ref.Workspace | undefined;
    setWorkspace(value?: Ref.Workspace): ListDeploymentsRequest;
    hasWorkspace(): boolean;
    clearWorkspace(): ListDeploymentsRequest;

    getOrder(): OperationOrder | undefined;
    setOrder(value?: OperationOrder): ListDeploymentsRequest;
    hasOrder(): boolean;
    clearOrder(): ListDeploymentsRequest;

    getStatusList(): Array<StatusFilter>;
    setStatusList(value: Array<StatusFilter>): ListDeploymentsRequest;
    clearStatusList(): ListDeploymentsRequest;
    addStatus(value?: StatusFilter, index?: number): StatusFilter;

    getPhysicalState(): Operation.PhysicalState;
    setPhysicalState(value: Operation.PhysicalState): ListDeploymentsRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListDeploymentsRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ListDeploymentsRequest): ListDeploymentsRequest.AsObject;
    static serializeBinaryToWriter(message: ListDeploymentsRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListDeploymentsRequest;
    static deserializeBinaryFromReader(message: ListDeploymentsRequest, reader: jspb.BinaryReader): ListDeploymentsRequest;
  }

  export namespace ListDeploymentsRequest {
    export type AsObject = {
      application?: Ref.Application.AsObject,
      workspace?: Ref.Workspace.AsObject,
      order?: OperationOrder.AsObject,
      statusList: Array<StatusFilter.AsObject>,
      physicalState: Operation.PhysicalState,
    }
  }


  export class ListDeploymentsResponse extends jspb.Message {
    getDeploymentsList(): Array<UI.DeploymentBundle>;
    setDeploymentsList(value: Array<UI.DeploymentBundle>): ListDeploymentsResponse;
    clearDeploymentsList(): ListDeploymentsResponse;
    addDeployments(value?: UI.DeploymentBundle, index?: number): UI.DeploymentBundle;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListDeploymentsResponse.AsObject;
    static toObject(includeInstance: boolean, msg: ListDeploymentsResponse): ListDeploymentsResponse.AsObject;
    static serializeBinaryToWriter(message: ListDeploymentsResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListDeploymentsResponse;
    static deserializeBinaryFromReader(message: ListDeploymentsResponse, reader: jspb.BinaryReader): ListDeploymentsResponse;
  }

  export namespace ListDeploymentsResponse {
    export type AsObject = {
      deploymentsList: Array<UI.DeploymentBundle.AsObject>,
    }
  }


  export class DeploymentBundle extends jspb.Message {
    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): DeploymentBundle;
    hasDeployment(): boolean;
    clearDeployment(): DeploymentBundle;

    getArtifact(): PushedArtifact | undefined;
    setArtifact(value?: PushedArtifact): DeploymentBundle;
    hasArtifact(): boolean;
    clearArtifact(): DeploymentBundle;

    getBuild(): Build | undefined;
    setBuild(value?: Build): DeploymentBundle;
    hasBuild(): boolean;
    clearBuild(): DeploymentBundle;

    getDeployUrl(): string;
    setDeployUrl(value: string): DeploymentBundle;

    getJobDataSourceRef(): Job.DataSource.Ref | undefined;
    setJobDataSourceRef(value?: Job.DataSource.Ref): DeploymentBundle;
    hasJobDataSourceRef(): boolean;
    clearJobDataSourceRef(): DeploymentBundle;

    getLatestStatusReport(): StatusReport | undefined;
    setLatestStatusReport(value?: StatusReport): DeploymentBundle;
    hasLatestStatusReport(): boolean;
    clearLatestStatusReport(): DeploymentBundle;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeploymentBundle.AsObject;
    static toObject(includeInstance: boolean, msg: DeploymentBundle): DeploymentBundle.AsObject;
    static serializeBinaryToWriter(message: DeploymentBundle, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeploymentBundle;
    static deserializeBinaryFromReader(message: DeploymentBundle, reader: jspb.BinaryReader): DeploymentBundle;
  }

  export namespace DeploymentBundle {
    export type AsObject = {
      deployment?: Deployment.AsObject,
      artifact?: PushedArtifact.AsObject,
      build?: Build.AsObject,
      deployUrl: string,
      jobDataSourceRef?: Job.DataSource.Ref.AsObject,
      latestStatusReport?: StatusReport.AsObject,
    }
  }


  export class ListReleasesRequest extends jspb.Message {
    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): ListReleasesRequest;
    hasApplication(): boolean;
    clearApplication(): ListReleasesRequest;

    getWorkspace(): Ref.Workspace | undefined;
    setWorkspace(value?: Ref.Workspace): ListReleasesRequest;
    hasWorkspace(): boolean;
    clearWorkspace(): ListReleasesRequest;

    getOrder(): OperationOrder | undefined;
    setOrder(value?: OperationOrder): ListReleasesRequest;
    hasOrder(): boolean;
    clearOrder(): ListReleasesRequest;

    getStatusList(): Array<StatusFilter>;
    setStatusList(value: Array<StatusFilter>): ListReleasesRequest;
    clearStatusList(): ListReleasesRequest;
    addStatus(value?: StatusFilter, index?: number): StatusFilter;

    getPhysicalState(): Operation.PhysicalState;
    setPhysicalState(value: Operation.PhysicalState): ListReleasesRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListReleasesRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ListReleasesRequest): ListReleasesRequest.AsObject;
    static serializeBinaryToWriter(message: ListReleasesRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListReleasesRequest;
    static deserializeBinaryFromReader(message: ListReleasesRequest, reader: jspb.BinaryReader): ListReleasesRequest;
  }

  export namespace ListReleasesRequest {
    export type AsObject = {
      application?: Ref.Application.AsObject,
      workspace?: Ref.Workspace.AsObject,
      order?: OperationOrder.AsObject,
      statusList: Array<StatusFilter.AsObject>,
      physicalState: Operation.PhysicalState,
    }
  }


  export class ListReleasesResponse extends jspb.Message {
    getReleasesList(): Array<UI.ReleaseBundle>;
    setReleasesList(value: Array<UI.ReleaseBundle>): ListReleasesResponse;
    clearReleasesList(): ListReleasesResponse;
    addReleases(value?: UI.ReleaseBundle, index?: number): UI.ReleaseBundle;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListReleasesResponse.AsObject;
    static toObject(includeInstance: boolean, msg: ListReleasesResponse): ListReleasesResponse.AsObject;
    static serializeBinaryToWriter(message: ListReleasesResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListReleasesResponse;
    static deserializeBinaryFromReader(message: ListReleasesResponse, reader: jspb.BinaryReader): ListReleasesResponse;
  }

  export namespace ListReleasesResponse {
    export type AsObject = {
      releasesList: Array<UI.ReleaseBundle.AsObject>,
    }
  }


  export class ReleaseBundle extends jspb.Message {
    getRelease(): Release | undefined;
    setRelease(value?: Release): ReleaseBundle;
    hasRelease(): boolean;
    clearRelease(): ReleaseBundle;

    getLatestStatusReport(): StatusReport | undefined;
    setLatestStatusReport(value?: StatusReport): ReleaseBundle;
    hasLatestStatusReport(): boolean;
    clearLatestStatusReport(): ReleaseBundle;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ReleaseBundle.AsObject;
    static toObject(includeInstance: boolean, msg: ReleaseBundle): ReleaseBundle.AsObject;
    static serializeBinaryToWriter(message: ReleaseBundle, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ReleaseBundle;
    static deserializeBinaryFromReader(message: ReleaseBundle, reader: jspb.BinaryReader): ReleaseBundle;
  }

  export namespace ReleaseBundle {
    export type AsObject = {
      release?: Release.AsObject,
      latestStatusReport?: StatusReport.AsObject,
    }
  }

}

export class GetVersionInfoResponse extends jspb.Message {
  getInfo(): VersionInfo | undefined;
  setInfo(value?: VersionInfo): GetVersionInfoResponse;
  hasInfo(): boolean;
  clearInfo(): GetVersionInfoResponse;

  getServerFeatures(): ServerFeatures | undefined;
  setServerFeatures(value?: ServerFeatures): GetVersionInfoResponse;
  hasServerFeatures(): boolean;
  clearServerFeatures(): GetVersionInfoResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVersionInfoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetVersionInfoResponse): GetVersionInfoResponse.AsObject;
  static serializeBinaryToWriter(message: GetVersionInfoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVersionInfoResponse;
  static deserializeBinaryFromReader(message: GetVersionInfoResponse, reader: jspb.BinaryReader): GetVersionInfoResponse;
}

export namespace GetVersionInfoResponse {
  export type AsObject = {
    info?: VersionInfo.AsObject,
    serverFeatures?: ServerFeatures.AsObject,
  }
}

export class ServerFeatures extends jspb.Message {
  getFeaturesList(): Array<ServerFeatures.feature>;
  setFeaturesList(value: Array<ServerFeatures.feature>): ServerFeatures;
  clearFeaturesList(): ServerFeatures;
  addFeatures(value: ServerFeatures.feature, index?: number): ServerFeatures;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerFeatures.AsObject;
  static toObject(includeInstance: boolean, msg: ServerFeatures): ServerFeatures.AsObject;
  static serializeBinaryToWriter(message: ServerFeatures, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerFeatures;
  static deserializeBinaryFromReader(message: ServerFeatures, reader: jspb.BinaryReader): ServerFeatures;
}

export namespace ServerFeatures {
  export type AsObject = {
    featuresList: Array<ServerFeatures.feature>,
  }

  export enum feature { 
    FEATURE_UNSPECIFIED = 0,
    FEATURE_INLINE_KEEPALIVES = 1,
  }
}

export class VersionInfo extends jspb.Message {
  getApi(): VersionInfo.ProtocolVersion | undefined;
  setApi(value?: VersionInfo.ProtocolVersion): VersionInfo;
  hasApi(): boolean;
  clearApi(): VersionInfo;

  getEntrypoint(): VersionInfo.ProtocolVersion | undefined;
  setEntrypoint(value?: VersionInfo.ProtocolVersion): VersionInfo;
  hasEntrypoint(): boolean;
  clearEntrypoint(): VersionInfo;

  getVersion(): string;
  setVersion(value: string): VersionInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VersionInfo.AsObject;
  static toObject(includeInstance: boolean, msg: VersionInfo): VersionInfo.AsObject;
  static serializeBinaryToWriter(message: VersionInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VersionInfo;
  static deserializeBinaryFromReader(message: VersionInfo, reader: jspb.BinaryReader): VersionInfo;
}

export namespace VersionInfo {
  export type AsObject = {
    api?: VersionInfo.ProtocolVersion.AsObject,
    entrypoint?: VersionInfo.ProtocolVersion.AsObject,
    version: string,
  }

  export class ProtocolVersion extends jspb.Message {
    getCurrent(): number;
    setCurrent(value: number): ProtocolVersion;

    getMinimum(): number;
    setMinimum(value: number): ProtocolVersion;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ProtocolVersion.AsObject;
    static toObject(includeInstance: boolean, msg: ProtocolVersion): ProtocolVersion.AsObject;
    static serializeBinaryToWriter(message: ProtocolVersion, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ProtocolVersion;
    static deserializeBinaryFromReader(message: ProtocolVersion, reader: jspb.BinaryReader): ProtocolVersion;
  }

  export namespace ProtocolVersion {
    export type AsObject = {
      current: number,
      minimum: number,
    }
  }

}

export class Application extends jspb.Message {
  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): Application;
  hasProject(): boolean;
  clearProject(): Application;

  getName(): string;
  setName(value: string): Application;

  getFileChangeSignal(): string;
  setFileChangeSignal(value: string): Application;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Application.AsObject;
  static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
  static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Application;
  static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
}

export namespace Application {
  export type AsObject = {
    project?: Ref.Project.AsObject,
    name: string,
    fileChangeSignal: string,
  }
}

export class Variable extends jspb.Message {
  getName(): string;
  setName(value: string): Variable;

  getStr(): string;
  setStr(value: string): Variable;

  getBool(): boolean;
  setBool(value: boolean): Variable;

  getNum(): number;
  setNum(value: number): Variable;

  getHcl(): string;
  setHcl(value: string): Variable;

  getCli(): google_protobuf_empty_pb.Empty | undefined;
  setCli(value?: google_protobuf_empty_pb.Empty): Variable;
  hasCli(): boolean;
  clearCli(): Variable;

  getFile(): Variable.File | undefined;
  setFile(value?: Variable.File): Variable;
  hasFile(): boolean;
  clearFile(): Variable;

  getEnv(): google_protobuf_empty_pb.Empty | undefined;
  setEnv(value?: google_protobuf_empty_pb.Empty): Variable;
  hasEnv(): boolean;
  clearEnv(): Variable;

  getVcs(): Variable.VCS | undefined;
  setVcs(value?: Variable.VCS): Variable;
  hasVcs(): boolean;
  clearVcs(): Variable;

  getServer(): google_protobuf_empty_pb.Empty | undefined;
  setServer(value?: google_protobuf_empty_pb.Empty): Variable;
  hasServer(): boolean;
  clearServer(): Variable;

  getDynamic(): google_protobuf_empty_pb.Empty | undefined;
  setDynamic(value?: google_protobuf_empty_pb.Empty): Variable;
  hasDynamic(): boolean;
  clearDynamic(): Variable;

  getFinalValue(): Variable.FinalValue | undefined;
  setFinalValue(value?: Variable.FinalValue): Variable;
  hasFinalValue(): boolean;
  clearFinalValue(): Variable;

  getSensitive(): boolean;
  setSensitive(value: boolean): Variable;

  getValueCase(): Variable.ValueCase;

  getSourceCase(): Variable.SourceCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Variable.AsObject;
  static toObject(includeInstance: boolean, msg: Variable): Variable.AsObject;
  static serializeBinaryToWriter(message: Variable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Variable;
  static deserializeBinaryFromReader(message: Variable, reader: jspb.BinaryReader): Variable;
}

export namespace Variable {
  export type AsObject = {
    name: string,
    str: string,
    bool: boolean,
    num: number,
    hcl: string,
    cli?: google_protobuf_empty_pb.Empty.AsObject,
    file?: Variable.File.AsObject,
    env?: google_protobuf_empty_pb.Empty.AsObject,
    vcs?: Variable.VCS.AsObject,
    server?: google_protobuf_empty_pb.Empty.AsObject,
    dynamic?: google_protobuf_empty_pb.Empty.AsObject,
    finalValue?: Variable.FinalValue.AsObject,
    sensitive: boolean,
  }

  export class File extends jspb.Message {
    getFilename(): string;
    setFilename(value: string): File;

    getHclRange(): Variable.HclRange | undefined;
    setHclRange(value?: Variable.HclRange): File;
    hasHclRange(): boolean;
    clearHclRange(): File;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): File.AsObject;
    static toObject(includeInstance: boolean, msg: File): File.AsObject;
    static serializeBinaryToWriter(message: File, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): File;
    static deserializeBinaryFromReader(message: File, reader: jspb.BinaryReader): File;
  }

  export namespace File {
    export type AsObject = {
      filename: string,
      hclRange?: Variable.HclRange.AsObject,
    }
  }


  export class VCS extends jspb.Message {
    getFilename(): string;
    setFilename(value: string): VCS;

    getHclRange(): Variable.HclRange | undefined;
    setHclRange(value?: Variable.HclRange): VCS;
    hasHclRange(): boolean;
    clearHclRange(): VCS;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): VCS.AsObject;
    static toObject(includeInstance: boolean, msg: VCS): VCS.AsObject;
    static serializeBinaryToWriter(message: VCS, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): VCS;
    static deserializeBinaryFromReader(message: VCS, reader: jspb.BinaryReader): VCS;
  }

  export namespace VCS {
    export type AsObject = {
      filename: string,
      hclRange?: Variable.HclRange.AsObject,
    }
  }


  export class HclPos extends jspb.Message {
    getLine(): number;
    setLine(value: number): HclPos;

    getColumn(): number;
    setColumn(value: number): HclPos;

    getByte(): number;
    setByte(value: number): HclPos;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HclPos.AsObject;
    static toObject(includeInstance: boolean, msg: HclPos): HclPos.AsObject;
    static serializeBinaryToWriter(message: HclPos, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HclPos;
    static deserializeBinaryFromReader(message: HclPos, reader: jspb.BinaryReader): HclPos;
  }

  export namespace HclPos {
    export type AsObject = {
      line: number,
      column: number,
      pb_byte: number,
    }
  }


  export class HclRange extends jspb.Message {
    getFilename(): string;
    setFilename(value: string): HclRange;

    getStart(): Variable.HclPos | undefined;
    setStart(value?: Variable.HclPos): HclRange;
    hasStart(): boolean;
    clearStart(): HclRange;

    getEnd(): Variable.HclPos | undefined;
    setEnd(value?: Variable.HclPos): HclRange;
    hasEnd(): boolean;
    clearEnd(): HclRange;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HclRange.AsObject;
    static toObject(includeInstance: boolean, msg: HclRange): HclRange.AsObject;
    static serializeBinaryToWriter(message: HclRange, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HclRange;
    static deserializeBinaryFromReader(message: HclRange, reader: jspb.BinaryReader): HclRange;
  }

  export namespace HclRange {
    export type AsObject = {
      filename: string,
      start?: Variable.HclPos.AsObject,
      end?: Variable.HclPos.AsObject,
    }
  }


  export class FinalValue extends jspb.Message {
    getSensitive(): string;
    setSensitive(value: string): FinalValue;

    getStr(): string;
    setStr(value: string): FinalValue;

    getBool(): boolean;
    setBool(value: boolean): FinalValue;

    getNum(): number;
    setNum(value: number): FinalValue;

    getHcl(): string;
    setHcl(value: string): FinalValue;

    getSource(): Variable.FinalValue.Source;
    setSource(value: Variable.FinalValue.Source): FinalValue;

    getValueCase(): FinalValue.ValueCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): FinalValue.AsObject;
    static toObject(includeInstance: boolean, msg: FinalValue): FinalValue.AsObject;
    static serializeBinaryToWriter(message: FinalValue, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): FinalValue;
    static deserializeBinaryFromReader(message: FinalValue, reader: jspb.BinaryReader): FinalValue;
  }

  export namespace FinalValue {
    export type AsObject = {
      sensitive: string,
      str: string,
      bool: boolean,
      num: number,
      hcl: string,
      source: Variable.FinalValue.Source,
    }

    export enum Source { 
      UNKNOWN = 0,
      DEFAULT = 1,
      FILE = 2,
      CLI = 3,
      ENV = 4,
      VCS = 5,
      SERVER = 6,
      DYNAMIC = 7,
    }

    export enum ValueCase { 
      VALUE_NOT_SET = 0,
      SENSITIVE = 1,
      STR = 2,
      BOOL = 3,
      NUM = 4,
      HCL = 5,
    }
  }


  export enum ValueCase { 
    VALUE_NOT_SET = 0,
    STR = 2,
    BOOL = 9,
    NUM = 10,
    HCL = 3,
  }

  export enum SourceCase { 
    SOURCE_NOT_SET = 0,
    CLI = 4,
    FILE = 5,
    ENV = 6,
    VCS = 7,
    SERVER = 8,
    DYNAMIC = 11,
  }
}

export class Project extends jspb.Message {
  getName(): string;
  setName(value: string): Project;

  getApplicationsList(): Array<Application>;
  setApplicationsList(value: Array<Application>): Project;
  clearApplicationsList(): Project;
  addApplications(value?: Application, index?: number): Application;

  getRemoteEnabled(): boolean;
  setRemoteEnabled(value: boolean): Project;

  getDataSource(): Job.DataSource | undefined;
  setDataSource(value?: Job.DataSource): Project;
  hasDataSource(): boolean;
  clearDataSource(): Project;

  getDataSourcePoll(): Project.Poll | undefined;
  setDataSourcePoll(value?: Project.Poll): Project;
  hasDataSourcePoll(): boolean;
  clearDataSourcePoll(): Project;

  getWaypointHcl(): Uint8Array | string;
  getWaypointHcl_asU8(): Uint8Array;
  getWaypointHcl_asB64(): string;
  setWaypointHcl(value: Uint8Array | string): Project;

  getWaypointHclFormat(): Hcl.Format;
  setWaypointHclFormat(value: Hcl.Format): Project;

  getFileChangeSignal(): string;
  setFileChangeSignal(value: string): Project;

  getVariablesList(): Array<Variable>;
  setVariablesList(value: Array<Variable>): Project;
  clearVariablesList(): Project;
  addVariables(value?: Variable, index?: number): Variable;

  getStatusReportPoll(): Project.AppStatusPoll | undefined;
  setStatusReportPoll(value?: Project.AppStatusPoll): Project;
  hasStatusReportPoll(): boolean;
  clearStatusReportPoll(): Project;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Project.AsObject;
  static toObject(includeInstance: boolean, msg: Project): Project.AsObject;
  static serializeBinaryToWriter(message: Project, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Project;
  static deserializeBinaryFromReader(message: Project, reader: jspb.BinaryReader): Project;
}

export namespace Project {
  export type AsObject = {
    name: string,
    applicationsList: Array<Application.AsObject>,
    remoteEnabled: boolean,
    dataSource?: Job.DataSource.AsObject,
    dataSourcePoll?: Project.Poll.AsObject,
    waypointHcl: Uint8Array | string,
    waypointHclFormat: Hcl.Format,
    fileChangeSignal: string,
    variablesList: Array<Variable.AsObject>,
    statusReportPoll?: Project.AppStatusPoll.AsObject,
  }

  export class Poll extends jspb.Message {
    getEnabled(): boolean;
    setEnabled(value: boolean): Poll;

    getInterval(): string;
    setInterval(value: string): Poll;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Poll.AsObject;
    static toObject(includeInstance: boolean, msg: Poll): Poll.AsObject;
    static serializeBinaryToWriter(message: Poll, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Poll;
    static deserializeBinaryFromReader(message: Poll, reader: jspb.BinaryReader): Poll;
  }

  export namespace Poll {
    export type AsObject = {
      enabled: boolean,
      interval: string,
    }
  }


  export class AppStatusPoll extends jspb.Message {
    getEnabled(): boolean;
    setEnabled(value: boolean): AppStatusPoll;

    getInterval(): string;
    setInterval(value: string): AppStatusPoll;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppStatusPoll.AsObject;
    static toObject(includeInstance: boolean, msg: AppStatusPoll): AppStatusPoll.AsObject;
    static serializeBinaryToWriter(message: AppStatusPoll, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppStatusPoll;
    static deserializeBinaryFromReader(message: AppStatusPoll, reader: jspb.BinaryReader): AppStatusPoll;
  }

  export namespace AppStatusPoll {
    export type AsObject = {
      enabled: boolean,
      interval: string,
    }
  }

}

export class Workspace extends jspb.Message {
  getName(): string;
  setName(value: string): Workspace;

  getProjectsList(): Array<Workspace.Project>;
  setProjectsList(value: Array<Workspace.Project>): Workspace;
  clearProjectsList(): Workspace;
  addProjects(value?: Workspace.Project, index?: number): Workspace.Project;

  getActiveTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setActiveTime(value?: google_protobuf_timestamp_pb.Timestamp): Workspace;
  hasActiveTime(): boolean;
  clearActiveTime(): Workspace;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Workspace.AsObject;
  static toObject(includeInstance: boolean, msg: Workspace): Workspace.AsObject;
  static serializeBinaryToWriter(message: Workspace, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Workspace;
  static deserializeBinaryFromReader(message: Workspace, reader: jspb.BinaryReader): Workspace;
}

export namespace Workspace {
  export type AsObject = {
    name: string,
    projectsList: Array<Workspace.Project.AsObject>,
    activeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }

  export class Project extends jspb.Message {
    getProject(): Ref.Project | undefined;
    setProject(value?: Ref.Project): Project;
    hasProject(): boolean;
    clearProject(): Project;

    getWorkspace(): Ref.Workspace | undefined;
    setWorkspace(value?: Ref.Workspace): Project;
    hasWorkspace(): boolean;
    clearWorkspace(): Project;

    getDataSourceRef(): Job.DataSource.Ref | undefined;
    setDataSourceRef(value?: Job.DataSource.Ref): Project;
    hasDataSourceRef(): boolean;
    clearDataSourceRef(): Project;

    getActiveTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setActiveTime(value?: google_protobuf_timestamp_pb.Timestamp): Project;
    hasActiveTime(): boolean;
    clearActiveTime(): Project;

    getApplicationsList(): Array<Workspace.Application>;
    setApplicationsList(value: Array<Workspace.Application>): Project;
    clearApplicationsList(): Project;
    addApplications(value?: Workspace.Application, index?: number): Workspace.Application;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Project.AsObject;
    static toObject(includeInstance: boolean, msg: Project): Project.AsObject;
    static serializeBinaryToWriter(message: Project, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Project;
    static deserializeBinaryFromReader(message: Project, reader: jspb.BinaryReader): Project;
  }

  export namespace Project {
    export type AsObject = {
      project?: Ref.Project.AsObject,
      workspace?: Ref.Workspace.AsObject,
      dataSourceRef?: Job.DataSource.Ref.AsObject,
      activeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      applicationsList: Array<Workspace.Application.AsObject>,
    }
  }


  export class Application extends jspb.Message {
    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): Application;
    hasApplication(): boolean;
    clearApplication(): Application;

    getActiveTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setActiveTime(value?: google_protobuf_timestamp_pb.Timestamp): Application;
    hasActiveTime(): boolean;
    clearActiveTime(): Application;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Application.AsObject;
    static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
    static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Application;
    static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
  }

  export namespace Application {
    export type AsObject = {
      application?: Ref.Application.AsObject,
      activeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
  }

}

export class User extends jspb.Message {
  getId(): string;
  setId(value: string): User;

  getUsername(): string;
  setUsername(value: string): User;

  getDisplay(): string;
  setDisplay(value: string): User;

  getEmail(): string;
  setEmail(value: string): User;

  getLinksList(): Array<User.Link>;
  setLinksList(value: Array<User.Link>): User;
  clearLinksList(): User;
  addLinks(value?: User.Link, index?: number): User.Link;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: string,
    username: string,
    display: string,
    email: string,
    linksList: Array<User.Link.AsObject>,
  }

  export class Link extends jspb.Message {
    getOidc(): User.Link.OIDC | undefined;
    setOidc(value?: User.Link.OIDC): Link;
    hasOidc(): boolean;
    clearOidc(): Link;

    getMethodCase(): Link.MethodCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Link.AsObject;
    static toObject(includeInstance: boolean, msg: Link): Link.AsObject;
    static serializeBinaryToWriter(message: Link, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Link;
    static deserializeBinaryFromReader(message: Link, reader: jspb.BinaryReader): Link;
  }

  export namespace Link {
    export type AsObject = {
      oidc?: User.Link.OIDC.AsObject,
    }

    export class OIDC extends jspb.Message {
      getIss(): string;
      setIss(value: string): OIDC;

      getSub(): string;
      setSub(value: string): OIDC;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): OIDC.AsObject;
      static toObject(includeInstance: boolean, msg: OIDC): OIDC.AsObject;
      static serializeBinaryToWriter(message: OIDC, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): OIDC;
      static deserializeBinaryFromReader(message: OIDC, reader: jspb.BinaryReader): OIDC;
    }

    export namespace OIDC {
      export type AsObject = {
        iss: string,
        sub: string,
      }
    }


    export enum MethodCase { 
      METHOD_NOT_SET = 0,
      OIDC = 1,
    }
  }

}

export class Ref extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Ref.AsObject;
  static toObject(includeInstance: boolean, msg: Ref): Ref.AsObject;
  static serializeBinaryToWriter(message: Ref, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Ref;
  static deserializeBinaryFromReader(message: Ref, reader: jspb.BinaryReader): Ref;
}

export namespace Ref {
  export type AsObject = {
  }

  export class Global extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Global.AsObject;
    static toObject(includeInstance: boolean, msg: Global): Global.AsObject;
    static serializeBinaryToWriter(message: Global, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Global;
    static deserializeBinaryFromReader(message: Global, reader: jspb.BinaryReader): Global;
  }

  export namespace Global {
    export type AsObject = {
    }
  }


  export class Application extends jspb.Message {
    getApplication(): string;
    setApplication(value: string): Application;

    getProject(): string;
    setProject(value: string): Application;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Application.AsObject;
    static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
    static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Application;
    static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
  }

  export namespace Application {
    export type AsObject = {
      application: string,
      project: string,
    }
  }


  export class Project extends jspb.Message {
    getProject(): string;
    setProject(value: string): Project;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Project.AsObject;
    static toObject(includeInstance: boolean, msg: Project): Project.AsObject;
    static serializeBinaryToWriter(message: Project, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Project;
    static deserializeBinaryFromReader(message: Project, reader: jspb.BinaryReader): Project;
  }

  export namespace Project {
    export type AsObject = {
      project: string,
    }
  }


  export class Workspace extends jspb.Message {
    getWorkspace(): string;
    setWorkspace(value: string): Workspace;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Workspace.AsObject;
    static toObject(includeInstance: boolean, msg: Workspace): Workspace.AsObject;
    static serializeBinaryToWriter(message: Workspace, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Workspace;
    static deserializeBinaryFromReader(message: Workspace, reader: jspb.BinaryReader): Workspace;
  }

  export namespace Workspace {
    export type AsObject = {
      workspace: string,
    }
  }


  export class Deployment extends jspb.Message {
    getLatest(): boolean;
    setLatest(value: boolean): Deployment;

    getSequence(): number;
    setSequence(value: number): Deployment;

    getRefCase(): Deployment.RefCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Deployment.AsObject;
    static toObject(includeInstance: boolean, msg: Deployment): Deployment.AsObject;
    static serializeBinaryToWriter(message: Deployment, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Deployment;
    static deserializeBinaryFromReader(message: Deployment, reader: jspb.BinaryReader): Deployment;
  }

  export namespace Deployment {
    export type AsObject = {
      latest: boolean,
      sequence: number,
    }

    export enum RefCase { 
      REF_NOT_SET = 0,
      LATEST = 1,
      SEQUENCE = 2,
    }
  }


  export class Artifact extends jspb.Message {
    getLatest(): boolean;
    setLatest(value: boolean): Artifact;

    getSequence(): number;
    setSequence(value: number): Artifact;

    getRefCase(): Artifact.RefCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Artifact.AsObject;
    static toObject(includeInstance: boolean, msg: Artifact): Artifact.AsObject;
    static serializeBinaryToWriter(message: Artifact, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Artifact;
    static deserializeBinaryFromReader(message: Artifact, reader: jspb.BinaryReader): Artifact;
  }

  export namespace Artifact {
    export type AsObject = {
      latest: boolean,
      sequence: number,
    }

    export enum RefCase { 
      REF_NOT_SET = 0,
      LATEST = 1,
      SEQUENCE = 2,
    }
  }


  export class Component extends jspb.Message {
    getType(): Component.Type;
    setType(value: Component.Type): Component;

    getName(): string;
    setName(value: string): Component;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Component.AsObject;
    static toObject(includeInstance: boolean, msg: Component): Component.AsObject;
    static serializeBinaryToWriter(message: Component, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Component;
    static deserializeBinaryFromReader(message: Component, reader: jspb.BinaryReader): Component;
  }

  export namespace Component {
    export type AsObject = {
      type: Component.Type,
      name: string,
    }
  }


  export class Operation extends jspb.Message {
    getId(): string;
    setId(value: string): Operation;

    getSequence(): Ref.OperationSeq | undefined;
    setSequence(value?: Ref.OperationSeq): Operation;
    hasSequence(): boolean;
    clearSequence(): Operation;

    getTargetCase(): Operation.TargetCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Operation.AsObject;
    static toObject(includeInstance: boolean, msg: Operation): Operation.AsObject;
    static serializeBinaryToWriter(message: Operation, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Operation;
    static deserializeBinaryFromReader(message: Operation, reader: jspb.BinaryReader): Operation;
  }

  export namespace Operation {
    export type AsObject = {
      id: string,
      sequence?: Ref.OperationSeq.AsObject,
    }

    export enum TargetCase { 
      TARGET_NOT_SET = 0,
      ID = 1,
      SEQUENCE = 2,
    }
  }


  export class OperationSeq extends jspb.Message {
    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): OperationSeq;
    hasApplication(): boolean;
    clearApplication(): OperationSeq;

    getNumber(): number;
    setNumber(value: number): OperationSeq;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OperationSeq.AsObject;
    static toObject(includeInstance: boolean, msg: OperationSeq): OperationSeq.AsObject;
    static serializeBinaryToWriter(message: OperationSeq, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OperationSeq;
    static deserializeBinaryFromReader(message: OperationSeq, reader: jspb.BinaryReader): OperationSeq;
  }

  export namespace OperationSeq {
    export type AsObject = {
      application?: Ref.Application.AsObject,
      number: number,
    }
  }


  export class Runner extends jspb.Message {
    getAny(): Ref.RunnerAny | undefined;
    setAny(value?: Ref.RunnerAny): Runner;
    hasAny(): boolean;
    clearAny(): Runner;

    getId(): Ref.RunnerId | undefined;
    setId(value?: Ref.RunnerId): Runner;
    hasId(): boolean;
    clearId(): Runner;

    getLabels(): Ref.RunnerLabels | undefined;
    setLabels(value?: Ref.RunnerLabels): Runner;
    hasLabels(): boolean;
    clearLabels(): Runner;

    getTargetCase(): Runner.TargetCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Runner.AsObject;
    static toObject(includeInstance: boolean, msg: Runner): Runner.AsObject;
    static serializeBinaryToWriter(message: Runner, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Runner;
    static deserializeBinaryFromReader(message: Runner, reader: jspb.BinaryReader): Runner;
  }

  export namespace Runner {
    export type AsObject = {
      any?: Ref.RunnerAny.AsObject,
      id?: Ref.RunnerId.AsObject,
      labels?: Ref.RunnerLabels.AsObject,
    }

    export enum TargetCase { 
      TARGET_NOT_SET = 0,
      ANY = 1,
      ID = 2,
      LABELS = 3,
    }
  }


  export class RunnerAny extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RunnerAny.AsObject;
    static toObject(includeInstance: boolean, msg: RunnerAny): RunnerAny.AsObject;
    static serializeBinaryToWriter(message: RunnerAny, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RunnerAny;
    static deserializeBinaryFromReader(message: RunnerAny, reader: jspb.BinaryReader): RunnerAny;
  }

  export namespace RunnerAny {
    export type AsObject = {
    }
  }


  export class RunnerId extends jspb.Message {
    getId(): string;
    setId(value: string): RunnerId;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RunnerId.AsObject;
    static toObject(includeInstance: boolean, msg: RunnerId): RunnerId.AsObject;
    static serializeBinaryToWriter(message: RunnerId, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RunnerId;
    static deserializeBinaryFromReader(message: RunnerId, reader: jspb.BinaryReader): RunnerId;
  }

  export namespace RunnerId {
    export type AsObject = {
      id: string,
    }
  }


  export class RunnerLabels extends jspb.Message {
    getLabelsMap(): jspb.Map<string, string>;
    clearLabelsMap(): RunnerLabels;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RunnerLabels.AsObject;
    static toObject(includeInstance: boolean, msg: RunnerLabels): RunnerLabels.AsObject;
    static serializeBinaryToWriter(message: RunnerLabels, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RunnerLabels;
    static deserializeBinaryFromReader(message: RunnerLabels, reader: jspb.BinaryReader): RunnerLabels;
  }

  export namespace RunnerLabels {
    export type AsObject = {
      labelsMap: Array<[string, string]>,
    }
  }


  export class User extends jspb.Message {
    getId(): Ref.UserId | undefined;
    setId(value?: Ref.UserId): User;
    hasId(): boolean;
    clearId(): User;

    getUsername(): Ref.UserUsername | undefined;
    setUsername(value?: Ref.UserUsername): User;
    hasUsername(): boolean;
    clearUsername(): User;

    getRefCase(): User.RefCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): User.AsObject;
    static toObject(includeInstance: boolean, msg: User): User.AsObject;
    static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): User;
    static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
  }

  export namespace User {
    export type AsObject = {
      id?: Ref.UserId.AsObject,
      username?: Ref.UserUsername.AsObject,
    }

    export enum RefCase { 
      REF_NOT_SET = 0,
      ID = 1,
      USERNAME = 2,
    }
  }


  export class UserId extends jspb.Message {
    getId(): string;
    setId(value: string): UserId;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UserId.AsObject;
    static toObject(includeInstance: boolean, msg: UserId): UserId.AsObject;
    static serializeBinaryToWriter(message: UserId, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UserId;
    static deserializeBinaryFromReader(message: UserId, reader: jspb.BinaryReader): UserId;
  }

  export namespace UserId {
    export type AsObject = {
      id: string,
    }
  }


  export class UserUsername extends jspb.Message {
    getUsername(): string;
    setUsername(value: string): UserUsername;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UserUsername.AsObject;
    static toObject(includeInstance: boolean, msg: UserUsername): UserUsername.AsObject;
    static serializeBinaryToWriter(message: UserUsername, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UserUsername;
    static deserializeBinaryFromReader(message: UserUsername, reader: jspb.BinaryReader): UserUsername;
  }

  export namespace UserUsername {
    export type AsObject = {
      username: string,
    }
  }


  export class AuthMethod extends jspb.Message {
    getName(): string;
    setName(value: string): AuthMethod;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AuthMethod.AsObject;
    static toObject(includeInstance: boolean, msg: AuthMethod): AuthMethod.AsObject;
    static serializeBinaryToWriter(message: AuthMethod, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AuthMethod;
    static deserializeBinaryFromReader(message: AuthMethod, reader: jspb.BinaryReader): AuthMethod;
  }

  export namespace AuthMethod {
    export type AsObject = {
      name: string,
    }
  }


  export class DeclaredResource extends jspb.Message {
    getName(): string;
    setName(value: string): DeclaredResource;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeclaredResource.AsObject;
    static toObject(includeInstance: boolean, msg: DeclaredResource): DeclaredResource.AsObject;
    static serializeBinaryToWriter(message: DeclaredResource, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeclaredResource;
    static deserializeBinaryFromReader(message: DeclaredResource, reader: jspb.BinaryReader): DeclaredResource;
  }

  export namespace DeclaredResource {
    export type AsObject = {
      name: string,
    }
  }


  export class OnDemandRunnerConfig extends jspb.Message {
    getId(): string;
    setId(value: string): OnDemandRunnerConfig;

    getName(): string;
    setName(value: string): OnDemandRunnerConfig;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OnDemandRunnerConfig.AsObject;
    static toObject(includeInstance: boolean, msg: OnDemandRunnerConfig): OnDemandRunnerConfig.AsObject;
    static serializeBinaryToWriter(message: OnDemandRunnerConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OnDemandRunnerConfig;
    static deserializeBinaryFromReader(message: OnDemandRunnerConfig, reader: jspb.BinaryReader): OnDemandRunnerConfig;
  }

  export namespace OnDemandRunnerConfig {
    export type AsObject = {
      id: string,
      name: string,
    }
  }


  export class Trigger extends jspb.Message {
    getId(): string;
    setId(value: string): Trigger;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Trigger.AsObject;
    static toObject(includeInstance: boolean, msg: Trigger): Trigger.AsObject;
    static serializeBinaryToWriter(message: Trigger, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Trigger;
    static deserializeBinaryFromReader(message: Trigger, reader: jspb.BinaryReader): Trigger;
  }

  export namespace Trigger {
    export type AsObject = {
      id: string,
    }
  }


  export class Job extends jspb.Message {
    getId(): string;
    setId(value: string): Job;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Job.AsObject;
    static toObject(includeInstance: boolean, msg: Job): Job.AsObject;
    static serializeBinaryToWriter(message: Job, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Job;
    static deserializeBinaryFromReader(message: Job, reader: jspb.BinaryReader): Job;
  }

  export namespace Job {
    export type AsObject = {
      id: string,
    }
  }


  export class Task extends jspb.Message {
    getId(): string;
    setId(value: string): Task;

    getJobId(): string;
    setJobId(value: string): Task;

    getRefCase(): Task.RefCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Task.AsObject;
    static toObject(includeInstance: boolean, msg: Task): Task.AsObject;
    static serializeBinaryToWriter(message: Task, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Task;
    static deserializeBinaryFromReader(message: Task, reader: jspb.BinaryReader): Task;
  }

  export namespace Task {
    export type AsObject = {
      id: string,
      jobId: string,
    }

    export enum RefCase { 
      REF_NOT_SET = 0,
      ID = 1,
      JOB_ID = 2,
    }
  }


  export class Pipeline extends jspb.Message {
    getId(): Ref.PipelineId | undefined;
    setId(value?: Ref.PipelineId): Pipeline;
    hasId(): boolean;
    clearId(): Pipeline;

    getOwner(): Ref.PipelineOwner | undefined;
    setOwner(value?: Ref.PipelineOwner): Pipeline;
    hasOwner(): boolean;
    clearOwner(): Pipeline;

    getRefCase(): Pipeline.RefCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Pipeline.AsObject;
    static toObject(includeInstance: boolean, msg: Pipeline): Pipeline.AsObject;
    static serializeBinaryToWriter(message: Pipeline, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Pipeline;
    static deserializeBinaryFromReader(message: Pipeline, reader: jspb.BinaryReader): Pipeline;
  }

  export namespace Pipeline {
    export type AsObject = {
      id?: Ref.PipelineId.AsObject,
      owner?: Ref.PipelineOwner.AsObject,
    }

    export enum RefCase { 
      REF_NOT_SET = 0,
      ID = 1,
      OWNER = 2,
    }
  }


  export class PipelineId extends jspb.Message {
    getId(): string;
    setId(value: string): PipelineId;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PipelineId.AsObject;
    static toObject(includeInstance: boolean, msg: PipelineId): PipelineId.AsObject;
    static serializeBinaryToWriter(message: PipelineId, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PipelineId;
    static deserializeBinaryFromReader(message: PipelineId, reader: jspb.BinaryReader): PipelineId;
  }

  export namespace PipelineId {
    export type AsObject = {
      id: string,
    }
  }


  export class PipelineOwner extends jspb.Message {
    getProject(): Ref.Project | undefined;
    setProject(value?: Ref.Project): PipelineOwner;
    hasProject(): boolean;
    clearProject(): PipelineOwner;

    getPipelineName(): string;
    setPipelineName(value: string): PipelineOwner;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PipelineOwner.AsObject;
    static toObject(includeInstance: boolean, msg: PipelineOwner): PipelineOwner.AsObject;
    static serializeBinaryToWriter(message: PipelineOwner, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PipelineOwner;
    static deserializeBinaryFromReader(message: PipelineOwner, reader: jspb.BinaryReader): PipelineOwner;
  }

  export namespace PipelineOwner {
    export type AsObject = {
      project?: Ref.Project.AsObject,
      pipelineName: string,
    }
  }


  export class PipelineStep extends jspb.Message {
    getPipeline(): string;
    setPipeline(value: string): PipelineStep;

    getStep(): string;
    setStep(value: string): PipelineStep;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PipelineStep.AsObject;
    static toObject(includeInstance: boolean, msg: PipelineStep): PipelineStep.AsObject;
    static serializeBinaryToWriter(message: PipelineStep, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PipelineStep;
    static deserializeBinaryFromReader(message: PipelineStep, reader: jspb.BinaryReader): PipelineStep;
  }

  export namespace PipelineStep {
    export type AsObject = {
      pipeline: string,
      step: string,
    }
  }

}

export class Component extends jspb.Message {
  getType(): Component.Type;
  setType(value: Component.Type): Component;

  getName(): string;
  setName(value: string): Component;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Component.AsObject;
  static toObject(includeInstance: boolean, msg: Component): Component.AsObject;
  static serializeBinaryToWriter(message: Component, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Component;
  static deserializeBinaryFromReader(message: Component, reader: jspb.BinaryReader): Component;
}

export namespace Component {
  export type AsObject = {
    type: Component.Type,
    name: string,
  }

  export enum Type { 
    UNKNOWN = 0,
    BUILDER = 1,
    REGISTRY = 2,
    PLATFORM = 3,
    RELEASEMANAGER = 4,
  }
}

export class Status extends jspb.Message {
  getState(): Status.State;
  setState(value: Status.State): Status;

  getDetails(): string;
  setDetails(value: string): Status;

  getError(): google_rpc_status_pb.Status | undefined;
  setError(value?: google_rpc_status_pb.Status): Status;
  hasError(): boolean;
  clearError(): Status;

  getStartTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setStartTime(value?: google_protobuf_timestamp_pb.Timestamp): Status;
  hasStartTime(): boolean;
  clearStartTime(): Status;

  getCompleteTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCompleteTime(value?: google_protobuf_timestamp_pb.Timestamp): Status;
  hasCompleteTime(): boolean;
  clearCompleteTime(): Status;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Status.AsObject;
  static toObject(includeInstance: boolean, msg: Status): Status.AsObject;
  static serializeBinaryToWriter(message: Status, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Status;
  static deserializeBinaryFromReader(message: Status, reader: jspb.BinaryReader): Status;
}

export namespace Status {
  export type AsObject = {
    state: Status.State,
    details: string,
    error?: google_rpc_status_pb.Status.AsObject,
    startTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    completeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }

  export enum State { 
    UNKNOWN = 0,
    RUNNING = 1,
    SUCCESS = 2,
    ERROR = 3,
  }
}

export class StatusFilter extends jspb.Message {
  getFiltersList(): Array<StatusFilter.Filter>;
  setFiltersList(value: Array<StatusFilter.Filter>): StatusFilter;
  clearFiltersList(): StatusFilter;
  addFilters(value?: StatusFilter.Filter, index?: number): StatusFilter.Filter;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusFilter.AsObject;
  static toObject(includeInstance: boolean, msg: StatusFilter): StatusFilter.AsObject;
  static serializeBinaryToWriter(message: StatusFilter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusFilter;
  static deserializeBinaryFromReader(message: StatusFilter, reader: jspb.BinaryReader): StatusFilter;
}

export namespace StatusFilter {
  export type AsObject = {
    filtersList: Array<StatusFilter.Filter.AsObject>,
  }

  export class Filter extends jspb.Message {
    getState(): Status.State;
    setState(value: Status.State): Filter;

    getFilterCase(): Filter.FilterCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Filter.AsObject;
    static toObject(includeInstance: boolean, msg: Filter): Filter.AsObject;
    static serializeBinaryToWriter(message: Filter, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Filter;
    static deserializeBinaryFromReader(message: Filter, reader: jspb.BinaryReader): Filter;
  }

  export namespace Filter {
    export type AsObject = {
      state: Status.State,
    }

    export enum FilterCase { 
      FILTER_NOT_SET = 0,
      STATE = 2,
    }
  }

}

export class Operation extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Operation.AsObject;
  static toObject(includeInstance: boolean, msg: Operation): Operation.AsObject;
  static serializeBinaryToWriter(message: Operation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Operation;
  static deserializeBinaryFromReader(message: Operation, reader: jspb.BinaryReader): Operation;
}

export namespace Operation {
  export type AsObject = {
  }

  export enum PhysicalState { 
    UNKNOWN = 0,
    PENDING = 1,
    CREATED = 3,
    DESTROYED = 4,
  }
}

export class OperationOrder extends jspb.Message {
  getOrder(): OperationOrder.Order;
  setOrder(value: OperationOrder.Order): OperationOrder;

  getDesc(): boolean;
  setDesc(value: boolean): OperationOrder;

  getLimit(): number;
  setLimit(value: number): OperationOrder;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OperationOrder.AsObject;
  static toObject(includeInstance: boolean, msg: OperationOrder): OperationOrder.AsObject;
  static serializeBinaryToWriter(message: OperationOrder, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OperationOrder;
  static deserializeBinaryFromReader(message: OperationOrder, reader: jspb.BinaryReader): OperationOrder;
}

export namespace OperationOrder {
  export type AsObject = {
    order: OperationOrder.Order,
    desc: boolean,
    limit: number,
  }

  export enum Order { 
    UNSET = 0,
    START_TIME = 1,
    COMPLETE_TIME = 2,
  }
}

export class Generation extends jspb.Message {
  getId(): string;
  setId(value: string): Generation;

  getInitialSequence(): number;
  setInitialSequence(value: number): Generation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Generation.AsObject;
  static toObject(includeInstance: boolean, msg: Generation): Generation.AsObject;
  static serializeBinaryToWriter(message: Generation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Generation;
  static deserializeBinaryFromReader(message: Generation, reader: jspb.BinaryReader): Generation;
}

export namespace Generation {
  export type AsObject = {
    id: string,
    initialSequence: number,
  }
}

export class DeclaredResource extends jspb.Message {
  getName(): string;
  setName(value: string): DeclaredResource;

  getType(): string;
  setType(value: string): DeclaredResource;

  getPlatform(): string;
  setPlatform(value: string): DeclaredResource;

  getState(): any_pb.Any | undefined;
  setState(value?: any_pb.Any): DeclaredResource;
  hasState(): boolean;
  clearState(): DeclaredResource;

  getStateJson(): string;
  setStateJson(value: string): DeclaredResource;

  getCategoryDisplayHint(): ResourceCategoryDisplayHint;
  setCategoryDisplayHint(value: ResourceCategoryDisplayHint): DeclaredResource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeclaredResource.AsObject;
  static toObject(includeInstance: boolean, msg: DeclaredResource): DeclaredResource.AsObject;
  static serializeBinaryToWriter(message: DeclaredResource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeclaredResource;
  static deserializeBinaryFromReader(message: DeclaredResource, reader: jspb.BinaryReader): DeclaredResource;
}

export namespace DeclaredResource {
  export type AsObject = {
    name: string,
    type: string,
    platform: string,
    state?: any_pb.Any.AsObject,
    stateJson: string,
    categoryDisplayHint: ResourceCategoryDisplayHint,
  }
}

export class TaskLaunchInfo extends jspb.Message {
  getOciUrl(): string;
  setOciUrl(value: string): TaskLaunchInfo;

  getEnvironmentVariablesMap(): jspb.Map<string, string>;
  clearEnvironmentVariablesMap(): TaskLaunchInfo;

  getEntrypointList(): Array<string>;
  setEntrypointList(value: Array<string>): TaskLaunchInfo;
  clearEntrypointList(): TaskLaunchInfo;
  addEntrypoint(value: string, index?: number): TaskLaunchInfo;

  getArgumentsList(): Array<string>;
  setArgumentsList(value: Array<string>): TaskLaunchInfo;
  clearArgumentsList(): TaskLaunchInfo;
  addArguments(value: string, index?: number): TaskLaunchInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TaskLaunchInfo.AsObject;
  static toObject(includeInstance: boolean, msg: TaskLaunchInfo): TaskLaunchInfo.AsObject;
  static serializeBinaryToWriter(message: TaskLaunchInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TaskLaunchInfo;
  static deserializeBinaryFromReader(message: TaskLaunchInfo, reader: jspb.BinaryReader): TaskLaunchInfo;
}

export namespace TaskLaunchInfo {
  export type AsObject = {
    ociUrl: string,
    environmentVariablesMap: Array<[string, string]>,
    entrypointList: Array<string>,
    argumentsList: Array<string>,
  }
}

export class Task extends jspb.Message {
  getId(): string;
  setId(value: string): Task;

  getTaskJob(): Ref.Job | undefined;
  setTaskJob(value?: Ref.Job): Task;
  hasTaskJob(): boolean;
  clearTaskJob(): Task;

  getWatchJob(): Ref.Job | undefined;
  setWatchJob(value?: Ref.Job): Task;
  hasWatchJob(): boolean;
  clearWatchJob(): Task;

  getStartJob(): Ref.Job | undefined;
  setStartJob(value?: Ref.Job): Task;
  hasStartJob(): boolean;
  clearStartJob(): Task;

  getStopJob(): Ref.Job | undefined;
  setStopJob(value?: Ref.Job): Task;
  hasStopJob(): boolean;
  clearStopJob(): Task;

  getStateJson(): string;
  setStateJson(value: string): Task;

  getState(): any_pb.Any | undefined;
  setState(value?: any_pb.Any): Task;
  hasState(): boolean;
  clearState(): Task;

  getJobState(): Task.State;
  setJobState(value: Task.State): Task;

  getResourceName(): string;
  setResourceName(value: string): Task;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Task.AsObject;
  static toObject(includeInstance: boolean, msg: Task): Task.AsObject;
  static serializeBinaryToWriter(message: Task, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Task;
  static deserializeBinaryFromReader(message: Task, reader: jspb.BinaryReader): Task;
}

export namespace Task {
  export type AsObject = {
    id: string,
    taskJob?: Ref.Job.AsObject,
    watchJob?: Ref.Job.AsObject,
    startJob?: Ref.Job.AsObject,
    stopJob?: Ref.Job.AsObject,
    stateJson: string,
    state?: any_pb.Any.AsObject,
    jobState: Task.State,
    resourceName: string,
  }

  export enum State { 
    UNKNOWN = 0,
    PENDING = 1,
    STARTING = 2,
    STARTED = 3,
    RUNNING = 4,
    COMPLETED = 5,
    STOPPING = 6,
    STOPPED = 7,
    MIXED = 8,
  }
}

export class UpsertTaskRequest extends jspb.Message {
  getTask(): Task | undefined;
  setTask(value?: Task): UpsertTaskRequest;
  hasTask(): boolean;
  clearTask(): UpsertTaskRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertTaskRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertTaskRequest): UpsertTaskRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertTaskRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertTaskRequest;
  static deserializeBinaryFromReader(message: UpsertTaskRequest, reader: jspb.BinaryReader): UpsertTaskRequest;
}

export namespace UpsertTaskRequest {
  export type AsObject = {
    task?: Task.AsObject,
  }
}

export class UpsertTaskResponse extends jspb.Message {
  getTask(): Task | undefined;
  setTask(value?: Task): UpsertTaskResponse;
  hasTask(): boolean;
  clearTask(): UpsertTaskResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertTaskResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertTaskResponse): UpsertTaskResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertTaskResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertTaskResponse;
  static deserializeBinaryFromReader(message: UpsertTaskResponse, reader: jspb.BinaryReader): UpsertTaskResponse;
}

export namespace UpsertTaskResponse {
  export type AsObject = {
    task?: Task.AsObject,
  }
}

export class GetTaskRequest extends jspb.Message {
  getRef(): Ref.Task | undefined;
  setRef(value?: Ref.Task): GetTaskRequest;
  hasRef(): boolean;
  clearRef(): GetTaskRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTaskRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTaskRequest): GetTaskRequest.AsObject;
  static serializeBinaryToWriter(message: GetTaskRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTaskRequest;
  static deserializeBinaryFromReader(message: GetTaskRequest, reader: jspb.BinaryReader): GetTaskRequest;
}

export namespace GetTaskRequest {
  export type AsObject = {
    ref?: Ref.Task.AsObject,
  }
}

export class GetTaskResponse extends jspb.Message {
  getTask(): Task | undefined;
  setTask(value?: Task): GetTaskResponse;
  hasTask(): boolean;
  clearTask(): GetTaskResponse;

  getTaskJob(): Job | undefined;
  setTaskJob(value?: Job): GetTaskResponse;
  hasTaskJob(): boolean;
  clearTaskJob(): GetTaskResponse;

  getStartJob(): Job | undefined;
  setStartJob(value?: Job): GetTaskResponse;
  hasStartJob(): boolean;
  clearStartJob(): GetTaskResponse;

  getStopJob(): Job | undefined;
  setStopJob(value?: Job): GetTaskResponse;
  hasStopJob(): boolean;
  clearStopJob(): GetTaskResponse;

  getWatchJob(): Job | undefined;
  setWatchJob(value?: Job): GetTaskResponse;
  hasWatchJob(): boolean;
  clearWatchJob(): GetTaskResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTaskResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetTaskResponse): GetTaskResponse.AsObject;
  static serializeBinaryToWriter(message: GetTaskResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTaskResponse;
  static deserializeBinaryFromReader(message: GetTaskResponse, reader: jspb.BinaryReader): GetTaskResponse;
}

export namespace GetTaskResponse {
  export type AsObject = {
    task?: Task.AsObject,
    taskJob?: Job.AsObject,
    startJob?: Job.AsObject,
    stopJob?: Job.AsObject,
    watchJob?: Job.AsObject,
  }
}

export class DeleteTaskRequest extends jspb.Message {
  getRef(): Ref.Task | undefined;
  setRef(value?: Ref.Task): DeleteTaskRequest;
  hasRef(): boolean;
  clearRef(): DeleteTaskRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTaskRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteTaskRequest): DeleteTaskRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteTaskRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTaskRequest;
  static deserializeBinaryFromReader(message: DeleteTaskRequest, reader: jspb.BinaryReader): DeleteTaskRequest;
}

export namespace DeleteTaskRequest {
  export type AsObject = {
    ref?: Ref.Task.AsObject,
  }
}

export class ListTaskRequest extends jspb.Message {
  getTaskstateList(): Array<Task.State>;
  setTaskstateList(value: Array<Task.State>): ListTaskRequest;
  clearTaskstateList(): ListTaskRequest;
  addTaskstate(value: Task.State, index?: number): ListTaskRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTaskRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListTaskRequest): ListTaskRequest.AsObject;
  static serializeBinaryToWriter(message: ListTaskRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTaskRequest;
  static deserializeBinaryFromReader(message: ListTaskRequest, reader: jspb.BinaryReader): ListTaskRequest;
}

export namespace ListTaskRequest {
  export type AsObject = {
    taskstateList: Array<Task.State>,
  }
}

export class ListTaskResponse extends jspb.Message {
  getTasksList(): Array<GetTaskResponse>;
  setTasksList(value: Array<GetTaskResponse>): ListTaskResponse;
  clearTasksList(): ListTaskResponse;
  addTasks(value?: GetTaskResponse, index?: number): GetTaskResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTaskResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTaskResponse): ListTaskResponse.AsObject;
  static serializeBinaryToWriter(message: ListTaskResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTaskResponse;
  static deserializeBinaryFromReader(message: ListTaskResponse, reader: jspb.BinaryReader): ListTaskResponse;
}

export namespace ListTaskResponse {
  export type AsObject = {
    tasksList: Array<GetTaskResponse.AsObject>,
  }
}

export class CancelTaskRequest extends jspb.Message {
  getRef(): Ref.Task | undefined;
  setRef(value?: Ref.Task): CancelTaskRequest;
  hasRef(): boolean;
  clearRef(): CancelTaskRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelTaskRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CancelTaskRequest): CancelTaskRequest.AsObject;
  static serializeBinaryToWriter(message: CancelTaskRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelTaskRequest;
  static deserializeBinaryFromReader(message: CancelTaskRequest, reader: jspb.BinaryReader): CancelTaskRequest;
}

export namespace CancelTaskRequest {
  export type AsObject = {
    ref?: Ref.Task.AsObject,
  }
}

export class ErrorDetail extends jspb.Message {
  getKey(): string;
  setKey(value: string): ErrorDetail;

  getValue(): string;
  setValue(value: string): ErrorDetail;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ErrorDetail.AsObject;
  static toObject(includeInstance: boolean, msg: ErrorDetail): ErrorDetail.AsObject;
  static serializeBinaryToWriter(message: ErrorDetail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ErrorDetail;
  static deserializeBinaryFromReader(message: ErrorDetail, reader: jspb.BinaryReader): ErrorDetail;
}

export namespace ErrorDetail {
  export type AsObject = {
    key: string,
    value: string,
  }
}

export class GetUserRequest extends jspb.Message {
  getUser(): Ref.User | undefined;
  setUser(value?: Ref.User): GetUserRequest;
  hasUser(): boolean;
  clearUser(): GetUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserRequest): GetUserRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserRequest;
  static deserializeBinaryFromReader(message: GetUserRequest, reader: jspb.BinaryReader): GetUserRequest;
}

export namespace GetUserRequest {
  export type AsObject = {
    user?: Ref.User.AsObject,
  }
}

export class GetUserResponse extends jspb.Message {
  getUser(): User | undefined;
  setUser(value?: User): GetUserResponse;
  hasUser(): boolean;
  clearUser(): GetUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserResponse): GetUserResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserResponse;
  static deserializeBinaryFromReader(message: GetUserResponse, reader: jspb.BinaryReader): GetUserResponse;
}

export namespace GetUserResponse {
  export type AsObject = {
    user?: User.AsObject,
  }
}

export class ListUsersResponse extends jspb.Message {
  getUsersList(): Array<User>;
  setUsersList(value: Array<User>): ListUsersResponse;
  clearUsersList(): ListUsersResponse;
  addUsers(value?: User, index?: number): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUsersResponse): ListUsersResponse.AsObject;
  static serializeBinaryToWriter(message: ListUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUsersResponse;
  static deserializeBinaryFromReader(message: ListUsersResponse, reader: jspb.BinaryReader): ListUsersResponse;
}

export namespace ListUsersResponse {
  export type AsObject = {
    usersList: Array<User.AsObject>,
  }
}

export class UpdateUserRequest extends jspb.Message {
  getUser(): User | undefined;
  setUser(value?: User): UpdateUserRequest;
  hasUser(): boolean;
  clearUser(): UpdateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserRequest): UpdateUserRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserRequest;
  static deserializeBinaryFromReader(message: UpdateUserRequest, reader: jspb.BinaryReader): UpdateUserRequest;
}

export namespace UpdateUserRequest {
  export type AsObject = {
    user?: User.AsObject,
  }
}

export class UpdateUserResponse extends jspb.Message {
  getUser(): User | undefined;
  setUser(value?: User): UpdateUserResponse;
  hasUser(): boolean;
  clearUser(): UpdateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserResponse): UpdateUserResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserResponse;
  static deserializeBinaryFromReader(message: UpdateUserResponse, reader: jspb.BinaryReader): UpdateUserResponse;
}

export namespace UpdateUserResponse {
  export type AsObject = {
    user?: User.AsObject,
  }
}

export class DeleteUserRequest extends jspb.Message {
  getUser(): Ref.User | undefined;
  setUser(value?: Ref.User): DeleteUserRequest;
  hasUser(): boolean;
  clearUser(): DeleteUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUserRequest): DeleteUserRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUserRequest;
  static deserializeBinaryFromReader(message: DeleteUserRequest, reader: jspb.BinaryReader): DeleteUserRequest;
}

export namespace DeleteUserRequest {
  export type AsObject = {
    user?: Ref.User.AsObject,
  }
}

export class AuthMethod extends jspb.Message {
  getName(): string;
  setName(value: string): AuthMethod;

  getDisplayName(): string;
  setDisplayName(value: string): AuthMethod;

  getDescription(): string;
  setDescription(value: string): AuthMethod;

  getAccessSelector(): string;
  setAccessSelector(value: string): AuthMethod;

  getOidc(): AuthMethod.OIDC | undefined;
  setOidc(value?: AuthMethod.OIDC): AuthMethod;
  hasOidc(): boolean;
  clearOidc(): AuthMethod;

  getMethodCase(): AuthMethod.MethodCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthMethod.AsObject;
  static toObject(includeInstance: boolean, msg: AuthMethod): AuthMethod.AsObject;
  static serializeBinaryToWriter(message: AuthMethod, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthMethod;
  static deserializeBinaryFromReader(message: AuthMethod, reader: jspb.BinaryReader): AuthMethod;
}

export namespace AuthMethod {
  export type AsObject = {
    name: string,
    displayName: string,
    description: string,
    accessSelector: string,
    oidc?: AuthMethod.OIDC.AsObject,
  }

  export class OIDC extends jspb.Message {
    getClientId(): string;
    setClientId(value: string): OIDC;

    getClientSecret(): string;
    setClientSecret(value: string): OIDC;

    getScopesList(): Array<string>;
    setScopesList(value: Array<string>): OIDC;
    clearScopesList(): OIDC;
    addScopes(value: string, index?: number): OIDC;

    getAudsList(): Array<string>;
    setAudsList(value: Array<string>): OIDC;
    clearAudsList(): OIDC;
    addAuds(value: string, index?: number): OIDC;

    getAllowedRedirectUrisList(): Array<string>;
    setAllowedRedirectUrisList(value: Array<string>): OIDC;
    clearAllowedRedirectUrisList(): OIDC;
    addAllowedRedirectUris(value: string, index?: number): OIDC;

    getSigningAlgsList(): Array<string>;
    setSigningAlgsList(value: Array<string>): OIDC;
    clearSigningAlgsList(): OIDC;
    addSigningAlgs(value: string, index?: number): OIDC;

    getDiscoveryUrl(): string;
    setDiscoveryUrl(value: string): OIDC;

    getDiscoveryCaPemList(): Array<string>;
    setDiscoveryCaPemList(value: Array<string>): OIDC;
    clearDiscoveryCaPemList(): OIDC;
    addDiscoveryCaPem(value: string, index?: number): OIDC;

    getClaimMappingsMap(): jspb.Map<string, string>;
    clearClaimMappingsMap(): OIDC;

    getListClaimMappingsMap(): jspb.Map<string, string>;
    clearListClaimMappingsMap(): OIDC;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OIDC.AsObject;
    static toObject(includeInstance: boolean, msg: OIDC): OIDC.AsObject;
    static serializeBinaryToWriter(message: OIDC, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OIDC;
    static deserializeBinaryFromReader(message: OIDC, reader: jspb.BinaryReader): OIDC;
  }

  export namespace OIDC {
    export type AsObject = {
      clientId: string,
      clientSecret: string,
      scopesList: Array<string>,
      audsList: Array<string>,
      allowedRedirectUrisList: Array<string>,
      signingAlgsList: Array<string>,
      discoveryUrl: string,
      discoveryCaPemList: Array<string>,
      claimMappingsMap: Array<[string, string]>,
      listClaimMappingsMap: Array<[string, string]>,
    }
  }


  export enum MethodCase { 
    METHOD_NOT_SET = 0,
    OIDC = 20,
  }
}

export class OIDCAuthMethod extends jspb.Message {
  getName(): string;
  setName(value: string): OIDCAuthMethod;

  getDisplayName(): string;
  setDisplayName(value: string): OIDCAuthMethod;

  getKind(): OIDCAuthMethod.Kind;
  setKind(value: OIDCAuthMethod.Kind): OIDCAuthMethod;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCAuthMethod.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCAuthMethod): OIDCAuthMethod.AsObject;
  static serializeBinaryToWriter(message: OIDCAuthMethod, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCAuthMethod;
  static deserializeBinaryFromReader(message: OIDCAuthMethod, reader: jspb.BinaryReader): OIDCAuthMethod;
}

export namespace OIDCAuthMethod {
  export type AsObject = {
    name: string,
    displayName: string,
    kind: OIDCAuthMethod.Kind,
  }

  export enum Kind { 
    UNKNOWN = 0,
    GITHUB = 1,
    GOOGLE = 2,
  }
}

export class UpsertAuthMethodRequest extends jspb.Message {
  getAuthMethod(): AuthMethod | undefined;
  setAuthMethod(value?: AuthMethod): UpsertAuthMethodRequest;
  hasAuthMethod(): boolean;
  clearAuthMethod(): UpsertAuthMethodRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertAuthMethodRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertAuthMethodRequest): UpsertAuthMethodRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertAuthMethodRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertAuthMethodRequest;
  static deserializeBinaryFromReader(message: UpsertAuthMethodRequest, reader: jspb.BinaryReader): UpsertAuthMethodRequest;
}

export namespace UpsertAuthMethodRequest {
  export type AsObject = {
    authMethod?: AuthMethod.AsObject,
  }
}

export class UpsertAuthMethodResponse extends jspb.Message {
  getAuthMethod(): AuthMethod | undefined;
  setAuthMethod(value?: AuthMethod): UpsertAuthMethodResponse;
  hasAuthMethod(): boolean;
  clearAuthMethod(): UpsertAuthMethodResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertAuthMethodResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertAuthMethodResponse): UpsertAuthMethodResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertAuthMethodResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertAuthMethodResponse;
  static deserializeBinaryFromReader(message: UpsertAuthMethodResponse, reader: jspb.BinaryReader): UpsertAuthMethodResponse;
}

export namespace UpsertAuthMethodResponse {
  export type AsObject = {
    authMethod?: AuthMethod.AsObject,
  }
}

export class GetAuthMethodRequest extends jspb.Message {
  getAuthMethod(): Ref.AuthMethod | undefined;
  setAuthMethod(value?: Ref.AuthMethod): GetAuthMethodRequest;
  hasAuthMethod(): boolean;
  clearAuthMethod(): GetAuthMethodRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthMethodRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthMethodRequest): GetAuthMethodRequest.AsObject;
  static serializeBinaryToWriter(message: GetAuthMethodRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthMethodRequest;
  static deserializeBinaryFromReader(message: GetAuthMethodRequest, reader: jspb.BinaryReader): GetAuthMethodRequest;
}

export namespace GetAuthMethodRequest {
  export type AsObject = {
    authMethod?: Ref.AuthMethod.AsObject,
  }
}

export class GetAuthMethodResponse extends jspb.Message {
  getAuthMethod(): AuthMethod | undefined;
  setAuthMethod(value?: AuthMethod): GetAuthMethodResponse;
  hasAuthMethod(): boolean;
  clearAuthMethod(): GetAuthMethodResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthMethodResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthMethodResponse): GetAuthMethodResponse.AsObject;
  static serializeBinaryToWriter(message: GetAuthMethodResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthMethodResponse;
  static deserializeBinaryFromReader(message: GetAuthMethodResponse, reader: jspb.BinaryReader): GetAuthMethodResponse;
}

export namespace GetAuthMethodResponse {
  export type AsObject = {
    authMethod?: AuthMethod.AsObject,
  }
}

export class DeleteAuthMethodRequest extends jspb.Message {
  getAuthMethod(): Ref.AuthMethod | undefined;
  setAuthMethod(value?: Ref.AuthMethod): DeleteAuthMethodRequest;
  hasAuthMethod(): boolean;
  clearAuthMethod(): DeleteAuthMethodRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAuthMethodRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAuthMethodRequest): DeleteAuthMethodRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteAuthMethodRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAuthMethodRequest;
  static deserializeBinaryFromReader(message: DeleteAuthMethodRequest, reader: jspb.BinaryReader): DeleteAuthMethodRequest;
}

export namespace DeleteAuthMethodRequest {
  export type AsObject = {
    authMethod?: Ref.AuthMethod.AsObject,
  }
}

export class ListAuthMethodsResponse extends jspb.Message {
  getAuthMethodsList(): Array<AuthMethod>;
  setAuthMethodsList(value: Array<AuthMethod>): ListAuthMethodsResponse;
  clearAuthMethodsList(): ListAuthMethodsResponse;
  addAuthMethods(value?: AuthMethod, index?: number): AuthMethod;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAuthMethodsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAuthMethodsResponse): ListAuthMethodsResponse.AsObject;
  static serializeBinaryToWriter(message: ListAuthMethodsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAuthMethodsResponse;
  static deserializeBinaryFromReader(message: ListAuthMethodsResponse, reader: jspb.BinaryReader): ListAuthMethodsResponse;
}

export namespace ListAuthMethodsResponse {
  export type AsObject = {
    authMethodsList: Array<AuthMethod.AsObject>,
  }
}

export class ListOIDCAuthMethodsResponse extends jspb.Message {
  getAuthMethodsList(): Array<OIDCAuthMethod>;
  setAuthMethodsList(value: Array<OIDCAuthMethod>): ListOIDCAuthMethodsResponse;
  clearAuthMethodsList(): ListOIDCAuthMethodsResponse;
  addAuthMethods(value?: OIDCAuthMethod, index?: number): OIDCAuthMethod;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOIDCAuthMethodsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOIDCAuthMethodsResponse): ListOIDCAuthMethodsResponse.AsObject;
  static serializeBinaryToWriter(message: ListOIDCAuthMethodsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOIDCAuthMethodsResponse;
  static deserializeBinaryFromReader(message: ListOIDCAuthMethodsResponse, reader: jspb.BinaryReader): ListOIDCAuthMethodsResponse;
}

export namespace ListOIDCAuthMethodsResponse {
  export type AsObject = {
    authMethodsList: Array<OIDCAuthMethod.AsObject>,
  }
}

export class GetOIDCAuthURLRequest extends jspb.Message {
  getAuthMethod(): Ref.AuthMethod | undefined;
  setAuthMethod(value?: Ref.AuthMethod): GetOIDCAuthURLRequest;
  hasAuthMethod(): boolean;
  clearAuthMethod(): GetOIDCAuthURLRequest;

  getRedirectUri(): string;
  setRedirectUri(value: string): GetOIDCAuthURLRequest;

  getNonce(): string;
  setNonce(value: string): GetOIDCAuthURLRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOIDCAuthURLRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOIDCAuthURLRequest): GetOIDCAuthURLRequest.AsObject;
  static serializeBinaryToWriter(message: GetOIDCAuthURLRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOIDCAuthURLRequest;
  static deserializeBinaryFromReader(message: GetOIDCAuthURLRequest, reader: jspb.BinaryReader): GetOIDCAuthURLRequest;
}

export namespace GetOIDCAuthURLRequest {
  export type AsObject = {
    authMethod?: Ref.AuthMethod.AsObject,
    redirectUri: string,
    nonce: string,
  }
}

export class GetOIDCAuthURLResponse extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): GetOIDCAuthURLResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOIDCAuthURLResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOIDCAuthURLResponse): GetOIDCAuthURLResponse.AsObject;
  static serializeBinaryToWriter(message: GetOIDCAuthURLResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOIDCAuthURLResponse;
  static deserializeBinaryFromReader(message: GetOIDCAuthURLResponse, reader: jspb.BinaryReader): GetOIDCAuthURLResponse;
}

export namespace GetOIDCAuthURLResponse {
  export type AsObject = {
    url: string,
  }
}

export class CompleteOIDCAuthRequest extends jspb.Message {
  getAuthMethod(): Ref.AuthMethod | undefined;
  setAuthMethod(value?: Ref.AuthMethod): CompleteOIDCAuthRequest;
  hasAuthMethod(): boolean;
  clearAuthMethod(): CompleteOIDCAuthRequest;

  getRedirectUri(): string;
  setRedirectUri(value: string): CompleteOIDCAuthRequest;

  getState(): string;
  setState(value: string): CompleteOIDCAuthRequest;

  getNonce(): string;
  setNonce(value: string): CompleteOIDCAuthRequest;

  getCode(): string;
  setCode(value: string): CompleteOIDCAuthRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteOIDCAuthRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteOIDCAuthRequest): CompleteOIDCAuthRequest.AsObject;
  static serializeBinaryToWriter(message: CompleteOIDCAuthRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteOIDCAuthRequest;
  static deserializeBinaryFromReader(message: CompleteOIDCAuthRequest, reader: jspb.BinaryReader): CompleteOIDCAuthRequest;
}

export namespace CompleteOIDCAuthRequest {
  export type AsObject = {
    authMethod?: Ref.AuthMethod.AsObject,
    redirectUri: string,
    state: string,
    nonce: string,
    code: string,
  }
}

export class CompleteOIDCAuthResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): CompleteOIDCAuthResponse;

  getUser(): User | undefined;
  setUser(value?: User): CompleteOIDCAuthResponse;
  hasUser(): boolean;
  clearUser(): CompleteOIDCAuthResponse;

  getIdClaimsJson(): string;
  setIdClaimsJson(value: string): CompleteOIDCAuthResponse;

  getUserClaimsJson(): string;
  setUserClaimsJson(value: string): CompleteOIDCAuthResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CompleteOIDCAuthResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CompleteOIDCAuthResponse): CompleteOIDCAuthResponse.AsObject;
  static serializeBinaryToWriter(message: CompleteOIDCAuthResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CompleteOIDCAuthResponse;
  static deserializeBinaryFromReader(message: CompleteOIDCAuthResponse, reader: jspb.BinaryReader): CompleteOIDCAuthResponse;
}

export namespace CompleteOIDCAuthResponse {
  export type AsObject = {
    token: string,
    user?: User.AsObject,
    idClaimsJson: string,
    userClaimsJson: string,
  }
}

export class QueueJobRequest extends jspb.Message {
  getJob(): Job | undefined;
  setJob(value?: Job): QueueJobRequest;
  hasJob(): boolean;
  clearJob(): QueueJobRequest;

  getExpiresIn(): string;
  setExpiresIn(value: string): QueueJobRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueueJobRequest.AsObject;
  static toObject(includeInstance: boolean, msg: QueueJobRequest): QueueJobRequest.AsObject;
  static serializeBinaryToWriter(message: QueueJobRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueueJobRequest;
  static deserializeBinaryFromReader(message: QueueJobRequest, reader: jspb.BinaryReader): QueueJobRequest;
}

export namespace QueueJobRequest {
  export type AsObject = {
    job?: Job.AsObject,
    expiresIn: string,
  }
}

export class QueueJobResponse extends jspb.Message {
  getJobId(): string;
  setJobId(value: string): QueueJobResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueueJobResponse.AsObject;
  static toObject(includeInstance: boolean, msg: QueueJobResponse): QueueJobResponse.AsObject;
  static serializeBinaryToWriter(message: QueueJobResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueueJobResponse;
  static deserializeBinaryFromReader(message: QueueJobResponse, reader: jspb.BinaryReader): QueueJobResponse;
}

export namespace QueueJobResponse {
  export type AsObject = {
    jobId: string,
  }
}

export class CancelJobRequest extends jspb.Message {
  getJobId(): string;
  setJobId(value: string): CancelJobRequest;

  getForce(): boolean;
  setForce(value: boolean): CancelJobRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelJobRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CancelJobRequest): CancelJobRequest.AsObject;
  static serializeBinaryToWriter(message: CancelJobRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelJobRequest;
  static deserializeBinaryFromReader(message: CancelJobRequest, reader: jspb.BinaryReader): CancelJobRequest;
}

export namespace CancelJobRequest {
  export type AsObject = {
    jobId: string,
    force: boolean,
  }
}

export class ValidateJobRequest extends jspb.Message {
  getJob(): Job | undefined;
  setJob(value?: Job): ValidateJobRequest;
  hasJob(): boolean;
  clearJob(): ValidateJobRequest;

  getDisableAssign(): boolean;
  setDisableAssign(value: boolean): ValidateJobRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateJobRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateJobRequest): ValidateJobRequest.AsObject;
  static serializeBinaryToWriter(message: ValidateJobRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateJobRequest;
  static deserializeBinaryFromReader(message: ValidateJobRequest, reader: jspb.BinaryReader): ValidateJobRequest;
}

export namespace ValidateJobRequest {
  export type AsObject = {
    job?: Job.AsObject,
    disableAssign: boolean,
  }
}

export class ValidateJobResponse extends jspb.Message {
  getValid(): boolean;
  setValid(value: boolean): ValidateJobResponse;

  getValidationError(): google_rpc_status_pb.Status | undefined;
  setValidationError(value?: google_rpc_status_pb.Status): ValidateJobResponse;
  hasValidationError(): boolean;
  clearValidationError(): ValidateJobResponse;

  getAssignable(): boolean;
  setAssignable(value: boolean): ValidateJobResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateJobResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateJobResponse): ValidateJobResponse.AsObject;
  static serializeBinaryToWriter(message: ValidateJobResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateJobResponse;
  static deserializeBinaryFromReader(message: ValidateJobResponse, reader: jspb.BinaryReader): ValidateJobResponse;
}

export namespace ValidateJobResponse {
  export type AsObject = {
    valid: boolean,
    validationError?: google_rpc_status_pb.Status.AsObject,
    assignable: boolean,
  }
}

export class Job extends jspb.Message {
  getId(): string;
  setId(value: string): Job;

  getSingletonId(): string;
  setSingletonId(value: string): Job;

  getDependsOnList(): Array<string>;
  setDependsOnList(value: Array<string>): Job;
  clearDependsOnList(): Job;
  addDependsOn(value: string, index?: number): Job;

  getDependsOnAllowFailureList(): Array<string>;
  setDependsOnAllowFailureList(value: Array<string>): Job;
  clearDependsOnAllowFailureList(): Job;
  addDependsOnAllowFailure(value: string, index?: number): Job;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): Job;
  hasApplication(): boolean;
  clearApplication(): Job;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): Job;
  hasWorkspace(): boolean;
  clearWorkspace(): Job;

  getTargetRunner(): Ref.Runner | undefined;
  setTargetRunner(value?: Ref.Runner): Job;
  hasTargetRunner(): boolean;
  clearTargetRunner(): Job;

  getOndemandRunner(): Ref.OnDemandRunnerConfig | undefined;
  setOndemandRunner(value?: Ref.OnDemandRunnerConfig): Job;
  hasOndemandRunner(): boolean;
  clearOndemandRunner(): Job;

  getOndemandRunnerTask(): Job.TaskOverride | undefined;
  setOndemandRunnerTask(value?: Job.TaskOverride): Job;
  hasOndemandRunnerTask(): boolean;
  clearOndemandRunnerTask(): Job;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Job;

  getDataSource(): Job.DataSource | undefined;
  setDataSource(value?: Job.DataSource): Job;
  hasDataSource(): boolean;
  clearDataSource(): Job;

  getDataSourceOverridesMap(): jspb.Map<string, string>;
  clearDataSourceOverridesMap(): Job;

  getWaypointHcl(): Hcl | undefined;
  setWaypointHcl(value?: Hcl): Job;
  hasWaypointHcl(): boolean;
  clearWaypointHcl(): Job;

  getVariablesList(): Array<Variable>;
  setVariablesList(value: Array<Variable>): Job;
  clearVariablesList(): Job;
  addVariables(value?: Variable, index?: number): Variable;

  getNoop(): Job.Noop | undefined;
  setNoop(value?: Job.Noop): Job;
  hasNoop(): boolean;
  clearNoop(): Job;

  getBuild(): Job.BuildOp | undefined;
  setBuild(value?: Job.BuildOp): Job;
  hasBuild(): boolean;
  clearBuild(): Job;

  getPush(): Job.PushOp | undefined;
  setPush(value?: Job.PushOp): Job;
  hasPush(): boolean;
  clearPush(): Job;

  getDeploy(): Job.DeployOp | undefined;
  setDeploy(value?: Job.DeployOp): Job;
  hasDeploy(): boolean;
  clearDeploy(): Job;

  getDestroy(): Job.DestroyOp | undefined;
  setDestroy(value?: Job.DestroyOp): Job;
  hasDestroy(): boolean;
  clearDestroy(): Job;

  getRelease(): Job.ReleaseOp | undefined;
  setRelease(value?: Job.ReleaseOp): Job;
  hasRelease(): boolean;
  clearRelease(): Job;

  getValidate(): Job.ValidateOp | undefined;
  setValidate(value?: Job.ValidateOp): Job;
  hasValidate(): boolean;
  clearValidate(): Job;

  getAuth(): Job.AuthOp | undefined;
  setAuth(value?: Job.AuthOp): Job;
  hasAuth(): boolean;
  clearAuth(): Job;

  getDocs(): Job.DocsOp | undefined;
  setDocs(value?: Job.DocsOp): Job;
  hasDocs(): boolean;
  clearDocs(): Job;

  getConfigSync(): Job.ConfigSyncOp | undefined;
  setConfigSync(value?: Job.ConfigSyncOp): Job;
  hasConfigSync(): boolean;
  clearConfigSync(): Job;

  getExec(): Job.ExecOp | undefined;
  setExec(value?: Job.ExecOp): Job;
  hasExec(): boolean;
  clearExec(): Job;

  getUp(): Job.UpOp | undefined;
  setUp(value?: Job.UpOp): Job;
  hasUp(): boolean;
  clearUp(): Job;

  getLogs(): Job.LogsOp | undefined;
  setLogs(value?: Job.LogsOp): Job;
  hasLogs(): boolean;
  clearLogs(): Job;

  getQueueProject(): Job.QueueProjectOp | undefined;
  setQueueProject(value?: Job.QueueProjectOp): Job;
  hasQueueProject(): boolean;
  clearQueueProject(): Job;

  getPoll(): Job.PollOp | undefined;
  setPoll(value?: Job.PollOp): Job;
  hasPoll(): boolean;
  clearPoll(): Job;

  getStatusReport(): Job.StatusReportOp | undefined;
  setStatusReport(value?: Job.StatusReportOp): Job;
  hasStatusReport(): boolean;
  clearStatusReport(): Job;

  getStartTask(): Job.StartTaskLaunchOp | undefined;
  setStartTask(value?: Job.StartTaskLaunchOp): Job;
  hasStartTask(): boolean;
  clearStartTask(): Job;

  getStopTask(): Job.StopTaskLaunchOp | undefined;
  setStopTask(value?: Job.StopTaskLaunchOp): Job;
  hasStopTask(): boolean;
  clearStopTask(): Job;

  getInit(): Job.InitOp | undefined;
  setInit(value?: Job.InitOp): Job;
  hasInit(): boolean;
  clearInit(): Job;

  getWatchTask(): Job.WatchTaskOp | undefined;
  setWatchTask(value?: Job.WatchTaskOp): Job;
  hasWatchTask(): boolean;
  clearWatchTask(): Job;

  getPipelineStep(): Job.PipelineStepOp | undefined;
  setPipelineStep(value?: Job.PipelineStepOp): Job;
  hasPipelineStep(): boolean;
  clearPipelineStep(): Job;

  getState(): Job.State;
  setState(value: Job.State): Job;

  getAssignedRunner(): Ref.RunnerId | undefined;
  setAssignedRunner(value?: Ref.RunnerId): Job;
  hasAssignedRunner(): boolean;
  clearAssignedRunner(): Job;

  getQueueTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setQueueTime(value?: google_protobuf_timestamp_pb.Timestamp): Job;
  hasQueueTime(): boolean;
  clearQueueTime(): Job;

  getAssignTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAssignTime(value?: google_protobuf_timestamp_pb.Timestamp): Job;
  hasAssignTime(): boolean;
  clearAssignTime(): Job;

  getAckTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAckTime(value?: google_protobuf_timestamp_pb.Timestamp): Job;
  hasAckTime(): boolean;
  clearAckTime(): Job;

  getCompleteTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCompleteTime(value?: google_protobuf_timestamp_pb.Timestamp): Job;
  hasCompleteTime(): boolean;
  clearCompleteTime(): Job;

  getDataSourceRef(): Job.DataSource.Ref | undefined;
  setDataSourceRef(value?: Job.DataSource.Ref): Job;
  hasDataSourceRef(): boolean;
  clearDataSourceRef(): Job;

  getVariableFinalValuesMap(): jspb.Map<string, Variable.FinalValue>;
  clearVariableFinalValuesMap(): Job;

  getConfig(): Job.Config | undefined;
  setConfig(value?: Job.Config): Job;
  hasConfig(): boolean;
  clearConfig(): Job;

  getError(): google_rpc_status_pb.Status | undefined;
  setError(value?: google_rpc_status_pb.Status): Job;
  hasError(): boolean;
  clearError(): Job;

  getResult(): Job.Result | undefined;
  setResult(value?: Job.Result): Job;
  hasResult(): boolean;
  clearResult(): Job;

  getCancelTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCancelTime(value?: google_protobuf_timestamp_pb.Timestamp): Job;
  hasCancelTime(): boolean;
  clearCancelTime(): Job;

  getExpireTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpireTime(value?: google_protobuf_timestamp_pb.Timestamp): Job;
  hasExpireTime(): boolean;
  clearExpireTime(): Job;

  getTask(): Ref.Task | undefined;
  setTask(value?: Ref.Task): Job;
  hasTask(): boolean;
  clearTask(): Job;

  getOperationCase(): Job.OperationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Job.AsObject;
  static toObject(includeInstance: boolean, msg: Job): Job.AsObject;
  static serializeBinaryToWriter(message: Job, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Job;
  static deserializeBinaryFromReader(message: Job, reader: jspb.BinaryReader): Job;
}

export namespace Job {
  export type AsObject = {
    id: string,
    singletonId: string,
    dependsOnList: Array<string>,
    dependsOnAllowFailureList: Array<string>,
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    targetRunner?: Ref.Runner.AsObject,
    ondemandRunner?: Ref.OnDemandRunnerConfig.AsObject,
    ondemandRunnerTask?: Job.TaskOverride.AsObject,
    labelsMap: Array<[string, string]>,
    dataSource?: Job.DataSource.AsObject,
    dataSourceOverridesMap: Array<[string, string]>,
    waypointHcl?: Hcl.AsObject,
    variablesList: Array<Variable.AsObject>,
    noop?: Job.Noop.AsObject,
    build?: Job.BuildOp.AsObject,
    push?: Job.PushOp.AsObject,
    deploy?: Job.DeployOp.AsObject,
    destroy?: Job.DestroyOp.AsObject,
    release?: Job.ReleaseOp.AsObject,
    validate?: Job.ValidateOp.AsObject,
    auth?: Job.AuthOp.AsObject,
    docs?: Job.DocsOp.AsObject,
    configSync?: Job.ConfigSyncOp.AsObject,
    exec?: Job.ExecOp.AsObject,
    up?: Job.UpOp.AsObject,
    logs?: Job.LogsOp.AsObject,
    queueProject?: Job.QueueProjectOp.AsObject,
    poll?: Job.PollOp.AsObject,
    statusReport?: Job.StatusReportOp.AsObject,
    startTask?: Job.StartTaskLaunchOp.AsObject,
    stopTask?: Job.StopTaskLaunchOp.AsObject,
    init?: Job.InitOp.AsObject,
    watchTask?: Job.WatchTaskOp.AsObject,
    pipelineStep?: Job.PipelineStepOp.AsObject,
    state: Job.State,
    assignedRunner?: Ref.RunnerId.AsObject,
    queueTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    assignTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    ackTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    completeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    dataSourceRef?: Job.DataSource.Ref.AsObject,
    variableFinalValuesMap: Array<[string, Variable.FinalValue.AsObject]>,
    config?: Job.Config.AsObject,
    error?: google_rpc_status_pb.Status.AsObject,
    result?: Job.Result.AsObject,
    cancelTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    expireTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    task?: Ref.Task.AsObject,
  }

  export class TaskOverride extends jspb.Message {
    getLaunchInfo(): TaskLaunchInfo | undefined;
    setLaunchInfo(value?: TaskLaunchInfo): TaskOverride;
    hasLaunchInfo(): boolean;
    clearLaunchInfo(): TaskOverride;

    getSkipOperation(): boolean;
    setSkipOperation(value: boolean): TaskOverride;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TaskOverride.AsObject;
    static toObject(includeInstance: boolean, msg: TaskOverride): TaskOverride.AsObject;
    static serializeBinaryToWriter(message: TaskOverride, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TaskOverride;
    static deserializeBinaryFromReader(message: TaskOverride, reader: jspb.BinaryReader): TaskOverride;
  }

  export namespace TaskOverride {
    export type AsObject = {
      launchInfo?: TaskLaunchInfo.AsObject,
      skipOperation: boolean,
    }
  }


  export class Result extends jspb.Message {
    getBuild(): Job.BuildResult | undefined;
    setBuild(value?: Job.BuildResult): Result;
    hasBuild(): boolean;
    clearBuild(): Result;

    getPush(): Job.PushResult | undefined;
    setPush(value?: Job.PushResult): Result;
    hasPush(): boolean;
    clearPush(): Result;

    getDeploy(): Job.DeployResult | undefined;
    setDeploy(value?: Job.DeployResult): Result;
    hasDeploy(): boolean;
    clearDeploy(): Result;

    getRelease(): Job.ReleaseResult | undefined;
    setRelease(value?: Job.ReleaseResult): Result;
    hasRelease(): boolean;
    clearRelease(): Result;

    getValidate(): Job.ValidateResult | undefined;
    setValidate(value?: Job.ValidateResult): Result;
    hasValidate(): boolean;
    clearValidate(): Result;

    getAuth(): Job.AuthResult | undefined;
    setAuth(value?: Job.AuthResult): Result;
    hasAuth(): boolean;
    clearAuth(): Result;

    getDocs(): Job.DocsResult | undefined;
    setDocs(value?: Job.DocsResult): Result;
    hasDocs(): boolean;
    clearDocs(): Result;

    getConfigSync(): Job.ConfigSyncResult | undefined;
    setConfigSync(value?: Job.ConfigSyncResult): Result;
    hasConfigSync(): boolean;
    clearConfigSync(): Result;

    getUp(): Job.UpResult | undefined;
    setUp(value?: Job.UpResult): Result;
    hasUp(): boolean;
    clearUp(): Result;

    getQueueProject(): Job.QueueProjectResult | undefined;
    setQueueProject(value?: Job.QueueProjectResult): Result;
    hasQueueProject(): boolean;
    clearQueueProject(): Result;

    getPoll(): Job.PollResult | undefined;
    setPoll(value?: Job.PollResult): Result;
    hasPoll(): boolean;
    clearPoll(): Result;

    getStatusReport(): Job.StatusReportResult | undefined;
    setStatusReport(value?: Job.StatusReportResult): Result;
    hasStatusReport(): boolean;
    clearStatusReport(): Result;

    getStartTask(): Job.StartTaskResult | undefined;
    setStartTask(value?: Job.StartTaskResult): Result;
    hasStartTask(): boolean;
    clearStartTask(): Result;

    getInit(): Job.InitResult | undefined;
    setInit(value?: Job.InitResult): Result;
    hasInit(): boolean;
    clearInit(): Result;

    getWatchTask(): Job.WatchTaskResult | undefined;
    setWatchTask(value?: Job.WatchTaskResult): Result;
    hasWatchTask(): boolean;
    clearWatchTask(): Result;

    getPipelineStep(): Job.PipelineStepResult | undefined;
    setPipelineStep(value?: Job.PipelineStepResult): Result;
    hasPipelineStep(): boolean;
    clearPipelineStep(): Result;

    getPipelineConfigSync(): Job.PipelineConfigSyncResult | undefined;
    setPipelineConfigSync(value?: Job.PipelineConfigSyncResult): Result;
    hasPipelineConfigSync(): boolean;
    clearPipelineConfigSync(): Result;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Result.AsObject;
    static toObject(includeInstance: boolean, msg: Result): Result.AsObject;
    static serializeBinaryToWriter(message: Result, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Result;
    static deserializeBinaryFromReader(message: Result, reader: jspb.BinaryReader): Result;
  }

  export namespace Result {
    export type AsObject = {
      build?: Job.BuildResult.AsObject,
      push?: Job.PushResult.AsObject,
      deploy?: Job.DeployResult.AsObject,
      release?: Job.ReleaseResult.AsObject,
      validate?: Job.ValidateResult.AsObject,
      auth?: Job.AuthResult.AsObject,
      docs?: Job.DocsResult.AsObject,
      configSync?: Job.ConfigSyncResult.AsObject,
      up?: Job.UpResult.AsObject,
      queueProject?: Job.QueueProjectResult.AsObject,
      poll?: Job.PollResult.AsObject,
      statusReport?: Job.StatusReportResult.AsObject,
      startTask?: Job.StartTaskResult.AsObject,
      init?: Job.InitResult.AsObject,
      watchTask?: Job.WatchTaskResult.AsObject,
      pipelineStep?: Job.PipelineStepResult.AsObject,
      pipelineConfigSync?: Job.PipelineConfigSyncResult.AsObject,
    }
  }


  export class Config extends jspb.Message {
    getSource(): Job.Config.Source;
    setSource(value: Job.Config.Source): Config;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Config.AsObject;
    static toObject(includeInstance: boolean, msg: Config): Config.AsObject;
    static serializeBinaryToWriter(message: Config, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Config;
    static deserializeBinaryFromReader(message: Config, reader: jspb.BinaryReader): Config;
  }

  export namespace Config {
    export type AsObject = {
      source: Job.Config.Source,
    }

    export enum Source { 
      UNKNOWN = 0,
      FILE = 1,
      SERVER = 2,
      JOB = 3,
    }
  }


  export class DataSource extends jspb.Message {
    getLocal(): Job.Local | undefined;
    setLocal(value?: Job.Local): DataSource;
    hasLocal(): boolean;
    clearLocal(): DataSource;

    getGit(): Job.Git | undefined;
    setGit(value?: Job.Git): DataSource;
    hasGit(): boolean;
    clearGit(): DataSource;

    getSourceCase(): DataSource.SourceCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DataSource.AsObject;
    static toObject(includeInstance: boolean, msg: DataSource): DataSource.AsObject;
    static serializeBinaryToWriter(message: DataSource, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DataSource;
    static deserializeBinaryFromReader(message: DataSource, reader: jspb.BinaryReader): DataSource;
  }

  export namespace DataSource {
    export type AsObject = {
      local?: Job.Local.AsObject,
      git?: Job.Git.AsObject,
    }

    export class Ref extends jspb.Message {
      getUnknown(): google_protobuf_empty_pb.Empty | undefined;
      setUnknown(value?: google_protobuf_empty_pb.Empty): Ref;
      hasUnknown(): boolean;
      clearUnknown(): Ref;

      getGit(): Job.Git.Ref | undefined;
      setGit(value?: Job.Git.Ref): Ref;
      hasGit(): boolean;
      clearGit(): Ref;

      getRefCase(): Ref.RefCase;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Ref.AsObject;
      static toObject(includeInstance: boolean, msg: Ref): Ref.AsObject;
      static serializeBinaryToWriter(message: Ref, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Ref;
      static deserializeBinaryFromReader(message: Ref, reader: jspb.BinaryReader): Ref;
    }

    export namespace Ref {
      export type AsObject = {
        unknown?: google_protobuf_empty_pb.Empty.AsObject,
        git?: Job.Git.Ref.AsObject,
      }

      export enum RefCase { 
        REF_NOT_SET = 0,
        UNKNOWN = 1,
        GIT = 2,
      }
    }


    export enum SourceCase { 
      SOURCE_NOT_SET = 0,
      LOCAL = 1,
      GIT = 2,
    }
  }


  export class Local extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Local.AsObject;
    static toObject(includeInstance: boolean, msg: Local): Local.AsObject;
    static serializeBinaryToWriter(message: Local, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Local;
    static deserializeBinaryFromReader(message: Local, reader: jspb.BinaryReader): Local;
  }

  export namespace Local {
    export type AsObject = {
    }
  }


  export class Git extends jspb.Message {
    getUrl(): string;
    setUrl(value: string): Git;

    getRef(): string;
    setRef(value: string): Git;

    getPath(): string;
    setPath(value: string): Git;

    getIgnoreChangesOutsidePath(): boolean;
    setIgnoreChangesOutsidePath(value: boolean): Git;

    getRecurseSubmodules(): number;
    setRecurseSubmodules(value: number): Git;

    getBasic(): Job.Git.Basic | undefined;
    setBasic(value?: Job.Git.Basic): Git;
    hasBasic(): boolean;
    clearBasic(): Git;

    getSsh(): Job.Git.SSH | undefined;
    setSsh(value?: Job.Git.SSH): Git;
    hasSsh(): boolean;
    clearSsh(): Git;

    getAuthCase(): Git.AuthCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Git.AsObject;
    static toObject(includeInstance: boolean, msg: Git): Git.AsObject;
    static serializeBinaryToWriter(message: Git, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Git;
    static deserializeBinaryFromReader(message: Git, reader: jspb.BinaryReader): Git;
  }

  export namespace Git {
    export type AsObject = {
      url: string,
      ref: string,
      path: string,
      ignoreChangesOutsidePath: boolean,
      recurseSubmodules: number,
      basic?: Job.Git.Basic.AsObject,
      ssh?: Job.Git.SSH.AsObject,
    }

    export class Basic extends jspb.Message {
      getUsername(): string;
      setUsername(value: string): Basic;

      getPassword(): string;
      setPassword(value: string): Basic;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Basic.AsObject;
      static toObject(includeInstance: boolean, msg: Basic): Basic.AsObject;
      static serializeBinaryToWriter(message: Basic, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Basic;
      static deserializeBinaryFromReader(message: Basic, reader: jspb.BinaryReader): Basic;
    }

    export namespace Basic {
      export type AsObject = {
        username: string,
        password: string,
      }
    }


    export class SSH extends jspb.Message {
      getPrivateKeyPem(): Uint8Array | string;
      getPrivateKeyPem_asU8(): Uint8Array;
      getPrivateKeyPem_asB64(): string;
      setPrivateKeyPem(value: Uint8Array | string): SSH;

      getPassword(): string;
      setPassword(value: string): SSH;

      getUser(): string;
      setUser(value: string): SSH;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): SSH.AsObject;
      static toObject(includeInstance: boolean, msg: SSH): SSH.AsObject;
      static serializeBinaryToWriter(message: SSH, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): SSH;
      static deserializeBinaryFromReader(message: SSH, reader: jspb.BinaryReader): SSH;
    }

    export namespace SSH {
      export type AsObject = {
        privateKeyPem: Uint8Array | string,
        password: string,
        user: string,
      }
    }


    export class Ref extends jspb.Message {
      getCommit(): string;
      setCommit(value: string): Ref;

      getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
      setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): Ref;
      hasTimestamp(): boolean;
      clearTimestamp(): Ref;

      getCommitMessage(): string;
      setCommitMessage(value: string): Ref;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Ref.AsObject;
      static toObject(includeInstance: boolean, msg: Ref): Ref.AsObject;
      static serializeBinaryToWriter(message: Ref, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Ref;
      static deserializeBinaryFromReader(message: Ref, reader: jspb.BinaryReader): Ref;
    }

    export namespace Ref {
      export type AsObject = {
        commit: string,
        timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        commitMessage: string,
      }
    }


    export enum AuthCase { 
      AUTH_NOT_SET = 0,
      BASIC = 4,
      SSH = 5,
    }
  }


  export class Noop extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Noop.AsObject;
    static toObject(includeInstance: boolean, msg: Noop): Noop.AsObject;
    static serializeBinaryToWriter(message: Noop, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Noop;
    static deserializeBinaryFromReader(message: Noop, reader: jspb.BinaryReader): Noop;
  }

  export namespace Noop {
    export type AsObject = {
    }
  }


  export class UpOp extends jspb.Message {
    getRelease(): Job.ReleaseOp | undefined;
    setRelease(value?: Job.ReleaseOp): UpOp;
    hasRelease(): boolean;
    clearRelease(): UpOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpOp.AsObject;
    static toObject(includeInstance: boolean, msg: UpOp): UpOp.AsObject;
    static serializeBinaryToWriter(message: UpOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpOp;
    static deserializeBinaryFromReader(message: UpOp, reader: jspb.BinaryReader): UpOp;
  }

  export namespace UpOp {
    export type AsObject = {
      release?: Job.ReleaseOp.AsObject,
    }
  }


  export class UpResult extends jspb.Message {
    getReleaseUrl(): string;
    setReleaseUrl(value: string): UpResult;

    getAppUrl(): string;
    setAppUrl(value: string): UpResult;

    getDeployUrl(): string;
    setDeployUrl(value: string): UpResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpResult.AsObject;
    static toObject(includeInstance: boolean, msg: UpResult): UpResult.AsObject;
    static serializeBinaryToWriter(message: UpResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpResult;
    static deserializeBinaryFromReader(message: UpResult, reader: jspb.BinaryReader): UpResult;
  }

  export namespace UpResult {
    export type AsObject = {
      releaseUrl: string,
      appUrl: string,
      deployUrl: string,
    }
  }


  export class ValidateOp extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ValidateOp.AsObject;
    static toObject(includeInstance: boolean, msg: ValidateOp): ValidateOp.AsObject;
    static serializeBinaryToWriter(message: ValidateOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ValidateOp;
    static deserializeBinaryFromReader(message: ValidateOp, reader: jspb.BinaryReader): ValidateOp;
  }

  export namespace ValidateOp {
    export type AsObject = {
    }
  }


  export class ValidateResult extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ValidateResult.AsObject;
    static toObject(includeInstance: boolean, msg: ValidateResult): ValidateResult.AsObject;
    static serializeBinaryToWriter(message: ValidateResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ValidateResult;
    static deserializeBinaryFromReader(message: ValidateResult, reader: jspb.BinaryReader): ValidateResult;
  }

  export namespace ValidateResult {
    export type AsObject = {
    }
  }


  export class AuthOp extends jspb.Message {
    getCheckOnly(): boolean;
    setCheckOnly(value: boolean): AuthOp;

    getComponent(): Ref.Component | undefined;
    setComponent(value?: Ref.Component): AuthOp;
    hasComponent(): boolean;
    clearComponent(): AuthOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AuthOp.AsObject;
    static toObject(includeInstance: boolean, msg: AuthOp): AuthOp.AsObject;
    static serializeBinaryToWriter(message: AuthOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AuthOp;
    static deserializeBinaryFromReader(message: AuthOp, reader: jspb.BinaryReader): AuthOp;
  }

  export namespace AuthOp {
    export type AsObject = {
      checkOnly: boolean,
      component?: Ref.Component.AsObject,
    }
  }


  export class AuthResult extends jspb.Message {
    getResultsList(): Array<Job.AuthResult.Result>;
    setResultsList(value: Array<Job.AuthResult.Result>): AuthResult;
    clearResultsList(): AuthResult;
    addResults(value?: Job.AuthResult.Result, index?: number): Job.AuthResult.Result;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AuthResult.AsObject;
    static toObject(includeInstance: boolean, msg: AuthResult): AuthResult.AsObject;
    static serializeBinaryToWriter(message: AuthResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AuthResult;
    static deserializeBinaryFromReader(message: AuthResult, reader: jspb.BinaryReader): AuthResult;
  }

  export namespace AuthResult {
    export type AsObject = {
      resultsList: Array<Job.AuthResult.Result.AsObject>,
    }

    export class Result extends jspb.Message {
      getComponent(): Component | undefined;
      setComponent(value?: Component): Result;
      hasComponent(): boolean;
      clearComponent(): Result;

      getCheckResult(): boolean;
      setCheckResult(value: boolean): Result;

      getCheckError(): google_rpc_status_pb.Status | undefined;
      setCheckError(value?: google_rpc_status_pb.Status): Result;
      hasCheckError(): boolean;
      clearCheckError(): Result;

      getAuthCompleted(): boolean;
      setAuthCompleted(value: boolean): Result;

      getAuthError(): google_rpc_status_pb.Status | undefined;
      setAuthError(value?: google_rpc_status_pb.Status): Result;
      hasAuthError(): boolean;
      clearAuthError(): Result;

      getAuthSupported(): boolean;
      setAuthSupported(value: boolean): Result;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Result.AsObject;
      static toObject(includeInstance: boolean, msg: Result): Result.AsObject;
      static serializeBinaryToWriter(message: Result, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Result;
      static deserializeBinaryFromReader(message: Result, reader: jspb.BinaryReader): Result;
    }

    export namespace Result {
      export type AsObject = {
        component?: Component.AsObject,
        checkResult: boolean,
        checkError?: google_rpc_status_pb.Status.AsObject,
        authCompleted: boolean,
        authError?: google_rpc_status_pb.Status.AsObject,
        authSupported: boolean,
      }
    }

  }


  export class BuildOp extends jspb.Message {
    getDisablePush(): boolean;
    setDisablePush(value: boolean): BuildOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BuildOp.AsObject;
    static toObject(includeInstance: boolean, msg: BuildOp): BuildOp.AsObject;
    static serializeBinaryToWriter(message: BuildOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BuildOp;
    static deserializeBinaryFromReader(message: BuildOp, reader: jspb.BinaryReader): BuildOp;
  }

  export namespace BuildOp {
    export type AsObject = {
      disablePush: boolean,
    }
  }


  export class BuildResult extends jspb.Message {
    getBuild(): Build | undefined;
    setBuild(value?: Build): BuildResult;
    hasBuild(): boolean;
    clearBuild(): BuildResult;

    getPush(): PushedArtifact | undefined;
    setPush(value?: PushedArtifact): BuildResult;
    hasPush(): boolean;
    clearPush(): BuildResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BuildResult.AsObject;
    static toObject(includeInstance: boolean, msg: BuildResult): BuildResult.AsObject;
    static serializeBinaryToWriter(message: BuildResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BuildResult;
    static deserializeBinaryFromReader(message: BuildResult, reader: jspb.BinaryReader): BuildResult;
  }

  export namespace BuildResult {
    export type AsObject = {
      build?: Build.AsObject,
      push?: PushedArtifact.AsObject,
    }
  }


  export class PushOp extends jspb.Message {
    getBuild(): Build | undefined;
    setBuild(value?: Build): PushOp;
    hasBuild(): boolean;
    clearBuild(): PushOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PushOp.AsObject;
    static toObject(includeInstance: boolean, msg: PushOp): PushOp.AsObject;
    static serializeBinaryToWriter(message: PushOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PushOp;
    static deserializeBinaryFromReader(message: PushOp, reader: jspb.BinaryReader): PushOp;
  }

  export namespace PushOp {
    export type AsObject = {
      build?: Build.AsObject,
    }
  }


  export class PushResult extends jspb.Message {
    getArtifact(): PushedArtifact | undefined;
    setArtifact(value?: PushedArtifact): PushResult;
    hasArtifact(): boolean;
    clearArtifact(): PushResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PushResult.AsObject;
    static toObject(includeInstance: boolean, msg: PushResult): PushResult.AsObject;
    static serializeBinaryToWriter(message: PushResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PushResult;
    static deserializeBinaryFromReader(message: PushResult, reader: jspb.BinaryReader): PushResult;
  }

  export namespace PushResult {
    export type AsObject = {
      artifact?: PushedArtifact.AsObject,
    }
  }


  export class DeployOp extends jspb.Message {
    getArtifact(): PushedArtifact | undefined;
    setArtifact(value?: PushedArtifact): DeployOp;
    hasArtifact(): boolean;
    clearArtifact(): DeployOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeployOp.AsObject;
    static toObject(includeInstance: boolean, msg: DeployOp): DeployOp.AsObject;
    static serializeBinaryToWriter(message: DeployOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeployOp;
    static deserializeBinaryFromReader(message: DeployOp, reader: jspb.BinaryReader): DeployOp;
  }

  export namespace DeployOp {
    export type AsObject = {
      artifact?: PushedArtifact.AsObject,
    }
  }


  export class DeployResult extends jspb.Message {
    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): DeployResult;
    hasDeployment(): boolean;
    clearDeployment(): DeployResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeployResult.AsObject;
    static toObject(includeInstance: boolean, msg: DeployResult): DeployResult.AsObject;
    static serializeBinaryToWriter(message: DeployResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeployResult;
    static deserializeBinaryFromReader(message: DeployResult, reader: jspb.BinaryReader): DeployResult;
  }

  export namespace DeployResult {
    export type AsObject = {
      deployment?: Deployment.AsObject,
    }
  }


  export class ExecOp extends jspb.Message {
    getInstanceId(): string;
    setInstanceId(value: string): ExecOp;

    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): ExecOp;
    hasDeployment(): boolean;
    clearDeployment(): ExecOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ExecOp.AsObject;
    static toObject(includeInstance: boolean, msg: ExecOp): ExecOp.AsObject;
    static serializeBinaryToWriter(message: ExecOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ExecOp;
    static deserializeBinaryFromReader(message: ExecOp, reader: jspb.BinaryReader): ExecOp;
  }

  export namespace ExecOp {
    export type AsObject = {
      instanceId: string,
      deployment?: Deployment.AsObject,
    }
  }


  export class ExecResult extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ExecResult.AsObject;
    static toObject(includeInstance: boolean, msg: ExecResult): ExecResult.AsObject;
    static serializeBinaryToWriter(message: ExecResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ExecResult;
    static deserializeBinaryFromReader(message: ExecResult, reader: jspb.BinaryReader): ExecResult;
  }

  export namespace ExecResult {
    export type AsObject = {
    }
  }


  export class LogsOp extends jspb.Message {
    getInstanceId(): string;
    setInstanceId(value: string): LogsOp;

    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): LogsOp;
    hasDeployment(): boolean;
    clearDeployment(): LogsOp;

    getStartTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setStartTime(value?: google_protobuf_timestamp_pb.Timestamp): LogsOp;
    hasStartTime(): boolean;
    clearStartTime(): LogsOp;

    getLimit(): number;
    setLimit(value: number): LogsOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LogsOp.AsObject;
    static toObject(includeInstance: boolean, msg: LogsOp): LogsOp.AsObject;
    static serializeBinaryToWriter(message: LogsOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LogsOp;
    static deserializeBinaryFromReader(message: LogsOp, reader: jspb.BinaryReader): LogsOp;
  }

  export namespace LogsOp {
    export type AsObject = {
      instanceId: string,
      deployment?: Deployment.AsObject,
      startTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      limit: number,
    }
  }


  export class DestroyOp extends jspb.Message {
    getWorkspace(): google_protobuf_empty_pb.Empty | undefined;
    setWorkspace(value?: google_protobuf_empty_pb.Empty): DestroyOp;
    hasWorkspace(): boolean;
    clearWorkspace(): DestroyOp;

    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): DestroyOp;
    hasDeployment(): boolean;
    clearDeployment(): DestroyOp;

    getTargetCase(): DestroyOp.TargetCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DestroyOp.AsObject;
    static toObject(includeInstance: boolean, msg: DestroyOp): DestroyOp.AsObject;
    static serializeBinaryToWriter(message: DestroyOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DestroyOp;
    static deserializeBinaryFromReader(message: DestroyOp, reader: jspb.BinaryReader): DestroyOp;
  }

  export namespace DestroyOp {
    export type AsObject = {
      workspace?: google_protobuf_empty_pb.Empty.AsObject,
      deployment?: Deployment.AsObject,
    }

    export enum TargetCase { 
      TARGET_NOT_SET = 0,
      WORKSPACE = 1,
      DEPLOYMENT = 2,
    }
  }


  export class ReleaseOp extends jspb.Message {
    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): ReleaseOp;
    hasDeployment(): boolean;
    clearDeployment(): ReleaseOp;

    getPrune(): boolean;
    setPrune(value: boolean): ReleaseOp;

    getPruneRetain(): number;
    setPruneRetain(value: number): ReleaseOp;

    getPruneRetainOverride(): boolean;
    setPruneRetainOverride(value: boolean): ReleaseOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ReleaseOp.AsObject;
    static toObject(includeInstance: boolean, msg: ReleaseOp): ReleaseOp.AsObject;
    static serializeBinaryToWriter(message: ReleaseOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ReleaseOp;
    static deserializeBinaryFromReader(message: ReleaseOp, reader: jspb.BinaryReader): ReleaseOp;
  }

  export namespace ReleaseOp {
    export type AsObject = {
      deployment?: Deployment.AsObject,
      prune: boolean,
      pruneRetain: number,
      pruneRetainOverride: boolean,
    }
  }


  export class ReleaseResult extends jspb.Message {
    getRelease(): Release | undefined;
    setRelease(value?: Release): ReleaseResult;
    hasRelease(): boolean;
    clearRelease(): ReleaseResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ReleaseResult.AsObject;
    static toObject(includeInstance: boolean, msg: ReleaseResult): ReleaseResult.AsObject;
    static serializeBinaryToWriter(message: ReleaseResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ReleaseResult;
    static deserializeBinaryFromReader(message: ReleaseResult, reader: jspb.BinaryReader): ReleaseResult;
  }

  export namespace ReleaseResult {
    export type AsObject = {
      release?: Release.AsObject,
    }
  }


  export class TaskPluginParams extends jspb.Message {
    getPluginType(): string;
    setPluginType(value: string): TaskPluginParams;

    getHclConfig(): Uint8Array | string;
    getHclConfig_asU8(): Uint8Array;
    getHclConfig_asB64(): string;
    setHclConfig(value: Uint8Array | string): TaskPluginParams;

    getHclFormat(): Hcl.Format;
    setHclFormat(value: Hcl.Format): TaskPluginParams;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TaskPluginParams.AsObject;
    static toObject(includeInstance: boolean, msg: TaskPluginParams): TaskPluginParams.AsObject;
    static serializeBinaryToWriter(message: TaskPluginParams, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TaskPluginParams;
    static deserializeBinaryFromReader(message: TaskPluginParams, reader: jspb.BinaryReader): TaskPluginParams;
  }

  export namespace TaskPluginParams {
    export type AsObject = {
      pluginType: string,
      hclConfig: Uint8Array | string,
      hclFormat: Hcl.Format,
    }
  }


  export class StartTaskLaunchOp extends jspb.Message {
    getInfo(): TaskLaunchInfo | undefined;
    setInfo(value?: TaskLaunchInfo): StartTaskLaunchOp;
    hasInfo(): boolean;
    clearInfo(): StartTaskLaunchOp;

    getParams(): Job.TaskPluginParams | undefined;
    setParams(value?: Job.TaskPluginParams): StartTaskLaunchOp;
    hasParams(): boolean;
    clearParams(): StartTaskLaunchOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StartTaskLaunchOp.AsObject;
    static toObject(includeInstance: boolean, msg: StartTaskLaunchOp): StartTaskLaunchOp.AsObject;
    static serializeBinaryToWriter(message: StartTaskLaunchOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StartTaskLaunchOp;
    static deserializeBinaryFromReader(message: StartTaskLaunchOp, reader: jspb.BinaryReader): StartTaskLaunchOp;
  }

  export namespace StartTaskLaunchOp {
    export type AsObject = {
      info?: TaskLaunchInfo.AsObject,
      params?: Job.TaskPluginParams.AsObject,
    }
  }


  export class StartTaskResult extends jspb.Message {
    getState(): any_pb.Any | undefined;
    setState(value?: any_pb.Any): StartTaskResult;
    hasState(): boolean;
    clearState(): StartTaskResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StartTaskResult.AsObject;
    static toObject(includeInstance: boolean, msg: StartTaskResult): StartTaskResult.AsObject;
    static serializeBinaryToWriter(message: StartTaskResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StartTaskResult;
    static deserializeBinaryFromReader(message: StartTaskResult, reader: jspb.BinaryReader): StartTaskResult;
  }

  export namespace StartTaskResult {
    export type AsObject = {
      state?: any_pb.Any.AsObject,
    }
  }


  export class StopTaskLaunchOp extends jspb.Message {
    getParams(): Job.TaskPluginParams | undefined;
    setParams(value?: Job.TaskPluginParams): StopTaskLaunchOp;
    hasParams(): boolean;
    clearParams(): StopTaskLaunchOp;

    getDirect(): any_pb.Any | undefined;
    setDirect(value?: any_pb.Any): StopTaskLaunchOp;
    hasDirect(): boolean;
    clearDirect(): StopTaskLaunchOp;

    getStartJobId(): string;
    setStartJobId(value: string): StopTaskLaunchOp;

    getStateCase(): StopTaskLaunchOp.StateCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StopTaskLaunchOp.AsObject;
    static toObject(includeInstance: boolean, msg: StopTaskLaunchOp): StopTaskLaunchOp.AsObject;
    static serializeBinaryToWriter(message: StopTaskLaunchOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StopTaskLaunchOp;
    static deserializeBinaryFromReader(message: StopTaskLaunchOp, reader: jspb.BinaryReader): StopTaskLaunchOp;
  }

  export namespace StopTaskLaunchOp {
    export type AsObject = {
      params?: Job.TaskPluginParams.AsObject,
      direct?: any_pb.Any.AsObject,
      startJobId: string,
    }

    export enum StateCase { 
      STATE_NOT_SET = 0,
      DIRECT = 1,
      START_JOB_ID = 3,
    }
  }


  export class WatchTaskOp extends jspb.Message {
    getStartJob(): Ref.Job | undefined;
    setStartJob(value?: Ref.Job): WatchTaskOp;
    hasStartJob(): boolean;
    clearStartJob(): WatchTaskOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): WatchTaskOp.AsObject;
    static toObject(includeInstance: boolean, msg: WatchTaskOp): WatchTaskOp.AsObject;
    static serializeBinaryToWriter(message: WatchTaskOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): WatchTaskOp;
    static deserializeBinaryFromReader(message: WatchTaskOp, reader: jspb.BinaryReader): WatchTaskOp;
  }

  export namespace WatchTaskOp {
    export type AsObject = {
      startJob?: Ref.Job.AsObject,
    }
  }


  export class WatchTaskResult extends jspb.Message {
    getExitCode(): number;
    setExitCode(value: number): WatchTaskResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): WatchTaskResult.AsObject;
    static toObject(includeInstance: boolean, msg: WatchTaskResult): WatchTaskResult.AsObject;
    static serializeBinaryToWriter(message: WatchTaskResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): WatchTaskResult;
    static deserializeBinaryFromReader(message: WatchTaskResult, reader: jspb.BinaryReader): WatchTaskResult;
  }

  export namespace WatchTaskResult {
    export type AsObject = {
      exitCode: number,
    }
  }


  export class StatusReportOp extends jspb.Message {
    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): StatusReportOp;
    hasDeployment(): boolean;
    clearDeployment(): StatusReportOp;

    getRelease(): Release | undefined;
    setRelease(value?: Release): StatusReportOp;
    hasRelease(): boolean;
    clearRelease(): StatusReportOp;

    getTargetCase(): StatusReportOp.TargetCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StatusReportOp.AsObject;
    static toObject(includeInstance: boolean, msg: StatusReportOp): StatusReportOp.AsObject;
    static serializeBinaryToWriter(message: StatusReportOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StatusReportOp;
    static deserializeBinaryFromReader(message: StatusReportOp, reader: jspb.BinaryReader): StatusReportOp;
  }

  export namespace StatusReportOp {
    export type AsObject = {
      deployment?: Deployment.AsObject,
      release?: Release.AsObject,
    }

    export enum TargetCase { 
      TARGET_NOT_SET = 0,
      DEPLOYMENT = 1,
      RELEASE = 2,
    }
  }


  export class StatusReportResult extends jspb.Message {
    getStatusReport(): StatusReport | undefined;
    setStatusReport(value?: StatusReport): StatusReportResult;
    hasStatusReport(): boolean;
    clearStatusReport(): StatusReportResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StatusReportResult.AsObject;
    static toObject(includeInstance: boolean, msg: StatusReportResult): StatusReportResult.AsObject;
    static serializeBinaryToWriter(message: StatusReportResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StatusReportResult;
    static deserializeBinaryFromReader(message: StatusReportResult, reader: jspb.BinaryReader): StatusReportResult;
  }

  export namespace StatusReportResult {
    export type AsObject = {
      statusReport?: StatusReport.AsObject,
    }
  }


  export class DocsOp extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DocsOp.AsObject;
    static toObject(includeInstance: boolean, msg: DocsOp): DocsOp.AsObject;
    static serializeBinaryToWriter(message: DocsOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DocsOp;
    static deserializeBinaryFromReader(message: DocsOp, reader: jspb.BinaryReader): DocsOp;
  }

  export namespace DocsOp {
    export type AsObject = {
    }
  }


  export class DocsResult extends jspb.Message {
    getResultsList(): Array<Job.DocsResult.Result>;
    setResultsList(value: Array<Job.DocsResult.Result>): DocsResult;
    clearResultsList(): DocsResult;
    addResults(value?: Job.DocsResult.Result, index?: number): Job.DocsResult.Result;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DocsResult.AsObject;
    static toObject(includeInstance: boolean, msg: DocsResult): DocsResult.AsObject;
    static serializeBinaryToWriter(message: DocsResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DocsResult;
    static deserializeBinaryFromReader(message: DocsResult, reader: jspb.BinaryReader): DocsResult;
  }

  export namespace DocsResult {
    export type AsObject = {
      resultsList: Array<Job.DocsResult.Result.AsObject>,
    }

    export class Result extends jspb.Message {
      getComponent(): Component | undefined;
      setComponent(value?: Component): Result;
      hasComponent(): boolean;
      clearComponent(): Result;

      getDocs(): Documentation | undefined;
      setDocs(value?: Documentation): Result;
      hasDocs(): boolean;
      clearDocs(): Result;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Result.AsObject;
      static toObject(includeInstance: boolean, msg: Result): Result.AsObject;
      static serializeBinaryToWriter(message: Result, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Result;
      static deserializeBinaryFromReader(message: Result, reader: jspb.BinaryReader): Result;
    }

    export namespace Result {
      export type AsObject = {
        component?: Component.AsObject,
        docs?: Documentation.AsObject,
      }
    }

  }


  export class ConfigSyncOp extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConfigSyncOp.AsObject;
    static toObject(includeInstance: boolean, msg: ConfigSyncOp): ConfigSyncOp.AsObject;
    static serializeBinaryToWriter(message: ConfigSyncOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConfigSyncOp;
    static deserializeBinaryFromReader(message: ConfigSyncOp, reader: jspb.BinaryReader): ConfigSyncOp;
  }

  export namespace ConfigSyncOp {
    export type AsObject = {
    }
  }


  export class ConfigSyncResult extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConfigSyncResult.AsObject;
    static toObject(includeInstance: boolean, msg: ConfigSyncResult): ConfigSyncResult.AsObject;
    static serializeBinaryToWriter(message: ConfigSyncResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConfigSyncResult;
    static deserializeBinaryFromReader(message: ConfigSyncResult, reader: jspb.BinaryReader): ConfigSyncResult;
  }

  export namespace ConfigSyncResult {
    export type AsObject = {
    }
  }


  export class PipelineConfigSyncResult extends jspb.Message {
    getSyncedPipelinesMap(): jspb.Map<string, Ref.Pipeline>;
    clearSyncedPipelinesMap(): PipelineConfigSyncResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PipelineConfigSyncResult.AsObject;
    static toObject(includeInstance: boolean, msg: PipelineConfigSyncResult): PipelineConfigSyncResult.AsObject;
    static serializeBinaryToWriter(message: PipelineConfigSyncResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PipelineConfigSyncResult;
    static deserializeBinaryFromReader(message: PipelineConfigSyncResult, reader: jspb.BinaryReader): PipelineConfigSyncResult;
  }

  export namespace PipelineConfigSyncResult {
    export type AsObject = {
      syncedPipelinesMap: Array<[string, Ref.Pipeline.AsObject]>,
    }
  }


  export class PollOp extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PollOp.AsObject;
    static toObject(includeInstance: boolean, msg: PollOp): PollOp.AsObject;
    static serializeBinaryToWriter(message: PollOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PollOp;
    static deserializeBinaryFromReader(message: PollOp, reader: jspb.BinaryReader): PollOp;
  }

  export namespace PollOp {
    export type AsObject = {
    }
  }


  export class PollResult extends jspb.Message {
    getJobId(): string;
    setJobId(value: string): PollResult;

    getOldRef(): Job.DataSource.Ref | undefined;
    setOldRef(value?: Job.DataSource.Ref): PollResult;
    hasOldRef(): boolean;
    clearOldRef(): PollResult;

    getNewRef(): Job.DataSource.Ref | undefined;
    setNewRef(value?: Job.DataSource.Ref): PollResult;
    hasNewRef(): boolean;
    clearNewRef(): PollResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PollResult.AsObject;
    static toObject(includeInstance: boolean, msg: PollResult): PollResult.AsObject;
    static serializeBinaryToWriter(message: PollResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PollResult;
    static deserializeBinaryFromReader(message: PollResult, reader: jspb.BinaryReader): PollResult;
  }

  export namespace PollResult {
    export type AsObject = {
      jobId: string,
      oldRef?: Job.DataSource.Ref.AsObject,
      newRef?: Job.DataSource.Ref.AsObject,
    }
  }


  export class InitOp extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InitOp.AsObject;
    static toObject(includeInstance: boolean, msg: InitOp): InitOp.AsObject;
    static serializeBinaryToWriter(message: InitOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InitOp;
    static deserializeBinaryFromReader(message: InitOp, reader: jspb.BinaryReader): InitOp;
  }

  export namespace InitOp {
    export type AsObject = {
    }
  }


  export class InitResult extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InitResult.AsObject;
    static toObject(includeInstance: boolean, msg: InitResult): InitResult.AsObject;
    static serializeBinaryToWriter(message: InitResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InitResult;
    static deserializeBinaryFromReader(message: InitResult, reader: jspb.BinaryReader): InitResult;
  }

  export namespace InitResult {
    export type AsObject = {
    }
  }


  export class PipelineStepOp extends jspb.Message {
    getStep(): Pipeline.Step | undefined;
    setStep(value?: Pipeline.Step): PipelineStepOp;
    hasStep(): boolean;
    clearStep(): PipelineStepOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PipelineStepOp.AsObject;
    static toObject(includeInstance: boolean, msg: PipelineStepOp): PipelineStepOp.AsObject;
    static serializeBinaryToWriter(message: PipelineStepOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PipelineStepOp;
    static deserializeBinaryFromReader(message: PipelineStepOp, reader: jspb.BinaryReader): PipelineStepOp;
  }

  export namespace PipelineStepOp {
    export type AsObject = {
      step?: Pipeline.Step.AsObject,
    }
  }


  export class PipelineStepResult extends jspb.Message {
    getResult(): google_rpc_status_pb.Status | undefined;
    setResult(value?: google_rpc_status_pb.Status): PipelineStepResult;
    hasResult(): boolean;
    clearResult(): PipelineStepResult;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PipelineStepResult.AsObject;
    static toObject(includeInstance: boolean, msg: PipelineStepResult): PipelineStepResult.AsObject;
    static serializeBinaryToWriter(message: PipelineStepResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PipelineStepResult;
    static deserializeBinaryFromReader(message: PipelineStepResult, reader: jspb.BinaryReader): PipelineStepResult;
  }

  export namespace PipelineStepResult {
    export type AsObject = {
      result?: google_rpc_status_pb.Status.AsObject,
    }
  }


  export class QueueProjectOp extends jspb.Message {
    getJobTemplate(): Job | undefined;
    setJobTemplate(value?: Job): QueueProjectOp;
    hasJobTemplate(): boolean;
    clearJobTemplate(): QueueProjectOp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueueProjectOp.AsObject;
    static toObject(includeInstance: boolean, msg: QueueProjectOp): QueueProjectOp.AsObject;
    static serializeBinaryToWriter(message: QueueProjectOp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueueProjectOp;
    static deserializeBinaryFromReader(message: QueueProjectOp, reader: jspb.BinaryReader): QueueProjectOp;
  }

  export namespace QueueProjectOp {
    export type AsObject = {
      jobTemplate?: Job.AsObject,
    }
  }


  export class QueueProjectResult extends jspb.Message {
    getApplicationsList(): Array<Job.QueueProjectResult.Application>;
    setApplicationsList(value: Array<Job.QueueProjectResult.Application>): QueueProjectResult;
    clearApplicationsList(): QueueProjectResult;
    addApplications(value?: Job.QueueProjectResult.Application, index?: number): Job.QueueProjectResult.Application;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueueProjectResult.AsObject;
    static toObject(includeInstance: boolean, msg: QueueProjectResult): QueueProjectResult.AsObject;
    static serializeBinaryToWriter(message: QueueProjectResult, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueueProjectResult;
    static deserializeBinaryFromReader(message: QueueProjectResult, reader: jspb.BinaryReader): QueueProjectResult;
  }

  export namespace QueueProjectResult {
    export type AsObject = {
      applicationsList: Array<Job.QueueProjectResult.Application.AsObject>,
    }

    export class Application extends jspb.Message {
      getApplication(): Ref.Application | undefined;
      setApplication(value?: Ref.Application): Application;
      hasApplication(): boolean;
      clearApplication(): Application;

      getJobId(): string;
      setJobId(value: string): Application;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Application.AsObject;
      static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
      static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Application;
      static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
    }

    export namespace Application {
      export type AsObject = {
        application?: Ref.Application.AsObject,
        jobId: string,
      }
    }

  }


  export enum State { 
    UNKNOWN = 0,
    QUEUED = 1,
    WAITING = 2,
    RUNNING = 3,
    ERROR = 4,
    SUCCESS = 5,
  }

  export enum OperationCase { 
    OPERATION_NOT_SET = 0,
    NOOP = 50,
    BUILD = 51,
    PUSH = 52,
    DEPLOY = 53,
    DESTROY = 54,
    RELEASE = 55,
    VALIDATE = 56,
    AUTH = 57,
    DOCS = 58,
    CONFIG_SYNC = 59,
    EXEC = 60,
    UP = 61,
    LOGS = 62,
    QUEUE_PROJECT = 63,
    POLL = 64,
    STATUS_REPORT = 65,
    START_TASK = 66,
    STOP_TASK = 67,
    INIT = 68,
    WATCH_TASK = 69,
    PIPELINE_STEP = 70,
  }
}

export class Documentation extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): Documentation;

  getExample(): string;
  setExample(value: string): Documentation;

  getInput(): string;
  setInput(value: string): Documentation;

  getOutput(): string;
  setOutput(value: string): Documentation;

  getFieldsMap(): jspb.Map<string, Documentation.Field>;
  clearFieldsMap(): Documentation;

  getMappersList(): Array<Documentation.Mapper>;
  setMappersList(value: Array<Documentation.Mapper>): Documentation;
  clearMappersList(): Documentation;
  addMappers(value?: Documentation.Mapper, index?: number): Documentation.Mapper;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Documentation.AsObject;
  static toObject(includeInstance: boolean, msg: Documentation): Documentation.AsObject;
  static serializeBinaryToWriter(message: Documentation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Documentation;
  static deserializeBinaryFromReader(message: Documentation, reader: jspb.BinaryReader): Documentation;
}

export namespace Documentation {
  export type AsObject = {
    description: string,
    example: string,
    input: string,
    output: string,
    fieldsMap: Array<[string, Documentation.Field.AsObject]>,
    mappersList: Array<Documentation.Mapper.AsObject>,
  }

  export class Field extends jspb.Message {
    getName(): string;
    setName(value: string): Field;

    getSynopsis(): string;
    setSynopsis(value: string): Field;

    getSummary(): string;
    setSummary(value: string): Field;

    getOptional(): boolean;
    setOptional(value: boolean): Field;

    getEnvVar(): string;
    setEnvVar(value: string): Field;

    getType(): string;
    setType(value: string): Field;

    getDefault(): string;
    setDefault(value: string): Field;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Field.AsObject;
    static toObject(includeInstance: boolean, msg: Field): Field.AsObject;
    static serializeBinaryToWriter(message: Field, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Field;
    static deserializeBinaryFromReader(message: Field, reader: jspb.BinaryReader): Field;
  }

  export namespace Field {
    export type AsObject = {
      name: string,
      synopsis: string,
      summary: string,
      optional: boolean,
      envVar: string,
      type: string,
      pb_default: string,
    }
  }


  export class Mapper extends jspb.Message {
    getInput(): string;
    setInput(value: string): Mapper;

    getOutput(): string;
    setOutput(value: string): Mapper;

    getDescription(): string;
    setDescription(value: string): Mapper;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Mapper.AsObject;
    static toObject(includeInstance: boolean, msg: Mapper): Mapper.AsObject;
    static serializeBinaryToWriter(message: Mapper, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Mapper;
    static deserializeBinaryFromReader(message: Mapper, reader: jspb.BinaryReader): Mapper;
  }

  export namespace Mapper {
    export type AsObject = {
      input: string,
      output: string,
      description: string,
    }
  }

}

export class GetJobRequest extends jspb.Message {
  getJobId(): string;
  setJobId(value: string): GetJobRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetJobRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetJobRequest): GetJobRequest.AsObject;
  static serializeBinaryToWriter(message: GetJobRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetJobRequest;
  static deserializeBinaryFromReader(message: GetJobRequest, reader: jspb.BinaryReader): GetJobRequest;
}

export namespace GetJobRequest {
  export type AsObject = {
    jobId: string,
  }
}

export class ListJobsRequest extends jspb.Message {
  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ListJobsRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ListJobsRequest;

  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): ListJobsRequest;
  hasProject(): boolean;
  clearProject(): ListJobsRequest;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListJobsRequest;
  hasApplication(): boolean;
  clearApplication(): ListJobsRequest;

  getTargetrunner(): Ref.Runner | undefined;
  setTargetrunner(value?: Ref.Runner): ListJobsRequest;
  hasTargetrunner(): boolean;
  clearTargetrunner(): ListJobsRequest;

  getJobstateList(): Array<Job.State>;
  setJobstateList(value: Array<Job.State>): ListJobsRequest;
  clearJobstateList(): ListJobsRequest;
  addJobstate(value: Job.State, index?: number): ListJobsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListJobsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListJobsRequest): ListJobsRequest.AsObject;
  static serializeBinaryToWriter(message: ListJobsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListJobsRequest;
  static deserializeBinaryFromReader(message: ListJobsRequest, reader: jspb.BinaryReader): ListJobsRequest;
}

export namespace ListJobsRequest {
  export type AsObject = {
    workspace?: Ref.Workspace.AsObject,
    project?: Ref.Project.AsObject,
    application?: Ref.Application.AsObject,
    targetrunner?: Ref.Runner.AsObject,
    jobstateList: Array<Job.State>,
  }
}

export class ListJobsResponse extends jspb.Message {
  getJobsList(): Array<Job>;
  setJobsList(value: Array<Job>): ListJobsResponse;
  clearJobsList(): ListJobsResponse;
  addJobs(value?: Job, index?: number): Job;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListJobsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListJobsResponse): ListJobsResponse.AsObject;
  static serializeBinaryToWriter(message: ListJobsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListJobsResponse;
  static deserializeBinaryFromReader(message: ListJobsResponse, reader: jspb.BinaryReader): ListJobsResponse;
}

export namespace ListJobsResponse {
  export type AsObject = {
    jobsList: Array<Job.AsObject>,
  }
}

export class GetJobStreamRequest extends jspb.Message {
  getJobId(): string;
  setJobId(value: string): GetJobStreamRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetJobStreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetJobStreamRequest): GetJobStreamRequest.AsObject;
  static serializeBinaryToWriter(message: GetJobStreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetJobStreamRequest;
  static deserializeBinaryFromReader(message: GetJobStreamRequest, reader: jspb.BinaryReader): GetJobStreamRequest;
}

export namespace GetJobStreamRequest {
  export type AsObject = {
    jobId: string,
  }
}

export class GetJobStreamResponse extends jspb.Message {
  getOpen(): GetJobStreamResponse.Open | undefined;
  setOpen(value?: GetJobStreamResponse.Open): GetJobStreamResponse;
  hasOpen(): boolean;
  clearOpen(): GetJobStreamResponse;

  getState(): GetJobStreamResponse.State | undefined;
  setState(value?: GetJobStreamResponse.State): GetJobStreamResponse;
  hasState(): boolean;
  clearState(): GetJobStreamResponse;

  getJob(): GetJobStreamResponse.JobChange | undefined;
  setJob(value?: GetJobStreamResponse.JobChange): GetJobStreamResponse;
  hasJob(): boolean;
  clearJob(): GetJobStreamResponse;

  getTerminal(): GetJobStreamResponse.Terminal | undefined;
  setTerminal(value?: GetJobStreamResponse.Terminal): GetJobStreamResponse;
  hasTerminal(): boolean;
  clearTerminal(): GetJobStreamResponse;

  getDownload(): GetJobStreamResponse.Download | undefined;
  setDownload(value?: GetJobStreamResponse.Download): GetJobStreamResponse;
  hasDownload(): boolean;
  clearDownload(): GetJobStreamResponse;

  getError(): GetJobStreamResponse.Error | undefined;
  setError(value?: GetJobStreamResponse.Error): GetJobStreamResponse;
  hasError(): boolean;
  clearError(): GetJobStreamResponse;

  getComplete(): GetJobStreamResponse.Complete | undefined;
  setComplete(value?: GetJobStreamResponse.Complete): GetJobStreamResponse;
  hasComplete(): boolean;
  clearComplete(): GetJobStreamResponse;

  getEventCase(): GetJobStreamResponse.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetJobStreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetJobStreamResponse): GetJobStreamResponse.AsObject;
  static serializeBinaryToWriter(message: GetJobStreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetJobStreamResponse;
  static deserializeBinaryFromReader(message: GetJobStreamResponse, reader: jspb.BinaryReader): GetJobStreamResponse;
}

export namespace GetJobStreamResponse {
  export type AsObject = {
    open?: GetJobStreamResponse.Open.AsObject,
    state?: GetJobStreamResponse.State.AsObject,
    job?: GetJobStreamResponse.JobChange.AsObject,
    terminal?: GetJobStreamResponse.Terminal.AsObject,
    download?: GetJobStreamResponse.Download.AsObject,
    error?: GetJobStreamResponse.Error.AsObject,
    complete?: GetJobStreamResponse.Complete.AsObject,
  }

  export class Open extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Open.AsObject;
    static toObject(includeInstance: boolean, msg: Open): Open.AsObject;
    static serializeBinaryToWriter(message: Open, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Open;
    static deserializeBinaryFromReader(message: Open, reader: jspb.BinaryReader): Open;
  }

  export namespace Open {
    export type AsObject = {
    }
  }


  export class State extends jspb.Message {
    getPrevious(): Job.State;
    setPrevious(value: Job.State): State;

    getCurrent(): Job.State;
    setCurrent(value: Job.State): State;

    getJob(): Job | undefined;
    setJob(value?: Job): State;
    hasJob(): boolean;
    clearJob(): State;

    getCanceling(): boolean;
    setCanceling(value: boolean): State;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): State.AsObject;
    static toObject(includeInstance: boolean, msg: State): State.AsObject;
    static serializeBinaryToWriter(message: State, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): State;
    static deserializeBinaryFromReader(message: State, reader: jspb.BinaryReader): State;
  }

  export namespace State {
    export type AsObject = {
      previous: Job.State,
      current: Job.State,
      job?: Job.AsObject,
      canceling: boolean,
    }
  }


  export class JobChange extends jspb.Message {
    getJob(): Job | undefined;
    setJob(value?: Job): JobChange;
    hasJob(): boolean;
    clearJob(): JobChange;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): JobChange.AsObject;
    static toObject(includeInstance: boolean, msg: JobChange): JobChange.AsObject;
    static serializeBinaryToWriter(message: JobChange, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): JobChange;
    static deserializeBinaryFromReader(message: JobChange, reader: jspb.BinaryReader): JobChange;
  }

  export namespace JobChange {
    export type AsObject = {
      job?: Job.AsObject,
    }
  }


  export class Download extends jspb.Message {
    getDataSourceRef(): Job.DataSource.Ref | undefined;
    setDataSourceRef(value?: Job.DataSource.Ref): Download;
    hasDataSourceRef(): boolean;
    clearDataSourceRef(): Download;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Download.AsObject;
    static toObject(includeInstance: boolean, msg: Download): Download.AsObject;
    static serializeBinaryToWriter(message: Download, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Download;
    static deserializeBinaryFromReader(message: Download, reader: jspb.BinaryReader): Download;
  }

  export namespace Download {
    export type AsObject = {
      dataSourceRef?: Job.DataSource.Ref.AsObject,
    }
  }


  export class Terminal extends jspb.Message {
    getEventsList(): Array<GetJobStreamResponse.Terminal.Event>;
    setEventsList(value: Array<GetJobStreamResponse.Terminal.Event>): Terminal;
    clearEventsList(): Terminal;
    addEvents(value?: GetJobStreamResponse.Terminal.Event, index?: number): GetJobStreamResponse.Terminal.Event;

    getBuffered(): boolean;
    setBuffered(value: boolean): Terminal;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Terminal.AsObject;
    static toObject(includeInstance: boolean, msg: Terminal): Terminal.AsObject;
    static serializeBinaryToWriter(message: Terminal, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Terminal;
    static deserializeBinaryFromReader(message: Terminal, reader: jspb.BinaryReader): Terminal;
  }

  export namespace Terminal {
    export type AsObject = {
      eventsList: Array<GetJobStreamResponse.Terminal.Event.AsObject>,
      buffered: boolean,
    }

    export class Event extends jspb.Message {
      getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
      setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): Event;
      hasTimestamp(): boolean;
      clearTimestamp(): Event;

      getLine(): GetJobStreamResponse.Terminal.Event.Line | undefined;
      setLine(value?: GetJobStreamResponse.Terminal.Event.Line): Event;
      hasLine(): boolean;
      clearLine(): Event;

      getStatus(): GetJobStreamResponse.Terminal.Event.Status | undefined;
      setStatus(value?: GetJobStreamResponse.Terminal.Event.Status): Event;
      hasStatus(): boolean;
      clearStatus(): Event;

      getNamedValues(): GetJobStreamResponse.Terminal.Event.NamedValues | undefined;
      setNamedValues(value?: GetJobStreamResponse.Terminal.Event.NamedValues): Event;
      hasNamedValues(): boolean;
      clearNamedValues(): Event;

      getRaw(): GetJobStreamResponse.Terminal.Event.Raw | undefined;
      setRaw(value?: GetJobStreamResponse.Terminal.Event.Raw): Event;
      hasRaw(): boolean;
      clearRaw(): Event;

      getTable(): GetJobStreamResponse.Terminal.Event.Table | undefined;
      setTable(value?: GetJobStreamResponse.Terminal.Event.Table): Event;
      hasTable(): boolean;
      clearTable(): Event;

      getStepGroup(): GetJobStreamResponse.Terminal.Event.StepGroup | undefined;
      setStepGroup(value?: GetJobStreamResponse.Terminal.Event.StepGroup): Event;
      hasStepGroup(): boolean;
      clearStepGroup(): Event;

      getStep(): GetJobStreamResponse.Terminal.Event.Step | undefined;
      setStep(value?: GetJobStreamResponse.Terminal.Event.Step): Event;
      hasStep(): boolean;
      clearStep(): Event;

      getEventCase(): Event.EventCase;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Event.AsObject;
      static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
      static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Event;
      static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
    }

    export namespace Event {
      export type AsObject = {
        timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        line?: GetJobStreamResponse.Terminal.Event.Line.AsObject,
        status?: GetJobStreamResponse.Terminal.Event.Status.AsObject,
        namedValues?: GetJobStreamResponse.Terminal.Event.NamedValues.AsObject,
        raw?: GetJobStreamResponse.Terminal.Event.Raw.AsObject,
        table?: GetJobStreamResponse.Terminal.Event.Table.AsObject,
        stepGroup?: GetJobStreamResponse.Terminal.Event.StepGroup.AsObject,
        step?: GetJobStreamResponse.Terminal.Event.Step.AsObject,
      }

      export class Status extends jspb.Message {
        getStatus(): string;
        setStatus(value: string): Status;

        getMsg(): string;
        setMsg(value: string): Status;

        getStep(): boolean;
        setStep(value: boolean): Status;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): Status.AsObject;
        static toObject(includeInstance: boolean, msg: Status): Status.AsObject;
        static serializeBinaryToWriter(message: Status, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): Status;
        static deserializeBinaryFromReader(message: Status, reader: jspb.BinaryReader): Status;
      }

      export namespace Status {
        export type AsObject = {
          status: string,
          msg: string,
          step: boolean,
        }
      }


      export class Line extends jspb.Message {
        getMsg(): string;
        setMsg(value: string): Line;

        getStyle(): string;
        setStyle(value: string): Line;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): Line.AsObject;
        static toObject(includeInstance: boolean, msg: Line): Line.AsObject;
        static serializeBinaryToWriter(message: Line, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): Line;
        static deserializeBinaryFromReader(message: Line, reader: jspb.BinaryReader): Line;
      }

      export namespace Line {
        export type AsObject = {
          msg: string,
          style: string,
        }
      }


      export class Raw extends jspb.Message {
        getData(): Uint8Array | string;
        getData_asU8(): Uint8Array;
        getData_asB64(): string;
        setData(value: Uint8Array | string): Raw;

        getStderr(): boolean;
        setStderr(value: boolean): Raw;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): Raw.AsObject;
        static toObject(includeInstance: boolean, msg: Raw): Raw.AsObject;
        static serializeBinaryToWriter(message: Raw, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): Raw;
        static deserializeBinaryFromReader(message: Raw, reader: jspb.BinaryReader): Raw;
      }

      export namespace Raw {
        export type AsObject = {
          data: Uint8Array | string,
          stderr: boolean,
        }
      }


      export class NamedValue extends jspb.Message {
        getName(): string;
        setName(value: string): NamedValue;

        getValue(): string;
        setValue(value: string): NamedValue;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): NamedValue.AsObject;
        static toObject(includeInstance: boolean, msg: NamedValue): NamedValue.AsObject;
        static serializeBinaryToWriter(message: NamedValue, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): NamedValue;
        static deserializeBinaryFromReader(message: NamedValue, reader: jspb.BinaryReader): NamedValue;
      }

      export namespace NamedValue {
        export type AsObject = {
          name: string,
          value: string,
        }
      }


      export class NamedValues extends jspb.Message {
        getValuesList(): Array<GetJobStreamResponse.Terminal.Event.NamedValue>;
        setValuesList(value: Array<GetJobStreamResponse.Terminal.Event.NamedValue>): NamedValues;
        clearValuesList(): NamedValues;
        addValues(value?: GetJobStreamResponse.Terminal.Event.NamedValue, index?: number): GetJobStreamResponse.Terminal.Event.NamedValue;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): NamedValues.AsObject;
        static toObject(includeInstance: boolean, msg: NamedValues): NamedValues.AsObject;
        static serializeBinaryToWriter(message: NamedValues, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): NamedValues;
        static deserializeBinaryFromReader(message: NamedValues, reader: jspb.BinaryReader): NamedValues;
      }

      export namespace NamedValues {
        export type AsObject = {
          valuesList: Array<GetJobStreamResponse.Terminal.Event.NamedValue.AsObject>,
        }
      }


      export class TableEntry extends jspb.Message {
        getValue(): string;
        setValue(value: string): TableEntry;

        getColor(): string;
        setColor(value: string): TableEntry;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): TableEntry.AsObject;
        static toObject(includeInstance: boolean, msg: TableEntry): TableEntry.AsObject;
        static serializeBinaryToWriter(message: TableEntry, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): TableEntry;
        static deserializeBinaryFromReader(message: TableEntry, reader: jspb.BinaryReader): TableEntry;
      }

      export namespace TableEntry {
        export type AsObject = {
          value: string,
          color: string,
        }
      }


      export class TableRow extends jspb.Message {
        getEntriesList(): Array<GetJobStreamResponse.Terminal.Event.TableEntry>;
        setEntriesList(value: Array<GetJobStreamResponse.Terminal.Event.TableEntry>): TableRow;
        clearEntriesList(): TableRow;
        addEntries(value?: GetJobStreamResponse.Terminal.Event.TableEntry, index?: number): GetJobStreamResponse.Terminal.Event.TableEntry;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): TableRow.AsObject;
        static toObject(includeInstance: boolean, msg: TableRow): TableRow.AsObject;
        static serializeBinaryToWriter(message: TableRow, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): TableRow;
        static deserializeBinaryFromReader(message: TableRow, reader: jspb.BinaryReader): TableRow;
      }

      export namespace TableRow {
        export type AsObject = {
          entriesList: Array<GetJobStreamResponse.Terminal.Event.TableEntry.AsObject>,
        }
      }


      export class Table extends jspb.Message {
        getHeadersList(): Array<string>;
        setHeadersList(value: Array<string>): Table;
        clearHeadersList(): Table;
        addHeaders(value: string, index?: number): Table;

        getRowsList(): Array<GetJobStreamResponse.Terminal.Event.TableRow>;
        setRowsList(value: Array<GetJobStreamResponse.Terminal.Event.TableRow>): Table;
        clearRowsList(): Table;
        addRows(value?: GetJobStreamResponse.Terminal.Event.TableRow, index?: number): GetJobStreamResponse.Terminal.Event.TableRow;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): Table.AsObject;
        static toObject(includeInstance: boolean, msg: Table): Table.AsObject;
        static serializeBinaryToWriter(message: Table, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): Table;
        static deserializeBinaryFromReader(message: Table, reader: jspb.BinaryReader): Table;
      }

      export namespace Table {
        export type AsObject = {
          headersList: Array<string>,
          rowsList: Array<GetJobStreamResponse.Terminal.Event.TableRow.AsObject>,
        }
      }


      export class StepGroup extends jspb.Message {
        getClose(): boolean;
        setClose(value: boolean): StepGroup;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): StepGroup.AsObject;
        static toObject(includeInstance: boolean, msg: StepGroup): StepGroup.AsObject;
        static serializeBinaryToWriter(message: StepGroup, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): StepGroup;
        static deserializeBinaryFromReader(message: StepGroup, reader: jspb.BinaryReader): StepGroup;
      }

      export namespace StepGroup {
        export type AsObject = {
          close: boolean,
        }
      }


      export class Step extends jspb.Message {
        getId(): number;
        setId(value: number): Step;

        getClose(): boolean;
        setClose(value: boolean): Step;

        getMsg(): string;
        setMsg(value: string): Step;

        getStatus(): string;
        setStatus(value: string): Step;

        getOutput(): Uint8Array | string;
        getOutput_asU8(): Uint8Array;
        getOutput_asB64(): string;
        setOutput(value: Uint8Array | string): Step;

        serializeBinary(): Uint8Array;
        toObject(includeInstance?: boolean): Step.AsObject;
        static toObject(includeInstance: boolean, msg: Step): Step.AsObject;
        static serializeBinaryToWriter(message: Step, writer: jspb.BinaryWriter): void;
        static deserializeBinary(bytes: Uint8Array): Step;
        static deserializeBinaryFromReader(message: Step, reader: jspb.BinaryReader): Step;
      }

      export namespace Step {
        export type AsObject = {
          id: number,
          close: boolean,
          msg: string,
          status: string,
          output: Uint8Array | string,
        }
      }


      export enum EventCase { 
        EVENT_NOT_SET = 0,
        LINE = 2,
        STATUS = 3,
        NAMED_VALUES = 4,
        RAW = 5,
        TABLE = 6,
        STEP_GROUP = 7,
        STEP = 8,
      }
    }

  }


  export class Error extends jspb.Message {
    getError(): google_rpc_status_pb.Status | undefined;
    setError(value?: google_rpc_status_pb.Status): Error;
    hasError(): boolean;
    clearError(): Error;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Error.AsObject;
    static toObject(includeInstance: boolean, msg: Error): Error.AsObject;
    static serializeBinaryToWriter(message: Error, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Error;
    static deserializeBinaryFromReader(message: Error, reader: jspb.BinaryReader): Error;
  }

  export namespace Error {
    export type AsObject = {
      error?: google_rpc_status_pb.Status.AsObject,
    }
  }


  export class Complete extends jspb.Message {
    getError(): google_rpc_status_pb.Status | undefined;
    setError(value?: google_rpc_status_pb.Status): Complete;
    hasError(): boolean;
    clearError(): Complete;

    getResult(): Job.Result | undefined;
    setResult(value?: Job.Result): Complete;
    hasResult(): boolean;
    clearResult(): Complete;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Complete.AsObject;
    static toObject(includeInstance: boolean, msg: Complete): Complete.AsObject;
    static serializeBinaryToWriter(message: Complete, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Complete;
    static deserializeBinaryFromReader(message: Complete, reader: jspb.BinaryReader): Complete;
  }

  export namespace Complete {
    export type AsObject = {
      error?: google_rpc_status_pb.Status.AsObject,
      result?: Job.Result.AsObject,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    OPEN = 1,
    STATE = 2,
    JOB = 7,
    TERMINAL = 3,
    DOWNLOAD = 6,
    ERROR = 4,
    COMPLETE = 5,
  }
}

export class Runner extends jspb.Message {
  getId(): string;
  setId(value: string): Runner;

  getByIdOnly(): boolean;
  setByIdOnly(value: boolean): Runner;

  getOdr(): Runner.ODR | undefined;
  setOdr(value?: Runner.ODR): Runner;
  hasOdr(): boolean;
  clearOdr(): Runner;

  getLocal(): Runner.Local | undefined;
  setLocal(value?: Runner.Local): Runner;
  hasLocal(): boolean;
  clearLocal(): Runner;

  getRemote(): Runner.Remote | undefined;
  setRemote(value?: Runner.Remote): Runner;
  hasRemote(): boolean;
  clearRemote(): Runner;

  getDeprecatedIsOdr(): boolean;
  setDeprecatedIsOdr(value: boolean): Runner;

  getComponentsList(): Array<Component>;
  setComponentsList(value: Array<Component>): Runner;
  clearComponentsList(): Runner;
  addComponents(value?: Component, index?: number): Component;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Runner;

  getOnline(): boolean;
  setOnline(value: boolean): Runner;

  getFirstSeen(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setFirstSeen(value?: google_protobuf_timestamp_pb.Timestamp): Runner;
  hasFirstSeen(): boolean;
  clearFirstSeen(): Runner;

  getLastSeen(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastSeen(value?: google_protobuf_timestamp_pb.Timestamp): Runner;
  hasLastSeen(): boolean;
  clearLastSeen(): Runner;

  getAdoptionState(): Runner.AdoptionState;
  setAdoptionState(value: Runner.AdoptionState): Runner;

  getKindCase(): Runner.KindCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Runner.AsObject;
  static toObject(includeInstance: boolean, msg: Runner): Runner.AsObject;
  static serializeBinaryToWriter(message: Runner, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Runner;
  static deserializeBinaryFromReader(message: Runner, reader: jspb.BinaryReader): Runner;
}

export namespace Runner {
  export type AsObject = {
    id: string,
    byIdOnly: boolean,
    odr?: Runner.ODR.AsObject,
    local?: Runner.Local.AsObject,
    remote?: Runner.Remote.AsObject,
    deprecatedIsOdr: boolean,
    componentsList: Array<Component.AsObject>,
    labelsMap: Array<[string, string]>,
    online: boolean,
    firstSeen?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastSeen?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    adoptionState: Runner.AdoptionState,
  }

  export class ODR extends jspb.Message {
    getProfileId(): string;
    setProfileId(value: string): ODR;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ODR.AsObject;
    static toObject(includeInstance: boolean, msg: ODR): ODR.AsObject;
    static serializeBinaryToWriter(message: ODR, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ODR;
    static deserializeBinaryFromReader(message: ODR, reader: jspb.BinaryReader): ODR;
  }

  export namespace ODR {
    export type AsObject = {
      profileId: string,
    }
  }


  export class Local extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Local.AsObject;
    static toObject(includeInstance: boolean, msg: Local): Local.AsObject;
    static serializeBinaryToWriter(message: Local, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Local;
    static deserializeBinaryFromReader(message: Local, reader: jspb.BinaryReader): Local;
  }

  export namespace Local {
    export type AsObject = {
    }
  }


  export class Remote extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Remote.AsObject;
    static toObject(includeInstance: boolean, msg: Remote): Remote.AsObject;
    static serializeBinaryToWriter(message: Remote, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Remote;
    static deserializeBinaryFromReader(message: Remote, reader: jspb.BinaryReader): Remote;
  }

  export namespace Remote {
    export type AsObject = {
    }
  }


  export enum AdoptionState { 
    PENDING = 0,
    PREADOPTED = 1,
    ADOPTED = 2,
    REJECTED = 3,
  }

  export enum KindCase { 
    KIND_NOT_SET = 0,
    ODR = 5,
    LOCAL = 6,
    REMOTE = 7,
    DEPRECATED_IS_ODR = 4,
  }
}

export class RunnerTokenRequest extends jspb.Message {
  getRunner(): Runner | undefined;
  setRunner(value?: Runner): RunnerTokenRequest;
  hasRunner(): boolean;
  clearRunner(): RunnerTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerTokenRequest): RunnerTokenRequest.AsObject;
  static serializeBinaryToWriter(message: RunnerTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerTokenRequest;
  static deserializeBinaryFromReader(message: RunnerTokenRequest, reader: jspb.BinaryReader): RunnerTokenRequest;
}

export namespace RunnerTokenRequest {
  export type AsObject = {
    runner?: Runner.AsObject,
  }
}

export class RunnerTokenResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): RunnerTokenResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerTokenResponse): RunnerTokenResponse.AsObject;
  static serializeBinaryToWriter(message: RunnerTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerTokenResponse;
  static deserializeBinaryFromReader(message: RunnerTokenResponse, reader: jspb.BinaryReader): RunnerTokenResponse;
}

export namespace RunnerTokenResponse {
  export type AsObject = {
    token: string,
  }
}

export class RunnerConfigRequest extends jspb.Message {
  getOpen(): RunnerConfigRequest.Open | undefined;
  setOpen(value?: RunnerConfigRequest.Open): RunnerConfigRequest;
  hasOpen(): boolean;
  clearOpen(): RunnerConfigRequest;

  getEventCase(): RunnerConfigRequest.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerConfigRequest): RunnerConfigRequest.AsObject;
  static serializeBinaryToWriter(message: RunnerConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerConfigRequest;
  static deserializeBinaryFromReader(message: RunnerConfigRequest, reader: jspb.BinaryReader): RunnerConfigRequest;
}

export namespace RunnerConfigRequest {
  export type AsObject = {
    open?: RunnerConfigRequest.Open.AsObject,
  }

  export class Open extends jspb.Message {
    getRunner(): Runner | undefined;
    setRunner(value?: Runner): Open;
    hasRunner(): boolean;
    clearRunner(): Open;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Open.AsObject;
    static toObject(includeInstance: boolean, msg: Open): Open.AsObject;
    static serializeBinaryToWriter(message: Open, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Open;
    static deserializeBinaryFromReader(message: Open, reader: jspb.BinaryReader): Open;
  }

  export namespace Open {
    export type AsObject = {
      runner?: Runner.AsObject,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    OPEN = 1,
  }
}

export class RunnerConfigResponse extends jspb.Message {
  getConfig(): RunnerConfig | undefined;
  setConfig(value?: RunnerConfig): RunnerConfigResponse;
  hasConfig(): boolean;
  clearConfig(): RunnerConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerConfigResponse): RunnerConfigResponse.AsObject;
  static serializeBinaryToWriter(message: RunnerConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerConfigResponse;
  static deserializeBinaryFromReader(message: RunnerConfigResponse, reader: jspb.BinaryReader): RunnerConfigResponse;
}

export namespace RunnerConfigResponse {
  export type AsObject = {
    config?: RunnerConfig.AsObject,
  }
}

export class RunnerConfig extends jspb.Message {
  getConfigVarsList(): Array<ConfigVar>;
  setConfigVarsList(value: Array<ConfigVar>): RunnerConfig;
  clearConfigVarsList(): RunnerConfig;
  addConfigVars(value?: ConfigVar, index?: number): ConfigVar;

  getConfigSourcesList(): Array<ConfigSource>;
  setConfigSourcesList(value: Array<ConfigSource>): RunnerConfig;
  clearConfigSourcesList(): RunnerConfig;
  addConfigSources(value?: ConfigSource, index?: number): ConfigSource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerConfig.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerConfig): RunnerConfig.AsObject;
  static serializeBinaryToWriter(message: RunnerConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerConfig;
  static deserializeBinaryFromReader(message: RunnerConfig, reader: jspb.BinaryReader): RunnerConfig;
}

export namespace RunnerConfig {
  export type AsObject = {
    configVarsList: Array<ConfigVar.AsObject>,
    configSourcesList: Array<ConfigSource.AsObject>,
  }
}

export class RunnerJobStreamRequest extends jspb.Message {
  getRequest(): RunnerJobStreamRequest.Request | undefined;
  setRequest(value?: RunnerJobStreamRequest.Request): RunnerJobStreamRequest;
  hasRequest(): boolean;
  clearRequest(): RunnerJobStreamRequest;

  getAck(): RunnerJobStreamRequest.Ack | undefined;
  setAck(value?: RunnerJobStreamRequest.Ack): RunnerJobStreamRequest;
  hasAck(): boolean;
  clearAck(): RunnerJobStreamRequest;

  getComplete(): RunnerJobStreamRequest.Complete | undefined;
  setComplete(value?: RunnerJobStreamRequest.Complete): RunnerJobStreamRequest;
  hasComplete(): boolean;
  clearComplete(): RunnerJobStreamRequest;

  getError(): RunnerJobStreamRequest.Error | undefined;
  setError(value?: RunnerJobStreamRequest.Error): RunnerJobStreamRequest;
  hasError(): boolean;
  clearError(): RunnerJobStreamRequest;

  getTerminal(): GetJobStreamResponse.Terminal | undefined;
  setTerminal(value?: GetJobStreamResponse.Terminal): RunnerJobStreamRequest;
  hasTerminal(): boolean;
  clearTerminal(): RunnerJobStreamRequest;

  getDownload(): GetJobStreamResponse.Download | undefined;
  setDownload(value?: GetJobStreamResponse.Download): RunnerJobStreamRequest;
  hasDownload(): boolean;
  clearDownload(): RunnerJobStreamRequest;

  getConfigLoad(): RunnerJobStreamRequest.ConfigLoad | undefined;
  setConfigLoad(value?: RunnerJobStreamRequest.ConfigLoad): RunnerJobStreamRequest;
  hasConfigLoad(): boolean;
  clearConfigLoad(): RunnerJobStreamRequest;

  getVariableValuesSet(): RunnerJobStreamRequest.VariableValuesSet | undefined;
  setVariableValuesSet(value?: RunnerJobStreamRequest.VariableValuesSet): RunnerJobStreamRequest;
  hasVariableValuesSet(): boolean;
  clearVariableValuesSet(): RunnerJobStreamRequest;

  getHeartbeat(): RunnerJobStreamRequest.Heartbeat | undefined;
  setHeartbeat(value?: RunnerJobStreamRequest.Heartbeat): RunnerJobStreamRequest;
  hasHeartbeat(): boolean;
  clearHeartbeat(): RunnerJobStreamRequest;

  getEventCase(): RunnerJobStreamRequest.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerJobStreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerJobStreamRequest): RunnerJobStreamRequest.AsObject;
  static serializeBinaryToWriter(message: RunnerJobStreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerJobStreamRequest;
  static deserializeBinaryFromReader(message: RunnerJobStreamRequest, reader: jspb.BinaryReader): RunnerJobStreamRequest;
}

export namespace RunnerJobStreamRequest {
  export type AsObject = {
    request?: RunnerJobStreamRequest.Request.AsObject,
    ack?: RunnerJobStreamRequest.Ack.AsObject,
    complete?: RunnerJobStreamRequest.Complete.AsObject,
    error?: RunnerJobStreamRequest.Error.AsObject,
    terminal?: GetJobStreamResponse.Terminal.AsObject,
    download?: GetJobStreamResponse.Download.AsObject,
    configLoad?: RunnerJobStreamRequest.ConfigLoad.AsObject,
    variableValuesSet?: RunnerJobStreamRequest.VariableValuesSet.AsObject,
    heartbeat?: RunnerJobStreamRequest.Heartbeat.AsObject,
  }

  export class Request extends jspb.Message {
    getRunnerId(): string;
    setRunnerId(value: string): Request;

    getReattachJobId(): string;
    setReattachJobId(value: string): Request;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Request.AsObject;
    static toObject(includeInstance: boolean, msg: Request): Request.AsObject;
    static serializeBinaryToWriter(message: Request, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Request;
    static deserializeBinaryFromReader(message: Request, reader: jspb.BinaryReader): Request;
  }

  export namespace Request {
    export type AsObject = {
      runnerId: string,
      reattachJobId: string,
    }
  }


  export class Ack extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Ack.AsObject;
    static toObject(includeInstance: boolean, msg: Ack): Ack.AsObject;
    static serializeBinaryToWriter(message: Ack, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Ack;
    static deserializeBinaryFromReader(message: Ack, reader: jspb.BinaryReader): Ack;
  }

  export namespace Ack {
    export type AsObject = {
    }
  }


  export class Complete extends jspb.Message {
    getResult(): Job.Result | undefined;
    setResult(value?: Job.Result): Complete;
    hasResult(): boolean;
    clearResult(): Complete;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Complete.AsObject;
    static toObject(includeInstance: boolean, msg: Complete): Complete.AsObject;
    static serializeBinaryToWriter(message: Complete, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Complete;
    static deserializeBinaryFromReader(message: Complete, reader: jspb.BinaryReader): Complete;
  }

  export namespace Complete {
    export type AsObject = {
      result?: Job.Result.AsObject,
    }
  }


  export class Error extends jspb.Message {
    getError(): google_rpc_status_pb.Status | undefined;
    setError(value?: google_rpc_status_pb.Status): Error;
    hasError(): boolean;
    clearError(): Error;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Error.AsObject;
    static toObject(includeInstance: boolean, msg: Error): Error.AsObject;
    static serializeBinaryToWriter(message: Error, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Error;
    static deserializeBinaryFromReader(message: Error, reader: jspb.BinaryReader): Error;
  }

  export namespace Error {
    export type AsObject = {
      error?: google_rpc_status_pb.Status.AsObject,
    }
  }


  export class Heartbeat extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Heartbeat.AsObject;
    static toObject(includeInstance: boolean, msg: Heartbeat): Heartbeat.AsObject;
    static serializeBinaryToWriter(message: Heartbeat, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Heartbeat;
    static deserializeBinaryFromReader(message: Heartbeat, reader: jspb.BinaryReader): Heartbeat;
  }

  export namespace Heartbeat {
    export type AsObject = {
    }
  }


  export class ConfigLoad extends jspb.Message {
    getConfig(): Job.Config | undefined;
    setConfig(value?: Job.Config): ConfigLoad;
    hasConfig(): boolean;
    clearConfig(): ConfigLoad;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConfigLoad.AsObject;
    static toObject(includeInstance: boolean, msg: ConfigLoad): ConfigLoad.AsObject;
    static serializeBinaryToWriter(message: ConfigLoad, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConfigLoad;
    static deserializeBinaryFromReader(message: ConfigLoad, reader: jspb.BinaryReader): ConfigLoad;
  }

  export namespace ConfigLoad {
    export type AsObject = {
      config?: Job.Config.AsObject,
    }
  }


  export class VariableValuesSet extends jspb.Message {
    getFinalValuesMap(): jspb.Map<string, Variable.FinalValue>;
    clearFinalValuesMap(): VariableValuesSet;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): VariableValuesSet.AsObject;
    static toObject(includeInstance: boolean, msg: VariableValuesSet): VariableValuesSet.AsObject;
    static serializeBinaryToWriter(message: VariableValuesSet, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): VariableValuesSet;
    static deserializeBinaryFromReader(message: VariableValuesSet, reader: jspb.BinaryReader): VariableValuesSet;
  }

  export namespace VariableValuesSet {
    export type AsObject = {
      finalValuesMap: Array<[string, Variable.FinalValue.AsObject]>,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    REQUEST = 1,
    ACK = 2,
    COMPLETE = 3,
    ERROR = 4,
    TERMINAL = 5,
    DOWNLOAD = 7,
    CONFIG_LOAD = 8,
    VARIABLE_VALUES_SET = 9,
    HEARTBEAT = 6,
  }
}

export class RunnerJobStreamResponse extends jspb.Message {
  getAssignment(): RunnerJobStreamResponse.JobAssignment | undefined;
  setAssignment(value?: RunnerJobStreamResponse.JobAssignment): RunnerJobStreamResponse;
  hasAssignment(): boolean;
  clearAssignment(): RunnerJobStreamResponse;

  getCancel(): RunnerJobStreamResponse.JobCancel | undefined;
  setCancel(value?: RunnerJobStreamResponse.JobCancel): RunnerJobStreamResponse;
  hasCancel(): boolean;
  clearCancel(): RunnerJobStreamResponse;

  getEventCase(): RunnerJobStreamResponse.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerJobStreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerJobStreamResponse): RunnerJobStreamResponse.AsObject;
  static serializeBinaryToWriter(message: RunnerJobStreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerJobStreamResponse;
  static deserializeBinaryFromReader(message: RunnerJobStreamResponse, reader: jspb.BinaryReader): RunnerJobStreamResponse;
}

export namespace RunnerJobStreamResponse {
  export type AsObject = {
    assignment?: RunnerJobStreamResponse.JobAssignment.AsObject,
    cancel?: RunnerJobStreamResponse.JobCancel.AsObject,
  }

  export class JobAssignment extends jspb.Message {
    getJob(): Job | undefined;
    setJob(value?: Job): JobAssignment;
    hasJob(): boolean;
    clearJob(): JobAssignment;

    getConfigSourcesList(): Array<ConfigSource>;
    setConfigSourcesList(value: Array<ConfigSource>): JobAssignment;
    clearConfigSourcesList(): JobAssignment;
    addConfigSources(value?: ConfigSource, index?: number): ConfigSource;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): JobAssignment.AsObject;
    static toObject(includeInstance: boolean, msg: JobAssignment): JobAssignment.AsObject;
    static serializeBinaryToWriter(message: JobAssignment, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): JobAssignment;
    static deserializeBinaryFromReader(message: JobAssignment, reader: jspb.BinaryReader): JobAssignment;
  }

  export namespace JobAssignment {
    export type AsObject = {
      job?: Job.AsObject,
      configSourcesList: Array<ConfigSource.AsObject>,
    }
  }


  export class JobCancel extends jspb.Message {
    getForce(): boolean;
    setForce(value: boolean): JobCancel;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): JobCancel.AsObject;
    static toObject(includeInstance: boolean, msg: JobCancel): JobCancel.AsObject;
    static serializeBinaryToWriter(message: JobCancel, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): JobCancel;
    static deserializeBinaryFromReader(message: JobCancel, reader: jspb.BinaryReader): JobCancel;
  }

  export namespace JobCancel {
    export type AsObject = {
      force: boolean,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    ASSIGNMENT = 1,
    CANCEL = 2,
  }
}

export class RunnerGetDeploymentConfigRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerGetDeploymentConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerGetDeploymentConfigRequest): RunnerGetDeploymentConfigRequest.AsObject;
  static serializeBinaryToWriter(message: RunnerGetDeploymentConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerGetDeploymentConfigRequest;
  static deserializeBinaryFromReader(message: RunnerGetDeploymentConfigRequest, reader: jspb.BinaryReader): RunnerGetDeploymentConfigRequest;
}

export namespace RunnerGetDeploymentConfigRequest {
  export type AsObject = {
  }
}

export class RunnerGetDeploymentConfigResponse extends jspb.Message {
  getServerAddr(): string;
  setServerAddr(value: string): RunnerGetDeploymentConfigResponse;

  getServerTls(): boolean;
  setServerTls(value: boolean): RunnerGetDeploymentConfigResponse;

  getServerTlsSkipVerify(): boolean;
  setServerTlsSkipVerify(value: boolean): RunnerGetDeploymentConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunnerGetDeploymentConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RunnerGetDeploymentConfigResponse): RunnerGetDeploymentConfigResponse.AsObject;
  static serializeBinaryToWriter(message: RunnerGetDeploymentConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunnerGetDeploymentConfigResponse;
  static deserializeBinaryFromReader(message: RunnerGetDeploymentConfigResponse, reader: jspb.BinaryReader): RunnerGetDeploymentConfigResponse;
}

export namespace RunnerGetDeploymentConfigResponse {
  export type AsObject = {
    serverAddr: string,
    serverTls: boolean,
    serverTlsSkipVerify: boolean,
  }
}

export class GetRunnerRequest extends jspb.Message {
  getRunnerId(): string;
  setRunnerId(value: string): GetRunnerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRunnerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRunnerRequest): GetRunnerRequest.AsObject;
  static serializeBinaryToWriter(message: GetRunnerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRunnerRequest;
  static deserializeBinaryFromReader(message: GetRunnerRequest, reader: jspb.BinaryReader): GetRunnerRequest;
}

export namespace GetRunnerRequest {
  export type AsObject = {
    runnerId: string,
  }
}

export class ListRunnersRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListRunnersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListRunnersRequest): ListRunnersRequest.AsObject;
  static serializeBinaryToWriter(message: ListRunnersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRunnersRequest;
  static deserializeBinaryFromReader(message: ListRunnersRequest, reader: jspb.BinaryReader): ListRunnersRequest;
}

export namespace ListRunnersRequest {
  export type AsObject = {
  }
}

export class ListRunnersResponse extends jspb.Message {
  getRunnersList(): Array<Runner>;
  setRunnersList(value: Array<Runner>): ListRunnersResponse;
  clearRunnersList(): ListRunnersResponse;
  addRunners(value?: Runner, index?: number): Runner;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListRunnersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListRunnersResponse): ListRunnersResponse.AsObject;
  static serializeBinaryToWriter(message: ListRunnersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListRunnersResponse;
  static deserializeBinaryFromReader(message: ListRunnersResponse, reader: jspb.BinaryReader): ListRunnersResponse;
}

export namespace ListRunnersResponse {
  export type AsObject = {
    runnersList: Array<Runner.AsObject>,
  }
}

export class AdoptRunnerRequest extends jspb.Message {
  getRunnerId(): string;
  setRunnerId(value: string): AdoptRunnerRequest;

  getAdopt(): boolean;
  setAdopt(value: boolean): AdoptRunnerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AdoptRunnerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AdoptRunnerRequest): AdoptRunnerRequest.AsObject;
  static serializeBinaryToWriter(message: AdoptRunnerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AdoptRunnerRequest;
  static deserializeBinaryFromReader(message: AdoptRunnerRequest, reader: jspb.BinaryReader): AdoptRunnerRequest;
}

export namespace AdoptRunnerRequest {
  export type AsObject = {
    runnerId: string,
    adopt: boolean,
  }
}

export class ForgetRunnerRequest extends jspb.Message {
  getRunnerId(): string;
  setRunnerId(value: string): ForgetRunnerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ForgetRunnerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ForgetRunnerRequest): ForgetRunnerRequest.AsObject;
  static serializeBinaryToWriter(message: ForgetRunnerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ForgetRunnerRequest;
  static deserializeBinaryFromReader(message: ForgetRunnerRequest, reader: jspb.BinaryReader): ForgetRunnerRequest;
}

export namespace ForgetRunnerRequest {
  export type AsObject = {
    runnerId: string,
  }
}

export class SetServerConfigRequest extends jspb.Message {
  getConfig(): ServerConfig | undefined;
  setConfig(value?: ServerConfig): SetServerConfigRequest;
  hasConfig(): boolean;
  clearConfig(): SetServerConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetServerConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetServerConfigRequest): SetServerConfigRequest.AsObject;
  static serializeBinaryToWriter(message: SetServerConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetServerConfigRequest;
  static deserializeBinaryFromReader(message: SetServerConfigRequest, reader: jspb.BinaryReader): SetServerConfigRequest;
}

export namespace SetServerConfigRequest {
  export type AsObject = {
    config?: ServerConfig.AsObject,
  }
}

export class GetServerConfigResponse extends jspb.Message {
  getConfig(): ServerConfig | undefined;
  setConfig(value?: ServerConfig): GetServerConfigResponse;
  hasConfig(): boolean;
  clearConfig(): GetServerConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetServerConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetServerConfigResponse): GetServerConfigResponse.AsObject;
  static serializeBinaryToWriter(message: GetServerConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetServerConfigResponse;
  static deserializeBinaryFromReader(message: GetServerConfigResponse, reader: jspb.BinaryReader): GetServerConfigResponse;
}

export namespace GetServerConfigResponse {
  export type AsObject = {
    config?: ServerConfig.AsObject,
  }
}

export class ServerConfig extends jspb.Message {
  getAdvertiseAddrsList(): Array<ServerConfig.AdvertiseAddr>;
  setAdvertiseAddrsList(value: Array<ServerConfig.AdvertiseAddr>): ServerConfig;
  clearAdvertiseAddrsList(): ServerConfig;
  addAdvertiseAddrs(value?: ServerConfig.AdvertiseAddr, index?: number): ServerConfig.AdvertiseAddr;

  getPlatform(): string;
  setPlatform(value: string): ServerConfig;

  getCookie(): string;
  setCookie(value: string): ServerConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ServerConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ServerConfig): ServerConfig.AsObject;
  static serializeBinaryToWriter(message: ServerConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ServerConfig;
  static deserializeBinaryFromReader(message: ServerConfig, reader: jspb.BinaryReader): ServerConfig;
}

export namespace ServerConfig {
  export type AsObject = {
    advertiseAddrsList: Array<ServerConfig.AdvertiseAddr.AsObject>,
    platform: string,
    cookie: string,
  }

  export class AdvertiseAddr extends jspb.Message {
    getAddr(): string;
    setAddr(value: string): AdvertiseAddr;

    getTls(): boolean;
    setTls(value: boolean): AdvertiseAddr;

    getTlsSkipVerify(): boolean;
    setTlsSkipVerify(value: boolean): AdvertiseAddr;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AdvertiseAddr.AsObject;
    static toObject(includeInstance: boolean, msg: AdvertiseAddr): AdvertiseAddr.AsObject;
    static serializeBinaryToWriter(message: AdvertiseAddr, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AdvertiseAddr;
    static deserializeBinaryFromReader(message: AdvertiseAddr, reader: jspb.BinaryReader): AdvertiseAddr;
  }

  export namespace AdvertiseAddr {
    export type AsObject = {
      addr: string,
      tls: boolean,
      tlsSkipVerify: boolean,
    }
  }

}

export class CreateHostnameRequest extends jspb.Message {
  getHostname(): string;
  setHostname(value: string): CreateHostnameRequest;

  getTarget(): Hostname.Target | undefined;
  setTarget(value?: Hostname.Target): CreateHostnameRequest;
  hasTarget(): boolean;
  clearTarget(): CreateHostnameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateHostnameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateHostnameRequest): CreateHostnameRequest.AsObject;
  static serializeBinaryToWriter(message: CreateHostnameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateHostnameRequest;
  static deserializeBinaryFromReader(message: CreateHostnameRequest, reader: jspb.BinaryReader): CreateHostnameRequest;
}

export namespace CreateHostnameRequest {
  export type AsObject = {
    hostname: string,
    target?: Hostname.Target.AsObject,
  }
}

export class CreateHostnameResponse extends jspb.Message {
  getHostname(): Hostname | undefined;
  setHostname(value?: Hostname): CreateHostnameResponse;
  hasHostname(): boolean;
  clearHostname(): CreateHostnameResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateHostnameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateHostnameResponse): CreateHostnameResponse.AsObject;
  static serializeBinaryToWriter(message: CreateHostnameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateHostnameResponse;
  static deserializeBinaryFromReader(message: CreateHostnameResponse, reader: jspb.BinaryReader): CreateHostnameResponse;
}

export namespace CreateHostnameResponse {
  export type AsObject = {
    hostname?: Hostname.AsObject,
  }
}

export class ListHostnamesRequest extends jspb.Message {
  getTarget(): Hostname.Target | undefined;
  setTarget(value?: Hostname.Target): ListHostnamesRequest;
  hasTarget(): boolean;
  clearTarget(): ListHostnamesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHostnamesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListHostnamesRequest): ListHostnamesRequest.AsObject;
  static serializeBinaryToWriter(message: ListHostnamesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHostnamesRequest;
  static deserializeBinaryFromReader(message: ListHostnamesRequest, reader: jspb.BinaryReader): ListHostnamesRequest;
}

export namespace ListHostnamesRequest {
  export type AsObject = {
    target?: Hostname.Target.AsObject,
  }
}

export class ListHostnamesResponse extends jspb.Message {
  getHostnamesList(): Array<Hostname>;
  setHostnamesList(value: Array<Hostname>): ListHostnamesResponse;
  clearHostnamesList(): ListHostnamesResponse;
  addHostnames(value?: Hostname, index?: number): Hostname;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHostnamesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListHostnamesResponse): ListHostnamesResponse.AsObject;
  static serializeBinaryToWriter(message: ListHostnamesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHostnamesResponse;
  static deserializeBinaryFromReader(message: ListHostnamesResponse, reader: jspb.BinaryReader): ListHostnamesResponse;
}

export namespace ListHostnamesResponse {
  export type AsObject = {
    hostnamesList: Array<Hostname.AsObject>,
  }
}

export class DeleteHostnameRequest extends jspb.Message {
  getHostname(): string;
  setHostname(value: string): DeleteHostnameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteHostnameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteHostnameRequest): DeleteHostnameRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteHostnameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteHostnameRequest;
  static deserializeBinaryFromReader(message: DeleteHostnameRequest, reader: jspb.BinaryReader): DeleteHostnameRequest;
}

export namespace DeleteHostnameRequest {
  export type AsObject = {
    hostname: string,
  }
}

export class Hostname extends jspb.Message {
  getHostname(): string;
  setHostname(value: string): Hostname;

  getFqdn(): string;
  setFqdn(value: string): Hostname;

  getTargetLabelsMap(): jspb.Map<string, string>;
  clearTargetLabelsMap(): Hostname;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Hostname.AsObject;
  static toObject(includeInstance: boolean, msg: Hostname): Hostname.AsObject;
  static serializeBinaryToWriter(message: Hostname, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Hostname;
  static deserializeBinaryFromReader(message: Hostname, reader: jspb.BinaryReader): Hostname;
}

export namespace Hostname {
  export type AsObject = {
    hostname: string,
    fqdn: string,
    targetLabelsMap: Array<[string, string]>,
  }

  export class Target extends jspb.Message {
    getApplication(): Hostname.TargetApp | undefined;
    setApplication(value?: Hostname.TargetApp): Target;
    hasApplication(): boolean;
    clearApplication(): Target;

    getTargetCase(): Target.TargetCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Target.AsObject;
    static toObject(includeInstance: boolean, msg: Target): Target.AsObject;
    static serializeBinaryToWriter(message: Target, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Target;
    static deserializeBinaryFromReader(message: Target, reader: jspb.BinaryReader): Target;
  }

  export namespace Target {
    export type AsObject = {
      application?: Hostname.TargetApp.AsObject,
    }

    export enum TargetCase { 
      TARGET_NOT_SET = 0,
      APPLICATION = 20,
    }
  }


  export class TargetApp extends jspb.Message {
    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): TargetApp;
    hasApplication(): boolean;
    clearApplication(): TargetApp;

    getWorkspace(): Ref.Workspace | undefined;
    setWorkspace(value?: Ref.Workspace): TargetApp;
    hasWorkspace(): boolean;
    clearWorkspace(): TargetApp;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TargetApp.AsObject;
    static toObject(includeInstance: boolean, msg: TargetApp): TargetApp.AsObject;
    static serializeBinaryToWriter(message: TargetApp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TargetApp;
    static deserializeBinaryFromReader(message: TargetApp, reader: jspb.BinaryReader): TargetApp;
  }

  export namespace TargetApp {
    export type AsObject = {
      application?: Ref.Application.AsObject,
      workspace?: Ref.Workspace.AsObject,
    }
  }

}

export class ListWorkspacesRequest extends jspb.Message {
  getGlobal(): google_protobuf_empty_pb.Empty | undefined;
  setGlobal(value?: google_protobuf_empty_pb.Empty): ListWorkspacesRequest;
  hasGlobal(): boolean;
  clearGlobal(): ListWorkspacesRequest;

  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): ListWorkspacesRequest;
  hasProject(): boolean;
  clearProject(): ListWorkspacesRequest;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListWorkspacesRequest;
  hasApplication(): boolean;
  clearApplication(): ListWorkspacesRequest;

  getScopeCase(): ListWorkspacesRequest.ScopeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWorkspacesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListWorkspacesRequest): ListWorkspacesRequest.AsObject;
  static serializeBinaryToWriter(message: ListWorkspacesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWorkspacesRequest;
  static deserializeBinaryFromReader(message: ListWorkspacesRequest, reader: jspb.BinaryReader): ListWorkspacesRequest;
}

export namespace ListWorkspacesRequest {
  export type AsObject = {
    global?: google_protobuf_empty_pb.Empty.AsObject,
    project?: Ref.Project.AsObject,
    application?: Ref.Application.AsObject,
  }

  export enum ScopeCase { 
    SCOPE_NOT_SET = 0,
    GLOBAL = 1,
    PROJECT = 2,
    APPLICATION = 3,
  }
}

export class ListWorkspacesResponse extends jspb.Message {
  getWorkspacesList(): Array<Workspace>;
  setWorkspacesList(value: Array<Workspace>): ListWorkspacesResponse;
  clearWorkspacesList(): ListWorkspacesResponse;
  addWorkspaces(value?: Workspace, index?: number): Workspace;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListWorkspacesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListWorkspacesResponse): ListWorkspacesResponse.AsObject;
  static serializeBinaryToWriter(message: ListWorkspacesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListWorkspacesResponse;
  static deserializeBinaryFromReader(message: ListWorkspacesResponse, reader: jspb.BinaryReader): ListWorkspacesResponse;
}

export namespace ListWorkspacesResponse {
  export type AsObject = {
    workspacesList: Array<Workspace.AsObject>,
  }
}

export class GetWorkspaceRequest extends jspb.Message {
  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): GetWorkspaceRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): GetWorkspaceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetWorkspaceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetWorkspaceRequest): GetWorkspaceRequest.AsObject;
  static serializeBinaryToWriter(message: GetWorkspaceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetWorkspaceRequest;
  static deserializeBinaryFromReader(message: GetWorkspaceRequest, reader: jspb.BinaryReader): GetWorkspaceRequest;
}

export namespace GetWorkspaceRequest {
  export type AsObject = {
    workspace?: Ref.Workspace.AsObject,
  }
}

export class GetWorkspaceResponse extends jspb.Message {
  getWorkspace(): Workspace | undefined;
  setWorkspace(value?: Workspace): GetWorkspaceResponse;
  hasWorkspace(): boolean;
  clearWorkspace(): GetWorkspaceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetWorkspaceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetWorkspaceResponse): GetWorkspaceResponse.AsObject;
  static serializeBinaryToWriter(message: GetWorkspaceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetWorkspaceResponse;
  static deserializeBinaryFromReader(message: GetWorkspaceResponse, reader: jspb.BinaryReader): GetWorkspaceResponse;
}

export namespace GetWorkspaceResponse {
  export type AsObject = {
    workspace?: Workspace.AsObject,
  }
}

export class UpsertWorkspaceRequest extends jspb.Message {
  getWorkspace(): Workspace | undefined;
  setWorkspace(value?: Workspace): UpsertWorkspaceRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): UpsertWorkspaceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertWorkspaceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertWorkspaceRequest): UpsertWorkspaceRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertWorkspaceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertWorkspaceRequest;
  static deserializeBinaryFromReader(message: UpsertWorkspaceRequest, reader: jspb.BinaryReader): UpsertWorkspaceRequest;
}

export namespace UpsertWorkspaceRequest {
  export type AsObject = {
    workspace?: Workspace.AsObject,
  }
}

export class UpsertWorkspaceResponse extends jspb.Message {
  getWorkspace(): Workspace | undefined;
  setWorkspace(value?: Workspace): UpsertWorkspaceResponse;
  hasWorkspace(): boolean;
  clearWorkspace(): UpsertWorkspaceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertWorkspaceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertWorkspaceResponse): UpsertWorkspaceResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertWorkspaceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertWorkspaceResponse;
  static deserializeBinaryFromReader(message: UpsertWorkspaceResponse, reader: jspb.BinaryReader): UpsertWorkspaceResponse;
}

export namespace UpsertWorkspaceResponse {
  export type AsObject = {
    workspace?: Workspace.AsObject,
  }
}

export class UpsertProjectRequest extends jspb.Message {
  getProject(): Project | undefined;
  setProject(value?: Project): UpsertProjectRequest;
  hasProject(): boolean;
  clearProject(): UpsertProjectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertProjectRequest): UpsertProjectRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertProjectRequest;
  static deserializeBinaryFromReader(message: UpsertProjectRequest, reader: jspb.BinaryReader): UpsertProjectRequest;
}

export namespace UpsertProjectRequest {
  export type AsObject = {
    project?: Project.AsObject,
  }
}

export class UpsertProjectResponse extends jspb.Message {
  getProject(): Project | undefined;
  setProject(value?: Project): UpsertProjectResponse;
  hasProject(): boolean;
  clearProject(): UpsertProjectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertProjectResponse): UpsertProjectResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertProjectResponse;
  static deserializeBinaryFromReader(message: UpsertProjectResponse, reader: jspb.BinaryReader): UpsertProjectResponse;
}

export namespace UpsertProjectResponse {
  export type AsObject = {
    project?: Project.AsObject,
  }
}

export class GetProjectRequest extends jspb.Message {
  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): GetProjectRequest;
  hasProject(): boolean;
  clearProject(): GetProjectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectRequest): GetProjectRequest.AsObject;
  static serializeBinaryToWriter(message: GetProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectRequest;
  static deserializeBinaryFromReader(message: GetProjectRequest, reader: jspb.BinaryReader): GetProjectRequest;
}

export namespace GetProjectRequest {
  export type AsObject = {
    project?: Ref.Project.AsObject,
  }
}

export class GetProjectResponse extends jspb.Message {
  getProject(): Project | undefined;
  setProject(value?: Project): GetProjectResponse;
  hasProject(): boolean;
  clearProject(): GetProjectResponse;

  getWorkspacesList(): Array<Workspace.Project>;
  setWorkspacesList(value: Array<Workspace.Project>): GetProjectResponse;
  clearWorkspacesList(): GetProjectResponse;
  addWorkspaces(value?: Workspace.Project, index?: number): Workspace.Project;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectResponse): GetProjectResponse.AsObject;
  static serializeBinaryToWriter(message: GetProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectResponse;
  static deserializeBinaryFromReader(message: GetProjectResponse, reader: jspb.BinaryReader): GetProjectResponse;
}

export namespace GetProjectResponse {
  export type AsObject = {
    project?: Project.AsObject,
    workspacesList: Array<Workspace.Project.AsObject>,
  }
}

export class ListProjectsResponse extends jspb.Message {
  getProjectsList(): Array<Ref.Project>;
  setProjectsList(value: Array<Ref.Project>): ListProjectsResponse;
  clearProjectsList(): ListProjectsResponse;
  addProjects(value?: Ref.Project, index?: number): Ref.Project;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectsResponse): ListProjectsResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectsResponse;
  static deserializeBinaryFromReader(message: ListProjectsResponse, reader: jspb.BinaryReader): ListProjectsResponse;
}

export namespace ListProjectsResponse {
  export type AsObject = {
    projectsList: Array<Ref.Project.AsObject>,
  }
}

export class GetApplicationRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): GetApplicationRequest;
  hasApplication(): boolean;
  clearApplication(): GetApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationRequest): GetApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: GetApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationRequest;
  static deserializeBinaryFromReader(message: GetApplicationRequest, reader: jspb.BinaryReader): GetApplicationRequest;
}

export namespace GetApplicationRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
  }
}

export class GetApplicationResponse extends jspb.Message {
  getApplication(): Application | undefined;
  setApplication(value?: Application): GetApplicationResponse;
  hasApplication(): boolean;
  clearApplication(): GetApplicationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationResponse): GetApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: GetApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationResponse;
  static deserializeBinaryFromReader(message: GetApplicationResponse, reader: jspb.BinaryReader): GetApplicationResponse;
}

export namespace GetApplicationResponse {
  export type AsObject = {
    application?: Application.AsObject,
  }
}

export class UpsertApplicationRequest extends jspb.Message {
  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): UpsertApplicationRequest;
  hasProject(): boolean;
  clearProject(): UpsertApplicationRequest;

  getName(): string;
  setName(value: string): UpsertApplicationRequest;

  getFileChangeSignal(): string;
  setFileChangeSignal(value: string): UpsertApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertApplicationRequest): UpsertApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertApplicationRequest;
  static deserializeBinaryFromReader(message: UpsertApplicationRequest, reader: jspb.BinaryReader): UpsertApplicationRequest;
}

export namespace UpsertApplicationRequest {
  export type AsObject = {
    project?: Ref.Project.AsObject,
    name: string,
    fileChangeSignal: string,
  }
}

export class UpsertApplicationResponse extends jspb.Message {
  getApplication(): Application | undefined;
  setApplication(value?: Application): UpsertApplicationResponse;
  hasApplication(): boolean;
  clearApplication(): UpsertApplicationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertApplicationResponse): UpsertApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertApplicationResponse;
  static deserializeBinaryFromReader(message: UpsertApplicationResponse, reader: jspb.BinaryReader): UpsertApplicationResponse;
}

export namespace UpsertApplicationResponse {
  export type AsObject = {
    application?: Application.AsObject,
  }
}

export class UpsertBuildRequest extends jspb.Message {
  getBuild(): Build | undefined;
  setBuild(value?: Build): UpsertBuildRequest;
  hasBuild(): boolean;
  clearBuild(): UpsertBuildRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertBuildRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertBuildRequest): UpsertBuildRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertBuildRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertBuildRequest;
  static deserializeBinaryFromReader(message: UpsertBuildRequest, reader: jspb.BinaryReader): UpsertBuildRequest;
}

export namespace UpsertBuildRequest {
  export type AsObject = {
    build?: Build.AsObject,
  }
}

export class UpsertBuildResponse extends jspb.Message {
  getBuild(): Build | undefined;
  setBuild(value?: Build): UpsertBuildResponse;
  hasBuild(): boolean;
  clearBuild(): UpsertBuildResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertBuildResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertBuildResponse): UpsertBuildResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertBuildResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertBuildResponse;
  static deserializeBinaryFromReader(message: UpsertBuildResponse, reader: jspb.BinaryReader): UpsertBuildResponse;
}

export namespace UpsertBuildResponse {
  export type AsObject = {
    build?: Build.AsObject,
  }
}

export class ListBuildsRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListBuildsRequest;
  hasApplication(): boolean;
  clearApplication(): ListBuildsRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ListBuildsRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ListBuildsRequest;

  getOrder(): OperationOrder | undefined;
  setOrder(value?: OperationOrder): ListBuildsRequest;
  hasOrder(): boolean;
  clearOrder(): ListBuildsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListBuildsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListBuildsRequest): ListBuildsRequest.AsObject;
  static serializeBinaryToWriter(message: ListBuildsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListBuildsRequest;
  static deserializeBinaryFromReader(message: ListBuildsRequest, reader: jspb.BinaryReader): ListBuildsRequest;
}

export namespace ListBuildsRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    order?: OperationOrder.AsObject,
  }
}

export class ListBuildsResponse extends jspb.Message {
  getBuildsList(): Array<Build>;
  setBuildsList(value: Array<Build>): ListBuildsResponse;
  clearBuildsList(): ListBuildsResponse;
  addBuilds(value?: Build, index?: number): Build;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListBuildsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListBuildsResponse): ListBuildsResponse.AsObject;
  static serializeBinaryToWriter(message: ListBuildsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListBuildsResponse;
  static deserializeBinaryFromReader(message: ListBuildsResponse, reader: jspb.BinaryReader): ListBuildsResponse;
}

export namespace ListBuildsResponse {
  export type AsObject = {
    buildsList: Array<Build.AsObject>,
  }
}

export class GetLatestBuildRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): GetLatestBuildRequest;
  hasApplication(): boolean;
  clearApplication(): GetLatestBuildRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): GetLatestBuildRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): GetLatestBuildRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLatestBuildRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLatestBuildRequest): GetLatestBuildRequest.AsObject;
  static serializeBinaryToWriter(message: GetLatestBuildRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLatestBuildRequest;
  static deserializeBinaryFromReader(message: GetLatestBuildRequest, reader: jspb.BinaryReader): GetLatestBuildRequest;
}

export namespace GetLatestBuildRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
  }
}

export class GetBuildRequest extends jspb.Message {
  getRef(): Ref.Operation | undefined;
  setRef(value?: Ref.Operation): GetBuildRequest;
  hasRef(): boolean;
  clearRef(): GetBuildRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBuildRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetBuildRequest): GetBuildRequest.AsObject;
  static serializeBinaryToWriter(message: GetBuildRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBuildRequest;
  static deserializeBinaryFromReader(message: GetBuildRequest, reader: jspb.BinaryReader): GetBuildRequest;
}

export namespace GetBuildRequest {
  export type AsObject = {
    ref?: Ref.Operation.AsObject,
  }
}

export class Build extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): Build;
  hasApplication(): boolean;
  clearApplication(): Build;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): Build;
  hasWorkspace(): boolean;
  clearWorkspace(): Build;

  getSequence(): number;
  setSequence(value: number): Build;

  getId(): string;
  setId(value: string): Build;

  getStatus(): Status | undefined;
  setStatus(value?: Status): Build;
  hasStatus(): boolean;
  clearStatus(): Build;

  getComponent(): Component | undefined;
  setComponent(value?: Component): Build;
  hasComponent(): boolean;
  clearComponent(): Build;

  getArtifact(): Artifact | undefined;
  setArtifact(value?: Artifact): Build;
  hasArtifact(): boolean;
  clearArtifact(): Build;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Build;

  getTemplateData(): Uint8Array | string;
  getTemplateData_asU8(): Uint8Array;
  getTemplateData_asB64(): string;
  setTemplateData(value: Uint8Array | string): Build;

  getJobId(): string;
  setJobId(value: string): Build;

  getPreload(): Build.Preload | undefined;
  setPreload(value?: Build.Preload): Build;
  hasPreload(): boolean;
  clearPreload(): Build;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Build.AsObject;
  static toObject(includeInstance: boolean, msg: Build): Build.AsObject;
  static serializeBinaryToWriter(message: Build, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Build;
  static deserializeBinaryFromReader(message: Build, reader: jspb.BinaryReader): Build;
}

export namespace Build {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    sequence: number,
    id: string,
    status?: Status.AsObject,
    component?: Component.AsObject,
    artifact?: Artifact.AsObject,
    labelsMap: Array<[string, string]>,
    templateData: Uint8Array | string,
    jobId: string,
    preload?: Build.Preload.AsObject,
  }

  export class Preload extends jspb.Message {
    getJobDataSourceRef(): Job.DataSource.Ref | undefined;
    setJobDataSourceRef(value?: Job.DataSource.Ref): Preload;
    hasJobDataSourceRef(): boolean;
    clearJobDataSourceRef(): Preload;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Preload.AsObject;
    static toObject(includeInstance: boolean, msg: Preload): Preload.AsObject;
    static serializeBinaryToWriter(message: Preload, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Preload;
    static deserializeBinaryFromReader(message: Preload, reader: jspb.BinaryReader): Preload;
  }

  export namespace Preload {
    export type AsObject = {
      jobDataSourceRef?: Job.DataSource.Ref.AsObject,
    }
  }

}

export class Artifact extends jspb.Message {
  getArtifact(): any_pb.Any | undefined;
  setArtifact(value?: any_pb.Any): Artifact;
  hasArtifact(): boolean;
  clearArtifact(): Artifact;

  getArtifactJson(): string;
  setArtifactJson(value: string): Artifact;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Artifact.AsObject;
  static toObject(includeInstance: boolean, msg: Artifact): Artifact.AsObject;
  static serializeBinaryToWriter(message: Artifact, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Artifact;
  static deserializeBinaryFromReader(message: Artifact, reader: jspb.BinaryReader): Artifact;
}

export namespace Artifact {
  export type AsObject = {
    artifact?: any_pb.Any.AsObject,
    artifactJson: string,
  }
}

export class OnDemandRunnerConfig extends jspb.Message {
  getId(): string;
  setId(value: string): OnDemandRunnerConfig;

  getName(): string;
  setName(value: string): OnDemandRunnerConfig;

  getTargetRunner(): Ref.Runner | undefined;
  setTargetRunner(value?: Ref.Runner): OnDemandRunnerConfig;
  hasTargetRunner(): boolean;
  clearTargetRunner(): OnDemandRunnerConfig;

  getOciUrl(): string;
  setOciUrl(value: string): OnDemandRunnerConfig;

  getEnvironmentVariablesMap(): jspb.Map<string, string>;
  clearEnvironmentVariablesMap(): OnDemandRunnerConfig;

  getPluginType(): string;
  setPluginType(value: string): OnDemandRunnerConfig;

  getPluginConfig(): Uint8Array | string;
  getPluginConfig_asU8(): Uint8Array;
  getPluginConfig_asB64(): string;
  setPluginConfig(value: Uint8Array | string): OnDemandRunnerConfig;

  getConfigFormat(): Hcl.Format;
  setConfigFormat(value: Hcl.Format): OnDemandRunnerConfig;

  getDefault(): boolean;
  setDefault(value: boolean): OnDemandRunnerConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OnDemandRunnerConfig.AsObject;
  static toObject(includeInstance: boolean, msg: OnDemandRunnerConfig): OnDemandRunnerConfig.AsObject;
  static serializeBinaryToWriter(message: OnDemandRunnerConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OnDemandRunnerConfig;
  static deserializeBinaryFromReader(message: OnDemandRunnerConfig, reader: jspb.BinaryReader): OnDemandRunnerConfig;
}

export namespace OnDemandRunnerConfig {
  export type AsObject = {
    id: string,
    name: string,
    targetRunner?: Ref.Runner.AsObject,
    ociUrl: string,
    environmentVariablesMap: Array<[string, string]>,
    pluginType: string,
    pluginConfig: Uint8Array | string,
    configFormat: Hcl.Format,
    pb_default: boolean,
  }
}

export class UpsertOnDemandRunnerConfigRequest extends jspb.Message {
  getConfig(): OnDemandRunnerConfig | undefined;
  setConfig(value?: OnDemandRunnerConfig): UpsertOnDemandRunnerConfigRequest;
  hasConfig(): boolean;
  clearConfig(): UpsertOnDemandRunnerConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertOnDemandRunnerConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertOnDemandRunnerConfigRequest): UpsertOnDemandRunnerConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertOnDemandRunnerConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertOnDemandRunnerConfigRequest;
  static deserializeBinaryFromReader(message: UpsertOnDemandRunnerConfigRequest, reader: jspb.BinaryReader): UpsertOnDemandRunnerConfigRequest;
}

export namespace UpsertOnDemandRunnerConfigRequest {
  export type AsObject = {
    config?: OnDemandRunnerConfig.AsObject,
  }
}

export class UpsertOnDemandRunnerConfigResponse extends jspb.Message {
  getConfig(): OnDemandRunnerConfig | undefined;
  setConfig(value?: OnDemandRunnerConfig): UpsertOnDemandRunnerConfigResponse;
  hasConfig(): boolean;
  clearConfig(): UpsertOnDemandRunnerConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertOnDemandRunnerConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertOnDemandRunnerConfigResponse): UpsertOnDemandRunnerConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertOnDemandRunnerConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertOnDemandRunnerConfigResponse;
  static deserializeBinaryFromReader(message: UpsertOnDemandRunnerConfigResponse, reader: jspb.BinaryReader): UpsertOnDemandRunnerConfigResponse;
}

export namespace UpsertOnDemandRunnerConfigResponse {
  export type AsObject = {
    config?: OnDemandRunnerConfig.AsObject,
  }
}

export class GetOnDemandRunnerConfigRequest extends jspb.Message {
  getConfig(): Ref.OnDemandRunnerConfig | undefined;
  setConfig(value?: Ref.OnDemandRunnerConfig): GetOnDemandRunnerConfigRequest;
  hasConfig(): boolean;
  clearConfig(): GetOnDemandRunnerConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOnDemandRunnerConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOnDemandRunnerConfigRequest): GetOnDemandRunnerConfigRequest.AsObject;
  static serializeBinaryToWriter(message: GetOnDemandRunnerConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOnDemandRunnerConfigRequest;
  static deserializeBinaryFromReader(message: GetOnDemandRunnerConfigRequest, reader: jspb.BinaryReader): GetOnDemandRunnerConfigRequest;
}

export namespace GetOnDemandRunnerConfigRequest {
  export type AsObject = {
    config?: Ref.OnDemandRunnerConfig.AsObject,
  }
}

export class GetOnDemandRunnerConfigResponse extends jspb.Message {
  getConfig(): OnDemandRunnerConfig | undefined;
  setConfig(value?: OnDemandRunnerConfig): GetOnDemandRunnerConfigResponse;
  hasConfig(): boolean;
  clearConfig(): GetOnDemandRunnerConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOnDemandRunnerConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOnDemandRunnerConfigResponse): GetOnDemandRunnerConfigResponse.AsObject;
  static serializeBinaryToWriter(message: GetOnDemandRunnerConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOnDemandRunnerConfigResponse;
  static deserializeBinaryFromReader(message: GetOnDemandRunnerConfigResponse, reader: jspb.BinaryReader): GetOnDemandRunnerConfigResponse;
}

export namespace GetOnDemandRunnerConfigResponse {
  export type AsObject = {
    config?: OnDemandRunnerConfig.AsObject,
  }
}

export class DeleteOnDemandRunnerConfigRequest extends jspb.Message {
  getConfig(): Ref.OnDemandRunnerConfig | undefined;
  setConfig(value?: Ref.OnDemandRunnerConfig): DeleteOnDemandRunnerConfigRequest;
  hasConfig(): boolean;
  clearConfig(): DeleteOnDemandRunnerConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOnDemandRunnerConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOnDemandRunnerConfigRequest): DeleteOnDemandRunnerConfigRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteOnDemandRunnerConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOnDemandRunnerConfigRequest;
  static deserializeBinaryFromReader(message: DeleteOnDemandRunnerConfigRequest, reader: jspb.BinaryReader): DeleteOnDemandRunnerConfigRequest;
}

export namespace DeleteOnDemandRunnerConfigRequest {
  export type AsObject = {
    config?: Ref.OnDemandRunnerConfig.AsObject,
  }
}

export class DeleteOnDemandRunnerConfigResponse extends jspb.Message {
  getConfig(): OnDemandRunnerConfig | undefined;
  setConfig(value?: OnDemandRunnerConfig): DeleteOnDemandRunnerConfigResponse;
  hasConfig(): boolean;
  clearConfig(): DeleteOnDemandRunnerConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOnDemandRunnerConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOnDemandRunnerConfigResponse): DeleteOnDemandRunnerConfigResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteOnDemandRunnerConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOnDemandRunnerConfigResponse;
  static deserializeBinaryFromReader(message: DeleteOnDemandRunnerConfigResponse, reader: jspb.BinaryReader): DeleteOnDemandRunnerConfigResponse;
}

export namespace DeleteOnDemandRunnerConfigResponse {
  export type AsObject = {
    config?: OnDemandRunnerConfig.AsObject,
  }
}

export class ListOnDemandRunnerConfigsResponse extends jspb.Message {
  getConfigsList(): Array<OnDemandRunnerConfig>;
  setConfigsList(value: Array<OnDemandRunnerConfig>): ListOnDemandRunnerConfigsResponse;
  clearConfigsList(): ListOnDemandRunnerConfigsResponse;
  addConfigs(value?: OnDemandRunnerConfig, index?: number): OnDemandRunnerConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOnDemandRunnerConfigsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOnDemandRunnerConfigsResponse): ListOnDemandRunnerConfigsResponse.AsObject;
  static serializeBinaryToWriter(message: ListOnDemandRunnerConfigsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOnDemandRunnerConfigsResponse;
  static deserializeBinaryFromReader(message: ListOnDemandRunnerConfigsResponse, reader: jspb.BinaryReader): ListOnDemandRunnerConfigsResponse;
}

export namespace ListOnDemandRunnerConfigsResponse {
  export type AsObject = {
    configsList: Array<OnDemandRunnerConfig.AsObject>,
  }
}

export class UpsertPushedArtifactRequest extends jspb.Message {
  getArtifact(): PushedArtifact | undefined;
  setArtifact(value?: PushedArtifact): UpsertPushedArtifactRequest;
  hasArtifact(): boolean;
  clearArtifact(): UpsertPushedArtifactRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertPushedArtifactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertPushedArtifactRequest): UpsertPushedArtifactRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertPushedArtifactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertPushedArtifactRequest;
  static deserializeBinaryFromReader(message: UpsertPushedArtifactRequest, reader: jspb.BinaryReader): UpsertPushedArtifactRequest;
}

export namespace UpsertPushedArtifactRequest {
  export type AsObject = {
    artifact?: PushedArtifact.AsObject,
  }
}

export class UpsertPushedArtifactResponse extends jspb.Message {
  getArtifact(): PushedArtifact | undefined;
  setArtifact(value?: PushedArtifact): UpsertPushedArtifactResponse;
  hasArtifact(): boolean;
  clearArtifact(): UpsertPushedArtifactResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertPushedArtifactResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertPushedArtifactResponse): UpsertPushedArtifactResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertPushedArtifactResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertPushedArtifactResponse;
  static deserializeBinaryFromReader(message: UpsertPushedArtifactResponse, reader: jspb.BinaryReader): UpsertPushedArtifactResponse;
}

export namespace UpsertPushedArtifactResponse {
  export type AsObject = {
    artifact?: PushedArtifact.AsObject,
  }
}

export class GetLatestPushedArtifactRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): GetLatestPushedArtifactRequest;
  hasApplication(): boolean;
  clearApplication(): GetLatestPushedArtifactRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): GetLatestPushedArtifactRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): GetLatestPushedArtifactRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLatestPushedArtifactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLatestPushedArtifactRequest): GetLatestPushedArtifactRequest.AsObject;
  static serializeBinaryToWriter(message: GetLatestPushedArtifactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLatestPushedArtifactRequest;
  static deserializeBinaryFromReader(message: GetLatestPushedArtifactRequest, reader: jspb.BinaryReader): GetLatestPushedArtifactRequest;
}

export namespace GetLatestPushedArtifactRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
  }
}

export class GetPushedArtifactRequest extends jspb.Message {
  getRef(): Ref.Operation | undefined;
  setRef(value?: Ref.Operation): GetPushedArtifactRequest;
  hasRef(): boolean;
  clearRef(): GetPushedArtifactRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPushedArtifactRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPushedArtifactRequest): GetPushedArtifactRequest.AsObject;
  static serializeBinaryToWriter(message: GetPushedArtifactRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPushedArtifactRequest;
  static deserializeBinaryFromReader(message: GetPushedArtifactRequest, reader: jspb.BinaryReader): GetPushedArtifactRequest;
}

export namespace GetPushedArtifactRequest {
  export type AsObject = {
    ref?: Ref.Operation.AsObject,
  }
}

export class ListPushedArtifactsRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListPushedArtifactsRequest;
  hasApplication(): boolean;
  clearApplication(): ListPushedArtifactsRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ListPushedArtifactsRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ListPushedArtifactsRequest;

  getStatusList(): Array<StatusFilter>;
  setStatusList(value: Array<StatusFilter>): ListPushedArtifactsRequest;
  clearStatusList(): ListPushedArtifactsRequest;
  addStatus(value?: StatusFilter, index?: number): StatusFilter;

  getOrder(): OperationOrder | undefined;
  setOrder(value?: OperationOrder): ListPushedArtifactsRequest;
  hasOrder(): boolean;
  clearOrder(): ListPushedArtifactsRequest;

  getIncludeBuild(): boolean;
  setIncludeBuild(value: boolean): ListPushedArtifactsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPushedArtifactsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListPushedArtifactsRequest): ListPushedArtifactsRequest.AsObject;
  static serializeBinaryToWriter(message: ListPushedArtifactsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPushedArtifactsRequest;
  static deserializeBinaryFromReader(message: ListPushedArtifactsRequest, reader: jspb.BinaryReader): ListPushedArtifactsRequest;
}

export namespace ListPushedArtifactsRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    statusList: Array<StatusFilter.AsObject>,
    order?: OperationOrder.AsObject,
    includeBuild: boolean,
  }
}

export class ListPushedArtifactsResponse extends jspb.Message {
  getArtifactsList(): Array<PushedArtifact>;
  setArtifactsList(value: Array<PushedArtifact>): ListPushedArtifactsResponse;
  clearArtifactsList(): ListPushedArtifactsResponse;
  addArtifacts(value?: PushedArtifact, index?: number): PushedArtifact;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPushedArtifactsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListPushedArtifactsResponse): ListPushedArtifactsResponse.AsObject;
  static serializeBinaryToWriter(message: ListPushedArtifactsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPushedArtifactsResponse;
  static deserializeBinaryFromReader(message: ListPushedArtifactsResponse, reader: jspb.BinaryReader): ListPushedArtifactsResponse;
}

export namespace ListPushedArtifactsResponse {
  export type AsObject = {
    artifactsList: Array<PushedArtifact.AsObject>,
  }
}

export class PushedArtifact extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): PushedArtifact;
  hasApplication(): boolean;
  clearApplication(): PushedArtifact;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): PushedArtifact;
  hasWorkspace(): boolean;
  clearWorkspace(): PushedArtifact;

  getSequence(): number;
  setSequence(value: number): PushedArtifact;

  getId(): string;
  setId(value: string): PushedArtifact;

  getStatus(): Status | undefined;
  setStatus(value?: Status): PushedArtifact;
  hasStatus(): boolean;
  clearStatus(): PushedArtifact;

  getComponent(): Component | undefined;
  setComponent(value?: Component): PushedArtifact;
  hasComponent(): boolean;
  clearComponent(): PushedArtifact;

  getArtifact(): Artifact | undefined;
  setArtifact(value?: Artifact): PushedArtifact;
  hasArtifact(): boolean;
  clearArtifact(): PushedArtifact;

  getBuildId(): string;
  setBuildId(value: string): PushedArtifact;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): PushedArtifact;

  getTemplateData(): Uint8Array | string;
  getTemplateData_asU8(): Uint8Array;
  getTemplateData_asB64(): string;
  setTemplateData(value: Uint8Array | string): PushedArtifact;

  getBuild(): Build | undefined;
  setBuild(value?: Build): PushedArtifact;
  hasBuild(): boolean;
  clearBuild(): PushedArtifact;

  getJobId(): string;
  setJobId(value: string): PushedArtifact;

  getPreload(): PushedArtifact.Preload | undefined;
  setPreload(value?: PushedArtifact.Preload): PushedArtifact;
  hasPreload(): boolean;
  clearPreload(): PushedArtifact;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PushedArtifact.AsObject;
  static toObject(includeInstance: boolean, msg: PushedArtifact): PushedArtifact.AsObject;
  static serializeBinaryToWriter(message: PushedArtifact, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PushedArtifact;
  static deserializeBinaryFromReader(message: PushedArtifact, reader: jspb.BinaryReader): PushedArtifact;
}

export namespace PushedArtifact {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    sequence: number,
    id: string,
    status?: Status.AsObject,
    component?: Component.AsObject,
    artifact?: Artifact.AsObject,
    buildId: string,
    labelsMap: Array<[string, string]>,
    templateData: Uint8Array | string,
    build?: Build.AsObject,
    jobId: string,
    preload?: PushedArtifact.Preload.AsObject,
  }

  export class Preload extends jspb.Message {
    getJobDataSourceRef(): Job.DataSource.Ref | undefined;
    setJobDataSourceRef(value?: Job.DataSource.Ref): Preload;
    hasJobDataSourceRef(): boolean;
    clearJobDataSourceRef(): Preload;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Preload.AsObject;
    static toObject(includeInstance: boolean, msg: Preload): Preload.AsObject;
    static serializeBinaryToWriter(message: Preload, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Preload;
    static deserializeBinaryFromReader(message: Preload, reader: jspb.BinaryReader): Preload;
  }

  export namespace Preload {
    export type AsObject = {
      jobDataSourceRef?: Job.DataSource.Ref.AsObject,
    }
  }

}

export class GetDeploymentRequest extends jspb.Message {
  getRef(): Ref.Operation | undefined;
  setRef(value?: Ref.Operation): GetDeploymentRequest;
  hasRef(): boolean;
  clearRef(): GetDeploymentRequest;

  getLoadDetails(): Deployment.LoadDetails;
  setLoadDetails(value: Deployment.LoadDetails): GetDeploymentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDeploymentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDeploymentRequest): GetDeploymentRequest.AsObject;
  static serializeBinaryToWriter(message: GetDeploymentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDeploymentRequest;
  static deserializeBinaryFromReader(message: GetDeploymentRequest, reader: jspb.BinaryReader): GetDeploymentRequest;
}

export namespace GetDeploymentRequest {
  export type AsObject = {
    ref?: Ref.Operation.AsObject,
    loadDetails: Deployment.LoadDetails,
  }
}

export class UpsertDeploymentRequest extends jspb.Message {
  getDeployment(): Deployment | undefined;
  setDeployment(value?: Deployment): UpsertDeploymentRequest;
  hasDeployment(): boolean;
  clearDeployment(): UpsertDeploymentRequest;

  getAutoHostname(): UpsertDeploymentRequest.Tristate;
  setAutoHostname(value: UpsertDeploymentRequest.Tristate): UpsertDeploymentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertDeploymentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertDeploymentRequest): UpsertDeploymentRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertDeploymentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertDeploymentRequest;
  static deserializeBinaryFromReader(message: UpsertDeploymentRequest, reader: jspb.BinaryReader): UpsertDeploymentRequest;
}

export namespace UpsertDeploymentRequest {
  export type AsObject = {
    deployment?: Deployment.AsObject,
    autoHostname: UpsertDeploymentRequest.Tristate,
  }

  export enum Tristate { 
    UNSET = 0,
    TRUE = 1,
    FALSE = 2,
  }
}

export class UpsertDeploymentResponse extends jspb.Message {
  getDeployment(): Deployment | undefined;
  setDeployment(value?: Deployment): UpsertDeploymentResponse;
  hasDeployment(): boolean;
  clearDeployment(): UpsertDeploymentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertDeploymentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertDeploymentResponse): UpsertDeploymentResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertDeploymentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertDeploymentResponse;
  static deserializeBinaryFromReader(message: UpsertDeploymentResponse, reader: jspb.BinaryReader): UpsertDeploymentResponse;
}

export namespace UpsertDeploymentResponse {
  export type AsObject = {
    deployment?: Deployment.AsObject,
  }
}

export class ListDeploymentsRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListDeploymentsRequest;
  hasApplication(): boolean;
  clearApplication(): ListDeploymentsRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ListDeploymentsRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ListDeploymentsRequest;

  getStatusList(): Array<StatusFilter>;
  setStatusList(value: Array<StatusFilter>): ListDeploymentsRequest;
  clearStatusList(): ListDeploymentsRequest;
  addStatus(value?: StatusFilter, index?: number): StatusFilter;

  getPhysicalState(): Operation.PhysicalState;
  setPhysicalState(value: Operation.PhysicalState): ListDeploymentsRequest;

  getOrder(): OperationOrder | undefined;
  setOrder(value?: OperationOrder): ListDeploymentsRequest;
  hasOrder(): boolean;
  clearOrder(): ListDeploymentsRequest;

  getLoadDetails(): Deployment.LoadDetails;
  setLoadDetails(value: Deployment.LoadDetails): ListDeploymentsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeploymentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeploymentsRequest): ListDeploymentsRequest.AsObject;
  static serializeBinaryToWriter(message: ListDeploymentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeploymentsRequest;
  static deserializeBinaryFromReader(message: ListDeploymentsRequest, reader: jspb.BinaryReader): ListDeploymentsRequest;
}

export namespace ListDeploymentsRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    statusList: Array<StatusFilter.AsObject>,
    physicalState: Operation.PhysicalState,
    order?: OperationOrder.AsObject,
    loadDetails: Deployment.LoadDetails,
  }
}

export class ListDeploymentsResponse extends jspb.Message {
  getDeploymentsList(): Array<Deployment>;
  setDeploymentsList(value: Array<Deployment>): ListDeploymentsResponse;
  clearDeploymentsList(): ListDeploymentsResponse;
  addDeployments(value?: Deployment, index?: number): Deployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeploymentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeploymentsResponse): ListDeploymentsResponse.AsObject;
  static serializeBinaryToWriter(message: ListDeploymentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeploymentsResponse;
  static deserializeBinaryFromReader(message: ListDeploymentsResponse, reader: jspb.BinaryReader): ListDeploymentsResponse;
}

export namespace ListDeploymentsResponse {
  export type AsObject = {
    deploymentsList: Array<Deployment.AsObject>,
  }
}

export class Deployment extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): Deployment;
  hasApplication(): boolean;
  clearApplication(): Deployment;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): Deployment;
  hasWorkspace(): boolean;
  clearWorkspace(): Deployment;

  getSequence(): number;
  setSequence(value: number): Deployment;

  getId(): string;
  setId(value: string): Deployment;

  getUrl(): string;
  setUrl(value: string): Deployment;

  getGeneration(): Generation | undefined;
  setGeneration(value?: Generation): Deployment;
  hasGeneration(): boolean;
  clearGeneration(): Deployment;

  getState(): Operation.PhysicalState;
  setState(value: Operation.PhysicalState): Deployment;

  getStatus(): Status | undefined;
  setStatus(value?: Status): Deployment;
  hasStatus(): boolean;
  clearStatus(): Deployment;

  getComponent(): Component | undefined;
  setComponent(value?: Component): Deployment;
  hasComponent(): boolean;
  clearComponent(): Deployment;

  getArtifactId(): string;
  setArtifactId(value: string): Deployment;

  getDeployment(): any_pb.Any | undefined;
  setDeployment(value?: any_pb.Any): Deployment;
  hasDeployment(): boolean;
  clearDeployment(): Deployment;

  getDeploymentJson(): string;
  setDeploymentJson(value: string): Deployment;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Deployment;

  getTemplateData(): Uint8Array | string;
  getTemplateData_asU8(): Uint8Array;
  getTemplateData_asB64(): string;
  setTemplateData(value: Uint8Array | string): Deployment;

  getJobId(): string;
  setJobId(value: string): Deployment;

  getHasEntrypointConfig(): boolean;
  setHasEntrypointConfig(value: boolean): Deployment;

  getHasExecPlugin(): boolean;
  setHasExecPlugin(value: boolean): Deployment;

  getHasLogsPlugin(): boolean;
  setHasLogsPlugin(value: boolean): Deployment;

  getDeclaredResourcesList(): Array<DeclaredResource>;
  setDeclaredResourcesList(value: Array<DeclaredResource>): Deployment;
  clearDeclaredResourcesList(): Deployment;
  addDeclaredResources(value?: DeclaredResource, index?: number): DeclaredResource;

  getPreload(): Deployment.Preload | undefined;
  setPreload(value?: Deployment.Preload): Deployment;
  hasPreload(): boolean;
  clearPreload(): Deployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Deployment.AsObject;
  static toObject(includeInstance: boolean, msg: Deployment): Deployment.AsObject;
  static serializeBinaryToWriter(message: Deployment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Deployment;
  static deserializeBinaryFromReader(message: Deployment, reader: jspb.BinaryReader): Deployment;
}

export namespace Deployment {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    sequence: number,
    id: string,
    url: string,
    generation?: Generation.AsObject,
    state: Operation.PhysicalState,
    status?: Status.AsObject,
    component?: Component.AsObject,
    artifactId: string,
    deployment?: any_pb.Any.AsObject,
    deploymentJson: string,
    labelsMap: Array<[string, string]>,
    templateData: Uint8Array | string,
    jobId: string,
    hasEntrypointConfig: boolean,
    hasExecPlugin: boolean,
    hasLogsPlugin: boolean,
    declaredResourcesList: Array<DeclaredResource.AsObject>,
    preload?: Deployment.Preload.AsObject,
  }

  export class Preload extends jspb.Message {
    getArtifact(): PushedArtifact | undefined;
    setArtifact(value?: PushedArtifact): Preload;
    hasArtifact(): boolean;
    clearArtifact(): Preload;

    getBuild(): Build | undefined;
    setBuild(value?: Build): Preload;
    hasBuild(): boolean;
    clearBuild(): Preload;

    getDeployUrl(): string;
    setDeployUrl(value: string): Preload;

    getJobDataSourceRef(): Job.DataSource.Ref | undefined;
    setJobDataSourceRef(value?: Job.DataSource.Ref): Preload;
    hasJobDataSourceRef(): boolean;
    clearJobDataSourceRef(): Preload;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Preload.AsObject;
    static toObject(includeInstance: boolean, msg: Preload): Preload.AsObject;
    static serializeBinaryToWriter(message: Preload, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Preload;
    static deserializeBinaryFromReader(message: Preload, reader: jspb.BinaryReader): Preload;
  }

  export namespace Preload {
    export type AsObject = {
      artifact?: PushedArtifact.AsObject,
      build?: Build.AsObject,
      deployUrl: string,
      jobDataSourceRef?: Job.DataSource.Ref.AsObject,
    }
  }


  export enum LoadDetails { 
    NONE = 0,
    ARTIFACT = 1,
    BUILD = 2,
  }
}

export class DeploymentExpanded extends jspb.Message {
  getDeployment(): Deployment | undefined;
  setDeployment(value?: Deployment): DeploymentExpanded;
  hasDeployment(): boolean;
  clearDeployment(): DeploymentExpanded;

  getLatestStatusReport(): StatusReport | undefined;
  setLatestStatusReport(value?: StatusReport): DeploymentExpanded;
  hasLatestStatusReport(): boolean;
  clearLatestStatusReport(): DeploymentExpanded;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentExpanded.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentExpanded): DeploymentExpanded.AsObject;
  static serializeBinaryToWriter(message: DeploymentExpanded, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentExpanded;
  static deserializeBinaryFromReader(message: DeploymentExpanded, reader: jspb.BinaryReader): DeploymentExpanded;
}

export namespace DeploymentExpanded {
  export type AsObject = {
    deployment?: Deployment.AsObject,
    latestStatusReport?: StatusReport.AsObject,
  }
}

export class ListInstancesRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): ListInstancesRequest;

  getApplication(): ListInstancesRequest.Application | undefined;
  setApplication(value?: ListInstancesRequest.Application): ListInstancesRequest;
  hasApplication(): boolean;
  clearApplication(): ListInstancesRequest;

  getWaitTimeout(): string;
  setWaitTimeout(value: string): ListInstancesRequest;

  getScopeCase(): ListInstancesRequest.ScopeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInstancesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListInstancesRequest): ListInstancesRequest.AsObject;
  static serializeBinaryToWriter(message: ListInstancesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInstancesRequest;
  static deserializeBinaryFromReader(message: ListInstancesRequest, reader: jspb.BinaryReader): ListInstancesRequest;
}

export namespace ListInstancesRequest {
  export type AsObject = {
    deploymentId: string,
    application?: ListInstancesRequest.Application.AsObject,
    waitTimeout: string,
  }

  export class Application extends jspb.Message {
    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): Application;
    hasApplication(): boolean;
    clearApplication(): Application;

    getWorkspace(): Ref.Workspace | undefined;
    setWorkspace(value?: Ref.Workspace): Application;
    hasWorkspace(): boolean;
    clearWorkspace(): Application;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Application.AsObject;
    static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
    static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Application;
    static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
  }

  export namespace Application {
    export type AsObject = {
      application?: Ref.Application.AsObject,
      workspace?: Ref.Workspace.AsObject,
    }
  }


  export enum ScopeCase { 
    SCOPE_NOT_SET = 0,
    DEPLOYMENT_ID = 1,
    APPLICATION = 2,
  }
}

export class ListInstancesResponse extends jspb.Message {
  getInstancesList(): Array<Instance>;
  setInstancesList(value: Array<Instance>): ListInstancesResponse;
  clearInstancesList(): ListInstancesResponse;
  addInstances(value?: Instance, index?: number): Instance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInstancesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListInstancesResponse): ListInstancesResponse.AsObject;
  static serializeBinaryToWriter(message: ListInstancesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInstancesResponse;
  static deserializeBinaryFromReader(message: ListInstancesResponse, reader: jspb.BinaryReader): ListInstancesResponse;
}

export namespace ListInstancesResponse {
  export type AsObject = {
    instancesList: Array<Instance.AsObject>,
  }
}

export class Instance extends jspb.Message {
  getId(): string;
  setId(value: string): Instance;

  getDeploymentId(): string;
  setDeploymentId(value: string): Instance;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): Instance;
  hasApplication(): boolean;
  clearApplication(): Instance;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): Instance;
  hasWorkspace(): boolean;
  clearWorkspace(): Instance;

  getType(): Instance.Type;
  setType(value: Instance.Type): Instance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Instance.AsObject;
  static toObject(includeInstance: boolean, msg: Instance): Instance.AsObject;
  static serializeBinaryToWriter(message: Instance, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Instance;
  static deserializeBinaryFromReader(message: Instance, reader: jspb.BinaryReader): Instance;
}

export namespace Instance {
  export type AsObject = {
    id: string,
    deploymentId: string,
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    type: Instance.Type,
  }

  export enum Type { 
    LONG_RUNNING = 0,
    ON_DEMAND = 1,
    VIRTUAL = 2,
  }
}

export class UpsertReleaseRequest extends jspb.Message {
  getRelease(): Release | undefined;
  setRelease(value?: Release): UpsertReleaseRequest;
  hasRelease(): boolean;
  clearRelease(): UpsertReleaseRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertReleaseRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertReleaseRequest): UpsertReleaseRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertReleaseRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertReleaseRequest;
  static deserializeBinaryFromReader(message: UpsertReleaseRequest, reader: jspb.BinaryReader): UpsertReleaseRequest;
}

export namespace UpsertReleaseRequest {
  export type AsObject = {
    release?: Release.AsObject,
  }
}

export class UpsertReleaseResponse extends jspb.Message {
  getRelease(): Release | undefined;
  setRelease(value?: Release): UpsertReleaseResponse;
  hasRelease(): boolean;
  clearRelease(): UpsertReleaseResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertReleaseResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertReleaseResponse): UpsertReleaseResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertReleaseResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertReleaseResponse;
  static deserializeBinaryFromReader(message: UpsertReleaseResponse, reader: jspb.BinaryReader): UpsertReleaseResponse;
}

export namespace UpsertReleaseResponse {
  export type AsObject = {
    release?: Release.AsObject,
  }
}

export class GetLatestReleaseRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): GetLatestReleaseRequest;
  hasApplication(): boolean;
  clearApplication(): GetLatestReleaseRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): GetLatestReleaseRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): GetLatestReleaseRequest;

  getLoadDetails(): Release.LoadDetails;
  setLoadDetails(value: Release.LoadDetails): GetLatestReleaseRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLatestReleaseRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLatestReleaseRequest): GetLatestReleaseRequest.AsObject;
  static serializeBinaryToWriter(message: GetLatestReleaseRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLatestReleaseRequest;
  static deserializeBinaryFromReader(message: GetLatestReleaseRequest, reader: jspb.BinaryReader): GetLatestReleaseRequest;
}

export namespace GetLatestReleaseRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    loadDetails: Release.LoadDetails,
  }
}

export class ListReleasesRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListReleasesRequest;
  hasApplication(): boolean;
  clearApplication(): ListReleasesRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ListReleasesRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ListReleasesRequest;

  getStatusList(): Array<StatusFilter>;
  setStatusList(value: Array<StatusFilter>): ListReleasesRequest;
  clearStatusList(): ListReleasesRequest;
  addStatus(value?: StatusFilter, index?: number): StatusFilter;

  getPhysicalState(): Operation.PhysicalState;
  setPhysicalState(value: Operation.PhysicalState): ListReleasesRequest;

  getOrder(): OperationOrder | undefined;
  setOrder(value?: OperationOrder): ListReleasesRequest;
  hasOrder(): boolean;
  clearOrder(): ListReleasesRequest;

  getLoadDetails(): Release.LoadDetails;
  setLoadDetails(value: Release.LoadDetails): ListReleasesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListReleasesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListReleasesRequest): ListReleasesRequest.AsObject;
  static serializeBinaryToWriter(message: ListReleasesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListReleasesRequest;
  static deserializeBinaryFromReader(message: ListReleasesRequest, reader: jspb.BinaryReader): ListReleasesRequest;
}

export namespace ListReleasesRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    statusList: Array<StatusFilter.AsObject>,
    physicalState: Operation.PhysicalState,
    order?: OperationOrder.AsObject,
    loadDetails: Release.LoadDetails,
  }
}

export class ListReleasesResponse extends jspb.Message {
  getReleasesList(): Array<Release>;
  setReleasesList(value: Array<Release>): ListReleasesResponse;
  clearReleasesList(): ListReleasesResponse;
  addReleases(value?: Release, index?: number): Release;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListReleasesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListReleasesResponse): ListReleasesResponse.AsObject;
  static serializeBinaryToWriter(message: ListReleasesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListReleasesResponse;
  static deserializeBinaryFromReader(message: ListReleasesResponse, reader: jspb.BinaryReader): ListReleasesResponse;
}

export namespace ListReleasesResponse {
  export type AsObject = {
    releasesList: Array<Release.AsObject>,
  }
}

export class GetReleaseRequest extends jspb.Message {
  getRef(): Ref.Operation | undefined;
  setRef(value?: Ref.Operation): GetReleaseRequest;
  hasRef(): boolean;
  clearRef(): GetReleaseRequest;

  getLoadDetails(): Release.LoadDetails;
  setLoadDetails(value: Release.LoadDetails): GetReleaseRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetReleaseRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetReleaseRequest): GetReleaseRequest.AsObject;
  static serializeBinaryToWriter(message: GetReleaseRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetReleaseRequest;
  static deserializeBinaryFromReader(message: GetReleaseRequest, reader: jspb.BinaryReader): GetReleaseRequest;
}

export namespace GetReleaseRequest {
  export type AsObject = {
    ref?: Ref.Operation.AsObject,
    loadDetails: Release.LoadDetails,
  }
}

export class Release extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): Release;
  hasApplication(): boolean;
  clearApplication(): Release;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): Release;
  hasWorkspace(): boolean;
  clearWorkspace(): Release;

  getSequence(): number;
  setSequence(value: number): Release;

  getId(): string;
  setId(value: string): Release;

  getStatus(): Status | undefined;
  setStatus(value?: Status): Release;
  hasStatus(): boolean;
  clearStatus(): Release;

  getState(): Operation.PhysicalState;
  setState(value: Operation.PhysicalState): Release;

  getComponent(): Component | undefined;
  setComponent(value?: Component): Release;
  hasComponent(): boolean;
  clearComponent(): Release;

  getRelease(): any_pb.Any | undefined;
  setRelease(value?: any_pb.Any): Release;
  hasRelease(): boolean;
  clearRelease(): Release;

  getReleaseJson(): string;
  setReleaseJson(value: string): Release;

  getDeploymentId(): string;
  setDeploymentId(value: string): Release;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Release;

  getTemplateData(): Uint8Array | string;
  getTemplateData_asU8(): Uint8Array;
  getTemplateData_asB64(): string;
  setTemplateData(value: Uint8Array | string): Release;

  getUrl(): string;
  setUrl(value: string): Release;

  getJobId(): string;
  setJobId(value: string): Release;

  getUnimplemented(): boolean;
  setUnimplemented(value: boolean): Release;

  getDeclaredResourcesList(): Array<DeclaredResource>;
  setDeclaredResourcesList(value: Array<DeclaredResource>): Release;
  clearDeclaredResourcesList(): Release;
  addDeclaredResources(value?: DeclaredResource, index?: number): DeclaredResource;

  getPreload(): Release.Preload | undefined;
  setPreload(value?: Release.Preload): Release;
  hasPreload(): boolean;
  clearPreload(): Release;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Release.AsObject;
  static toObject(includeInstance: boolean, msg: Release): Release.AsObject;
  static serializeBinaryToWriter(message: Release, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Release;
  static deserializeBinaryFromReader(message: Release, reader: jspb.BinaryReader): Release;
}

export namespace Release {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    sequence: number,
    id: string,
    status?: Status.AsObject,
    state: Operation.PhysicalState,
    component?: Component.AsObject,
    release?: any_pb.Any.AsObject,
    releaseJson: string,
    deploymentId: string,
    labelsMap: Array<[string, string]>,
    templateData: Uint8Array | string,
    url: string,
    jobId: string,
    unimplemented: boolean,
    declaredResourcesList: Array<DeclaredResource.AsObject>,
    preload?: Release.Preload.AsObject,
  }

  export class Preload extends jspb.Message {
    getDeployment(): Deployment | undefined;
    setDeployment(value?: Deployment): Preload;
    hasDeployment(): boolean;
    clearDeployment(): Preload;

    getArtifact(): PushedArtifact | undefined;
    setArtifact(value?: PushedArtifact): Preload;
    hasArtifact(): boolean;
    clearArtifact(): Preload;

    getBuild(): Build | undefined;
    setBuild(value?: Build): Preload;
    hasBuild(): boolean;
    clearBuild(): Preload;

    getJobDataSourceRef(): Job.DataSource.Ref | undefined;
    setJobDataSourceRef(value?: Job.DataSource.Ref): Preload;
    hasJobDataSourceRef(): boolean;
    clearJobDataSourceRef(): Preload;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Preload.AsObject;
    static toObject(includeInstance: boolean, msg: Preload): Preload.AsObject;
    static serializeBinaryToWriter(message: Preload, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Preload;
    static deserializeBinaryFromReader(message: Preload, reader: jspb.BinaryReader): Preload;
  }

  export namespace Preload {
    export type AsObject = {
      deployment?: Deployment.AsObject,
      artifact?: PushedArtifact.AsObject,
      build?: Build.AsObject,
      jobDataSourceRef?: Job.DataSource.Ref.AsObject,
    }
  }


  export enum LoadDetails { 
    NONE = 0,
    DEPLOYMENT = 1,
    ARTIFACT = 2,
    BUILD = 3,
  }
}

export class ReleaseExpanded extends jspb.Message {
  getRelease(): Release | undefined;
  setRelease(value?: Release): ReleaseExpanded;
  hasRelease(): boolean;
  clearRelease(): ReleaseExpanded;

  getLatestStatusReport(): StatusReport | undefined;
  setLatestStatusReport(value?: StatusReport): ReleaseExpanded;
  hasLatestStatusReport(): boolean;
  clearLatestStatusReport(): ReleaseExpanded;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReleaseExpanded.AsObject;
  static toObject(includeInstance: boolean, msg: ReleaseExpanded): ReleaseExpanded.AsObject;
  static serializeBinaryToWriter(message: ReleaseExpanded, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReleaseExpanded;
  static deserializeBinaryFromReader(message: ReleaseExpanded, reader: jspb.BinaryReader): ReleaseExpanded;
}

export namespace ReleaseExpanded {
  export type AsObject = {
    release?: Release.AsObject,
    latestStatusReport?: StatusReport.AsObject,
  }
}

export class Trigger extends jspb.Message {
  getId(): string;
  setId(value: string): Trigger;

  getName(): string;
  setName(value: string): Trigger;

  getDescription(): string;
  setDescription(value: string): Trigger;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): Trigger;
  clearTagsList(): Trigger;
  addTags(value: string, index?: number): Trigger;

  getActiveTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setActiveTime(value?: google_protobuf_timestamp_pb.Timestamp): Trigger;
  hasActiveTime(): boolean;
  clearActiveTime(): Trigger;

  getAuthenticated(): boolean;
  setAuthenticated(value: boolean): Trigger;

  getBuild(): Job.BuildOp | undefined;
  setBuild(value?: Job.BuildOp): Trigger;
  hasBuild(): boolean;
  clearBuild(): Trigger;

  getPush(): Job.PushOp | undefined;
  setPush(value?: Job.PushOp): Trigger;
  hasPush(): boolean;
  clearPush(): Trigger;

  getDeploy(): Job.DeployOp | undefined;
  setDeploy(value?: Job.DeployOp): Trigger;
  hasDeploy(): boolean;
  clearDeploy(): Trigger;

  getDestroy(): Job.DestroyOp | undefined;
  setDestroy(value?: Job.DestroyOp): Trigger;
  hasDestroy(): boolean;
  clearDestroy(): Trigger;

  getRelease(): Job.ReleaseOp | undefined;
  setRelease(value?: Job.ReleaseOp): Trigger;
  hasRelease(): boolean;
  clearRelease(): Trigger;

  getUp(): Job.UpOp | undefined;
  setUp(value?: Job.UpOp): Trigger;
  hasUp(): boolean;
  clearUp(): Trigger;

  getInit(): Job.InitOp | undefined;
  setInit(value?: Job.InitOp): Trigger;
  hasInit(): boolean;
  clearInit(): Trigger;

  getStatusReport(): Job.StatusReportOp | undefined;
  setStatusReport(value?: Job.StatusReportOp): Trigger;
  hasStatusReport(): boolean;
  clearStatusReport(): Trigger;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): Trigger;
  hasWorkspace(): boolean;
  clearWorkspace(): Trigger;

  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): Trigger;
  hasProject(): boolean;
  clearProject(): Trigger;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): Trigger;
  hasApplication(): boolean;
  clearApplication(): Trigger;

  getOperationCase(): Trigger.OperationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Trigger.AsObject;
  static toObject(includeInstance: boolean, msg: Trigger): Trigger.AsObject;
  static serializeBinaryToWriter(message: Trigger, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Trigger;
  static deserializeBinaryFromReader(message: Trigger, reader: jspb.BinaryReader): Trigger;
}

export namespace Trigger {
  export type AsObject = {
    id: string,
    name: string,
    description: string,
    tagsList: Array<string>,
    activeTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    authenticated: boolean,
    build?: Job.BuildOp.AsObject,
    push?: Job.PushOp.AsObject,
    deploy?: Job.DeployOp.AsObject,
    destroy?: Job.DestroyOp.AsObject,
    release?: Job.ReleaseOp.AsObject,
    up?: Job.UpOp.AsObject,
    init?: Job.InitOp.AsObject,
    statusReport?: Job.StatusReportOp.AsObject,
    workspace?: Ref.Workspace.AsObject,
    project?: Ref.Project.AsObject,
    application?: Ref.Application.AsObject,
  }

  export enum OperationCase { 
    OPERATION_NOT_SET = 0,
    BUILD = 50,
    PUSH = 51,
    DEPLOY = 52,
    DESTROY = 53,
    RELEASE = 54,
    UP = 55,
    INIT = 56,
    STATUS_REPORT = 57,
  }
}

export class UpsertTriggerRequest extends jspb.Message {
  getTrigger(): Trigger | undefined;
  setTrigger(value?: Trigger): UpsertTriggerRequest;
  hasTrigger(): boolean;
  clearTrigger(): UpsertTriggerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertTriggerRequest): UpsertTriggerRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertTriggerRequest;
  static deserializeBinaryFromReader(message: UpsertTriggerRequest, reader: jspb.BinaryReader): UpsertTriggerRequest;
}

export namespace UpsertTriggerRequest {
  export type AsObject = {
    trigger?: Trigger.AsObject,
  }
}

export class UpsertTriggerResponse extends jspb.Message {
  getTrigger(): Trigger | undefined;
  setTrigger(value?: Trigger): UpsertTriggerResponse;
  hasTrigger(): boolean;
  clearTrigger(): UpsertTriggerResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertTriggerResponse): UpsertTriggerResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertTriggerResponse;
  static deserializeBinaryFromReader(message: UpsertTriggerResponse, reader: jspb.BinaryReader): UpsertTriggerResponse;
}

export namespace UpsertTriggerResponse {
  export type AsObject = {
    trigger?: Trigger.AsObject,
  }
}

export class GetTriggerRequest extends jspb.Message {
  getRef(): Ref.Trigger | undefined;
  setRef(value?: Ref.Trigger): GetTriggerRequest;
  hasRef(): boolean;
  clearRef(): GetTriggerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTriggerRequest): GetTriggerRequest.AsObject;
  static serializeBinaryToWriter(message: GetTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTriggerRequest;
  static deserializeBinaryFromReader(message: GetTriggerRequest, reader: jspb.BinaryReader): GetTriggerRequest;
}

export namespace GetTriggerRequest {
  export type AsObject = {
    ref?: Ref.Trigger.AsObject,
  }
}

export class GetTriggerResponse extends jspb.Message {
  getTrigger(): Trigger | undefined;
  setTrigger(value?: Trigger): GetTriggerResponse;
  hasTrigger(): boolean;
  clearTrigger(): GetTriggerResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetTriggerResponse): GetTriggerResponse.AsObject;
  static serializeBinaryToWriter(message: GetTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTriggerResponse;
  static deserializeBinaryFromReader(message: GetTriggerResponse, reader: jspb.BinaryReader): GetTriggerResponse;
}

export namespace GetTriggerResponse {
  export type AsObject = {
    trigger?: Trigger.AsObject,
  }
}

export class DeleteTriggerRequest extends jspb.Message {
  getRef(): Ref.Trigger | undefined;
  setRef(value?: Ref.Trigger): DeleteTriggerRequest;
  hasRef(): boolean;
  clearRef(): DeleteTriggerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteTriggerRequest): DeleteTriggerRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTriggerRequest;
  static deserializeBinaryFromReader(message: DeleteTriggerRequest, reader: jspb.BinaryReader): DeleteTriggerRequest;
}

export namespace DeleteTriggerRequest {
  export type AsObject = {
    ref?: Ref.Trigger.AsObject,
  }
}

export class RunTriggerRequest extends jspb.Message {
  getRef(): Ref.Trigger | undefined;
  setRef(value?: Ref.Trigger): RunTriggerRequest;
  hasRef(): boolean;
  clearRef(): RunTriggerRequest;

  getVariableOverridesList(): Array<Variable>;
  setVariableOverridesList(value: Array<Variable>): RunTriggerRequest;
  clearVariableOverridesList(): RunTriggerRequest;
  addVariableOverrides(value?: Variable, index?: number): Variable;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RunTriggerRequest): RunTriggerRequest.AsObject;
  static serializeBinaryToWriter(message: RunTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunTriggerRequest;
  static deserializeBinaryFromReader(message: RunTriggerRequest, reader: jspb.BinaryReader): RunTriggerRequest;
}

export namespace RunTriggerRequest {
  export type AsObject = {
    ref?: Ref.Trigger.AsObject,
    variableOverridesList: Array<Variable.AsObject>,
  }
}

export class RunTriggerResponse extends jspb.Message {
  getJobIdsList(): Array<string>;
  setJobIdsList(value: Array<string>): RunTriggerResponse;
  clearJobIdsList(): RunTriggerResponse;
  addJobIds(value: string, index?: number): RunTriggerResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RunTriggerResponse): RunTriggerResponse.AsObject;
  static serializeBinaryToWriter(message: RunTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunTriggerResponse;
  static deserializeBinaryFromReader(message: RunTriggerResponse, reader: jspb.BinaryReader): RunTriggerResponse;
}

export namespace RunTriggerResponse {
  export type AsObject = {
    jobIdsList: Array<string>,
  }
}

export class ListTriggerRequest extends jspb.Message {
  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ListTriggerRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ListTriggerRequest;

  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): ListTriggerRequest;
  hasProject(): boolean;
  clearProject(): ListTriggerRequest;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListTriggerRequest;
  hasApplication(): boolean;
  clearApplication(): ListTriggerRequest;

  getTagsList(): Array<string>;
  setTagsList(value: Array<string>): ListTriggerRequest;
  clearTagsList(): ListTriggerRequest;
  addTags(value: string, index?: number): ListTriggerRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTriggerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListTriggerRequest): ListTriggerRequest.AsObject;
  static serializeBinaryToWriter(message: ListTriggerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTriggerRequest;
  static deserializeBinaryFromReader(message: ListTriggerRequest, reader: jspb.BinaryReader): ListTriggerRequest;
}

export namespace ListTriggerRequest {
  export type AsObject = {
    workspace?: Ref.Workspace.AsObject,
    project?: Ref.Project.AsObject,
    application?: Ref.Application.AsObject,
    tagsList: Array<string>,
  }
}

export class ListTriggerResponse extends jspb.Message {
  getTriggersList(): Array<Trigger>;
  setTriggersList(value: Array<Trigger>): ListTriggerResponse;
  clearTriggersList(): ListTriggerResponse;
  addTriggers(value?: Trigger, index?: number): Trigger;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTriggerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTriggerResponse): ListTriggerResponse.AsObject;
  static serializeBinaryToWriter(message: ListTriggerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTriggerResponse;
  static deserializeBinaryFromReader(message: ListTriggerResponse, reader: jspb.BinaryReader): ListTriggerResponse;
}

export namespace ListTriggerResponse {
  export type AsObject = {
    triggersList: Array<Trigger.AsObject>,
  }
}

export class UpsertStatusReportRequest extends jspb.Message {
  getStatusReport(): StatusReport | undefined;
  setStatusReport(value?: StatusReport): UpsertStatusReportRequest;
  hasStatusReport(): boolean;
  clearStatusReport(): UpsertStatusReportRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertStatusReportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertStatusReportRequest): UpsertStatusReportRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertStatusReportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertStatusReportRequest;
  static deserializeBinaryFromReader(message: UpsertStatusReportRequest, reader: jspb.BinaryReader): UpsertStatusReportRequest;
}

export namespace UpsertStatusReportRequest {
  export type AsObject = {
    statusReport?: StatusReport.AsObject,
  }
}

export class UpsertStatusReportResponse extends jspb.Message {
  getStatusReport(): StatusReport | undefined;
  setStatusReport(value?: StatusReport): UpsertStatusReportResponse;
  hasStatusReport(): boolean;
  clearStatusReport(): UpsertStatusReportResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertStatusReportResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertStatusReportResponse): UpsertStatusReportResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertStatusReportResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertStatusReportResponse;
  static deserializeBinaryFromReader(message: UpsertStatusReportResponse, reader: jspb.BinaryReader): UpsertStatusReportResponse;
}

export namespace UpsertStatusReportResponse {
  export type AsObject = {
    statusReport?: StatusReport.AsObject,
  }
}

export class GetLatestStatusReportRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): GetLatestStatusReportRequest;
  hasApplication(): boolean;
  clearApplication(): GetLatestStatusReportRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): GetLatestStatusReportRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): GetLatestStatusReportRequest;

  getAny(): google_protobuf_empty_pb.Empty | undefined;
  setAny(value?: google_protobuf_empty_pb.Empty): GetLatestStatusReportRequest;
  hasAny(): boolean;
  clearAny(): GetLatestStatusReportRequest;

  getDeploymentAny(): google_protobuf_empty_pb.Empty | undefined;
  setDeploymentAny(value?: google_protobuf_empty_pb.Empty): GetLatestStatusReportRequest;
  hasDeploymentAny(): boolean;
  clearDeploymentAny(): GetLatestStatusReportRequest;

  getDeploymentId(): string;
  setDeploymentId(value: string): GetLatestStatusReportRequest;

  getReleaseAny(): google_protobuf_empty_pb.Empty | undefined;
  setReleaseAny(value?: google_protobuf_empty_pb.Empty): GetLatestStatusReportRequest;
  hasReleaseAny(): boolean;
  clearReleaseAny(): GetLatestStatusReportRequest;

  getReleaseId(): string;
  setReleaseId(value: string): GetLatestStatusReportRequest;

  getTargetCase(): GetLatestStatusReportRequest.TargetCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLatestStatusReportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLatestStatusReportRequest): GetLatestStatusReportRequest.AsObject;
  static serializeBinaryToWriter(message: GetLatestStatusReportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLatestStatusReportRequest;
  static deserializeBinaryFromReader(message: GetLatestStatusReportRequest, reader: jspb.BinaryReader): GetLatestStatusReportRequest;
}

export namespace GetLatestStatusReportRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    any?: google_protobuf_empty_pb.Empty.AsObject,
    deploymentAny?: google_protobuf_empty_pb.Empty.AsObject,
    deploymentId: string,
    releaseAny?: google_protobuf_empty_pb.Empty.AsObject,
    releaseId: string,
  }

  export enum TargetCase { 
    TARGET_NOT_SET = 0,
    ANY = 3,
    DEPLOYMENT_ANY = 4,
    DEPLOYMENT_ID = 5,
    RELEASE_ANY = 6,
    RELEASE_ID = 7,
  }
}

export class ListStatusReportsRequest extends jspb.Message {
  getStatusList(): Array<StatusFilter>;
  setStatusList(value: Array<StatusFilter>): ListStatusReportsRequest;
  clearStatusList(): ListStatusReportsRequest;
  addStatus(value?: StatusFilter, index?: number): StatusFilter;

  getOrder(): OperationOrder | undefined;
  setOrder(value?: OperationOrder): ListStatusReportsRequest;
  hasOrder(): boolean;
  clearOrder(): ListStatusReportsRequest;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ListStatusReportsRequest;
  hasApplication(): boolean;
  clearApplication(): ListStatusReportsRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ListStatusReportsRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ListStatusReportsRequest;

  getDeployment(): Ref.Operation | undefined;
  setDeployment(value?: Ref.Operation): ListStatusReportsRequest;
  hasDeployment(): boolean;
  clearDeployment(): ListStatusReportsRequest;

  getRelease(): Ref.Operation | undefined;
  setRelease(value?: Ref.Operation): ListStatusReportsRequest;
  hasRelease(): boolean;
  clearRelease(): ListStatusReportsRequest;

  getTargetCase(): ListStatusReportsRequest.TargetCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListStatusReportsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListStatusReportsRequest): ListStatusReportsRequest.AsObject;
  static serializeBinaryToWriter(message: ListStatusReportsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListStatusReportsRequest;
  static deserializeBinaryFromReader(message: ListStatusReportsRequest, reader: jspb.BinaryReader): ListStatusReportsRequest;
}

export namespace ListStatusReportsRequest {
  export type AsObject = {
    statusList: Array<StatusFilter.AsObject>,
    order?: OperationOrder.AsObject,
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    deployment?: Ref.Operation.AsObject,
    release?: Ref.Operation.AsObject,
  }

  export enum TargetCase { 
    TARGET_NOT_SET = 0,
    DEPLOYMENT = 5,
    RELEASE = 6,
  }
}

export class ListStatusReportsResponse extends jspb.Message {
  getStatusReportsList(): Array<StatusReport>;
  setStatusReportsList(value: Array<StatusReport>): ListStatusReportsResponse;
  clearStatusReportsList(): ListStatusReportsResponse;
  addStatusReports(value?: StatusReport, index?: number): StatusReport;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListStatusReportsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListStatusReportsResponse): ListStatusReportsResponse.AsObject;
  static serializeBinaryToWriter(message: ListStatusReportsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListStatusReportsResponse;
  static deserializeBinaryFromReader(message: ListStatusReportsResponse, reader: jspb.BinaryReader): ListStatusReportsResponse;
}

export namespace ListStatusReportsResponse {
  export type AsObject = {
    statusReportsList: Array<StatusReport.AsObject>,
  }
}

export class GetStatusReportRequest extends jspb.Message {
  getRef(): Ref.Operation | undefined;
  setRef(value?: Ref.Operation): GetStatusReportRequest;
  hasRef(): boolean;
  clearRef(): GetStatusReportRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetStatusReportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetStatusReportRequest): GetStatusReportRequest.AsObject;
  static serializeBinaryToWriter(message: GetStatusReportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetStatusReportRequest;
  static deserializeBinaryFromReader(message: GetStatusReportRequest, reader: jspb.BinaryReader): GetStatusReportRequest;
}

export namespace GetStatusReportRequest {
  export type AsObject = {
    ref?: Ref.Operation.AsObject,
  }
}

export class ExpediteStatusReportRequest extends jspb.Message {
  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ExpediteStatusReportRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ExpediteStatusReportRequest;

  getDeployment(): Ref.Operation | undefined;
  setDeployment(value?: Ref.Operation): ExpediteStatusReportRequest;
  hasDeployment(): boolean;
  clearDeployment(): ExpediteStatusReportRequest;

  getRelease(): Ref.Operation | undefined;
  setRelease(value?: Ref.Operation): ExpediteStatusReportRequest;
  hasRelease(): boolean;
  clearRelease(): ExpediteStatusReportRequest;

  getTargetCase(): ExpediteStatusReportRequest.TargetCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExpediteStatusReportRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExpediteStatusReportRequest): ExpediteStatusReportRequest.AsObject;
  static serializeBinaryToWriter(message: ExpediteStatusReportRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExpediteStatusReportRequest;
  static deserializeBinaryFromReader(message: ExpediteStatusReportRequest, reader: jspb.BinaryReader): ExpediteStatusReportRequest;
}

export namespace ExpediteStatusReportRequest {
  export type AsObject = {
    workspace?: Ref.Workspace.AsObject,
    deployment?: Ref.Operation.AsObject,
    release?: Ref.Operation.AsObject,
  }

  export enum TargetCase { 
    TARGET_NOT_SET = 0,
    DEPLOYMENT = 2,
    RELEASE = 3,
  }
}

export class ExpediteStatusReportResponse extends jspb.Message {
  getJobId(): string;
  setJobId(value: string): ExpediteStatusReportResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExpediteStatusReportResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExpediteStatusReportResponse): ExpediteStatusReportResponse.AsObject;
  static serializeBinaryToWriter(message: ExpediteStatusReportResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExpediteStatusReportResponse;
  static deserializeBinaryFromReader(message: ExpediteStatusReportResponse, reader: jspb.BinaryReader): ExpediteStatusReportResponse;
}

export namespace ExpediteStatusReportResponse {
  export type AsObject = {
    jobId: string,
  }
}

export class StatusReport extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): StatusReport;
  hasApplication(): boolean;
  clearApplication(): StatusReport;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): StatusReport;
  hasWorkspace(): boolean;
  clearWorkspace(): StatusReport;

  getDeploymentId(): string;
  setDeploymentId(value: string): StatusReport;

  getReleaseId(): string;
  setReleaseId(value: string): StatusReport;

  getStatus(): Status | undefined;
  setStatus(value?: Status): StatusReport;
  hasStatus(): boolean;
  clearStatus(): StatusReport;

  getId(): string;
  setId(value: string): StatusReport;

  getStatusReport(): any_pb.Any | undefined;
  setStatusReport(value?: any_pb.Any): StatusReport;
  hasStatusReport(): boolean;
  clearStatusReport(): StatusReport;

  getStatusReportJson(): string;
  setStatusReportJson(value: string): StatusReport;

  getHealth(): StatusReport.Health | undefined;
  setHealth(value?: StatusReport.Health): StatusReport;
  hasHealth(): boolean;
  clearHealth(): StatusReport;

  getGeneratedTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setGeneratedTime(value?: google_protobuf_timestamp_pb.Timestamp): StatusReport;
  hasGeneratedTime(): boolean;
  clearGeneratedTime(): StatusReport;

  getExternal(): boolean;
  setExternal(value: boolean): StatusReport;

  getResourcesList(): Array<StatusReport.Resource>;
  setResourcesList(value: Array<StatusReport.Resource>): StatusReport;
  clearResourcesList(): StatusReport;
  addResources(value?: StatusReport.Resource, index?: number): StatusReport.Resource;

  getDeprecatedResourcesHealthList(): Array<StatusReport.Health>;
  setDeprecatedResourcesHealthList(value: Array<StatusReport.Health>): StatusReport;
  clearDeprecatedResourcesHealthList(): StatusReport;
  addDeprecatedResourcesHealth(value?: StatusReport.Health, index?: number): StatusReport.Health;

  getInstancesCount(): number;
  setInstancesCount(value: number): StatusReport;

  getTargetIdCase(): StatusReport.TargetIdCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusReport.AsObject;
  static toObject(includeInstance: boolean, msg: StatusReport): StatusReport.AsObject;
  static serializeBinaryToWriter(message: StatusReport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusReport;
  static deserializeBinaryFromReader(message: StatusReport, reader: jspb.BinaryReader): StatusReport;
}

export namespace StatusReport {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    workspace?: Ref.Workspace.AsObject,
    deploymentId: string,
    releaseId: string,
    status?: Status.AsObject,
    id: string,
    statusReport?: any_pb.Any.AsObject,
    statusReportJson: string,
    health?: StatusReport.Health.AsObject,
    generatedTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    external: boolean,
    resourcesList: Array<StatusReport.Resource.AsObject>,
    deprecatedResourcesHealthList: Array<StatusReport.Health.AsObject>,
    instancesCount: number,
  }

  export class Resource extends jspb.Message {
    getId(): string;
    setId(value: string): Resource;

    getDeclaredResource(): Ref.DeclaredResource | undefined;
    setDeclaredResource(value?: Ref.DeclaredResource): Resource;
    hasDeclaredResource(): boolean;
    clearDeclaredResource(): Resource;

    getParentResourceId(): string;
    setParentResourceId(value: string): Resource;

    getName(): string;
    setName(value: string): Resource;

    getPlatform(): string;
    setPlatform(value: string): Resource;

    getType(): string;
    setType(value: string): Resource;

    getPlatformUrl(): string;
    setPlatformUrl(value: string): Resource;

    getCategoryDisplayHint(): ResourceCategoryDisplayHint;
    setCategoryDisplayHint(value: ResourceCategoryDisplayHint): Resource;

    getCreatedTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setCreatedTime(value?: google_protobuf_timestamp_pb.Timestamp): Resource;
    hasCreatedTime(): boolean;
    clearCreatedTime(): Resource;

    getStateJson(): string;
    setStateJson(value: string): Resource;

    getHealth(): StatusReport.Resource.Health;
    setHealth(value: StatusReport.Resource.Health): Resource;

    getHealthMessage(): string;
    setHealthMessage(value: string): Resource;

    getDeprecatedHealth(): StatusReport.Health | undefined;
    setDeprecatedHealth(value?: StatusReport.Health): Resource;
    hasDeprecatedHealth(): boolean;
    clearDeprecatedHealth(): Resource;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Resource.AsObject;
    static toObject(includeInstance: boolean, msg: Resource): Resource.AsObject;
    static serializeBinaryToWriter(message: Resource, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Resource;
    static deserializeBinaryFromReader(message: Resource, reader: jspb.BinaryReader): Resource;
  }

  export namespace Resource {
    export type AsObject = {
      id: string,
      declaredResource?: Ref.DeclaredResource.AsObject,
      parentResourceId: string,
      name: string,
      platform: string,
      type: string,
      platformUrl: string,
      categoryDisplayHint: ResourceCategoryDisplayHint,
      createdTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      stateJson: string,
      health: StatusReport.Resource.Health,
      healthMessage: string,
      deprecatedHealth?: StatusReport.Health.AsObject,
    }

    export enum Health { 
      UNKNOWN = 0,
      ALIVE = 1,
      READY = 2,
      DOWN = 3,
      MISSING = 5,
      PARTIAL = 4,
    }
  }


  export class Health extends jspb.Message {
    getHealthStatus(): string;
    setHealthStatus(value: string): Health;

    getHealthMessage(): string;
    setHealthMessage(value: string): Health;

    getDeprecatedName(): string;
    setDeprecatedName(value: string): Health;

    getDeprecatedId(): string;
    setDeprecatedId(value: string): Health;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Health.AsObject;
    static toObject(includeInstance: boolean, msg: Health): Health.AsObject;
    static serializeBinaryToWriter(message: Health, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Health;
    static deserializeBinaryFromReader(message: Health, reader: jspb.BinaryReader): Health;
  }

  export namespace Health {
    export type AsObject = {
      healthStatus: string,
      healthMessage: string,
      deprecatedName: string,
      deprecatedId: string,
    }
  }


  export enum TargetIdCase { 
    TARGET_ID_NOT_SET = 0,
    DEPLOYMENT_ID = 3,
    RELEASE_ID = 4,
  }
}

export class GetLogStreamRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): GetLogStreamRequest;

  getApplication(): GetLogStreamRequest.Application | undefined;
  setApplication(value?: GetLogStreamRequest.Application): GetLogStreamRequest;
  hasApplication(): boolean;
  clearApplication(): GetLogStreamRequest;

  getLimitBacklog(): number;
  setLimitBacklog(value: number): GetLogStreamRequest;

  getScopeCase(): GetLogStreamRequest.ScopeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLogStreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLogStreamRequest): GetLogStreamRequest.AsObject;
  static serializeBinaryToWriter(message: GetLogStreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLogStreamRequest;
  static deserializeBinaryFromReader(message: GetLogStreamRequest, reader: jspb.BinaryReader): GetLogStreamRequest;
}

export namespace GetLogStreamRequest {
  export type AsObject = {
    deploymentId: string,
    application?: GetLogStreamRequest.Application.AsObject,
    limitBacklog: number,
  }

  export class Application extends jspb.Message {
    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): Application;
    hasApplication(): boolean;
    clearApplication(): Application;

    getWorkspace(): Ref.Workspace | undefined;
    setWorkspace(value?: Ref.Workspace): Application;
    hasWorkspace(): boolean;
    clearWorkspace(): Application;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Application.AsObject;
    static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
    static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Application;
    static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
  }

  export namespace Application {
    export type AsObject = {
      application?: Ref.Application.AsObject,
      workspace?: Ref.Workspace.AsObject,
    }
  }


  export enum ScopeCase { 
    SCOPE_NOT_SET = 0,
    DEPLOYMENT_ID = 1,
    APPLICATION = 2,
  }
}

export class LogBatch extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): LogBatch;

  getInstanceId(): string;
  setInstanceId(value: string): LogBatch;

  getLinesList(): Array<LogBatch.Entry>;
  setLinesList(value: Array<LogBatch.Entry>): LogBatch;
  clearLinesList(): LogBatch;
  addLines(value?: LogBatch.Entry, index?: number): LogBatch.Entry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogBatch.AsObject;
  static toObject(includeInstance: boolean, msg: LogBatch): LogBatch.AsObject;
  static serializeBinaryToWriter(message: LogBatch, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogBatch;
  static deserializeBinaryFromReader(message: LogBatch, reader: jspb.BinaryReader): LogBatch;
}

export namespace LogBatch {
  export type AsObject = {
    deploymentId: string,
    instanceId: string,
    linesList: Array<LogBatch.Entry.AsObject>,
  }

  export class Entry extends jspb.Message {
    getSource(): LogBatch.Entry.Source;
    setSource(value: LogBatch.Entry.Source): Entry;

    getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): Entry;
    hasTimestamp(): boolean;
    clearTimestamp(): Entry;

    getLine(): string;
    setLine(value: string): Entry;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Entry.AsObject;
    static toObject(includeInstance: boolean, msg: Entry): Entry.AsObject;
    static serializeBinaryToWriter(message: Entry, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Entry;
    static deserializeBinaryFromReader(message: Entry, reader: jspb.BinaryReader): Entry;
  }

  export namespace Entry {
    export type AsObject = {
      source: LogBatch.Entry.Source,
      timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      line: string,
    }

    export enum Source { 
      APP = 0,
      ENTRYPOINT = 1,
    }
  }

}

export class ConfigVar extends jspb.Message {
  getTarget(): ConfigVar.Target | undefined;
  setTarget(value?: ConfigVar.Target): ConfigVar;
  hasTarget(): boolean;
  clearTarget(): ConfigVar;

  getName(): string;
  setName(value: string): ConfigVar;

  getUnset(): google_protobuf_empty_pb.Empty | undefined;
  setUnset(value?: google_protobuf_empty_pb.Empty): ConfigVar;
  hasUnset(): boolean;
  clearUnset(): ConfigVar;

  getStatic(): string;
  setStatic(value: string): ConfigVar;

  getDynamic(): ConfigVar.DynamicVal | undefined;
  setDynamic(value?: ConfigVar.DynamicVal): ConfigVar;
  hasDynamic(): boolean;
  clearDynamic(): ConfigVar;

  getInternal(): boolean;
  setInternal(value: boolean): ConfigVar;

  getNameIsPath(): boolean;
  setNameIsPath(value: boolean): ConfigVar;

  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ConfigVar;
  hasApplication(): boolean;
  clearApplication(): ConfigVar;

  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): ConfigVar;
  hasProject(): boolean;
  clearProject(): ConfigVar;

  getRunner(): Ref.Runner | undefined;
  setRunner(value?: Ref.Runner): ConfigVar;
  hasRunner(): boolean;
  clearRunner(): ConfigVar;

  getValueCase(): ConfigVar.ValueCase;

  getUnusedScopeCase(): ConfigVar.UnusedScopeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigVar.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigVar): ConfigVar.AsObject;
  static serializeBinaryToWriter(message: ConfigVar, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigVar;
  static deserializeBinaryFromReader(message: ConfigVar, reader: jspb.BinaryReader): ConfigVar;
}

export namespace ConfigVar {
  export type AsObject = {
    target?: ConfigVar.Target.AsObject,
    name: string,
    unset?: google_protobuf_empty_pb.Empty.AsObject,
    pb_static: string,
    dynamic?: ConfigVar.DynamicVal.AsObject,
    internal: boolean,
    nameIsPath: boolean,
    application?: Ref.Application.AsObject,
    project?: Ref.Project.AsObject,
    runner?: Ref.Runner.AsObject,
  }

  export class DynamicVal extends jspb.Message {
    getFrom(): string;
    setFrom(value: string): DynamicVal;

    getConfigMap(): jspb.Map<string, string>;
    clearConfigMap(): DynamicVal;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DynamicVal.AsObject;
    static toObject(includeInstance: boolean, msg: DynamicVal): DynamicVal.AsObject;
    static serializeBinaryToWriter(message: DynamicVal, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DynamicVal;
    static deserializeBinaryFromReader(message: DynamicVal, reader: jspb.BinaryReader): DynamicVal;
  }

  export namespace DynamicVal {
    export type AsObject = {
      from: string,
      configMap: Array<[string, string]>,
    }
  }


  export class Target extends jspb.Message {
    getGlobal(): Ref.Global | undefined;
    setGlobal(value?: Ref.Global): Target;
    hasGlobal(): boolean;
    clearGlobal(): Target;

    getProject(): Ref.Project | undefined;
    setProject(value?: Ref.Project): Target;
    hasProject(): boolean;
    clearProject(): Target;

    getApplication(): Ref.Application | undefined;
    setApplication(value?: Ref.Application): Target;
    hasApplication(): boolean;
    clearApplication(): Target;

    getWorkspace(): Ref.Workspace | undefined;
    setWorkspace(value?: Ref.Workspace): Target;
    hasWorkspace(): boolean;
    clearWorkspace(): Target;

    getLabelSelector(): string;
    setLabelSelector(value: string): Target;

    getRunner(): Ref.Runner | undefined;
    setRunner(value?: Ref.Runner): Target;
    hasRunner(): boolean;
    clearRunner(): Target;

    getAppScopeCase(): Target.AppScopeCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Target.AsObject;
    static toObject(includeInstance: boolean, msg: Target): Target.AsObject;
    static serializeBinaryToWriter(message: Target, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Target;
    static deserializeBinaryFromReader(message: Target, reader: jspb.BinaryReader): Target;
  }

  export namespace Target {
    export type AsObject = {
      global?: Ref.Global.AsObject,
      project?: Ref.Project.AsObject,
      application?: Ref.Application.AsObject,
      workspace?: Ref.Workspace.AsObject,
      labelSelector: string,
      runner?: Ref.Runner.AsObject,
    }

    export enum AppScopeCase { 
      APP_SCOPE_NOT_SET = 0,
      GLOBAL = 1,
      PROJECT = 2,
      APPLICATION = 3,
    }
  }


  export enum ValueCase { 
    VALUE_NOT_SET = 0,
    UNSET = 7,
    STATIC = 2,
    DYNAMIC = 6,
  }

  export enum UnusedScopeCase { 
    UNUSED_SCOPE_NOT_SET = 0,
    APPLICATION = 3,
    PROJECT = 4,
    RUNNER = 5,
  }
}

export class ConfigSetRequest extends jspb.Message {
  getVariablesList(): Array<ConfigVar>;
  setVariablesList(value: Array<ConfigVar>): ConfigSetRequest;
  clearVariablesList(): ConfigSetRequest;
  addVariables(value?: ConfigVar, index?: number): ConfigVar;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigSetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigSetRequest): ConfigSetRequest.AsObject;
  static serializeBinaryToWriter(message: ConfigSetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigSetRequest;
  static deserializeBinaryFromReader(message: ConfigSetRequest, reader: jspb.BinaryReader): ConfigSetRequest;
}

export namespace ConfigSetRequest {
  export type AsObject = {
    variablesList: Array<ConfigVar.AsObject>,
  }
}

export class ConfigSetResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigSetResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigSetResponse): ConfigSetResponse.AsObject;
  static serializeBinaryToWriter(message: ConfigSetResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigSetResponse;
  static deserializeBinaryFromReader(message: ConfigSetResponse, reader: jspb.BinaryReader): ConfigSetResponse;
}

export namespace ConfigSetResponse {
  export type AsObject = {
  }
}

export class ConfigGetRequest extends jspb.Message {
  getApplication(): Ref.Application | undefined;
  setApplication(value?: Ref.Application): ConfigGetRequest;
  hasApplication(): boolean;
  clearApplication(): ConfigGetRequest;

  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): ConfigGetRequest;
  hasProject(): boolean;
  clearProject(): ConfigGetRequest;

  getRunner(): Ref.RunnerId | undefined;
  setRunner(value?: Ref.RunnerId): ConfigGetRequest;
  hasRunner(): boolean;
  clearRunner(): ConfigGetRequest;

  getWorkspace(): Ref.Workspace | undefined;
  setWorkspace(value?: Ref.Workspace): ConfigGetRequest;
  hasWorkspace(): boolean;
  clearWorkspace(): ConfigGetRequest;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): ConfigGetRequest;

  getPrefix(): string;
  setPrefix(value: string): ConfigGetRequest;

  getScopeCase(): ConfigGetRequest.ScopeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigGetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigGetRequest): ConfigGetRequest.AsObject;
  static serializeBinaryToWriter(message: ConfigGetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigGetRequest;
  static deserializeBinaryFromReader(message: ConfigGetRequest, reader: jspb.BinaryReader): ConfigGetRequest;
}

export namespace ConfigGetRequest {
  export type AsObject = {
    application?: Ref.Application.AsObject,
    project?: Ref.Project.AsObject,
    runner?: Ref.RunnerId.AsObject,
    workspace?: Ref.Workspace.AsObject,
    labelsMap: Array<[string, string]>,
    prefix: string,
  }

  export enum ScopeCase { 
    SCOPE_NOT_SET = 0,
    APPLICATION = 2,
    PROJECT = 3,
  }
}

export class ConfigGetResponse extends jspb.Message {
  getVariablesList(): Array<ConfigVar>;
  setVariablesList(value: Array<ConfigVar>): ConfigGetResponse;
  clearVariablesList(): ConfigGetResponse;
  addVariables(value?: ConfigVar, index?: number): ConfigVar;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigGetResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigGetResponse): ConfigGetResponse.AsObject;
  static serializeBinaryToWriter(message: ConfigGetResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigGetResponse;
  static deserializeBinaryFromReader(message: ConfigGetResponse, reader: jspb.BinaryReader): ConfigGetResponse;
}

export namespace ConfigGetResponse {
  export type AsObject = {
    variablesList: Array<ConfigVar.AsObject>,
  }
}

export class ConfigSource extends jspb.Message {
  getDelete(): boolean;
  setDelete(value: boolean): ConfigSource;

  getGlobal(): Ref.Global | undefined;
  setGlobal(value?: Ref.Global): ConfigSource;
  hasGlobal(): boolean;
  clearGlobal(): ConfigSource;

  getType(): string;
  setType(value: string): ConfigSource;

  getConfigMap(): jspb.Map<string, string>;
  clearConfigMap(): ConfigSource;

  getHash(): number;
  setHash(value: number): ConfigSource;

  getScopeCase(): ConfigSource.ScopeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigSource.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigSource): ConfigSource.AsObject;
  static serializeBinaryToWriter(message: ConfigSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigSource;
  static deserializeBinaryFromReader(message: ConfigSource, reader: jspb.BinaryReader): ConfigSource;
}

export namespace ConfigSource {
  export type AsObject = {
    pb_delete: boolean,
    global?: Ref.Global.AsObject,
    type: string,
    configMap: Array<[string, string]>,
    hash: number,
  }

  export enum ScopeCase { 
    SCOPE_NOT_SET = 0,
    GLOBAL = 50,
  }
}

export class SetConfigSourceRequest extends jspb.Message {
  getConfigSource(): ConfigSource | undefined;
  setConfigSource(value?: ConfigSource): SetConfigSourceRequest;
  hasConfigSource(): boolean;
  clearConfigSource(): SetConfigSourceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetConfigSourceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetConfigSourceRequest): SetConfigSourceRequest.AsObject;
  static serializeBinaryToWriter(message: SetConfigSourceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetConfigSourceRequest;
  static deserializeBinaryFromReader(message: SetConfigSourceRequest, reader: jspb.BinaryReader): SetConfigSourceRequest;
}

export namespace SetConfigSourceRequest {
  export type AsObject = {
    configSource?: ConfigSource.AsObject,
  }
}

export class GetConfigSourceRequest extends jspb.Message {
  getGlobal(): Ref.Global | undefined;
  setGlobal(value?: Ref.Global): GetConfigSourceRequest;
  hasGlobal(): boolean;
  clearGlobal(): GetConfigSourceRequest;

  getType(): string;
  setType(value: string): GetConfigSourceRequest;

  getScopeCase(): GetConfigSourceRequest.ScopeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigSourceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigSourceRequest): GetConfigSourceRequest.AsObject;
  static serializeBinaryToWriter(message: GetConfigSourceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigSourceRequest;
  static deserializeBinaryFromReader(message: GetConfigSourceRequest, reader: jspb.BinaryReader): GetConfigSourceRequest;
}

export namespace GetConfigSourceRequest {
  export type AsObject = {
    global?: Ref.Global.AsObject,
    type: string,
  }

  export enum ScopeCase { 
    SCOPE_NOT_SET = 0,
    GLOBAL = 50,
  }
}

export class GetConfigSourceResponse extends jspb.Message {
  getConfigSourcesList(): Array<ConfigSource>;
  setConfigSourcesList(value: Array<ConfigSource>): GetConfigSourceResponse;
  clearConfigSourcesList(): GetConfigSourceResponse;
  addConfigSources(value?: ConfigSource, index?: number): ConfigSource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigSourceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigSourceResponse): GetConfigSourceResponse.AsObject;
  static serializeBinaryToWriter(message: GetConfigSourceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigSourceResponse;
  static deserializeBinaryFromReader(message: GetConfigSourceResponse, reader: jspb.BinaryReader): GetConfigSourceResponse;
}

export namespace GetConfigSourceResponse {
  export type AsObject = {
    configSourcesList: Array<ConfigSource.AsObject>,
  }
}

export class ExecStreamRequest extends jspb.Message {
  getStart(): ExecStreamRequest.Start | undefined;
  setStart(value?: ExecStreamRequest.Start): ExecStreamRequest;
  hasStart(): boolean;
  clearStart(): ExecStreamRequest;

  getInput(): ExecStreamRequest.Input | undefined;
  setInput(value?: ExecStreamRequest.Input): ExecStreamRequest;
  hasInput(): boolean;
  clearInput(): ExecStreamRequest;

  getWinch(): ExecStreamRequest.WindowSize | undefined;
  setWinch(value?: ExecStreamRequest.WindowSize): ExecStreamRequest;
  hasWinch(): boolean;
  clearWinch(): ExecStreamRequest;

  getInputEof(): google_protobuf_empty_pb.Empty | undefined;
  setInputEof(value?: google_protobuf_empty_pb.Empty): ExecStreamRequest;
  hasInputEof(): boolean;
  clearInputEof(): ExecStreamRequest;

  getEventCase(): ExecStreamRequest.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecStreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExecStreamRequest): ExecStreamRequest.AsObject;
  static serializeBinaryToWriter(message: ExecStreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecStreamRequest;
  static deserializeBinaryFromReader(message: ExecStreamRequest, reader: jspb.BinaryReader): ExecStreamRequest;
}

export namespace ExecStreamRequest {
  export type AsObject = {
    start?: ExecStreamRequest.Start.AsObject,
    input?: ExecStreamRequest.Input.AsObject,
    winch?: ExecStreamRequest.WindowSize.AsObject,
    inputEof?: google_protobuf_empty_pb.Empty.AsObject,
  }

  export class Start extends jspb.Message {
    getDeploymentId(): string;
    setDeploymentId(value: string): Start;

    getInstanceId(): string;
    setInstanceId(value: string): Start;

    getArgsList(): Array<string>;
    setArgsList(value: Array<string>): Start;
    clearArgsList(): Start;
    addArgs(value: string, index?: number): Start;

    getPty(): ExecStreamRequest.PTY | undefined;
    setPty(value?: ExecStreamRequest.PTY): Start;
    hasPty(): boolean;
    clearPty(): Start;

    getTargetCase(): Start.TargetCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Start.AsObject;
    static toObject(includeInstance: boolean, msg: Start): Start.AsObject;
    static serializeBinaryToWriter(message: Start, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Start;
    static deserializeBinaryFromReader(message: Start, reader: jspb.BinaryReader): Start;
  }

  export namespace Start {
    export type AsObject = {
      deploymentId: string,
      instanceId: string,
      argsList: Array<string>,
      pty?: ExecStreamRequest.PTY.AsObject,
    }

    export enum TargetCase { 
      TARGET_NOT_SET = 0,
      DEPLOYMENT_ID = 1,
      INSTANCE_ID = 4,
    }
  }


  export class Input extends jspb.Message {
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): Input;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Input.AsObject;
    static toObject(includeInstance: boolean, msg: Input): Input.AsObject;
    static serializeBinaryToWriter(message: Input, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Input;
    static deserializeBinaryFromReader(message: Input, reader: jspb.BinaryReader): Input;
  }

  export namespace Input {
    export type AsObject = {
      data: Uint8Array | string,
    }
  }


  export class PTY extends jspb.Message {
    getEnable(): boolean;
    setEnable(value: boolean): PTY;

    getTerm(): string;
    setTerm(value: string): PTY;

    getWindowSize(): ExecStreamRequest.WindowSize | undefined;
    setWindowSize(value?: ExecStreamRequest.WindowSize): PTY;
    hasWindowSize(): boolean;
    clearWindowSize(): PTY;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PTY.AsObject;
    static toObject(includeInstance: boolean, msg: PTY): PTY.AsObject;
    static serializeBinaryToWriter(message: PTY, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PTY;
    static deserializeBinaryFromReader(message: PTY, reader: jspb.BinaryReader): PTY;
  }

  export namespace PTY {
    export type AsObject = {
      enable: boolean,
      term: string,
      windowSize?: ExecStreamRequest.WindowSize.AsObject,
    }
  }


  export class WindowSize extends jspb.Message {
    getRows(): number;
    setRows(value: number): WindowSize;

    getCols(): number;
    setCols(value: number): WindowSize;

    getWidth(): number;
    setWidth(value: number): WindowSize;

    getHeight(): number;
    setHeight(value: number): WindowSize;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): WindowSize.AsObject;
    static toObject(includeInstance: boolean, msg: WindowSize): WindowSize.AsObject;
    static serializeBinaryToWriter(message: WindowSize, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): WindowSize;
    static deserializeBinaryFromReader(message: WindowSize, reader: jspb.BinaryReader): WindowSize;
  }

  export namespace WindowSize {
    export type AsObject = {
      rows: number,
      cols: number,
      width: number,
      height: number,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    START = 1,
    INPUT = 2,
    WINCH = 3,
    INPUT_EOF = 4,
  }
}

export class ExecStreamResponse extends jspb.Message {
  getOpen(): ExecStreamResponse.Open | undefined;
  setOpen(value?: ExecStreamResponse.Open): ExecStreamResponse;
  hasOpen(): boolean;
  clearOpen(): ExecStreamResponse;

  getOutput(): ExecStreamResponse.Output | undefined;
  setOutput(value?: ExecStreamResponse.Output): ExecStreamResponse;
  hasOutput(): boolean;
  clearOutput(): ExecStreamResponse;

  getExit(): ExecStreamResponse.Exit | undefined;
  setExit(value?: ExecStreamResponse.Exit): ExecStreamResponse;
  hasExit(): boolean;
  clearExit(): ExecStreamResponse;

  getEventCase(): ExecStreamResponse.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecStreamResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExecStreamResponse): ExecStreamResponse.AsObject;
  static serializeBinaryToWriter(message: ExecStreamResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecStreamResponse;
  static deserializeBinaryFromReader(message: ExecStreamResponse, reader: jspb.BinaryReader): ExecStreamResponse;
}

export namespace ExecStreamResponse {
  export type AsObject = {
    open?: ExecStreamResponse.Open.AsObject,
    output?: ExecStreamResponse.Output.AsObject,
    exit?: ExecStreamResponse.Exit.AsObject,
  }

  export class Open extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Open.AsObject;
    static toObject(includeInstance: boolean, msg: Open): Open.AsObject;
    static serializeBinaryToWriter(message: Open, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Open;
    static deserializeBinaryFromReader(message: Open, reader: jspb.BinaryReader): Open;
  }

  export namespace Open {
    export type AsObject = {
    }
  }


  export class Exit extends jspb.Message {
    getCode(): number;
    setCode(value: number): Exit;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Exit.AsObject;
    static toObject(includeInstance: boolean, msg: Exit): Exit.AsObject;
    static serializeBinaryToWriter(message: Exit, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Exit;
    static deserializeBinaryFromReader(message: Exit, reader: jspb.BinaryReader): Exit;
  }

  export namespace Exit {
    export type AsObject = {
      code: number,
    }
  }


  export class Output extends jspb.Message {
    getChannel(): ExecStreamResponse.Output.Channel;
    setChannel(value: ExecStreamResponse.Output.Channel): Output;

    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): Output;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Output.AsObject;
    static toObject(includeInstance: boolean, msg: Output): Output.AsObject;
    static serializeBinaryToWriter(message: Output, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Output;
    static deserializeBinaryFromReader(message: Output, reader: jspb.BinaryReader): Output;
  }

  export namespace Output {
    export type AsObject = {
      channel: ExecStreamResponse.Output.Channel,
      data: Uint8Array | string,
    }

    export enum Channel { 
      UNKNOWN = 0,
      STDOUT = 1,
      STDERR = 2,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    OPEN = 3,
    OUTPUT = 1,
    EXIT = 2,
  }
}

export class EntrypointConfigRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): EntrypointConfigRequest;

  getInstanceId(): string;
  setInstanceId(value: string): EntrypointConfigRequest;

  getType(): Instance.Type;
  setType(value: Instance.Type): EntrypointConfigRequest;

  getDisableExec(): boolean;
  setDisableExec(value: boolean): EntrypointConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EntrypointConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EntrypointConfigRequest): EntrypointConfigRequest.AsObject;
  static serializeBinaryToWriter(message: EntrypointConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EntrypointConfigRequest;
  static deserializeBinaryFromReader(message: EntrypointConfigRequest, reader: jspb.BinaryReader): EntrypointConfigRequest;
}

export namespace EntrypointConfigRequest {
  export type AsObject = {
    deploymentId: string,
    instanceId: string,
    type: Instance.Type,
    disableExec: boolean,
  }
}

export class EntrypointConfigResponse extends jspb.Message {
  getConfig(): EntrypointConfig | undefined;
  setConfig(value?: EntrypointConfig): EntrypointConfigResponse;
  hasConfig(): boolean;
  clearConfig(): EntrypointConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EntrypointConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EntrypointConfigResponse): EntrypointConfigResponse.AsObject;
  static serializeBinaryToWriter(message: EntrypointConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EntrypointConfigResponse;
  static deserializeBinaryFromReader(message: EntrypointConfigResponse, reader: jspb.BinaryReader): EntrypointConfigResponse;
}

export namespace EntrypointConfigResponse {
  export type AsObject = {
    config?: EntrypointConfig.AsObject,
  }
}

export class EntrypointConfig extends jspb.Message {
  getExecList(): Array<EntrypointConfig.Exec>;
  setExecList(value: Array<EntrypointConfig.Exec>): EntrypointConfig;
  clearExecList(): EntrypointConfig;
  addExec(value?: EntrypointConfig.Exec, index?: number): EntrypointConfig.Exec;

  getEnvVarsList(): Array<ConfigVar>;
  setEnvVarsList(value: Array<ConfigVar>): EntrypointConfig;
  clearEnvVarsList(): EntrypointConfig;
  addEnvVars(value?: ConfigVar, index?: number): ConfigVar;

  getConfigSourcesList(): Array<ConfigSource>;
  setConfigSourcesList(value: Array<ConfigSource>): EntrypointConfig;
  clearConfigSourcesList(): EntrypointConfig;
  addConfigSources(value?: ConfigSource, index?: number): ConfigSource;

  getFileChangeSignal(): string;
  setFileChangeSignal(value: string): EntrypointConfig;

  getUrlService(): EntrypointConfig.URLService | undefined;
  setUrlService(value?: EntrypointConfig.URLService): EntrypointConfig;
  hasUrlService(): boolean;
  clearUrlService(): EntrypointConfig;

  getDeployment(): EntrypointConfig.DeploymentInfo | undefined;
  setDeployment(value?: EntrypointConfig.DeploymentInfo): EntrypointConfig;
  hasDeployment(): boolean;
  clearDeployment(): EntrypointConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EntrypointConfig.AsObject;
  static toObject(includeInstance: boolean, msg: EntrypointConfig): EntrypointConfig.AsObject;
  static serializeBinaryToWriter(message: EntrypointConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EntrypointConfig;
  static deserializeBinaryFromReader(message: EntrypointConfig, reader: jspb.BinaryReader): EntrypointConfig;
}

export namespace EntrypointConfig {
  export type AsObject = {
    execList: Array<EntrypointConfig.Exec.AsObject>,
    envVarsList: Array<ConfigVar.AsObject>,
    configSourcesList: Array<ConfigSource.AsObject>,
    fileChangeSignal: string,
    urlService?: EntrypointConfig.URLService.AsObject,
    deployment?: EntrypointConfig.DeploymentInfo.AsObject,
  }

  export class Exec extends jspb.Message {
    getIndex(): number;
    setIndex(value: number): Exec;

    getArgsList(): Array<string>;
    setArgsList(value: Array<string>): Exec;
    clearArgsList(): Exec;
    addArgs(value: string, index?: number): Exec;

    getPty(): ExecStreamRequest.PTY | undefined;
    setPty(value?: ExecStreamRequest.PTY): Exec;
    hasPty(): boolean;
    clearPty(): Exec;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Exec.AsObject;
    static toObject(includeInstance: boolean, msg: Exec): Exec.AsObject;
    static serializeBinaryToWriter(message: Exec, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Exec;
    static deserializeBinaryFromReader(message: Exec, reader: jspb.BinaryReader): Exec;
  }

  export namespace Exec {
    export type AsObject = {
      index: number,
      argsList: Array<string>,
      pty?: ExecStreamRequest.PTY.AsObject,
    }
  }


  export class URLService extends jspb.Message {
    getControlAddr(): string;
    setControlAddr(value: string): URLService;

    getToken(): string;
    setToken(value: string): URLService;

    getLabels(): string;
    setLabels(value: string): URLService;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): URLService.AsObject;
    static toObject(includeInstance: boolean, msg: URLService): URLService.AsObject;
    static serializeBinaryToWriter(message: URLService, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): URLService;
    static deserializeBinaryFromReader(message: URLService, reader: jspb.BinaryReader): URLService;
  }

  export namespace URLService {
    export type AsObject = {
      controlAddr: string,
      token: string,
      labels: string,
    }
  }


  export class DeploymentInfo extends jspb.Message {
    getComponent(): Component | undefined;
    setComponent(value?: Component): DeploymentInfo;
    hasComponent(): boolean;
    clearComponent(): DeploymentInfo;

    getLabelsMap(): jspb.Map<string, string>;
    clearLabelsMap(): DeploymentInfo;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeploymentInfo.AsObject;
    static toObject(includeInstance: boolean, msg: DeploymentInfo): DeploymentInfo.AsObject;
    static serializeBinaryToWriter(message: DeploymentInfo, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeploymentInfo;
    static deserializeBinaryFromReader(message: DeploymentInfo, reader: jspb.BinaryReader): DeploymentInfo;
  }

  export namespace DeploymentInfo {
    export type AsObject = {
      component?: Component.AsObject,
      labelsMap: Array<[string, string]>,
    }
  }

}

export class EntrypointLogBatch extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): EntrypointLogBatch;

  getLinesList(): Array<LogBatch.Entry>;
  setLinesList(value: Array<LogBatch.Entry>): EntrypointLogBatch;
  clearLinesList(): EntrypointLogBatch;
  addLines(value?: LogBatch.Entry, index?: number): LogBatch.Entry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EntrypointLogBatch.AsObject;
  static toObject(includeInstance: boolean, msg: EntrypointLogBatch): EntrypointLogBatch.AsObject;
  static serializeBinaryToWriter(message: EntrypointLogBatch, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EntrypointLogBatch;
  static deserializeBinaryFromReader(message: EntrypointLogBatch, reader: jspb.BinaryReader): EntrypointLogBatch;
}

export namespace EntrypointLogBatch {
  export type AsObject = {
    instanceId: string,
    linesList: Array<LogBatch.Entry.AsObject>,
  }
}

export class EntrypointExecRequest extends jspb.Message {
  getOpen(): EntrypointExecRequest.Open | undefined;
  setOpen(value?: EntrypointExecRequest.Open): EntrypointExecRequest;
  hasOpen(): boolean;
  clearOpen(): EntrypointExecRequest;

  getExit(): EntrypointExecRequest.Exit | undefined;
  setExit(value?: EntrypointExecRequest.Exit): EntrypointExecRequest;
  hasExit(): boolean;
  clearExit(): EntrypointExecRequest;

  getOutput(): EntrypointExecRequest.Output | undefined;
  setOutput(value?: EntrypointExecRequest.Output): EntrypointExecRequest;
  hasOutput(): boolean;
  clearOutput(): EntrypointExecRequest;

  getError(): EntrypointExecRequest.Error | undefined;
  setError(value?: EntrypointExecRequest.Error): EntrypointExecRequest;
  hasError(): boolean;
  clearError(): EntrypointExecRequest;

  getEventCase(): EntrypointExecRequest.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EntrypointExecRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EntrypointExecRequest): EntrypointExecRequest.AsObject;
  static serializeBinaryToWriter(message: EntrypointExecRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EntrypointExecRequest;
  static deserializeBinaryFromReader(message: EntrypointExecRequest, reader: jspb.BinaryReader): EntrypointExecRequest;
}

export namespace EntrypointExecRequest {
  export type AsObject = {
    open?: EntrypointExecRequest.Open.AsObject,
    exit?: EntrypointExecRequest.Exit.AsObject,
    output?: EntrypointExecRequest.Output.AsObject,
    error?: EntrypointExecRequest.Error.AsObject,
  }

  export class Open extends jspb.Message {
    getInstanceId(): string;
    setInstanceId(value: string): Open;

    getIndex(): number;
    setIndex(value: number): Open;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Open.AsObject;
    static toObject(includeInstance: boolean, msg: Open): Open.AsObject;
    static serializeBinaryToWriter(message: Open, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Open;
    static deserializeBinaryFromReader(message: Open, reader: jspb.BinaryReader): Open;
  }

  export namespace Open {
    export type AsObject = {
      instanceId: string,
      index: number,
    }
  }


  export class Exit extends jspb.Message {
    getCode(): number;
    setCode(value: number): Exit;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Exit.AsObject;
    static toObject(includeInstance: boolean, msg: Exit): Exit.AsObject;
    static serializeBinaryToWriter(message: Exit, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Exit;
    static deserializeBinaryFromReader(message: Exit, reader: jspb.BinaryReader): Exit;
  }

  export namespace Exit {
    export type AsObject = {
      code: number,
    }
  }


  export class Output extends jspb.Message {
    getChannel(): EntrypointExecRequest.Output.Channel;
    setChannel(value: EntrypointExecRequest.Output.Channel): Output;

    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): Output;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Output.AsObject;
    static toObject(includeInstance: boolean, msg: Output): Output.AsObject;
    static serializeBinaryToWriter(message: Output, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Output;
    static deserializeBinaryFromReader(message: Output, reader: jspb.BinaryReader): Output;
  }

  export namespace Output {
    export type AsObject = {
      channel: EntrypointExecRequest.Output.Channel,
      data: Uint8Array | string,
    }

    export enum Channel { 
      UNKNOWN = 0,
      STDOUT = 1,
      STDERR = 2,
    }
  }


  export class Error extends jspb.Message {
    getError(): google_rpc_status_pb.Status | undefined;
    setError(value?: google_rpc_status_pb.Status): Error;
    hasError(): boolean;
    clearError(): Error;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Error.AsObject;
    static toObject(includeInstance: boolean, msg: Error): Error.AsObject;
    static serializeBinaryToWriter(message: Error, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Error;
    static deserializeBinaryFromReader(message: Error, reader: jspb.BinaryReader): Error;
  }

  export namespace Error {
    export type AsObject = {
      error?: google_rpc_status_pb.Status.AsObject,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    OPEN = 1,
    EXIT = 2,
    OUTPUT = 3,
    ERROR = 4,
  }
}

export class EntrypointExecResponse extends jspb.Message {
  getInput(): Uint8Array | string;
  getInput_asU8(): Uint8Array;
  getInput_asB64(): string;
  setInput(value: Uint8Array | string): EntrypointExecResponse;

  getInputEof(): google_protobuf_empty_pb.Empty | undefined;
  setInputEof(value?: google_protobuf_empty_pb.Empty): EntrypointExecResponse;
  hasInputEof(): boolean;
  clearInputEof(): EntrypointExecResponse;

  getWinch(): ExecStreamRequest.WindowSize | undefined;
  setWinch(value?: ExecStreamRequest.WindowSize): EntrypointExecResponse;
  hasWinch(): boolean;
  clearWinch(): EntrypointExecResponse;

  getOpened(): boolean;
  setOpened(value: boolean): EntrypointExecResponse;

  getEventCase(): EntrypointExecResponse.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EntrypointExecResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EntrypointExecResponse): EntrypointExecResponse.AsObject;
  static serializeBinaryToWriter(message: EntrypointExecResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EntrypointExecResponse;
  static deserializeBinaryFromReader(message: EntrypointExecResponse, reader: jspb.BinaryReader): EntrypointExecResponse;
}

export namespace EntrypointExecResponse {
  export type AsObject = {
    input: Uint8Array | string,
    inputEof?: google_protobuf_empty_pb.Empty.AsObject,
    winch?: ExecStreamRequest.WindowSize.AsObject,
    opened: boolean,
  }

  export enum EventCase { 
    EVENT_NOT_SET = 0,
    INPUT = 1,
    INPUT_EOF = 4,
    WINCH = 2,
    OPENED = 3,
  }
}

export class Pipeline extends jspb.Message {
  getId(): string;
  setId(value: string): Pipeline;

  getName(): string;
  setName(value: string): Pipeline;

  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): Pipeline;
  hasProject(): boolean;
  clearProject(): Pipeline;

  getStepsMap(): jspb.Map<string, Pipeline.Step>;
  clearStepsMap(): Pipeline;

  getOwnerCase(): Pipeline.OwnerCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Pipeline.AsObject;
  static toObject(includeInstance: boolean, msg: Pipeline): Pipeline.AsObject;
  static serializeBinaryToWriter(message: Pipeline, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Pipeline;
  static deserializeBinaryFromReader(message: Pipeline, reader: jspb.BinaryReader): Pipeline;
}

export namespace Pipeline {
  export type AsObject = {
    id: string,
    name: string,
    project?: Ref.Project.AsObject,
    stepsMap: Array<[string, Pipeline.Step.AsObject]>,
  }

  export class Step extends jspb.Message {
    getName(): string;
    setName(value: string): Step;

    getDependsOnList(): Array<string>;
    setDependsOnList(value: Array<string>): Step;
    clearDependsOnList(): Step;
    addDependsOn(value: string, index?: number): Step;

    getExec(): Pipeline.Step.Exec | undefined;
    setExec(value?: Pipeline.Step.Exec): Step;
    hasExec(): boolean;
    clearExec(): Step;

    getBuild(): Pipeline.Step.Build | undefined;
    setBuild(value?: Pipeline.Step.Build): Step;
    hasBuild(): boolean;
    clearBuild(): Step;

    getDeploy(): Pipeline.Step.Deploy | undefined;
    setDeploy(value?: Pipeline.Step.Deploy): Step;
    hasDeploy(): boolean;
    clearDeploy(): Step;

    getRelease(): Pipeline.Step.Release | undefined;
    setRelease(value?: Pipeline.Step.Release): Step;
    hasRelease(): boolean;
    clearRelease(): Step;

    getUp(): Pipeline.Step.Up | undefined;
    setUp(value?: Pipeline.Step.Up): Step;
    hasUp(): boolean;
    clearUp(): Step;

    getImage(): string;
    setImage(value: string): Step;

    getKindCase(): Step.KindCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Step.AsObject;
    static toObject(includeInstance: boolean, msg: Step): Step.AsObject;
    static serializeBinaryToWriter(message: Step, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Step;
    static deserializeBinaryFromReader(message: Step, reader: jspb.BinaryReader): Step;
  }

  export namespace Step {
    export type AsObject = {
      name: string,
      dependsOnList: Array<string>,
      exec?: Pipeline.Step.Exec.AsObject,
      build?: Pipeline.Step.Build.AsObject,
      deploy?: Pipeline.Step.Deploy.AsObject,
      release?: Pipeline.Step.Release.AsObject,
      up?: Pipeline.Step.Up.AsObject,
      image: string,
    }

    export class Exec extends jspb.Message {
      getImage(): string;
      setImage(value: string): Exec;

      getCommand(): string;
      setCommand(value: string): Exec;

      getArgsList(): Array<string>;
      setArgsList(value: Array<string>): Exec;
      clearArgsList(): Exec;
      addArgs(value: string, index?: number): Exec;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Exec.AsObject;
      static toObject(includeInstance: boolean, msg: Exec): Exec.AsObject;
      static serializeBinaryToWriter(message: Exec, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Exec;
      static deserializeBinaryFromReader(message: Exec, reader: jspb.BinaryReader): Exec;
    }

    export namespace Exec {
      export type AsObject = {
        image: string,
        command: string,
        argsList: Array<string>,
      }
    }


    export class Build extends jspb.Message {
      getDisablePush(): boolean;
      setDisablePush(value: boolean): Build;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Build.AsObject;
      static toObject(includeInstance: boolean, msg: Build): Build.AsObject;
      static serializeBinaryToWriter(message: Build, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Build;
      static deserializeBinaryFromReader(message: Build, reader: jspb.BinaryReader): Build;
    }

    export namespace Build {
      export type AsObject = {
        disablePush: boolean,
      }
    }


    export class Deploy extends jspb.Message {
      getRelease(): boolean;
      setRelease(value: boolean): Deploy;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Deploy.AsObject;
      static toObject(includeInstance: boolean, msg: Deploy): Deploy.AsObject;
      static serializeBinaryToWriter(message: Deploy, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Deploy;
      static deserializeBinaryFromReader(message: Deploy, reader: jspb.BinaryReader): Deploy;
    }

    export namespace Deploy {
      export type AsObject = {
        release: boolean,
      }
    }


    export class Release extends jspb.Message {
      getDeployment(): Ref.Deployment | undefined;
      setDeployment(value?: Ref.Deployment): Release;
      hasDeployment(): boolean;
      clearDeployment(): Release;

      getPrune(): boolean;
      setPrune(value: boolean): Release;

      getPruneRetain(): number;
      setPruneRetain(value: number): Release;

      getPruneRetainOverride(): boolean;
      setPruneRetainOverride(value: boolean): Release;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Release.AsObject;
      static toObject(includeInstance: boolean, msg: Release): Release.AsObject;
      static serializeBinaryToWriter(message: Release, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Release;
      static deserializeBinaryFromReader(message: Release, reader: jspb.BinaryReader): Release;
    }

    export namespace Release {
      export type AsObject = {
        deployment?: Ref.Deployment.AsObject,
        prune: boolean,
        pruneRetain: number,
        pruneRetainOverride: boolean,
      }
    }


    export class Up extends jspb.Message {
      getPrune(): boolean;
      setPrune(value: boolean): Up;

      getPruneRetain(): number;
      setPruneRetain(value: number): Up;

      getPruneRetainOverride(): boolean;
      setPruneRetainOverride(value: boolean): Up;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Up.AsObject;
      static toObject(includeInstance: boolean, msg: Up): Up.AsObject;
      static serializeBinaryToWriter(message: Up, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Up;
      static deserializeBinaryFromReader(message: Up, reader: jspb.BinaryReader): Up;
    }

    export namespace Up {
      export type AsObject = {
        prune: boolean,
        pruneRetain: number,
        pruneRetainOverride: boolean,
      }
    }


    export enum KindCase { 
      KIND_NOT_SET = 0,
      EXEC = 3,
      BUILD = 4,
      DEPLOY = 6,
      RELEASE = 7,
      UP = 8,
    }
  }


  export enum OwnerCase { 
    OWNER_NOT_SET = 0,
    PROJECT = 3,
  }
}

export class UpsertPipelineRequest extends jspb.Message {
  getPipeline(): Pipeline | undefined;
  setPipeline(value?: Pipeline): UpsertPipelineRequest;
  hasPipeline(): boolean;
  clearPipeline(): UpsertPipelineRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertPipelineRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertPipelineRequest): UpsertPipelineRequest.AsObject;
  static serializeBinaryToWriter(message: UpsertPipelineRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertPipelineRequest;
  static deserializeBinaryFromReader(message: UpsertPipelineRequest, reader: jspb.BinaryReader): UpsertPipelineRequest;
}

export namespace UpsertPipelineRequest {
  export type AsObject = {
    pipeline?: Pipeline.AsObject,
  }
}

export class UpsertPipelineResponse extends jspb.Message {
  getPipeline(): Pipeline | undefined;
  setPipeline(value?: Pipeline): UpsertPipelineResponse;
  hasPipeline(): boolean;
  clearPipeline(): UpsertPipelineResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpsertPipelineResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpsertPipelineResponse): UpsertPipelineResponse.AsObject;
  static serializeBinaryToWriter(message: UpsertPipelineResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpsertPipelineResponse;
  static deserializeBinaryFromReader(message: UpsertPipelineResponse, reader: jspb.BinaryReader): UpsertPipelineResponse;
}

export namespace UpsertPipelineResponse {
  export type AsObject = {
    pipeline?: Pipeline.AsObject,
  }
}

export class GetPipelineRequest extends jspb.Message {
  getPipeline(): Ref.Pipeline | undefined;
  setPipeline(value?: Ref.Pipeline): GetPipelineRequest;
  hasPipeline(): boolean;
  clearPipeline(): GetPipelineRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPipelineRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPipelineRequest): GetPipelineRequest.AsObject;
  static serializeBinaryToWriter(message: GetPipelineRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPipelineRequest;
  static deserializeBinaryFromReader(message: GetPipelineRequest, reader: jspb.BinaryReader): GetPipelineRequest;
}

export namespace GetPipelineRequest {
  export type AsObject = {
    pipeline?: Ref.Pipeline.AsObject,
  }
}

export class GetPipelineResponse extends jspb.Message {
  getPipeline(): Pipeline | undefined;
  setPipeline(value?: Pipeline): GetPipelineResponse;
  hasPipeline(): boolean;
  clearPipeline(): GetPipelineResponse;

  getRootStep(): string;
  setRootStep(value: string): GetPipelineResponse;

  getGraph(): GetPipelineResponse.Graph | undefined;
  setGraph(value?: GetPipelineResponse.Graph): GetPipelineResponse;
  hasGraph(): boolean;
  clearGraph(): GetPipelineResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPipelineResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPipelineResponse): GetPipelineResponse.AsObject;
  static serializeBinaryToWriter(message: GetPipelineResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPipelineResponse;
  static deserializeBinaryFromReader(message: GetPipelineResponse, reader: jspb.BinaryReader): GetPipelineResponse;
}

export namespace GetPipelineResponse {
  export type AsObject = {
    pipeline?: Pipeline.AsObject,
    rootStep: string,
    graph?: GetPipelineResponse.Graph.AsObject,
  }

  export class Graph extends jspb.Message {
    getContent(): Uint8Array | string;
    getContent_asU8(): Uint8Array;
    getContent_asB64(): string;
    setContent(value: Uint8Array | string): Graph;

    getFormat(): GetPipelineResponse.Graph.Format;
    setFormat(value: GetPipelineResponse.Graph.Format): Graph;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Graph.AsObject;
    static toObject(includeInstance: boolean, msg: Graph): Graph.AsObject;
    static serializeBinaryToWriter(message: Graph, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Graph;
    static deserializeBinaryFromReader(message: Graph, reader: jspb.BinaryReader): Graph;
  }

  export namespace Graph {
    export type AsObject = {
      content: Uint8Array | string,
      format: GetPipelineResponse.Graph.Format,
    }

    export enum Format { 
      INVALID = 0,
      MERMAID = 1,
    }
  }

}

export class RunPipelineRequest extends jspb.Message {
  getPipeline(): Ref.Pipeline | undefined;
  setPipeline(value?: Ref.Pipeline): RunPipelineRequest;
  hasPipeline(): boolean;
  clearPipeline(): RunPipelineRequest;

  getJobTemplate(): Job | undefined;
  setJobTemplate(value?: Job): RunPipelineRequest;
  hasJobTemplate(): boolean;
  clearJobTemplate(): RunPipelineRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunPipelineRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RunPipelineRequest): RunPipelineRequest.AsObject;
  static serializeBinaryToWriter(message: RunPipelineRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunPipelineRequest;
  static deserializeBinaryFromReader(message: RunPipelineRequest, reader: jspb.BinaryReader): RunPipelineRequest;
}

export namespace RunPipelineRequest {
  export type AsObject = {
    pipeline?: Ref.Pipeline.AsObject,
    jobTemplate?: Job.AsObject,
  }
}

export class RunPipelineResponse extends jspb.Message {
  getJobId(): string;
  setJobId(value: string): RunPipelineResponse;

  getAllJobIdsList(): Array<string>;
  setAllJobIdsList(value: Array<string>): RunPipelineResponse;
  clearAllJobIdsList(): RunPipelineResponse;
  addAllJobIds(value: string, index?: number): RunPipelineResponse;

  getJobMapMap(): jspb.Map<string, Ref.PipelineStep>;
  clearJobMapMap(): RunPipelineResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RunPipelineResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RunPipelineResponse): RunPipelineResponse.AsObject;
  static serializeBinaryToWriter(message: RunPipelineResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RunPipelineResponse;
  static deserializeBinaryFromReader(message: RunPipelineResponse, reader: jspb.BinaryReader): RunPipelineResponse;
}

export namespace RunPipelineResponse {
  export type AsObject = {
    jobId: string,
    allJobIdsList: Array<string>,
    jobMapMap: Array<[string, Ref.PipelineStep.AsObject]>,
  }
}

export class ListPipelinesRequest extends jspb.Message {
  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): ListPipelinesRequest;
  hasProject(): boolean;
  clearProject(): ListPipelinesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPipelinesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListPipelinesRequest): ListPipelinesRequest.AsObject;
  static serializeBinaryToWriter(message: ListPipelinesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPipelinesRequest;
  static deserializeBinaryFromReader(message: ListPipelinesRequest, reader: jspb.BinaryReader): ListPipelinesRequest;
}

export namespace ListPipelinesRequest {
  export type AsObject = {
    project?: Ref.Project.AsObject,
  }
}

export class ListPipelinesResponse extends jspb.Message {
  getPipelinesList(): Array<Pipeline>;
  setPipelinesList(value: Array<Pipeline>): ListPipelinesResponse;
  clearPipelinesList(): ListPipelinesResponse;
  addPipelines(value?: Pipeline, index?: number): Pipeline;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPipelinesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListPipelinesResponse): ListPipelinesResponse.AsObject;
  static serializeBinaryToWriter(message: ListPipelinesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPipelinesResponse;
  static deserializeBinaryFromReader(message: ListPipelinesResponse, reader: jspb.BinaryReader): ListPipelinesResponse;
}

export namespace ListPipelinesResponse {
  export type AsObject = {
    pipelinesList: Array<Pipeline.AsObject>,
  }
}

export class ConfigSyncPipelineRequest extends jspb.Message {
  getProject(): Ref.Project | undefined;
  setProject(value?: Ref.Project): ConfigSyncPipelineRequest;
  hasProject(): boolean;
  clearProject(): ConfigSyncPipelineRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigSyncPipelineRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigSyncPipelineRequest): ConfigSyncPipelineRequest.AsObject;
  static serializeBinaryToWriter(message: ConfigSyncPipelineRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigSyncPipelineRequest;
  static deserializeBinaryFromReader(message: ConfigSyncPipelineRequest, reader: jspb.BinaryReader): ConfigSyncPipelineRequest;
}

export namespace ConfigSyncPipelineRequest {
  export type AsObject = {
    project?: Ref.Project.AsObject,
  }
}

export class ConfigSyncPipelineResponse extends jspb.Message {
  getSyncedPipelinesMap(): jspb.Map<string, Ref.Pipeline>;
  clearSyncedPipelinesMap(): ConfigSyncPipelineResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConfigSyncPipelineResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ConfigSyncPipelineResponse): ConfigSyncPipelineResponse.AsObject;
  static serializeBinaryToWriter(message: ConfigSyncPipelineResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConfigSyncPipelineResponse;
  static deserializeBinaryFromReader(message: ConfigSyncPipelineResponse, reader: jspb.BinaryReader): ConfigSyncPipelineResponse;
}

export namespace ConfigSyncPipelineResponse {
  export type AsObject = {
    syncedPipelinesMap: Array<[string, Ref.Pipeline.AsObject]>,
  }
}

export class TokenTransport extends jspb.Message {
  getBody(): Uint8Array | string;
  getBody_asU8(): Uint8Array;
  getBody_asB64(): string;
  setBody(value: Uint8Array | string): TokenTransport;

  getSignature(): Uint8Array | string;
  getSignature_asU8(): Uint8Array;
  getSignature_asB64(): string;
  setSignature(value: Uint8Array | string): TokenTransport;

  getKeyId(): string;
  setKeyId(value: string): TokenTransport;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): TokenTransport;

  getOauthCreds(): TokenTransport.OAuthCredentials | undefined;
  setOauthCreds(value?: TokenTransport.OAuthCredentials): TokenTransport;
  hasOauthCreds(): boolean;
  clearOauthCreds(): TokenTransport;

  getExternalCredsCase(): TokenTransport.ExternalCredsCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TokenTransport.AsObject;
  static toObject(includeInstance: boolean, msg: TokenTransport): TokenTransport.AsObject;
  static serializeBinaryToWriter(message: TokenTransport, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TokenTransport;
  static deserializeBinaryFromReader(message: TokenTransport, reader: jspb.BinaryReader): TokenTransport;
}

export namespace TokenTransport {
  export type AsObject = {
    body: Uint8Array | string,
    signature: Uint8Array | string,
    keyId: string,
    metadataMap: Array<[string, string]>,
    oauthCreds?: TokenTransport.OAuthCredentials.AsObject,
  }

  export class OAuthCredentials extends jspb.Message {
    getUrl(): string;
    setUrl(value: string): OAuthCredentials;

    getClientId(): string;
    setClientId(value: string): OAuthCredentials;

    getClientSecret(): string;
    setClientSecret(value: string): OAuthCredentials;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OAuthCredentials.AsObject;
    static toObject(includeInstance: boolean, msg: OAuthCredentials): OAuthCredentials.AsObject;
    static serializeBinaryToWriter(message: OAuthCredentials, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OAuthCredentials;
    static deserializeBinaryFromReader(message: OAuthCredentials, reader: jspb.BinaryReader): OAuthCredentials;
  }

  export namespace OAuthCredentials {
    export type AsObject = {
      url: string,
      clientId: string,
      clientSecret: string,
    }
  }


  export enum ExternalCredsCase { 
    EXTERNAL_CREDS_NOT_SET = 0,
    OAUTH_CREDS = 12,
  }
}

export class Token extends jspb.Message {
  getAccessorId(): Uint8Array | string;
  getAccessorId_asU8(): Uint8Array;
  getAccessorId_asB64(): string;
  setAccessorId(value: Uint8Array | string): Token;

  getValidUntil(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setValidUntil(value?: google_protobuf_timestamp_pb.Timestamp): Token;
  hasValidUntil(): boolean;
  clearValidUntil(): Token;

  getIssuedTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setIssuedTime(value?: google_protobuf_timestamp_pb.Timestamp): Token;
  hasIssuedTime(): boolean;
  clearIssuedTime(): Token;

  getLogin(): Token.Login | undefined;
  setLogin(value?: Token.Login): Token;
  hasLogin(): boolean;
  clearLogin(): Token;

  getRunner(): Token.Runner | undefined;
  setRunner(value?: Token.Runner): Token;
  hasRunner(): boolean;
  clearRunner(): Token;

  getInvite(): Token.Invite | undefined;
  setInvite(value?: Token.Invite): Token;
  hasInvite(): boolean;
  clearInvite(): Token;

  getTrigger(): Token.Trigger | undefined;
  setTrigger(value?: Token.Trigger): Token;
  hasTrigger(): boolean;
  clearTrigger(): Token;

  getUnusedUser(): string;
  setUnusedUser(value: string): Token;

  getUnusedLogin(): boolean;
  setUnusedLogin(value: boolean): Token;

  getUnusedInvite(): boolean;
  setUnusedInvite(value: boolean): Token;

  getUnusedEntrypoint(): Token.Entrypoint | undefined;
  setUnusedEntrypoint(value?: Token.Entrypoint): Token;
  hasUnusedEntrypoint(): boolean;
  clearUnusedEntrypoint(): Token;

  getKindCase(): Token.KindCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Token.AsObject;
  static toObject(includeInstance: boolean, msg: Token): Token.AsObject;
  static serializeBinaryToWriter(message: Token, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Token;
  static deserializeBinaryFromReader(message: Token, reader: jspb.BinaryReader): Token;
}

export namespace Token {
  export type AsObject = {
    accessorId: Uint8Array | string,
    validUntil?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    issuedTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    login?: Token.Login.AsObject,
    runner?: Token.Runner.AsObject,
    invite?: Token.Invite.AsObject,
    trigger?: Token.Trigger.AsObject,
    unusedUser: string,
    unusedLogin: boolean,
    unusedInvite: boolean,
    unusedEntrypoint?: Token.Entrypoint.AsObject,
  }

  export class Login extends jspb.Message {
    getUserId(): string;
    setUserId(value: string): Login;

    getEntrypoint(): Token.Entrypoint | undefined;
    setEntrypoint(value?: Token.Entrypoint): Login;
    hasEntrypoint(): boolean;
    clearEntrypoint(): Login;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Login.AsObject;
    static toObject(includeInstance: boolean, msg: Login): Login.AsObject;
    static serializeBinaryToWriter(message: Login, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Login;
    static deserializeBinaryFromReader(message: Login, reader: jspb.BinaryReader): Login;
  }

  export namespace Login {
    export type AsObject = {
      userId: string,
      entrypoint?: Token.Entrypoint.AsObject,
    }
  }


  export class Runner extends jspb.Message {
    getId(): string;
    setId(value: string): Runner;

    getLabelHash(): number;
    setLabelHash(value: number): Runner;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Runner.AsObject;
    static toObject(includeInstance: boolean, msg: Runner): Runner.AsObject;
    static serializeBinaryToWriter(message: Runner, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Runner;
    static deserializeBinaryFromReader(message: Runner, reader: jspb.BinaryReader): Runner;
  }

  export namespace Runner {
    export type AsObject = {
      id: string,
      labelHash: number,
    }
  }


  export class Invite extends jspb.Message {
    getFromUserId(): string;
    setFromUserId(value: string): Invite;

    getLogin(): Token.Login | undefined;
    setLogin(value?: Token.Login): Invite;
    hasLogin(): boolean;
    clearLogin(): Invite;

    getSignup(): Token.Invite.Signup | undefined;
    setSignup(value?: Token.Invite.Signup): Invite;
    hasSignup(): boolean;
    clearSignup(): Invite;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Invite.AsObject;
    static toObject(includeInstance: boolean, msg: Invite): Invite.AsObject;
    static serializeBinaryToWriter(message: Invite, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Invite;
    static deserializeBinaryFromReader(message: Invite, reader: jspb.BinaryReader): Invite;
  }

  export namespace Invite {
    export type AsObject = {
      fromUserId: string,
      login?: Token.Login.AsObject,
      signup?: Token.Invite.Signup.AsObject,
    }

    export class Signup extends jspb.Message {
      getInitialUsername(): string;
      setInitialUsername(value: string): Signup;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Signup.AsObject;
      static toObject(includeInstance: boolean, msg: Signup): Signup.AsObject;
      static serializeBinaryToWriter(message: Signup, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Signup;
      static deserializeBinaryFromReader(message: Signup, reader: jspb.BinaryReader): Signup;
    }

    export namespace Signup {
      export type AsObject = {
        initialUsername: string,
      }
    }

  }


  export class Entrypoint extends jspb.Message {
    getDeploymentId(): string;
    setDeploymentId(value: string): Entrypoint;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Entrypoint.AsObject;
    static toObject(includeInstance: boolean, msg: Entrypoint): Entrypoint.AsObject;
    static serializeBinaryToWriter(message: Entrypoint, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Entrypoint;
    static deserializeBinaryFromReader(message: Entrypoint, reader: jspb.BinaryReader): Entrypoint;
  }

  export namespace Entrypoint {
    export type AsObject = {
      deploymentId: string,
    }
  }


  export class Trigger extends jspb.Message {
    getFromUserId(): string;
    setFromUserId(value: string): Trigger;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Trigger.AsObject;
    static toObject(includeInstance: boolean, msg: Trigger): Trigger.AsObject;
    static serializeBinaryToWriter(message: Trigger, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Trigger;
    static deserializeBinaryFromReader(message: Trigger, reader: jspb.BinaryReader): Trigger;
  }

  export namespace Trigger {
    export type AsObject = {
      fromUserId: string,
    }
  }


  export enum KindCase { 
    KIND_NOT_SET = 0,
    LOGIN = 8,
    RUNNER = 11,
    INVITE = 9,
    TRIGGER = 10,
  }
}

export class HMACKey extends jspb.Message {
  getId(): string;
  setId(value: string): HMACKey;

  getKey(): Uint8Array | string;
  getKey_asU8(): Uint8Array;
  getKey_asB64(): string;
  setKey(value: Uint8Array | string): HMACKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HMACKey.AsObject;
  static toObject(includeInstance: boolean, msg: HMACKey): HMACKey.AsObject;
  static serializeBinaryToWriter(message: HMACKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HMACKey;
  static deserializeBinaryFromReader(message: HMACKey, reader: jspb.BinaryReader): HMACKey;
}

export namespace HMACKey {
  export type AsObject = {
    id: string,
    key: Uint8Array | string,
  }
}

export class DecodeTokenRequest extends jspb.Message {
  getToken(): string;
  setToken(value: string): DecodeTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DecodeTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DecodeTokenRequest): DecodeTokenRequest.AsObject;
  static serializeBinaryToWriter(message: DecodeTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DecodeTokenRequest;
  static deserializeBinaryFromReader(message: DecodeTokenRequest, reader: jspb.BinaryReader): DecodeTokenRequest;
}

export namespace DecodeTokenRequest {
  export type AsObject = {
    token: string,
  }
}

export class DecodeTokenResponse extends jspb.Message {
  getToken(): Token | undefined;
  setToken(value?: Token): DecodeTokenResponse;
  hasToken(): boolean;
  clearToken(): DecodeTokenResponse;

  getTransport(): TokenTransport | undefined;
  setTransport(value?: TokenTransport): DecodeTokenResponse;
  hasTransport(): boolean;
  clearTransport(): DecodeTokenResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DecodeTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DecodeTokenResponse): DecodeTokenResponse.AsObject;
  static serializeBinaryToWriter(message: DecodeTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DecodeTokenResponse;
  static deserializeBinaryFromReader(message: DecodeTokenResponse, reader: jspb.BinaryReader): DecodeTokenResponse;
}

export namespace DecodeTokenResponse {
  export type AsObject = {
    token?: Token.AsObject,
    transport?: TokenTransport.AsObject,
  }
}

export class LoginTokenRequest extends jspb.Message {
  getDuration(): string;
  setDuration(value: string): LoginTokenRequest;

  getUser(): Ref.User | undefined;
  setUser(value?: Ref.User): LoginTokenRequest;
  hasUser(): boolean;
  clearUser(): LoginTokenRequest;

  getTrigger(): boolean;
  setTrigger(value: boolean): LoginTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LoginTokenRequest): LoginTokenRequest.AsObject;
  static serializeBinaryToWriter(message: LoginTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginTokenRequest;
  static deserializeBinaryFromReader(message: LoginTokenRequest, reader: jspb.BinaryReader): LoginTokenRequest;
}

export namespace LoginTokenRequest {
  export type AsObject = {
    duration: string,
    user?: Ref.User.AsObject,
    trigger: boolean,
  }
}

export class GenerateRunnerTokenRequest extends jspb.Message {
  getDuration(): string;
  setDuration(value: string): GenerateRunnerTokenRequest;

  getId(): string;
  setId(value: string): GenerateRunnerTokenRequest;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): GenerateRunnerTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateRunnerTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateRunnerTokenRequest): GenerateRunnerTokenRequest.AsObject;
  static serializeBinaryToWriter(message: GenerateRunnerTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateRunnerTokenRequest;
  static deserializeBinaryFromReader(message: GenerateRunnerTokenRequest, reader: jspb.BinaryReader): GenerateRunnerTokenRequest;
}

export namespace GenerateRunnerTokenRequest {
  export type AsObject = {
    duration: string,
    id: string,
    labelsMap: Array<[string, string]>,
  }
}

export class InviteTokenRequest extends jspb.Message {
  getDuration(): string;
  setDuration(value: string): InviteTokenRequest;

  getLogin(): Token.Login | undefined;
  setLogin(value?: Token.Login): InviteTokenRequest;
  hasLogin(): boolean;
  clearLogin(): InviteTokenRequest;

  getSignup(): Token.Invite.Signup | undefined;
  setSignup(value?: Token.Invite.Signup): InviteTokenRequest;
  hasSignup(): boolean;
  clearSignup(): InviteTokenRequest;

  getUnusedEntrypoint(): Token.Entrypoint | undefined;
  setUnusedEntrypoint(value?: Token.Entrypoint): InviteTokenRequest;
  hasUnusedEntrypoint(): boolean;
  clearUnusedEntrypoint(): InviteTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InviteTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: InviteTokenRequest): InviteTokenRequest.AsObject;
  static serializeBinaryToWriter(message: InviteTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InviteTokenRequest;
  static deserializeBinaryFromReader(message: InviteTokenRequest, reader: jspb.BinaryReader): InviteTokenRequest;
}

export namespace InviteTokenRequest {
  export type AsObject = {
    duration: string,
    login?: Token.Login.AsObject,
    signup?: Token.Invite.Signup.AsObject,
    unusedEntrypoint?: Token.Entrypoint.AsObject,
  }
}

export class NewTokenResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): NewTokenResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: NewTokenResponse): NewTokenResponse.AsObject;
  static serializeBinaryToWriter(message: NewTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewTokenResponse;
  static deserializeBinaryFromReader(message: NewTokenResponse, reader: jspb.BinaryReader): NewTokenResponse;
}

export namespace NewTokenResponse {
  export type AsObject = {
    token: string,
  }
}

export class ConvertInviteTokenRequest extends jspb.Message {
  getToken(): string;
  setToken(value: string): ConvertInviteTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ConvertInviteTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ConvertInviteTokenRequest): ConvertInviteTokenRequest.AsObject;
  static serializeBinaryToWriter(message: ConvertInviteTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ConvertInviteTokenRequest;
  static deserializeBinaryFromReader(message: ConvertInviteTokenRequest, reader: jspb.BinaryReader): ConvertInviteTokenRequest;
}

export namespace ConvertInviteTokenRequest {
  export type AsObject = {
    token: string,
  }
}

export class CreateSnapshotResponse extends jspb.Message {
  getOpen(): CreateSnapshotResponse.Open | undefined;
  setOpen(value?: CreateSnapshotResponse.Open): CreateSnapshotResponse;
  hasOpen(): boolean;
  clearOpen(): CreateSnapshotResponse;

  getChunk(): Uint8Array | string;
  getChunk_asU8(): Uint8Array;
  getChunk_asB64(): string;
  setChunk(value: Uint8Array | string): CreateSnapshotResponse;

  getEventCase(): CreateSnapshotResponse.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSnapshotResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSnapshotResponse): CreateSnapshotResponse.AsObject;
  static serializeBinaryToWriter(message: CreateSnapshotResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSnapshotResponse;
  static deserializeBinaryFromReader(message: CreateSnapshotResponse, reader: jspb.BinaryReader): CreateSnapshotResponse;
}

export namespace CreateSnapshotResponse {
  export type AsObject = {
    open?: CreateSnapshotResponse.Open.AsObject,
    chunk: Uint8Array | string,
  }

  export class Open extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Open.AsObject;
    static toObject(includeInstance: boolean, msg: Open): Open.AsObject;
    static serializeBinaryToWriter(message: Open, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Open;
    static deserializeBinaryFromReader(message: Open, reader: jspb.BinaryReader): Open;
  }

  export namespace Open {
    export type AsObject = {
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    OPEN = 1,
    CHUNK = 2,
  }
}

export class RestoreSnapshotRequest extends jspb.Message {
  getOpen(): RestoreSnapshotRequest.Open | undefined;
  setOpen(value?: RestoreSnapshotRequest.Open): RestoreSnapshotRequest;
  hasOpen(): boolean;
  clearOpen(): RestoreSnapshotRequest;

  getChunk(): Uint8Array | string;
  getChunk_asU8(): Uint8Array;
  getChunk_asB64(): string;
  setChunk(value: Uint8Array | string): RestoreSnapshotRequest;

  getEventCase(): RestoreSnapshotRequest.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RestoreSnapshotRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RestoreSnapshotRequest): RestoreSnapshotRequest.AsObject;
  static serializeBinaryToWriter(message: RestoreSnapshotRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RestoreSnapshotRequest;
  static deserializeBinaryFromReader(message: RestoreSnapshotRequest, reader: jspb.BinaryReader): RestoreSnapshotRequest;
}

export namespace RestoreSnapshotRequest {
  export type AsObject = {
    open?: RestoreSnapshotRequest.Open.AsObject,
    chunk: Uint8Array | string,
  }

  export class Open extends jspb.Message {
    getExit(): boolean;
    setExit(value: boolean): Open;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Open.AsObject;
    static toObject(includeInstance: boolean, msg: Open): Open.AsObject;
    static serializeBinaryToWriter(message: Open, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Open;
    static deserializeBinaryFromReader(message: Open, reader: jspb.BinaryReader): Open;
  }

  export namespace Open {
    export type AsObject = {
      exit: boolean,
    }
  }


  export enum EventCase { 
    EVENT_NOT_SET = 0,
    OPEN = 1,
    CHUNK = 2,
  }
}

export class Snapshot extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Snapshot.AsObject;
  static toObject(includeInstance: boolean, msg: Snapshot): Snapshot.AsObject;
  static serializeBinaryToWriter(message: Snapshot, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Snapshot;
  static deserializeBinaryFromReader(message: Snapshot, reader: jspb.BinaryReader): Snapshot;
}

export namespace Snapshot {
  export type AsObject = {
  }

  export class Header extends jspb.Message {
    getVersion(): VersionInfo | undefined;
    setVersion(value?: VersionInfo): Header;
    hasVersion(): boolean;
    clearVersion(): Header;

    getFormat(): Snapshot.Header.Format;
    setFormat(value: Snapshot.Header.Format): Header;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Header.AsObject;
    static toObject(includeInstance: boolean, msg: Header): Header.AsObject;
    static serializeBinaryToWriter(message: Header, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Header;
    static deserializeBinaryFromReader(message: Header, reader: jspb.BinaryReader): Header;
  }

  export namespace Header {
    export type AsObject = {
      version?: VersionInfo.AsObject,
      format: Snapshot.Header.Format,
    }

    export enum Format { 
      UNKNOWN = 0,
      BOLT = 1,
    }
  }


  export class Trailer extends jspb.Message {
    getSha256(): string;
    setSha256(value: string): Trailer;

    getChecksumCase(): Trailer.ChecksumCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Trailer.AsObject;
    static toObject(includeInstance: boolean, msg: Trailer): Trailer.AsObject;
    static serializeBinaryToWriter(message: Trailer, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Trailer;
    static deserializeBinaryFromReader(message: Trailer, reader: jspb.BinaryReader): Trailer;
  }

  export namespace Trailer {
    export type AsObject = {
      sha256: string,
    }

    export enum ChecksumCase { 
      CHECKSUM_NOT_SET = 0,
      SHA256 = 1,
    }
  }


  export class BoltChunk extends jspb.Message {
    getBucket(): string;
    setBucket(value: string): BoltChunk;

    getItemsMap(): jspb.Map<string, Uint8Array | string>;
    clearItemsMap(): BoltChunk;

    getFinal(): boolean;
    setFinal(value: boolean): BoltChunk;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BoltChunk.AsObject;
    static toObject(includeInstance: boolean, msg: BoltChunk): BoltChunk.AsObject;
    static serializeBinaryToWriter(message: BoltChunk, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BoltChunk;
    static deserializeBinaryFromReader(message: BoltChunk, reader: jspb.BinaryReader): BoltChunk;
  }

  export namespace BoltChunk {
    export type AsObject = {
      bucket: string,
      itemsMap: Array<[string, Uint8Array | string]>,
      pb_final: boolean,
    }
  }

}

export class Hcl extends jspb.Message {
  getContents(): Uint8Array | string;
  getContents_asU8(): Uint8Array;
  getContents_asB64(): string;
  setContents(value: Uint8Array | string): Hcl;

  getFormat(): Hcl.Format;
  setFormat(value: Hcl.Format): Hcl;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Hcl.AsObject;
  static toObject(includeInstance: boolean, msg: Hcl): Hcl.AsObject;
  static serializeBinaryToWriter(message: Hcl, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Hcl;
  static deserializeBinaryFromReader(message: Hcl, reader: jspb.BinaryReader): Hcl;
}

export namespace Hcl {
  export type AsObject = {
    contents: Uint8Array | string,
    format: Hcl.Format,
  }

  export enum Format { 
    HCL = 0,
    JSON = 1,
  }
}

export class WaypointHclFmtRequest extends jspb.Message {
  getWaypointHcl(): Uint8Array | string;
  getWaypointHcl_asU8(): Uint8Array;
  getWaypointHcl_asB64(): string;
  setWaypointHcl(value: Uint8Array | string): WaypointHclFmtRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WaypointHclFmtRequest.AsObject;
  static toObject(includeInstance: boolean, msg: WaypointHclFmtRequest): WaypointHclFmtRequest.AsObject;
  static serializeBinaryToWriter(message: WaypointHclFmtRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WaypointHclFmtRequest;
  static deserializeBinaryFromReader(message: WaypointHclFmtRequest, reader: jspb.BinaryReader): WaypointHclFmtRequest;
}

export namespace WaypointHclFmtRequest {
  export type AsObject = {
    waypointHcl: Uint8Array | string,
  }
}

export class WaypointHclFmtResponse extends jspb.Message {
  getWaypointHcl(): Uint8Array | string;
  getWaypointHcl_asU8(): Uint8Array;
  getWaypointHcl_asB64(): string;
  setWaypointHcl(value: Uint8Array | string): WaypointHclFmtResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WaypointHclFmtResponse.AsObject;
  static toObject(includeInstance: boolean, msg: WaypointHclFmtResponse): WaypointHclFmtResponse.AsObject;
  static serializeBinaryToWriter(message: WaypointHclFmtResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WaypointHclFmtResponse;
  static deserializeBinaryFromReader(message: WaypointHclFmtResponse, reader: jspb.BinaryReader): WaypointHclFmtResponse;
}

export namespace WaypointHclFmtResponse {
  export type AsObject = {
    waypointHcl: Uint8Array | string,
  }
}

export class InlineKeepalive extends jspb.Message {
  getSignature(): string;
  setSignature(value: string): InlineKeepalive;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InlineKeepalive.AsObject;
  static toObject(includeInstance: boolean, msg: InlineKeepalive): InlineKeepalive.AsObject;
  static serializeBinaryToWriter(message: InlineKeepalive, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InlineKeepalive;
  static deserializeBinaryFromReader(message: InlineKeepalive, reader: jspb.BinaryReader): InlineKeepalive;
}

export namespace InlineKeepalive {
  export type AsObject = {
    signature: string,
  }
}

export enum ResourceCategoryDisplayHint { 
  UNKNOWN = 0,
  OTHER = 1,
  INSTANCE = 2,
  INSTANCE_MANAGER = 3,
  ROUTER = 4,
  POLICY = 5,
  CONFIG = 6,
  FUNCTION = 7,
  STORAGE = 8,
}
