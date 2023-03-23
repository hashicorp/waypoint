/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// package: google.api
// file: google/api/service.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../../google/api/annotations_pb";
import * as google_api_auth_pb from "../../google/api/auth_pb";
import * as google_api_backend_pb from "../../google/api/backend_pb";
import * as google_api_billing_pb from "../../google/api/billing_pb";
import * as google_api_context_pb from "../../google/api/context_pb";
import * as google_api_control_pb from "../../google/api/control_pb";
import * as google_api_documentation_pb from "../../google/api/documentation_pb";
import * as google_api_endpoint_pb from "../../google/api/endpoint_pb";
import * as google_api_http_pb from "../../google/api/http_pb";
import * as google_api_log_pb from "../../google/api/log_pb";
import * as google_api_logging_pb from "../../google/api/logging_pb";
import * as google_api_metric_pb from "../../google/api/metric_pb";
import * as google_api_monitored_resource_pb from "../../google/api/monitored_resource_pb";
import * as google_api_monitoring_pb from "../../google/api/monitoring_pb";
import * as google_api_quota_pb from "../../google/api/quota_pb";
import * as google_api_source_info_pb from "../../google/api/source_info_pb";
import * as google_api_system_parameter_pb from "../../google/api/system_parameter_pb";
import * as google_api_usage_pb from "../../google/api/usage_pb";
import * as google_protobuf_api_pb from "google-protobuf/google/protobuf/api_pb";
import * as google_protobuf_type_pb from "google-protobuf/google/protobuf/type_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class Service extends jspb.Message {
  hasConfigVersion(): boolean;
  clearConfigVersion(): void;
  getConfigVersion(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setConfigVersion(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  getName(): string;
  setName(value: string): void;

  getId(): string;
  setId(value: string): void;

  getTitle(): string;
  setTitle(value: string): void;

  getProducerProjectId(): string;
  setProducerProjectId(value: string): void;

  clearApisList(): void;
  getApisList(): Array<google_protobuf_api_pb.Api>;
  setApisList(value: Array<google_protobuf_api_pb.Api>): void;
  addApis(value?: google_protobuf_api_pb.Api, index?: number): google_protobuf_api_pb.Api;

  clearTypesList(): void;
  getTypesList(): Array<google_protobuf_type_pb.Type>;
  setTypesList(value: Array<google_protobuf_type_pb.Type>): void;
  addTypes(value?: google_protobuf_type_pb.Type, index?: number): google_protobuf_type_pb.Type;

  clearEnumsList(): void;
  getEnumsList(): Array<google_protobuf_type_pb.Enum>;
  setEnumsList(value: Array<google_protobuf_type_pb.Enum>): void;
  addEnums(value?: google_protobuf_type_pb.Enum, index?: number): google_protobuf_type_pb.Enum;

  hasDocumentation(): boolean;
  clearDocumentation(): void;
  getDocumentation(): google_api_documentation_pb.Documentation | undefined;
  setDocumentation(value?: google_api_documentation_pb.Documentation): void;

  hasBackend(): boolean;
  clearBackend(): void;
  getBackend(): google_api_backend_pb.Backend | undefined;
  setBackend(value?: google_api_backend_pb.Backend): void;

  hasHttp(): boolean;
  clearHttp(): void;
  getHttp(): google_api_http_pb.Http | undefined;
  setHttp(value?: google_api_http_pb.Http): void;

  hasQuota(): boolean;
  clearQuota(): void;
  getQuota(): google_api_quota_pb.Quota | undefined;
  setQuota(value?: google_api_quota_pb.Quota): void;

  hasAuthentication(): boolean;
  clearAuthentication(): void;
  getAuthentication(): google_api_auth_pb.Authentication | undefined;
  setAuthentication(value?: google_api_auth_pb.Authentication): void;

  hasContext(): boolean;
  clearContext(): void;
  getContext(): google_api_context_pb.Context | undefined;
  setContext(value?: google_api_context_pb.Context): void;

  hasUsage(): boolean;
  clearUsage(): void;
  getUsage(): google_api_usage_pb.Usage | undefined;
  setUsage(value?: google_api_usage_pb.Usage): void;

  clearEndpointsList(): void;
  getEndpointsList(): Array<google_api_endpoint_pb.Endpoint>;
  setEndpointsList(value: Array<google_api_endpoint_pb.Endpoint>): void;
  addEndpoints(value?: google_api_endpoint_pb.Endpoint, index?: number): google_api_endpoint_pb.Endpoint;

  hasControl(): boolean;
  clearControl(): void;
  getControl(): google_api_control_pb.Control | undefined;
  setControl(value?: google_api_control_pb.Control): void;

  clearLogsList(): void;
  getLogsList(): Array<google_api_log_pb.LogDescriptor>;
  setLogsList(value: Array<google_api_log_pb.LogDescriptor>): void;
  addLogs(value?: google_api_log_pb.LogDescriptor, index?: number): google_api_log_pb.LogDescriptor;

  clearMetricsList(): void;
  getMetricsList(): Array<google_api_metric_pb.MetricDescriptor>;
  setMetricsList(value: Array<google_api_metric_pb.MetricDescriptor>): void;
  addMetrics(value?: google_api_metric_pb.MetricDescriptor, index?: number): google_api_metric_pb.MetricDescriptor;

  clearMonitoredResourcesList(): void;
  getMonitoredResourcesList(): Array<google_api_monitored_resource_pb.MonitoredResourceDescriptor>;
  setMonitoredResourcesList(value: Array<google_api_monitored_resource_pb.MonitoredResourceDescriptor>): void;
  addMonitoredResources(value?: google_api_monitored_resource_pb.MonitoredResourceDescriptor, index?: number): google_api_monitored_resource_pb.MonitoredResourceDescriptor;

  hasBilling(): boolean;
  clearBilling(): void;
  getBilling(): google_api_billing_pb.Billing | undefined;
  setBilling(value?: google_api_billing_pb.Billing): void;

  hasLogging(): boolean;
  clearLogging(): void;
  getLogging(): google_api_logging_pb.Logging | undefined;
  setLogging(value?: google_api_logging_pb.Logging): void;

  hasMonitoring(): boolean;
  clearMonitoring(): void;
  getMonitoring(): google_api_monitoring_pb.Monitoring | undefined;
  setMonitoring(value?: google_api_monitoring_pb.Monitoring): void;

  hasSystemParameters(): boolean;
  clearSystemParameters(): void;
  getSystemParameters(): google_api_system_parameter_pb.SystemParameters | undefined;
  setSystemParameters(value?: google_api_system_parameter_pb.SystemParameters): void;

  hasSourceInfo(): boolean;
  clearSourceInfo(): void;
  getSourceInfo(): google_api_source_info_pb.SourceInfo | undefined;
  setSourceInfo(value?: google_api_source_info_pb.SourceInfo): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Service.AsObject;
  static toObject(includeInstance: boolean, msg: Service): Service.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Service, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Service;
  static deserializeBinaryFromReader(message: Service, reader: jspb.BinaryReader): Service;
}

export namespace Service {
  export type AsObject = {
    configVersion?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    name: string,
    id: string,
    title: string,
    producerProjectId: string,
    apisList: Array<google_protobuf_api_pb.Api.AsObject>,
    typesList: Array<google_protobuf_type_pb.Type.AsObject>,
    enumsList: Array<google_protobuf_type_pb.Enum.AsObject>,
    documentation?: google_api_documentation_pb.Documentation.AsObject,
    backend?: google_api_backend_pb.Backend.AsObject,
    http?: google_api_http_pb.Http.AsObject,
    quota?: google_api_quota_pb.Quota.AsObject,
    authentication?: google_api_auth_pb.Authentication.AsObject,
    context?: google_api_context_pb.Context.AsObject,
    usage?: google_api_usage_pb.Usage.AsObject,
    endpointsList: Array<google_api_endpoint_pb.Endpoint.AsObject>,
    control?: google_api_control_pb.Control.AsObject,
    logsList: Array<google_api_log_pb.LogDescriptor.AsObject>,
    metricsList: Array<google_api_metric_pb.MetricDescriptor.AsObject>,
    monitoredResourcesList: Array<google_api_monitored_resource_pb.MonitoredResourceDescriptor.AsObject>,
    billing?: google_api_billing_pb.Billing.AsObject,
    logging?: google_api_logging_pb.Logging.AsObject,
    monitoring?: google_api_monitoring_pb.Monitoring.AsObject,
    systemParameters?: google_api_system_parameter_pb.SystemParameters.AsObject,
    sourceInfo?: google_api_source_info_pb.SourceInfo.AsObject,
  }
}

