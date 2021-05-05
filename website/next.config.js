const withHashicorp = require('@hashicorp/nextjs-scripts')
const redirects = require('./redirects')

// log out our primary environment variables for clarity in build logs
console.log(`HASHI_ENV: ${process.env.HASHI_ENV}`)
console.log(`NODE_ENV: ${process.env.NODE_ENV}`)

// add a X-Robots-Tag noindex HTTP header
// prevent indexing for tip.waypointproject.io
let customHeaders = []
const robotsHeader = { key: 'X-Robots-Tag', value: 'noindex' }
if (process.env.VERCEL_GIT_COMMIT_REF == 'main') {
  customHeaders.push(
    {
      source: '/',
      headers: [robotsHeader],
    },
    {
      source: '/:all*',
      headers: [robotsHeader],
    }
  )
}

module.exports = withHashicorp({
  defaultLayout: true,
  transpileModules: [
    'is-absolute-url',
    '@hashicorp/react-.*',
    '@hashicorp/versioned-docs',
  ],
})({
  redirects() {
    return redirects
  },
  headers() {
    return Promise.resolve(customHeaders)
  },
  svgo: { plugins: [{ removeViewBox: false }] },
  env: {
    HASHI_ENV: process.env.HASHI_ENV || 'development',
    SEGMENT_WRITE_KEY: '9mlIVayJbNtJW2EOdAFKHNKcdLAgEDlV',
    BUGSNAG_CLIENT_KEY: '98922c3298fff145a2d154ad2e6d4e6a',
    BUGSNAG_SERVER_KEY: '45f0129bdbe991d7fdcd0338a1a4f1d7',
  },
})
