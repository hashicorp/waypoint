import styles from './HomePage.module.css'
import InfoGrid from 'components/info-grid'
import AnimatedStepsList from 'components/animated-steps-list'
import HomepageSection from 'components/homepage-section'
import HomepageHero from 'components/homepage-hero'
import BrandedCta from 'components/branded-cta'

export default function HomePage() {
  return (
    <div className={styles.homePage}>
      <HomepageHero
        title="Build. Deploy. Release."
        subtitle="Waypoint provides a modern workflow for deploying your development code on your development platform."
        description="Waypoint does not run your software. It provides a single configuration file and API to manage and observe deployments across environments and platforms, from your local workstation to your CI environment."
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
            },
          ]}
        />
      </HomepageSection>

      <HomepageSection title="Features" theme="gray"></HomepageSection>

      <HomepageSection title="Why Waypoint" theme="light">
        <InfoGrid
          items={[
            {
              icon: require('./img/info.svg'),
              title: 'Confidence in deployments',
              description:
                'Validate deployments across distinct and complex environments with common tooling',
            },
            {
              icon: require('./img/info.svg'),
              title: 'Consistency of workflows',
              description:
                'A consistent workflow for build, deploy, and release across common developer platforms',
            },
            {
              icon: require('./img/info.svg'),
              title: 'Extensibility with the ecosystem',
              description:
                'Extend workflows across the ecosystem via built-in plugins and an extensible interface',
            },
          ]}
        />
        <img
          className={styles.whyWaypointDiagram}
          src={require('./img/why-waypoint-diagram.svg')}
          alt="Waypoint Diagram"
        />
      </HomepageSection>

      <BrandedCta
        heading="Ready to get started?"
        content="TODO: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
        links={[
          { text: 'Download', url: '/download', type: 'download' },
          { text: 'Explore documentation', url: '/docs' },
        ]}
      />
    </div>
  )
}
