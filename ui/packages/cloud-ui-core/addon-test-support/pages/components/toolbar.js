import { isPresent } from 'ember-cli-page-object';

export default {
  rendersToolbar: isPresent(' [data-test-toolbar] '),
  rendersActions: isPresent(' [data-test-toolbar-actions] '),
  rendersFilters: isPresent(' [data-test-toolbar-filters] '),
};
