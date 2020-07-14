/**
 * @fileoverview gRPC-Web generated client stub for hashicorp.waypoint
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as internal_server_proto_server_pb from 'waypoint-pb';


export class WaypointClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoListBuilds = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.ListBuildsResponse,
    (request: internal_server_proto_server_pb.ListBuildsRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.ListBuildsResponse.deserializeBinary
  );

  listBuilds(
    request: internal_server_proto_server_pb.ListBuildsRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.ListBuildsResponse>;

  listBuilds(
    request: internal_server_proto_server_pb.ListBuildsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ListBuildsResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.ListBuildsResponse>;

  listBuilds(
    request: internal_server_proto_server_pb.ListBuildsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ListBuildsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListBuilds',
        request,
        metadata || {},
        this.methodInfoListBuilds,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListBuilds',
    request,
    metadata || {},
    this.methodInfoListBuilds);
  }

  methodInfoListPushedArtifacts = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.ListPushedArtifactsResponse,
    (request: internal_server_proto_server_pb.ListPushedArtifactsRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.ListPushedArtifactsResponse.deserializeBinary
  );

  listPushedArtifacts(
    request: internal_server_proto_server_pb.ListPushedArtifactsRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.ListPushedArtifactsResponse>;

  listPushedArtifacts(
    request: internal_server_proto_server_pb.ListPushedArtifactsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ListPushedArtifactsResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.ListPushedArtifactsResponse>;

  listPushedArtifacts(
    request: internal_server_proto_server_pb.ListPushedArtifactsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ListPushedArtifactsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListPushedArtifacts',
        request,
        metadata || {},
        this.methodInfoListPushedArtifacts,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListPushedArtifacts',
    request,
    metadata || {},
    this.methodInfoListPushedArtifacts);
  }

  methodInfoListDeployments = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.ListDeploymentsResponse,
    (request: internal_server_proto_server_pb.ListDeploymentsRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.ListDeploymentsResponse.deserializeBinary
  );

  listDeployments(
    request: internal_server_proto_server_pb.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.ListDeploymentsResponse>;

  listDeployments(
    request: internal_server_proto_server_pb.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ListDeploymentsResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.ListDeploymentsResponse>;

  listDeployments(
    request: internal_server_proto_server_pb.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ListDeploymentsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListDeployments',
        request,
        metadata || {},
        this.methodInfoListDeployments,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListDeployments',
    request,
    metadata || {},
    this.methodInfoListDeployments);
  }

  methodInfoGetDeployment = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.Deployment,
    (request: internal_server_proto_server_pb.GetDeploymentRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.Deployment.deserializeBinary
  );

  getDeployment(
    request: internal_server_proto_server_pb.GetDeploymentRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.Deployment>;

  getDeployment(
    request: internal_server_proto_server_pb.GetDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Deployment) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.Deployment>;

  getDeployment(
    request: internal_server_proto_server_pb.GetDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Deployment) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetDeployment',
        request,
        metadata || {},
        this.methodInfoGetDeployment,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetDeployment',
    request,
    metadata || {},
    this.methodInfoGetDeployment);
  }

  methodInfoGetLatestBuild = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.Build,
    (request: internal_server_proto_server_pb.GetLatestBuildRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.Build.deserializeBinary
  );

  getLatestBuild(
    request: internal_server_proto_server_pb.GetLatestBuildRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.Build>;

  getLatestBuild(
    request: internal_server_proto_server_pb.GetLatestBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Build) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.Build>;

  getLatestBuild(
    request: internal_server_proto_server_pb.GetLatestBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Build) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetLatestBuild',
        request,
        metadata || {},
        this.methodInfoGetLatestBuild,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetLatestBuild',
    request,
    metadata || {},
    this.methodInfoGetLatestBuild);
  }

  methodInfoGetLatestPushedArtifact = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.PushedArtifact,
    (request: internal_server_proto_server_pb.GetLatestPushedArtifactRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.PushedArtifact.deserializeBinary
  );

  getLatestPushedArtifact(
    request: internal_server_proto_server_pb.GetLatestPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.PushedArtifact>;

  getLatestPushedArtifact(
    request: internal_server_proto_server_pb.GetLatestPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.PushedArtifact) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.PushedArtifact>;

  getLatestPushedArtifact(
    request: internal_server_proto_server_pb.GetLatestPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.PushedArtifact) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetLatestPushedArtifact',
        request,
        metadata || {},
        this.methodInfoGetLatestPushedArtifact,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetLatestPushedArtifact',
    request,
    metadata || {},
    this.methodInfoGetLatestPushedArtifact);
  }

  methodInfoGetLogStream = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.LogBatch,
    (request: internal_server_proto_server_pb.GetLogStreamRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.LogBatch.deserializeBinary
  );

  getLogStream(
    request: internal_server_proto_server_pb.GetLogStreamRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/hashicorp.waypoint.Waypoint/GetLogStream',
      request,
      metadata || {},
      this.methodInfoGetLogStream);
  }

  methodInfoSetConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.ConfigSetResponse,
    (request: internal_server_proto_server_pb.ConfigSetRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.ConfigSetResponse.deserializeBinary
  );

  setConfig(
    request: internal_server_proto_server_pb.ConfigSetRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.ConfigSetResponse>;

  setConfig(
    request: internal_server_proto_server_pb.ConfigSetRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ConfigSetResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.ConfigSetResponse>;

  setConfig(
    request: internal_server_proto_server_pb.ConfigSetRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ConfigSetResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/SetConfig',
        request,
        metadata || {},
        this.methodInfoSetConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/SetConfig',
    request,
    metadata || {},
    this.methodInfoSetConfig);
  }

  methodInfoGetConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.ConfigGetResponse,
    (request: internal_server_proto_server_pb.ConfigGetRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.ConfigGetResponse.deserializeBinary
  );

  getConfig(
    request: internal_server_proto_server_pb.ConfigGetRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.ConfigGetResponse>;

  getConfig(
    request: internal_server_proto_server_pb.ConfigGetRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ConfigGetResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.ConfigGetResponse>;

  getConfig(
    request: internal_server_proto_server_pb.ConfigGetRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ConfigGetResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetConfig',
        request,
        metadata || {},
        this.methodInfoGetConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetConfig',
    request,
    metadata || {},
    this.methodInfoGetConfig);
  }

  methodInfoQueueJob = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.QueueJobResponse,
    (request: internal_server_proto_server_pb.QueueJobRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.QueueJobResponse.deserializeBinary
  );

  queueJob(
    request: internal_server_proto_server_pb.QueueJobRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.QueueJobResponse>;

  queueJob(
    request: internal_server_proto_server_pb.QueueJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.QueueJobResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.QueueJobResponse>;

  queueJob(
    request: internal_server_proto_server_pb.QueueJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.QueueJobResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/QueueJob',
        request,
        metadata || {},
        this.methodInfoQueueJob,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/QueueJob',
    request,
    metadata || {},
    this.methodInfoQueueJob);
  }

  methodInfoGetJob = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.Job,
    (request: internal_server_proto_server_pb.GetJobRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.Job.deserializeBinary
  );

  getJob(
    request: internal_server_proto_server_pb.GetJobRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.Job>;

  getJob(
    request: internal_server_proto_server_pb.GetJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Job) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.Job>;

  getJob(
    request: internal_server_proto_server_pb.GetJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Job) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetJob',
        request,
        metadata || {},
        this.methodInfoGetJob,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetJob',
    request,
    metadata || {},
    this.methodInfoGetJob);
  }

  methodInfoValidateJob = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.ValidateJobResponse,
    (request: internal_server_proto_server_pb.ValidateJobRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.ValidateJobResponse.deserializeBinary
  );

  validateJob(
    request: internal_server_proto_server_pb.ValidateJobRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.ValidateJobResponse>;

  validateJob(
    request: internal_server_proto_server_pb.ValidateJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ValidateJobResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.ValidateJobResponse>;

  validateJob(
    request: internal_server_proto_server_pb.ValidateJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.ValidateJobResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ValidateJob',
        request,
        metadata || {},
        this.methodInfoValidateJob,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ValidateJob',
    request,
    metadata || {},
    this.methodInfoValidateJob);
  }

  methodInfoGetJobStream = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.GetJobStreamResponse,
    (request: internal_server_proto_server_pb.GetJobStreamRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.GetJobStreamResponse.deserializeBinary
  );

  getJobStream(
    request: internal_server_proto_server_pb.GetJobStreamRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/hashicorp.waypoint.Waypoint/GetJobStream',
      request,
      metadata || {},
      this.methodInfoGetJobStream);
  }

  methodInfoGetRunner = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.Runner,
    (request: internal_server_proto_server_pb.GetRunnerRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.Runner.deserializeBinary
  );

  getRunner(
    request: internal_server_proto_server_pb.GetRunnerRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.Runner>;

  getRunner(
    request: internal_server_proto_server_pb.GetRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Runner) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.Runner>;

  getRunner(
    request: internal_server_proto_server_pb.GetRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.Runner) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetRunner',
        request,
        metadata || {},
        this.methodInfoGetRunner,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetRunner',
    request,
    metadata || {},
    this.methodInfoGetRunner);
  }

  methodInfoSetServerConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    google_protobuf_empty_pb.Empty,
    (request: internal_server_proto_server_pb.SetServerConfigRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  setServerConfig(
    request: internal_server_proto_server_pb.SetServerConfigRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  setServerConfig(
    request: internal_server_proto_server_pb.SetServerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  setServerConfig(
    request: internal_server_proto_server_pb.SetServerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/SetServerConfig',
        request,
        metadata || {},
        this.methodInfoSetServerConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/SetServerConfig',
    request,
    metadata || {},
    this.methodInfoSetServerConfig);
  }

  methodInfoRunnerGetDeploymentConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.RunnerGetDeploymentConfigResponse,
    (request: internal_server_proto_server_pb.RunnerGetDeploymentConfigRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.RunnerGetDeploymentConfigResponse.deserializeBinary
  );

  runnerGetDeploymentConfig(
    request: internal_server_proto_server_pb.RunnerGetDeploymentConfigRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.RunnerGetDeploymentConfigResponse>;

  runnerGetDeploymentConfig(
    request: internal_server_proto_server_pb.RunnerGetDeploymentConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.RunnerGetDeploymentConfigResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.RunnerGetDeploymentConfigResponse>;

  runnerGetDeploymentConfig(
    request: internal_server_proto_server_pb.RunnerGetDeploymentConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.RunnerGetDeploymentConfigResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/RunnerGetDeploymentConfig',
        request,
        metadata || {},
        this.methodInfoRunnerGetDeploymentConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/RunnerGetDeploymentConfig',
    request,
    metadata || {},
    this.methodInfoRunnerGetDeploymentConfig);
  }

  methodInfoEntrypointConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.EntrypointConfigResponse,
    (request: internal_server_proto_server_pb.EntrypointConfigRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.EntrypointConfigResponse.deserializeBinary
  );

  entrypointConfig(
    request: internal_server_proto_server_pb.EntrypointConfigRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/hashicorp.waypoint.Waypoint/EntrypointConfig',
      request,
      metadata || {},
      this.methodInfoEntrypointConfig);
  }

  methodInfoUpsertBuild = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.UpsertBuildResponse,
    (request: internal_server_proto_server_pb.UpsertBuildRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.UpsertBuildResponse.deserializeBinary
  );

  upsertBuild(
    request: internal_server_proto_server_pb.UpsertBuildRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.UpsertBuildResponse>;

  upsertBuild(
    request: internal_server_proto_server_pb.UpsertBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertBuildResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.UpsertBuildResponse>;

  upsertBuild(
    request: internal_server_proto_server_pb.UpsertBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertBuildResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertBuild',
        request,
        metadata || {},
        this.methodInfoUpsertBuild,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertBuild',
    request,
    metadata || {},
    this.methodInfoUpsertBuild);
  }

  methodInfoUpsertPushedArtifact = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.UpsertPushedArtifactResponse,
    (request: internal_server_proto_server_pb.UpsertPushedArtifactRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.UpsertPushedArtifactResponse.deserializeBinary
  );

  upsertPushedArtifact(
    request: internal_server_proto_server_pb.UpsertPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.UpsertPushedArtifactResponse>;

  upsertPushedArtifact(
    request: internal_server_proto_server_pb.UpsertPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertPushedArtifactResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.UpsertPushedArtifactResponse>;

  upsertPushedArtifact(
    request: internal_server_proto_server_pb.UpsertPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertPushedArtifactResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertPushedArtifact',
        request,
        metadata || {},
        this.methodInfoUpsertPushedArtifact,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertPushedArtifact',
    request,
    metadata || {},
    this.methodInfoUpsertPushedArtifact);
  }

  methodInfoUpsertDeployment = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.UpsertDeploymentResponse,
    (request: internal_server_proto_server_pb.UpsertDeploymentRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.UpsertDeploymentResponse.deserializeBinary
  );

  upsertDeployment(
    request: internal_server_proto_server_pb.UpsertDeploymentRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.UpsertDeploymentResponse>;

  upsertDeployment(
    request: internal_server_proto_server_pb.UpsertDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertDeploymentResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.UpsertDeploymentResponse>;

  upsertDeployment(
    request: internal_server_proto_server_pb.UpsertDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertDeploymentResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertDeployment',
        request,
        metadata || {},
        this.methodInfoUpsertDeployment,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertDeployment',
    request,
    metadata || {},
    this.methodInfoUpsertDeployment);
  }

  methodInfoUpsertRelease = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.UpsertReleaseResponse,
    (request: internal_server_proto_server_pb.UpsertReleaseRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.UpsertReleaseResponse.deserializeBinary
  );

  upsertRelease(
    request: internal_server_proto_server_pb.UpsertReleaseRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.UpsertReleaseResponse>;

  upsertRelease(
    request: internal_server_proto_server_pb.UpsertReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertReleaseResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.UpsertReleaseResponse>;

  upsertRelease(
    request: internal_server_proto_server_pb.UpsertReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.UpsertReleaseResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertRelease',
        request,
        metadata || {},
        this.methodInfoUpsertRelease,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertRelease',
    request,
    metadata || {},
    this.methodInfoUpsertRelease);
  }

  methodInfoGenerateInviteToken = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.NewTokenResponse,
    (request: internal_server_proto_server_pb.InviteTokenRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  generateInviteToken(
    request: internal_server_proto_server_pb.InviteTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.NewTokenResponse>;

  generateInviteToken(
    request: internal_server_proto_server_pb.InviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.NewTokenResponse>;

  generateInviteToken(
    request: internal_server_proto_server_pb.InviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GenerateInviteToken',
        request,
        metadata || {},
        this.methodInfoGenerateInviteToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GenerateInviteToken',
    request,
    metadata || {},
    this.methodInfoGenerateInviteToken);
  }

  methodInfoGenerateLoginToken = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.NewTokenResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  generateLoginToken(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.NewTokenResponse>;

  generateLoginToken(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.NewTokenResponse>;

  generateLoginToken(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GenerateLoginToken',
        request,
        metadata || {},
        this.methodInfoGenerateLoginToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GenerateLoginToken',
    request,
    metadata || {},
    this.methodInfoGenerateLoginToken);
  }

  methodInfoConvertInviteToken = new grpcWeb.AbstractClientBase.MethodInfo(
    internal_server_proto_server_pb.NewTokenResponse,
    (request: internal_server_proto_server_pb.ConvertInviteTokenRequest) => {
      return request.serializeBinary();
    },
    internal_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  convertInviteToken(
    request: internal_server_proto_server_pb.ConvertInviteTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<internal_server_proto_server_pb.NewTokenResponse>;

  convertInviteToken(
    request: internal_server_proto_server_pb.ConvertInviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<internal_server_proto_server_pb.NewTokenResponse>;

  convertInviteToken(
    request: internal_server_proto_server_pb.ConvertInviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: internal_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ConvertInviteToken',
        request,
        metadata || {},
        this.methodInfoConvertInviteToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ConvertInviteToken',
    request,
    metadata || {},
    this.methodInfoConvertInviteToken);
  }

}

