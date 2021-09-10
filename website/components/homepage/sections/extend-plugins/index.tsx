import { useInView } from 'react-intersection-observer'
import classNames from 'classnames'
import InlineSvg from '@hashicorp/react-inline-svg'
import MediaObject from 'components/homepage/media-object'
import Section, {
  SectionHeading,
  SectionDescription,
} from 'components/homepage/section'
import s from './style.module.css'

export default function SectionExtendPlugins() {
  const { ref, inView } = useInView({
    threshold: 0.8,
    triggerOnce: true,
    delay: 0.5,
  })
  return (
    <Section className={s.extendPlugins}>
      <div className={s.inner}>
        <div className={s.content}>
          <SectionHeading>Extend Waypoint with plugins</SectionHeading>
          <SectionDescription>
            Extend workflows via built-in plugins and an extensible interface.
            Supports custom builders, deployment platforms, registries, release
            managers, and more
          </SectionDescription>
          <div className={s.contentBlocks}>
            <MediaObject
              stacked={true}
              icon={require('../icons/box.svg?include')}
              heading="Available Plugins"
              description="View a list of existing HashiCorp maintained plugins"
            />
            <MediaObject
              stacked={true}
              icon={require('../icons/code-union.svg?include')}
              heading="Creating Waypoint Plugins"
              description="Learn to extend Waypoint for your project’s needs"
            />
          </div>
        </div>
        <div
          className={classNames(s.media, {
            [s.visible]: inView,
          })}
          ref={ref}
        >
          <InlineSvg src={require('./plugins.svg?include')} />
        </div>
      </div>
    </Section>
  )
}
