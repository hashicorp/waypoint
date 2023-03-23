/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  /*
   * This simulates the data plugins store for build artifacts. It is
   * made available on the protobuf as artifact_json. For testing
   * purposes, assign a POJO to this field and it’ll be automatically
   * converted to JSON during Mirage serialization.
   *
   * @example
   * let artifact = server.create('artifact', {
   *   artifact: {
   *     some_plugin_field: 'some-plugin-value'
   *   }
   * });
   * artifact.toProtobuf().getArtifactJson();
   * // => '{"some_plugin_field":"some-plugin-value"}'
   */
  artifact: () => ({}),
});
