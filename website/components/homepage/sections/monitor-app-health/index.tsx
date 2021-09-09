import InlineSvg from '@hashicorp/react-inline-svg'
import classNames from 'classnames'
import { useInView } from 'react-intersection-observer'
import MediaObject from 'components/homepage/media-object'
import Section, {
  SectionHeading,
  SectionDescription,
} from 'components/homepage/section'
import s from './style.module.css'

export default function SectionMonitorAppHealth() {
  const { ref, inView } = useInView({
    threshold: 0.85,
    triggerOnce: true,
  })
  return (
    <Section className={s.monitorAppHealth}>
      <div
        className={classNames(s.media, {
          [s.visible]: inView,
        })}
        ref={ref}
      >
        <InlineSvg src={require('./graphic.svg?include')} />
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
