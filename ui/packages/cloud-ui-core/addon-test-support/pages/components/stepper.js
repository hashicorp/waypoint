import { collection, hasClass, notHasClass, isPresent, text } from 'ember-cli-page-object';

export const stepperSelector = '[ data-test-stepper-container ]';
export const stepSelector = '[ data-test-stepper-step-container ]';
export const stepIconSelector = '[ data-test-stepper-step-icon ]';
export const stepLabelSelector = '[ data-test-stepper-step-label ]';
export const stepContentSelector = '[ data-test-stepper-step-content-container ]';
export default {
  showsContainer: isPresent(stepperSelector),
  steps: collection(stepSelector, {
    isIconComplete: hasClass('text--success', stepIconSelector),
    isIconIncomplete: hasClass('text--muted', stepIconSelector),
    isLabelComplete: hasClass('text--muted', stepLabelSelector),
    isLabelIncomplete: notHasClass('text--muted', stepLabelSelector),
    label: text(stepLabelSelector),
    content: text(stepContentSelector),
  }),
};
