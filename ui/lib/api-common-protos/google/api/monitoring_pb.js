/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// source: google/api/monitoring.proto
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

var google_api_annotations_pb = require('../../google/api/annotations_pb.js');
goog.object.extend(proto, google_api_annotations_pb);
goog.exportSymbol('proto.google.api.Monitoring', null, global);
goog.exportSymbol('proto.google.api.Monitoring.MonitoringDestination', null, global);
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
proto.google.api.Monitoring = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.google.api.Monitoring.repeatedFields_, null);
};
goog.inherits(proto.google.api.Monitoring, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.google.api.Monitoring.displayName = 'proto.google.api.Monitoring';
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
proto.google.api.Monitoring.MonitoringDestination = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.google.api.Monitoring.MonitoringDestination.repeatedFields_, null);
};
goog.inherits(proto.google.api.Monitoring.MonitoringDestination, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.google.api.Monitoring.MonitoringDestination.displayName = 'proto.google.api.Monitoring.MonitoringDestination';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.google.api.Monitoring.repeatedFields_ = [1,2];



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
proto.google.api.Monitoring.prototype.toObject = function(opt_includeInstance) {
  return proto.google.api.Monitoring.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.google.api.Monitoring} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.google.api.Monitoring.toObject = function(includeInstance, msg) {
  var f, obj = {
    producerDestinationsList: jspb.Message.toObjectList(msg.getProducerDestinationsList(),
    proto.google.api.Monitoring.MonitoringDestination.toObject, includeInstance),
    consumerDestinationsList: jspb.Message.toObjectList(msg.getConsumerDestinationsList(),
    proto.google.api.Monitoring.MonitoringDestination.toObject, includeInstance)
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
 * @return {!proto.google.api.Monitoring}
 */
proto.google.api.Monitoring.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.google.api.Monitoring;
  return proto.google.api.Monitoring.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.google.api.Monitoring} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.google.api.Monitoring}
 */
proto.google.api.Monitoring.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.google.api.Monitoring.MonitoringDestination;
      reader.readMessage(value,proto.google.api.Monitoring.MonitoringDestination.deserializeBinaryFromReader);
      msg.addProducerDestinations(value);
      break;
    case 2:
      var value = new proto.google.api.Monitoring.MonitoringDestination;
      reader.readMessage(value,proto.google.api.Monitoring.MonitoringDestination.deserializeBinaryFromReader);
      msg.addConsumerDestinations(value);
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
proto.google.api.Monitoring.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.google.api.Monitoring.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.google.api.Monitoring} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.google.api.Monitoring.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProducerDestinationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.google.api.Monitoring.MonitoringDestination.serializeBinaryToWriter
    );
  }
  f = message.getConsumerDestinationsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      proto.google.api.Monitoring.MonitoringDestination.serializeBinaryToWriter
    );
  }
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.google.api.Monitoring.MonitoringDestination.repeatedFields_ = [2];



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
proto.google.api.Monitoring.MonitoringDestination.prototype.toObject = function(opt_includeInstance) {
  return proto.google.api.Monitoring.MonitoringDestination.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.google.api.Monitoring.MonitoringDestination} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.google.api.Monitoring.MonitoringDestination.toObject = function(includeInstance, msg) {
  var f, obj = {
    monitoredResource: jspb.Message.getFieldWithDefault(msg, 1, ""),
    metricsList: (f = jspb.Message.getRepeatedField(msg, 2)) == null ? undefined : f
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
 * @return {!proto.google.api.Monitoring.MonitoringDestination}
 */
proto.google.api.Monitoring.MonitoringDestination.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.google.api.Monitoring.MonitoringDestination;
  return proto.google.api.Monitoring.MonitoringDestination.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.google.api.Monitoring.MonitoringDestination} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.google.api.Monitoring.MonitoringDestination}
 */
proto.google.api.Monitoring.MonitoringDestination.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setMonitoredResource(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.addMetrics(value);
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
proto.google.api.Monitoring.MonitoringDestination.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.google.api.Monitoring.MonitoringDestination.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.google.api.Monitoring.MonitoringDestination} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.google.api.Monitoring.MonitoringDestination.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMonitoredResource();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getMetricsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      2,
      f
    );
  }
};


/**
 * optional string monitored_resource = 1;
 * @return {string}
 */
proto.google.api.Monitoring.MonitoringDestination.prototype.getMonitoredResource = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.google.api.Monitoring.MonitoringDestination} returns this
 */
proto.google.api.Monitoring.MonitoringDestination.prototype.setMonitoredResource = function(value) {
  return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * repeated string metrics = 2;
 * @return {!Array<string>}
 */
proto.google.api.Monitoring.MonitoringDestination.prototype.getMetricsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 2));
};


/**
 * @param {!Array<string>} value
 * @return {!proto.google.api.Monitoring.MonitoringDestination} returns this
 */
proto.google.api.Monitoring.MonitoringDestination.prototype.setMetricsList = function(value) {
  return jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 * @return {!proto.google.api.Monitoring.MonitoringDestination} returns this
 */
proto.google.api.Monitoring.MonitoringDestination.prototype.addMetrics = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.google.api.Monitoring.MonitoringDestination} returns this
 */
proto.google.api.Monitoring.MonitoringDestination.prototype.clearMetricsList = function() {
  return this.setMetricsList([]);
};


/**
 * repeated MonitoringDestination producer_destinations = 1;
 * @return {!Array<!proto.google.api.Monitoring.MonitoringDestination>}
 */
proto.google.api.Monitoring.prototype.getProducerDestinationsList = function() {
  return /** @type{!Array<!proto.google.api.Monitoring.MonitoringDestination>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.google.api.Monitoring.MonitoringDestination, 1));
};


/**
 * @param {!Array<!proto.google.api.Monitoring.MonitoringDestination>} value
 * @return {!proto.google.api.Monitoring} returns this
*/
proto.google.api.Monitoring.prototype.setProducerDestinationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.google.api.Monitoring.MonitoringDestination=} opt_value
 * @param {number=} opt_index
 * @return {!proto.google.api.Monitoring.MonitoringDestination}
 */
proto.google.api.Monitoring.prototype.addProducerDestinations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.google.api.Monitoring.MonitoringDestination, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.google.api.Monitoring} returns this
 */
proto.google.api.Monitoring.prototype.clearProducerDestinationsList = function() {
  return this.setProducerDestinationsList([]);
};


/**
 * repeated MonitoringDestination consumer_destinations = 2;
 * @return {!Array<!proto.google.api.Monitoring.MonitoringDestination>}
 */
proto.google.api.Monitoring.prototype.getConsumerDestinationsList = function() {
  return /** @type{!Array<!proto.google.api.Monitoring.MonitoringDestination>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.google.api.Monitoring.MonitoringDestination, 2));
};


/**
 * @param {!Array<!proto.google.api.Monitoring.MonitoringDestination>} value
 * @return {!proto.google.api.Monitoring} returns this
*/
proto.google.api.Monitoring.prototype.setConsumerDestinationsList = function(value) {
  return jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.google.api.Monitoring.MonitoringDestination=} opt_value
 * @param {number=} opt_index
 * @return {!proto.google.api.Monitoring.MonitoringDestination}
 */
proto.google.api.Monitoring.prototype.addConsumerDestinations = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.google.api.Monitoring.MonitoringDestination, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.google.api.Monitoring} returns this
 */
proto.google.api.Monitoring.prototype.clearConsumerDestinationsList = function() {
  return this.setConsumerDestinationsList([]);
};


goog.object.extend(exports, proto.google.api);
