import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select, text } from '@storybook/addon-knobs';

export default {
  title: 'ModalDialog',
  component: 'ModalDialog',
  decorators: [withKnobs],
};

let IS_ACTIVE = [true, false];
let VARIANT = [null, 'delete', 'edit', 'error'];
export let ModalDialog = () => ({
  template: hbs`
    <div class='pdsApp'>
      <Button
        @variant='warning'
        id={{returnFocusTo}}
        type='button'
      > Open Modal Dialog
      </Button>
      <div class='pdsModalDialogs'></div>
      <ModalDialog
        @returnFocusTo={{returnFocusTo}}
        @isActive={{isActive}}
        @variant={{variant}}
        as |MD|
      >
        <MD.Header>
          Delete HashiCorp Virtual Network?
        </MD.Header>
        <MD.Body>
          Deleting will delete this network and its related peering connections. Any resources associated must have already been destroyed before this network can be deleted.<br /><br />
          To recreate this network, you will need to create a new HashiCorp Virtual Network and reconfigure any peering connections needed.<br /><br />
          <div class='hcpForm'>
            <div>
              <label for='name'>
                Confirm HVN Name:
              </label>
              <Input type='text' id='name' name='name' data-test-network-name />
            </div>
          </div>
        </MD.Body>
        <MD.Footer as |F|>
          <F.Actions>
            <Button
              aria-label='delete hvn network'
              @variant='warning'
              {{on 'click' modalDialogAction}}
            >
              Delete
            </Button>
          </F.Actions>
          <F.Cancel>
            Cancel
          </F.Cancel>
        </MD.Footer>
      </ModalDialog>
    </div>
  `,

  context: {
    returnFocusTo: text('@returnFocusTo', 'some-dasherized-name'),
    isActive: select('@isActive', IS_ACTIVE, true),
    variant: select('@variant', VARIANT, 'delete'),
    modalDialogAction: function() {
      alert('what is up, yo?!');
    },
  },
});
