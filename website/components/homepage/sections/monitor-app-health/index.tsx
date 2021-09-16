import classNames from 'classnames'
import { useInView } from 'react-intersection-observer'
import Section, {
  SectionHeading,
  SectionDescription,
} from 'components/homepage/section'
import GraphicSvg from './graphic'
import Features, { FeaturesProps } from 'components/homepage/features'
import s from './style.module.css'

interface SectionMonitorAppHealthProps {
  heading: string
  description: string
  features: FeaturesProps
}

export default function SectionMonitorAppHealth({
  heading,
  description,
  features,
}: SectionMonitorAppHealthProps): JSX.Element {
  const { ref, inView } = useInView({
    threshold: 0.5,
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
        <GraphicSvg />
      </div>
      <div className={s.content}>
        <SectionHeading>{heading}</SectionHeading>
        <SectionDescription>{description}</SectionDescription>
        <div className={s.contentMediaObject}>
          <Features items={features} />
        </div>
      </div>
    </Section>
  )
}
