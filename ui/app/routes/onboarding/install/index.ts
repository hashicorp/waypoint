/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { UAParser } from 'ua-parser-js';

export default class OnboardingInstallIndex extends Route {
  redirect(): void {
    let parser = new UAParser();

    switch (parser.getResult().os.name) {
      case 'Mac OS':
        this.transitionTo('onboarding.install.homebrew');
        return;
      // There isn't yet a chocolatey package for Waypoint
      // case 'Windows':
      //   return this.transitionTo('onboarding.install.chocolatey');
      case 'Debian':
      case 'Ubuntu':
        this.transitionTo('onboarding.install.linux.ubuntu');
        return;
      case 'CentOS':
        this.transitionTo('onboarding.install.linux.centos');
        return;
      case 'Fedora':
        this.transitionTo('onboarding.install.linux.fedora');
        return;
      default:
        this.transitionTo('onboarding.install.manual');
    }
  }
}
