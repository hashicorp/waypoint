import Route from '@ember/routing/route';

export default class OnboardingInstallLinuxIndex extends Route {
  redirect() {
    return this.transitionTo('onboarding.install.linux.ubuntu');
  }
}
