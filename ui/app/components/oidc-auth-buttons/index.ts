import Component from '@glimmer/component';
import { ListOIDCAuthMethodsResponse } from 'waypoint-pb';
import { tracked } from '@glimmer/tracking';

interface OIDCAuthButtonsArgs {
  model: ListOIDCAuthMethodsResponse.AsObject;
}

export default class OIDCAuthButtonsComponent extends Component<OIDCAuthButtonsArgs> {
  @tracked model!: ListOIDCAuthMethodsResponse.AsObject;
}
