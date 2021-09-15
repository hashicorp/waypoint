export default {
  getStartedDescription:
    "Follow step-by-step tutorials on AWS, Azure, GCP, and localhost.",
  getStartedLinks: [
    {
      label: "Deploy to Docker",
      href: "https://learn.hashicorp.com/collections/waypoint/get-started-docker",
    },
    {
      label: "Deploy to Kubernetes",
      href: "https://learn.hashicorp.com/collections/waypoint/get-started-kubernetes",
    },
    {
      label: "Deploy to AWS",
      href: "https://learn.hashicorp.com/collections/waypoint/deploy-aws",
    },
    {
      label: "View all Waypoint tutorials",
      href: "https://learn.hashicorp.com/waypoint",
    },
  ],
  logo: (
    <img
      style={{ width: "140px" }}
      alt="Waypoint"
      src={require("./download-logo.svg")}
    />
  ),
  product: "waypoint",
  tutorialLink: {
    href: "https://learn.hashicorp.com/waypoint",
    label: "View Tutorials at HashiCorp Learn",
  },
  // HashiCorp officially supported package managers
  // NOTE: This information will be moving into the releases API
  packageManagers: [
    {
      label: "Homebrew",
      commands: [
        "brew tap hashicorp/tap",
        "brew install hashicorp/tap/waypoint",
      ],
      os: "darwin",
    },
    {
      label: "Ubuntu/Debian",
      commands: [
        "curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -",
        'sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"',
        "sudo apt-get update && sudo apt-get install waypoint",
      ],
      os: "linux",
    },
    {
      label: "CentOS/RHEL",
      commands: [
        "sudo yum install -y yum-utils",
        "sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo",
        "sudo yum -y install waypoint",
      ],
      os: "linux",
    },
    {
      label: "Fedora",
      commands: [
        "sudo dnf install -y dnf-plugins-core",
        "sudo dnf config-manager --add-repo https://rpm.releases.hashicorp.com/fedora/hashicorp.repo",
        "sudo dnf -y install waypoint",
      ],
      os: "linux",
    },
    {
      label: "Amazon Linux",
      commands: [
        "sudo yum install -y yum-utils",
        "sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo",
        "sudo yum -y install waypoint",
      ],
      os: "linux",
    },
  ],
};
