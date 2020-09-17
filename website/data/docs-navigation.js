// The root folder for this documentation category is `pages/docs`
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
          'switching-builder',
          'update-app',
          'view-exec-app',
          'logging-app',
          'summary',
        ],
      },
    ],
  },
  'troubleshooting',
  '---',
  { title: 'External Link', href: 'https://www.hashicorp.com' },
]
