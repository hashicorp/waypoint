import styles from './HomePage.module.css'
import InfoGrid from 'components/info-grid'
import AnimatedStepsList from 'components/animated-steps-list'
import HomepageSection from 'components/homepage-section'
import HomepageHero from 'components/homepage-hero'
import BrandedCta from 'components/branded-cta'
import WaypointDiagram from 'components/waypoint-diagram'
import FeaturesList from 'components/features-list'

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
            url:
              'https://learn.hashicorp.com/collections/waypoint/getting-started',
            type: 'outbound',
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
              logos: require('./img/step-logos/build.svg'),
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
                    color: 'navy',
                    code: "│ [exporter] Adding layer 'ruby:ruby'",
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code: '│ [exporter] Adding 1/1 app layer(s)',
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code: "│ [exporter] Reusing layer 'launcher'",
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code: "│ [exporter] Reusing layer 'config'",
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code:
                      "│ [exporter] Adding label 'io.buildpacks.lifecycle.metadata'",
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code:
                      "│ [exporter] Adding label 'io.buildpacks.build.metadata'",
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code:
                      "│ [exporter] Adding label 'io.buildpacks.project.metadata'",
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code: '│ [exporter] *** Images (512c587cc97c):',
                  },
                  {
                    frames: 5,
                    color: 'navy',
                    code:
                      '│ [exporter]       index.docker.io/library/example-ruby:latest',
                  },
                  {
                    frames: 5,
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
              logos: require('./img/step-logos/deploy.svg'),
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
              logos: require('./img/step-logos/release.svg'),
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
                      'Deployment URL: https://immensely-guided-stag-5.alpha.waypoint.run',
                    color: 'white',
                  },
                ],
              },
            },
          ]}
        />
      </HomepageSection>

      <HomepageSection title="Features" theme="gray">
        <FeaturesList />
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
            url:
              'https://learn.hashicorp.com/collections/waypoint/getting-started',
            type: 'outbound',
          },
          { text: 'Explore documentation', url: '/docs' },
        ]}
      />
    </div>
  )
}
