import { fillable } from 'ember-cli-page-object';

export default {
  confirm: fillable('[ data-test-modal-confirm ]'),
  async confirmDelete() {
    await this.confirm('DELETE');
  },
};
