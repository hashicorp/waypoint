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
        description={
          <>
            <p>
              Waypoint allows developers to define their application build,
              deploy, and release lifecycle as code, reducing the time to
              deliver deployments through a consistent and repeatable workflow.
            </p>
            <br />
            <p>
              As we prepare for the 0.1 release at HashiConf in October - we
              invite you to take it for a test drive by exploring the links
              below.
            </p>
          </>
        }
        links={[
          {
            text: 'Download',
            url: 'https://go.hashi.co/waypoint-beta-binaries',
            type: 'download',
          },
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
