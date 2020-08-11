import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { guidFor } from '@ember/object/internals';
import { task } from 'ember-concurrency-decorators';
import { timeout } from 'ember-concurrency';
import { COPY_ICON_TYPE, COPY_ICON_CLASS, COPIED_ICON_TYPE, COPIED_ICON_CLASS } from './consts';

/**
 *
 * `CodeBlock` wraps code or terminal examples with the proper styling.
 *
 *
 * ```
 * <CodeBlock>rm -rf .</CodeBlock>
 * ```
 *
 * @class CodeBlock
 *
 */

export default class CodeBlockComponent extends Component {
  @tracked codeBlockId = guidFor(this);
  @tracked copyIconType = COPY_ICON_TYPE;
  @tracked copyIconClass = COPY_ICON_CLASS;

  @task
  *copied() {
    let originalIconType = this.copyIconType;
    let originalIconClass = this.copyIconClass;

    this.copyIconType = COPIED_ICON_TYPE;
    this.copyIconClass = COPIED_ICON_CLASS;

    yield timeout(1000);

    this.copyIconType = originalIconType;
    this.copyIconClass = originalIconClass;
  }
}
