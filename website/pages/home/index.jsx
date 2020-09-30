import styles from './HomePage.module.css'
import InfoGrid from 'components/info-grid'
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

      <HomepageSection
        title="TODO, Interactive Code Here"
        theme="light"
      ></HomepageSection>

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
