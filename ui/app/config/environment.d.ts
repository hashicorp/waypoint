/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default config;

/**
 * Type declarations for
 *    import config from './config/environment'
 *
 * For now these need to be managed by the developer
 * since different ember addons can materialize new entries.
 */
declare const config: {
  apiAddress: string;
  environment: string;
  modulePrefix: string;
  podModulePrefix: string;
  locationType: string;
  rootURL: string;
};
