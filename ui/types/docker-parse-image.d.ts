/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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
