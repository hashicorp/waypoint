'use strict';

const browsers = ['last 1 Chrome versions', 'last 1 Firefox versions', 'last 1 Safari versions'];

const isCI = Boolean(process.env.CI);
const isProduction = process.env.EMBER_ENV === 'production';

module.exports = {
  browsers,
};
