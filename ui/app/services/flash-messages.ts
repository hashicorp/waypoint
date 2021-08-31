import FlashMessages from 'ember-cli-flash/services/flash-messages';

type BaseFlashFunction = FlashMessages['info'];
type FlashFunctionParams = Parameters<BaseFlashFunction>;

type FlashFunction = (
  message: FlashFunctionParams[0],
  // This allows the “throw any property you like on the flash object”
  // behavior of ember-cli-flash to type check.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  options?: FlashFunctionParams[1] & Record<string, any>
) => ReturnType<BaseFlashFunction>;

class PdsFlashMessages extends FlashMessages {
  // This is the one custom convenience method we register in
  // config/environment.js.
  //
  // The superclass will dynamically initialize this method when the
  // service is created (thus the extra exclamation to reassure
  // TypeScript).
  error!: FlashFunction;
}

export default PdsFlashMessages;
