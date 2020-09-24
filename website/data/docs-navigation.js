// The root folder for this documentation category is `content/docs`
//
// - A string refers to the name of a file
// - A "category" value refers to the name of a directory
// - All directories must have an "index.mdx" file to serve as
//   the landing page for the category, or a "name" property to
//   serve as the category title in the sidebar

export default [
  {
    category: 'getting-started',
    content: [
      {
        category: 'docker-example-app',
        content: [
          'install-prereqs',
          'install-waypoint',
          'init-waypoint',
          'deploy-app',
          'update-app',
          'view-exec-app',
          'ui',
          'summary',
        ],
      },
      {
        category: 'k8s-example-app',
        content: [
          'install-prereqs',
          'install-waypoint',
          'init-waypoint',
          'deploy-app',
          'update-app',
          'view-exec-app',
          'logging-app',
          'ui',
          'summary',
        ],
      },
      {
        category: 'nomad-example-app',
        content: [
          'install-prereqs',
          'install-waypoint',
          'init-waypoint',
          'deploy-app',
          'update-app',
          'view-exec-app',
          'logging-app',
          'ui',
          'summary',
        ],
      },
    ],
  },
  'glossary',
  'troubleshooting',
  '-----------',
  'url',
  'logs',
  'exec',
  'config',
  'workspaces',
  '-----------',
  {
    category: 'entrypoint',
    content: ['disable'],
  },
]
