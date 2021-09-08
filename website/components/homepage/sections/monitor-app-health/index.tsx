import MediaObject from 'components/homepage/media-object'
import Section, { SectionHeading, SectionDescription } from '../../section'
import s from './style.module.css'

export default function SectionMonitorAppHealth() {
  return (
    <Section className={s.monitorAppHealth}>
      <div className={s.media}>
        <img src={require('./img/monitor-app-health.png')} alt="" />
      </div>
      <div className={s.content}>
        <SectionHeading>Monitor app health on any cloud</SectionHeading>
        <SectionDescription>
          One place to monitor the entire lifecycle of your applications, no
          matter where you deploy to. View Logs, Builds, Releasese and even run
          Exec commands from the Waypoint UI
        </SectionDescription>
        <div className={s.contentMediaObject}>
          <MediaObject
            icon={require('../icons/eye.svg?include')}
            heading="A single pane of glass"
            description="View all deployments, regardless of target from one location"
          />
        </div>
      </div>
    </Section>
  )
}
