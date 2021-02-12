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
  'getting-started',
  {
    category: 'upgrading',
    content: [
      'compatibility',
      'protocol-table',
      'release-notifications',
      {
        category: 'version-guides',
        content: ['0.2.0'],
      },
    ],
  },
  '-----------',
  {
    category: 'lifecycle',
    content: ['build', 'deploy', 'release', 'hooks'],
  },
  {
    category: 'waypoint-hcl',
    content: [
      {
        category: 'variables',
        content: ['artifact', 'deploy', 'entrypoint', 'path'],
      },
      {
        category: 'functions',
        content: ['all', 'template'],
      },
      {
        category: 'syntax',
        content: ['expressions', 'json'],
      },
      'app',
      'build',
      'config',
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
        content: ['maintenance', 'production', 'security'],
      },
    ],
  },
  'url',
  'logs',
  'exec',
  {
    category: 'app-config',
    content: ['dynamic'],
  },
  'workspaces',
  'plugins',
  '-----------',
  {
    category: 'extending-waypoint',
    content: [
      {
        category: 'creating-plugins',
        content: [
          'main',
          'configuration',
          'build-interface',
          'compiling',
          'example-application',
          'testing',
        ],
      },
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
    ],
  },
  {
    category: 'entrypoint',
    content: ['disable'],
  },
  {
    category: 'automating-execution',
    content: ['github-actions', 'gitlab-cicd', 'circle-ci', 'jenkins'],
  },
  'troubleshooting',
  'glossary',
  '-----------',
  {
    category: 'internals',
    content: ['architecture', 'execution'],
  },
  'roadmap',
]
