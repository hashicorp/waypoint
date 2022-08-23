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
import * as pkg_server_proto_server_pb from 'waypoint-pb';


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

  methodDescriptorGetVersionInfo = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetVersionInfo',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.GetVersionInfoResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetVersionInfoResponse.deserializeBinary
  );

  getVersionInfo(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetVersionInfoResponse>;

  getVersionInfo(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetVersionInfoResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetVersionInfoResponse>;

  getVersionInfo(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetVersionInfoResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetVersionInfo',
        request,
        metadata || {},
        this.methodDescriptorGetVersionInfo,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetVersionInfo',
    request,
    metadata || {},
    this.methodDescriptorGetVersionInfo);
  }

  methodDescriptorListOIDCAuthMethods = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListOIDCAuthMethods',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.ListOIDCAuthMethodsResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListOIDCAuthMethodsResponse.deserializeBinary
  );

  listOIDCAuthMethods(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListOIDCAuthMethodsResponse>;

  listOIDCAuthMethods(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListOIDCAuthMethodsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListOIDCAuthMethodsResponse>;

  listOIDCAuthMethods(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListOIDCAuthMethodsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListOIDCAuthMethods',
        request,
        metadata || {},
        this.methodDescriptorListOIDCAuthMethods,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListOIDCAuthMethods',
    request,
    metadata || {},
    this.methodDescriptorListOIDCAuthMethods);
  }

  methodDescriptorGetOIDCAuthURL = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetOIDCAuthURL',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetOIDCAuthURLRequest,
    pkg_server_proto_server_pb.GetOIDCAuthURLResponse,
    (request: pkg_server_proto_server_pb.GetOIDCAuthURLRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetOIDCAuthURLResponse.deserializeBinary
  );

  getOIDCAuthURL(
    request: pkg_server_proto_server_pb.GetOIDCAuthURLRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetOIDCAuthURLResponse>;

  getOIDCAuthURL(
    request: pkg_server_proto_server_pb.GetOIDCAuthURLRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetOIDCAuthURLResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetOIDCAuthURLResponse>;

  getOIDCAuthURL(
    request: pkg_server_proto_server_pb.GetOIDCAuthURLRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetOIDCAuthURLResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetOIDCAuthURL',
        request,
        metadata || {},
        this.methodDescriptorGetOIDCAuthURL,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetOIDCAuthURL',
    request,
    metadata || {},
    this.methodDescriptorGetOIDCAuthURL);
  }

  methodDescriptorCompleteOIDCAuth = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/CompleteOIDCAuth',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.CompleteOIDCAuthRequest,
    pkg_server_proto_server_pb.CompleteOIDCAuthResponse,
    (request: pkg_server_proto_server_pb.CompleteOIDCAuthRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.CompleteOIDCAuthResponse.deserializeBinary
  );

  completeOIDCAuth(
    request: pkg_server_proto_server_pb.CompleteOIDCAuthRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.CompleteOIDCAuthResponse>;

  completeOIDCAuth(
    request: pkg_server_proto_server_pb.CompleteOIDCAuthRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.CompleteOIDCAuthResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.CompleteOIDCAuthResponse>;

  completeOIDCAuth(
    request: pkg_server_proto_server_pb.CompleteOIDCAuthRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.CompleteOIDCAuthResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/CompleteOIDCAuth',
        request,
        metadata || {},
        this.methodDescriptorCompleteOIDCAuth,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/CompleteOIDCAuth',
    request,
    metadata || {},
    this.methodDescriptorCompleteOIDCAuth);
  }

  methodDescriptorNoAuthRunTrigger = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/NoAuthRunTrigger',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.RunTriggerRequest,
    pkg_server_proto_server_pb.RunTriggerResponse,
    (request: pkg_server_proto_server_pb.RunTriggerRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.RunTriggerResponse.deserializeBinary
  );

  noAuthRunTrigger(
    request: pkg_server_proto_server_pb.RunTriggerRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.RunTriggerResponse>;

  noAuthRunTrigger(
    request: pkg_server_proto_server_pb.RunTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunTriggerResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.RunTriggerResponse>;

  noAuthRunTrigger(
    request: pkg_server_proto_server_pb.RunTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunTriggerResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/NoAuthRunTrigger',
        request,
        metadata || {},
        this.methodDescriptorNoAuthRunTrigger,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/NoAuthRunTrigger',
    request,
    metadata || {},
    this.methodDescriptorNoAuthRunTrigger);
  }

  methodDescriptorGetUser = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetUser',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetUserRequest,
    pkg_server_proto_server_pb.GetUserResponse,
    (request: pkg_server_proto_server_pb.GetUserRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetUserResponse.deserializeBinary
  );

  getUser(
    request: pkg_server_proto_server_pb.GetUserRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetUserResponse>;

  getUser(
    request: pkg_server_proto_server_pb.GetUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetUserResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetUserResponse>;

  getUser(
    request: pkg_server_proto_server_pb.GetUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetUserResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetUser',
        request,
        metadata || {},
        this.methodDescriptorGetUser,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetUser',
    request,
    metadata || {},
    this.methodDescriptorGetUser);
  }

  methodDescriptorListUsers = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListUsers',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.ListUsersResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListUsersResponse.deserializeBinary
  );

  listUsers(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListUsersResponse>;

  listUsers(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListUsersResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListUsersResponse>;

  listUsers(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListUsersResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListUsers',
        request,
        metadata || {},
        this.methodDescriptorListUsers,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListUsers',
    request,
    metadata || {},
    this.methodDescriptorListUsers);
  }

  methodDescriptorUpdateUser = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpdateUser',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpdateUserRequest,
    pkg_server_proto_server_pb.UpdateUserResponse,
    (request: pkg_server_proto_server_pb.UpdateUserRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpdateUserResponse.deserializeBinary
  );

  updateUser(
    request: pkg_server_proto_server_pb.UpdateUserRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpdateUserResponse>;

  updateUser(
    request: pkg_server_proto_server_pb.UpdateUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpdateUserResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpdateUserResponse>;

  updateUser(
    request: pkg_server_proto_server_pb.UpdateUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpdateUserResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpdateUser',
        request,
        metadata || {},
        this.methodDescriptorUpdateUser,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpdateUser',
    request,
    metadata || {},
    this.methodDescriptorUpdateUser);
  }

  methodDescriptorDeleteUser = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/DeleteUser',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.DeleteUserRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.DeleteUserRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  deleteUser(
    request: pkg_server_proto_server_pb.DeleteUserRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  deleteUser(
    request: pkg_server_proto_server_pb.DeleteUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  deleteUser(
    request: pkg_server_proto_server_pb.DeleteUserRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/DeleteUser',
        request,
        metadata || {},
        this.methodDescriptorDeleteUser,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/DeleteUser',
    request,
    metadata || {},
    this.methodDescriptorDeleteUser);
  }

  methodDescriptorUpsertAuthMethod = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertAuthMethod',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertAuthMethodRequest,
    pkg_server_proto_server_pb.UpsertAuthMethodResponse,
    (request: pkg_server_proto_server_pb.UpsertAuthMethodRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertAuthMethodResponse.deserializeBinary
  );

  upsertAuthMethod(
    request: pkg_server_proto_server_pb.UpsertAuthMethodRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertAuthMethodResponse>;

  upsertAuthMethod(
    request: pkg_server_proto_server_pb.UpsertAuthMethodRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertAuthMethodResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertAuthMethodResponse>;

  upsertAuthMethod(
    request: pkg_server_proto_server_pb.UpsertAuthMethodRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertAuthMethodResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertAuthMethod',
        request,
        metadata || {},
        this.methodDescriptorUpsertAuthMethod,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertAuthMethod',
    request,
    metadata || {},
    this.methodDescriptorUpsertAuthMethod);
  }

  methodDescriptorGetAuthMethod = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetAuthMethod',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetAuthMethodRequest,
    pkg_server_proto_server_pb.GetAuthMethodResponse,
    (request: pkg_server_proto_server_pb.GetAuthMethodRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetAuthMethodResponse.deserializeBinary
  );

  getAuthMethod(
    request: pkg_server_proto_server_pb.GetAuthMethodRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetAuthMethodResponse>;

  getAuthMethod(
    request: pkg_server_proto_server_pb.GetAuthMethodRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetAuthMethodResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetAuthMethodResponse>;

  getAuthMethod(
    request: pkg_server_proto_server_pb.GetAuthMethodRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetAuthMethodResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetAuthMethod',
        request,
        metadata || {},
        this.methodDescriptorGetAuthMethod,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetAuthMethod',
    request,
    metadata || {},
    this.methodDescriptorGetAuthMethod);
  }

  methodDescriptorListAuthMethods = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListAuthMethods',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.ListAuthMethodsResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListAuthMethodsResponse.deserializeBinary
  );

  listAuthMethods(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListAuthMethodsResponse>;

  listAuthMethods(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListAuthMethodsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListAuthMethodsResponse>;

  listAuthMethods(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListAuthMethodsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListAuthMethods',
        request,
        metadata || {},
        this.methodDescriptorListAuthMethods,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListAuthMethods',
    request,
    metadata || {},
    this.methodDescriptorListAuthMethods);
  }

  methodDescriptorDeleteAuthMethod = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/DeleteAuthMethod',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.DeleteAuthMethodRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.DeleteAuthMethodRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  deleteAuthMethod(
    request: pkg_server_proto_server_pb.DeleteAuthMethodRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  deleteAuthMethod(
    request: pkg_server_proto_server_pb.DeleteAuthMethodRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  deleteAuthMethod(
    request: pkg_server_proto_server_pb.DeleteAuthMethodRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/DeleteAuthMethod',
        request,
        metadata || {},
        this.methodDescriptorDeleteAuthMethod,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/DeleteAuthMethod',
    request,
    metadata || {},
    this.methodDescriptorDeleteAuthMethod);
  }

  methodDescriptorListWorkspaces = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListWorkspaces',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListWorkspacesRequest,
    pkg_server_proto_server_pb.ListWorkspacesResponse,
    (request: pkg_server_proto_server_pb.ListWorkspacesRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListWorkspacesResponse.deserializeBinary
  );

  listWorkspaces(
    request: pkg_server_proto_server_pb.ListWorkspacesRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListWorkspacesResponse>;

  listWorkspaces(
    request: pkg_server_proto_server_pb.ListWorkspacesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListWorkspacesResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListWorkspacesResponse>;

  listWorkspaces(
    request: pkg_server_proto_server_pb.ListWorkspacesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListWorkspacesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListWorkspaces',
        request,
        metadata || {},
        this.methodDescriptorListWorkspaces,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListWorkspaces',
    request,
    metadata || {},
    this.methodDescriptorListWorkspaces);
  }

  methodDescriptorGetWorkspace = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetWorkspace',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetWorkspaceRequest,
    pkg_server_proto_server_pb.GetWorkspaceResponse,
    (request: pkg_server_proto_server_pb.GetWorkspaceRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetWorkspaceResponse.deserializeBinary
  );

  getWorkspace(
    request: pkg_server_proto_server_pb.GetWorkspaceRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetWorkspaceResponse>;

  getWorkspace(
    request: pkg_server_proto_server_pb.GetWorkspaceRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetWorkspaceResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetWorkspaceResponse>;

  getWorkspace(
    request: pkg_server_proto_server_pb.GetWorkspaceRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetWorkspaceResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetWorkspace',
        request,
        metadata || {},
        this.methodDescriptorGetWorkspace,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetWorkspace',
    request,
    metadata || {},
    this.methodDescriptorGetWorkspace);
  }

  methodDescriptorUpsertWorkspace = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertWorkspace',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertWorkspaceRequest,
    pkg_server_proto_server_pb.UpsertWorkspaceResponse,
    (request: pkg_server_proto_server_pb.UpsertWorkspaceRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertWorkspaceResponse.deserializeBinary
  );

  upsertWorkspace(
    request: pkg_server_proto_server_pb.UpsertWorkspaceRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertWorkspaceResponse>;

  upsertWorkspace(
    request: pkg_server_proto_server_pb.UpsertWorkspaceRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertWorkspaceResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertWorkspaceResponse>;

  upsertWorkspace(
    request: pkg_server_proto_server_pb.UpsertWorkspaceRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertWorkspaceResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertWorkspace',
        request,
        metadata || {},
        this.methodDescriptorUpsertWorkspace,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertWorkspace',
    request,
    metadata || {},
    this.methodDescriptorUpsertWorkspace);
  }

  methodDescriptorUpsertProject = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertProject',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertProjectRequest,
    pkg_server_proto_server_pb.UpsertProjectResponse,
    (request: pkg_server_proto_server_pb.UpsertProjectRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertProjectResponse.deserializeBinary
  );

  upsertProject(
    request: pkg_server_proto_server_pb.UpsertProjectRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertProjectResponse>;

  upsertProject(
    request: pkg_server_proto_server_pb.UpsertProjectRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertProjectResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertProjectResponse>;

  upsertProject(
    request: pkg_server_proto_server_pb.UpsertProjectRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertProjectResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertProject',
        request,
        metadata || {},
        this.methodDescriptorUpsertProject,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertProject',
    request,
    metadata || {},
    this.methodDescriptorUpsertProject);
  }

  methodDescriptorGetProject = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetProject',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetProjectRequest,
    pkg_server_proto_server_pb.GetProjectResponse,
    (request: pkg_server_proto_server_pb.GetProjectRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetProjectResponse.deserializeBinary
  );

  getProject(
    request: pkg_server_proto_server_pb.GetProjectRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetProjectResponse>;

  getProject(
    request: pkg_server_proto_server_pb.GetProjectRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetProjectResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetProjectResponse>;

  getProject(
    request: pkg_server_proto_server_pb.GetProjectRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetProjectResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetProject',
        request,
        metadata || {},
        this.methodDescriptorGetProject,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetProject',
    request,
    metadata || {},
    this.methodDescriptorGetProject);
  }

  methodDescriptorListProjects = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListProjects',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.ListProjectsResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListProjectsResponse.deserializeBinary
  );

  listProjects(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListProjectsResponse>;

  listProjects(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListProjectsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListProjectsResponse>;

  listProjects(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListProjectsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListProjects',
        request,
        metadata || {},
        this.methodDescriptorListProjects,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListProjects',
    request,
    metadata || {},
    this.methodDescriptorListProjects);
  }

  methodDescriptorGetApplication = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetApplication',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetApplicationRequest,
    pkg_server_proto_server_pb.GetApplicationResponse,
    (request: pkg_server_proto_server_pb.GetApplicationRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetApplicationResponse.deserializeBinary
  );

  getApplication(
    request: pkg_server_proto_server_pb.GetApplicationRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetApplicationResponse>;

  getApplication(
    request: pkg_server_proto_server_pb.GetApplicationRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetApplicationResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetApplicationResponse>;

  getApplication(
    request: pkg_server_proto_server_pb.GetApplicationRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetApplicationResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetApplication',
        request,
        metadata || {},
        this.methodDescriptorGetApplication,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetApplication',
    request,
    metadata || {},
    this.methodDescriptorGetApplication);
  }

  methodDescriptorUpsertApplication = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertApplication',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertApplicationRequest,
    pkg_server_proto_server_pb.UpsertApplicationResponse,
    (request: pkg_server_proto_server_pb.UpsertApplicationRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertApplicationResponse.deserializeBinary
  );

  upsertApplication(
    request: pkg_server_proto_server_pb.UpsertApplicationRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertApplicationResponse>;

  upsertApplication(
    request: pkg_server_proto_server_pb.UpsertApplicationRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertApplicationResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertApplicationResponse>;

  upsertApplication(
    request: pkg_server_proto_server_pb.UpsertApplicationRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertApplicationResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertApplication',
        request,
        metadata || {},
        this.methodDescriptorUpsertApplication,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertApplication',
    request,
    metadata || {},
    this.methodDescriptorUpsertApplication);
  }

  methodDescriptorListBuilds = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListBuilds',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListBuildsRequest,
    pkg_server_proto_server_pb.ListBuildsResponse,
    (request: pkg_server_proto_server_pb.ListBuildsRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListBuildsResponse.deserializeBinary
  );

  listBuilds(
    request: pkg_server_proto_server_pb.ListBuildsRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListBuildsResponse>;

  listBuilds(
    request: pkg_server_proto_server_pb.ListBuildsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListBuildsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListBuildsResponse>;

  listBuilds(
    request: pkg_server_proto_server_pb.ListBuildsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListBuildsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListBuilds',
        request,
        metadata || {},
        this.methodDescriptorListBuilds,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListBuilds',
    request,
    metadata || {},
    this.methodDescriptorListBuilds);
  }

  methodDescriptorGetBuild = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetBuild',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetBuildRequest,
    pkg_server_proto_server_pb.Build,
    (request: pkg_server_proto_server_pb.GetBuildRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.Build.deserializeBinary
  );

  getBuild(
    request: pkg_server_proto_server_pb.GetBuildRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.Build>;

  getBuild(
    request: pkg_server_proto_server_pb.GetBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Build) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.Build>;

  getBuild(
    request: pkg_server_proto_server_pb.GetBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Build) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetBuild',
        request,
        metadata || {},
        this.methodDescriptorGetBuild,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetBuild',
    request,
    metadata || {},
    this.methodDescriptorGetBuild);
  }

  methodDescriptorGetLatestBuild = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetLatestBuild',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetLatestBuildRequest,
    pkg_server_proto_server_pb.Build,
    (request: pkg_server_proto_server_pb.GetLatestBuildRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.Build.deserializeBinary
  );

  getLatestBuild(
    request: pkg_server_proto_server_pb.GetLatestBuildRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.Build>;

  getLatestBuild(
    request: pkg_server_proto_server_pb.GetLatestBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Build) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.Build>;

  getLatestBuild(
    request: pkg_server_proto_server_pb.GetLatestBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Build) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetLatestBuild',
        request,
        metadata || {},
        this.methodDescriptorGetLatestBuild,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetLatestBuild',
    request,
    metadata || {},
    this.methodDescriptorGetLatestBuild);
  }

  methodDescriptorListPushedArtifacts = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListPushedArtifacts',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListPushedArtifactsRequest,
    pkg_server_proto_server_pb.ListPushedArtifactsResponse,
    (request: pkg_server_proto_server_pb.ListPushedArtifactsRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListPushedArtifactsResponse.deserializeBinary
  );

  listPushedArtifacts(
    request: pkg_server_proto_server_pb.ListPushedArtifactsRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListPushedArtifactsResponse>;

  listPushedArtifacts(
    request: pkg_server_proto_server_pb.ListPushedArtifactsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListPushedArtifactsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListPushedArtifactsResponse>;

  listPushedArtifacts(
    request: pkg_server_proto_server_pb.ListPushedArtifactsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListPushedArtifactsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListPushedArtifacts',
        request,
        metadata || {},
        this.methodDescriptorListPushedArtifacts,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListPushedArtifacts',
    request,
    metadata || {},
    this.methodDescriptorListPushedArtifacts);
  }

  methodDescriptorGetPushedArtifact = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetPushedArtifact',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetPushedArtifactRequest,
    pkg_server_proto_server_pb.PushedArtifact,
    (request: pkg_server_proto_server_pb.GetPushedArtifactRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.PushedArtifact.deserializeBinary
  );

  getPushedArtifact(
    request: pkg_server_proto_server_pb.GetPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.PushedArtifact>;

  getPushedArtifact(
    request: pkg_server_proto_server_pb.GetPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.PushedArtifact) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.PushedArtifact>;

  getPushedArtifact(
    request: pkg_server_proto_server_pb.GetPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.PushedArtifact) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetPushedArtifact',
        request,
        metadata || {},
        this.methodDescriptorGetPushedArtifact,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetPushedArtifact',
    request,
    metadata || {},
    this.methodDescriptorGetPushedArtifact);
  }

  methodDescriptorGetLatestPushedArtifact = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetLatestPushedArtifact',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetLatestPushedArtifactRequest,
    pkg_server_proto_server_pb.PushedArtifact,
    (request: pkg_server_proto_server_pb.GetLatestPushedArtifactRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.PushedArtifact.deserializeBinary
  );

  getLatestPushedArtifact(
    request: pkg_server_proto_server_pb.GetLatestPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.PushedArtifact>;

  getLatestPushedArtifact(
    request: pkg_server_proto_server_pb.GetLatestPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.PushedArtifact) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.PushedArtifact>;

  getLatestPushedArtifact(
    request: pkg_server_proto_server_pb.GetLatestPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.PushedArtifact) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetLatestPushedArtifact',
        request,
        metadata || {},
        this.methodDescriptorGetLatestPushedArtifact,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetLatestPushedArtifact',
    request,
    metadata || {},
    this.methodDescriptorGetLatestPushedArtifact);
  }

  methodDescriptorListDeployments = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListDeployments',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListDeploymentsRequest,
    pkg_server_proto_server_pb.ListDeploymentsResponse,
    (request: pkg_server_proto_server_pb.ListDeploymentsRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListDeploymentsResponse.deserializeBinary
  );

  listDeployments(
    request: pkg_server_proto_server_pb.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListDeploymentsResponse>;

  listDeployments(
    request: pkg_server_proto_server_pb.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListDeploymentsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListDeploymentsResponse>;

  listDeployments(
    request: pkg_server_proto_server_pb.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListDeploymentsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListDeployments',
        request,
        metadata || {},
        this.methodDescriptorListDeployments,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListDeployments',
    request,
    metadata || {},
    this.methodDescriptorListDeployments);
  }

  methodDescriptorGetDeployment = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetDeployment',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetDeploymentRequest,
    pkg_server_proto_server_pb.Deployment,
    (request: pkg_server_proto_server_pb.GetDeploymentRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.Deployment.deserializeBinary
  );

  getDeployment(
    request: pkg_server_proto_server_pb.GetDeploymentRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.Deployment>;

  getDeployment(
    request: pkg_server_proto_server_pb.GetDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Deployment) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.Deployment>;

  getDeployment(
    request: pkg_server_proto_server_pb.GetDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Deployment) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetDeployment',
        request,
        metadata || {},
        this.methodDescriptorGetDeployment,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetDeployment',
    request,
    metadata || {},
    this.methodDescriptorGetDeployment);
  }

  methodDescriptorListInstances = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListInstances',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListInstancesRequest,
    pkg_server_proto_server_pb.ListInstancesResponse,
    (request: pkg_server_proto_server_pb.ListInstancesRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListInstancesResponse.deserializeBinary
  );

  listInstances(
    request: pkg_server_proto_server_pb.ListInstancesRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListInstancesResponse>;

  listInstances(
    request: pkg_server_proto_server_pb.ListInstancesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListInstancesResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListInstancesResponse>;

  listInstances(
    request: pkg_server_proto_server_pb.ListInstancesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListInstancesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListInstances',
        request,
        metadata || {},
        this.methodDescriptorListInstances,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListInstances',
    request,
    metadata || {},
    this.methodDescriptorListInstances);
  }

  methodDescriptorListReleases = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListReleases',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListReleasesRequest,
    pkg_server_proto_server_pb.ListReleasesResponse,
    (request: pkg_server_proto_server_pb.ListReleasesRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListReleasesResponse.deserializeBinary
  );

  listReleases(
    request: pkg_server_proto_server_pb.ListReleasesRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListReleasesResponse>;

  listReleases(
    request: pkg_server_proto_server_pb.ListReleasesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListReleasesResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListReleasesResponse>;

  listReleases(
    request: pkg_server_proto_server_pb.ListReleasesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListReleasesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListReleases',
        request,
        metadata || {},
        this.methodDescriptorListReleases,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListReleases',
    request,
    metadata || {},
    this.methodDescriptorListReleases);
  }

  methodDescriptorGetRelease = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetRelease',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetReleaseRequest,
    pkg_server_proto_server_pb.Release,
    (request: pkg_server_proto_server_pb.GetReleaseRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.Release.deserializeBinary
  );

  getRelease(
    request: pkg_server_proto_server_pb.GetReleaseRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.Release>;

  getRelease(
    request: pkg_server_proto_server_pb.GetReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Release) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.Release>;

  getRelease(
    request: pkg_server_proto_server_pb.GetReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Release) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetRelease',
        request,
        metadata || {},
        this.methodDescriptorGetRelease,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetRelease',
    request,
    metadata || {},
    this.methodDescriptorGetRelease);
  }

  methodDescriptorGetLatestRelease = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetLatestRelease',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetLatestReleaseRequest,
    pkg_server_proto_server_pb.Release,
    (request: pkg_server_proto_server_pb.GetLatestReleaseRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.Release.deserializeBinary
  );

  getLatestRelease(
    request: pkg_server_proto_server_pb.GetLatestReleaseRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.Release>;

  getLatestRelease(
    request: pkg_server_proto_server_pb.GetLatestReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Release) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.Release>;

  getLatestRelease(
    request: pkg_server_proto_server_pb.GetLatestReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Release) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetLatestRelease',
        request,
        metadata || {},
        this.methodDescriptorGetLatestRelease,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetLatestRelease',
    request,
    metadata || {},
    this.methodDescriptorGetLatestRelease);
  }

  methodDescriptorGetStatusReport = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetStatusReport',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetStatusReportRequest,
    pkg_server_proto_server_pb.StatusReport,
    (request: pkg_server_proto_server_pb.GetStatusReportRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.StatusReport.deserializeBinary
  );

  getStatusReport(
    request: pkg_server_proto_server_pb.GetStatusReportRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.StatusReport>;

  getStatusReport(
    request: pkg_server_proto_server_pb.GetStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.StatusReport) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.StatusReport>;

  getStatusReport(
    request: pkg_server_proto_server_pb.GetStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.StatusReport) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetStatusReport',
        request,
        metadata || {},
        this.methodDescriptorGetStatusReport,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetStatusReport',
    request,
    metadata || {},
    this.methodDescriptorGetStatusReport);
  }

  methodDescriptorGetLatestStatusReport = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetLatestStatusReport',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetLatestStatusReportRequest,
    pkg_server_proto_server_pb.StatusReport,
    (request: pkg_server_proto_server_pb.GetLatestStatusReportRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.StatusReport.deserializeBinary
  );

  getLatestStatusReport(
    request: pkg_server_proto_server_pb.GetLatestStatusReportRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.StatusReport>;

  getLatestStatusReport(
    request: pkg_server_proto_server_pb.GetLatestStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.StatusReport) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.StatusReport>;

  getLatestStatusReport(
    request: pkg_server_proto_server_pb.GetLatestStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.StatusReport) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetLatestStatusReport',
        request,
        metadata || {},
        this.methodDescriptorGetLatestStatusReport,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetLatestStatusReport',
    request,
    metadata || {},
    this.methodDescriptorGetLatestStatusReport);
  }

  methodDescriptorListStatusReports = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListStatusReports',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListStatusReportsRequest,
    pkg_server_proto_server_pb.ListStatusReportsResponse,
    (request: pkg_server_proto_server_pb.ListStatusReportsRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListStatusReportsResponse.deserializeBinary
  );

  listStatusReports(
    request: pkg_server_proto_server_pb.ListStatusReportsRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListStatusReportsResponse>;

  listStatusReports(
    request: pkg_server_proto_server_pb.ListStatusReportsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListStatusReportsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListStatusReportsResponse>;

  listStatusReports(
    request: pkg_server_proto_server_pb.ListStatusReportsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListStatusReportsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListStatusReports',
        request,
        metadata || {},
        this.methodDescriptorListStatusReports,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListStatusReports',
    request,
    metadata || {},
    this.methodDescriptorListStatusReports);
  }

  methodDescriptorExpediteStatusReport = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ExpediteStatusReport',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ExpediteStatusReportRequest,
    pkg_server_proto_server_pb.ExpediteStatusReportResponse,
    (request: pkg_server_proto_server_pb.ExpediteStatusReportRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ExpediteStatusReportResponse.deserializeBinary
  );

  expediteStatusReport(
    request: pkg_server_proto_server_pb.ExpediteStatusReportRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ExpediteStatusReportResponse>;

  expediteStatusReport(
    request: pkg_server_proto_server_pb.ExpediteStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ExpediteStatusReportResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ExpediteStatusReportResponse>;

  expediteStatusReport(
    request: pkg_server_proto_server_pb.ExpediteStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ExpediteStatusReportResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ExpediteStatusReport',
        request,
        metadata || {},
        this.methodDescriptorExpediteStatusReport,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ExpediteStatusReport',
    request,
    metadata || {},
    this.methodDescriptorExpediteStatusReport);
  }

  methodDescriptorGetLogStream = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetLogStream',
    grpcWeb.MethodType.SERVER_STREAMING,
    pkg_server_proto_server_pb.GetLogStreamRequest,
    pkg_server_proto_server_pb.LogBatch,
    (request: pkg_server_proto_server_pb.GetLogStreamRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.LogBatch.deserializeBinary
  );

  getLogStream(
    request: pkg_server_proto_server_pb.GetLogStreamRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.LogBatch> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/hashicorp.waypoint.Waypoint/GetLogStream',
      request,
      metadata || {},
      this.methodDescriptorGetLogStream);
  }

  methodDescriptorSetConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/SetConfig',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ConfigSetRequest,
    pkg_server_proto_server_pb.ConfigSetResponse,
    (request: pkg_server_proto_server_pb.ConfigSetRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ConfigSetResponse.deserializeBinary
  );

  setConfig(
    request: pkg_server_proto_server_pb.ConfigSetRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ConfigSetResponse>;

  setConfig(
    request: pkg_server_proto_server_pb.ConfigSetRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ConfigSetResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ConfigSetResponse>;

  setConfig(
    request: pkg_server_proto_server_pb.ConfigSetRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ConfigSetResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/SetConfig',
        request,
        metadata || {},
        this.methodDescriptorSetConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/SetConfig',
    request,
    metadata || {},
    this.methodDescriptorSetConfig);
  }

  methodDescriptorGetConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetConfig',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ConfigGetRequest,
    pkg_server_proto_server_pb.ConfigGetResponse,
    (request: pkg_server_proto_server_pb.ConfigGetRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ConfigGetResponse.deserializeBinary
  );

  getConfig(
    request: pkg_server_proto_server_pb.ConfigGetRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ConfigGetResponse>;

  getConfig(
    request: pkg_server_proto_server_pb.ConfigGetRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ConfigGetResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ConfigGetResponse>;

  getConfig(
    request: pkg_server_proto_server_pb.ConfigGetRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ConfigGetResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetConfig',
        request,
        metadata || {},
        this.methodDescriptorGetConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetConfig',
    request,
    metadata || {},
    this.methodDescriptorGetConfig);
  }

  methodDescriptorSetConfigSource = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/SetConfigSource',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.SetConfigSourceRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.SetConfigSourceRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  setConfigSource(
    request: pkg_server_proto_server_pb.SetConfigSourceRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  setConfigSource(
    request: pkg_server_proto_server_pb.SetConfigSourceRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  setConfigSource(
    request: pkg_server_proto_server_pb.SetConfigSourceRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/SetConfigSource',
        request,
        metadata || {},
        this.methodDescriptorSetConfigSource,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/SetConfigSource',
    request,
    metadata || {},
    this.methodDescriptorSetConfigSource);
  }

  methodDescriptorGetConfigSource = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetConfigSource',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetConfigSourceRequest,
    pkg_server_proto_server_pb.GetConfigSourceResponse,
    (request: pkg_server_proto_server_pb.GetConfigSourceRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetConfigSourceResponse.deserializeBinary
  );

  getConfigSource(
    request: pkg_server_proto_server_pb.GetConfigSourceRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetConfigSourceResponse>;

  getConfigSource(
    request: pkg_server_proto_server_pb.GetConfigSourceRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetConfigSourceResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetConfigSourceResponse>;

  getConfigSource(
    request: pkg_server_proto_server_pb.GetConfigSourceRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetConfigSourceResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetConfigSource',
        request,
        metadata || {},
        this.methodDescriptorGetConfigSource,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetConfigSource',
    request,
    metadata || {},
    this.methodDescriptorGetConfigSource);
  }

  methodDescriptorCreateHostname = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/CreateHostname',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.CreateHostnameRequest,
    pkg_server_proto_server_pb.CreateHostnameResponse,
    (request: pkg_server_proto_server_pb.CreateHostnameRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.CreateHostnameResponse.deserializeBinary
  );

  createHostname(
    request: pkg_server_proto_server_pb.CreateHostnameRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.CreateHostnameResponse>;

  createHostname(
    request: pkg_server_proto_server_pb.CreateHostnameRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.CreateHostnameResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.CreateHostnameResponse>;

  createHostname(
    request: pkg_server_proto_server_pb.CreateHostnameRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.CreateHostnameResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/CreateHostname',
        request,
        metadata || {},
        this.methodDescriptorCreateHostname,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/CreateHostname',
    request,
    metadata || {},
    this.methodDescriptorCreateHostname);
  }

  methodDescriptorDeleteHostname = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/DeleteHostname',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.DeleteHostnameRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.DeleteHostnameRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  deleteHostname(
    request: pkg_server_proto_server_pb.DeleteHostnameRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  deleteHostname(
    request: pkg_server_proto_server_pb.DeleteHostnameRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  deleteHostname(
    request: pkg_server_proto_server_pb.DeleteHostnameRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/DeleteHostname',
        request,
        metadata || {},
        this.methodDescriptorDeleteHostname,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/DeleteHostname',
    request,
    metadata || {},
    this.methodDescriptorDeleteHostname);
  }

  methodDescriptorListHostnames = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListHostnames',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListHostnamesRequest,
    pkg_server_proto_server_pb.ListHostnamesResponse,
    (request: pkg_server_proto_server_pb.ListHostnamesRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListHostnamesResponse.deserializeBinary
  );

  listHostnames(
    request: pkg_server_proto_server_pb.ListHostnamesRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListHostnamesResponse>;

  listHostnames(
    request: pkg_server_proto_server_pb.ListHostnamesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListHostnamesResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListHostnamesResponse>;

  listHostnames(
    request: pkg_server_proto_server_pb.ListHostnamesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListHostnamesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListHostnames',
        request,
        metadata || {},
        this.methodDescriptorListHostnames,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListHostnames',
    request,
    metadata || {},
    this.methodDescriptorListHostnames);
  }

  methodDescriptorQueueJob = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/QueueJob',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.QueueJobRequest,
    pkg_server_proto_server_pb.QueueJobResponse,
    (request: pkg_server_proto_server_pb.QueueJobRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.QueueJobResponse.deserializeBinary
  );

  queueJob(
    request: pkg_server_proto_server_pb.QueueJobRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.QueueJobResponse>;

  queueJob(
    request: pkg_server_proto_server_pb.QueueJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.QueueJobResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.QueueJobResponse>;

  queueJob(
    request: pkg_server_proto_server_pb.QueueJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.QueueJobResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/QueueJob',
        request,
        metadata || {},
        this.methodDescriptorQueueJob,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/QueueJob',
    request,
    metadata || {},
    this.methodDescriptorQueueJob);
  }

  methodDescriptorCancelJob = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/CancelJob',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.CancelJobRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.CancelJobRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  cancelJob(
    request: pkg_server_proto_server_pb.CancelJobRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  cancelJob(
    request: pkg_server_proto_server_pb.CancelJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  cancelJob(
    request: pkg_server_proto_server_pb.CancelJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/CancelJob',
        request,
        metadata || {},
        this.methodDescriptorCancelJob,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/CancelJob',
    request,
    metadata || {},
    this.methodDescriptorCancelJob);
  }

  methodDescriptorGetJob = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetJob',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetJobRequest,
    pkg_server_proto_server_pb.Job,
    (request: pkg_server_proto_server_pb.GetJobRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.Job.deserializeBinary
  );

  getJob(
    request: pkg_server_proto_server_pb.GetJobRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.Job>;

  getJob(
    request: pkg_server_proto_server_pb.GetJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Job) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.Job>;

  getJob(
    request: pkg_server_proto_server_pb.GetJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Job) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetJob',
        request,
        metadata || {},
        this.methodDescriptorGetJob,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetJob',
    request,
    metadata || {},
    this.methodDescriptorGetJob);
  }

  methodDescriptorListJobs = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListJobs',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListJobsRequest,
    pkg_server_proto_server_pb.ListJobsResponse,
    (request: pkg_server_proto_server_pb.ListJobsRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListJobsResponse.deserializeBinary
  );

  listJobs(
    request: pkg_server_proto_server_pb.ListJobsRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListJobsResponse>;

  listJobs(
    request: pkg_server_proto_server_pb.ListJobsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListJobsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListJobsResponse>;

  listJobs(
    request: pkg_server_proto_server_pb.ListJobsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListJobsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListJobs',
        request,
        metadata || {},
        this.methodDescriptorListJobs,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListJobs',
    request,
    metadata || {},
    this.methodDescriptorListJobs);
  }

  methodDescriptorValidateJob = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ValidateJob',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ValidateJobRequest,
    pkg_server_proto_server_pb.ValidateJobResponse,
    (request: pkg_server_proto_server_pb.ValidateJobRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ValidateJobResponse.deserializeBinary
  );

  validateJob(
    request: pkg_server_proto_server_pb.ValidateJobRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ValidateJobResponse>;

  validateJob(
    request: pkg_server_proto_server_pb.ValidateJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ValidateJobResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ValidateJobResponse>;

  validateJob(
    request: pkg_server_proto_server_pb.ValidateJobRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ValidateJobResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ValidateJob',
        request,
        metadata || {},
        this.methodDescriptorValidateJob,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ValidateJob',
    request,
    metadata || {},
    this.methodDescriptorValidateJob);
  }

  methodDescriptorGetJobStream = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetJobStream',
    grpcWeb.MethodType.SERVER_STREAMING,
    pkg_server_proto_server_pb.GetJobStreamRequest,
    pkg_server_proto_server_pb.GetJobStreamResponse,
    (request: pkg_server_proto_server_pb.GetJobStreamRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetJobStreamResponse.deserializeBinary
  );

  getJobStream(
    request: pkg_server_proto_server_pb.GetJobStreamRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetJobStreamResponse> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/hashicorp.waypoint.Waypoint/GetJobStream',
      request,
      metadata || {},
      this.methodDescriptorGetJobStream);
  }

  methodDescriptorGetRunner = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetRunner',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetRunnerRequest,
    pkg_server_proto_server_pb.Runner,
    (request: pkg_server_proto_server_pb.GetRunnerRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.Runner.deserializeBinary
  );

  getRunner(
    request: pkg_server_proto_server_pb.GetRunnerRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.Runner>;

  getRunner(
    request: pkg_server_proto_server_pb.GetRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Runner) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.Runner>;

  getRunner(
    request: pkg_server_proto_server_pb.GetRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.Runner) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetRunner',
        request,
        metadata || {},
        this.methodDescriptorGetRunner,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetRunner',
    request,
    metadata || {},
    this.methodDescriptorGetRunner);
  }

  methodDescriptorListRunners = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListRunners',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListRunnersRequest,
    pkg_server_proto_server_pb.ListRunnersResponse,
    (request: pkg_server_proto_server_pb.ListRunnersRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListRunnersResponse.deserializeBinary
  );

  listRunners(
    request: pkg_server_proto_server_pb.ListRunnersRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListRunnersResponse>;

  listRunners(
    request: pkg_server_proto_server_pb.ListRunnersRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListRunnersResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListRunnersResponse>;

  listRunners(
    request: pkg_server_proto_server_pb.ListRunnersRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListRunnersResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListRunners',
        request,
        metadata || {},
        this.methodDescriptorListRunners,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListRunners',
    request,
    metadata || {},
    this.methodDescriptorListRunners);
  }

  methodDescriptorAdoptRunner = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/AdoptRunner',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.AdoptRunnerRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.AdoptRunnerRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  adoptRunner(
    request: pkg_server_proto_server_pb.AdoptRunnerRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  adoptRunner(
    request: pkg_server_proto_server_pb.AdoptRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  adoptRunner(
    request: pkg_server_proto_server_pb.AdoptRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/AdoptRunner',
        request,
        metadata || {},
        this.methodDescriptorAdoptRunner,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/AdoptRunner',
    request,
    metadata || {},
    this.methodDescriptorAdoptRunner);
  }

  methodDescriptorForgetRunner = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ForgetRunner',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ForgetRunnerRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.ForgetRunnerRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  forgetRunner(
    request: pkg_server_proto_server_pb.ForgetRunnerRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  forgetRunner(
    request: pkg_server_proto_server_pb.ForgetRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  forgetRunner(
    request: pkg_server_proto_server_pb.ForgetRunnerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ForgetRunner',
        request,
        metadata || {},
        this.methodDescriptorForgetRunner,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ForgetRunner',
    request,
    metadata || {},
    this.methodDescriptorForgetRunner);
  }

  methodDescriptorGetServerConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetServerConfig',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.GetServerConfigResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetServerConfigResponse.deserializeBinary
  );

  getServerConfig(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetServerConfigResponse>;

  getServerConfig(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetServerConfigResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetServerConfigResponse>;

  getServerConfig(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetServerConfigResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetServerConfig',
        request,
        metadata || {},
        this.methodDescriptorGetServerConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetServerConfig',
    request,
    metadata || {},
    this.methodDescriptorGetServerConfig);
  }

  methodDescriptorSetServerConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/SetServerConfig',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.SetServerConfigRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.SetServerConfigRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  setServerConfig(
    request: pkg_server_proto_server_pb.SetServerConfigRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  setServerConfig(
    request: pkg_server_proto_server_pb.SetServerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  setServerConfig(
    request: pkg_server_proto_server_pb.SetServerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/SetServerConfig',
        request,
        metadata || {},
        this.methodDescriptorSetServerConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/SetServerConfig',
    request,
    metadata || {},
    this.methodDescriptorSetServerConfig);
  }

  methodDescriptorCreateSnapshot = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/CreateSnapshot',
    grpcWeb.MethodType.SERVER_STREAMING,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.CreateSnapshotResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.CreateSnapshotResponse.deserializeBinary
  );

  createSnapshot(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.CreateSnapshotResponse> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/hashicorp.waypoint.Waypoint/CreateSnapshot',
      request,
      metadata || {},
      this.methodDescriptorCreateSnapshot);
  }

  methodDescriptorBootstrapToken = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/BootstrapToken',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.NewTokenResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  bootstrapToken(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.NewTokenResponse>;

  bootstrapToken(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.NewTokenResponse>;

  bootstrapToken(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/BootstrapToken',
        request,
        metadata || {},
        this.methodDescriptorBootstrapToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/BootstrapToken',
    request,
    metadata || {},
    this.methodDescriptorBootstrapToken);
  }

  methodDescriptorDecodeToken = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/DecodeToken',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.DecodeTokenRequest,
    pkg_server_proto_server_pb.DecodeTokenResponse,
    (request: pkg_server_proto_server_pb.DecodeTokenRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.DecodeTokenResponse.deserializeBinary
  );

  decodeToken(
    request: pkg_server_proto_server_pb.DecodeTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.DecodeTokenResponse>;

  decodeToken(
    request: pkg_server_proto_server_pb.DecodeTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.DecodeTokenResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.DecodeTokenResponse>;

  decodeToken(
    request: pkg_server_proto_server_pb.DecodeTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.DecodeTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/DecodeToken',
        request,
        metadata || {},
        this.methodDescriptorDecodeToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/DecodeToken',
    request,
    metadata || {},
    this.methodDescriptorDecodeToken);
  }

  methodDescriptorGenerateInviteToken = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GenerateInviteToken',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.InviteTokenRequest,
    pkg_server_proto_server_pb.NewTokenResponse,
    (request: pkg_server_proto_server_pb.InviteTokenRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  generateInviteToken(
    request: pkg_server_proto_server_pb.InviteTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.NewTokenResponse>;

  generateInviteToken(
    request: pkg_server_proto_server_pb.InviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.NewTokenResponse>;

  generateInviteToken(
    request: pkg_server_proto_server_pb.InviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GenerateInviteToken',
        request,
        metadata || {},
        this.methodDescriptorGenerateInviteToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GenerateInviteToken',
    request,
    metadata || {},
    this.methodDescriptorGenerateInviteToken);
  }

  methodDescriptorGenerateLoginToken = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GenerateLoginToken',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.LoginTokenRequest,
    pkg_server_proto_server_pb.NewTokenResponse,
    (request: pkg_server_proto_server_pb.LoginTokenRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  generateLoginToken(
    request: pkg_server_proto_server_pb.LoginTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.NewTokenResponse>;

  generateLoginToken(
    request: pkg_server_proto_server_pb.LoginTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.NewTokenResponse>;

  generateLoginToken(
    request: pkg_server_proto_server_pb.LoginTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GenerateLoginToken',
        request,
        metadata || {},
        this.methodDescriptorGenerateLoginToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GenerateLoginToken',
    request,
    metadata || {},
    this.methodDescriptorGenerateLoginToken);
  }

  methodDescriptorGenerateRunnerToken = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GenerateRunnerToken',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GenerateRunnerTokenRequest,
    pkg_server_proto_server_pb.NewTokenResponse,
    (request: pkg_server_proto_server_pb.GenerateRunnerTokenRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  generateRunnerToken(
    request: pkg_server_proto_server_pb.GenerateRunnerTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.NewTokenResponse>;

  generateRunnerToken(
    request: pkg_server_proto_server_pb.GenerateRunnerTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.NewTokenResponse>;

  generateRunnerToken(
    request: pkg_server_proto_server_pb.GenerateRunnerTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GenerateRunnerToken',
        request,
        metadata || {},
        this.methodDescriptorGenerateRunnerToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GenerateRunnerToken',
    request,
    metadata || {},
    this.methodDescriptorGenerateRunnerToken);
  }

  methodDescriptorConvertInviteToken = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ConvertInviteToken',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ConvertInviteTokenRequest,
    pkg_server_proto_server_pb.NewTokenResponse,
    (request: pkg_server_proto_server_pb.ConvertInviteTokenRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.NewTokenResponse.deserializeBinary
  );

  convertInviteToken(
    request: pkg_server_proto_server_pb.ConvertInviteTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.NewTokenResponse>;

  convertInviteToken(
    request: pkg_server_proto_server_pb.ConvertInviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.NewTokenResponse>;

  convertInviteToken(
    request: pkg_server_proto_server_pb.ConvertInviteTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.NewTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ConvertInviteToken',
        request,
        metadata || {},
        this.methodDescriptorConvertInviteToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ConvertInviteToken',
    request,
    metadata || {},
    this.methodDescriptorConvertInviteToken);
  }

  methodDescriptorRunnerToken = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/RunnerToken',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.RunnerTokenRequest,
    pkg_server_proto_server_pb.RunnerTokenResponse,
    (request: pkg_server_proto_server_pb.RunnerTokenRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.RunnerTokenResponse.deserializeBinary
  );

  runnerToken(
    request: pkg_server_proto_server_pb.RunnerTokenRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.RunnerTokenResponse>;

  runnerToken(
    request: pkg_server_proto_server_pb.RunnerTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunnerTokenResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.RunnerTokenResponse>;

  runnerToken(
    request: pkg_server_proto_server_pb.RunnerTokenRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunnerTokenResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/RunnerToken',
        request,
        metadata || {},
        this.methodDescriptorRunnerToken,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/RunnerToken',
    request,
    metadata || {},
    this.methodDescriptorRunnerToken);
  }

  methodDescriptorRunnerGetDeploymentConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/RunnerGetDeploymentConfig',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.RunnerGetDeploymentConfigRequest,
    pkg_server_proto_server_pb.RunnerGetDeploymentConfigResponse,
    (request: pkg_server_proto_server_pb.RunnerGetDeploymentConfigRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.RunnerGetDeploymentConfigResponse.deserializeBinary
  );

  runnerGetDeploymentConfig(
    request: pkg_server_proto_server_pb.RunnerGetDeploymentConfigRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.RunnerGetDeploymentConfigResponse>;

  runnerGetDeploymentConfig(
    request: pkg_server_proto_server_pb.RunnerGetDeploymentConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunnerGetDeploymentConfigResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.RunnerGetDeploymentConfigResponse>;

  runnerGetDeploymentConfig(
    request: pkg_server_proto_server_pb.RunnerGetDeploymentConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunnerGetDeploymentConfigResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/RunnerGetDeploymentConfig',
        request,
        metadata || {},
        this.methodDescriptorRunnerGetDeploymentConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/RunnerGetDeploymentConfig',
    request,
    metadata || {},
    this.methodDescriptorRunnerGetDeploymentConfig);
  }

  methodDescriptorEntrypointConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/EntrypointConfig',
    grpcWeb.MethodType.SERVER_STREAMING,
    pkg_server_proto_server_pb.EntrypointConfigRequest,
    pkg_server_proto_server_pb.EntrypointConfigResponse,
    (request: pkg_server_proto_server_pb.EntrypointConfigRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.EntrypointConfigResponse.deserializeBinary
  );

  entrypointConfig(
    request: pkg_server_proto_server_pb.EntrypointConfigRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.EntrypointConfigResponse> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/hashicorp.waypoint.Waypoint/EntrypointConfig',
      request,
      metadata || {},
      this.methodDescriptorEntrypointConfig);
  }

  methodDescriptorWaypointHclFmt = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/WaypointHclFmt',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.WaypointHclFmtRequest,
    pkg_server_proto_server_pb.WaypointHclFmtResponse,
    (request: pkg_server_proto_server_pb.WaypointHclFmtRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.WaypointHclFmtResponse.deserializeBinary
  );

  waypointHclFmt(
    request: pkg_server_proto_server_pb.WaypointHclFmtRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.WaypointHclFmtResponse>;

  waypointHclFmt(
    request: pkg_server_proto_server_pb.WaypointHclFmtRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.WaypointHclFmtResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.WaypointHclFmtResponse>;

  waypointHclFmt(
    request: pkg_server_proto_server_pb.WaypointHclFmtRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.WaypointHclFmtResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/WaypointHclFmt',
        request,
        metadata || {},
        this.methodDescriptorWaypointHclFmt,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/WaypointHclFmt',
    request,
    metadata || {},
    this.methodDescriptorWaypointHclFmt);
  }

  methodDescriptorUpsertOnDemandRunnerConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertOnDemandRunnerConfig',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigRequest,
    pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigResponse,
    (request: pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigResponse.deserializeBinary
  );

  upsertOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigResponse>;

  upsertOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigResponse>;

  upsertOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertOnDemandRunnerConfigResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertOnDemandRunnerConfig',
        request,
        metadata || {},
        this.methodDescriptorUpsertOnDemandRunnerConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertOnDemandRunnerConfig',
    request,
    metadata || {},
    this.methodDescriptorUpsertOnDemandRunnerConfig);
  }

  methodDescriptorGetOnDemandRunnerConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetOnDemandRunnerConfig',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetOnDemandRunnerConfigRequest,
    pkg_server_proto_server_pb.GetOnDemandRunnerConfigResponse,
    (request: pkg_server_proto_server_pb.GetOnDemandRunnerConfigRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetOnDemandRunnerConfigResponse.deserializeBinary
  );

  getOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.GetOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetOnDemandRunnerConfigResponse>;

  getOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.GetOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetOnDemandRunnerConfigResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetOnDemandRunnerConfigResponse>;

  getOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.GetOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetOnDemandRunnerConfigResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetOnDemandRunnerConfig',
        request,
        metadata || {},
        this.methodDescriptorGetOnDemandRunnerConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetOnDemandRunnerConfig',
    request,
    metadata || {},
    this.methodDescriptorGetOnDemandRunnerConfig);
  }

  methodDescriptorDeleteOnDemandRunnerConfig = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/DeleteOnDemandRunnerConfig',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigRequest,
    pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigResponse,
    (request: pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigResponse.deserializeBinary
  );

  deleteOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigResponse>;

  deleteOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigResponse>;

  deleteOnDemandRunnerConfig(
    request: pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.DeleteOnDemandRunnerConfigResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/DeleteOnDemandRunnerConfig',
        request,
        metadata || {},
        this.methodDescriptorDeleteOnDemandRunnerConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/DeleteOnDemandRunnerConfig',
    request,
    metadata || {},
    this.methodDescriptorDeleteOnDemandRunnerConfig);
  }

  methodDescriptorListOnDemandRunnerConfigs = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListOnDemandRunnerConfigs',
    grpcWeb.MethodType.UNARY,
    google_protobuf_empty_pb.Empty,
    pkg_server_proto_server_pb.ListOnDemandRunnerConfigsResponse,
    (request: google_protobuf_empty_pb.Empty) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListOnDemandRunnerConfigsResponse.deserializeBinary
  );

  listOnDemandRunnerConfigs(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListOnDemandRunnerConfigsResponse>;

  listOnDemandRunnerConfigs(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListOnDemandRunnerConfigsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListOnDemandRunnerConfigsResponse>;

  listOnDemandRunnerConfigs(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListOnDemandRunnerConfigsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListOnDemandRunnerConfigs',
        request,
        metadata || {},
        this.methodDescriptorListOnDemandRunnerConfigs,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListOnDemandRunnerConfigs',
    request,
    metadata || {},
    this.methodDescriptorListOnDemandRunnerConfigs);
  }

  methodDescriptorUpsertBuild = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertBuild',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertBuildRequest,
    pkg_server_proto_server_pb.UpsertBuildResponse,
    (request: pkg_server_proto_server_pb.UpsertBuildRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertBuildResponse.deserializeBinary
  );

  upsertBuild(
    request: pkg_server_proto_server_pb.UpsertBuildRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertBuildResponse>;

  upsertBuild(
    request: pkg_server_proto_server_pb.UpsertBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertBuildResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertBuildResponse>;

  upsertBuild(
    request: pkg_server_proto_server_pb.UpsertBuildRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertBuildResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertBuild',
        request,
        metadata || {},
        this.methodDescriptorUpsertBuild,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertBuild',
    request,
    metadata || {},
    this.methodDescriptorUpsertBuild);
  }

  methodDescriptorUpsertPushedArtifact = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertPushedArtifact',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertPushedArtifactRequest,
    pkg_server_proto_server_pb.UpsertPushedArtifactResponse,
    (request: pkg_server_proto_server_pb.UpsertPushedArtifactRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertPushedArtifactResponse.deserializeBinary
  );

  upsertPushedArtifact(
    request: pkg_server_proto_server_pb.UpsertPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertPushedArtifactResponse>;

  upsertPushedArtifact(
    request: pkg_server_proto_server_pb.UpsertPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertPushedArtifactResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertPushedArtifactResponse>;

  upsertPushedArtifact(
    request: pkg_server_proto_server_pb.UpsertPushedArtifactRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertPushedArtifactResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertPushedArtifact',
        request,
        metadata || {},
        this.methodDescriptorUpsertPushedArtifact,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertPushedArtifact',
    request,
    metadata || {},
    this.methodDescriptorUpsertPushedArtifact);
  }

  methodDescriptorUpsertDeployment = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertDeployment',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertDeploymentRequest,
    pkg_server_proto_server_pb.UpsertDeploymentResponse,
    (request: pkg_server_proto_server_pb.UpsertDeploymentRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertDeploymentResponse.deserializeBinary
  );

  upsertDeployment(
    request: pkg_server_proto_server_pb.UpsertDeploymentRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertDeploymentResponse>;

  upsertDeployment(
    request: pkg_server_proto_server_pb.UpsertDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertDeploymentResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertDeploymentResponse>;

  upsertDeployment(
    request: pkg_server_proto_server_pb.UpsertDeploymentRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertDeploymentResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertDeployment',
        request,
        metadata || {},
        this.methodDescriptorUpsertDeployment,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertDeployment',
    request,
    metadata || {},
    this.methodDescriptorUpsertDeployment);
  }

  methodDescriptorUpsertRelease = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertRelease',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertReleaseRequest,
    pkg_server_proto_server_pb.UpsertReleaseResponse,
    (request: pkg_server_proto_server_pb.UpsertReleaseRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertReleaseResponse.deserializeBinary
  );

  upsertRelease(
    request: pkg_server_proto_server_pb.UpsertReleaseRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertReleaseResponse>;

  upsertRelease(
    request: pkg_server_proto_server_pb.UpsertReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertReleaseResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertReleaseResponse>;

  upsertRelease(
    request: pkg_server_proto_server_pb.UpsertReleaseRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertReleaseResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertRelease',
        request,
        metadata || {},
        this.methodDescriptorUpsertRelease,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertRelease',
    request,
    metadata || {},
    this.methodDescriptorUpsertRelease);
  }

  methodDescriptorUpsertStatusReport = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertStatusReport',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertStatusReportRequest,
    pkg_server_proto_server_pb.UpsertStatusReportResponse,
    (request: pkg_server_proto_server_pb.UpsertStatusReportRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertStatusReportResponse.deserializeBinary
  );

  upsertStatusReport(
    request: pkg_server_proto_server_pb.UpsertStatusReportRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertStatusReportResponse>;

  upsertStatusReport(
    request: pkg_server_proto_server_pb.UpsertStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertStatusReportResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertStatusReportResponse>;

  upsertStatusReport(
    request: pkg_server_proto_server_pb.UpsertStatusReportRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertStatusReportResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertStatusReport',
        request,
        metadata || {},
        this.methodDescriptorUpsertStatusReport,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertStatusReport',
    request,
    metadata || {},
    this.methodDescriptorUpsertStatusReport);
  }

  methodDescriptorGetTask = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetTask',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetTaskRequest,
    pkg_server_proto_server_pb.GetTaskResponse,
    (request: pkg_server_proto_server_pb.GetTaskRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetTaskResponse.deserializeBinary
  );

  getTask(
    request: pkg_server_proto_server_pb.GetTaskRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetTaskResponse>;

  getTask(
    request: pkg_server_proto_server_pb.GetTaskRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetTaskResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetTaskResponse>;

  getTask(
    request: pkg_server_proto_server_pb.GetTaskRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetTaskResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetTask',
        request,
        metadata || {},
        this.methodDescriptorGetTask,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetTask',
    request,
    metadata || {},
    this.methodDescriptorGetTask);
  }

  methodDescriptorListTask = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListTask',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListTaskRequest,
    pkg_server_proto_server_pb.ListTaskResponse,
    (request: pkg_server_proto_server_pb.ListTaskRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListTaskResponse.deserializeBinary
  );

  listTask(
    request: pkg_server_proto_server_pb.ListTaskRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListTaskResponse>;

  listTask(
    request: pkg_server_proto_server_pb.ListTaskRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListTaskResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListTaskResponse>;

  listTask(
    request: pkg_server_proto_server_pb.ListTaskRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListTaskResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListTask',
        request,
        metadata || {},
        this.methodDescriptorListTask,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListTask',
    request,
    metadata || {},
    this.methodDescriptorListTask);
  }

  methodDescriptorCancelTask = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/CancelTask',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.CancelTaskRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.CancelTaskRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  cancelTask(
    request: pkg_server_proto_server_pb.CancelTaskRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  cancelTask(
    request: pkg_server_proto_server_pb.CancelTaskRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  cancelTask(
    request: pkg_server_proto_server_pb.CancelTaskRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/CancelTask',
        request,
        metadata || {},
        this.methodDescriptorCancelTask,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/CancelTask',
    request,
    metadata || {},
    this.methodDescriptorCancelTask);
  }

  methodDescriptorUpsertTrigger = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertTrigger',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertTriggerRequest,
    pkg_server_proto_server_pb.UpsertTriggerResponse,
    (request: pkg_server_proto_server_pb.UpsertTriggerRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertTriggerResponse.deserializeBinary
  );

  upsertTrigger(
    request: pkg_server_proto_server_pb.UpsertTriggerRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertTriggerResponse>;

  upsertTrigger(
    request: pkg_server_proto_server_pb.UpsertTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertTriggerResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertTriggerResponse>;

  upsertTrigger(
    request: pkg_server_proto_server_pb.UpsertTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertTriggerResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertTrigger',
        request,
        metadata || {},
        this.methodDescriptorUpsertTrigger,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertTrigger',
    request,
    metadata || {},
    this.methodDescriptorUpsertTrigger);
  }

  methodDescriptorGetTrigger = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetTrigger',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetTriggerRequest,
    pkg_server_proto_server_pb.GetTriggerResponse,
    (request: pkg_server_proto_server_pb.GetTriggerRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetTriggerResponse.deserializeBinary
  );

  getTrigger(
    request: pkg_server_proto_server_pb.GetTriggerRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetTriggerResponse>;

  getTrigger(
    request: pkg_server_proto_server_pb.GetTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetTriggerResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetTriggerResponse>;

  getTrigger(
    request: pkg_server_proto_server_pb.GetTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetTriggerResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetTrigger',
        request,
        metadata || {},
        this.methodDescriptorGetTrigger,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetTrigger',
    request,
    metadata || {},
    this.methodDescriptorGetTrigger);
  }

  methodDescriptorDeleteTrigger = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/DeleteTrigger',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.DeleteTriggerRequest,
    google_protobuf_empty_pb.Empty,
    (request: pkg_server_proto_server_pb.DeleteTriggerRequest) => {
      return request.serializeBinary();
    },
    google_protobuf_empty_pb.Empty.deserializeBinary
  );

  deleteTrigger(
    request: pkg_server_proto_server_pb.DeleteTriggerRequest,
    metadata: grpcWeb.Metadata | null): Promise<google_protobuf_empty_pb.Empty>;

  deleteTrigger(
    request: pkg_server_proto_server_pb.DeleteTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  deleteTrigger(
    request: pkg_server_proto_server_pb.DeleteTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: google_protobuf_empty_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/DeleteTrigger',
        request,
        metadata || {},
        this.methodDescriptorDeleteTrigger,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/DeleteTrigger',
    request,
    metadata || {},
    this.methodDescriptorDeleteTrigger);
  }

  methodDescriptorListTriggers = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListTriggers',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListTriggerRequest,
    pkg_server_proto_server_pb.ListTriggerResponse,
    (request: pkg_server_proto_server_pb.ListTriggerRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListTriggerResponse.deserializeBinary
  );

  listTriggers(
    request: pkg_server_proto_server_pb.ListTriggerRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListTriggerResponse>;

  listTriggers(
    request: pkg_server_proto_server_pb.ListTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListTriggerResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListTriggerResponse>;

  listTriggers(
    request: pkg_server_proto_server_pb.ListTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListTriggerResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListTriggers',
        request,
        metadata || {},
        this.methodDescriptorListTriggers,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListTriggers',
    request,
    metadata || {},
    this.methodDescriptorListTriggers);
  }

  methodDescriptorRunTrigger = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/RunTrigger',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.RunTriggerRequest,
    pkg_server_proto_server_pb.RunTriggerResponse,
    (request: pkg_server_proto_server_pb.RunTriggerRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.RunTriggerResponse.deserializeBinary
  );

  runTrigger(
    request: pkg_server_proto_server_pb.RunTriggerRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.RunTriggerResponse>;

  runTrigger(
    request: pkg_server_proto_server_pb.RunTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunTriggerResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.RunTriggerResponse>;

  runTrigger(
    request: pkg_server_proto_server_pb.RunTriggerRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunTriggerResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/RunTrigger',
        request,
        metadata || {},
        this.methodDescriptorRunTrigger,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/RunTrigger',
    request,
    metadata || {},
    this.methodDescriptorRunTrigger);
  }

  methodDescriptorUpsertPipeline = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UpsertPipeline',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UpsertPipelineRequest,
    pkg_server_proto_server_pb.UpsertPipelineResponse,
    (request: pkg_server_proto_server_pb.UpsertPipelineRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UpsertPipelineResponse.deserializeBinary
  );

  upsertPipeline(
    request: pkg_server_proto_server_pb.UpsertPipelineRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UpsertPipelineResponse>;

  upsertPipeline(
    request: pkg_server_proto_server_pb.UpsertPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertPipelineResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UpsertPipelineResponse>;

  upsertPipeline(
    request: pkg_server_proto_server_pb.UpsertPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UpsertPipelineResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UpsertPipeline',
        request,
        metadata || {},
        this.methodDescriptorUpsertPipeline,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UpsertPipeline',
    request,
    metadata || {},
    this.methodDescriptorUpsertPipeline);
  }

  methodDescriptorRunPipeline = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/RunPipeline',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.RunPipelineRequest,
    pkg_server_proto_server_pb.RunPipelineResponse,
    (request: pkg_server_proto_server_pb.RunPipelineRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.RunPipelineResponse.deserializeBinary
  );

  runPipeline(
    request: pkg_server_proto_server_pb.RunPipelineRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.RunPipelineResponse>;

  runPipeline(
    request: pkg_server_proto_server_pb.RunPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunPipelineResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.RunPipelineResponse>;

  runPipeline(
    request: pkg_server_proto_server_pb.RunPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.RunPipelineResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/RunPipeline',
        request,
        metadata || {},
        this.methodDescriptorRunPipeline,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/RunPipeline',
    request,
    metadata || {},
    this.methodDescriptorRunPipeline);
  }

  methodDescriptorGetPipeline = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/GetPipeline',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.GetPipelineRequest,
    pkg_server_proto_server_pb.GetPipelineResponse,
    (request: pkg_server_proto_server_pb.GetPipelineRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.GetPipelineResponse.deserializeBinary
  );

  getPipeline(
    request: pkg_server_proto_server_pb.GetPipelineRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.GetPipelineResponse>;

  getPipeline(
    request: pkg_server_proto_server_pb.GetPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetPipelineResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.GetPipelineResponse>;

  getPipeline(
    request: pkg_server_proto_server_pb.GetPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.GetPipelineResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/GetPipeline',
        request,
        metadata || {},
        this.methodDescriptorGetPipeline,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/GetPipeline',
    request,
    metadata || {},
    this.methodDescriptorGetPipeline);
  }

  methodDescriptorListPipelines = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ListPipelines',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ListPipelinesRequest,
    pkg_server_proto_server_pb.ListPipelinesResponse,
    (request: pkg_server_proto_server_pb.ListPipelinesRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ListPipelinesResponse.deserializeBinary
  );

  listPipelines(
    request: pkg_server_proto_server_pb.ListPipelinesRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ListPipelinesResponse>;

  listPipelines(
    request: pkg_server_proto_server_pb.ListPipelinesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListPipelinesResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ListPipelinesResponse>;

  listPipelines(
    request: pkg_server_proto_server_pb.ListPipelinesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ListPipelinesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ListPipelines',
        request,
        metadata || {},
        this.methodDescriptorListPipelines,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ListPipelines',
    request,
    metadata || {},
    this.methodDescriptorListPipelines);
  }

  methodDescriptorConfigSyncPipeline = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/ConfigSyncPipeline',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.ConfigSyncPipelineRequest,
    pkg_server_proto_server_pb.ConfigSyncPipelineResponse,
    (request: pkg_server_proto_server_pb.ConfigSyncPipelineRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.ConfigSyncPipelineResponse.deserializeBinary
  );

  configSyncPipeline(
    request: pkg_server_proto_server_pb.ConfigSyncPipelineRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.ConfigSyncPipelineResponse>;

  configSyncPipeline(
    request: pkg_server_proto_server_pb.ConfigSyncPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ConfigSyncPipelineResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.ConfigSyncPipelineResponse>;

  configSyncPipeline(
    request: pkg_server_proto_server_pb.ConfigSyncPipelineRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.ConfigSyncPipelineResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/ConfigSyncPipeline',
        request,
        metadata || {},
        this.methodDescriptorConfigSyncPipeline,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/ConfigSyncPipeline',
    request,
    metadata || {},
    this.methodDescriptorConfigSyncPipeline);
  }

  methodDescriptorUI_GetProject = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UI_GetProject',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UI.GetProjectRequest,
    pkg_server_proto_server_pb.UI.GetProjectResponse,
    (request: pkg_server_proto_server_pb.UI.GetProjectRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UI.GetProjectResponse.deserializeBinary
  );

  uI_GetProject(
    request: pkg_server_proto_server_pb.UI.GetProjectRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UI.GetProjectResponse>;

  uI_GetProject(
    request: pkg_server_proto_server_pb.UI.GetProjectRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UI.GetProjectResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UI.GetProjectResponse>;

  uI_GetProject(
    request: pkg_server_proto_server_pb.UI.GetProjectRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UI.GetProjectResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UI_GetProject',
        request,
        metadata || {},
        this.methodDescriptorUI_GetProject,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UI_GetProject',
    request,
    metadata || {},
    this.methodDescriptorUI_GetProject);
  }

  methodDescriptorUI_ListDeployments = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UI_ListDeployments',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UI.ListDeploymentsRequest,
    pkg_server_proto_server_pb.UI.ListDeploymentsResponse,
    (request: pkg_server_proto_server_pb.UI.ListDeploymentsRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UI.ListDeploymentsResponse.deserializeBinary
  );

  uI_ListDeployments(
    request: pkg_server_proto_server_pb.UI.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UI.ListDeploymentsResponse>;

  uI_ListDeployments(
    request: pkg_server_proto_server_pb.UI.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UI.ListDeploymentsResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UI.ListDeploymentsResponse>;

  uI_ListDeployments(
    request: pkg_server_proto_server_pb.UI.ListDeploymentsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UI.ListDeploymentsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UI_ListDeployments',
        request,
        metadata || {},
        this.methodDescriptorUI_ListDeployments,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UI_ListDeployments',
    request,
    metadata || {},
    this.methodDescriptorUI_ListDeployments);
  }

  methodDescriptorUI_ListReleases = new grpcWeb.MethodDescriptor(
    '/hashicorp.waypoint.Waypoint/UI_ListReleases',
    grpcWeb.MethodType.UNARY,
    pkg_server_proto_server_pb.UI.ListReleasesRequest,
    pkg_server_proto_server_pb.UI.ListReleasesResponse,
    (request: pkg_server_proto_server_pb.UI.ListReleasesRequest) => {
      return request.serializeBinary();
    },
    pkg_server_proto_server_pb.UI.ListReleasesResponse.deserializeBinary
  );

  uI_ListReleases(
    request: pkg_server_proto_server_pb.UI.ListReleasesRequest,
    metadata: grpcWeb.Metadata | null): Promise<pkg_server_proto_server_pb.UI.ListReleasesResponse>;

  uI_ListReleases(
    request: pkg_server_proto_server_pb.UI.ListReleasesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UI.ListReleasesResponse) => void): grpcWeb.ClientReadableStream<pkg_server_proto_server_pb.UI.ListReleasesResponse>;

  uI_ListReleases(
    request: pkg_server_proto_server_pb.UI.ListReleasesRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: pkg_server_proto_server_pb.UI.ListReleasesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/hashicorp.waypoint.Waypoint/UI_ListReleases',
        request,
        metadata || {},
        this.methodDescriptorUI_ListReleases,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/hashicorp.waypoint.Waypoint/UI_ListReleases',
    request,
    metadata || {},
    this.methodDescriptorUI_ListReleases);
  }

}

