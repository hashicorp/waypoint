/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Response } from 'miragejs';

export default interface RouteHandler {
  serialize(response: unknown, serializerType: string): Response;
}
