// The root folder for this documentation category is `content/docs`
//
// - A string refers to the name of a file
// - A "category" value refers to the name of a directory
// - All directories must have an "index.mdx" file to serve as
//   the landing page for the category, or a "name" property to
//   serve as the category title in the sidebar

export default [
  {
    category: 'intro',
    content: [
      {
        category: 'vs',
        content: ['helm', 'paas', 'kubernetes'],
      },
    ],
  },
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
  {
    category: 'upgrading',
    content: ['compatibility', 'protocol-table'],
  },
  '-----------',
  {
    category: 'lifecycle',
    content: ['build', 'deploy', 'release', 'hooks'],
  },
  {
    category: 'waypoint-hcl',
    content: [
      'app',
      'build',
      'deploy',
      'hook',
      'plugin',
      'registry',
      'release',
      'url',
      'use',
    ],
  },
  {
    category: 'server',
    content: [
      'auth',
      {
        category: 'run',
        content: ['maintenance', 'security'],
      },
    ],
  },
  'url',
  'logs',
  'exec',
  'app-config',
  'workspaces',
  'plugins',
  '-----------',
  {
    category: 'entrypoint',
    content: ['disable'],
  },
  {
    category: 'automating-execution',
    content: ['github-actions', 'circle-ci'],
  },
  'troubleshooting',
  'glossary',
  '-----------',
  {
    category: 'internals',
    content: ['architecture', 'execution'],
  },
  {
    category: 'extending-waypoint',
    content: [
      'main-func',
      'passing-values',
      {
        category: 'plugin-interfaces',
        content: [
          'authenticator',
          'configurable',
          'configurable-notify',
          'builder',
          'registry',
          'platform',
          'release-manager',
          'destroy',
          'default-parameters',
        ],
      },
      'example-plugin',
    ],
  },
]
