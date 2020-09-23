import nextAuthApiRoute from 'lib/next-auth-utils/config'

export default (req, res) =>
  nextAuthApiRoute(
    req,
    res
  )({
    environments: { production: ['Okta'], preview: ['Auth0', 'Okta'] },
    pages: {
      error: '/signin-error', // Error code passed in query string as ?error=
    },
  })
