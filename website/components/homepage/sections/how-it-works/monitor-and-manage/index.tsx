import { useInView } from 'react-intersection-observer'
import InlineSvg from '@hashicorp/react-inline-svg'
import classNames from 'classnames'
import NumberedBlock from 'components/homepage/numbered-block'
import MediaObject from 'components/homepage/media-object'
import s from './style.module.css'

export default function MonitorAndManage() {
  const { ref, inView } = useInView({
    threshold: 1,
    triggerOnce: true,
    delay: 1000,
  })
  return (
    <div className={s.root}>
      <div className={s.content}>
        <div className={s.contentInner}>
          <NumberedBlock index="3" heading="Monitor and manage in one place">
            <MediaObject
              icon={require('../../icons/sliders.svg?include')}
              heading="One place for all your deployments"
              description="No matter where your developers are deploying to, monitor the activity through Waypoint’s aggregated logs and activity UI."
            />
          </NumberedBlock>
          <InlineSvg className={s.logos} src={require('./logos.svg?include')} />
        </div>
      </div>
      <div
        ref={ref}
        className={classNames(s.media, {
          [s.visible]: inView,
        })}
      >
        <InlineSvg src={require('./graphic.svg?include')} />
      </div>
    </div>
  )
}
