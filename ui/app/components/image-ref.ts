/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Ember from 'ember';
import Component from '@glimmer/component';
import { Ref } from 'docker-parse-image';
import { TaskGenerator, task, timeout } from 'ember-concurrency';

type Args = {
  imageRef?: Ref;
};

export default class extends Component<Args> {
  get uri(): string {
    if (!this.args.imageRef) {
      return '';
    }

    let { registry, namespace, repository } = this.args.imageRef;

    return [registry, namespace, repository].filter(Boolean).join('/');
  }

  get hasTag(): boolean {
    return !!this.args.imageRef?.tag;
  }

  get tagIsDigest(): boolean {
    return this.args.imageRef?.tag?.includes(':') ?? false;
  }

  get presentableTag(): string | undefined {
    let tag = this.args.imageRef?.tag;

    if (!tag) {
      return;
    }

    if (this.tagIsDigest) {
      let [alg, digest] = tag.split(':');
      return `${alg}:${digest.substr(0, 7)}`;
    }

    return tag;
  }

  @task({ restartable: true })
  *displayCopySuccess(): TaskGenerator<void> {
    let duration = Ember.testing ? 0 : 2000;
    yield timeout(duration);
  }
}
