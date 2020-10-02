import HomepageHero from 'components/homepage-hero'

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
            url:
              'http://ihngtake2gyn8nbyfgtgvu449dnsbrgopvukjdbntyndmlv7tb.s3-website-us-east-1.amazonaws.com/waypoint/',
            type: 'download',
          },
          {
            text: 'Get Started',
            url: '/docs/getting-started',
            type: 'inbound',
          },
        ]}
      />
    </div>
  )
}
