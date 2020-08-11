import Component from '@glimmer/component';
import { run } from '@ember/runloop';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 *
 * `LoadingElapsed` displays the time since being rendered on screen.
 *
 *
 * ```
 * <Loading::Elapsed></Loading::Elapsed>
 * ```
 *
 * @class LoadingElapsed
 *
 */

export default class LoadingElapsedComponent extends Component {
  /**
   *
   * An optional timestamp that will be used to calculate elapsed time since.
   * @argument startTime;
   * @type {?number}
   *
   */

  timeout = null;
  renderTime = new Date().getTime();

  @tracked elapsed = '--:--';

  @action
  startTimer() {
    this.timeout = this.poll();
  }

  @action
  destroyTimer() {
    clearTimeout(this.timeout);
  }

  poll() {
    return setTimeout(() => {
      run(() => {
        this.elapsed = this.getElapsed();
        this.timeout = this.poll();
      });
    }, 1000);
  }

  get startTime() {
    return this.args.startTime || this.renderTime;
  }

  /**
   * Get the elapsed time since the argument "startTime" or since render time.
   * @method Typography#getElapsed
   * @return {string} The elapsed time string.
   */
  getElapsed() {
    let updatedTime = new Date().getTime();
    let difference = updatedTime - this.startTime;
    let hours = Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    let minutes = Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60));
    let seconds = Math.floor((difference % (1000 * 60)) / 1000);
    hours = hours < 10 ? '0' + hours : hours;
    minutes = minutes < 10 ? '0' + minutes : minutes;
    seconds = seconds < 10 ? '0' + seconds : seconds;
    if (hours == 0) {
      return `${minutes}:${seconds}`;
    }
    return `${hours}:${minutes}:${seconds}`;
  }
}
