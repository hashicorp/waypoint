import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 *
 *  <br />
 *  ##I. Introduction
 *  `<ModalDialog>` renders modal dialog content into the DOM element with the selector `.pdsModalDialogs` which is specified in `<Layout>`.  This is achieved through the use of the <a href="https://github.com/emberjs/ember.js/blob/master/packages/%40ember/-internals/glimmer/lib/syntax/in-element.ts" target="_blank">in-element</a> Ember helper.<br /><br />
 *
 *
 *  ##II. Accessibilty
 *  `<ModalDialog>` has been developed in consideration of the latest accessibility guidelines, including:
 *  <ul>
 *    <li>The modal dialog contains `role=dialog`</li>
 *    <li>When a user first opens the modal dialog, the `<header>` inside of the modal dialog is focused</li>
 *    <li>Users can cycle focus through the items inside of the modal dialog (i.e. input fields and buttons) by clicking the `tab` key. This is achieved through the <a href="https://josemarluedke.github.io/ember-focus-trap/" target="_blank">ember-focus-trap</a> helper.</li>
 *    <li>Only the content within the modal dialog are reachable by the user when the modal dialog is displayed.  This is achieved by setting all elements that are outside of the modal dialog to `inert` through the use of <a href="https://github.com/WICG/inert" target="_blank">wicg-inert polyfill</a></li>
 *  </ul>
 *
 *
 *  ##III. Ember Arguments
 *  `<ModalDialog>` accepts 3 arguments:
 *    <ol>
 *      <li>`@returnFocusTo`</li>
 *      <li>`@isActive`</li>
 *      <li>`@onActiveChange`</li>
 *      <li>`@variant`</li>
 *    </ol>
 *
 *
 *  ##IV. Ember Component Structure (Contextual Components)
 *  `<ModalDialog>` accepts 3 main contextual components:
 *    <ol>
 *      <li>`<ModalDialog::Header>`</li>
 *      <li>`<ModalDialog::Body>`</li>
 *      <li>`<ModalDialog::Footer>` which accepts two sub contextual components:<br />
          `<ModalDialog::Footer::Actions>` for actions passed down from the parent component.<br />
          `<ModalDialog::Footer::Cancel>` to render the cancel button that closes the `<ModalDialog>`.</br>
        </li>
 *    </ol>
 *
 *
 *  ##V. Ember Code Sample
 *
 *
 * ```
 <ModalDialog
   @returnFocusTo='modal-dialog-open-button'
   @isActive={{this.modalShowing}}
   @onActiveChange={{fn this.setModalShowing}}
   @variant='delete'
   as |MD|
 >
   <MD.Header>
     {{t 'components.page.hvns.delete.title' htfmlSafe=true}}
   </MD.Header>
   <MD.Body>
     {{t 'components.page.hvns.delete.description' modelName=@model.name htmlSafe=true}}
     <div class='hcpForm'>
       <div>
         <label for='name'>
           {{t 'components.page.hvns.delete.form.label.confirm-network-name'}}
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
         {{on 'click' @deleteHVN}}
       >
         {{t 'components.form.delete'}}
       </Button>
     </F.Actions>
     <F.Cancel>
       {{t 'components.form.cancel'}}
     </F.Cancel>
   </MD.Footer>
 </ModalDialog>
 * ```
 * <br /><br />
 * ##VI. Render Sample and Argument Definitions
 * @class ModalDialog
 *
 */

export default class ModalDialogComponent extends Component {
  /**
   * `@returnFocusTo` represents the id of the button that opens the modal dialog.
   *
   * This value is needed in order for `.focus()` to return to this button when the modal dialog is closed.
   *
   * You can pass-in dasherized string value, for example:<br /> `@returnFocusTo='some-dasherized-name'`
   * @argument @returnFocusTo
   * @type {$string}
   */

  /**
   * `@isActive` is either `true` or `false`.  This boolean determines whether or not the modal dialog will render to the DOM.  If `@isActive={{true}}`, the modal dialog will render to the DOM.
   * @argument @isActive
   * @type {?boolean}
   */

  /**
   * `@onActiveChange` is a function that recieves either `true` or `false`, this is needed to update the value of what gets passed for `isActive` because the value can be changed internally by pressing ESC.
   * @argument @onActiveChange
   * @type {?function}
   */

  /**
   * `@variant` represents the type of modal dialog theme you can choose from.
   *
   *  Currently the custom themes are, `delete`, `error`, or `edit`
   *
   *  Each theme customizes the following:
   <ul>
      <li>sets a maximum width on the modal dialog</li>
      <li>includes the corresponding icon for edit and delete</li>
      <li>styles the modal header</li>
      <li>sets a maximum width on the modal header title text</li>
    </ul>
   *
   *
   * `null` if no value is specified, the maximum modal width is 655px and no icon is displayed
   *
   *
   * `edit` maximum modal width is 500px
   *
   *
   * `delete` maximum modal width is 655px
   *
   *  *<i>In all cases, ellispis will appear in the header title if the text exceeds the width on the title</i>
   * @argument @variant
   * @type {$string}
   */

  get headerIconType() {
    if (this.args.variant === 'delete') {
      return 'alert-triangle';
    }

    if (this.args.variant === 'edit') {
      return 'edit';
    }

    if (this.args.variant === 'error') {
      return 'cancel-square-fill';
    }
    return '';
  }

  @action
  activeChanged(val) {
    if (this.args.onActiveChange && typeof this.args.onActiveChange === 'function') {
      this.args.onActiveChange(val);
    }
  }

  get modalDialogContainer() {
    let layoutModalElement = document.querySelector('.pdsApp .pdsModalDialogs');
    return layoutModalElement;
  }

  //actions are in alphabetical order
  @action
  closeModalDialog() {
    if (this.isDestroyed || this.isDestroying) return;
    this.activeChanged(false);
  }

  @action
  onKeyDown(event) {
    if (event.keyCode === 27) {
      this.closeModalDialog();
    } else {
      return;
    }
  }

  @action
  openModalDialog() {
    document.addEventListener('keydown', this.onKeyDown);
    this.toggleInert();
  }

  @action
  teardownModal() {
    let openButton = document.getElementById(this.args.returnFocusTo);
    this.toggleInert();
    document.removeEventListener('keydown', this.onKeyDown);
    if (openButton) {
      openButton.focus();
    }
  }

  @action
  toggleInert() {
    let appChildrenElements = [];

    //because .children returns an HTMLCollection - 🤯
    for (let child of document.querySelector('.pdsApp').children) {
      appChildrenElements.push(child);
    }

    let isActive = this.args.isActive;

    let landmarkElements = appChildrenElements.filter(child => child.className !== 'pdsModalDialogs'); //this grabs all the siblings to the modal dialog parent container
    landmarkElements.forEach(function(landmarkElement) {
      landmarkElement.inert = !!isActive;
    });
  }
}
