import Component from '@glimmer/component';

/**
 *
 * `StepperStep` renders a single step that should be used within the Stepper
 *     component
 *
 *
 * ```
 * <Stepper::Step
 *   @completed={{true}}
 *   as |S|
 * >
 *   <S.StepLabel>
 *     Label
 *   </S.StepLabel>
 *   <S.StepContent>
 *     Content
 *   </S.StepContent>
 * </Stepper::Step>
 * ```
 *
 * @class StepperStep
 * @yield {StepperStepLabel} StepLabel `Stepper::StepLabel` component
 * @yield {StepperStepContent} StepContent `Stepper::StepContent` component
 *
 */

export default class StepperStepComponent extends Component {
  /**
   * A boolean value to visually render the completed state.
   * @argument completed
   * @type {boolean}
   */
}
