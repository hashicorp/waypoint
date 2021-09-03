import Hero from 'components/homepage/hero'
import NumberedBlock from 'components/homepage/numbered-block'
import MediaObject from 'components/homepage/media-object'
import Section, {
  SectionHeading,
  SectionDescription,
} from 'components/homepage/section'
import s from './style.module.css'

export default function HomePage() {
  return (
    <>
      <Hero />
      <Section className={s.sectionHowItWorks} id="how-it-works">
        <SectionHeading>How it works</SectionHeading>
        <NumberedBlock index="1" heading="Configure your app for Waypoint">
          <MediaObject
            icon={require('./icons/edit-pencil.svg?include')}
            heading="Writing waypoint.hcl files"
            description="Your waypoint.hcl file defines how Waypoint builds, deploys, and releases a project."
          />
          <MediaObject
            icon={require('./icons/layout.svg?include')}
            heading="Sample Waypoint files"
            description="View sample waypoint.hcl files to see how straight-forward it is to configure your deployments"
          />
        </NumberedBlock>
        <NumberedBlock index="2" heading="Build and deploy">
          <MediaObject
            icon={require('./icons/file-plus.svg?include')}
            heading="One simple command"
            description="Perform the build, deploy, and release steps for the app all from one simple command. Or instrument your Waypoint deployments through Remote or Git operations"
          />
        </NumberedBlock>
        <NumberedBlock index="3" heading="Monitor and manage in one place">
          <MediaObject
            icon={require('./icons/sliders.svg?include')}
            heading="One place for all your deployments"
            description="No matter where your developers are deploying to, monitor the activity through Waypoint’s aggregated logs and activity UI."
          />
        </NumberedBlock>
      </Section>
      <Section>
        <SectionHeading>Monitor app health on any cloud</SectionHeading>
        <SectionDescription>
          One place to monitor the entire lifecycle of your applications, no
          matter where you deploy to. View Logs, Builds, Releasese and even run
          Exec commands from the Waypoint UI
        </SectionDescription>
        <MediaObject
          icon={require('./icons/eye.svg?include')}
          heading="A single pane of glass"
          description="View all deployments, regardless of target from one location"
        />
      </Section>
      <Section>
        <SectionHeading>Extend Waypoint with plugins</SectionHeading>
        <SectionDescription>
          Extend workflows via built-in plugins and an extensible interface.
          Supports custom builders, deployment platforms, registries, release
          managers, and more
        </SectionDescription>
        <MediaObject
          stacked={true}
          icon={require('./icons/box.svg?include')}
          heading="Available Plugins"
          description="View a list of existing HashiCorp maintained plugins"
        />
        <MediaObject
          stacked={true}
          icon={require('./icons/code-union.svg?include')}
          heading="Creating Waypoint Plugins"
          description="Learn to extend Waypoint for your project’s needs"
        />
      </Section>
      <Section>
        <SectionHeading>One workflow that scales</SectionHeading>
        <SectionDescription>
          By creating one common workflow to enable developers to deploy; teams
          of every size can take advantage of Waypoint. Use plugins to
          automatically detect your tools, or for established projects, layer in
          your existing configuration like Dockerfiles and YAML.
        </SectionDescription>
      </Section>
    </>
  )
}
