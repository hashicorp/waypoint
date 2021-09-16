import { useInView } from 'react-intersection-observer'
import InlineSvg from '@hashicorp/react-inline-svg'
import classNames from 'classnames'
import NumberedBlock from 'components/homepage/numbered-block'
import Features, { FeaturesProps } from 'components/homepage/features'
import s from './style.module.css'

interface MonitorAndManageProps {
  heading: string
  features: FeaturesProps
}

export default function MonitorAndManage({
  heading,
  features,
}: MonitorAndManageProps): JSX.Element {
  const { ref, inView } = useInView({
    threshold: 0.5,
    triggerOnce: true,
    delay: 200,
  })
  return (
    <div className={s.root}>
      <div className={s.content}>
        <div className={s.contentInner}>
          <NumberedBlock index="3" heading={heading}>
            <Features items={features} />
            <InlineSvg
              className={s.logos}
              src={require('./logos.svg?include')}
            />
          </NumberedBlock>
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
