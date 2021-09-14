import { useInView } from 'react-intersection-observer'
import classNames from 'classnames'
import MediaObject from 'components/homepage/media-object'
import Section, {
  SectionHeading,
  SectionDescription,
} from 'components/homepage/section'
import PluginsSvg from './plugins'
import s from './style.module.css'

export default function SectionExtendPlugins() {
  const { ref, inView } = useInView({
    threshold: 0.5,
    triggerOnce: true,
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
              icon={require('components/homepage/icons/box.svg?include')}
              heading="Available plugins"
              description="View a list of existing HashiCorp maintained plugins"
              link={{
                url: '/',
                text: 'Plugins',
              }}
            />
            <MediaObject
              stacked={true}
              icon={require('components/homepage/icons/code-union.svg?include')}
              heading="Creating Waypoint plugins"
              description="Learn to extend Waypoint for your project’s needs"
              link={{
                url: '/',
                text: 'Create',
              }}
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
