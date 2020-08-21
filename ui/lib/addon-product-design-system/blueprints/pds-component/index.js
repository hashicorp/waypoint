const { classify, dasherize } = require('ember-cli-string-utils');

const PREFIX = 'pds';

module.exports = {
  description: 'Generates a PDS component',

  fileMapTokens(options) {
    let { locals } = options;
    return {
      __emberPrefix__: () => locals.emberPrefix,
      __sassPrefix__: () => PREFIX,
    };
  },

  files() {
    let _files = this._super.files.apply(this, arguments);

    return _files.filter(file => {
      // exclude dot files
      if (file.includes('/.') || file.startsWith('.')) {
        return false;
      }
      return true;
    });
  },

  locals(options) {
    // destructure options into local vars
    let { entity, taskOptions } = options;
    let { name } = entity;
    let { dummy } = taskOptions;

    // used for Ember component namespacing
    let emberPrefix = PREFIX;

    if (dummy) {
      emberPrefix = 'docs';
    }

    let classyModule = classify(name);
    let classyNamespace = classify(emberPrefix);

    let tagName = `${classyNamespace}::${classyModule}`;

    // <div class=".<%= cssClassName %>">
    let cssClassName = `${PREFIX}${classyModule}`;

    // @use '<%= sassModule %>';
    let sassModule = `${PREFIX}/components/${dasherize(name)}`;

    // export default <%= jsClass %> extends Component {}
    let jsClass = `${classyNamespace}${classyModule}`;

    let _locals = {
      PREFIX,
      cssClassName,
      emberPrefix,
      jsClass,
      sassModule,
      tagName,
      dummy,
    };

    return _locals;
  },
};
