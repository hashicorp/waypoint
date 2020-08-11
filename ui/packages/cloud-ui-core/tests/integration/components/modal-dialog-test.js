import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, triggerKeyEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import modalDialog, {
  HEADER_SELECTOR,
  CLOSE_BUTTON_SELECTOR,
  CANCEL_BUTTON_SELECTOR,
} from 'cloud-ui-core/test-support/pages/components/modal-dialog';
import { create } from 'ember-cli-page-object';

let component = create(modalDialog);

module('Integration | Component | Modal Dialog', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    await render(hbs`
      <button id="modal-dialog-open-button">Test</button>
      <div class="pdsApp">
        <div class="pdsModalDialogs"></div>
        <ModalDialog
          @returnFocusTo='modal-dialog-open-button'
          @isActive={{true}}
          @variant='delete' as |MD|
        >
          <MD.Header>
            Modal Title
          </MD.Header>
          <MD.Body>
            Some Content
          </MD.Body>
          <MD.Footer as |F|>
            <F.Actions>
              <button>some actions</button>
            </F.Actions>
            <F.Cancel>
              Cancel Button Text
            </F.Cancel>
          </MD.Footer>
        </ModalDialog>
      </div>
    `);
    assert.ok(modalDialog.isPresent, 'renders the modal dialog');
  });

  test('focus trap focuses on header and it cycles elements on tab', async function(assert) {
    assert.expect(2);

    await render(hbs`
      <div class="pdsApp">
      test
      <a></a>
      <div class="pdsModalDialogs"></div>
      <div class="pdsAppHeader"></div>
      <div class="pdsSidebar"></div>
      <div class="pdsAppContent"></div>
      <div class="pdsAppFooter"></div>
      <div class="pdsDrawer"></div>
      <div class="pdsModalDialogs"></div>
      <ModalDialog
        @returnFocusTo='modal-dialog-open-button'
        @isActive={{true}}
        @variant='delete' as |MD|
      >
          <MD.Header>
            Modal Title
          </MD.Header>
          <MD.Body>
            <input type='text' class='input-1' />
            <input type='text' class='input-2' />
          </MD.Body>
          <MD.Footer as |F|>
            <F.Actions>
              <button class="action-button">some actions</button>
            </F.Actions>
            <F.Cancel>
              Cancel Button Text
            </F.Cancel>
          </MD.Footer>
        </ModalDialog>
      </div>
    `);

    let TAB_KEY_CODE = 9;

    await triggerKeyEvent(document.querySelector(HEADER_SELECTOR), 'keydown', TAB_KEY_CODE);
    assert.ok(document.activeElement.matches(HEADER_SELECTOR), 'the header is focused');

    await triggerKeyEvent(document.querySelector(CANCEL_BUTTON_SELECTOR), 'keydown', TAB_KEY_CODE);
    assert.ok(document.activeElement.matches(CLOSE_BUTTON_SELECTOR), 'the close button has focus');
  });

  test('it sets sibling elements to inert', async function(assert) {
    await render(hbs`
      <div class="pdsApp">
      test
      <a></a>
      <div class="pdsModalDialogs"></div>
      <div class="pdsAppHeader"></div>
      <div class="pdsSidebar"></div>
      <div class="pdsAppContent"></div>
      <div class="pdsAppFooter"></div>
      <div class="pdsDrawer"></div>
      <div class="pdsModalDialogs"></div>
      <ModalDialog
        @returnFocusTo='modal-dialog-open-button'
        @isActive={{true}}
        @variant='delete' as |MD|
      >
          <MD.Header>
            Modal Title
          </MD.Header>
          <MD.Body>
            Some Content
          </MD.Body>
          <MD.Footer as |F|>
            <F.Actions>
              <button>some actions</button>
            </F.Actions>
            <F.Cancel>
              Cancel Button Text
            </F.Cancel>
          </MD.Footer>
        </ModalDialog>
      </div>
    `);
    let appModal = document.querySelector('.pdsModalDialogs');

    let appChildrenElements = [];

    //because .children returns an HTMLCollection - 🤯
    for (let child of document.querySelector('.pdsApp').children) {
      appChildrenElements.push(child);
    }

    let landmarkElements = appChildrenElements.filter(child => child.className !== 'pdsModalDialogs');

    let landmarkElementName;

    landmarkElements.forEach(function(landmarkElement) {
      if (landmarkElement.parentElement === appModal.parentElement) {
        if (landmarkElement.className) {
          landmarkElementName = landmarkElement.className;
        } else {
          landmarkElementName = landmarkElement.tagName;
        }
        assert.ok(landmarkElement.inert, `${landmarkElementName} is inert`);
      }
    });
  });

  test('it enables focus on header', async function(assert) {
    assert.expect(1);
    await render(hbs`
      <div class="pdsApp">
        <div class="pdsModalDialogs"></div>
        <ModalDialog
          @returnFocusTo='modal-dialog-open-button'
          @isActive={{true}}
          @variant='delete' as |MD|
        >
          <MD.Header>
            Modal Title
          </MD.Header>
          <MD.Body>
            Some Content
          </MD.Body>
          <MD.Footer as |F|>
            <F.Actions>
              <button>some actions</button>
            </F.Actions>
            <F.Cancel>
              Cancel Button Text
            </F.Cancel>
          </MD.Footer>
        </ModalDialog>
      </div>
    `);
    assert.ok(document.activeElement.matches(HEADER_SELECTOR), 'The header has focus');
  });

  test('it closes when escape button is pressed and open button is focused', async function(assert) {
    assert.expect(2);
    this.isActive = true;
    await render(hbs`
      <div class="pdsApp">
        <button id='modal-dialog-open-button'>Open Modal</button>
        <div class="pdsModalDialogs"></div>
        <ModalDialog
          @returnFocusTo='modal-dialog-open-button'
          @isActive={{this.isActive}}
          @onActiveChange={{fn (mut this.isActive)}}
          @variant='delete' as |MD|
        >
          <MD.Header>
            Modal Title
          </MD.Header>
          <MD.Body>
            Some Content
          </MD.Body>
          <MD.Footer as |F|>
            <F.Actions>
              <button>some actions</button>
            </F.Actions>
            <F.Cancel>
              Cancel Button Text
            </F.Cancel>
          </MD.Footer>
        </ModalDialog>
      </div>
    `);

    await triggerKeyEvent(document.querySelector(CLOSE_BUTTON_SELECTOR), 'keydown', 27);
    assert.notOk(component.isPresent, 'the modal does not exist');
    assert.ok(
      document.activeElement.matches('#modal-dialog-open-button'),
      'The open button is focused when escape is pressed'
    );
  });

  test('it closes when cancel button is pressed and open button is focused', async function(assert) {
    assert.expect(2);
    this.isActive = true;
    await render(hbs`
      <div class="pdsApp">
        <button id='modal-dialog-open-button'>Open Modal</button>
        <div class="pdsModalDialogs"></div>
        <ModalDialog
          @returnFocusTo='modal-dialog-open-button'
          @isActive={{this.isActive}}
          @onActiveChange={{fn (mut this.isActive)}}
          @variant='delete' as |MD|
        >
          <MD.Header>
            Modal Title
          </MD.Header>
          <MD.Body>
            Some Content
          </MD.Body>
          <MD.Footer as |F|>
            <F.Actions>
              <button>some actions</button>
            </F.Actions>
            <F.Cancel>
              Cancel Button Text
            </F.Cancel>
          </MD.Footer>
        </ModalDialog>
      </div>
    `);

    await component.cancel();
    assert.notOk(component.isPresent, 'the modal does not exist');
    assert.ok(
      document.activeElement.matches('#modal-dialog-open-button'),
      'The open button is focused when cancel is pressed'
    );
  });
});
