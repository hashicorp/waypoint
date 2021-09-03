import NumberedBlock from '../../numbered-block'
import MediaObject from '../../media-object'
import Section, { SectionHeading } from '../../section'
import s from './style.module.css'

export default function SectionHowItWorks() {
  return (
    <Section className={s.howItWorks} id="how-it-works">
      <div className={s.inner}>
        <SectionHeading>How it works</SectionHeading>
        <ol>
          <li>
            <NumberedBlock index="1" heading="Configure your app for Waypoint">
              <MediaObject
                icon={require('../icons/edit-pencil.svg?include')}
                heading="Writing waypoint.hcl files"
                description="Your waypoint.hcl file defines how Waypoint builds, deploys, and releases a project."
              />
              <MediaObject
                icon={require('../icons/layout.svg?include')}
                heading="Sample Waypoint files"
                description="View sample waypoint.hcl files to see how straight-forward it is to configure your deployments"
              />
            </NumberedBlock>
            <img
              src={require('./img/configure-your-app-for-waypoint.png')}
              alt=""
            />
          </li>
          <li>
            <NumberedBlock index="2" heading="Build and deploy">
              <MediaObject
                icon={require('../icons/file-plus.svg?include')}
                heading="One simple command"
                description="Perform the build, deploy, and release steps for the app all from one simple command. Or instrument your Waypoint deployments through Remote or Git operations"
              />
            </NumberedBlock>
            <img src={require('./img/build-and-deploy.png')} alt="" />
          </li>
          <li>
            <NumberedBlock index="3" heading="Monitor and manage in one place">
              <MediaObject
                icon={require('../icons/sliders.svg?include')}
                heading="One place for all your deployments"
                description="No matter where your developers are deploying to, monitor the activity through Waypoint’s aggregated logs and activity UI."
              />
            </NumberedBlock>
            <img src={require('./img/monitor-and-manage.png')} alt="" />
          </li>
        </ol>
      </div>
    </Section>
  )
}
