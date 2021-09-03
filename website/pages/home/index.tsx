import Hero from 'components/homepage/hero'
import SectionHowItWorks from 'components/homepage/sections/how-it-works'
import SectionMonitorAppHealth from 'components/homepage/sections/monitor-app-health'
import SectionExtendPlugins from 'components/homepage/sections/extend-plugins'
import SectionWorkflowThatScales from 'components/homepage/sections/workflow-that-scales'

export default function HomePage() {
  return (
    <>
      <Hero />
      <SectionHowItWorks />
      <SectionMonitorAppHealth />
      <SectionExtendPlugins />
      <SectionWorkflowThatScales />
    </>
  )
}
