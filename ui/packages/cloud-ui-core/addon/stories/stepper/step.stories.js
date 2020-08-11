import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select } from '@storybook/addon-knobs';

export default {
  title: 'Stepper/Step',
  component: 'StepperStep',
  decorators: [withKnobs],
};

export let StepperStep = () => ({
  template: hbs`
    <Stepper::Step @completed={{completed}} as |S|>
      <S.StepLabel>
        Label
      </S.StepLabel>
      <S.StepContent>
        Content
      </S.StepContent>
    </Stepper::Step>
  `,
  context: {
    completed: select('Completed', [true, false], false),
  }
});