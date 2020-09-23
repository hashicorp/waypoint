import NextAuth from 'next-auth'
import NextAuthProviders from 'next-auth/providers'

function getSiteProvidersForEnvironment() {
  let providers = []
  switch (process.env.HASHI_ENV) {
    case 'production':
      providers = ['Okta']
      break
    case 'preview':
      providers = ['Auth0', 'Okta']
      break
    case 'development':
      break
    default:
      break
  }
  return providers
}

export const customAuthPages = {
  pages: {
    error: '/signin-error', // Error code passed in query string as ?error=
  },
}

export const siteAuthProviders = {
  providers: getSiteProvidersForEnvironment().map((ap) =>
    NextAuthProviders[ap](formatProviderConfig(ap))
  ),
}

function formatProviderConfig(ap) {
  const apName = ap.toUpperCase()
  const config = {
    clientId: process.env[`${apName}_CLIENT_ID`],
    clientSecret: process.env[`${apName}_CLIENT_SECRET`],
    domain: process.env[`${apName}_DOMAIN`],
  }
  console.log(config)
  return config
}

export default (req, res) =>
  NextAuth(req, res, { ...siteAuthProviders, ...customAuthPages })
