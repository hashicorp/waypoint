/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

declare module 'docker-parse-image' {
  export default function parse(s: string): Ref;

  export interface Ref {
    registry?: string;
    namespace?: string;
    repository?: string;
    tag?: string;
  }
}
