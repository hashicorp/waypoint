import { useInView } from 'react-intersection-observer'
import classNames from 'classnames'
import Section, {
  SectionHeading,
  SectionDescription,
} from 'components/homepage/section'
import PluginsSvg from './plugins'
import Features, { FeaturesProps } from 'components/homepage/features'
import s from './style.module.css'

interface SectionExtendPluginsProps {
  heading: string
  description: string
  features: FeaturesProps
}

export default function SectionExtendPlugins({
  heading,
  description,
  features,
}: SectionExtendPluginsProps): JSX.Element {
  const { ref, inView } = useInView({
    threshold: 0.5,
    triggerOnce: true,
  })
  return (
    <Section className={s.extendPlugins}>
      <div className={s.inner}>
        <div className={s.content}>
          <SectionHeading>{heading}</SectionHeading>
          <SectionDescription>{description}</SectionDescription>
          <div className={s.contentBlocks}>
            <Features
              items={features.map((feature) => {
                return {
                  stacked: true,
                  ...feature,
                }
              })}
            />
          </div>
        </div>
        <div
          className={classNames(s.media, {
            [s.visible]: inView,
          })}
          ref={ref}
        >
          <PluginsSvg />
        </div>
      </div>
    </Section>
  )
}
