'use strict';

const EmberApp = require('ember-cli/lib/broccoli/ember-app');

const ENV = EmberApp.env();
const isProd = ENV.environment === 'production';
const isTest = ENV.environment === 'test';

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
    svg: {
      paths: ['node_modules/@hashicorp/structure-icons/dist', 'public/images', 'public/images/icons'],
      optimize: false,
    },
    svgJar: {
      sourceDirs: ['node_modules/@hashicorp/structure-icons/dist', 'public/images', 'public/images/icons'],
    },
    autoImport: {
      // allows use of a CSP without 'unsafe-eval' directive
      forbidEval: true,
      // uncomment this to see a tree view of what ember-cli-auto-import builds
      // webpack: {
      //   plugins: [new (require('webpack-bundle-analyzer').BundleAnalyzerPlugin)()],
      // },
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
  return app.toTree();
};
