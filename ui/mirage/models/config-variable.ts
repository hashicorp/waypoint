/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Model, belongsTo } from 'ember-cli-mirage';
import { ConfigVar } from 'waypoint-pb';

export default Model.extend({
  project: belongsTo(),

  toProtobuf(): ConfigVar {
    let result = new ConfigVar();

    result.setProject(this.project?.toProtobufRef());
    result.setName(this.name);
    result.setStatic(this.pb_static);
    result.setInternal(this.internal);
    result.setNameIsPath(this.nameIsPath);
    if (this.dynamic) {
      let dynamicVal = new ConfigVar.DynamicVal();
      dynamicVal.setFrom(this.dynamic.from);
      // DynamicVal.ConfigMap is a map type and has no setter,
      // and the native JS Map type is not supported: https://github.com/protocolbuffers/protobuf/issues/2789
      // so we need to use the getter then modify the keys/values
      // https://developers.google.com/protocol-buffers/docs/reference/javascript-generated#map
      let configMap = dynamicVal.getConfigMap();
      this.dynamic.configMap.forEach(([key, value]) => {
        configMap.set(key, value);
      });
      result.setDynamic(dynamicVal);
    }

    return result;
  },
});
