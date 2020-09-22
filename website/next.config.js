const withHashicorp = require('@hashicorp/nextjs-scripts')
const redirects = require('./redirects')

// log out our primary environment variables for clarity in build logs
console.log(`HASHI_ENV: ${process.env.HASHI_ENV}`)
console.log(`NODE_ENV: ${process.env.NODE_ENV}`)
console.log(`OKTA_DOMAIN: ${process.env.OKTA_DOMAIN}`)

module.exports = withHashicorp({
  defaultLayout: true,
  transpileModules: ['is-absolute-url', '@hashicorp/react-mega-nav'],
})({
  redirects() {
    return redirects
  },
  env: {
    HASHI_ENV: process.env.HASHI_ENV || 'development',
    SEGMENT_WRITE_KEY: '9mlIVayJbNtJW2EOdAFKHNKcdLAgEDlV',
    BUGSNAG_CLIENT_KEY: '98922c3298fff145a2d154ad2e6d4e6a',
    BUGSNAG_SERVER_KEY: '45f0129bdbe991d7fdcd0338a1a4f1d7',
    NEXT_PUBLIC_OKTA_CLIENT_ID: process.env.OKTA_CLIENT_ID,
    NEXT_PUBLIC_OKTA_DOMAIN: process.env.OKTA_DOMAIN,
  },
})
