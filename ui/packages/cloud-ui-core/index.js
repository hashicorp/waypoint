'use strict';
let Funnel = require('broccoli-funnel');
let mergeTrees = require('broccoli-merge-trees');
let path = require('path');
//let stew = require('broccoli-stew');

module.exports = {

  name: require('./package').name,
  isDevelopingAddon() {
    return true;
  },

  included: function(/* app */) {
    this._super.included.apply(this, arguments);
  },

  treeForStyles: function() {
    let PDSStyles = path.join(path.dirname(require.resolve('product-design-system')), `app/styles/pds`);
    let stylesPath = path.join(__dirname, `addon/styles`);
    let PDS = new Funnel(PDSStyles, {
      destDir: 'pds',
      annotation: 'PDS',
    });

    let Core = new Funnel(stylesPath, {
      destDir: 'app/styles/hcp',
      annotation: 'hcp styles from cloud-ui-core',
      exclude: ['addon.scss']
    });
    let mergedTrees =  mergeTrees([PDS, Core], { overwrite: true });

    // uncomment this line to see DEBUG output folder
    //return stew.debug(mergedTrees);
    return mergedTrees;

  },
};
