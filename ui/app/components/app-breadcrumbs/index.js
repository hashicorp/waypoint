import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { reads } from '@ember/object/computed';
import { tagName, classNames } from '@ember-decorators/component';
import classic from 'ember-classic-decorator';

@classic
@tagName('ul')
@classNames('breadcrumbs', 'pds-breadcrumbs')
export default class AppBreadcrumbs extends Component {
  @service('breadcrumbs') breadcrumbsService;

  @reads('breadcrumbsService.breadcrumbs') breadcrumbs;
}
