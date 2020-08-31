'use strict';

const EmberAddon = require('ember-cli/lib/broccoli/ember-addon');

module.exports = function(defaults) {
  let app = new EmberAddon(defaults, {
    // See https://github.com/adopted-ember-addons/ember-cli-sass
    sassOptions: {
      // TBD...
    },
    'ember-cli-storybook': {
      enableAddonDocsIntegration: true,
      componentFilePathPatterns: [
        'addon/components/**/*.js'
      ],
    },
    hinting: false,
    svgJar: {
      optimizer: false, // prevents optimizing pre-optimized icons
      sourceDirs: [
        '../../node_modules/@hashicorp/structure-icons/dist',
        './tests/dummy/public',
      ],
    },
  });

  /*
    This build file specifies the options for the dummy test app of this
    addon, located in `/tests/dummy`
    This build file does *not* influence how the addon or the app using it
    behave. You most likely want to be modifying `./index.js` or app's build file
  */

  return app.toTree();
};
