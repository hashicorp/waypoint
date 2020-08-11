import { attribute, collection, fillable, isPresent, property, text, value } from 'ember-cli-page-object';

export default {
  isRendered: isPresent('[ data-test-select ]'),
  fill: fillable('[ data-test-select ]'),
  value: value('[ data-test-select ]'),
  options: collection('[ data-test-select-option ]', {
    selected: property('selected'),
    value: attribute('value'),
    label: text(),
  }),
  selectedOption() {
    return this.options.findOneBy('selected', true);
  },
};
