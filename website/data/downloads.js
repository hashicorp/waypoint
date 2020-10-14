export const packageManagers = {
  homebrew: {
    label: 'Homebrew',
    url: '#',
    commands: ['brew tap hashicorp/tap', 'brew install hashicorp/tap/waypoint'],
  },
  ubuntu: {
    label: 'Ubuntu/Debian',
    commands: [
      'curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -',
      'sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"',
      'sudo apt-get update && sudo apt-get install waypoint',
    ],
  },
  centos: {
    label: 'CentOS/RHEL',
    commands: [
      'sudo yum install -y yum-utils',
      'sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo',
      'sudo yum -y install waypoint',
    ],
  },
  fedora: {
    label: 'Fedora',
    commands: [
      'sudo dnf install -y dnf-plugins-core',
      'sudo dnf config-manager --add-repo https://rpm.releases.hashicorp.com/fedora/hashicorp.repo',
      'sudo dnf -y install waypoint',
    ],
  },
  amazonLinux: {
    label: 'Amazon Linux',
    commands: [
      'sudo yum install -y yum-utils',
      'sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo',
      'sudo yum -y install waypoint',
    ],
  },
}

export const packageManagersByOs = {
  darwin: packageManagers.homebrew,
  linux: [
    packageManagers.ubuntu,
    packageManagers.centos,
    packageManagers.fedora,
    packageManagers.amazonLinux,
  ],
}

export const getStartedLinks = [
  {
    label: 'Deploy to Docker',
    href: 'https://learn.hashicorp.com/collections/waypoint/get-started-docker',
  },
  {
    label: 'Deploy to Kubernetes',
    href:
      'https://learn.hashicorp.com/collections/waypoint/get-started-kubernetes',
  },
  {
    label: 'Deploy to AWS',
    href: 'https://learn.hashicorp.com/collections/waypoint/deploy-aws',
  },
  {
    label: 'View all Waypoint tutorials',
    href: 'https://learn.hashicorp.com/waypoint',
  },
]
