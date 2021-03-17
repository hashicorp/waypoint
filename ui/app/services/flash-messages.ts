import FlashMessages from 'ember-cli-flash/services/flash-messages';

export default class PdsFlashMessages extends FlashMessages {
  remove(id) {
    let message = this.queue.findBy('id', id);
    if (message) {
      message._teardown();
    }
  }
}
