import Route from '@ember/routing/route';
import { UAParser } from 'ua-parser-js';

export default class OnboardingInstallIndex extends Route {
  redirect() {
    let parser = new UAParser();

    switch (parser.getResult().os.name) {
      case 'Mac OS':
        return this.transitionTo('onboarding.install.homebrew');
      // There isn't yet a chocolatey package for Waypoint
      // case 'Windows':
      //   return this.transitionTo('onboarding.install.chocolatey');
      case 'Debian':
      case 'Ubuntu':
        return this.transitionTo('onboarding.install.linux.ubuntu');
      case 'CentOS':
        return this.transitionTo('onboarding.install.linux.centos');
      case 'Fedora':
        return this.transitionTo('onboarding.install.linux.fedora');
      default:
        return this.transitionTo('onboarding.install.manual');
    }
  }
}
