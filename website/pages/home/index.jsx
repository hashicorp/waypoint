import styles from './HomePage.module.css'
import InfoGrid from 'components/info-grid'
import AnimatedStepsList from 'components/animated-steps-list'
import HomepageSection from 'components/homepage-section'
import HomepageHero from 'components/homepage-hero'
import BrandedCta from 'components/branded-cta'
import WaypointDiagram from 'components/waypoint-diagram'
import Features from 'components/features'
import Terminal from 'components/terminal'

export default function HomePage() {
  return (
    <div className={styles.homePage}>
      <HomepageHero
        title="Build. Deploy. Release."
        subtitle="Waypoint provides a modern workflow for build, deploy, and release across platforms."
        description="Waypoint does not run your software. It provides you a single configuration file and API to manage and observe deployments across environments and platforms, from your local workstation to your CI environment."
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
                    Waypoint builds applications for your language or framework,
                    from default compilation for common frameworks using
                    Buildpacks to fine grained control with custom Dockerfiles.
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
                    code: [
                      '» Building',
                      '» Building .',
                      '» Building . .',
                      '» Building . . .',
                    ],
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
                    code: [
                      '⠋ Building image',
                      '⠙ Building image',
                      '⠹ Building image',
                      '⠸ Building image',
                      '⠼ Building image',
                      '⠴ Building image',
                      '⠦ Building image',
                      '⠧ Building image',
                      '⠇ Building image',
                      '⠏ Building image',
                      '⠋ Building image',
                      '⠙ Building image',
                      '⠹ Building image',
                      '⠸ Building image',
                      '⠼ Building image',
                      '⠴ Building image',
                      '⠦ Building image',
                      '⠧ Building image',
                      '⠇ Building image',
                      '⠏ Building image',
                      '  Building image',
                    ],
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
                      'Tagging Docker image: example-ruby:latest => gcr.io/wp-dev-277323/example-ruby:latest',
                    color: 'gray',
                  },
                  {
                    frames: 5,
                    code:
                      'Docker image pushed: gcr.io/wp-dev-277323/example-ruby:latest',
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
                    variety of platforms, from Kubernetes to static site hosts.
                  </p>
                  <p>
                    It configures your target platform to be accessible and
                    starts the service, making it available for traffic at the
                    release stage.
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
                  url: require('./img/step-logos/netlify.svg'),
                  alt: 'Netlify',
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
                    code: [
                      '» Deploying',
                      '» Deploying .',
                      '» Deploying . .',
                      '» Deploying . . .',
                    ],
                  },
                  {
                    frames: 5,
                    color: 'white',
                    code:
                      '» Configuring https://kubernetes.docker.internal:6443 in namespace default',
                  },
                  {
                    frames: 2,
                    code: [
                      '⠋ Waiting on deployment to become available: 1/1/0',
                      '⠙ Waiting on deployment to become available: 1/1/0',
                      '⠹ Waiting on deployment to become available: 1/1/0',
                      '⠸ Waiting on deployment to become available: 1/1/0',
                      '⠼ Waiting on deployment to become available: 1/1/0',
                      '⠴ Waiting on deployment to become available: 1/1/0',
                      '⠦ Waiting on deployment to become available: 1/1/0',
                      '⠧ Waiting on deployment to become available: 1/1/0',
                      '⠇ Waiting on deployment to become available: 1/1/0',
                      '⠏ Waiting on deployment to become available: 1/1/0',
                      '⠋ Waiting on deployment to become available: 1/1/0',
                      '⠙ Waiting on deployment to become available: 1/1/0',
                      '⠹ Waiting on deployment to become available: 1/1/0',
                      '⠸ Waiting on deployment to become available: 1/1/0',
                      '⠼ Waiting on deployment to become available: 1/1/0',
                      '⠴ Waiting on deployment to become available: 1/1/0',
                      '⠦ Waiting on deployment to become available: 1/1/0',
                      '⠧ Waiting on deployment to become available: 1/1/0',
                      '⠇ Waiting on deployment to become available: 1/1/0',
                      '⠏ Waiting on deployment to become available: 1/1/0',
                      '✓ Waiting on deployment to become available: 1/1/0',
                    ],
                  },
                  {
                    frames: 10,
                    code: '✓ Deployment successfully rolled out!',
                  },
                ],
              },
            },
            {
              name: 'Release',
              description: (
                <>
                  <p>
                    Your deployment, a running copy of the application you built
                    and stored the artifact for, is released with a deployment
                    specific routable URL for previews.
                  </p>
                  <p>
                    In addition, if your application is configured with a
                    release step, it will automatically graduate or make the
                    release available based on an extensible plugin interface.
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
                    code: [
                      '» Releasing',
                      '» Releasing .',
                      '» Releasing . .',
                      '» Releasing . . .',
                    ],
                  },
                  {
                    frames: 5,
                    code: '✓ Service successfully configured!',
                  },
                  { code: '' },
                  {
                    frames: 4,
                    color: 'white',
                    code: [
                      '» Pruning old deployments',
                      '» Pruning old deployments .',
                      '» Pruning old deployments . .',
                      '» Pruning old deployments . . .',
                    ],
                  },
                  {
                    frames: 5,
                    code: 'Deployment: 01EJCSFNDDD15P2BXBW2KCYVB2',
                    color: 'navy',
                  },
                  { code: '' },
                  {
                    frames: 5,
                    code:
                      'The deploy was successful! A Waypoint deployment URL is shown below. This can be used internally to check your deployment and is not meant for external traffic. You can manage this hostname using "waypoint hostname"',
                    color: 'gray',
                  },
                  { code: '' },
                  {
                    frames: 1,
                    code: 'Release URL: https://www.example.com',
                    color: 'white',
                  },
                  {
                    frames: 1,
                    code:
                      'Deployment URL: https://immensely-guided-stag--v5.waypoint.run',
                    color: 'white',
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
              learnMoreLink: 'https://waypointproject.io/docs/logs',
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
                        '[11] * Version 3.11.2 (ruby 2.6.6-p146), codename: Love Song',
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
              content: (
                <Terminal
                  lines={[
                    { code: '$ waypoint exec bash' },
                    {
                      code: 'Connected to deployment v18',
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
              learnMoreLink: 'https://waypointproject.io/docs/url',
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
                    { code: '' },
                    { code: '» Releasing...', color: 'white' },
                    {
                      code: '✓ Service successfully configured!',
                      color: 'navy',
                    },
                    { code: '' },
                    { code: '» Pruning old deployments...', color: 'white' },
                    {
                      code: 'Deployment: 01EJYN2P2FG4N9CTXAYGGZW9W0',
                      color: 'white',
                    },
                    { code: '' },
                    {
                      code:
                        'The deploy was successful! A Waypoint deployment URL is shown below.',
                      color: 'white',
                    },
                    { code: '' },
                    {
                      code:
                        'Release URL: https://admittedly-poetic-joey.waypoint.run',
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
                  style={{ border: '1px solid rgba(174,176,183,.45)' }}
                  src={require('./img/web-ui.png')}
                  alt="Web UI"
                />
              ),
            },
            {
              title: 'CI/CD and Version Control Integration',
              description:
                'Integrate easily with existing CI/CD providers and version control providers like GitHub',
              learnMoreLink: '/docs/automating-execution/github-actions',
              content: (
                <Terminal
                  title="config.yaml"
                  lines={[
                    {
                      code: 'env:',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code:
                        'WAYPOINT_SERVER_TOKEN: ${{ secrets.WAYPOINT_SERVER_TOKEN }}',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: 'WAYPOINT_SERVER_ADDR: waypoint.example.com:9701',
                      color: 'white',
                    },
                    {
                      code: 'steps:',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: '- uses: actions/checkout@v2',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: '- uses: hashicorp/actions-setup-waypoint',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: 'with:',
                      color: 'white',
                    },
                    {
                      indent: 2,
                      code: "version: '0.1.0'",
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
                      indent: 1,
                      code: 'ctx context.Context,',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: 'log hclog.Logger,',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: 'deployment *Deployment,',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: 'ui terminal.UI,',
                      color: 'white',
                    },
                    {
                      code: ') error {',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: '',
                    },
                    {
                      indent: 1,
                      code: 'client, err := api.NewClient(api.DefaultConfig())',
                      color: 'white',
                    },
                    {
                      indent: 1,
                      code: 'if err != nil {',
                      color: 'gray',
                    },
                    {
                      indent: 2,
                      code: 'return err',
                      color: 'gray',
                    },
                    {
                      indent: 1,
                      code: '}',
                      color: 'gray',
                    },
                    {
                      indent: 1,
                      code: '',
                      color: 'gray',
                    },
                    {
                      indent: 1,
                      code: 'st.Update("Deleting job...")',
                      color: 'gray',
                    },
                    {
                      indent: 1,
                      code:
                        '_, _, err = client.Jobs().Deregister(deployment.Id, true, nil)',
                      color: 'navy',
                    },
                    {
                      indent: 1,
                      code: 'return err',
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
        content="Explore Waypoint documentation to deploy a simple application."
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
