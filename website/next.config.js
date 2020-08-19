const withHashicorp = require('@hashicorp/nextjs-scripts')
const path = require('path')

module.exports = withHashicorp({
  defaultLayout: true,
  transpileModules: ['is-absolute-url', '@hashicorp/react-mega-nav'],
  mdx: { resolveIncludes: path.join(__dirname, 'pages/partials') },
})({
  experimental: { modern: true },
  env: {
    HASHI_ENV: process.env.HASHI_ENV || 'development',
    SEGMENT_WRITE_KEY: 'xxx',
    BUGSNAG_CLIENT_KEY: 'xxx',
    BUGSNAG_SERVER_KEY: 'xxx',
  },
})
