const withHashicorp = require('@hashicorp/platform-nextjs-plugin')
const redirects = require('./redirects')

// log out our primary environment variables for clarity in build logs
console.log(`HASHI_ENV: ${process.env.HASHI_ENV}`)
console.log(`NODE_ENV: ${process.env.NODE_ENV}`)
console.log(`VERCEL_ENV: ${process.env.VERCEL_ENV}`)
console.log(`MKTG_CONTENT_API: ${process.env.MKTG_CONTENT_API}`)
console.log(`ENABLE_VERSIONED_DOCS: ${process.env.ENABLE_VERSIONED_DOCS}`)

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
  nextOptimizedImages: true,
})({
  webpack5: false,
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
    ENABLE_VERSIONED_DOCS: process.env.ENABLE_VERSIONED_DOCS || false,
  },
})
