import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select } from '@storybook/addon-knobs';

export default {
  title: 'Form Control Error',
  component: 'FormControlError',
  decorators: [withKnobs],
};

export let FormControlError = () => ({
  template: hbs`
    <form class='hcpForm'>
      <FormControlError id='name_error'>
        <:message>cannot be blank</:message>
      </FormControlError>
    </form>`,
  context: {},
});

let ERRORS = [true, false];

//Example of how FormControlError appears within an hcpForm
export let FormControlErrorInAnHcpCreateForm = () => ({
  template: hbs`

  <Box @padding='xs 2xl'>
    <FlexGrid as |G|>
      <G.Item @lg='8' @md='12'>
        <div class='hcpForm'>
          <section>
            <div>
              <label for='name'>
                {{t 'components.page.hvns.create.form.label.network-name'}}
              </label>
              <Input
                class="{{if errors 'error'}}"
                type='text'
                id='name'
                name='name'
                @value={{this.name}}
                aria-describedby='name_error'
                data-test-network-name
              />
              {{#if errors}}
                <FormControlError id='name_error'>
                  <:message>cannot be blank</:message>
                </FormControlError>
              {{/if}}
            </div>
          </section>
        </div>
      </G.Item>
    </FlexGrid>
  </Box>
`,
  context: {
    errors: select('errors', ERRORS, false),
  },
});
