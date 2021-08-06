// source: google/api/client.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js');
goog.object.extend(proto, google_protobuf_descriptor_pb);
goog.exportSymbol('proto.google.api.defaultHost', null, global);
goog.exportSymbol('proto.google.api.methodSignatureList', null, global);
goog.exportSymbol('proto.google.api.oauthScopes', null, global);

/**
 * A tuple of {field number, class constructor} for the extension
 * field named `defaultHost`.
 * @type {!jspb.ExtensionFieldInfo<string>}
 */
proto.google.api.defaultHost = new jspb.ExtensionFieldInfo(
    1049,
    {defaultHost: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.ServiceOptions.extensionsBinary[1049] = new jspb.ExtensionFieldBinaryInfo(
    proto.google.api.defaultHost,
    jspb.BinaryReader.prototype.readString,
    jspb.BinaryWriter.prototype.writeString,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.ServiceOptions.extensions[1049] = proto.google.api.defaultHost;


/**
 * A tuple of {field number, class constructor} for the extension
 * field named `oauthScopes`.
 * @type {!jspb.ExtensionFieldInfo<string>}
 */
proto.google.api.oauthScopes = new jspb.ExtensionFieldInfo(
    1050,
    {oauthScopes: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    0);

google_protobuf_descriptor_pb.ServiceOptions.extensionsBinary[1050] = new jspb.ExtensionFieldBinaryInfo(
    proto.google.api.oauthScopes,
    jspb.BinaryReader.prototype.readString,
    jspb.BinaryWriter.prototype.writeString,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.ServiceOptions.extensions[1050] = proto.google.api.oauthScopes;


/**
 * A tuple of {field number, class constructor} for the extension
 * field named `methodSignatureList`.
 * @type {!jspb.ExtensionFieldInfo<!Array<string>>}
 */
proto.google.api.methodSignatureList = new jspb.ExtensionFieldInfo(
    1051,
    {methodSignatureList: 0},
    null,
     /** @type {?function((boolean|undefined),!jspb.Message=): !Object} */ (
         null),
    1);

google_protobuf_descriptor_pb.MethodOptions.extensionsBinary[1051] = new jspb.ExtensionFieldBinaryInfo(
    proto.google.api.methodSignatureList,
    jspb.BinaryReader.prototype.readString,
    jspb.BinaryWriter.prototype.writeRepeatedString,
    undefined,
    undefined,
    false);
// This registers the extension field with the extended class, so that
// toObject() will function correctly.
google_protobuf_descriptor_pb.MethodOptions.extensions[1051] = proto.google.api.methodSignatureList;

goog.object.extend(exports, proto.google.api);
