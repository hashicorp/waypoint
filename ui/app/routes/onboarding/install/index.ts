import Route from '@ember/routing/route';

export default class OnboardingInstallIndex extends Route {
  redirect() {
    return this.transitionTo('onboarding.install.manual');
  }
}
