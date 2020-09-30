import HomepageHero from 'components/homepage-hero'
import BrandedCta from 'components/branded-cta'

export default function HomePage() {
  return (
    <div className="p-home">
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
