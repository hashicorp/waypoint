// source: internal/server/proto/server.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');
goog.object.extend(proto, google_protobuf_any_pb);
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');
goog.object.extend(proto, google_protobuf_empty_pb);
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');
goog.object.extend(proto, google_protobuf_timestamp_pb);
var google_rpc_status_pb = require('api-common-protos/google/rpc/status_pb.js');
goog.object.extend(proto, google_rpc_status_pb);
goog.exportSymbol('proto.hashicorp.waypoint.Application', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Artifact', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Build', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Component', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Component.Type', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConfigGetRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConfigGetRequest.ScopeCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConfigGetResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConfigSetRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConfigSetResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConfigVar', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConfigVar.ScopeCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ConvertInviteTokenRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Deployment', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Deployment.State', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointConfig', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointConfig.Exec', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointConfigRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointConfigResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecRequest.Error', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecRequest.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecRequest.Exit', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecRequest.Open', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecRequest.Output', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecRequest.Output.Channel', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointExecResponse.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.EntrypointLogBatch', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamRequest.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamRequest.Input', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamRequest.PTY', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamRequest.Start', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamRequest.WindowSize', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamResponse.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamResponse.Exit', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamResponse.Output', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ExecStreamResponse.Output.Channel', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetDeploymentRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Complete', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Error', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Open', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.State', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetLatestBuildRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetLatestPushedArtifactRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetLogStreamRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.GetRunnerRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.HMACKey', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.InviteTokenRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.BuildOp', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.BuildResult', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.DataSourceCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.DeployOp', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.DeployResult', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.DestroyDeployOp', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.Local', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.Noop', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.OperationCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.PushOp', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.PushResult', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.ReleaseOp', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.ReleaseResult', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.Result', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Job.State', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ListBuildsRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ListBuildsResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ListDeploymentsRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ListDeploymentsResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ListPushedArtifactsRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ListPushedArtifactsResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.LogBatch', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.LogBatch.Entry', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.NewTokenResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.OperationOrder', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.OperationOrder.Order', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Project', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.PushedArtifact', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.QueueJobRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.QueueJobResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref.Application', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref.Project', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref.Runner', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref.Runner.TargetCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref.RunnerAny', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref.RunnerId', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Ref.Workspace', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Release', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Release.Split', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Release.SplitTarget', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Runner', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerConfig', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerConfigRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerConfigRequest.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerConfigRequest.Open', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerConfigResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamRequest.Error', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamRequest.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamRequest.Request', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamResponse.EventCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ServerConfig', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.SetServerConfigRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Status', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Status.State', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.StatusFilter', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.StatusFilter.Filter', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.StatusFilter.Filter.FilterCase', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Token', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.TokenTransport', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertBuildRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertBuildResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertDeploymentRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertDeploymentResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertPushedArtifactRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertPushedArtifactResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertReleaseRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.UpsertReleaseResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ValidateJobRequest', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.ValidateJobResponse', null, global);
goog.exportSymbol('proto.hashicorp.waypoint.Workspace', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Application = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Application, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Application.displayName = 'proto.hashicorp.waypoint.Application';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Project = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Project, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Project.displayName = 'proto.hashicorp.waypoint.Project';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Workspace = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Workspace, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Workspace.displayName = 'proto.hashicorp.waypoint.Workspace';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Ref = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Ref, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Ref.displayName = 'proto.hashicorp.waypoint.Ref';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Ref.Application = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Ref.Application, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Ref.Application.displayName = 'proto.hashicorp.waypoint.Ref.Application';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Ref.Project = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Ref.Project, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Ref.Project.displayName = 'proto.hashicorp.waypoint.Ref.Project';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Ref.Workspace = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Ref.Workspace, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Ref.Workspace.displayName = 'proto.hashicorp.waypoint.Ref.Workspace';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Ref.Runner = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.Ref.Runner.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.Ref.Runner, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Ref.Runner.displayName = 'proto.hashicorp.waypoint.Ref.Runner';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Ref.RunnerId = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Ref.RunnerId, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Ref.RunnerId.displayName = 'proto.hashicorp.waypoint.Ref.RunnerId';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Ref.RunnerAny = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Ref.RunnerAny, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Ref.RunnerAny.displayName = 'proto.hashicorp.waypoint.Ref.RunnerAny';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Component = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Component, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Component.displayName = 'proto.hashicorp.waypoint.Component';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Status = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Status, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Status.displayName = 'proto.hashicorp.waypoint.Status';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.StatusFilter = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.StatusFilter.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.StatusFilter, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.StatusFilter.displayName = 'proto.hashicorp.waypoint.StatusFilter';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.StatusFilter.Filter = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.StatusFilter.Filter.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.StatusFilter.Filter, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.StatusFilter.Filter.displayName = 'proto.hashicorp.waypoint.StatusFilter.Filter';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.OperationOrder = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.OperationOrder, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.OperationOrder.displayName = 'proto.hashicorp.waypoint.OperationOrder';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.QueueJobRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.QueueJobRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.QueueJobRequest.displayName = 'proto.hashicorp.waypoint.QueueJobRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.QueueJobResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.QueueJobResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.QueueJobResponse.displayName = 'proto.hashicorp.waypoint.QueueJobResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ValidateJobRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ValidateJobRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ValidateJobRequest.displayName = 'proto.hashicorp.waypoint.ValidateJobRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ValidateJobResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ValidateJobResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ValidateJobResponse.displayName = 'proto.hashicorp.waypoint.ValidateJobResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.Job.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.Job, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.displayName = 'proto.hashicorp.waypoint.Job';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.Result = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.Result, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.Result.displayName = 'proto.hashicorp.waypoint.Job.Result';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.Local = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.Local, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.Local.displayName = 'proto.hashicorp.waypoint.Job.Local';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.Noop = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.Noop, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.Noop.displayName = 'proto.hashicorp.waypoint.Job.Noop';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.BuildOp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.BuildOp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.BuildOp.displayName = 'proto.hashicorp.waypoint.Job.BuildOp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.BuildResult = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.BuildResult, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.BuildResult.displayName = 'proto.hashicorp.waypoint.Job.BuildResult';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.PushOp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.PushOp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.PushOp.displayName = 'proto.hashicorp.waypoint.Job.PushOp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.PushResult = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.PushResult, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.PushResult.displayName = 'proto.hashicorp.waypoint.Job.PushResult';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.DeployOp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.DeployOp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.DeployOp.displayName = 'proto.hashicorp.waypoint.Job.DeployOp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.DeployResult = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.DeployResult, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.DeployResult.displayName = 'proto.hashicorp.waypoint.Job.DeployResult';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.DestroyDeployOp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.DestroyDeployOp.displayName = 'proto.hashicorp.waypoint.Job.DestroyDeployOp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.ReleaseOp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.ReleaseOp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.ReleaseOp.displayName = 'proto.hashicorp.waypoint.Job.ReleaseOp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Job.ReleaseResult = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Job.ReleaseResult, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Job.ReleaseResult.displayName = 'proto.hashicorp.waypoint.Job.ReleaseResult';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobRequest.displayName = 'proto.hashicorp.waypoint.GetJobRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamRequest.displayName = 'proto.hashicorp.waypoint.GetJobStreamRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Open = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Open, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Open.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Open';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.State, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.State.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.State';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Error, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Error.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Error';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetJobStreamResponse.Complete, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetJobStreamResponse.Complete.displayName = 'proto.hashicorp.waypoint.GetJobStreamResponse.Complete';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Runner = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.Runner.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.Runner, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Runner.displayName = 'proto.hashicorp.waypoint.Runner';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerConfigRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.RunnerConfigRequest.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.RunnerConfigRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerConfigRequest.displayName = 'proto.hashicorp.waypoint.RunnerConfigRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerConfigRequest.Open, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerConfigRequest.Open.displayName = 'proto.hashicorp.waypoint.RunnerConfigRequest.Open';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerConfigResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerConfigResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerConfigResponse.displayName = 'proto.hashicorp.waypoint.RunnerConfigResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.RunnerConfig.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerConfig.displayName = 'proto.hashicorp.waypoint.RunnerConfig';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.RunnerJobStreamRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerJobStreamRequest.displayName = 'proto.hashicorp.waypoint.RunnerJobStreamRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerJobStreamRequest.Request, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.displayName = 'proto.hashicorp.waypoint.RunnerJobStreamRequest.Request';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.displayName = 'proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.displayName = 'proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerJobStreamRequest.Error, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.displayName = 'proto.hashicorp.waypoint.RunnerJobStreamRequest.Error';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.RunnerJobStreamResponse.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.RunnerJobStreamResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerJobStreamResponse.displayName = 'proto.hashicorp.waypoint.RunnerJobStreamResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.displayName = 'proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.displayName = 'proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.displayName = 'proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetRunnerRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetRunnerRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetRunnerRequest.displayName = 'proto.hashicorp.waypoint.GetRunnerRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.SetServerConfigRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.SetServerConfigRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.SetServerConfigRequest.displayName = 'proto.hashicorp.waypoint.SetServerConfigRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ServerConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ServerConfig.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ServerConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ServerConfig.displayName = 'proto.hashicorp.waypoint.ServerConfig';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.displayName = 'proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertBuildRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertBuildRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertBuildRequest.displayName = 'proto.hashicorp.waypoint.UpsertBuildRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertBuildResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertBuildResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertBuildResponse.displayName = 'proto.hashicorp.waypoint.UpsertBuildResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ListBuildsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ListBuildsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ListBuildsRequest.displayName = 'proto.hashicorp.waypoint.ListBuildsRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ListBuildsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ListBuildsResponse.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ListBuildsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ListBuildsResponse.displayName = 'proto.hashicorp.waypoint.ListBuildsResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetLatestBuildRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetLatestBuildRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetLatestBuildRequest.displayName = 'proto.hashicorp.waypoint.GetLatestBuildRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Build = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Build, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Build.displayName = 'proto.hashicorp.waypoint.Build';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Artifact = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Artifact, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Artifact.displayName = 'proto.hashicorp.waypoint.Artifact';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertPushedArtifactRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertPushedArtifactRequest.displayName = 'proto.hashicorp.waypoint.UpsertPushedArtifactRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertPushedArtifactResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertPushedArtifactResponse.displayName = 'proto.hashicorp.waypoint.UpsertPushedArtifactResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetLatestPushedArtifactRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.displayName = 'proto.hashicorp.waypoint.GetLatestPushedArtifactRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ListPushedArtifactsRequest.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ListPushedArtifactsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ListPushedArtifactsRequest.displayName = 'proto.hashicorp.waypoint.ListPushedArtifactsRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ListPushedArtifactsResponse.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ListPushedArtifactsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ListPushedArtifactsResponse.displayName = 'proto.hashicorp.waypoint.ListPushedArtifactsResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.PushedArtifact = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.PushedArtifact, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.PushedArtifact.displayName = 'proto.hashicorp.waypoint.PushedArtifact';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetDeploymentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetDeploymentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetDeploymentRequest.displayName = 'proto.hashicorp.waypoint.GetDeploymentRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertDeploymentRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertDeploymentRequest.displayName = 'proto.hashicorp.waypoint.UpsertDeploymentRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertDeploymentResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertDeploymentResponse.displayName = 'proto.hashicorp.waypoint.UpsertDeploymentResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ListDeploymentsRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ListDeploymentsRequest.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ListDeploymentsRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ListDeploymentsRequest.displayName = 'proto.hashicorp.waypoint.ListDeploymentsRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ListDeploymentsResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ListDeploymentsResponse.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ListDeploymentsResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ListDeploymentsResponse.displayName = 'proto.hashicorp.waypoint.ListDeploymentsResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Deployment = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Deployment, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Deployment.displayName = 'proto.hashicorp.waypoint.Deployment';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertReleaseRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertReleaseRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertReleaseRequest.displayName = 'proto.hashicorp.waypoint.UpsertReleaseRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.UpsertReleaseResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.UpsertReleaseResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.UpsertReleaseResponse.displayName = 'proto.hashicorp.waypoint.UpsertReleaseResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Release = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Release, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Release.displayName = 'proto.hashicorp.waypoint.Release';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Release.Split = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.Release.Split.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.Release.Split, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Release.Split.displayName = 'proto.hashicorp.waypoint.Release.Split';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Release.SplitTarget = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Release.SplitTarget, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Release.SplitTarget.displayName = 'proto.hashicorp.waypoint.Release.SplitTarget';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.GetLogStreamRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.GetLogStreamRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.GetLogStreamRequest.displayName = 'proto.hashicorp.waypoint.GetLogStreamRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.LogBatch = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.LogBatch.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.LogBatch, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.LogBatch.displayName = 'proto.hashicorp.waypoint.LogBatch';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.LogBatch.Entry = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.LogBatch.Entry, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.LogBatch.Entry.displayName = 'proto.hashicorp.waypoint.LogBatch.Entry';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ConfigVar = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.ConfigVar.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.ConfigVar, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ConfigVar.displayName = 'proto.hashicorp.waypoint.ConfigVar';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ConfigSetRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ConfigSetRequest.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ConfigSetRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ConfigSetRequest.displayName = 'proto.hashicorp.waypoint.ConfigSetRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ConfigSetResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ConfigSetResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ConfigSetResponse.displayName = 'proto.hashicorp.waypoint.ConfigSetResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ConfigGetRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.ConfigGetRequest.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.ConfigGetRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ConfigGetRequest.displayName = 'proto.hashicorp.waypoint.ConfigGetRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ConfigGetResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ConfigGetResponse.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ConfigGetResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ConfigGetResponse.displayName = 'proto.hashicorp.waypoint.ConfigGetResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.ExecStreamRequest.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamRequest.displayName = 'proto.hashicorp.waypoint.ExecStreamRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.ExecStreamRequest.Start.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamRequest.Start, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamRequest.Start.displayName = 'proto.hashicorp.waypoint.ExecStreamRequest.Start';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamRequest.Input, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamRequest.Input.displayName = 'proto.hashicorp.waypoint.ExecStreamRequest.Input';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamRequest.PTY, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamRequest.PTY.displayName = 'proto.hashicorp.waypoint.ExecStreamRequest.PTY';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamRequest.WindowSize, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.displayName = 'proto.hashicorp.waypoint.ExecStreamRequest.WindowSize';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.ExecStreamResponse.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamResponse.displayName = 'proto.hashicorp.waypoint.ExecStreamResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamResponse.Exit, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamResponse.Exit.displayName = 'proto.hashicorp.waypoint.ExecStreamResponse.Exit';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ExecStreamResponse.Output, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ExecStreamResponse.Output.displayName = 'proto.hashicorp.waypoint.ExecStreamResponse.Output';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointConfigRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointConfigRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointConfigRequest.displayName = 'proto.hashicorp.waypoint.EntrypointConfigRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointConfigResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointConfigResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointConfigResponse.displayName = 'proto.hashicorp.waypoint.EntrypointConfigResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.EntrypointConfig.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointConfig.displayName = 'proto.hashicorp.waypoint.EntrypointConfig';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.EntrypointConfig.Exec.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointConfig.Exec, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointConfig.Exec.displayName = 'proto.hashicorp.waypoint.EntrypointConfig.Exec';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointLogBatch = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.hashicorp.waypoint.EntrypointLogBatch.repeatedFields_, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointLogBatch, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointLogBatch.displayName = 'proto.hashicorp.waypoint.EntrypointLogBatch';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointExecRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.EntrypointExecRequest.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointExecRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointExecRequest.displayName = 'proto.hashicorp.waypoint.EntrypointExecRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointExecRequest.Open, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointExecRequest.Open.displayName = 'proto.hashicorp.waypoint.EntrypointExecRequest.Open';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointExecRequest.Exit, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointExecRequest.Exit.displayName = 'proto.hashicorp.waypoint.EntrypointExecRequest.Exit';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointExecRequest.Output, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointExecRequest.Output.displayName = 'proto.hashicorp.waypoint.EntrypointExecRequest.Output';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointExecRequest.Error, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointExecRequest.Error.displayName = 'proto.hashicorp.waypoint.EntrypointExecRequest.Error';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.EntrypointExecResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_);
};
goog.inherits(proto.hashicorp.waypoint.EntrypointExecResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.EntrypointExecResponse.displayName = 'proto.hashicorp.waypoint.EntrypointExecResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.TokenTransport = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.TokenTransport, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.TokenTransport.displayName = 'proto.hashicorp.waypoint.TokenTransport';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.Token = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.Token, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.Token.displayName = 'proto.hashicorp.waypoint.Token';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.HMACKey = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.HMACKey, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.HMACKey.displayName = 'proto.hashicorp.waypoint.HMACKey';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.InviteTokenRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.InviteTokenRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.InviteTokenRequest.displayName = 'proto.hashicorp.waypoint.InviteTokenRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.NewTokenResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.NewTokenResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.NewTokenResponse.displayName = 'proto.hashicorp.waypoint.NewTokenResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.hashicorp.waypoint.ConvertInviteTokenRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.hashicorp.waypoint.ConvertInviteTokenRequest.displayName = 'proto.hashicorp.waypoint.ConvertInviteTokenRequest';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Application.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Application.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Application} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Application.toObject = function(includeInstance, msg) {
  var f, obj = {
    project: (f = msg.getProject()) && proto.hashicorp.waypoint.Ref.Project.toObject(includeInstance, f),
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Application}
 */
proto.hashicorp.waypoint.Application.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Application;
  return proto.hashicorp.waypoint.Application.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Application} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Application}
 */
proto.hashicorp.waypoint.Application.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = new proto.hashicorp.waypoint.Ref.Project;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Project.deserializeBinaryFromReader);
      msg.setProject(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Application.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Application.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Application} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Application.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProject();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Ref.Project.serializeBinaryToWriter
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional Ref.Project project = 2;
 * @return {?proto.hashicorp.waypoint.Ref.Project}
 */
proto.hashicorp.waypoint.Application.prototype.getProject = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Project} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Project, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Project|undefined} value
 * @return {!proto.hashicorp.waypoint.Application} returns this
*/
proto.hashicorp.waypoint.Application.prototype.setProject = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Application} returns this
 */
proto.hashicorp.waypoint.Application.prototype.clearProject = function() {
  return this.setProject(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Application.prototype.hasProject = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Application.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Application} returns this
 */
proto.hashicorp.waypoint.Application.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Project.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Project.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Project} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Project.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Project}
 */
proto.hashicorp.waypoint.Project.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Project;
  return proto.hashicorp.waypoint.Project.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Project} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Project}
 */
proto.hashicorp.waypoint.Project.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Project.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Project.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Project} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Project.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Project.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Project} returns this
 */
proto.hashicorp.waypoint.Project.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Workspace.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Workspace.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Workspace} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Workspace.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Workspace}
 */
proto.hashicorp.waypoint.Workspace.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Workspace;
  return proto.hashicorp.waypoint.Workspace.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Workspace} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Workspace}
 */
proto.hashicorp.waypoint.Workspace.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Workspace.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Workspace.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Workspace} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Workspace.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Workspace.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Workspace} returns this
 */
proto.hashicorp.waypoint.Workspace.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Ref.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Ref.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Ref} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Ref}
 */
proto.hashicorp.waypoint.Ref.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Ref;
  return proto.hashicorp.waypoint.Ref.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Ref} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Ref}
 */
proto.hashicorp.waypoint.Ref.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Ref.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Ref.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Ref} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Ref.Application.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Ref.Application.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Ref.Application} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Application.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: jspb.Message.getFieldWithDefault(msg, 1, ""),
    project: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.Ref.Application.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Ref.Application;
  return proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Ref.Application} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setApplication(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setProject(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Ref.Application.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Ref.Application} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getProject();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string application = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Ref.Application.prototype.getApplication = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Ref.Application} returns this
 */
proto.hashicorp.waypoint.Ref.Application.prototype.setApplication = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string project = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.Ref.Application.prototype.getProject = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Ref.Application} returns this
 */
proto.hashicorp.waypoint.Ref.Application.prototype.setProject = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Ref.Project.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Ref.Project.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Ref.Project} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Project.toObject = function(includeInstance, msg) {
  var f, obj = {
    project: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Ref.Project}
 */
proto.hashicorp.waypoint.Ref.Project.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Ref.Project;
  return proto.hashicorp.waypoint.Ref.Project.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Ref.Project} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Ref.Project}
 */
proto.hashicorp.waypoint.Ref.Project.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setProject(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Ref.Project.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Ref.Project.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Ref.Project} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Project.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProject();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string project = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Ref.Project.prototype.getProject = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Ref.Project} returns this
 */
proto.hashicorp.waypoint.Ref.Project.prototype.setProject = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Ref.Workspace.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Ref.Workspace.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Ref.Workspace} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Workspace.toObject = function(includeInstance, msg) {
  var f, obj = {
    workspace: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.Ref.Workspace.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Ref.Workspace;
  return proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Ref.Workspace} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setWorkspace(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Ref.Workspace.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Ref.Workspace} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getWorkspace();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string workspace = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Ref.Workspace.prototype.getWorkspace = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Ref.Workspace} returns this
 */
proto.hashicorp.waypoint.Ref.Workspace.prototype.setWorkspace = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.Ref.Runner.oneofGroups_ = [[1,2]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.Ref.Runner.TargetCase = {
  TARGET_NOT_SET: 0,
  ANY: 1,
  ID: 2
};

/**
 * @return {proto.hashicorp.waypoint.Ref.Runner.TargetCase}
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.getTargetCase = function() {
  return /** @type {proto.hashicorp.waypoint.Ref.Runner.TargetCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.Ref.Runner.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Ref.Runner.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Ref.Runner} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Runner.toObject = function(includeInstance, msg) {
  var f, obj = {
    any: (f = msg.getAny()) && proto.hashicorp.waypoint.Ref.RunnerAny.toObject(includeInstance, f),
    id: (f = msg.getId()) && proto.hashicorp.waypoint.Ref.RunnerId.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Ref.Runner}
 */
proto.hashicorp.waypoint.Ref.Runner.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Ref.Runner;
  return proto.hashicorp.waypoint.Ref.Runner.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Ref.Runner} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Ref.Runner}
 */
proto.hashicorp.waypoint.Ref.Runner.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Ref.RunnerAny;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.RunnerAny.deserializeBinaryFromReader);
      msg.setAny(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Ref.RunnerId;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.RunnerId.deserializeBinaryFromReader);
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Ref.Runner.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Ref.Runner} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.Runner.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAny();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Ref.RunnerAny.serializeBinaryToWriter
    );
  }
  f = message.getId();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Ref.RunnerId.serializeBinaryToWriter
    );
  }
};


/**
 * optional RunnerAny any = 1;
 * @return {?proto.hashicorp.waypoint.Ref.RunnerAny}
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.getAny = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.RunnerAny} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.RunnerAny, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.RunnerAny|undefined} value
 * @return {!proto.hashicorp.waypoint.Ref.Runner} returns this
*/
proto.hashicorp.waypoint.Ref.Runner.prototype.setAny = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.Ref.Runner.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Ref.Runner} returns this
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.clearAny = function() {
  return this.setAny(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.hasAny = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional RunnerId id = 2;
 * @return {?proto.hashicorp.waypoint.Ref.RunnerId}
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.getId = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.RunnerId} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.RunnerId, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.RunnerId|undefined} value
 * @return {!proto.hashicorp.waypoint.Ref.Runner} returns this
*/
proto.hashicorp.waypoint.Ref.Runner.prototype.setId = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.Ref.Runner.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Ref.Runner} returns this
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.clearId = function() {
  return this.setId(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Ref.Runner.prototype.hasId = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Ref.RunnerId.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Ref.RunnerId.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Ref.RunnerId} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.RunnerId.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Ref.RunnerId}
 */
proto.hashicorp.waypoint.Ref.RunnerId.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Ref.RunnerId;
  return proto.hashicorp.waypoint.Ref.RunnerId.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Ref.RunnerId} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Ref.RunnerId}
 */
proto.hashicorp.waypoint.Ref.RunnerId.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Ref.RunnerId.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Ref.RunnerId.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Ref.RunnerId} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.RunnerId.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Ref.RunnerId.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Ref.RunnerId} returns this
 */
proto.hashicorp.waypoint.Ref.RunnerId.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Ref.RunnerAny.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Ref.RunnerAny.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Ref.RunnerAny} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.RunnerAny.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Ref.RunnerAny}
 */
proto.hashicorp.waypoint.Ref.RunnerAny.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Ref.RunnerAny;
  return proto.hashicorp.waypoint.Ref.RunnerAny.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Ref.RunnerAny} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Ref.RunnerAny}
 */
proto.hashicorp.waypoint.Ref.RunnerAny.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Ref.RunnerAny.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Ref.RunnerAny.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Ref.RunnerAny} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Ref.RunnerAny.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Component.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Component.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Component} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Component.toObject = function(includeInstance, msg) {
  var f, obj = {
    type: jspb.Message.getFieldWithDefault(msg, 1, 0),
    name: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Component}
 */
proto.hashicorp.waypoint.Component.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Component;
  return proto.hashicorp.waypoint.Component.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Component} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Component}
 */
proto.hashicorp.waypoint.Component.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.hashicorp.waypoint.Component.Type} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Component.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Component.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Component} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Component.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.hashicorp.waypoint.Component.Type = {
  UNKNOWN: 0,
  BUILDER: 1,
  REGISTRY: 2
};

/**
 * optional Type type = 1;
 * @return {!proto.hashicorp.waypoint.Component.Type}
 */
proto.hashicorp.waypoint.Component.prototype.getType = function() {
  return /** @type {!proto.hashicorp.waypoint.Component.Type} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.Component.Type} value
 * @return {!proto.hashicorp.waypoint.Component} returns this
 */
proto.hashicorp.waypoint.Component.prototype.setType = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional string name = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.Component.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Component} returns this
 */
proto.hashicorp.waypoint.Component.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Status.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Status.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Status} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Status.toObject = function(includeInstance, msg) {
  var f, obj = {
    state: jspb.Message.getFieldWithDefault(msg, 1, 0),
    details: jspb.Message.getFieldWithDefault(msg, 2, ""),
    error: (f = msg.getError()) && google_rpc_status_pb.Status.toObject(includeInstance, f),
    startTime: (f = msg.getStartTime()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    completeTime: (f = msg.getCompleteTime()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Status}
 */
proto.hashicorp.waypoint.Status.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Status;
  return proto.hashicorp.waypoint.Status.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Status} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Status}
 */
proto.hashicorp.waypoint.Status.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.hashicorp.waypoint.Status.State} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setDetails(value);
      break;
    case 3:
      var value = new google_rpc_status_pb.Status;
      reader.readMessage(value,google_rpc_status_pb.Status.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    case 4:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setStartTime(value);
      break;
    case 5:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setCompleteTime(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Status.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Status.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Status} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Status.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getDetails();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_rpc_status_pb.Status.serializeBinaryToWriter
    );
  }
  f = message.getStartTime();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getCompleteTime();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.hashicorp.waypoint.Status.State = {
  UNKNOWN: 0,
  RUNNING: 1,
  SUCCESS: 2,
  ERROR: 3
};

/**
 * optional State state = 1;
 * @return {!proto.hashicorp.waypoint.Status.State}
 */
proto.hashicorp.waypoint.Status.prototype.getState = function() {
  return /** @type {!proto.hashicorp.waypoint.Status.State} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.Status.State} value
 * @return {!proto.hashicorp.waypoint.Status} returns this
 */
proto.hashicorp.waypoint.Status.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional string details = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.Status.prototype.getDetails = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Status} returns this
 */
proto.hashicorp.waypoint.Status.prototype.setDetails = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional google.rpc.Status error = 3;
 * @return {?proto.google.rpc.Status}
 */
proto.hashicorp.waypoint.Status.prototype.getError = function() {
  return /** @type{?proto.google.rpc.Status} */ (
    jspb.Message.getWrapperField(this, google_rpc_status_pb.Status, 3));
};


/**
 * @param {?proto.google.rpc.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.Status} returns this
*/
proto.hashicorp.waypoint.Status.prototype.setError = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Status} returns this
 */
proto.hashicorp.waypoint.Status.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Status.prototype.hasError = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional google.protobuf.Timestamp start_time = 4;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.Status.prototype.getStartTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 4));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.Status} returns this
*/
proto.hashicorp.waypoint.Status.prototype.setStartTime = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Status} returns this
 */
proto.hashicorp.waypoint.Status.prototype.clearStartTime = function() {
  return this.setStartTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Status.prototype.hasStartTime = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional google.protobuf.Timestamp complete_time = 5;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.Status.prototype.getCompleteTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 5));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.Status} returns this
*/
proto.hashicorp.waypoint.Status.prototype.setCompleteTime = function(value) {
  return jspb.Message.setWrapperField(this, 5, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Status} returns this
 */
proto.hashicorp.waypoint.Status.prototype.clearCompleteTime = function() {
  return this.setCompleteTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Status.prototype.hasCompleteTime = function() {
  return jspb.Message.getField(this, 5) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.StatusFilter.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.StatusFilter.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.StatusFilter.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.StatusFilter} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.StatusFilter.toObject = function(includeInstance, msg) {
  var f, obj = {
    filtersList: jspb.Message.toObjectList(msg.getFiltersList(),
    proto.hashicorp.waypoint.StatusFilter.Filter.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.StatusFilter}
 */
proto.hashicorp.waypoint.StatusFilter.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.StatusFilter;
  return proto.hashicorp.waypoint.StatusFilter.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.StatusFilter} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.StatusFilter}
 */
proto.hashicorp.waypoint.StatusFilter.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.StatusFilter.Filter;
      reader.readMessage(value,proto.hashicorp.waypoint.StatusFilter.Filter.deserializeBinaryFromReader);
      msg.addFilters(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.StatusFilter.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.StatusFilter.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.StatusFilter} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.StatusFilter.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFiltersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.StatusFilter.Filter.serializeBinaryToWriter
    );
  }
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.StatusFilter.Filter.oneofGroups_ = [[2]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.FilterCase = {
  FILTER_NOT_SET: 0,
  STATE: 2
};

/**
 * @return {proto.hashicorp.waypoint.StatusFilter.Filter.FilterCase}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.prototype.getFilterCase = function() {
  return /** @type {proto.hashicorp.waypoint.StatusFilter.Filter.FilterCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.StatusFilter.Filter.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.StatusFilter.Filter.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.StatusFilter.Filter} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.StatusFilter.Filter.toObject = function(includeInstance, msg) {
  var f, obj = {
    state: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.StatusFilter.Filter}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.StatusFilter.Filter;
  return proto.hashicorp.waypoint.StatusFilter.Filter.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.StatusFilter.Filter} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.StatusFilter.Filter}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = /** @type {!proto.hashicorp.waypoint.Status.State} */ (reader.readEnum());
      msg.setState(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.StatusFilter.Filter.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.StatusFilter.Filter} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.StatusFilter.Filter.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {!proto.hashicorp.waypoint.Status.State} */ (jspb.Message.getField(message, 2));
  if (f != null) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional Status.State state = 2;
 * @return {!proto.hashicorp.waypoint.Status.State}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.prototype.getState = function() {
  return /** @type {!proto.hashicorp.waypoint.Status.State} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.Status.State} value
 * @return {!proto.hashicorp.waypoint.StatusFilter.Filter} returns this
 */
proto.hashicorp.waypoint.StatusFilter.Filter.prototype.setState = function(value) {
  return jspb.Message.setOneofField(this, 2, proto.hashicorp.waypoint.StatusFilter.Filter.oneofGroups_[0], value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.hashicorp.waypoint.StatusFilter.Filter} returns this
 */
proto.hashicorp.waypoint.StatusFilter.Filter.prototype.clearState = function() {
  return jspb.Message.setOneofField(this, 2, proto.hashicorp.waypoint.StatusFilter.Filter.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.StatusFilter.Filter.prototype.hasState = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * repeated Filter filters = 1;
 * @return {!Array<!proto.hashicorp.waypoint.StatusFilter.Filter>}
 */
proto.hashicorp.waypoint.StatusFilter.prototype.getFiltersList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.StatusFilter.Filter>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.StatusFilter.Filter, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.StatusFilter.Filter>} value
 * @return {!proto.hashicorp.waypoint.StatusFilter} returns this
*/
proto.hashicorp.waypoint.StatusFilter.prototype.setFiltersList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.StatusFilter.Filter=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.StatusFilter.Filter}
 */
proto.hashicorp.waypoint.StatusFilter.prototype.addFilters = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.StatusFilter.Filter, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.StatusFilter} returns this
 */
proto.hashicorp.waypoint.StatusFilter.prototype.clearFiltersList = function() {
  return this.setFiltersList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.OperationOrder.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.OperationOrder.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.OperationOrder} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.OperationOrder.toObject = function(includeInstance, msg) {
  var f, obj = {
    order: jspb.Message.getFieldWithDefault(msg, 2, 0),
    desc: jspb.Message.getBooleanFieldWithDefault(msg, 3, false),
    limit: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.OperationOrder}
 */
proto.hashicorp.waypoint.OperationOrder.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.OperationOrder;
  return proto.hashicorp.waypoint.OperationOrder.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.OperationOrder} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.OperationOrder}
 */
proto.hashicorp.waypoint.OperationOrder.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = /** @type {!proto.hashicorp.waypoint.OperationOrder.Order} */ (reader.readEnum());
      msg.setOrder(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDesc(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setLimit(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.OperationOrder.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.OperationOrder.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.OperationOrder} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.OperationOrder.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOrder();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getDesc();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
  f = message.getLimit();
  if (f !== 0) {
    writer.writeUint32(
      4,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.hashicorp.waypoint.OperationOrder.Order = {
  UNSET: 0,
  START_TIME: 1,
  COMPLETE_TIME: 2
};

/**
 * optional Order order = 2;
 * @return {!proto.hashicorp.waypoint.OperationOrder.Order}
 */
proto.hashicorp.waypoint.OperationOrder.prototype.getOrder = function() {
  return /** @type {!proto.hashicorp.waypoint.OperationOrder.Order} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.OperationOrder.Order} value
 * @return {!proto.hashicorp.waypoint.OperationOrder} returns this
 */
proto.hashicorp.waypoint.OperationOrder.prototype.setOrder = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional bool desc = 3;
 * @return {boolean}
 */
proto.hashicorp.waypoint.OperationOrder.prototype.getDesc = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.OperationOrder} returns this
 */
proto.hashicorp.waypoint.OperationOrder.prototype.setDesc = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};


/**
 * optional uint32 limit = 4;
 * @return {number}
 */
proto.hashicorp.waypoint.OperationOrder.prototype.getLimit = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.OperationOrder} returns this
 */
proto.hashicorp.waypoint.OperationOrder.prototype.setLimit = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.QueueJobRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.QueueJobRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.QueueJobRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.QueueJobRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    job: (f = msg.getJob()) && proto.hashicorp.waypoint.Job.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.QueueJobRequest}
 */
proto.hashicorp.waypoint.QueueJobRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.QueueJobRequest;
  return proto.hashicorp.waypoint.QueueJobRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.QueueJobRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.QueueJobRequest}
 */
proto.hashicorp.waypoint.QueueJobRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Job;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.deserializeBinaryFromReader);
      msg.setJob(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.QueueJobRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.QueueJobRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.QueueJobRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.QueueJobRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getJob();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Job.serializeBinaryToWriter
    );
  }
};


/**
 * optional Job job = 1;
 * @return {?proto.hashicorp.waypoint.Job}
 */
proto.hashicorp.waypoint.QueueJobRequest.prototype.getJob = function() {
  return /** @type{?proto.hashicorp.waypoint.Job} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Job|undefined} value
 * @return {!proto.hashicorp.waypoint.QueueJobRequest} returns this
*/
proto.hashicorp.waypoint.QueueJobRequest.prototype.setJob = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.QueueJobRequest} returns this
 */
proto.hashicorp.waypoint.QueueJobRequest.prototype.clearJob = function() {
  return this.setJob(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.QueueJobRequest.prototype.hasJob = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.QueueJobResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.QueueJobResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.QueueJobResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.QueueJobResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    jobId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.QueueJobResponse}
 */
proto.hashicorp.waypoint.QueueJobResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.QueueJobResponse;
  return proto.hashicorp.waypoint.QueueJobResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.QueueJobResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.QueueJobResponse}
 */
proto.hashicorp.waypoint.QueueJobResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setJobId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.QueueJobResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.QueueJobResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.QueueJobResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.QueueJobResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getJobId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string job_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.QueueJobResponse.prototype.getJobId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.QueueJobResponse} returns this
 */
proto.hashicorp.waypoint.QueueJobResponse.prototype.setJobId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ValidateJobRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ValidateJobRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ValidateJobRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ValidateJobRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    job: (f = msg.getJob()) && proto.hashicorp.waypoint.Job.toObject(includeInstance, f),
    disableAssign: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ValidateJobRequest}
 */
proto.hashicorp.waypoint.ValidateJobRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ValidateJobRequest;
  return proto.hashicorp.waypoint.ValidateJobRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ValidateJobRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ValidateJobRequest}
 */
proto.hashicorp.waypoint.ValidateJobRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Job;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.deserializeBinaryFromReader);
      msg.setJob(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisableAssign(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ValidateJobRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ValidateJobRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ValidateJobRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ValidateJobRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getJob();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Job.serializeBinaryToWriter
    );
  }
  f = message.getDisableAssign();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * optional Job job = 1;
 * @return {?proto.hashicorp.waypoint.Job}
 */
proto.hashicorp.waypoint.ValidateJobRequest.prototype.getJob = function() {
  return /** @type{?proto.hashicorp.waypoint.Job} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Job|undefined} value
 * @return {!proto.hashicorp.waypoint.ValidateJobRequest} returns this
*/
proto.hashicorp.waypoint.ValidateJobRequest.prototype.setJob = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ValidateJobRequest} returns this
 */
proto.hashicorp.waypoint.ValidateJobRequest.prototype.clearJob = function() {
  return this.setJob(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ValidateJobRequest.prototype.hasJob = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional bool disable_assign = 2;
 * @return {boolean}
 */
proto.hashicorp.waypoint.ValidateJobRequest.prototype.getDisableAssign = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.ValidateJobRequest} returns this
 */
proto.hashicorp.waypoint.ValidateJobRequest.prototype.setDisableAssign = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ValidateJobResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ValidateJobResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ValidateJobResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    valid: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    validationError: (f = msg.getValidationError()) && google_rpc_status_pb.Status.toObject(includeInstance, f),
    assignable: jspb.Message.getBooleanFieldWithDefault(msg, 3, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ValidateJobResponse}
 */
proto.hashicorp.waypoint.ValidateJobResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ValidateJobResponse;
  return proto.hashicorp.waypoint.ValidateJobResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ValidateJobResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ValidateJobResponse}
 */
proto.hashicorp.waypoint.ValidateJobResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setValid(value);
      break;
    case 2:
      var value = new google_rpc_status_pb.Status;
      reader.readMessage(value,google_rpc_status_pb.Status.deserializeBinaryFromReader);
      msg.setValidationError(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setAssignable(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ValidateJobResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ValidateJobResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ValidateJobResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValid();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getValidationError();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      google_rpc_status_pb.Status.serializeBinaryToWriter
    );
  }
  f = message.getAssignable();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional bool valid = 1;
 * @return {boolean}
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.getValid = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.ValidateJobResponse} returns this
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.setValid = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional google.rpc.Status validation_error = 2;
 * @return {?proto.google.rpc.Status}
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.getValidationError = function() {
  return /** @type{?proto.google.rpc.Status} */ (
    jspb.Message.getWrapperField(this, google_rpc_status_pb.Status, 2));
};


/**
 * @param {?proto.google.rpc.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.ValidateJobResponse} returns this
*/
proto.hashicorp.waypoint.ValidateJobResponse.prototype.setValidationError = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ValidateJobResponse} returns this
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.clearValidationError = function() {
  return this.setValidationError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.hasValidationError = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional bool assignable = 3;
 * @return {boolean}
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.getAssignable = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.ValidateJobResponse} returns this
 */
proto.hashicorp.waypoint.ValidateJobResponse.prototype.setAssignable = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.Job.oneofGroups_ = [[20],[50,51,52,53,54,55]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.Job.DataSourceCase = {
  DATA_SOURCE_NOT_SET: 0,
  LOCAL: 20
};

/**
 * @return {proto.hashicorp.waypoint.Job.DataSourceCase}
 */
proto.hashicorp.waypoint.Job.prototype.getDataSourceCase = function() {
  return /** @type {proto.hashicorp.waypoint.Job.DataSourceCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.Job.oneofGroups_[0]));
};

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.Job.OperationCase = {
  OPERATION_NOT_SET: 0,
  NOOP: 50,
  BUILD: 51,
  PUSH: 52,
  DEPLOY: 53,
  DESTROY_DEPLOY: 54,
  RELEASE: 55
};

/**
 * @return {proto.hashicorp.waypoint.Job.OperationCase}
 */
proto.hashicorp.waypoint.Job.prototype.getOperationCase = function() {
  return /** @type {proto.hashicorp.waypoint.Job.OperationCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.Job.oneofGroups_[1]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f),
    targetRunner: (f = msg.getTargetRunner()) && proto.hashicorp.waypoint.Ref.Runner.toObject(includeInstance, f),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : [],
    local: (f = msg.getLocal()) && proto.hashicorp.waypoint.Job.Local.toObject(includeInstance, f),
    noop: (f = msg.getNoop()) && proto.hashicorp.waypoint.Job.Noop.toObject(includeInstance, f),
    build: (f = msg.getBuild()) && proto.hashicorp.waypoint.Job.BuildOp.toObject(includeInstance, f),
    push: (f = msg.getPush()) && proto.hashicorp.waypoint.Job.PushOp.toObject(includeInstance, f),
    deploy: (f = msg.getDeploy()) && proto.hashicorp.waypoint.Job.DeployOp.toObject(includeInstance, f),
    destroyDeploy: (f = msg.getDestroyDeploy()) && proto.hashicorp.waypoint.Job.DestroyDeployOp.toObject(includeInstance, f),
    release: (f = msg.getRelease()) && proto.hashicorp.waypoint.Job.ReleaseOp.toObject(includeInstance, f),
    state: jspb.Message.getFieldWithDefault(msg, 100, 0),
    assignedRunner: (f = msg.getAssignedRunner()) && proto.hashicorp.waypoint.Ref.RunnerId.toObject(includeInstance, f),
    queueTime: (f = msg.getQueueTime()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    assignTime: (f = msg.getAssignTime()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    ackTime: (f = msg.getAckTime()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    completeTime: (f = msg.getCompleteTime()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    error: (f = msg.getError()) && google_rpc_status_pb.Status.toObject(includeInstance, f),
    result: (f = msg.getResult()) && proto.hashicorp.waypoint.Job.Result.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job}
 */
proto.hashicorp.waypoint.Job.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job;
  return proto.hashicorp.waypoint.Job.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job}
 */
proto.hashicorp.waypoint.Job.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Ref.Runner;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Runner.deserializeBinaryFromReader);
      msg.setTargetRunner(value);
      break;
    case 5:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    case 20:
      var value = new proto.hashicorp.waypoint.Job.Local;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.Local.deserializeBinaryFromReader);
      msg.setLocal(value);
      break;
    case 50:
      var value = new proto.hashicorp.waypoint.Job.Noop;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.Noop.deserializeBinaryFromReader);
      msg.setNoop(value);
      break;
    case 51:
      var value = new proto.hashicorp.waypoint.Job.BuildOp;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.BuildOp.deserializeBinaryFromReader);
      msg.setBuild(value);
      break;
    case 52:
      var value = new proto.hashicorp.waypoint.Job.PushOp;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.PushOp.deserializeBinaryFromReader);
      msg.setPush(value);
      break;
    case 53:
      var value = new proto.hashicorp.waypoint.Job.DeployOp;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.DeployOp.deserializeBinaryFromReader);
      msg.setDeploy(value);
      break;
    case 54:
      var value = new proto.hashicorp.waypoint.Job.DestroyDeployOp;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.DestroyDeployOp.deserializeBinaryFromReader);
      msg.setDestroyDeploy(value);
      break;
    case 55:
      var value = new proto.hashicorp.waypoint.Job.ReleaseOp;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.ReleaseOp.deserializeBinaryFromReader);
      msg.setRelease(value);
      break;
    case 100:
      var value = /** @type {!proto.hashicorp.waypoint.Job.State} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 101:
      var value = new proto.hashicorp.waypoint.Ref.RunnerId;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.RunnerId.deserializeBinaryFromReader);
      msg.setAssignedRunner(value);
      break;
    case 102:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setQueueTime(value);
      break;
    case 103:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setAssignTime(value);
      break;
    case 104:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setAckTime(value);
      break;
    case 105:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setCompleteTime(value);
      break;
    case 106:
      var value = new google_rpc_status_pb.Status;
      reader.readMessage(value,google_rpc_status_pb.Status.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    case 107:
      var value = new proto.hashicorp.waypoint.Job.Result;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.Result.deserializeBinaryFromReader);
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
  f = message.getTargetRunner();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Ref.Runner.serializeBinaryToWriter
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(5, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getLocal();
  if (f != null) {
    writer.writeMessage(
      20,
      f,
      proto.hashicorp.waypoint.Job.Local.serializeBinaryToWriter
    );
  }
  f = message.getNoop();
  if (f != null) {
    writer.writeMessage(
      50,
      f,
      proto.hashicorp.waypoint.Job.Noop.serializeBinaryToWriter
    );
  }
  f = message.getBuild();
  if (f != null) {
    writer.writeMessage(
      51,
      f,
      proto.hashicorp.waypoint.Job.BuildOp.serializeBinaryToWriter
    );
  }
  f = message.getPush();
  if (f != null) {
    writer.writeMessage(
      52,
      f,
      proto.hashicorp.waypoint.Job.PushOp.serializeBinaryToWriter
    );
  }
  f = message.getDeploy();
  if (f != null) {
    writer.writeMessage(
      53,
      f,
      proto.hashicorp.waypoint.Job.DeployOp.serializeBinaryToWriter
    );
  }
  f = message.getDestroyDeploy();
  if (f != null) {
    writer.writeMessage(
      54,
      f,
      proto.hashicorp.waypoint.Job.DestroyDeployOp.serializeBinaryToWriter
    );
  }
  f = message.getRelease();
  if (f != null) {
    writer.writeMessage(
      55,
      f,
      proto.hashicorp.waypoint.Job.ReleaseOp.serializeBinaryToWriter
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      100,
      f
    );
  }
  f = message.getAssignedRunner();
  if (f != null) {
    writer.writeMessage(
      101,
      f,
      proto.hashicorp.waypoint.Ref.RunnerId.serializeBinaryToWriter
    );
  }
  f = message.getQueueTime();
  if (f != null) {
    writer.writeMessage(
      102,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getAssignTime();
  if (f != null) {
    writer.writeMessage(
      103,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getAckTime();
  if (f != null) {
    writer.writeMessage(
      104,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getCompleteTime();
  if (f != null) {
    writer.writeMessage(
      105,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      106,
      f,
      google_rpc_status_pb.Status.serializeBinaryToWriter
    );
  }
  f = message.getResult();
  if (f != null) {
    writer.writeMessage(
      107,
      f,
      proto.hashicorp.waypoint.Job.Result.serializeBinaryToWriter
    );
  }
};


/**
 * @enum {number}
 */
proto.hashicorp.waypoint.Job.State = {
  UNKNOWN: 0,
  QUEUED: 1,
  WAITING: 2,
  RUNNING: 3,
  ERROR: 4,
  SUCCESS: 5
};




if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.Result.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.Result.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.Result} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.Result.toObject = function(includeInstance, msg) {
  var f, obj = {
    build: (f = msg.getBuild()) && proto.hashicorp.waypoint.Job.BuildResult.toObject(includeInstance, f),
    push: (f = msg.getPush()) && proto.hashicorp.waypoint.Job.PushResult.toObject(includeInstance, f),
    deploy: (f = msg.getDeploy()) && proto.hashicorp.waypoint.Job.DeployResult.toObject(includeInstance, f),
    release: (f = msg.getRelease()) && proto.hashicorp.waypoint.Job.ReleaseResult.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.Result}
 */
proto.hashicorp.waypoint.Job.Result.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.Result;
  return proto.hashicorp.waypoint.Job.Result.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.Result} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.Result}
 */
proto.hashicorp.waypoint.Job.Result.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Job.BuildResult;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.BuildResult.deserializeBinaryFromReader);
      msg.setBuild(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Job.PushResult;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.PushResult.deserializeBinaryFromReader);
      msg.setPush(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Job.DeployResult;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.DeployResult.deserializeBinaryFromReader);
      msg.setDeploy(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Job.ReleaseResult;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.ReleaseResult.deserializeBinaryFromReader);
      msg.setRelease(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.Result.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.Result.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.Result} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.Result.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBuild();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Job.BuildResult.serializeBinaryToWriter
    );
  }
  f = message.getPush();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Job.PushResult.serializeBinaryToWriter
    );
  }
  f = message.getDeploy();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Job.DeployResult.serializeBinaryToWriter
    );
  }
  f = message.getRelease();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Job.ReleaseResult.serializeBinaryToWriter
    );
  }
};


/**
 * optional BuildResult build = 1;
 * @return {?proto.hashicorp.waypoint.Job.BuildResult}
 */
proto.hashicorp.waypoint.Job.Result.prototype.getBuild = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.BuildResult} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.BuildResult, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.BuildResult|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
*/
proto.hashicorp.waypoint.Job.Result.prototype.setBuild = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
 */
proto.hashicorp.waypoint.Job.Result.prototype.clearBuild = function() {
  return this.setBuild(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.Result.prototype.hasBuild = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional PushResult push = 2;
 * @return {?proto.hashicorp.waypoint.Job.PushResult}
 */
proto.hashicorp.waypoint.Job.Result.prototype.getPush = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.PushResult} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.PushResult, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.PushResult|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
*/
proto.hashicorp.waypoint.Job.Result.prototype.setPush = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
 */
proto.hashicorp.waypoint.Job.Result.prototype.clearPush = function() {
  return this.setPush(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.Result.prototype.hasPush = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional DeployResult deploy = 3;
 * @return {?proto.hashicorp.waypoint.Job.DeployResult}
 */
proto.hashicorp.waypoint.Job.Result.prototype.getDeploy = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.DeployResult} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.DeployResult, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.DeployResult|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
*/
proto.hashicorp.waypoint.Job.Result.prototype.setDeploy = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
 */
proto.hashicorp.waypoint.Job.Result.prototype.clearDeploy = function() {
  return this.setDeploy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.Result.prototype.hasDeploy = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional ReleaseResult release = 4;
 * @return {?proto.hashicorp.waypoint.Job.ReleaseResult}
 */
proto.hashicorp.waypoint.Job.Result.prototype.getRelease = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.ReleaseResult} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.ReleaseResult, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.ReleaseResult|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
*/
proto.hashicorp.waypoint.Job.Result.prototype.setRelease = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.Result} returns this
 */
proto.hashicorp.waypoint.Job.Result.prototype.clearRelease = function() {
  return this.setRelease(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.Result.prototype.hasRelease = function() {
  return jspb.Message.getField(this, 4) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.Local.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.Local.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.Local} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.Local.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.Local}
 */
proto.hashicorp.waypoint.Job.Local.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.Local;
  return proto.hashicorp.waypoint.Job.Local.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.Local} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.Local}
 */
proto.hashicorp.waypoint.Job.Local.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.Local.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.Local.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.Local} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.Local.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.Noop.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.Noop.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.Noop} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.Noop.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.Noop}
 */
proto.hashicorp.waypoint.Job.Noop.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.Noop;
  return proto.hashicorp.waypoint.Job.Noop.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.Noop} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.Noop}
 */
proto.hashicorp.waypoint.Job.Noop.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.Noop.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.Noop.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.Noop} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.Noop.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.BuildOp.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.BuildOp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.BuildOp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.BuildOp.toObject = function(includeInstance, msg) {
  var f, obj = {
    disablePush: jspb.Message.getBooleanFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.BuildOp}
 */
proto.hashicorp.waypoint.Job.BuildOp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.BuildOp;
  return proto.hashicorp.waypoint.Job.BuildOp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.BuildOp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.BuildOp}
 */
proto.hashicorp.waypoint.Job.BuildOp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setDisablePush(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.BuildOp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.BuildOp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.BuildOp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.BuildOp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDisablePush();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool disable_push = 1;
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.BuildOp.prototype.getDisablePush = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.Job.BuildOp} returns this
 */
proto.hashicorp.waypoint.Job.BuildOp.prototype.setDisablePush = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.BuildResult.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.BuildResult} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.BuildResult.toObject = function(includeInstance, msg) {
  var f, obj = {
    build: (f = msg.getBuild()) && proto.hashicorp.waypoint.Build.toObject(includeInstance, f),
    push: (f = msg.getPush()) && proto.hashicorp.waypoint.PushedArtifact.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.BuildResult}
 */
proto.hashicorp.waypoint.Job.BuildResult.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.BuildResult;
  return proto.hashicorp.waypoint.Job.BuildResult.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.BuildResult} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.BuildResult}
 */
proto.hashicorp.waypoint.Job.BuildResult.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Build;
      reader.readMessage(value,proto.hashicorp.waypoint.Build.deserializeBinaryFromReader);
      msg.setBuild(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.PushedArtifact;
      reader.readMessage(value,proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader);
      msg.setPush(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.BuildResult.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.BuildResult} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.BuildResult.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBuild();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Build.serializeBinaryToWriter
    );
  }
  f = message.getPush();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter
    );
  }
};


/**
 * optional Build build = 1;
 * @return {?proto.hashicorp.waypoint.Build}
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.getBuild = function() {
  return /** @type{?proto.hashicorp.waypoint.Build} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Build, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Build|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.BuildResult} returns this
*/
proto.hashicorp.waypoint.Job.BuildResult.prototype.setBuild = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.BuildResult} returns this
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.clearBuild = function() {
  return this.setBuild(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.hasBuild = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional PushedArtifact push = 2;
 * @return {?proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.getPush = function() {
  return /** @type{?proto.hashicorp.waypoint.PushedArtifact} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.PushedArtifact, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.PushedArtifact|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.BuildResult} returns this
*/
proto.hashicorp.waypoint.Job.BuildResult.prototype.setPush = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.BuildResult} returns this
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.clearPush = function() {
  return this.setPush(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.BuildResult.prototype.hasPush = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.PushOp.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.PushOp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.PushOp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.PushOp.toObject = function(includeInstance, msg) {
  var f, obj = {
    build: (f = msg.getBuild()) && proto.hashicorp.waypoint.Build.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.PushOp}
 */
proto.hashicorp.waypoint.Job.PushOp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.PushOp;
  return proto.hashicorp.waypoint.Job.PushOp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.PushOp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.PushOp}
 */
proto.hashicorp.waypoint.Job.PushOp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Build;
      reader.readMessage(value,proto.hashicorp.waypoint.Build.deserializeBinaryFromReader);
      msg.setBuild(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.PushOp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.PushOp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.PushOp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.PushOp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBuild();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Build.serializeBinaryToWriter
    );
  }
};


/**
 * optional Build build = 1;
 * @return {?proto.hashicorp.waypoint.Build}
 */
proto.hashicorp.waypoint.Job.PushOp.prototype.getBuild = function() {
  return /** @type{?proto.hashicorp.waypoint.Build} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Build, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Build|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.PushOp} returns this
*/
proto.hashicorp.waypoint.Job.PushOp.prototype.setBuild = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.PushOp} returns this
 */
proto.hashicorp.waypoint.Job.PushOp.prototype.clearBuild = function() {
  return this.setBuild(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.PushOp.prototype.hasBuild = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.PushResult.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.PushResult.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.PushResult} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.PushResult.toObject = function(includeInstance, msg) {
  var f, obj = {
    artifact: (f = msg.getArtifact()) && proto.hashicorp.waypoint.PushedArtifact.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.PushResult}
 */
proto.hashicorp.waypoint.Job.PushResult.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.PushResult;
  return proto.hashicorp.waypoint.Job.PushResult.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.PushResult} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.PushResult}
 */
proto.hashicorp.waypoint.Job.PushResult.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.PushedArtifact;
      reader.readMessage(value,proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader);
      msg.setArtifact(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.PushResult.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.PushResult.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.PushResult} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.PushResult.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getArtifact();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter
    );
  }
};


/**
 * optional PushedArtifact artifact = 1;
 * @return {?proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.Job.PushResult.prototype.getArtifact = function() {
  return /** @type{?proto.hashicorp.waypoint.PushedArtifact} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.PushedArtifact, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.PushedArtifact|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.PushResult} returns this
*/
proto.hashicorp.waypoint.Job.PushResult.prototype.setArtifact = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.PushResult} returns this
 */
proto.hashicorp.waypoint.Job.PushResult.prototype.clearArtifact = function() {
  return this.setArtifact(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.PushResult.prototype.hasArtifact = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.DeployOp.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.DeployOp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.DeployOp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.DeployOp.toObject = function(includeInstance, msg) {
  var f, obj = {
    artifact: (f = msg.getArtifact()) && proto.hashicorp.waypoint.PushedArtifact.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.DeployOp}
 */
proto.hashicorp.waypoint.Job.DeployOp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.DeployOp;
  return proto.hashicorp.waypoint.Job.DeployOp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.DeployOp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.DeployOp}
 */
proto.hashicorp.waypoint.Job.DeployOp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.PushedArtifact;
      reader.readMessage(value,proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader);
      msg.setArtifact(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.DeployOp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.DeployOp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.DeployOp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.DeployOp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getArtifact();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter
    );
  }
};


/**
 * optional PushedArtifact artifact = 1;
 * @return {?proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.Job.DeployOp.prototype.getArtifact = function() {
  return /** @type{?proto.hashicorp.waypoint.PushedArtifact} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.PushedArtifact, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.PushedArtifact|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.DeployOp} returns this
*/
proto.hashicorp.waypoint.Job.DeployOp.prototype.setArtifact = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.DeployOp} returns this
 */
proto.hashicorp.waypoint.Job.DeployOp.prototype.clearArtifact = function() {
  return this.setArtifact(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.DeployOp.prototype.hasArtifact = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.DeployResult.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.DeployResult.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.DeployResult} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.DeployResult.toObject = function(includeInstance, msg) {
  var f, obj = {
    deployment: (f = msg.getDeployment()) && proto.hashicorp.waypoint.Deployment.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.DeployResult}
 */
proto.hashicorp.waypoint.Job.DeployResult.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.DeployResult;
  return proto.hashicorp.waypoint.Job.DeployResult.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.DeployResult} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.DeployResult}
 */
proto.hashicorp.waypoint.Job.DeployResult.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Deployment;
      reader.readMessage(value,proto.hashicorp.waypoint.Deployment.deserializeBinaryFromReader);
      msg.setDeployment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.DeployResult.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.DeployResult.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.DeployResult} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.DeployResult.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeployment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Deployment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Deployment deployment = 1;
 * @return {?proto.hashicorp.waypoint.Deployment}
 */
proto.hashicorp.waypoint.Job.DeployResult.prototype.getDeployment = function() {
  return /** @type{?proto.hashicorp.waypoint.Deployment} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Deployment, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Deployment|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.DeployResult} returns this
*/
proto.hashicorp.waypoint.Job.DeployResult.prototype.setDeployment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.DeployResult} returns this
 */
proto.hashicorp.waypoint.Job.DeployResult.prototype.clearDeployment = function() {
  return this.setDeployment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.DeployResult.prototype.hasDeployment = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.DestroyDeployOp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.DestroyDeployOp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.toObject = function(includeInstance, msg) {
  var f, obj = {
    deployment: (f = msg.getDeployment()) && proto.hashicorp.waypoint.Deployment.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.DestroyDeployOp}
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.DestroyDeployOp;
  return proto.hashicorp.waypoint.Job.DestroyDeployOp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.DestroyDeployOp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.DestroyDeployOp}
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Deployment;
      reader.readMessage(value,proto.hashicorp.waypoint.Deployment.deserializeBinaryFromReader);
      msg.setDeployment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.DestroyDeployOp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.DestroyDeployOp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeployment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Deployment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Deployment deployment = 1;
 * @return {?proto.hashicorp.waypoint.Deployment}
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.prototype.getDeployment = function() {
  return /** @type{?proto.hashicorp.waypoint.Deployment} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Deployment, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Deployment|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.DestroyDeployOp} returns this
*/
proto.hashicorp.waypoint.Job.DestroyDeployOp.prototype.setDeployment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.DestroyDeployOp} returns this
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.prototype.clearDeployment = function() {
  return this.setDeployment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.DestroyDeployOp.prototype.hasDeployment = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.ReleaseOp.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.ReleaseOp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.ReleaseOp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.ReleaseOp.toObject = function(includeInstance, msg) {
  var f, obj = {
    trafficSplit: (f = msg.getTrafficSplit()) && proto.hashicorp.waypoint.Release.Split.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.ReleaseOp}
 */
proto.hashicorp.waypoint.Job.ReleaseOp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.ReleaseOp;
  return proto.hashicorp.waypoint.Job.ReleaseOp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.ReleaseOp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.ReleaseOp}
 */
proto.hashicorp.waypoint.Job.ReleaseOp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Release.Split;
      reader.readMessage(value,proto.hashicorp.waypoint.Release.Split.deserializeBinaryFromReader);
      msg.setTrafficSplit(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.ReleaseOp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.ReleaseOp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.ReleaseOp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.ReleaseOp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTrafficSplit();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Release.Split.serializeBinaryToWriter
    );
  }
};


/**
 * optional Release.Split traffic_split = 1;
 * @return {?proto.hashicorp.waypoint.Release.Split}
 */
proto.hashicorp.waypoint.Job.ReleaseOp.prototype.getTrafficSplit = function() {
  return /** @type{?proto.hashicorp.waypoint.Release.Split} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Release.Split, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Release.Split|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.ReleaseOp} returns this
*/
proto.hashicorp.waypoint.Job.ReleaseOp.prototype.setTrafficSplit = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.ReleaseOp} returns this
 */
proto.hashicorp.waypoint.Job.ReleaseOp.prototype.clearTrafficSplit = function() {
  return this.setTrafficSplit(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.ReleaseOp.prototype.hasTrafficSplit = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Job.ReleaseResult.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Job.ReleaseResult.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Job.ReleaseResult} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.ReleaseResult.toObject = function(includeInstance, msg) {
  var f, obj = {
    release: (f = msg.getRelease()) && proto.hashicorp.waypoint.Release.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Job.ReleaseResult}
 */
proto.hashicorp.waypoint.Job.ReleaseResult.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Job.ReleaseResult;
  return proto.hashicorp.waypoint.Job.ReleaseResult.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Job.ReleaseResult} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Job.ReleaseResult}
 */
proto.hashicorp.waypoint.Job.ReleaseResult.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Release;
      reader.readMessage(value,proto.hashicorp.waypoint.Release.deserializeBinaryFromReader);
      msg.setRelease(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Job.ReleaseResult.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Job.ReleaseResult.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Job.ReleaseResult} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Job.ReleaseResult.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRelease();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Release.serializeBinaryToWriter
    );
  }
};


/**
 * optional Release release = 1;
 * @return {?proto.hashicorp.waypoint.Release}
 */
proto.hashicorp.waypoint.Job.ReleaseResult.prototype.getRelease = function() {
  return /** @type{?proto.hashicorp.waypoint.Release} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Release, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Release|undefined} value
 * @return {!proto.hashicorp.waypoint.Job.ReleaseResult} returns this
*/
proto.hashicorp.waypoint.Job.ReleaseResult.prototype.setRelease = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job.ReleaseResult} returns this
 */
proto.hashicorp.waypoint.Job.ReleaseResult.prototype.clearRelease = function() {
  return this.setRelease(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.ReleaseResult.prototype.hasRelease = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Job.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Ref.Application application = 2;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.Job.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Ref.Workspace workspace = 3;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.Job.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Ref.Runner target_runner = 4;
 * @return {?proto.hashicorp.waypoint.Ref.Runner}
 */
proto.hashicorp.waypoint.Job.prototype.getTargetRunner = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Runner} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Runner, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Runner|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setTargetRunner = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearTargetRunner = function() {
  return this.setTargetRunner(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasTargetRunner = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * map<string, string> labels = 5;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.hashicorp.waypoint.Job.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 5, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;};


/**
 * optional Local local = 20;
 * @return {?proto.hashicorp.waypoint.Job.Local}
 */
proto.hashicorp.waypoint.Job.prototype.getLocal = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.Local} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.Local, 20));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.Local|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setLocal = function(value) {
  return jspb.Message.setOneofWrapperField(this, 20, proto.hashicorp.waypoint.Job.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearLocal = function() {
  return this.setLocal(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasLocal = function() {
  return jspb.Message.getField(this, 20) != null;
};


/**
 * optional Noop noop = 50;
 * @return {?proto.hashicorp.waypoint.Job.Noop}
 */
proto.hashicorp.waypoint.Job.prototype.getNoop = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.Noop} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.Noop, 50));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.Noop|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setNoop = function(value) {
  return jspb.Message.setOneofWrapperField(this, 50, proto.hashicorp.waypoint.Job.oneofGroups_[1], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearNoop = function() {
  return this.setNoop(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasNoop = function() {
  return jspb.Message.getField(this, 50) != null;
};


/**
 * optional BuildOp build = 51;
 * @return {?proto.hashicorp.waypoint.Job.BuildOp}
 */
proto.hashicorp.waypoint.Job.prototype.getBuild = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.BuildOp} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.BuildOp, 51));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.BuildOp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setBuild = function(value) {
  return jspb.Message.setOneofWrapperField(this, 51, proto.hashicorp.waypoint.Job.oneofGroups_[1], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearBuild = function() {
  return this.setBuild(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasBuild = function() {
  return jspb.Message.getField(this, 51) != null;
};


/**
 * optional PushOp push = 52;
 * @return {?proto.hashicorp.waypoint.Job.PushOp}
 */
proto.hashicorp.waypoint.Job.prototype.getPush = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.PushOp} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.PushOp, 52));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.PushOp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setPush = function(value) {
  return jspb.Message.setOneofWrapperField(this, 52, proto.hashicorp.waypoint.Job.oneofGroups_[1], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearPush = function() {
  return this.setPush(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasPush = function() {
  return jspb.Message.getField(this, 52) != null;
};


/**
 * optional DeployOp deploy = 53;
 * @return {?proto.hashicorp.waypoint.Job.DeployOp}
 */
proto.hashicorp.waypoint.Job.prototype.getDeploy = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.DeployOp} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.DeployOp, 53));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.DeployOp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setDeploy = function(value) {
  return jspb.Message.setOneofWrapperField(this, 53, proto.hashicorp.waypoint.Job.oneofGroups_[1], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearDeploy = function() {
  return this.setDeploy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasDeploy = function() {
  return jspb.Message.getField(this, 53) != null;
};


/**
 * optional DestroyDeployOp destroy_deploy = 54;
 * @return {?proto.hashicorp.waypoint.Job.DestroyDeployOp}
 */
proto.hashicorp.waypoint.Job.prototype.getDestroyDeploy = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.DestroyDeployOp} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.DestroyDeployOp, 54));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.DestroyDeployOp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setDestroyDeploy = function(value) {
  return jspb.Message.setOneofWrapperField(this, 54, proto.hashicorp.waypoint.Job.oneofGroups_[1], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearDestroyDeploy = function() {
  return this.setDestroyDeploy(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasDestroyDeploy = function() {
  return jspb.Message.getField(this, 54) != null;
};


/**
 * optional ReleaseOp release = 55;
 * @return {?proto.hashicorp.waypoint.Job.ReleaseOp}
 */
proto.hashicorp.waypoint.Job.prototype.getRelease = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.ReleaseOp} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.ReleaseOp, 55));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.ReleaseOp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setRelease = function(value) {
  return jspb.Message.setOneofWrapperField(this, 55, proto.hashicorp.waypoint.Job.oneofGroups_[1], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearRelease = function() {
  return this.setRelease(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasRelease = function() {
  return jspb.Message.getField(this, 55) != null;
};


/**
 * optional State state = 100;
 * @return {!proto.hashicorp.waypoint.Job.State}
 */
proto.hashicorp.waypoint.Job.prototype.getState = function() {
  return /** @type {!proto.hashicorp.waypoint.Job.State} */ (jspb.Message.getFieldWithDefault(this, 100, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.Job.State} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 100, value);
};


/**
 * optional Ref.RunnerId assigned_runner = 101;
 * @return {?proto.hashicorp.waypoint.Ref.RunnerId}
 */
proto.hashicorp.waypoint.Job.prototype.getAssignedRunner = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.RunnerId} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.RunnerId, 101));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.RunnerId|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setAssignedRunner = function(value) {
  return jspb.Message.setWrapperField(this, 101, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearAssignedRunner = function() {
  return this.setAssignedRunner(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasAssignedRunner = function() {
  return jspb.Message.getField(this, 101) != null;
};


/**
 * optional google.protobuf.Timestamp queue_time = 102;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.Job.prototype.getQueueTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 102));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setQueueTime = function(value) {
  return jspb.Message.setWrapperField(this, 102, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearQueueTime = function() {
  return this.setQueueTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasQueueTime = function() {
  return jspb.Message.getField(this, 102) != null;
};


/**
 * optional google.protobuf.Timestamp assign_time = 103;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.Job.prototype.getAssignTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 103));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setAssignTime = function(value) {
  return jspb.Message.setWrapperField(this, 103, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearAssignTime = function() {
  return this.setAssignTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasAssignTime = function() {
  return jspb.Message.getField(this, 103) != null;
};


/**
 * optional google.protobuf.Timestamp ack_time = 104;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.Job.prototype.getAckTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 104));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setAckTime = function(value) {
  return jspb.Message.setWrapperField(this, 104, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearAckTime = function() {
  return this.setAckTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasAckTime = function() {
  return jspb.Message.getField(this, 104) != null;
};


/**
 * optional google.protobuf.Timestamp complete_time = 105;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.Job.prototype.getCompleteTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 105));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setCompleteTime = function(value) {
  return jspb.Message.setWrapperField(this, 105, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearCompleteTime = function() {
  return this.setCompleteTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasCompleteTime = function() {
  return jspb.Message.getField(this, 105) != null;
};


/**
 * optional google.rpc.Status error = 106;
 * @return {?proto.google.rpc.Status}
 */
proto.hashicorp.waypoint.Job.prototype.getError = function() {
  return /** @type{?proto.google.rpc.Status} */ (
    jspb.Message.getWrapperField(this, google_rpc_status_pb.Status, 106));
};


/**
 * @param {?proto.google.rpc.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setError = function(value) {
  return jspb.Message.setWrapperField(this, 106, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasError = function() {
  return jspb.Message.getField(this, 106) != null;
};


/**
 * optional Result result = 107;
 * @return {?proto.hashicorp.waypoint.Job.Result}
 */
proto.hashicorp.waypoint.Job.prototype.getResult = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.Result} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.Result, 107));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.Result|undefined} value
 * @return {!proto.hashicorp.waypoint.Job} returns this
*/
proto.hashicorp.waypoint.Job.prototype.setResult = function(value) {
  return jspb.Message.setWrapperField(this, 107, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Job} returns this
 */
proto.hashicorp.waypoint.Job.prototype.clearResult = function() {
  return this.setResult(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Job.prototype.hasResult = function() {
  return jspb.Message.getField(this, 107) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    jobId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobRequest}
 */
proto.hashicorp.waypoint.GetJobRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobRequest;
  return proto.hashicorp.waypoint.GetJobRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobRequest}
 */
proto.hashicorp.waypoint.GetJobRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setJobId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getJobId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string job_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobRequest.prototype.getJobId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobRequest} returns this
 */
proto.hashicorp.waypoint.GetJobRequest.prototype.setJobId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    jobId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamRequest}
 */
proto.hashicorp.waypoint.GetJobStreamRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamRequest;
  return proto.hashicorp.waypoint.GetJobStreamRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamRequest}
 */
proto.hashicorp.waypoint.GetJobStreamRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setJobId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getJobId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string job_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamRequest.prototype.getJobId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamRequest} returns this
 */
proto.hashicorp.waypoint.GetJobStreamRequest.prototype.setJobId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_ = [[1,2,3,4,5]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.EventCase = {
  EVENT_NOT_SET: 0,
  OPEN: 1,
  STATE: 2,
  TERMINAL: 3,
  ERROR: 4,
  COMPLETE: 5
};

/**
 * @return {proto.hashicorp.waypoint.GetJobStreamResponse.EventCase}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.GetJobStreamResponse.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    open: (f = msg.getOpen()) && proto.hashicorp.waypoint.GetJobStreamResponse.Open.toObject(includeInstance, f),
    state: (f = msg.getState()) && proto.hashicorp.waypoint.GetJobStreamResponse.State.toObject(includeInstance, f),
    terminal: (f = msg.getTerminal()) && proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.toObject(includeInstance, f),
    error: (f = msg.getError()) && proto.hashicorp.waypoint.GetJobStreamResponse.Error.toObject(includeInstance, f),
    complete: (f = msg.getComplete()) && proto.hashicorp.waypoint.GetJobStreamResponse.Complete.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse;
  return proto.hashicorp.waypoint.GetJobStreamResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Open;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Open.deserializeBinaryFromReader);
      msg.setOpen(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.State;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.State.deserializeBinaryFromReader);
      msg.setState(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.deserializeBinaryFromReader);
      msg.setTerminal(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Error;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Error.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    case 5:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Complete;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Complete.deserializeBinaryFromReader);
      msg.setComplete(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOpen();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Open.serializeBinaryToWriter
    );
  }
  f = message.getState();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.State.serializeBinaryToWriter
    );
  }
  f = message.getTerminal();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.serializeBinaryToWriter
    );
  }
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Error.serializeBinaryToWriter
    );
  }
  f = message.getComplete();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Complete.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Open.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Open.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Open} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Open.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Open}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Open.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Open;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Open.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Open} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Open}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Open.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Open.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Open.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Open} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Open.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.State.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.State} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.toObject = function(includeInstance, msg) {
  var f, obj = {
    previous: jspb.Message.getFieldWithDefault(msg, 1, 0),
    current: jspb.Message.getFieldWithDefault(msg, 2, 0),
    job: (f = msg.getJob()) && proto.hashicorp.waypoint.Job.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.State}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.State;
  return proto.hashicorp.waypoint.GetJobStreamResponse.State.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.State} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.State}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.hashicorp.waypoint.Job.State} */ (reader.readEnum());
      msg.setPrevious(value);
      break;
    case 2:
      var value = /** @type {!proto.hashicorp.waypoint.Job.State} */ (reader.readEnum());
      msg.setCurrent(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Job;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.deserializeBinaryFromReader);
      msg.setJob(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.State.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.State} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPrevious();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getCurrent();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getJob();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Job.serializeBinaryToWriter
    );
  }
};


/**
 * optional Job.State previous = 1;
 * @return {!proto.hashicorp.waypoint.Job.State}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.getPrevious = function() {
  return /** @type {!proto.hashicorp.waypoint.Job.State} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.Job.State} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.State} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.setPrevious = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional Job.State current = 2;
 * @return {!proto.hashicorp.waypoint.Job.State}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.getCurrent = function() {
  return /** @type {!proto.hashicorp.waypoint.Job.State} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.Job.State} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.State} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.setCurrent = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional Job job = 3;
 * @return {?proto.hashicorp.waypoint.Job}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.getJob = function() {
  return /** @type{?proto.hashicorp.waypoint.Job} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Job|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.State} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.setJob = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.State} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.clearJob = function() {
  return this.setJob(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.State.prototype.hasJob = function() {
  return jspb.Message.getField(this, 3) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.toObject = function(includeInstance, msg) {
  var f, obj = {
    eventsList: jspb.Message.toObjectList(msg.getEventsList(),
    proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.toObject, includeInstance),
    buffered: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.deserializeBinaryFromReader);
      msg.addEvents(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setBuffered(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEventsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.serializeBinaryToWriter
    );
  }
  f = message.getBuffered();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_ = [[2,3,4,5,6]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.EventCase = {
  EVENT_NOT_SET: 0,
  LINE: 2,
  STATUS: 3,
  NAMED_VALUES: 4,
  RAW: 5,
  TABLE: 6
};

/**
 * @return {proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.EventCase}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.toObject = function(includeInstance, msg) {
  var f, obj = {
    timestamp: (f = msg.getTimestamp()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    line: (f = msg.getLine()) && proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.toObject(includeInstance, f),
    status: (f = msg.getStatus()) && proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.toObject(includeInstance, f),
    namedValues: (f = msg.getNamedValues()) && proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.toObject(includeInstance, f),
    raw: (f = msg.getRaw()) && proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.toObject(includeInstance, f),
    table: (f = msg.getTable()) && proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setTimestamp(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.deserializeBinaryFromReader);
      msg.setLine(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.deserializeBinaryFromReader);
      msg.setStatus(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.deserializeBinaryFromReader);
      msg.setNamedValues(value);
      break;
    case 5:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.deserializeBinaryFromReader);
      msg.setRaw(value);
      break;
    case 6:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.deserializeBinaryFromReader);
      msg.setTable(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTimestamp();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getLine();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.serializeBinaryToWriter
    );
  }
  f = message.getStatus();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.serializeBinaryToWriter
    );
  }
  f = message.getNamedValues();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.serializeBinaryToWriter
    );
  }
  f = message.getRaw();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.serializeBinaryToWriter
    );
  }
  f = message.getTable();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.toObject = function(includeInstance, msg) {
  var f, obj = {
    status: jspb.Message.getFieldWithDefault(msg, 1, ""),
    msg: jspb.Message.getFieldWithDefault(msg, 2, ""),
    step: jspb.Message.getBooleanFieldWithDefault(msg, 3, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setStatus(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setMsg(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setStep(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStatus();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getMsg();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getStep();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional string status = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.getStatus = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.setStatus = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string msg = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.getMsg = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.setMsg = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool step = 3;
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.getStep = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status.prototype.setStep = function(value) {
  return jspb.Message.setProto3BooleanField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.toObject = function(includeInstance, msg) {
  var f, obj = {
    msg: jspb.Message.getFieldWithDefault(msg, 1, ""),
    style: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setMsg(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setStyle(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMsg();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStyle();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string msg = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.prototype.getMsg = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.prototype.setMsg = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string style = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.prototype.getStyle = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line.prototype.setStyle = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.toObject = function(includeInstance, msg) {
  var f, obj = {
    data: msg.getData_asB64(),
    stderr: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setData(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setStderr(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getData_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      1,
      f
    );
  }
  f = message.getStderr();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * optional bytes data = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * optional bytes data = 1;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.getData_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getData()));
};


/**
 * optional bytes data = 1;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getData()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.setData = function(value) {
  return jspb.Message.setProto3BytesField(this, 1, value);
};


/**
 * optional bool stderr = 2;
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.getStderr = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw.prototype.setStderr = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.toObject = function(includeInstance, msg) {
  var f, obj = {
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    value: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getValue();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string value = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.prototype.setValue = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.toObject = function(includeInstance, msg) {
  var f, obj = {
    valuesList: jspb.Message.toObjectList(msg.getValuesList(),
    proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.deserializeBinaryFromReader);
      msg.addValues(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValuesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue.serializeBinaryToWriter
    );
  }
};


/**
 * repeated NamedValue values = 1;
 * @return {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue>}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.prototype.getValuesList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue>} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.prototype.setValuesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.prototype.addValues = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValue, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues.prototype.clearValuesList = function() {
  return this.setValuesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.toObject = function(includeInstance, msg) {
  var f, obj = {
    value: jspb.Message.getFieldWithDefault(msg, 1, ""),
    color: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setColor(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getValue();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getColor();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string value = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.prototype.setValue = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string color = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.prototype.getColor = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.prototype.setColor = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.toObject = function(includeInstance, msg) {
  var f, obj = {
    entriesList: jspb.Message.toObjectList(msg.getEntriesList(),
    proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.deserializeBinaryFromReader);
      msg.addEntries(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEntriesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry.serializeBinaryToWriter
    );
  }
};


/**
 * repeated TableEntry entries = 1;
 * @return {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry>}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.prototype.getEntriesList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry>} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.prototype.setEntriesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.prototype.addEntries = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.prototype.clearEntriesList = function() {
  return this.setEntriesList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.repeatedFields_ = [1,2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.toObject = function(includeInstance, msg) {
  var f, obj = {
    headersList: (f = jspb.Message.getRepeatedField(msg, 1)) == null ? undefined : f,
    rowsList: jspb.Message.toObjectList(msg.getRowsList(),
    proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addHeaders(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.deserializeBinaryFromReader);
      msg.addRows(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getHeadersList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
  f = message.getRowsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow.serializeBinaryToWriter
    );
  }
};


/**
 * repeated string headers = 1;
 * @return {!Array<string>}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.getHeadersList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.setHeadersList = function(value) {
  return jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.addHeaders = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.clearHeadersList = function() {
  return this.setHeadersList([]);
};


/**
 * repeated TableRow rows = 2;
 * @return {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow>}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.getRowsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow, 2));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow>} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.setRowsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.addRows = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableRow, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table.prototype.clearRowsList = function() {
  return this.setRowsList([]);
};


/**
 * optional google.protobuf.Timestamp timestamp = 1;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.getTimestamp = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 1));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.setTimestamp = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.clearTimestamp = function() {
  return this.setTimestamp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.hasTimestamp = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Line line = 2;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.getLine = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Line|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.setLine = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.clearLine = function() {
  return this.setLine(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.hasLine = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Status status = 3;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.getStatus = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.setStatus = function(value) {
  return jspb.Message.setOneofWrapperField(this, 3, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.clearStatus = function() {
  return this.setStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.hasStatus = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional NamedValues named_values = 4;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.getNamedValues = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.NamedValues|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.setNamedValues = function(value) {
  return jspb.Message.setOneofWrapperField(this, 4, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.clearNamedValues = function() {
  return this.setNamedValues(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.hasNamedValues = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional Raw raw = 5;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.getRaw = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw, 5));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Raw|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.setRaw = function(value) {
  return jspb.Message.setOneofWrapperField(this, 5, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.clearRaw = function() {
  return this.setRaw(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.hasRaw = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional Table table = 6;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.getTable = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table, 6));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.Table|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.setTable = function(value) {
  return jspb.Message.setOneofWrapperField(this, 6, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.clearTable = function() {
  return this.setTable(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.prototype.hasTable = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * repeated Event events = 1;
 * @return {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event>}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.getEventsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event>} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.setEventsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.addEvents = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.Event, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.clearEventsList = function() {
  return this.setEventsList([]);
};


/**
 * optional bool buffered = 2;
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.getBuffered = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.prototype.setBuffered = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Error.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Error} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.toObject = function(includeInstance, msg) {
  var f, obj = {
    error: (f = msg.getError()) && google_rpc_status_pb.Status.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Error}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Error;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Error.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Error} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Error}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_rpc_status_pb.Status;
      reader.readMessage(value,google_rpc_status_pb.Status.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Error.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Error} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_rpc_status_pb.Status.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.rpc.Status error = 1;
 * @return {?proto.google.rpc.Status}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.prototype.getError = function() {
  return /** @type{?proto.google.rpc.Status} */ (
    jspb.Message.getWrapperField(this, google_rpc_status_pb.Status, 1));
};


/**
 * @param {?proto.google.rpc.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Error} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Error.prototype.setError = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Error} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Error.prototype.hasError = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetJobStreamResponse.Complete.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.toObject = function(includeInstance, msg) {
  var f, obj = {
    error: (f = msg.getError()) && google_rpc_status_pb.Status.toObject(includeInstance, f),
    result: (f = msg.getResult()) && proto.hashicorp.waypoint.Job.Result.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetJobStreamResponse.Complete;
  return proto.hashicorp.waypoint.GetJobStreamResponse.Complete.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_rpc_status_pb.Status;
      reader.readMessage(value,google_rpc_status_pb.Status.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Job.Result;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.Result.deserializeBinaryFromReader);
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetJobStreamResponse.Complete.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_rpc_status_pb.Status.serializeBinaryToWriter
    );
  }
  f = message.getResult();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Job.Result.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.rpc.Status error = 1;
 * @return {?proto.google.rpc.Status}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.getError = function() {
  return /** @type{?proto.google.rpc.Status} */ (
    jspb.Message.getWrapperField(this, google_rpc_status_pb.Status, 1));
};


/**
 * @param {?proto.google.rpc.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.setError = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.hasError = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Job.Result result = 2;
 * @return {?proto.hashicorp.waypoint.Job.Result}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.getResult = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.Result} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.Result, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.Result|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.setResult = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse.Complete} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.clearResult = function() {
  return this.setResult(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.Complete.prototype.hasResult = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Open open = 1;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Open}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.getOpen = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Open} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Open, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Open|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.setOpen = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.clearOpen = function() {
  return this.setOpen(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.hasOpen = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional State state = 2;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.State}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.getState = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.State} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.State, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.State|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.setState = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.clearState = function() {
  return this.setState(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.hasState = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Terminal terminal = 3;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.getTerminal = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.setTerminal = function(value) {
  return jspb.Message.setOneofWrapperField(this, 3, proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.clearTerminal = function() {
  return this.setTerminal(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.hasTerminal = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Error error = 4;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Error}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.getError = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Error} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Error, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Error|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.setError = function(value) {
  return jspb.Message.setOneofWrapperField(this, 4, proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.hasError = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional Complete complete = 5;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Complete}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.getComplete = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Complete} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Complete, 5));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Complete|undefined} value
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
*/
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.setComplete = function(value) {
  return jspb.Message.setOneofWrapperField(this, 5, proto.hashicorp.waypoint.GetJobStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetJobStreamResponse} returns this
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.clearComplete = function() {
  return this.setComplete(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetJobStreamResponse.prototype.hasComplete = function() {
  return jspb.Message.getField(this, 5) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.Runner.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Runner.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Runner.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Runner} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Runner.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    byIdOnly: jspb.Message.getBooleanFieldWithDefault(msg, 2, false),
    componentsList: jspb.Message.toObjectList(msg.getComponentsList(),
    proto.hashicorp.waypoint.Component.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Runner}
 */
proto.hashicorp.waypoint.Runner.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Runner;
  return proto.hashicorp.waypoint.Runner.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Runner} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Runner}
 */
proto.hashicorp.waypoint.Runner.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setByIdOnly(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Component;
      reader.readMessage(value,proto.hashicorp.waypoint.Component.deserializeBinaryFromReader);
      msg.addComponents(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Runner.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Runner.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Runner} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Runner.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getByIdOnly();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
  f = message.getComponentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.hashicorp.waypoint.Component.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Runner.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Runner} returns this
 */
proto.hashicorp.waypoint.Runner.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bool by_id_only = 2;
 * @return {boolean}
 */
proto.hashicorp.waypoint.Runner.prototype.getByIdOnly = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.Runner} returns this
 */
proto.hashicorp.waypoint.Runner.prototype.setByIdOnly = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};


/**
 * repeated Component components = 3;
 * @return {!Array<!proto.hashicorp.waypoint.Component>}
 */
proto.hashicorp.waypoint.Runner.prototype.getComponentsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.Component>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.Component, 3));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.Component>} value
 * @return {!proto.hashicorp.waypoint.Runner} returns this
*/
proto.hashicorp.waypoint.Runner.prototype.setComponentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.hashicorp.waypoint.Component=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.Component}
 */
proto.hashicorp.waypoint.Runner.prototype.addComponents = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.hashicorp.waypoint.Component, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.Runner} returns this
 */
proto.hashicorp.waypoint.Runner.prototype.clearComponentsList = function() {
  return this.setComponentsList([]);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.RunnerConfigRequest.oneofGroups_ = [[1]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.EventCase = {
  EVENT_NOT_SET: 0,
  OPEN: 1
};

/**
 * @return {proto.hashicorp.waypoint.RunnerConfigRequest.EventCase}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.RunnerConfigRequest.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.RunnerConfigRequest.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerConfigRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerConfigRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfigRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    open: (f = msg.getOpen()) && proto.hashicorp.waypoint.RunnerConfigRequest.Open.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerConfigRequest;
  return proto.hashicorp.waypoint.RunnerConfigRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerConfigRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.RunnerConfigRequest.Open;
      reader.readMessage(value,proto.hashicorp.waypoint.RunnerConfigRequest.Open.deserializeBinaryFromReader);
      msg.setOpen(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerConfigRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerConfigRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfigRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOpen();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.RunnerConfigRequest.Open.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerConfigRequest.Open.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerConfigRequest.Open} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.toObject = function(includeInstance, msg) {
  var f, obj = {
    runner: (f = msg.getRunner()) && proto.hashicorp.waypoint.Runner.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest.Open}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerConfigRequest.Open;
  return proto.hashicorp.waypoint.RunnerConfigRequest.Open.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerConfigRequest.Open} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest.Open}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Runner;
      reader.readMessage(value,proto.hashicorp.waypoint.Runner.deserializeBinaryFromReader);
      msg.setRunner(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerConfigRequest.Open.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerConfigRequest.Open} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRunner();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Runner.serializeBinaryToWriter
    );
  }
};


/**
 * optional Runner runner = 1;
 * @return {?proto.hashicorp.waypoint.Runner}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.prototype.getRunner = function() {
  return /** @type{?proto.hashicorp.waypoint.Runner} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Runner, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Runner|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest.Open} returns this
*/
proto.hashicorp.waypoint.RunnerConfigRequest.Open.prototype.setRunner = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest.Open} returns this
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.prototype.clearRunner = function() {
  return this.setRunner(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.Open.prototype.hasRunner = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Open open = 1;
 * @return {?proto.hashicorp.waypoint.RunnerConfigRequest.Open}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.prototype.getOpen = function() {
  return /** @type{?proto.hashicorp.waypoint.RunnerConfigRequest.Open} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.RunnerConfigRequest.Open, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.RunnerConfigRequest.Open|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest} returns this
*/
proto.hashicorp.waypoint.RunnerConfigRequest.prototype.setOpen = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.RunnerConfigRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerConfigRequest} returns this
 */
proto.hashicorp.waypoint.RunnerConfigRequest.prototype.clearOpen = function() {
  return this.setOpen(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerConfigRequest.prototype.hasOpen = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerConfigResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerConfigResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerConfigResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfigResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    config: (f = msg.getConfig()) && proto.hashicorp.waypoint.RunnerConfig.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerConfigResponse}
 */
proto.hashicorp.waypoint.RunnerConfigResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerConfigResponse;
  return proto.hashicorp.waypoint.RunnerConfigResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerConfigResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerConfigResponse}
 */
proto.hashicorp.waypoint.RunnerConfigResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = new proto.hashicorp.waypoint.RunnerConfig;
      reader.readMessage(value,proto.hashicorp.waypoint.RunnerConfig.deserializeBinaryFromReader);
      msg.setConfig(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerConfigResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerConfigResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerConfigResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfigResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConfig();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.RunnerConfig.serializeBinaryToWriter
    );
  }
};


/**
 * optional RunnerConfig config = 2;
 * @return {?proto.hashicorp.waypoint.RunnerConfig}
 */
proto.hashicorp.waypoint.RunnerConfigResponse.prototype.getConfig = function() {
  return /** @type{?proto.hashicorp.waypoint.RunnerConfig} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.RunnerConfig, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.RunnerConfig|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerConfigResponse} returns this
*/
proto.hashicorp.waypoint.RunnerConfigResponse.prototype.setConfig = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerConfigResponse} returns this
 */
proto.hashicorp.waypoint.RunnerConfigResponse.prototype.clearConfig = function() {
  return this.setConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerConfigResponse.prototype.hasConfig = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.RunnerConfig.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    configVarsList: jspb.Message.toObjectList(msg.getConfigVarsList(),
    proto.hashicorp.waypoint.ConfigVar.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerConfig}
 */
proto.hashicorp.waypoint.RunnerConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerConfig;
  return proto.hashicorp.waypoint.RunnerConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerConfig}
 */
proto.hashicorp.waypoint.RunnerConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.ConfigVar;
      reader.readMessage(value,proto.hashicorp.waypoint.ConfigVar.deserializeBinaryFromReader);
      msg.addConfigVars(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConfigVarsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.ConfigVar.serializeBinaryToWriter
    );
  }
};


/**
 * repeated ConfigVar config_vars = 1;
 * @return {!Array<!proto.hashicorp.waypoint.ConfigVar>}
 */
proto.hashicorp.waypoint.RunnerConfig.prototype.getConfigVarsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.ConfigVar>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.ConfigVar, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.ConfigVar>} value
 * @return {!proto.hashicorp.waypoint.RunnerConfig} returns this
*/
proto.hashicorp.waypoint.RunnerConfig.prototype.setConfigVarsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.ConfigVar=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.ConfigVar}
 */
proto.hashicorp.waypoint.RunnerConfig.prototype.addConfigVars = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.ConfigVar, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.RunnerConfig} returns this
 */
proto.hashicorp.waypoint.RunnerConfig.prototype.clearConfigVarsList = function() {
  return this.setConfigVarsList([]);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_ = [[1,2,3,4,5]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.EventCase = {
  EVENT_NOT_SET: 0,
  REQUEST: 1,
  ACK: 2,
  COMPLETE: 3,
  ERROR: 4,
  TERMINAL: 5
};

/**
 * @return {proto.hashicorp.waypoint.RunnerJobStreamRequest.EventCase}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.RunnerJobStreamRequest.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    request: (f = msg.getRequest()) && proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.toObject(includeInstance, f),
    ack: (f = msg.getAck()) && proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.toObject(includeInstance, f),
    complete: (f = msg.getComplete()) && proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.toObject(includeInstance, f),
    error: (f = msg.getError()) && proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.toObject(includeInstance, f),
    terminal: (f = msg.getTerminal()) && proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerJobStreamRequest;
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Request;
      reader.readMessage(value,proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.deserializeBinaryFromReader);
      msg.setRequest(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack;
      reader.readMessage(value,proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.deserializeBinaryFromReader);
      msg.setAck(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete;
      reader.readMessage(value,proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.deserializeBinaryFromReader);
      msg.setComplete(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Error;
      reader.readMessage(value,proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    case 5:
      var value = new proto.hashicorp.waypoint.GetJobStreamResponse.Terminal;
      reader.readMessage(value,proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.deserializeBinaryFromReader);
      msg.setTerminal(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerJobStreamRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRequest();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.serializeBinaryToWriter
    );
  }
  f = message.getAck();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.serializeBinaryToWriter
    );
  }
  f = message.getComplete();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.serializeBinaryToWriter
    );
  }
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.serializeBinaryToWriter
    );
  }
  f = message.getTerminal();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.hashicorp.waypoint.GetJobStreamResponse.Terminal.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Request} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.toObject = function(includeInstance, msg) {
  var f, obj = {
    runnerId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Request}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Request;
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Request} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Request}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setRunnerId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Request} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRunnerId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string runner_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.prototype.getRunnerId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Request} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Request.prototype.setRunnerId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack;
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.toObject = function(includeInstance, msg) {
  var f, obj = {
    result: (f = msg.getResult()) && proto.hashicorp.waypoint.Job.Result.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete;
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Job.Result;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.Result.deserializeBinaryFromReader);
      msg.setResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getResult();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Job.Result.serializeBinaryToWriter
    );
  }
};


/**
 * optional Job.Result result = 1;
 * @return {?proto.hashicorp.waypoint.Job.Result}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.prototype.getResult = function() {
  return /** @type{?proto.hashicorp.waypoint.Job.Result} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job.Result, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Job.Result|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.prototype.setResult = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.prototype.clearResult = function() {
  return this.setResult(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete.prototype.hasResult = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Error} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.toObject = function(includeInstance, msg) {
  var f, obj = {
    error: (f = msg.getError()) && google_rpc_status_pb.Status.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Error}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerJobStreamRequest.Error;
  return proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Error} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Error}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_rpc_status_pb.Status;
      reader.readMessage(value,google_rpc_status_pb.Status.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Error} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_rpc_status_pb.Status.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.rpc.Status error = 1;
 * @return {?proto.google.rpc.Status}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.prototype.getError = function() {
  return /** @type{?proto.google.rpc.Status} */ (
    jspb.Message.getWrapperField(this, google_rpc_status_pb.Status, 1));
};


/**
 * @param {?proto.google.rpc.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Error} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.prototype.setError = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest.Error} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.Error.prototype.hasError = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Request request = 1;
 * @return {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Request}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.getRequest = function() {
  return /** @type{?proto.hashicorp.waypoint.RunnerJobStreamRequest.Request} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.RunnerJobStreamRequest.Request, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Request|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.setRequest = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.clearRequest = function() {
  return this.setRequest(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.hasRequest = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Ack ack = 2;
 * @return {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.getAck = function() {
  return /** @type{?proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Ack|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.setAck = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.clearAck = function() {
  return this.setAck(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.hasAck = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Complete complete = 3;
 * @return {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.getComplete = function() {
  return /** @type{?proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Complete|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.setComplete = function(value) {
  return jspb.Message.setOneofWrapperField(this, 3, proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.clearComplete = function() {
  return this.setComplete(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.hasComplete = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Error error = 4;
 * @return {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Error}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.getError = function() {
  return /** @type{?proto.hashicorp.waypoint.RunnerJobStreamRequest.Error} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.RunnerJobStreamRequest.Error, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.RunnerJobStreamRequest.Error|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.setError = function(value) {
  return jspb.Message.setOneofWrapperField(this, 4, proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.hasError = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional GetJobStreamResponse.Terminal terminal = 5;
 * @return {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.getTerminal = function() {
  return /** @type{?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.GetJobStreamResponse.Terminal, 5));
};


/**
 * @param {?proto.hashicorp.waypoint.GetJobStreamResponse.Terminal|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.setTerminal = function(value) {
  return jspb.Message.setOneofWrapperField(this, 5, proto.hashicorp.waypoint.RunnerJobStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamRequest} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.clearTerminal = function() {
  return this.setTerminal(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamRequest.prototype.hasTerminal = function() {
  return jspb.Message.getField(this, 5) != null;
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.oneofGroups_ = [[1]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.EventCase = {
  EVENT_NOT_SET: 0,
  ASSIGNMENT: 1
};

/**
 * @return {proto.hashicorp.waypoint.RunnerJobStreamResponse.EventCase}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.RunnerJobStreamResponse.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.RunnerJobStreamResponse.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerJobStreamResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    assignment: (f = msg.getAssignment()) && proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerJobStreamResponse;
  return proto.hashicorp.waypoint.RunnerJobStreamResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment;
      reader.readMessage(value,proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.deserializeBinaryFromReader);
      msg.setAssignment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerJobStreamResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAssignment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.toObject = function(includeInstance, msg) {
  var f, obj = {
    job: (f = msg.getJob()) && proto.hashicorp.waypoint.Job.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment;
  return proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Job;
      reader.readMessage(value,proto.hashicorp.waypoint.Job.deserializeBinaryFromReader);
      msg.setJob(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getJob();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Job.serializeBinaryToWriter
    );
  }
};


/**
 * optional Job job = 1;
 * @return {?proto.hashicorp.waypoint.Job}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.prototype.getJob = function() {
  return /** @type{?proto.hashicorp.waypoint.Job} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Job, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Job|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.prototype.setJob = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.prototype.clearJob = function() {
  return this.setJob(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment.prototype.hasJob = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional JobAssignment assignment = 1;
 * @return {?proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.prototype.getAssignment = function() {
  return /** @type{?proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.RunnerJobStreamResponse.JobAssignment|undefined} value
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse} returns this
*/
proto.hashicorp.waypoint.RunnerJobStreamResponse.prototype.setAssignment = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.RunnerJobStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.RunnerJobStreamResponse} returns this
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.prototype.clearAssignment = function() {
  return this.setAssignment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerJobStreamResponse.prototype.hasAssignment = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest;
  return proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    serverAddr: jspb.Message.getFieldWithDefault(msg, 1, ""),
    serverInsecure: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse;
  return proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setServerAddr(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setServerInsecure(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getServerAddr();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getServerInsecure();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * optional string server_addr = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.prototype.getServerAddr = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse} returns this
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.prototype.setServerAddr = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bool server_insecure = 2;
 * @return {boolean}
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.prototype.getServerInsecure = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse} returns this
 */
proto.hashicorp.waypoint.RunnerGetDeploymentConfigResponse.prototype.setServerInsecure = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetRunnerRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetRunnerRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetRunnerRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetRunnerRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    runnerId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetRunnerRequest}
 */
proto.hashicorp.waypoint.GetRunnerRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetRunnerRequest;
  return proto.hashicorp.waypoint.GetRunnerRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetRunnerRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetRunnerRequest}
 */
proto.hashicorp.waypoint.GetRunnerRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setRunnerId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetRunnerRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetRunnerRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetRunnerRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetRunnerRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRunnerId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string runner_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetRunnerRequest.prototype.getRunnerId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetRunnerRequest} returns this
 */
proto.hashicorp.waypoint.GetRunnerRequest.prototype.setRunnerId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.SetServerConfigRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.SetServerConfigRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.SetServerConfigRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.SetServerConfigRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    config: (f = msg.getConfig()) && proto.hashicorp.waypoint.ServerConfig.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.SetServerConfigRequest}
 */
proto.hashicorp.waypoint.SetServerConfigRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.SetServerConfigRequest;
  return proto.hashicorp.waypoint.SetServerConfigRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.SetServerConfigRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.SetServerConfigRequest}
 */
proto.hashicorp.waypoint.SetServerConfigRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.ServerConfig;
      reader.readMessage(value,proto.hashicorp.waypoint.ServerConfig.deserializeBinaryFromReader);
      msg.setConfig(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.SetServerConfigRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.SetServerConfigRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.SetServerConfigRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.SetServerConfigRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConfig();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.ServerConfig.serializeBinaryToWriter
    );
  }
};


/**
 * optional ServerConfig config = 1;
 * @return {?proto.hashicorp.waypoint.ServerConfig}
 */
proto.hashicorp.waypoint.SetServerConfigRequest.prototype.getConfig = function() {
  return /** @type{?proto.hashicorp.waypoint.ServerConfig} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ServerConfig, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.ServerConfig|undefined} value
 * @return {!proto.hashicorp.waypoint.SetServerConfigRequest} returns this
*/
proto.hashicorp.waypoint.SetServerConfigRequest.prototype.setConfig = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.SetServerConfigRequest} returns this
 */
proto.hashicorp.waypoint.SetServerConfigRequest.prototype.clearConfig = function() {
  return this.setConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.SetServerConfigRequest.prototype.hasConfig = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ServerConfig.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ServerConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ServerConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ServerConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ServerConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    advertiseAddrsList: jspb.Message.toObjectList(msg.getAdvertiseAddrsList(),
    proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ServerConfig}
 */
proto.hashicorp.waypoint.ServerConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ServerConfig;
  return proto.hashicorp.waypoint.ServerConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ServerConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ServerConfig}
 */
proto.hashicorp.waypoint.ServerConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr;
      reader.readMessage(value,proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.deserializeBinaryFromReader);
      msg.addAdvertiseAddrs(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ServerConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ServerConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ServerConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ServerConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAdvertiseAddrsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.toObject = function(includeInstance, msg) {
  var f, obj = {
    addr: jspb.Message.getFieldWithDefault(msg, 1, ""),
    insecure: jspb.Message.getBooleanFieldWithDefault(msg, 2, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr}
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr;
  return proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr}
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAddr(value);
      break;
    case 2:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setInsecure(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAddr();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getInsecure();
  if (f) {
    writer.writeBool(
      2,
      f
    );
  }
};


/**
 * optional string addr = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.prototype.getAddr = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr} returns this
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.prototype.setAddr = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bool insecure = 2;
 * @return {boolean}
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.prototype.getInsecure = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 2, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr} returns this
 */
proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr.prototype.setInsecure = function(value) {
  return jspb.Message.setProto3BooleanField(this, 2, value);
};


/**
 * repeated AdvertiseAddr advertise_addrs = 1;
 * @return {!Array<!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr>}
 */
proto.hashicorp.waypoint.ServerConfig.prototype.getAdvertiseAddrsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr>} value
 * @return {!proto.hashicorp.waypoint.ServerConfig} returns this
*/
proto.hashicorp.waypoint.ServerConfig.prototype.setAdvertiseAddrsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr}
 */
proto.hashicorp.waypoint.ServerConfig.prototype.addAdvertiseAddrs = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.ServerConfig.AdvertiseAddr, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ServerConfig} returns this
 */
proto.hashicorp.waypoint.ServerConfig.prototype.clearAdvertiseAddrsList = function() {
  return this.setAdvertiseAddrsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertBuildRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertBuildRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertBuildRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertBuildRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    build: (f = msg.getBuild()) && proto.hashicorp.waypoint.Build.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertBuildRequest}
 */
proto.hashicorp.waypoint.UpsertBuildRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertBuildRequest;
  return proto.hashicorp.waypoint.UpsertBuildRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertBuildRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertBuildRequest}
 */
proto.hashicorp.waypoint.UpsertBuildRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Build;
      reader.readMessage(value,proto.hashicorp.waypoint.Build.deserializeBinaryFromReader);
      msg.setBuild(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertBuildRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertBuildRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertBuildRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertBuildRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBuild();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Build.serializeBinaryToWriter
    );
  }
};


/**
 * optional Build build = 1;
 * @return {?proto.hashicorp.waypoint.Build}
 */
proto.hashicorp.waypoint.UpsertBuildRequest.prototype.getBuild = function() {
  return /** @type{?proto.hashicorp.waypoint.Build} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Build, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Build|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertBuildRequest} returns this
*/
proto.hashicorp.waypoint.UpsertBuildRequest.prototype.setBuild = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertBuildRequest} returns this
 */
proto.hashicorp.waypoint.UpsertBuildRequest.prototype.clearBuild = function() {
  return this.setBuild(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertBuildRequest.prototype.hasBuild = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertBuildResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertBuildResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertBuildResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertBuildResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    build: (f = msg.getBuild()) && proto.hashicorp.waypoint.Build.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertBuildResponse}
 */
proto.hashicorp.waypoint.UpsertBuildResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertBuildResponse;
  return proto.hashicorp.waypoint.UpsertBuildResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertBuildResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertBuildResponse}
 */
proto.hashicorp.waypoint.UpsertBuildResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Build;
      reader.readMessage(value,proto.hashicorp.waypoint.Build.deserializeBinaryFromReader);
      msg.setBuild(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertBuildResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertBuildResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertBuildResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertBuildResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBuild();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Build.serializeBinaryToWriter
    );
  }
};


/**
 * optional Build build = 1;
 * @return {?proto.hashicorp.waypoint.Build}
 */
proto.hashicorp.waypoint.UpsertBuildResponse.prototype.getBuild = function() {
  return /** @type{?proto.hashicorp.waypoint.Build} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Build, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Build|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertBuildResponse} returns this
*/
proto.hashicorp.waypoint.UpsertBuildResponse.prototype.setBuild = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertBuildResponse} returns this
 */
proto.hashicorp.waypoint.UpsertBuildResponse.prototype.clearBuild = function() {
  return this.setBuild(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertBuildResponse.prototype.hasBuild = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ListBuildsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ListBuildsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListBuildsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ListBuildsRequest}
 */
proto.hashicorp.waypoint.ListBuildsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ListBuildsRequest;
  return proto.hashicorp.waypoint.ListBuildsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ListBuildsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ListBuildsRequest}
 */
proto.hashicorp.waypoint.ListBuildsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ListBuildsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ListBuildsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListBuildsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
};


/**
 * optional Ref.Application application = 1;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.ListBuildsRequest} returns this
*/
proto.hashicorp.waypoint.ListBuildsRequest.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListBuildsRequest} returns this
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Ref.Workspace workspace = 2;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.ListBuildsRequest} returns this
*/
proto.hashicorp.waypoint.ListBuildsRequest.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListBuildsRequest} returns this
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListBuildsRequest.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ListBuildsResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ListBuildsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ListBuildsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ListBuildsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListBuildsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    buildsList: jspb.Message.toObjectList(msg.getBuildsList(),
    proto.hashicorp.waypoint.Build.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ListBuildsResponse}
 */
proto.hashicorp.waypoint.ListBuildsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ListBuildsResponse;
  return proto.hashicorp.waypoint.ListBuildsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ListBuildsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ListBuildsResponse}
 */
proto.hashicorp.waypoint.ListBuildsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Build;
      reader.readMessage(value,proto.hashicorp.waypoint.Build.deserializeBinaryFromReader);
      msg.addBuilds(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ListBuildsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ListBuildsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ListBuildsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListBuildsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBuildsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.Build.serializeBinaryToWriter
    );
  }
};


/**
 * repeated Build builds = 1;
 * @return {!Array<!proto.hashicorp.waypoint.Build>}
 */
proto.hashicorp.waypoint.ListBuildsResponse.prototype.getBuildsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.Build>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.Build, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.Build>} value
 * @return {!proto.hashicorp.waypoint.ListBuildsResponse} returns this
*/
proto.hashicorp.waypoint.ListBuildsResponse.prototype.setBuildsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.Build=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.Build}
 */
proto.hashicorp.waypoint.ListBuildsResponse.prototype.addBuilds = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.Build, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ListBuildsResponse} returns this
 */
proto.hashicorp.waypoint.ListBuildsResponse.prototype.clearBuildsList = function() {
  return this.setBuildsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetLatestBuildRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetLatestBuildRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetLatestBuildRequest}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetLatestBuildRequest;
  return proto.hashicorp.waypoint.GetLatestBuildRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetLatestBuildRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetLatestBuildRequest}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetLatestBuildRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetLatestBuildRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
};


/**
 * optional Ref.Application application = 1;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.GetLatestBuildRequest} returns this
*/
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetLatestBuildRequest} returns this
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Ref.Workspace workspace = 2;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.GetLatestBuildRequest} returns this
*/
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetLatestBuildRequest} returns this
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetLatestBuildRequest.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Build.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Build.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Build} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Build.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f),
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    status: (f = msg.getStatus()) && proto.hashicorp.waypoint.Status.toObject(includeInstance, f),
    component: (f = msg.getComponent()) && proto.hashicorp.waypoint.Component.toObject(includeInstance, f),
    artifact: (f = msg.getArtifact()) && proto.hashicorp.waypoint.Artifact.toObject(includeInstance, f),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Build}
 */
proto.hashicorp.waypoint.Build.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Build;
  return proto.hashicorp.waypoint.Build.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Build} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Build}
 */
proto.hashicorp.waypoint.Build.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 6:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 7:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Status;
      reader.readMessage(value,proto.hashicorp.waypoint.Status.deserializeBinaryFromReader);
      msg.setStatus(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Component;
      reader.readMessage(value,proto.hashicorp.waypoint.Component.deserializeBinaryFromReader);
      msg.setComponent(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Artifact;
      reader.readMessage(value,proto.hashicorp.waypoint.Artifact.deserializeBinaryFromReader);
      msg.setArtifact(value);
      break;
    case 5:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Build.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Build.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Build} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Build.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStatus();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Status.serializeBinaryToWriter
    );
  }
  f = message.getComponent();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Component.serializeBinaryToWriter
    );
  }
  f = message.getArtifact();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Artifact.serializeBinaryToWriter
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(5, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional Ref.Application application = 6;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.Build.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 6));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.Build} returns this
*/
proto.hashicorp.waypoint.Build.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Build} returns this
 */
proto.hashicorp.waypoint.Build.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Build.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional Ref.Workspace workspace = 7;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.Build.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 7));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.Build} returns this
*/
proto.hashicorp.waypoint.Build.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Build} returns this
 */
proto.hashicorp.waypoint.Build.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Build.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Build.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Build} returns this
 */
proto.hashicorp.waypoint.Build.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Status status = 2;
 * @return {?proto.hashicorp.waypoint.Status}
 */
proto.hashicorp.waypoint.Build.prototype.getStatus = function() {
  return /** @type{?proto.hashicorp.waypoint.Status} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Status, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.Build} returns this
*/
proto.hashicorp.waypoint.Build.prototype.setStatus = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Build} returns this
 */
proto.hashicorp.waypoint.Build.prototype.clearStatus = function() {
  return this.setStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Build.prototype.hasStatus = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Component component = 3;
 * @return {?proto.hashicorp.waypoint.Component}
 */
proto.hashicorp.waypoint.Build.prototype.getComponent = function() {
  return /** @type{?proto.hashicorp.waypoint.Component} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Component, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Component|undefined} value
 * @return {!proto.hashicorp.waypoint.Build} returns this
*/
proto.hashicorp.waypoint.Build.prototype.setComponent = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Build} returns this
 */
proto.hashicorp.waypoint.Build.prototype.clearComponent = function() {
  return this.setComponent(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Build.prototype.hasComponent = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Artifact artifact = 4;
 * @return {?proto.hashicorp.waypoint.Artifact}
 */
proto.hashicorp.waypoint.Build.prototype.getArtifact = function() {
  return /** @type{?proto.hashicorp.waypoint.Artifact} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Artifact, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Artifact|undefined} value
 * @return {!proto.hashicorp.waypoint.Build} returns this
*/
proto.hashicorp.waypoint.Build.prototype.setArtifact = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Build} returns this
 */
proto.hashicorp.waypoint.Build.prototype.clearArtifact = function() {
  return this.setArtifact(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Build.prototype.hasArtifact = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * map<string, string> labels = 5;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.hashicorp.waypoint.Build.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 5, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.hashicorp.waypoint.Build} returns this
 */
proto.hashicorp.waypoint.Build.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Artifact.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Artifact.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Artifact} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Artifact.toObject = function(includeInstance, msg) {
  var f, obj = {
    artifact: (f = msg.getArtifact()) && google_protobuf_any_pb.Any.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Artifact}
 */
proto.hashicorp.waypoint.Artifact.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Artifact;
  return proto.hashicorp.waypoint.Artifact.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Artifact} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Artifact}
 */
proto.hashicorp.waypoint.Artifact.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setArtifact(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Artifact.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Artifact.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Artifact} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Artifact.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getArtifact();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.protobuf.Any artifact = 1;
 * @return {?proto.google.protobuf.Any}
 */
proto.hashicorp.waypoint.Artifact.prototype.getArtifact = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 1));
};


/**
 * @param {?proto.google.protobuf.Any|undefined} value
 * @return {!proto.hashicorp.waypoint.Artifact} returns this
*/
proto.hashicorp.waypoint.Artifact.prototype.setArtifact = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Artifact} returns this
 */
proto.hashicorp.waypoint.Artifact.prototype.clearArtifact = function() {
  return this.setArtifact(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Artifact.prototype.hasArtifact = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertPushedArtifactRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertPushedArtifactRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    artifact: (f = msg.getArtifact()) && proto.hashicorp.waypoint.PushedArtifact.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactRequest}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertPushedArtifactRequest;
  return proto.hashicorp.waypoint.UpsertPushedArtifactRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertPushedArtifactRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactRequest}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.PushedArtifact;
      reader.readMessage(value,proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader);
      msg.setArtifact(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertPushedArtifactRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertPushedArtifactRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getArtifact();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter
    );
  }
};


/**
 * optional PushedArtifact artifact = 1;
 * @return {?proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.prototype.getArtifact = function() {
  return /** @type{?proto.hashicorp.waypoint.PushedArtifact} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.PushedArtifact, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.PushedArtifact|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactRequest} returns this
*/
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.prototype.setArtifact = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactRequest} returns this
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.prototype.clearArtifact = function() {
  return this.setArtifact(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactRequest.prototype.hasArtifact = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertPushedArtifactResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertPushedArtifactResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    artifact: (f = msg.getArtifact()) && proto.hashicorp.waypoint.PushedArtifact.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactResponse}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertPushedArtifactResponse;
  return proto.hashicorp.waypoint.UpsertPushedArtifactResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertPushedArtifactResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactResponse}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.PushedArtifact;
      reader.readMessage(value,proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader);
      msg.setArtifact(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertPushedArtifactResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertPushedArtifactResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getArtifact();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter
    );
  }
};


/**
 * optional PushedArtifact artifact = 1;
 * @return {?proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.prototype.getArtifact = function() {
  return /** @type{?proto.hashicorp.waypoint.PushedArtifact} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.PushedArtifact, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.PushedArtifact|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactResponse} returns this
*/
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.prototype.setArtifact = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertPushedArtifactResponse} returns this
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.prototype.clearArtifact = function() {
  return this.setArtifact(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertPushedArtifactResponse.prototype.hasArtifact = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetLatestPushedArtifactRequest;
  return proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
};


/**
 * optional Ref.Application application = 1;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest} returns this
*/
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest} returns this
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Ref.Workspace workspace = 2;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest} returns this
*/
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.GetLatestPushedArtifactRequest} returns this
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.GetLatestPushedArtifactRequest.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ListPushedArtifactsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f),
    statusList: jspb.Message.toObjectList(msg.getStatusList(),
    proto.hashicorp.waypoint.StatusFilter.toObject, includeInstance),
    order: (f = msg.getOrder()) && proto.hashicorp.waypoint.OperationOrder.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ListPushedArtifactsRequest;
  return proto.hashicorp.waypoint.ListPushedArtifactsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 3:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    case 1:
      var value = new proto.hashicorp.waypoint.StatusFilter;
      reader.readMessage(value,proto.hashicorp.waypoint.StatusFilter.deserializeBinaryFromReader);
      msg.addStatus(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.OperationOrder;
      reader.readMessage(value,proto.hashicorp.waypoint.OperationOrder.deserializeBinaryFromReader);
      msg.setOrder(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ListPushedArtifactsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
  f = message.getStatusList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.StatusFilter.serializeBinaryToWriter
    );
  }
  f = message.getOrder();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.OperationOrder.serializeBinaryToWriter
    );
  }
};


/**
 * optional Ref.Application application = 3;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
*/
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Ref.Workspace workspace = 4;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
*/
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * repeated StatusFilter status = 1;
 * @return {!Array<!proto.hashicorp.waypoint.StatusFilter>}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.getStatusList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.StatusFilter>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.StatusFilter, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.StatusFilter>} value
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
*/
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.setStatusList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.StatusFilter=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.StatusFilter}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.addStatus = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.StatusFilter, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.clearStatusList = function() {
  return this.setStatusList([]);
};


/**
 * optional OperationOrder order = 2;
 * @return {?proto.hashicorp.waypoint.OperationOrder}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.getOrder = function() {
  return /** @type{?proto.hashicorp.waypoint.OperationOrder} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.OperationOrder, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.OperationOrder|undefined} value
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
*/
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.setOrder = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsRequest} returns this
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.clearOrder = function() {
  return this.setOrder(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListPushedArtifactsRequest.prototype.hasOrder = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ListPushedArtifactsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ListPushedArtifactsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    artifactsList: jspb.Message.toObjectList(msg.getArtifactsList(),
    proto.hashicorp.waypoint.PushedArtifact.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsResponse}
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ListPushedArtifactsResponse;
  return proto.hashicorp.waypoint.ListPushedArtifactsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ListPushedArtifactsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsResponse}
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.PushedArtifact;
      reader.readMessage(value,proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader);
      msg.addArtifacts(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ListPushedArtifactsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ListPushedArtifactsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getArtifactsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter
    );
  }
};


/**
 * repeated PushedArtifact artifacts = 1;
 * @return {!Array<!proto.hashicorp.waypoint.PushedArtifact>}
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.prototype.getArtifactsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.PushedArtifact>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.PushedArtifact, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.PushedArtifact>} value
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsResponse} returns this
*/
proto.hashicorp.waypoint.ListPushedArtifactsResponse.prototype.setArtifactsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.PushedArtifact=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.prototype.addArtifacts = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.PushedArtifact, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ListPushedArtifactsResponse} returns this
 */
proto.hashicorp.waypoint.ListPushedArtifactsResponse.prototype.clearArtifactsList = function() {
  return this.setArtifactsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.PushedArtifact.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.PushedArtifact} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.PushedArtifact.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f),
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    status: (f = msg.getStatus()) && proto.hashicorp.waypoint.Status.toObject(includeInstance, f),
    component: (f = msg.getComponent()) && proto.hashicorp.waypoint.Component.toObject(includeInstance, f),
    artifact: (f = msg.getArtifact()) && proto.hashicorp.waypoint.Artifact.toObject(includeInstance, f),
    buildId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.PushedArtifact.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.PushedArtifact;
  return proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.PushedArtifact} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.PushedArtifact}
 */
proto.hashicorp.waypoint.PushedArtifact.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 7:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 8:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Status;
      reader.readMessage(value,proto.hashicorp.waypoint.Status.deserializeBinaryFromReader);
      msg.setStatus(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Component;
      reader.readMessage(value,proto.hashicorp.waypoint.Component.deserializeBinaryFromReader);
      msg.setComponent(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Artifact;
      reader.readMessage(value,proto.hashicorp.waypoint.Artifact.deserializeBinaryFromReader);
      msg.setArtifact(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setBuildId(value);
      break;
    case 6:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.PushedArtifact} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.PushedArtifact.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStatus();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Status.serializeBinaryToWriter
    );
  }
  f = message.getComponent();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Component.serializeBinaryToWriter
    );
  }
  f = message.getArtifact();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Artifact.serializeBinaryToWriter
    );
  }
  f = message.getBuildId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(6, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional Ref.Application application = 7;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 7));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
*/
proto.hashicorp.waypoint.PushedArtifact.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional Ref.Workspace workspace = 8;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 8));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
*/
proto.hashicorp.waypoint.PushedArtifact.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Status status = 2;
 * @return {?proto.hashicorp.waypoint.Status}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getStatus = function() {
  return /** @type{?proto.hashicorp.waypoint.Status} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Status, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
*/
proto.hashicorp.waypoint.PushedArtifact.prototype.setStatus = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.clearStatus = function() {
  return this.setStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.hasStatus = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Component component = 3;
 * @return {?proto.hashicorp.waypoint.Component}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getComponent = function() {
  return /** @type{?proto.hashicorp.waypoint.Component} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Component, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Component|undefined} value
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
*/
proto.hashicorp.waypoint.PushedArtifact.prototype.setComponent = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.clearComponent = function() {
  return this.setComponent(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.hasComponent = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Artifact artifact = 4;
 * @return {?proto.hashicorp.waypoint.Artifact}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getArtifact = function() {
  return /** @type{?proto.hashicorp.waypoint.Artifact} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Artifact, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Artifact|undefined} value
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
*/
proto.hashicorp.waypoint.PushedArtifact.prototype.setArtifact = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.clearArtifact = function() {
  return this.setArtifact(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.hasArtifact = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional string build_id = 5;
 * @return {string}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getBuildId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.setBuildId = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * map<string, string> labels = 6;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 6, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.hashicorp.waypoint.PushedArtifact} returns this
 */
proto.hashicorp.waypoint.PushedArtifact.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetDeploymentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetDeploymentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetDeploymentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetDeploymentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetDeploymentRequest}
 */
proto.hashicorp.waypoint.GetDeploymentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetDeploymentRequest;
  return proto.hashicorp.waypoint.GetDeploymentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetDeploymentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetDeploymentRequest}
 */
proto.hashicorp.waypoint.GetDeploymentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetDeploymentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetDeploymentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetDeploymentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetDeploymentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetDeploymentRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetDeploymentRequest} returns this
 */
proto.hashicorp.waypoint.GetDeploymentRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertDeploymentRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertDeploymentRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deployment: (f = msg.getDeployment()) && proto.hashicorp.waypoint.Deployment.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentRequest}
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertDeploymentRequest;
  return proto.hashicorp.waypoint.UpsertDeploymentRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertDeploymentRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentRequest}
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Deployment;
      reader.readMessage(value,proto.hashicorp.waypoint.Deployment.deserializeBinaryFromReader);
      msg.setDeployment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertDeploymentRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertDeploymentRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeployment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Deployment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Deployment deployment = 1;
 * @return {?proto.hashicorp.waypoint.Deployment}
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.prototype.getDeployment = function() {
  return /** @type{?proto.hashicorp.waypoint.Deployment} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Deployment, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Deployment|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentRequest} returns this
*/
proto.hashicorp.waypoint.UpsertDeploymentRequest.prototype.setDeployment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentRequest} returns this
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.prototype.clearDeployment = function() {
  return this.setDeployment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertDeploymentRequest.prototype.hasDeployment = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertDeploymentResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertDeploymentResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    deployment: (f = msg.getDeployment()) && proto.hashicorp.waypoint.Deployment.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentResponse}
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertDeploymentResponse;
  return proto.hashicorp.waypoint.UpsertDeploymentResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertDeploymentResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentResponse}
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Deployment;
      reader.readMessage(value,proto.hashicorp.waypoint.Deployment.deserializeBinaryFromReader);
      msg.setDeployment(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertDeploymentResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertDeploymentResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeployment();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Deployment.serializeBinaryToWriter
    );
  }
};


/**
 * optional Deployment deployment = 1;
 * @return {?proto.hashicorp.waypoint.Deployment}
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.prototype.getDeployment = function() {
  return /** @type{?proto.hashicorp.waypoint.Deployment} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Deployment, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Deployment|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentResponse} returns this
*/
proto.hashicorp.waypoint.UpsertDeploymentResponse.prototype.setDeployment = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertDeploymentResponse} returns this
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.prototype.clearDeployment = function() {
  return this.setDeployment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertDeploymentResponse.prototype.hasDeployment = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ListDeploymentsRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ListDeploymentsRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f),
    statusList: jspb.Message.toObjectList(msg.getStatusList(),
    proto.hashicorp.waypoint.StatusFilter.toObject, includeInstance),
    order: (f = msg.getOrder()) && proto.hashicorp.waypoint.OperationOrder.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ListDeploymentsRequest;
  return proto.hashicorp.waypoint.ListDeploymentsRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ListDeploymentsRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 3:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    case 1:
      var value = new proto.hashicorp.waypoint.StatusFilter;
      reader.readMessage(value,proto.hashicorp.waypoint.StatusFilter.deserializeBinaryFromReader);
      msg.addStatus(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.OperationOrder;
      reader.readMessage(value,proto.hashicorp.waypoint.OperationOrder.deserializeBinaryFromReader);
      msg.setOrder(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ListDeploymentsRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ListDeploymentsRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
  f = message.getStatusList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.StatusFilter.serializeBinaryToWriter
    );
  }
  f = message.getOrder();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.OperationOrder.serializeBinaryToWriter
    );
  }
};


/**
 * optional Ref.Application application = 3;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
*/
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Ref.Workspace workspace = 4;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
*/
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * repeated StatusFilter status = 1;
 * @return {!Array<!proto.hashicorp.waypoint.StatusFilter>}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.getStatusList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.StatusFilter>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.StatusFilter, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.StatusFilter>} value
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
*/
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.setStatusList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.StatusFilter=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.StatusFilter}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.addStatus = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.StatusFilter, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.clearStatusList = function() {
  return this.setStatusList([]);
};


/**
 * optional OperationOrder order = 2;
 * @return {?proto.hashicorp.waypoint.OperationOrder}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.getOrder = function() {
  return /** @type{?proto.hashicorp.waypoint.OperationOrder} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.OperationOrder, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.OperationOrder|undefined} value
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
*/
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.setOrder = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsRequest} returns this
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.clearOrder = function() {
  return this.setOrder(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ListDeploymentsRequest.prototype.hasOrder = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ListDeploymentsResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ListDeploymentsResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentsList: jspb.Message.toObjectList(msg.getDeploymentsList(),
    proto.hashicorp.waypoint.Deployment.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsResponse}
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ListDeploymentsResponse;
  return proto.hashicorp.waypoint.ListDeploymentsResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ListDeploymentsResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsResponse}
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Deployment;
      reader.readMessage(value,proto.hashicorp.waypoint.Deployment.deserializeBinaryFromReader);
      msg.addDeployments(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ListDeploymentsResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ListDeploymentsResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.Deployment.serializeBinaryToWriter
    );
  }
};


/**
 * repeated Deployment deployments = 1;
 * @return {!Array<!proto.hashicorp.waypoint.Deployment>}
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.prototype.getDeploymentsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.Deployment>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.Deployment, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.Deployment>} value
 * @return {!proto.hashicorp.waypoint.ListDeploymentsResponse} returns this
*/
proto.hashicorp.waypoint.ListDeploymentsResponse.prototype.setDeploymentsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.Deployment=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.Deployment}
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.prototype.addDeployments = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.Deployment, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ListDeploymentsResponse} returns this
 */
proto.hashicorp.waypoint.ListDeploymentsResponse.prototype.clearDeploymentsList = function() {
  return this.setDeploymentsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Deployment.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Deployment.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Deployment} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Deployment.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f),
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    state: jspb.Message.getFieldWithDefault(msg, 2, 0),
    status: (f = msg.getStatus()) && proto.hashicorp.waypoint.Status.toObject(includeInstance, f),
    component: (f = msg.getComponent()) && proto.hashicorp.waypoint.Component.toObject(includeInstance, f),
    artifactId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    deployment: (f = msg.getDeployment()) && google_protobuf_any_pb.Any.toObject(includeInstance, f),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Deployment}
 */
proto.hashicorp.waypoint.Deployment.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Deployment;
  return proto.hashicorp.waypoint.Deployment.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Deployment} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Deployment}
 */
proto.hashicorp.waypoint.Deployment.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 8:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 9:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {!proto.hashicorp.waypoint.Deployment.State} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Status;
      reader.readMessage(value,proto.hashicorp.waypoint.Status.deserializeBinaryFromReader);
      msg.setStatus(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Component;
      reader.readMessage(value,proto.hashicorp.waypoint.Component.deserializeBinaryFromReader);
      msg.setComponent(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setArtifactId(value);
      break;
    case 6:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setDeployment(value);
      break;
    case 7:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Deployment.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Deployment.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Deployment} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Deployment.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      9,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getStatus();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Status.serializeBinaryToWriter
    );
  }
  f = message.getComponent();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Component.serializeBinaryToWriter
    );
  }
  f = message.getArtifactId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getDeployment();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(7, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * @enum {number}
 */
proto.hashicorp.waypoint.Deployment.State = {
  UNKNOWN: 0,
  PENDING: 1,
  DEPLOY: 3,
  DESTROY: 4
};

/**
 * optional Ref.Application application = 8;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.Deployment.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 8));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
*/
proto.hashicorp.waypoint.Deployment.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Deployment.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional Ref.Workspace workspace = 9;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.Deployment.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 9));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
*/
proto.hashicorp.waypoint.Deployment.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 9, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Deployment.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 9) != null;
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Deployment.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional State state = 2;
 * @return {!proto.hashicorp.waypoint.Deployment.State}
 */
proto.hashicorp.waypoint.Deployment.prototype.getState = function() {
  return /** @type {!proto.hashicorp.waypoint.Deployment.State} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.Deployment.State} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.setState = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional Status status = 3;
 * @return {?proto.hashicorp.waypoint.Status}
 */
proto.hashicorp.waypoint.Deployment.prototype.getStatus = function() {
  return /** @type{?proto.hashicorp.waypoint.Status} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Status, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
*/
proto.hashicorp.waypoint.Deployment.prototype.setStatus = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.clearStatus = function() {
  return this.setStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Deployment.prototype.hasStatus = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Component component = 4;
 * @return {?proto.hashicorp.waypoint.Component}
 */
proto.hashicorp.waypoint.Deployment.prototype.getComponent = function() {
  return /** @type{?proto.hashicorp.waypoint.Component} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Component, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Component|undefined} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
*/
proto.hashicorp.waypoint.Deployment.prototype.setComponent = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.clearComponent = function() {
  return this.setComponent(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Deployment.prototype.hasComponent = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional string artifact_id = 5;
 * @return {string}
 */
proto.hashicorp.waypoint.Deployment.prototype.getArtifactId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.setArtifactId = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional google.protobuf.Any deployment = 6;
 * @return {?proto.google.protobuf.Any}
 */
proto.hashicorp.waypoint.Deployment.prototype.getDeployment = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 6));
};


/**
 * @param {?proto.google.protobuf.Any|undefined} value
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
*/
proto.hashicorp.waypoint.Deployment.prototype.setDeployment = function(value) {
  return jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.clearDeployment = function() {
  return this.setDeployment(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Deployment.prototype.hasDeployment = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * map<string, string> labels = 7;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.hashicorp.waypoint.Deployment.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 7, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.hashicorp.waypoint.Deployment} returns this
 */
proto.hashicorp.waypoint.Deployment.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertReleaseRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertReleaseRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    release: (f = msg.getRelease()) && proto.hashicorp.waypoint.Release.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertReleaseRequest}
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertReleaseRequest;
  return proto.hashicorp.waypoint.UpsertReleaseRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertReleaseRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertReleaseRequest}
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Release;
      reader.readMessage(value,proto.hashicorp.waypoint.Release.deserializeBinaryFromReader);
      msg.setRelease(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertReleaseRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertReleaseRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRelease();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Release.serializeBinaryToWriter
    );
  }
};


/**
 * optional Release release = 1;
 * @return {?proto.hashicorp.waypoint.Release}
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.prototype.getRelease = function() {
  return /** @type{?proto.hashicorp.waypoint.Release} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Release, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Release|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertReleaseRequest} returns this
*/
proto.hashicorp.waypoint.UpsertReleaseRequest.prototype.setRelease = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertReleaseRequest} returns this
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.prototype.clearRelease = function() {
  return this.setRelease(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertReleaseRequest.prototype.hasRelease = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.UpsertReleaseResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.UpsertReleaseResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    release: (f = msg.getRelease()) && proto.hashicorp.waypoint.Release.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.UpsertReleaseResponse}
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.UpsertReleaseResponse;
  return proto.hashicorp.waypoint.UpsertReleaseResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.UpsertReleaseResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.UpsertReleaseResponse}
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Release;
      reader.readMessage(value,proto.hashicorp.waypoint.Release.deserializeBinaryFromReader);
      msg.setRelease(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.UpsertReleaseResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.UpsertReleaseResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRelease();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.Release.serializeBinaryToWriter
    );
  }
};


/**
 * optional Release release = 1;
 * @return {?proto.hashicorp.waypoint.Release}
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.prototype.getRelease = function() {
  return /** @type{?proto.hashicorp.waypoint.Release} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Release, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.Release|undefined} value
 * @return {!proto.hashicorp.waypoint.UpsertReleaseResponse} returns this
*/
proto.hashicorp.waypoint.UpsertReleaseResponse.prototype.setRelease = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.UpsertReleaseResponse} returns this
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.prototype.clearRelease = function() {
  return this.setRelease(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.UpsertReleaseResponse.prototype.hasRelease = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Release.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Release.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Release} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Release.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    workspace: (f = msg.getWorkspace()) && proto.hashicorp.waypoint.Ref.Workspace.toObject(includeInstance, f),
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    status: (f = msg.getStatus()) && proto.hashicorp.waypoint.Status.toObject(includeInstance, f),
    component: (f = msg.getComponent()) && proto.hashicorp.waypoint.Component.toObject(includeInstance, f),
    release: (f = msg.getRelease()) && google_protobuf_any_pb.Any.toObject(includeInstance, f),
    trafficSplit: (f = msg.getTrafficSplit()) && proto.hashicorp.waypoint.Release.Split.toObject(includeInstance, f),
    labelsMap: (f = msg.getLabelsMap()) ? f.toObject(includeInstance, undefined) : [],
    url: jspb.Message.getFieldWithDefault(msg, 9, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Release}
 */
proto.hashicorp.waypoint.Release.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Release;
  return proto.hashicorp.waypoint.Release.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Release} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Release}
 */
proto.hashicorp.waypoint.Release.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 7:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 8:
      var value = new proto.hashicorp.waypoint.Ref.Workspace;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Workspace.deserializeBinaryFromReader);
      msg.setWorkspace(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.Status;
      reader.readMessage(value,proto.hashicorp.waypoint.Status.deserializeBinaryFromReader);
      msg.setStatus(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Component;
      reader.readMessage(value,proto.hashicorp.waypoint.Component.deserializeBinaryFromReader);
      msg.setComponent(value);
      break;
    case 4:
      var value = new google_protobuf_any_pb.Any;
      reader.readMessage(value,google_protobuf_any_pb.Any.deserializeBinaryFromReader);
      msg.setRelease(value);
      break;
    case 5:
      var value = new proto.hashicorp.waypoint.Release.Split;
      reader.readMessage(value,proto.hashicorp.waypoint.Release.Split.deserializeBinaryFromReader);
      msg.setTrafficSplit(value);
      break;
    case 6:
      var value = msg.getLabelsMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrl(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Release.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Release.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Release} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Release.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getWorkspace();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.hashicorp.waypoint.Ref.Workspace.serializeBinaryToWriter
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getStatus();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Status.serializeBinaryToWriter
    );
  }
  f = message.getComponent();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Component.serializeBinaryToWriter
    );
  }
  f = message.getRelease();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_any_pb.Any.serializeBinaryToWriter
    );
  }
  f = message.getTrafficSplit();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.hashicorp.waypoint.Release.Split.serializeBinaryToWriter
    );
  }
  f = message.getLabelsMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(6, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
  f = message.getUrl();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.Release.Split.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Release.Split.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Release.Split.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Release.Split} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Release.Split.toObject = function(includeInstance, msg) {
  var f, obj = {
    targetsList: jspb.Message.toObjectList(msg.getTargetsList(),
    proto.hashicorp.waypoint.Release.SplitTarget.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Release.Split}
 */
proto.hashicorp.waypoint.Release.Split.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Release.Split;
  return proto.hashicorp.waypoint.Release.Split.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Release.Split} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Release.Split}
 */
proto.hashicorp.waypoint.Release.Split.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.Release.SplitTarget;
      reader.readMessage(value,proto.hashicorp.waypoint.Release.SplitTarget.deserializeBinaryFromReader);
      msg.addTargets(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Release.Split.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Release.Split.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Release.Split} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Release.Split.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTargetsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.Release.SplitTarget.serializeBinaryToWriter
    );
  }
};


/**
 * repeated SplitTarget targets = 1;
 * @return {!Array<!proto.hashicorp.waypoint.Release.SplitTarget>}
 */
proto.hashicorp.waypoint.Release.Split.prototype.getTargetsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.Release.SplitTarget>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.Release.SplitTarget, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.Release.SplitTarget>} value
 * @return {!proto.hashicorp.waypoint.Release.Split} returns this
*/
proto.hashicorp.waypoint.Release.Split.prototype.setTargetsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.Release.SplitTarget=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.Release.SplitTarget}
 */
proto.hashicorp.waypoint.Release.Split.prototype.addTargets = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.Release.SplitTarget, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.Release.Split} returns this
 */
proto.hashicorp.waypoint.Release.Split.prototype.clearTargetsList = function() {
  return this.setTargetsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Release.SplitTarget.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Release.SplitTarget.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Release.SplitTarget} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Release.SplitTarget.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    percent: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Release.SplitTarget}
 */
proto.hashicorp.waypoint.Release.SplitTarget.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Release.SplitTarget;
  return proto.hashicorp.waypoint.Release.SplitTarget.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Release.SplitTarget} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Release.SplitTarget}
 */
proto.hashicorp.waypoint.Release.SplitTarget.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setPercent(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Release.SplitTarget.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Release.SplitTarget.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Release.SplitTarget} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Release.SplitTarget.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPercent();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Release.SplitTarget.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Release.SplitTarget} returns this
 */
proto.hashicorp.waypoint.Release.SplitTarget.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int32 percent = 2;
 * @return {number}
 */
proto.hashicorp.waypoint.Release.SplitTarget.prototype.getPercent = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.Release.SplitTarget} returns this
 */
proto.hashicorp.waypoint.Release.SplitTarget.prototype.setPercent = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional Ref.Application application = 7;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.Release.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 7));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
*/
proto.hashicorp.waypoint.Release.prototype.setApplication = function(value) {
  return jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Release.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional Ref.Workspace workspace = 8;
 * @return {?proto.hashicorp.waypoint.Ref.Workspace}
 */
proto.hashicorp.waypoint.Release.prototype.getWorkspace = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Workspace} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Workspace, 8));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Workspace|undefined} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
*/
proto.hashicorp.waypoint.Release.prototype.setWorkspace = function(value) {
  return jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.clearWorkspace = function() {
  return this.setWorkspace(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Release.prototype.hasWorkspace = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Release.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional Status status = 2;
 * @return {?proto.hashicorp.waypoint.Status}
 */
proto.hashicorp.waypoint.Release.prototype.getStatus = function() {
  return /** @type{?proto.hashicorp.waypoint.Status} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Status, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
*/
proto.hashicorp.waypoint.Release.prototype.setStatus = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.clearStatus = function() {
  return this.setStatus(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Release.prototype.hasStatus = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Component component = 3;
 * @return {?proto.hashicorp.waypoint.Component}
 */
proto.hashicorp.waypoint.Release.prototype.getComponent = function() {
  return /** @type{?proto.hashicorp.waypoint.Component} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Component, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Component|undefined} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
*/
proto.hashicorp.waypoint.Release.prototype.setComponent = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.clearComponent = function() {
  return this.setComponent(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Release.prototype.hasComponent = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional google.protobuf.Any release = 4;
 * @return {?proto.google.protobuf.Any}
 */
proto.hashicorp.waypoint.Release.prototype.getRelease = function() {
  return /** @type{?proto.google.protobuf.Any} */ (
    jspb.Message.getWrapperField(this, google_protobuf_any_pb.Any, 4));
};


/**
 * @param {?proto.google.protobuf.Any|undefined} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
*/
proto.hashicorp.waypoint.Release.prototype.setRelease = function(value) {
  return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.clearRelease = function() {
  return this.setRelease(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Release.prototype.hasRelease = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional Split traffic_split = 5;
 * @return {?proto.hashicorp.waypoint.Release.Split}
 */
proto.hashicorp.waypoint.Release.prototype.getTrafficSplit = function() {
  return /** @type{?proto.hashicorp.waypoint.Release.Split} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Release.Split, 5));
};


/**
 * @param {?proto.hashicorp.waypoint.Release.Split|undefined} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
*/
proto.hashicorp.waypoint.Release.prototype.setTrafficSplit = function(value) {
  return jspb.Message.setWrapperField(this, 5, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.clearTrafficSplit = function() {
  return this.setTrafficSplit(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Release.prototype.hasTrafficSplit = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * map<string, string> labels = 6;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.hashicorp.waypoint.Release.prototype.getLabelsMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 6, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.clearLabelsMap = function() {
  this.getLabelsMap().clear();
  return this;};


/**
 * optional string url = 9;
 * @return {string}
 */
proto.hashicorp.waypoint.Release.prototype.getUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Release} returns this
 */
proto.hashicorp.waypoint.Release.prototype.setUrl = function(value) {
  return jspb.Message.setProto3StringField(this, 9, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.GetLogStreamRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.GetLogStreamRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.GetLogStreamRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetLogStreamRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.GetLogStreamRequest}
 */
proto.hashicorp.waypoint.GetLogStreamRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.GetLogStreamRequest;
  return proto.hashicorp.waypoint.GetLogStreamRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.GetLogStreamRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.GetLogStreamRequest}
 */
proto.hashicorp.waypoint.GetLogStreamRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.GetLogStreamRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.GetLogStreamRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.GetLogStreamRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.GetLogStreamRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.GetLogStreamRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.GetLogStreamRequest} returns this
 */
proto.hashicorp.waypoint.GetLogStreamRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.LogBatch.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.LogBatch.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.LogBatch.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.LogBatch} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.LogBatch.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    instanceId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    linesList: jspb.Message.toObjectList(msg.getLinesList(),
    proto.hashicorp.waypoint.LogBatch.Entry.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.LogBatch}
 */
proto.hashicorp.waypoint.LogBatch.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.LogBatch;
  return proto.hashicorp.waypoint.LogBatch.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.LogBatch} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.LogBatch}
 */
proto.hashicorp.waypoint.LogBatch.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setInstanceId(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.LogBatch.Entry;
      reader.readMessage(value,proto.hashicorp.waypoint.LogBatch.Entry.deserializeBinaryFromReader);
      msg.addLines(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.LogBatch.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.LogBatch.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.LogBatch} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.LogBatch.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getInstanceId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLinesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      3,
      f,
      proto.hashicorp.waypoint.LogBatch.Entry.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.LogBatch.Entry.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.LogBatch.Entry.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.LogBatch.Entry} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.LogBatch.Entry.toObject = function(includeInstance, msg) {
  var f, obj = {
    timestamp: (f = msg.getTimestamp()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    line: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.LogBatch.Entry}
 */
proto.hashicorp.waypoint.LogBatch.Entry.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.LogBatch.Entry;
  return proto.hashicorp.waypoint.LogBatch.Entry.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.LogBatch.Entry} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.LogBatch.Entry}
 */
proto.hashicorp.waypoint.LogBatch.Entry.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setTimestamp(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setLine(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.LogBatch.Entry.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.LogBatch.Entry.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.LogBatch.Entry} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.LogBatch.Entry.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getTimestamp();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getLine();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional google.protobuf.Timestamp timestamp = 1;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.LogBatch.Entry.prototype.getTimestamp = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 1));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.LogBatch.Entry} returns this
*/
proto.hashicorp.waypoint.LogBatch.Entry.prototype.setTimestamp = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.LogBatch.Entry} returns this
 */
proto.hashicorp.waypoint.LogBatch.Entry.prototype.clearTimestamp = function() {
  return this.setTimestamp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.LogBatch.Entry.prototype.hasTimestamp = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional string line = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.LogBatch.Entry.prototype.getLine = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.LogBatch.Entry} returns this
 */
proto.hashicorp.waypoint.LogBatch.Entry.prototype.setLine = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.LogBatch.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.LogBatch} returns this
 */
proto.hashicorp.waypoint.LogBatch.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string instance_id = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.LogBatch.prototype.getInstanceId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.LogBatch} returns this
 */
proto.hashicorp.waypoint.LogBatch.prototype.setInstanceId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * repeated Entry lines = 3;
 * @return {!Array<!proto.hashicorp.waypoint.LogBatch.Entry>}
 */
proto.hashicorp.waypoint.LogBatch.prototype.getLinesList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.LogBatch.Entry>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.LogBatch.Entry, 3));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.LogBatch.Entry>} value
 * @return {!proto.hashicorp.waypoint.LogBatch} returns this
*/
proto.hashicorp.waypoint.LogBatch.prototype.setLinesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 3, value);
};


/**
 * @param {!proto.hashicorp.waypoint.LogBatch.Entry=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.LogBatch.Entry}
 */
proto.hashicorp.waypoint.LogBatch.prototype.addLines = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 3, opt_value, proto.hashicorp.waypoint.LogBatch.Entry, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.LogBatch} returns this
 */
proto.hashicorp.waypoint.LogBatch.prototype.clearLinesList = function() {
  return this.setLinesList([]);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.ConfigVar.oneofGroups_ = [[3,4,5]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.ConfigVar.ScopeCase = {
  SCOPE_NOT_SET: 0,
  APPLICATION: 3,
  PROJECT: 4,
  RUNNER: 5
};

/**
 * @return {proto.hashicorp.waypoint.ConfigVar.ScopeCase}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.getScopeCase = function() {
  return /** @type {proto.hashicorp.waypoint.ConfigVar.ScopeCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.ConfigVar.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ConfigVar.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ConfigVar} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigVar.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    project: (f = msg.getProject()) && proto.hashicorp.waypoint.Ref.Project.toObject(includeInstance, f),
    runner: (f = msg.getRunner()) && proto.hashicorp.waypoint.Ref.Runner.toObject(includeInstance, f),
    name: jspb.Message.getFieldWithDefault(msg, 1, ""),
    value: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ConfigVar}
 */
proto.hashicorp.waypoint.ConfigVar.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ConfigVar;
  return proto.hashicorp.waypoint.ConfigVar.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ConfigVar} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ConfigVar}
 */
proto.hashicorp.waypoint.ConfigVar.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 3:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Ref.Project;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Project.deserializeBinaryFromReader);
      msg.setProject(value);
      break;
    case 5:
      var value = new proto.hashicorp.waypoint.Ref.Runner;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Runner.deserializeBinaryFromReader);
      msg.setRunner(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ConfigVar.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ConfigVar} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigVar.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getProject();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Ref.Project.serializeBinaryToWriter
    );
  }
  f = message.getRunner();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.hashicorp.waypoint.Ref.Runner.serializeBinaryToWriter
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getValue();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional Ref.Application application = 3;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
*/
proto.hashicorp.waypoint.ConfigVar.prototype.setApplication = function(value) {
  return jspb.Message.setOneofWrapperField(this, 3, proto.hashicorp.waypoint.ConfigVar.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
 */
proto.hashicorp.waypoint.ConfigVar.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Ref.Project project = 4;
 * @return {?proto.hashicorp.waypoint.Ref.Project}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.getProject = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Project} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Project, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Project|undefined} value
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
*/
proto.hashicorp.waypoint.ConfigVar.prototype.setProject = function(value) {
  return jspb.Message.setOneofWrapperField(this, 4, proto.hashicorp.waypoint.ConfigVar.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
 */
proto.hashicorp.waypoint.ConfigVar.prototype.clearProject = function() {
  return this.setProject(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.hasProject = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional Ref.Runner runner = 5;
 * @return {?proto.hashicorp.waypoint.Ref.Runner}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.getRunner = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Runner} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Runner, 5));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Runner|undefined} value
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
*/
proto.hashicorp.waypoint.ConfigVar.prototype.setRunner = function(value) {
  return jspb.Message.setOneofWrapperField(this, 5, proto.hashicorp.waypoint.ConfigVar.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
 */
proto.hashicorp.waypoint.ConfigVar.prototype.clearRunner = function() {
  return this.setRunner(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.hasRunner = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional string name = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
 */
proto.hashicorp.waypoint.ConfigVar.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string value = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.ConfigVar.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.ConfigVar} returns this
 */
proto.hashicorp.waypoint.ConfigVar.prototype.setValue = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ConfigSetRequest.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ConfigSetRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ConfigSetRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ConfigSetRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigSetRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    variablesList: jspb.Message.toObjectList(msg.getVariablesList(),
    proto.hashicorp.waypoint.ConfigVar.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ConfigSetRequest}
 */
proto.hashicorp.waypoint.ConfigSetRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ConfigSetRequest;
  return proto.hashicorp.waypoint.ConfigSetRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ConfigSetRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ConfigSetRequest}
 */
proto.hashicorp.waypoint.ConfigSetRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.ConfigVar;
      reader.readMessage(value,proto.hashicorp.waypoint.ConfigVar.deserializeBinaryFromReader);
      msg.addVariables(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ConfigSetRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ConfigSetRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ConfigSetRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigSetRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getVariablesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.ConfigVar.serializeBinaryToWriter
    );
  }
};


/**
 * repeated ConfigVar variables = 1;
 * @return {!Array<!proto.hashicorp.waypoint.ConfigVar>}
 */
proto.hashicorp.waypoint.ConfigSetRequest.prototype.getVariablesList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.ConfigVar>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.ConfigVar, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.ConfigVar>} value
 * @return {!proto.hashicorp.waypoint.ConfigSetRequest} returns this
*/
proto.hashicorp.waypoint.ConfigSetRequest.prototype.setVariablesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.ConfigVar=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.ConfigVar}
 */
proto.hashicorp.waypoint.ConfigSetRequest.prototype.addVariables = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.ConfigVar, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ConfigSetRequest} returns this
 */
proto.hashicorp.waypoint.ConfigSetRequest.prototype.clearVariablesList = function() {
  return this.setVariablesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ConfigSetResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ConfigSetResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ConfigSetResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigSetResponse.toObject = function(includeInstance, msg) {
  var f, obj = {

  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ConfigSetResponse}
 */
proto.hashicorp.waypoint.ConfigSetResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ConfigSetResponse;
  return proto.hashicorp.waypoint.ConfigSetResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ConfigSetResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ConfigSetResponse}
 */
proto.hashicorp.waypoint.ConfigSetResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ConfigSetResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ConfigSetResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ConfigSetResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigSetResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.ConfigGetRequest.oneofGroups_ = [[2,3,4]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.ConfigGetRequest.ScopeCase = {
  SCOPE_NOT_SET: 0,
  APPLICATION: 2,
  PROJECT: 3,
  RUNNER: 4
};

/**
 * @return {proto.hashicorp.waypoint.ConfigGetRequest.ScopeCase}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.getScopeCase = function() {
  return /** @type {proto.hashicorp.waypoint.ConfigGetRequest.ScopeCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.ConfigGetRequest.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ConfigGetRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ConfigGetRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigGetRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    application: (f = msg.getApplication()) && proto.hashicorp.waypoint.Ref.Application.toObject(includeInstance, f),
    project: (f = msg.getProject()) && proto.hashicorp.waypoint.Ref.Project.toObject(includeInstance, f),
    runner: (f = msg.getRunner()) && proto.hashicorp.waypoint.Ref.RunnerId.toObject(includeInstance, f),
    prefix: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest}
 */
proto.hashicorp.waypoint.ConfigGetRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ConfigGetRequest;
  return proto.hashicorp.waypoint.ConfigGetRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ConfigGetRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest}
 */
proto.hashicorp.waypoint.ConfigGetRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = new proto.hashicorp.waypoint.Ref.Application;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Application.deserializeBinaryFromReader);
      msg.setApplication(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.Ref.Project;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.Project.deserializeBinaryFromReader);
      msg.setProject(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.Ref.RunnerId;
      reader.readMessage(value,proto.hashicorp.waypoint.Ref.RunnerId.deserializeBinaryFromReader);
      msg.setRunner(value);
      break;
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPrefix(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ConfigGetRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ConfigGetRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigGetRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getApplication();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.Ref.Application.serializeBinaryToWriter
    );
  }
  f = message.getProject();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.Ref.Project.serializeBinaryToWriter
    );
  }
  f = message.getRunner();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.Ref.RunnerId.serializeBinaryToWriter
    );
  }
  f = message.getPrefix();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional Ref.Application application = 2;
 * @return {?proto.hashicorp.waypoint.Ref.Application}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.getApplication = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Application} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Application, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Application|undefined} value
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest} returns this
*/
proto.hashicorp.waypoint.ConfigGetRequest.prototype.setApplication = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.ConfigGetRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest} returns this
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.clearApplication = function() {
  return this.setApplication(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.hasApplication = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Ref.Project project = 3;
 * @return {?proto.hashicorp.waypoint.Ref.Project}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.getProject = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.Project} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.Project, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.Project|undefined} value
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest} returns this
*/
proto.hashicorp.waypoint.ConfigGetRequest.prototype.setProject = function(value) {
  return jspb.Message.setOneofWrapperField(this, 3, proto.hashicorp.waypoint.ConfigGetRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest} returns this
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.clearProject = function() {
  return this.setProject(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.hasProject = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Ref.RunnerId runner = 4;
 * @return {?proto.hashicorp.waypoint.Ref.RunnerId}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.getRunner = function() {
  return /** @type{?proto.hashicorp.waypoint.Ref.RunnerId} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.Ref.RunnerId, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.Ref.RunnerId|undefined} value
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest} returns this
*/
proto.hashicorp.waypoint.ConfigGetRequest.prototype.setRunner = function(value) {
  return jspb.Message.setOneofWrapperField(this, 4, proto.hashicorp.waypoint.ConfigGetRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest} returns this
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.clearRunner = function() {
  return this.setRunner(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.hasRunner = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional string prefix = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.getPrefix = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.ConfigGetRequest} returns this
 */
proto.hashicorp.waypoint.ConfigGetRequest.prototype.setPrefix = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ConfigGetResponse.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ConfigGetResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ConfigGetResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ConfigGetResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigGetResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    variablesList: jspb.Message.toObjectList(msg.getVariablesList(),
    proto.hashicorp.waypoint.ConfigVar.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ConfigGetResponse}
 */
proto.hashicorp.waypoint.ConfigGetResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ConfigGetResponse;
  return proto.hashicorp.waypoint.ConfigGetResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ConfigGetResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ConfigGetResponse}
 */
proto.hashicorp.waypoint.ConfigGetResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.ConfigVar;
      reader.readMessage(value,proto.hashicorp.waypoint.ConfigVar.deserializeBinaryFromReader);
      msg.addVariables(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ConfigGetResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ConfigGetResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ConfigGetResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConfigGetResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getVariablesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.ConfigVar.serializeBinaryToWriter
    );
  }
};


/**
 * repeated ConfigVar variables = 1;
 * @return {!Array<!proto.hashicorp.waypoint.ConfigVar>}
 */
proto.hashicorp.waypoint.ConfigGetResponse.prototype.getVariablesList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.ConfigVar>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.ConfigVar, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.ConfigVar>} value
 * @return {!proto.hashicorp.waypoint.ConfigGetResponse} returns this
*/
proto.hashicorp.waypoint.ConfigGetResponse.prototype.setVariablesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.ConfigVar=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.ConfigVar}
 */
proto.hashicorp.waypoint.ConfigGetResponse.prototype.addVariables = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.ConfigVar, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ConfigGetResponse} returns this
 */
proto.hashicorp.waypoint.ConfigGetResponse.prototype.clearVariablesList = function() {
  return this.setVariablesList([]);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.ExecStreamRequest.oneofGroups_ = [[1,2,3]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.ExecStreamRequest.EventCase = {
  EVENT_NOT_SET: 0,
  START: 1,
  INPUT: 2,
  WINCH: 3
};

/**
 * @return {proto.hashicorp.waypoint.ExecStreamRequest.EventCase}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.ExecStreamRequest.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.ExecStreamRequest.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    start: (f = msg.getStart()) && proto.hashicorp.waypoint.ExecStreamRequest.Start.toObject(includeInstance, f),
    input: (f = msg.getInput()) && proto.hashicorp.waypoint.ExecStreamRequest.Input.toObject(includeInstance, f),
    winch: (f = msg.getWinch()) && proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest}
 */
proto.hashicorp.waypoint.ExecStreamRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamRequest;
  return proto.hashicorp.waypoint.ExecStreamRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest}
 */
proto.hashicorp.waypoint.ExecStreamRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.ExecStreamRequest.Start;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamRequest.Start.deserializeBinaryFromReader);
      msg.setStart(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.ExecStreamRequest.Input;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamRequest.Input.deserializeBinaryFromReader);
      msg.setInput(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.ExecStreamRequest.WindowSize;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.deserializeBinaryFromReader);
      msg.setWinch(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getStart();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.ExecStreamRequest.Start.serializeBinaryToWriter
    );
  }
  f = message.getInput();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.ExecStreamRequest.Input.serializeBinaryToWriter
    );
  }
  f = message.getWinch();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.serializeBinaryToWriter
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.repeatedFields_ = [2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamRequest.Start.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.Start} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    argsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f,
    pty: (f = msg.getPty()) && proto.hashicorp.waypoint.ExecStreamRequest.PTY.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamRequest.Start;
  return proto.hashicorp.waypoint.ExecStreamRequest.Start.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.Start} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addArgs(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.ExecStreamRequest.PTY;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamRequest.PTY.deserializeBinaryFromReader);
      msg.setPty(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamRequest.Start.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.Start} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getArgsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
  f = message.getPty();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.ExecStreamRequest.PTY.serializeBinaryToWriter
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated string args = 2;
 * @return {!Array<string>}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.getArgsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.setArgsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.addArgs = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.clearArgsList = function() {
  return this.setArgsList([]);
};


/**
 * optional PTY pty = 3;
 * @return {?proto.hashicorp.waypoint.ExecStreamRequest.PTY}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.getPty = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamRequest.PTY} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamRequest.PTY, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamRequest.PTY|undefined} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start} returns this
*/
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.setPty = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Start} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.clearPty = function() {
  return this.setPty(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Start.prototype.hasPty = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamRequest.Input.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.Input} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.toObject = function(includeInstance, msg) {
  var f, obj = {
    data: msg.getData_asB64()
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Input}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamRequest.Input;
  return proto.hashicorp.waypoint.ExecStreamRequest.Input.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.Input} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Input}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setData(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamRequest.Input.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.Input} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getData_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      1,
      f
    );
  }
};


/**
 * optional bytes data = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * optional bytes data = 1;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.prototype.getData_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getData()));
};


/**
 * optional bytes data = 1;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getData()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.Input} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.Input.prototype.setData = function(value) {
  return jspb.Message.setProto3BytesField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamRequest.PTY.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.PTY} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.toObject = function(includeInstance, msg) {
  var f, obj = {
    enable: jspb.Message.getBooleanFieldWithDefault(msg, 1, false),
    term: jspb.Message.getFieldWithDefault(msg, 2, ""),
    windowSize: (f = msg.getWindowSize()) && proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.PTY}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamRequest.PTY;
  return proto.hashicorp.waypoint.ExecStreamRequest.PTY.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.PTY} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.PTY}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setEnable(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setTerm(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.ExecStreamRequest.WindowSize;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.deserializeBinaryFromReader);
      msg.setWindowSize(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamRequest.PTY.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.PTY} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEnable();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getTerm();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getWindowSize();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.serializeBinaryToWriter
    );
  }
};


/**
 * optional bool enable = 1;
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.getEnable = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 1, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.PTY} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.setEnable = function(value) {
  return jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * optional string term = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.getTerm = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.PTY} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.setTerm = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional WindowSize window_size = 3;
 * @return {?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.getWindowSize = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamRequest.WindowSize, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize|undefined} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.PTY} returns this
*/
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.setWindowSize = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.PTY} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.clearWindowSize = function() {
  return this.setWindowSize(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamRequest.PTY.prototype.hasWindowSize = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.toObject = function(includeInstance, msg) {
  var f, obj = {
    rows: jspb.Message.getFieldWithDefault(msg, 1, 0),
    cols: jspb.Message.getFieldWithDefault(msg, 2, 0),
    width: jspb.Message.getFieldWithDefault(msg, 3, 0),
    height: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamRequest.WindowSize;
  return proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setRows(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setCols(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setWidth(value);
      break;
    case 4:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setHeight(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRows();
  if (f !== 0) {
    writer.writeInt32(
      1,
      f
    );
  }
  f = message.getCols();
  if (f !== 0) {
    writer.writeInt32(
      2,
      f
    );
  }
  f = message.getWidth();
  if (f !== 0) {
    writer.writeInt32(
      3,
      f
    );
  }
  f = message.getHeight();
  if (f !== 0) {
    writer.writeInt32(
      4,
      f
    );
  }
};


/**
 * optional int32 rows = 1;
 * @return {number}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.getRows = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.setRows = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional int32 cols = 2;
 * @return {number}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.getCols = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.setCols = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional int32 width = 3;
 * @return {number}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.getWidth = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.setWidth = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional int32 height = 4;
 * @return {number}
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.getHeight = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.prototype.setHeight = function(value) {
  return jspb.Message.setProto3IntField(this, 4, value);
};


/**
 * optional Start start = 1;
 * @return {?proto.hashicorp.waypoint.ExecStreamRequest.Start}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.getStart = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamRequest.Start} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamRequest.Start, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamRequest.Start|undefined} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest} returns this
*/
proto.hashicorp.waypoint.ExecStreamRequest.prototype.setStart = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.ExecStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.clearStart = function() {
  return this.setStart(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.hasStart = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Input input = 2;
 * @return {?proto.hashicorp.waypoint.ExecStreamRequest.Input}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.getInput = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamRequest.Input} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamRequest.Input, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamRequest.Input|undefined} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest} returns this
*/
proto.hashicorp.waypoint.ExecStreamRequest.prototype.setInput = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.ExecStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.clearInput = function() {
  return this.setInput(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.hasInput = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional WindowSize winch = 3;
 * @return {?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.getWinch = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamRequest.WindowSize, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize|undefined} value
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest} returns this
*/
proto.hashicorp.waypoint.ExecStreamRequest.prototype.setWinch = function(value) {
  return jspb.Message.setOneofWrapperField(this, 3, proto.hashicorp.waypoint.ExecStreamRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ExecStreamRequest} returns this
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.clearWinch = function() {
  return this.setWinch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamRequest.prototype.hasWinch = function() {
  return jspb.Message.getField(this, 3) != null;
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.ExecStreamResponse.oneofGroups_ = [[1,2]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.ExecStreamResponse.EventCase = {
  EVENT_NOT_SET: 0,
  OUTPUT: 1,
  EXIT: 2
};

/**
 * @return {proto.hashicorp.waypoint.ExecStreamResponse.EventCase}
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.ExecStreamResponse.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.ExecStreamResponse.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    output: (f = msg.getOutput()) && proto.hashicorp.waypoint.ExecStreamResponse.Output.toObject(includeInstance, f),
    exit: (f = msg.getExit()) && proto.hashicorp.waypoint.ExecStreamResponse.Exit.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse}
 */
proto.hashicorp.waypoint.ExecStreamResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamResponse;
  return proto.hashicorp.waypoint.ExecStreamResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse}
 */
proto.hashicorp.waypoint.ExecStreamResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.ExecStreamResponse.Output;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamResponse.Output.deserializeBinaryFromReader);
      msg.setOutput(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.ExecStreamResponse.Exit;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamResponse.Exit.deserializeBinaryFromReader);
      msg.setExit(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOutput();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.ExecStreamResponse.Output.serializeBinaryToWriter
    );
  }
  f = message.getExit();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.ExecStreamResponse.Exit.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamResponse.Exit.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse.Exit} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.toObject = function(includeInstance, msg) {
  var f, obj = {
    code: jspb.Message.getFieldWithDefault(msg, 1, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Exit}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamResponse.Exit;
  return proto.hashicorp.waypoint.ExecStreamResponse.Exit.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse.Exit} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Exit}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamResponse.Exit.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse.Exit} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCode();
  if (f !== 0) {
    writer.writeInt32(
      1,
      f
    );
  }
};


/**
 * optional int32 code = 1;
 * @return {number}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.prototype.getCode = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Exit} returns this
 */
proto.hashicorp.waypoint.ExecStreamResponse.Exit.prototype.setCode = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ExecStreamResponse.Output.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse.Output} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.toObject = function(includeInstance, msg) {
  var f, obj = {
    channel: jspb.Message.getFieldWithDefault(msg, 1, 0),
    data: msg.getData_asB64()
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Output}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ExecStreamResponse.Output;
  return proto.hashicorp.waypoint.ExecStreamResponse.Output.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse.Output} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Output}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.hashicorp.waypoint.ExecStreamResponse.Output.Channel} */ (reader.readEnum());
      msg.setChannel(value);
      break;
    case 2:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setData(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ExecStreamResponse.Output.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse.Output} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getChannel();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getData_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      2,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.Channel = {
  UNKNOWN: 0,
  STDOUT: 1,
  STDERR: 2
};

/**
 * optional Channel channel = 1;
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Output.Channel}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.getChannel = function() {
  return /** @type {!proto.hashicorp.waypoint.ExecStreamResponse.Output.Channel} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.ExecStreamResponse.Output.Channel} value
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Output} returns this
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.setChannel = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional bytes data = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * optional bytes data = 2;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.getData_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getData()));
};


/**
 * optional bytes data = 2;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getData()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse.Output} returns this
 */
proto.hashicorp.waypoint.ExecStreamResponse.Output.prototype.setData = function(value) {
  return jspb.Message.setProto3BytesField(this, 2, value);
};


/**
 * optional Output output = 1;
 * @return {?proto.hashicorp.waypoint.ExecStreamResponse.Output}
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.getOutput = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamResponse.Output} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamResponse.Output, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamResponse.Output|undefined} value
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse} returns this
*/
proto.hashicorp.waypoint.ExecStreamResponse.prototype.setOutput = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.ExecStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse} returns this
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.clearOutput = function() {
  return this.setOutput(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.hasOutput = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Exit exit = 2;
 * @return {?proto.hashicorp.waypoint.ExecStreamResponse.Exit}
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.getExit = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamResponse.Exit} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamResponse.Exit, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamResponse.Exit|undefined} value
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse} returns this
*/
proto.hashicorp.waypoint.ExecStreamResponse.prototype.setExit = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.ExecStreamResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.ExecStreamResponse} returns this
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.clearExit = function() {
  return this.setExit(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.ExecStreamResponse.prototype.hasExit = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointConfigRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointConfigRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    deploymentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    instanceId: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointConfigRequest}
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointConfigRequest;
  return proto.hashicorp.waypoint.EntrypointConfigRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointConfigRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointConfigRequest}
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDeploymentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setInstanceId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointConfigRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointConfigRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDeploymentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getInstanceId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string deployment_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.prototype.getDeploymentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfigRequest} returns this
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.prototype.setDeploymentId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string instance_id = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.prototype.getInstanceId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfigRequest} returns this
 */
proto.hashicorp.waypoint.EntrypointConfigRequest.prototype.setInstanceId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointConfigResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointConfigResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    config: (f = msg.getConfig()) && proto.hashicorp.waypoint.EntrypointConfig.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointConfigResponse}
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointConfigResponse;
  return proto.hashicorp.waypoint.EntrypointConfigResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointConfigResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointConfigResponse}
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 2:
      var value = new proto.hashicorp.waypoint.EntrypointConfig;
      reader.readMessage(value,proto.hashicorp.waypoint.EntrypointConfig.deserializeBinaryFromReader);
      msg.setConfig(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointConfigResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointConfigResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getConfig();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.EntrypointConfig.serializeBinaryToWriter
    );
  }
};


/**
 * optional EntrypointConfig config = 2;
 * @return {?proto.hashicorp.waypoint.EntrypointConfig}
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.prototype.getConfig = function() {
  return /** @type{?proto.hashicorp.waypoint.EntrypointConfig} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.EntrypointConfig, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.EntrypointConfig|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfigResponse} returns this
*/
proto.hashicorp.waypoint.EntrypointConfigResponse.prototype.setConfig = function(value) {
  return jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointConfigResponse} returns this
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.prototype.clearConfig = function() {
  return this.setConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointConfigResponse.prototype.hasConfig = function() {
  return jspb.Message.getField(this, 2) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.EntrypointConfig.repeatedFields_ = [1,2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    execList: jspb.Message.toObjectList(msg.getExecList(),
    proto.hashicorp.waypoint.EntrypointConfig.Exec.toObject, includeInstance),
    envVarsList: jspb.Message.toObjectList(msg.getEnvVarsList(),
    proto.hashicorp.waypoint.ConfigVar.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig}
 */
proto.hashicorp.waypoint.EntrypointConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointConfig;
  return proto.hashicorp.waypoint.EntrypointConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig}
 */
proto.hashicorp.waypoint.EntrypointConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.EntrypointConfig.Exec;
      reader.readMessage(value,proto.hashicorp.waypoint.EntrypointConfig.Exec.deserializeBinaryFromReader);
      msg.addExec(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.ConfigVar;
      reader.readMessage(value,proto.hashicorp.waypoint.ConfigVar.deserializeBinaryFromReader);
      msg.addEnvVars(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getExecList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.hashicorp.waypoint.EntrypointConfig.Exec.serializeBinaryToWriter
    );
  }
  f = message.getEnvVarsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto.hashicorp.waypoint.ConfigVar.serializeBinaryToWriter
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.repeatedFields_ = [2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointConfig.Exec.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointConfig.Exec} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.toObject = function(includeInstance, msg) {
  var f, obj = {
    index: jspb.Message.getFieldWithDefault(msg, 1, 0),
    argsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f,
    pty: (f = msg.getPty()) && proto.hashicorp.waypoint.ExecStreamRequest.PTY.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointConfig.Exec;
  return proto.hashicorp.waypoint.EntrypointConfig.Exec.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointConfig.Exec} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setIndex(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addArgs(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.ExecStreamRequest.PTY;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamRequest.PTY.deserializeBinaryFromReader);
      msg.setPty(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointConfig.Exec.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointConfig.Exec} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getIndex();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getArgsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
  f = message.getPty();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.ExecStreamRequest.PTY.serializeBinaryToWriter
    );
  }
};


/**
 * optional int64 index = 1;
 * @return {number}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.getIndex = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec} returns this
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.setIndex = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * repeated string args = 2;
 * @return {!Array<string>}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.getArgsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec} returns this
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.setArgsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec} returns this
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.addArgs = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec} returns this
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.clearArgsList = function() {
  return this.setArgsList([]);
};


/**
 * optional ExecStreamRequest.PTY pty = 3;
 * @return {?proto.hashicorp.waypoint.ExecStreamRequest.PTY}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.getPty = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamRequest.PTY} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamRequest.PTY, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamRequest.PTY|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec} returns this
*/
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.setPty = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec} returns this
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.clearPty = function() {
  return this.setPty(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointConfig.Exec.prototype.hasPty = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * repeated Exec exec = 1;
 * @return {!Array<!proto.hashicorp.waypoint.EntrypointConfig.Exec>}
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.getExecList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.EntrypointConfig.Exec>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.EntrypointConfig.Exec, 1));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.EntrypointConfig.Exec>} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfig} returns this
*/
proto.hashicorp.waypoint.EntrypointConfig.prototype.setExecList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.hashicorp.waypoint.EntrypointConfig.Exec=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.EntrypointConfig.Exec}
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.addExec = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.hashicorp.waypoint.EntrypointConfig.Exec, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig} returns this
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.clearExecList = function() {
  return this.setExecList([]);
};


/**
 * repeated ConfigVar env_vars = 2;
 * @return {!Array<!proto.hashicorp.waypoint.ConfigVar>}
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.getEnvVarsList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.ConfigVar>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.ConfigVar, 2));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.ConfigVar>} value
 * @return {!proto.hashicorp.waypoint.EntrypointConfig} returns this
*/
proto.hashicorp.waypoint.EntrypointConfig.prototype.setEnvVarsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.hashicorp.waypoint.ConfigVar=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.ConfigVar}
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.addEnvVars = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.hashicorp.waypoint.ConfigVar, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.EntrypointConfig} returns this
 */
proto.hashicorp.waypoint.EntrypointConfig.prototype.clearEnvVarsList = function() {
  return this.setEnvVarsList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.hashicorp.waypoint.EntrypointLogBatch.repeatedFields_ = [2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointLogBatch.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointLogBatch} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointLogBatch.toObject = function(includeInstance, msg) {
  var f, obj = {
    instanceId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    linesList: jspb.Message.toObjectList(msg.getLinesList(),
    proto.hashicorp.waypoint.LogBatch.Entry.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointLogBatch}
 */
proto.hashicorp.waypoint.EntrypointLogBatch.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointLogBatch;
  return proto.hashicorp.waypoint.EntrypointLogBatch.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointLogBatch} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointLogBatch}
 */
proto.hashicorp.waypoint.EntrypointLogBatch.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setInstanceId(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.LogBatch.Entry;
      reader.readMessage(value,proto.hashicorp.waypoint.LogBatch.Entry.deserializeBinaryFromReader);
      msg.addLines(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointLogBatch.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointLogBatch} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointLogBatch.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getInstanceId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLinesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto.hashicorp.waypoint.LogBatch.Entry.serializeBinaryToWriter
    );
  }
};


/**
 * optional string instance_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.getInstanceId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.EntrypointLogBatch} returns this
 */
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.setInstanceId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated LogBatch.Entry lines = 2;
 * @return {!Array<!proto.hashicorp.waypoint.LogBatch.Entry>}
 */
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.getLinesList = function() {
  return /** @type{!Array<!proto.hashicorp.waypoint.LogBatch.Entry>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.hashicorp.waypoint.LogBatch.Entry, 2));
};


/**
 * @param {!Array<!proto.hashicorp.waypoint.LogBatch.Entry>} value
 * @return {!proto.hashicorp.waypoint.EntrypointLogBatch} returns this
*/
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.setLinesList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.hashicorp.waypoint.LogBatch.Entry=} opt_value
 * @param {number=} opt_index
 * @return {!proto.hashicorp.waypoint.LogBatch.Entry}
 */
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.addLines = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.hashicorp.waypoint.LogBatch.Entry, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.hashicorp.waypoint.EntrypointLogBatch} returns this
 */
proto.hashicorp.waypoint.EntrypointLogBatch.prototype.clearLinesList = function() {
  return this.setLinesList([]);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.EntrypointExecRequest.oneofGroups_ = [[1,2,3,4]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.EventCase = {
  EVENT_NOT_SET: 0,
  OPEN: 1,
  EXIT: 2,
  OUTPUT: 3,
  ERROR: 4
};

/**
 * @return {proto.hashicorp.waypoint.EntrypointExecRequest.EventCase}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.EntrypointExecRequest.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.EntrypointExecRequest.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointExecRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    open: (f = msg.getOpen()) && proto.hashicorp.waypoint.EntrypointExecRequest.Open.toObject(includeInstance, f),
    exit: (f = msg.getExit()) && proto.hashicorp.waypoint.EntrypointExecRequest.Exit.toObject(includeInstance, f),
    output: (f = msg.getOutput()) && proto.hashicorp.waypoint.EntrypointExecRequest.Output.toObject(includeInstance, f),
    error: (f = msg.getError()) && proto.hashicorp.waypoint.EntrypointExecRequest.Error.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointExecRequest;
  return proto.hashicorp.waypoint.EntrypointExecRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.hashicorp.waypoint.EntrypointExecRequest.Open;
      reader.readMessage(value,proto.hashicorp.waypoint.EntrypointExecRequest.Open.deserializeBinaryFromReader);
      msg.setOpen(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.EntrypointExecRequest.Exit;
      reader.readMessage(value,proto.hashicorp.waypoint.EntrypointExecRequest.Exit.deserializeBinaryFromReader);
      msg.setExit(value);
      break;
    case 3:
      var value = new proto.hashicorp.waypoint.EntrypointExecRequest.Output;
      reader.readMessage(value,proto.hashicorp.waypoint.EntrypointExecRequest.Output.deserializeBinaryFromReader);
      msg.setOutput(value);
      break;
    case 4:
      var value = new proto.hashicorp.waypoint.EntrypointExecRequest.Error;
      reader.readMessage(value,proto.hashicorp.waypoint.EntrypointExecRequest.Error.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointExecRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOpen();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.hashicorp.waypoint.EntrypointExecRequest.Open.serializeBinaryToWriter
    );
  }
  f = message.getExit();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.EntrypointExecRequest.Exit.serializeBinaryToWriter
    );
  }
  f = message.getOutput();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.hashicorp.waypoint.EntrypointExecRequest.Output.serializeBinaryToWriter
    );
  }
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.hashicorp.waypoint.EntrypointExecRequest.Error.serializeBinaryToWriter
    );
  }
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointExecRequest.Open.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Open} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.toObject = function(includeInstance, msg) {
  var f, obj = {
    instanceId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    index: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Open}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointExecRequest.Open;
  return proto.hashicorp.waypoint.EntrypointExecRequest.Open.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Open} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Open}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setInstanceId(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setIndex(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointExecRequest.Open.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Open} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getInstanceId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getIndex();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
};


/**
 * optional string instance_id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.prototype.getInstanceId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Open} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.prototype.setInstanceId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional int64 index = 2;
 * @return {number}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.prototype.getIndex = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Open} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Open.prototype.setIndex = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointExecRequest.Exit.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Exit} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.toObject = function(includeInstance, msg) {
  var f, obj = {
    code: jspb.Message.getFieldWithDefault(msg, 1, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Exit}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointExecRequest.Exit;
  return proto.hashicorp.waypoint.EntrypointExecRequest.Exit.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Exit} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Exit}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt32());
      msg.setCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointExecRequest.Exit.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Exit} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCode();
  if (f !== 0) {
    writer.writeInt32(
      1,
      f
    );
  }
};


/**
 * optional int32 code = 1;
 * @return {number}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.prototype.getCode = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Exit} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Exit.prototype.setCode = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointExecRequest.Output.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Output} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.toObject = function(includeInstance, msg) {
  var f, obj = {
    channel: jspb.Message.getFieldWithDefault(msg, 1, 0),
    data: msg.getData_asB64()
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Output}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointExecRequest.Output;
  return proto.hashicorp.waypoint.EntrypointExecRequest.Output.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Output} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Output}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.hashicorp.waypoint.EntrypointExecRequest.Output.Channel} */ (reader.readEnum());
      msg.setChannel(value);
      break;
    case 2:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setData(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointExecRequest.Output.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Output} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getChannel();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getData_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      2,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.Channel = {
  UNKNOWN: 0,
  STDOUT: 1,
  STDERR: 2
};

/**
 * optional Channel channel = 1;
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Output.Channel}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.getChannel = function() {
  return /** @type {!proto.hashicorp.waypoint.EntrypointExecRequest.Output.Channel} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Output.Channel} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Output} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.setChannel = function(value) {
  return jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional bytes data = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.getData = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * optional bytes data = 2;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.getData_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getData()));
};


/**
 * optional bytes data = 2;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getData()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Output} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Output.prototype.setData = function(value) {
  return jspb.Message.setProto3BytesField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointExecRequest.Error.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Error} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.toObject = function(includeInstance, msg) {
  var f, obj = {
    error: (f = msg.getError()) && google_rpc_status_pb.Status.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Error}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointExecRequest.Error;
  return proto.hashicorp.waypoint.EntrypointExecRequest.Error.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Error} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Error}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new google_rpc_status_pb.Status;
      reader.readMessage(value,google_rpc_status_pb.Status.deserializeBinaryFromReader);
      msg.setError(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointExecRequest.Error.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointExecRequest.Error} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getError();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      google_rpc_status_pb.Status.serializeBinaryToWriter
    );
  }
};


/**
 * optional google.rpc.Status error = 1;
 * @return {?proto.google.rpc.Status}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.prototype.getError = function() {
  return /** @type{?proto.google.rpc.Status} */ (
    jspb.Message.getWrapperField(this, google_rpc_status_pb.Status, 1));
};


/**
 * @param {?proto.google.rpc.Status|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Error} returns this
*/
proto.hashicorp.waypoint.EntrypointExecRequest.Error.prototype.setError = function(value) {
  return jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest.Error} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.Error.prototype.hasError = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Open open = 1;
 * @return {?proto.hashicorp.waypoint.EntrypointExecRequest.Open}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.getOpen = function() {
  return /** @type{?proto.hashicorp.waypoint.EntrypointExecRequest.Open} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.EntrypointExecRequest.Open, 1));
};


/**
 * @param {?proto.hashicorp.waypoint.EntrypointExecRequest.Open|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
*/
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.setOpen = function(value) {
  return jspb.Message.setOneofWrapperField(this, 1, proto.hashicorp.waypoint.EntrypointExecRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.clearOpen = function() {
  return this.setOpen(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.hasOpen = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional Exit exit = 2;
 * @return {?proto.hashicorp.waypoint.EntrypointExecRequest.Exit}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.getExit = function() {
  return /** @type{?proto.hashicorp.waypoint.EntrypointExecRequest.Exit} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.EntrypointExecRequest.Exit, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.EntrypointExecRequest.Exit|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
*/
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.setExit = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.EntrypointExecRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.clearExit = function() {
  return this.setExit(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.hasExit = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional Output output = 3;
 * @return {?proto.hashicorp.waypoint.EntrypointExecRequest.Output}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.getOutput = function() {
  return /** @type{?proto.hashicorp.waypoint.EntrypointExecRequest.Output} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.EntrypointExecRequest.Output, 3));
};


/**
 * @param {?proto.hashicorp.waypoint.EntrypointExecRequest.Output|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
*/
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.setOutput = function(value) {
  return jspb.Message.setOneofWrapperField(this, 3, proto.hashicorp.waypoint.EntrypointExecRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.clearOutput = function() {
  return this.setOutput(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.hasOutput = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional Error error = 4;
 * @return {?proto.hashicorp.waypoint.EntrypointExecRequest.Error}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.getError = function() {
  return /** @type{?proto.hashicorp.waypoint.EntrypointExecRequest.Error} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.EntrypointExecRequest.Error, 4));
};


/**
 * @param {?proto.hashicorp.waypoint.EntrypointExecRequest.Error|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
*/
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.setError = function(value) {
  return jspb.Message.setOneofWrapperField(this, 4, proto.hashicorp.waypoint.EntrypointExecRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecRequest} returns this
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.clearError = function() {
  return this.setError(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecRequest.prototype.hasError = function() {
  return jspb.Message.getField(this, 4) != null;
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_ = [[1,2,3]];

/**
 * @enum {number}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.EventCase = {
  EVENT_NOT_SET: 0,
  INPUT: 1,
  WINCH: 2,
  OPENED: 3
};

/**
 * @return {proto.hashicorp.waypoint.EntrypointExecResponse.EventCase}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.getEventCase = function() {
  return /** @type {proto.hashicorp.waypoint.EntrypointExecResponse.EventCase} */(jspb.Message.computeOneofCase(this, proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.EntrypointExecResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.EntrypointExecResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    input: msg.getInput_asB64(),
    winch: (f = msg.getWinch()) && proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.toObject(includeInstance, f),
    opened: jspb.Message.getBooleanFieldWithDefault(msg, 3, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.EntrypointExecResponse;
  return proto.hashicorp.waypoint.EntrypointExecResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.EntrypointExecResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setInput(value);
      break;
    case 2:
      var value = new proto.hashicorp.waypoint.ExecStreamRequest.WindowSize;
      reader.readMessage(value,proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.deserializeBinaryFromReader);
      msg.setWinch(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setOpened(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.EntrypointExecResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.EntrypointExecResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.EntrypointExecResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {!(string|Uint8Array)} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeBytes(
      1,
      f
    );
  }
  f = message.getWinch();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.hashicorp.waypoint.ExecStreamRequest.WindowSize.serializeBinaryToWriter
    );
  }
  f = /** @type {boolean} */ (jspb.Message.getField(message, 3));
  if (f != null) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional bytes input = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.getInput = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * optional bytes input = 1;
 * This is a type-conversion wrapper around `getInput()`
 * @return {string}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.getInput_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getInput()));
};


/**
 * optional bytes input = 1;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getInput()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.getInput_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getInput()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse} returns this
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.setInput = function(value) {
  return jspb.Message.setOneofField(this, 1, proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_[0], value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse} returns this
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.clearInput = function() {
  return jspb.Message.setOneofField(this, 1, proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.hasInput = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional ExecStreamRequest.WindowSize winch = 2;
 * @return {?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.getWinch = function() {
  return /** @type{?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize} */ (
    jspb.Message.getWrapperField(this, proto.hashicorp.waypoint.ExecStreamRequest.WindowSize, 2));
};


/**
 * @param {?proto.hashicorp.waypoint.ExecStreamRequest.WindowSize|undefined} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse} returns this
*/
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.setWinch = function(value) {
  return jspb.Message.setOneofWrapperField(this, 2, proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse} returns this
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.clearWinch = function() {
  return this.setWinch(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.hasWinch = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional bool opened = 3;
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.getOpened = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 3, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse} returns this
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.setOpened = function(value) {
  return jspb.Message.setOneofField(this, 3, proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_[0], value);
};


/**
 * Clears the field making it undefined.
 * @return {!proto.hashicorp.waypoint.EntrypointExecResponse} returns this
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.clearOpened = function() {
  return jspb.Message.setOneofField(this, 3, proto.hashicorp.waypoint.EntrypointExecResponse.oneofGroups_[0], undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.EntrypointExecResponse.prototype.hasOpened = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.TokenTransport.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.TokenTransport} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.TokenTransport.toObject = function(includeInstance, msg) {
  var f, obj = {
    body: msg.getBody_asB64(),
    signature: msg.getSignature_asB64(),
    keyId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    metadataMap: (f = msg.getMetadataMap()) ? f.toObject(includeInstance, undefined) : []
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.TokenTransport}
 */
proto.hashicorp.waypoint.TokenTransport.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.TokenTransport;
  return proto.hashicorp.waypoint.TokenTransport.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.TokenTransport} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.TokenTransport}
 */
proto.hashicorp.waypoint.TokenTransport.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setBody(value);
      break;
    case 2:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setSignature(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setKeyId(value);
      break;
    case 4:
      var value = msg.getMetadataMap();
      reader.readMessage(value, function(message, reader) {
        jspb.Map.deserializeBinary(message, reader, jspb.BinaryReader.prototype.readString, jspb.BinaryReader.prototype.readString, null, "", "");
         });
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.TokenTransport.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.TokenTransport} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.TokenTransport.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBody_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      1,
      f
    );
  }
  f = message.getSignature_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      2,
      f
    );
  }
  f = message.getKeyId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getMetadataMap(true);
  if (f && f.getLength() > 0) {
    f.serializeBinary(4, writer, jspb.BinaryWriter.prototype.writeString, jspb.BinaryWriter.prototype.writeString);
  }
};


/**
 * optional bytes body = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getBody = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * optional bytes body = 1;
 * This is a type-conversion wrapper around `getBody()`
 * @return {string}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getBody_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getBody()));
};


/**
 * optional bytes body = 1;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getBody()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getBody_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getBody()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.TokenTransport} returns this
 */
proto.hashicorp.waypoint.TokenTransport.prototype.setBody = function(value) {
  return jspb.Message.setProto3BytesField(this, 1, value);
};


/**
 * optional bytes signature = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getSignature = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * optional bytes signature = 2;
 * This is a type-conversion wrapper around `getSignature()`
 * @return {string}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getSignature_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getSignature()));
};


/**
 * optional bytes signature = 2;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getSignature()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getSignature_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getSignature()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.TokenTransport} returns this
 */
proto.hashicorp.waypoint.TokenTransport.prototype.setSignature = function(value) {
  return jspb.Message.setProto3BytesField(this, 2, value);
};


/**
 * optional string key_id = 3;
 * @return {string}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getKeyId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.TokenTransport} returns this
 */
proto.hashicorp.waypoint.TokenTransport.prototype.setKeyId = function(value) {
  return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * map<string, string> metadata = 4;
 * @param {boolean=} opt_noLazyCreate Do not create the map if
 * empty, instead returning `undefined`
 * @return {!jspb.Map<string,string>}
 */
proto.hashicorp.waypoint.TokenTransport.prototype.getMetadataMap = function(opt_noLazyCreate) {
  return /** @type {!jspb.Map<string,string>} */ (
      jspb.Message.getMapField(this, 4, opt_noLazyCreate,
      null));
};


/**
 * Clears values from the map. The map will be non-null.
 * @return {!proto.hashicorp.waypoint.TokenTransport} returns this
 */
proto.hashicorp.waypoint.TokenTransport.prototype.clearMetadataMap = function() {
  this.getMetadataMap().clear();
  return this;};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.Token.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.Token.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.Token} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Token.toObject = function(includeInstance, msg) {
  var f, obj = {
    user: jspb.Message.getFieldWithDefault(msg, 1, ""),
    tokenId: msg.getTokenId_asB64(),
    validUntil: (f = msg.getValidUntil()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    login: jspb.Message.getBooleanFieldWithDefault(msg, 4, false),
    invite: jspb.Message.getBooleanFieldWithDefault(msg, 5, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.Token}
 */
proto.hashicorp.waypoint.Token.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.Token;
  return proto.hashicorp.waypoint.Token.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.Token} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.Token}
 */
proto.hashicorp.waypoint.Token.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUser(value);
      break;
    case 2:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setTokenId(value);
      break;
    case 3:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setValidUntil(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setLogin(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setInvite(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Token.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.Token.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.Token} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.Token.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUser();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getTokenId_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      2,
      f
    );
  }
  f = message.getValidUntil();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getLogin();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getInvite();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
};


/**
 * optional string user = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.Token.prototype.getUser = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.Token} returns this
 */
proto.hashicorp.waypoint.Token.prototype.setUser = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bytes token_id = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.Token.prototype.getTokenId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * optional bytes token_id = 2;
 * This is a type-conversion wrapper around `getTokenId()`
 * @return {string}
 */
proto.hashicorp.waypoint.Token.prototype.getTokenId_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getTokenId()));
};


/**
 * optional bytes token_id = 2;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getTokenId()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.Token.prototype.getTokenId_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getTokenId()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.Token} returns this
 */
proto.hashicorp.waypoint.Token.prototype.setTokenId = function(value) {
  return jspb.Message.setProto3BytesField(this, 2, value);
};


/**
 * optional google.protobuf.Timestamp valid_until = 3;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.hashicorp.waypoint.Token.prototype.getValidUntil = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 3));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.hashicorp.waypoint.Token} returns this
*/
proto.hashicorp.waypoint.Token.prototype.setValidUntil = function(value) {
  return jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.hashicorp.waypoint.Token} returns this
 */
proto.hashicorp.waypoint.Token.prototype.clearValidUntil = function() {
  return this.setValidUntil(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.hashicorp.waypoint.Token.prototype.hasValidUntil = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional bool login = 4;
 * @return {boolean}
 */
proto.hashicorp.waypoint.Token.prototype.getLogin = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 4, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.Token} returns this
 */
proto.hashicorp.waypoint.Token.prototype.setLogin = function(value) {
  return jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * optional bool invite = 5;
 * @return {boolean}
 */
proto.hashicorp.waypoint.Token.prototype.getInvite = function() {
  return /** @type {boolean} */ (jspb.Message.getBooleanFieldWithDefault(this, 5, false));
};


/**
 * @param {boolean} value
 * @return {!proto.hashicorp.waypoint.Token} returns this
 */
proto.hashicorp.waypoint.Token.prototype.setInvite = function(value) {
  return jspb.Message.setProto3BooleanField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.HMACKey.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.HMACKey.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.HMACKey} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.HMACKey.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    key: msg.getKey_asB64()
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.HMACKey}
 */
proto.hashicorp.waypoint.HMACKey.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.HMACKey;
  return proto.hashicorp.waypoint.HMACKey.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.HMACKey} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.HMACKey}
 */
proto.hashicorp.waypoint.HMACKey.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setKey(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.HMACKey.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.HMACKey.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.HMACKey} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.HMACKey.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getKey_asU8();
  if (f.length > 0) {
    writer.writeBytes(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.HMACKey.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.HMACKey} returns this
 */
proto.hashicorp.waypoint.HMACKey.prototype.setId = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional bytes key = 2;
 * @return {string}
 */
proto.hashicorp.waypoint.HMACKey.prototype.getKey = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * optional bytes key = 2;
 * This is a type-conversion wrapper around `getKey()`
 * @return {string}
 */
proto.hashicorp.waypoint.HMACKey.prototype.getKey_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getKey()));
};


/**
 * optional bytes key = 2;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getKey()`
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.HMACKey.prototype.getKey_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getKey()));
};


/**
 * @param {!(string|Uint8Array)} value
 * @return {!proto.hashicorp.waypoint.HMACKey} returns this
 */
proto.hashicorp.waypoint.HMACKey.prototype.setKey = function(value) {
  return jspb.Message.setProto3BytesField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.InviteTokenRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.InviteTokenRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.InviteTokenRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.InviteTokenRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    duration: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.InviteTokenRequest}
 */
proto.hashicorp.waypoint.InviteTokenRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.InviteTokenRequest;
  return proto.hashicorp.waypoint.InviteTokenRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.InviteTokenRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.InviteTokenRequest}
 */
proto.hashicorp.waypoint.InviteTokenRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setDuration(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.InviteTokenRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.InviteTokenRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.InviteTokenRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.InviteTokenRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getDuration();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string duration = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.InviteTokenRequest.prototype.getDuration = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.InviteTokenRequest} returns this
 */
proto.hashicorp.waypoint.InviteTokenRequest.prototype.setDuration = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.NewTokenResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.NewTokenResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.NewTokenResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.NewTokenResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    token: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.NewTokenResponse}
 */
proto.hashicorp.waypoint.NewTokenResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.NewTokenResponse;
  return proto.hashicorp.waypoint.NewTokenResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.NewTokenResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.NewTokenResponse}
 */
proto.hashicorp.waypoint.NewTokenResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setToken(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.NewTokenResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.NewTokenResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.NewTokenResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.NewTokenResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getToken();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string token = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.NewTokenResponse.prototype.getToken = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.NewTokenResponse} returns this
 */
proto.hashicorp.waypoint.NewTokenResponse.prototype.setToken = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.hashicorp.waypoint.ConvertInviteTokenRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.hashicorp.waypoint.ConvertInviteTokenRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    token: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.hashicorp.waypoint.ConvertInviteTokenRequest}
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.hashicorp.waypoint.ConvertInviteTokenRequest;
  return proto.hashicorp.waypoint.ConvertInviteTokenRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.hashicorp.waypoint.ConvertInviteTokenRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.hashicorp.waypoint.ConvertInviteTokenRequest}
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setToken(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.hashicorp.waypoint.ConvertInviteTokenRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.hashicorp.waypoint.ConvertInviteTokenRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getToken();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string token = 1;
 * @return {string}
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.prototype.getToken = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.hashicorp.waypoint.ConvertInviteTokenRequest} returns this
 */
proto.hashicorp.waypoint.ConvertInviteTokenRequest.prototype.setToken = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


goog.object.extend(exports, proto.hashicorp.waypoint);
