module.exports = [
  {
    source: '/waypoint/docs/kubernetes/:path*',
    destination: '/waypoint/docs/platforms/kubernetes/:path*',
    permanent: true,
  },
  {
    source: '/waypoint/docs/glossary',
    destination: '/waypoint/docs/resources/glossary',
    permanent: true,
  },
  {
    source: '/waypoint/docs/roadmap',
    destination: '/waypoint/docs/resources/roadmap',
    permanent: true,
  },
  {
    source: '/waypoint/docs/troubleshooting',
    destination: '/waypoint/docs/resources/troubleshooting',
    permanent: true,
  },
  {
    source: '/waypoint/docs/internals/:path*',
    destination: '/waypoint/docs/resources/internals/:path*',
    permanent: true,
  },
]
