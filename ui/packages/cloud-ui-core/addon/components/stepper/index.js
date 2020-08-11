import Component from '@glimmer/component';

/**
 *
 * `Stepper` displays seqential actions to take and can display if they have
 *     been completed. This component yields a Step component.
 *
 *
 * ```
 * <Stepper as |SR|>
 *   <SR.Step @completed={{true}} as |S|>
 *     <S.StepLabel>
 *       Label
 *     </S.StepLabel>
 *     <S.StepContent>
 *       Content
 *     </S.StepContent>
 *   </SR.Step>
 *   <SR.Step as |S|>
 *     <S.StepLabel>
 *       Label
 *     </S.StepLabel>
 *     <S.StepContent>
 *       Content
 *     </S.StepContent>
 *   </SR.Step>
 * </Stepper>
 * ```
 *
 * @class Stepper
 * @yield {StepperStep} Step `Stepper::Step` component
 *
 */

export default class StepperComponent extends Component {}
