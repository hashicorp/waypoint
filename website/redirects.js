module.exports = [
  // This is an example redirect, it can be removed once other redirects have been added
  {
    source: '/home',
    destination: '/',
    permanent: true,
  },
  {
    source: '/docs/kubernetes/:path*',
    destination: '/docs/platforms/kubernetes/:path*',
    permanent: true,
  },
  {
    source: '/docs/glossary',
    destination: '/docs/resources/glossary',
    permanent: true,
  },
  {
    source: '/docs/roadmap',
    destination: '/docs/resources/roadmap',
    permanent: true,
  },
  {
    source: '/docs/troubleshooting',
    destination: '/docs/resources/troubleshooting',
    permanent: true,
  },
  {
    source: '/docs/internals/:path*',
    destination: '/docs/resources/internals/:path*',
    permanent: true,
  },
]
