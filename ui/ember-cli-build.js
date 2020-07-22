'use strict';

const EmberApp = require('ember-cli/lib/broccoli/ember-app');

const ENV = EmberApp.env();
const isProd = ENV.environment === 'production';
const isTest = ENV.environment === 'test';

module.exports = function(defaults) {
  let app = new EmberApp(defaults, {
    'ember-cli-favicon': {
      enabled: true,
      iconPath: 'favicon.png', // icon path related to `public` folder

      // See the [favicons](https://github.com/itgalaxy/favicons) module for details on the available configuration options.
      faviconsConfig: {
        // these options are passed directly to the favicons module
        path: "http://localhost:4200",
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
          yandex: false
        }
      }
    },
    autoImport: {
      // allows use of a CSP without 'unsafe-eval' directive
      forbidEval: true,
      // uncomment this to see a tree view of what ember-cli-auto-import builds
      //webpack: {
      //plugins: [
      //new (require("webpack-bundle-analyzer")).BundleAnalyzerPlugin(),
      //],
      //},
    },
  });

  app.import('node_modules/google-protobuf/google-protobuf.js', {
    using: [
      { transformation: 'cjs', as: 'google-protobuf' },
    ],
  });

  app.import('node_modules/google-protobuf/google/protobuf/any_pb.js', {
    using: [
      { transformation: 'cjs', as: 'google-protobuf/google/protobuf/any_pb.js' },
    ],
  });

  app.import('node_modules/google-protobuf/google/protobuf/timestamp_pb.js', {
    using: [
      { transformation: 'cjs', as: 'google-protobuf/google/protobuf/timestamp_pb.js' },
    ],
  });

  app.import('node_modules/google-protobuf/google/protobuf/empty_pb.js', {
    using: [
      { transformation: 'cjs', as: 'google-protobuf/google/protobuf/empty_pb.js' },
    ],
  });

  // There is a known issue in the CJS transform that forces you to 
  // only import from the node_modules path. For this reason we
  // make a few packages of basically vendored generated clients/messages
  // https://github.com/rwjblue/ember-cli-cjs-transform/issues/284
  // In the future this could be an ember add-on or see a fix upstream
  // and move them back to `vendor/`
  app.import('node_modules/api-common-protos/google/rpc/status_pb.js', {
    using: [
      { transformation: 'cjs', as: 'api-common-protos/google/rpc/status_pb.js' },
    ],
  });

  // app.import('node_modules/waypoint-client/ServerServiceClientPb.ts', {
  //   using: [
  //     { transformation: 'cjs', as: 'lib/waypoint-client/ServerServiceClientPb' },
  //   ],
  // });

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
