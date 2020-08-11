import hbs from 'htmlbars-inline-precompile';

export default {
  title: 'Stepper',
  component: 'Stepper',
};

export let Stepper = () => ({
  template: hbs`
    <Stepper as |SR|>
      <SR.Step @completed={{true}} as |S|>
        <S.StepLabel>
          Label
        </S.StepLabel>
        <S.StepContent>
          Content
        </S.StepContent>
      </SR.Step>
      <SR.Step as |S|>
        <S.StepLabel>
          Label
        </S.StepLabel>
        <S.StepContent>
          Content
        </S.StepContent>
      </SR.Step>
    </Stepper>
  `,
  context: {
  }
});
