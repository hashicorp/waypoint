/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

'use strict';

const EmberApp = require('ember-cli/lib/broccoli/ember-app');
const webpack = require('webpack');

module.exports = function (defaults) {
  let app = new EmberApp(defaults, {
    babel: {
      sourceMaps: 'inline',
    },
    'ember-cli-favicon': {
      enabled: true,
      iconPath: 'favicon.png', // icon path related to `public` folder

      // See the [favicons](https://github.com/itgalaxy/favicons) module for details on the available configuration options.
      faviconsConfig: {
        // these options are passed directly to the favicons module
        path: '/',
        appName: 'Waypoint',
        appShortName: 'WP',
        developerName: 'HashiCorp',
        appleStatusBarStyle: 'black',
        icons: {
          favicons: true,
          android: true,
          appleIcon: true,
          firefox: true,
          windows: true,
          coast: false,
          appleStartup: false,
          yandex: false,
        },
      },
    },
    'ember-simple-auth': {
      useSessionSetupMethod: true,
    },
    sassOptions: {
      precision: 4,
      includePaths: ['./node_modules/@hashicorp/design-system-tokens/dist/products/css'],
    },
    svgJar: {
      sourceDirs: ['node_modules/@hashicorp/structure-icons/dist', 'public/images'],
    },
    autoImport: {
      // allows use of a CSP without 'unsafe-eval' directive
      forbidEval: true,

      webpack: {
        plugins: [
          // required for ansi-colors to work in the browser
          new webpack.ProvidePlugin({ process: 'process/browser' }),

          // uncomment this to see a tree view of what ember-cli-auto-import builds
          // new (require('webpack-bundle-analyzer').BundleAnalyzerPlugin)(),
        ],
      },
      skipBabel: [
        {
          package: 'waypoint-client',
          semverRange: '*',
        },
        {
          package: 'waypoint-pb',
          semverRange: '*',
        },
      ],
    },
  });

  // Use `app.import` to add additional libraries to the generated
  // output files.
  //
  // If you need to use different assets in different
  // environments, specify an object as the first parameter. That
  // object's keys should be the environment name and the values
  // should be the asset to use in that environment.
  //
  // If the library that you are including contains AMD or ES6
  // modules that you would like to import into your application
  // please specify an object with the list of modules as keys
  // along with the exports of each module as its value.

  // xterm.js styles https://xtermjs.org/
  app.import('node_modules/xterm/css/xterm.css');

  return app.toTree();
};
