/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ConfigVar, Project } from 'waypoint-pb';

import ApiService from 'waypoint/services/api';
import { BufferedChangeset } from 'ember-changeset/types';
import { Changeset } from 'ember-changeset';
import Component from '@glimmer/component';
import FlashMessagesService from 'waypoint/services/pds-flash-messages';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

type ConfigVarChangeset = Partial<BufferedChangeset> & ConfigVar.AsObject;

interface VariableArgs {
  variable: ConfigVar.AsObject;
  isEditing: boolean;
  isCreating: boolean;
  saveVariableSettings: (variable: ConfigVarChangeset, deleteVariable?: boolean) => Promise<Project.AsObject>;
  deleteVariable: (variable: ConfigVar.AsObject) => Promise<void>;
  cancelCreate: () => void;
}

export default class ProjectConfigVariablesListItemComponent extends Component<VariableArgs> {
  @service api!: ApiService;
  @service('pdsFlashMessages') flashMessages!: FlashMessagesService;

  initialVariable!: ConfigVar.AsObject;
  @tracked variable: ConfigVar.AsObject;
  @tracked changeset?: ConfigVarChangeset;
  @tracked isCreating: boolean;
  @tracked isEditing: boolean;

  constructor(owner: unknown, args: VariableArgs) {
    super(owner, args);
    let { variable, isEditing, isCreating } = args;
    this.variable = variable;
    this.isEditing = isEditing;
    this.isCreating = isCreating;

    if (this.isCreating || this.isEditing) {
      this.changeset = Changeset(this.variable) as ConfigVarChangeset;
    }
  }

  get isEditable(): boolean {
    // Temporarily making Dynamic vars uneditable
    return !this.variable.dynamic;
  }

  storeInitialVariable(): void {
    this.initialVariable = JSON.parse(JSON.stringify(this.variable));
  }

  @action
  async deleteVariable(variable: ConfigVar.AsObject): Promise<void> {
    await this.args.deleteVariable(variable);
  }

  @action
  editVariable(): void {
    this.isEditing = true;
    this.changeset = Changeset(this.variable) as ConfigVarChangeset;
  }

  @action
  async saveVariable(e: Event): Promise<void> {
    e.preventDefault();

    if (this.changeset === undefined) {
      return;
    }

    if (this.changeset.name === '' || this.changeset.pb_static === '') {
      this.flashMessages.error('Variable keys or values can not be empty');
      return;
    }
    if (this.variable.name !== this.changeset.name) {
      await this.args.saveVariableSettings(this.changeset, false);
      await this.args.deleteVariable(this.variable);
    } else {
      await this.args.saveVariableSettings(this.changeset, false);
    }
    this.changeset = undefined;
    this.isCreating = false;
    this.isEditing = false;
  }

  @action
  cancelCreate(): void {
    this.changeset = undefined;
    this.isCreating = false;
    this.isEditing = false;
    this.args.cancelCreate();
  }

  @action
  cancelEdit(): void {
    this.changeset = undefined;
    this.isCreating = false;
    this.isEditing = false;
  }
}
