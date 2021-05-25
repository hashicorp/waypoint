import styles from './style.module.css'
import InfoGrid from 'components/info-grid'
import AnimatedStepsList from 'components/animated-steps-list'
import HomepageSection from 'components/homepage-section'
import HomepageHero from 'components/homepage-hero'
import BrandedCta from 'components/branded-cta'
import WaypointDiagram from 'components/waypoint-diagram'
import Features from '@hashicorp/react-stepped-feature-list'
import Terminal from 'components/terminal'

const loadingDots = ['', '.', '. .', '. . .']
const spinner = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏']

export default function HomePage() {
  return (
    <div className={styles.homePage}>
      <HomepageHero
        title="Build. Deploy. Release."
        subtitle="Waypoint provides a modern workflow to build, deploy, and release across platforms."
        description="Waypoint uses a single configuration file and common workflow to manage and observe deployments across platforms such as Kubernetes, Nomad, EC2, Google Cloud Run, and more."
        links={[
          {
            text: 'Get Started',
            url: '/docs/getting-started',
            type: 'inbound',
          },
        ]}
      />

      <HomepageSection theme="light">
        <AnimatedStepsList
          terminalHeroState={{
            frameLength: 100,
            loop: true,
            lines: [
              {
                frames: 5,
                code: ['$ waypoint up', '$ waypoint up |'],
              },
            ],
          }}
          steps={[
            {
              name: 'Build',
              description: (
                <>
                  <p>
                    Waypoint builds applications for any language or framework.
                    You can use Buildpacks for automatically building common
                    frameworks or custom Dockerfiles or other build tools for
                    more fine-grained control.
                  </p>
                  <p>
                    The build step is where your application and assets are
                    compiled, validated, and an artifact is created.
                  </p>
                  <p>
                    This artifact can be published to a remote registry or
                    simply passed to the deploy step.
                  </p>
                </>
              ),
              logos: [
                {
                  url: require('./img/step-logos/angular.svg'),
                  alt: 'Angular',
                },
                {
                  url: require('./img/step-logos/react.svg'),
                  alt: 'React',
                },
                {
                  url: require('./img/step-logos/ruby.svg'),
                  alt: 'Ruby',
                },
                {
                  url: require('./img/step-logos/python.svg'),
                  alt: 'Python',
                },
                {
                  url: require('./img/step-logos/go.svg'),
                  alt: 'Go',
                },
                {
                  url: require('./img/step-logos/nodejs.svg'),
                  alt: 'Node.js',
                },
                {
                  url: require('./img/step-logos/nextjs.svg'),
                  alt: 'Next.js',
                },
                {
                  url: require('./img/step-logos/dotnet.svg'),
                  alt: '.NET',
                },
                {
                  url: require('./img/step-logos/and-more.svg'),
                  alt: 'and More',
                },
              ],
              terminal: {
                frameLength: 100,
                loop: false,
                lines: [
                  {
                    color: 'white',
                    frames: 4,
                    code: loadingDots.map((dots) => `» Building ${dots}`),
                  },
                  {
                    frames: 10,
                    code:
                      'Creating new buildpack-based image using builder: buildpacks:18',
                    color: 'gray',
                  },
                  {
                    frames: 5,
                    code: '✓ Creating pack client',
                  },
                  {
                    frames: 2,
                    code: spinner
                      .concat(spinner)
                      .map((step) => `${step} Building image`)
                      .concat('  Building image'),
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code: "│ [exporter] Adding layer 'ruby:ruby'",
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code: '│ [exporter] Adding 1/1 app layer(s)',
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code: "│ [exporter] Reusing layer 'launcher'",
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code: "│ [exporter] Reusing layer 'config'",
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code:
                      "│ [exporter] Adding label 'io.buildpacks.lifecycle.metadata'",
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code:
                      "│ [exporter] Adding label 'io.buildpacks.build.metadata'",
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code:
                      "│ [exporter] Adding label 'io.buildpacks.project.metadata'",
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code: '│ [exporter] *** Images (512c587cc97c):',
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code:
                      '│ [exporter]       index.docker.io/library/example-ruby:latest',
                  },
                  {
                    frames: 5,
                    short: true,
                    color: 'navy',
                    code: "│ [exporter] Reusing cache layer 'ruby:gems'",
                  },
                  { code: '' },
                  {
                    frames: 5,
                    code: '✓ Injecting entrypoint binary to image',
                  },
                  { code: '' },
                  {
                    frames: 5,
                    code: 'Generated new Docker image: example-ruby:latest',
                    color: 'gray',
                  },
                  {
                    frames: 5,
                    code:
                      'Tagging Docker image: example-ruby:latest => gcr.io/example/example-ruby:latest',
                    color: 'gray',
                  },
                  {
                    frames: 5,
                    code:
                      'Docker image pushed: gcr.io/example/example-ruby:latest',
                    color: 'white',
                  },
                ],
              },
            },
            {
              name: 'Deploy',
              description: (
                <>
                  <p>
                    Waypoint deploys artifacts created by the build step to a
                    variety of platforms, from Kubernetes to EC2 to static site
                    hosts.
                  </p>
                  <p>
                    It configures your target platform and prepares the new
                    application version to be publicly accessible. Deployments
                    are accessible via a preview URL prior to release.
                  </p>
                </>
              ),
              logos: [
                {
                  url: require('./img/step-logos/kubernetes.svg'),
                  alt: 'Kubernetes',
                },
                {
                  url: require('./img/step-logos/nomad.svg'),
                  alt: 'Nomad',
                },
                {
                  url: require('./img/step-logos/amazon-ecs.svg'),
                  alt: 'Amazon ECS',
                },
                {
                  url: require('./img/step-logos/azure-container-service.svg'),
                  alt: 'Azure Container Service',
                },
                {
                  url: require('./img/step-logos/docker.svg'),
                  alt: 'Docker',
                },
                {
                  url: require('./img/step-logos/cloud-run.svg'),
                  alt: 'Google Cloud Run',
                },
                {
                  url: require('./img/step-logos/and-more.svg'),
                  alt: 'and More',
                },
              ],
              terminal: {
                frameLength: 100,
                loop: false,
                lines: [
                  {
                    color: 'white',
                    frames: 4,
                    code: loadingDots.map((dots) => `» Deploying ${dots}`),
                  },
                  {
                    frames: 5,
                    code:
                      '✓ Kubernetes client connected to https://kubernetes.example.com:6443',
                  },
                  {
                    frames: 2,
                    color: 'gray',
                    code: spinner
                      .concat(spinner)
                      .map((step) => `${step} Preparing deployment`)
                      .concat('✓ Created deployment'),
                  },
                  {
                    frames: 2,
                    color: 'gray',
                    code: spinner
                      .concat(spinner)
                      .map(
                        (step) =>
                          `${step} Waiting on deployment to become available: 1/1/0`
                      )
                      .concat('✓ Deployment successfully rolled out!'),
                  },
                  { code: '' },
                  {
                    frames: 2,
                    code:
                      '\nThe deploy was successful! A Waypoint deployment URL is shown below. This can be used internally to check your deployment and is not meant for external traffic. You can manage this hostname using "waypoint hostname"',
                    color: 'gray',
                  },
                  { code: '' },
                  {
                    frames: 1,
                    code:
                      '\nDeployment URL: https://immensely-guided-stag--v5.waypoint.run',
                    color: 'green',
                  },
                ],
              },
            },
            {
              name: 'Release',
              description: (
                <>
                  <p>
                    Waypoint releases your staged deployments and makes them
                    accessible to the public. This works by updating load
                    balancers, configuring DNS, etc. The exact behavior depends
                    on your target platform.
                  </p>
                  <p>
                    The release step is pluggable, enabling you to drop in
                    custom release logic such as blue/green, service mesh usage,
                    and more.
                  </p>
                </>
              ),
              logos: [
                {
                  url: require('./img/step-logos/aws.svg'),
                  alt: 'Amazon Web Services',
                },
                {
                  url: require('./img/step-logos/azure.svg'),
                  alt: 'Microsoft Azure',
                },
                {
                  url: require('./img/step-logos/gcp.svg'),
                  alt: 'Google Cloud Platform',
                },
                {
                  url: require('./img/step-logos/terraform.svg'),
                  alt: 'Terraform',
                },
                {
                  url: require('./img/step-logos/circleci.svg'),
                  alt: 'CircleCI',
                },
                {
                  url: require('./img/step-logos/slack.svg'),
                  alt: 'Slack',
                },
                {
                  url: require('./img/step-logos/github.svg'),
                  alt: 'Github',
                },
                {
                  url: require('./img/step-logos/gitlab.svg'),
                  alt: 'Gitlab',
                },
                {
                  url: require('./img/step-logos/and-more.svg'),
                  alt: 'and More',
                },
              ],
              terminal: {
                frameLength: 100,
                loop: false,
                lines: [
                  {
                    frames: 4,
                    color: 'white',
                    code: loadingDots.map((dots) => `» Releasing ${dots}`),
                  },
                  {
                    frames: 5,
                    code:
                      '✓ Kubernetes client connected to https://kubernetes.example.com:6443',
                  },
                  {
                    frames: 2,
                    color: 'gray',
                    code: spinner
                      .concat(spinner)
                      .map((step) => `${step} Preparing service`)
                      .concat('✓ Service is ready!'),
                  },
                  { code: '' },
                  {
                    frames: 4,
                    color: 'white',
                    code: loadingDots.map(
                      (dots) => `» Pruning old deployments ${dots}`
                    ),
                  },
                  {
                    frames: 5,
                    code: '  Deployment: 01EJCSFNDDD15P2BXBW2KCYVB2',
                    color: 'navy',
                  },
                  { code: '' },
                  {
                    frames: 5,
                    code: '\nThe release was successful!',
                    color: 'green',
                  },
                  { code: '' },
                  {
                    frames: 1,
                    code: '\nRelease URL: https://www.example.com',
                    color: 'green',
                  },
                ],
              },
            },
          ]}
        />
      </HomepageSection>

      <HomepageSection title="Features" theme="gray">
        <Features
          features={[
            {
              title: 'Application Logs',
              description:
                'View log output for running applications and deployments',
              learnMoreLink: '/docs/logs',
              content: (
                <Terminal
                  lines={[
                    { code: '$ waypoint logs' },
                    {
                      code: '[11] Puma starting in cluster mode...',
                      color: 'gray',
                    },
                    {
                      code:
                        '[11] * Version 5.0.2 (ruby 2.7.1-p83), codename: Spoony Bard',
                      color: 'gray',
                    },
                    {
                      code: '[11] * Min threads: 5, max threads: 5',
                      color: 'gray',
                    },
                    { code: '[11] * Environment: production', color: 'gray' },
                    { code: '[11] * Process workers: 2', color: 'gray' },
                    { code: '[11] * Preloading application', color: 'gray' },
                    {
                      code: '[11] * Listening on tcp://0.0.0.0:3000',
                      color: 'gray',
                    },
                    {
                      code:
                        'I, [2020-09-23T19:38:59.250971 #17] INFO -- : [936a952c-76b1-41f0-a4fe-ae2b77afc398] Started GET "/" for 10.36.5.1 at 2020-09-23 19:38:59 +0000',
                      color: 'gray',
                    },
                  ]}
                />
              ),
            },
            {
              title: 'Live Exec',
              description:
                'Execute a command in the context of a running application',
              learnMoreLink: '/docs/exec',
              content: (
                <Terminal
                  lines={[
                    { code: '$ waypoint exec bash' },
                    {
                      code: 'Connected to deployment v18',
                      color: 'gray',
                    },
                    {
                      code: '$ rake db:migrate',
                      color: 'white',
                    },
                  ]}
                />
              ),
            },
            {
              title: 'Preview URLs',
              description:
                'Get publicly accessible preview URLs per-deployment',
              learnMoreLink: '/docs/url',
              content: (
                <Terminal
                  lines={[
                    { code: '$ waypoint deploy' },
                    { code: '' },
                    { code: '» Deploying...', color: 'white' },
                    {
                      code: '✓ Deployment successfully rolled out!',
                      color: 'navy',
                    },
                    { code: '\n' },
                    { code: '» Releasing...', color: 'white' },
                    {
                      code: '✓ Service successfully configured!',
                      color: 'navy',
                    },
                    { code: '\n' },
                    {
                      code:
                        'The deploy was successful! A Waypoint URL is shown below.',
                      color: 'white',
                    },
                    { code: '\n' },
                    {
                      code:
                        '   Release URL: https://admittedly-poetic-joey.waypoint.run',
                      color: 'white',
                    },
                    {
                      code:
                        'Deployment URL: https://admittedly-poetic-joey--v18.waypoint.run',
                      color: 'white',
                    },
                  ]}
                />
              ),
            },

            {
              title: 'Web UI',
              description:
                'View projects and applications being deployed by Waypoint in a web interface',
              content: (
                <img
                  style={{
                    height: '500px',
                    width: 'auto',
                  }}
                  src={require('./img/waypoint_ui@3x.png')}
                  alt="Web UI"
                />
              ),
            },
            {
              title: 'CI/CD and Version Control Integration',
              description:
                'Integrate with existing CI/CD providers and version control providers like GitHub, CircleCI, Jenkins, GitLab, and more',
              learnMoreLink: '/docs/automating-execution',
              content: (
                <Terminal
                  title="config.yaml"
                  lines={[
                    {
                      code: 'env:',
                      color: 'white',
                    },
                    {
                      code:
                        '  WAYPOINT_SERVER_TOKEN: ${{ secrets.WAYPOINT_SERVER_TOKEN }}',
                      color: 'white',
                    },
                    {
                      code: '  WAYPOINT_SERVER_ADDR: waypoint.example.com:9701',
                      color: 'white',
                    },
                    {
                      code: 'steps:',
                      color: 'white',
                    },
                    {
                      code: '  - uses: actions/checkout@v2',
                      color: 'white',
                    },
                    {
                      code: '  - uses: hashicorp/action-setup-waypoint',
                      color: 'white',
                    },
                    {
                      code: '  with:',
                      color: 'white',
                    },
                    {
                      code: "    version: '0.1.0'",
                      color: 'white',
                    },
                    {
                      code: '- run: waypoint init',
                      color: 'white',
                    },
                    {
                      code: '- run: waypoint up',
                      color: 'white',
                    },
                  ]}
                />
              ),
            },
            {
              title: 'Extensible Plugin Interface',
              description:
                'Easily extend Waypoint with custom support for platforms, build processes, and release systems.',
              learnMoreLink: '/docs/extending-waypoint',
              content: (
                <Terminal
                  title="plugin.go"
                  lines={[
                    {
                      code: '// Destroy deletes the Nomad job.',
                    },
                    {
                      code: 'func (p *Platform) Destroy(',
                      color: 'white',
                    },
                    {
                      code: '  ctx context.Context,',
                      color: 'white',
                    },
                    {
                      code: '  log hclog.Logger,',
                      color: 'white',
                    },
                    {
                      code: '  deployment *Deployment,',
                      color: 'white',
                    },
                    {
                      code: '  ui terminal.UI,',
                      color: 'white',
                    },
                    {
                      code: ') error {',
                      color: 'white',
                    },
                    {
                      code:
                        '  client, err := api.NewClient(api.DefaultConfig())',
                      color: 'white',
                    },
                    {
                      code: '  if err != nil {',
                      color: 'gray',
                    },
                    {
                      code: '    return err',
                      color: 'gray',
                    },
                    {
                      code: '  }',
                      color: 'gray',
                    },
                    {
                      code: '  ',
                      color: 'gray',
                    },
                    {
                      code: '  st.Update("Deleting job...")',
                      color: 'gray',
                    },
                    {
                      code:
                        '  _, _, err = client.Jobs().Deregister(deployment.Id, true, nil)',
                      color: 'navy',
                    },
                    {
                      code: '  return err',
                      color: 'gray',
                    },
                    {
                      code: '}',
                      color: 'white',
                    },
                  ]}
                />
              ),
            },
          ]}
        />
      </HomepageSection>

      <HomepageSection title="Why Waypoint" theme="light">
        <InfoGrid
          items={[
            {
              icon: require('./img/why-waypoint/workflow-consistency.svg'),
              title: 'Consistency of workflows',
              description:
                'Consistent workflow for build, deploy, and release across platforms',
            },
            {
              icon: require('./img/why-waypoint/deployment-confidence.svg'),
              title: 'Confidence in deployments',
              description:
                'Validate deployments across environments with common tooling',
            },
            {
              icon: require('./img/why-waypoint/ecosystem-extensibility.svg'),
              title: 'Extensibility with the ecosystem',
              description:
                'Extend workflows via built-in plugins and an extensible interface',
            },
          ]}
        />
        <WaypointDiagram className={styles.whyWaypointDiagram} />
      </HomepageSection>

      <BrandedCta
        heading="Ready to get started?"
        content="Start by following a tutorial to deploy a simple application with Waypoint or learn about how the project works by exploring the documentation."
        links={[
          {
            text: 'Get Started',
            url: '/docs/getting-started',
            type: 'inbound',
          },
          { text: 'Explore documentation', url: '/docs' },
        ]}
      />
    </div>
  )
}
