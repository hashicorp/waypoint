import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import stepperPageObject from 'cloud-ui-core/test-support/pages/components/stepper';

const stepper = create(stepperPageObject);

module('Integration | Component | Stepper', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <Stepper as |SR|>
        <SR.Step @completed={{true}} as |S|>
          <S.StepLabel>
            Complete Label
          </S.StepLabel>
          <S.StepContent>
            Content
          </S.StepContent>
        </SR.Step>
        <SR.Step as |S|>
          <S.StepLabel>
            Incomplete Label
          </S.StepLabel>
          <S.StepContent>
            Content
          </S.StepContent>
        </SR.Step>
      </Stepper>
    `);

    assert.ok(stepper.showsContainer, `the container renders`);
    assert.ok(stepper.steps[0].isIconComplete, `the first step icon is complete`);
    assert.ok(stepper.steps[0].isLabelComplete, `the first step label is complete`);
    assert.equal(stepper.steps[0].label, 'Complete Label', `the label renders and yields`);
    assert.equal(stepper.steps[0].content, 'Content', `the content renders and yields`);
    assert.ok(stepper.steps[1].isIconIncomplete, `the second step icon is incomplete`);
    assert.ok(stepper.steps[1].isLabelIncomplete, `the second step label is incomplete`);
    assert.equal(stepper.steps[1].label, 'Incomplete Label', `the label renders and yields`);
    assert.equal(stepper.steps[1].content, 'Content', `the content renders and yields`);
  });
});
