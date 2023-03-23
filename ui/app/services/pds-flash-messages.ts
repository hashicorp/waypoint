/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import FlashMessages from 'ember-cli-flash/services/flash-messages';

type BaseFlashFunction = FlashMessages['info'];
type FlashFunctionParams = Parameters<BaseFlashFunction>;

type FlashFunction = (
  message: FlashFunctionParams[0],
  // This allows the “throw any property you like on the flash object”
  // behavior of ember-cli-flash to type check.
  options?: FlashFunctionParams[1] & Record<string, unknown>
) => ReturnType<BaseFlashFunction>;

class PdsFlashMessages extends FlashMessages {
  // This is the one custom convenience method we register in
  // config/environment.js.
  //
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  error(..._args: Parameters<FlashFunction>): FlashMessages {
    // This is a stub implementation. The superclass will dynamically add the
    // real implementation at runtime.
    return this;
  }
}

export default PdsFlashMessages;

// DO NOT DELETE: this is how TypeScript knows how to look up your services.
declare module '@ember/service' {
  interface Registry {
    pdsFlashMessages: PdsFlashMessages;
  }
}
