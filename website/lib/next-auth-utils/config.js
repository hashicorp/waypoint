import NextAuth from 'next-auth'
import NextAuthProviders from 'next-auth/providers'

function formatProviderConfig(ap) {
  const apName = ap.toUpperCase()
  const config = {
    clientId: process.env[`${apName}_CLIENT_ID`],
    clientSecret: process.env[`${apName}_CLIENT_SECRET`],
    domain: process.env[`${apName}_DOMAIN`],
  }
  return NextAuthProviders[ap](config)
}

export default (req, res) => ({ environments, pages }) =>
  NextAuth(req, res, {
    providers:
      environments[process.env.HASHI_ENV]?.map(formatProviderConfig) || [],
    pages,
  })
