import Route from '@ember/routing/route';

export default class OnboardingIndex extends Route {
  redirect() {
    return this.transitionTo('onboarding.install');
  }
}
