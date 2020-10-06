// TODO: This file is entirely placeholder content

export const packageManagers = {
  homebrew: {
    label: 'Homebrew',
    url: '#',
    commands: ['brew install hashicorp/tap/waypoint'],
  },
  chocolatey: {
    label: 'Chocolatey',
    url: '#',
    commands: ['choco install waypoint'],
  },
  ubuntu: {
    label: 'Ubuntu/Debian',
    commands: ['command one', 'command two'],
  },
  centos: {
    label: 'CentOS/RHEL',
    commands: ['command one', 'command two'],
  },
  fedora: {
    label: 'Fedora',
    commands: ['command one'],
  },
  amazonLinux: {
    label: 'Amazon Linux',
    commands: ['command one', 'command two'],
  },
}

export const packageManagersByOs = {
  darwin: packageManagers.homebrew,
  windows: packageManagers.chocolatey,
  linux: [
    packageManagers.ubuntu,
    packageManagers.centos,
    packageManagers.fedora,
    packageManagers.amazonLinux,
  ],
}

export const containers = [
  {
    label: 'Docker',
    url: '#',
  },
]

export const tutorials = [
  {
    label: 'Windows',
    url: '#',
  },
  {
    label: 'macOS',
    url: '#',
  },
  {
    label: 'Linux',
    url: '#',
  },
]

export const getStartedLinks = [
  {
    label: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    href: '#1',
  },
  {
    label: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    href: '#2',
  },
  {
    label: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    href: '#3',
  },
]
