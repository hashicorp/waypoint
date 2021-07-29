import Route from '@ember/routing/route';

export default class OnboardingIndex extends Route {
  redirect(): void {
    this.transitionTo('onboarding.install');
  }
}
