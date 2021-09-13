import Hero from 'components/homepage/hero'
import Intro from 'components/homepage/intro'
import SectionHowItWorks from 'components/homepage/sections/how-it-works'
import SectionMonitorAppHealth from 'components/homepage/sections/monitor-app-health'
import SectionExtendPlugins from 'components/homepage/sections/extend-plugins'
// import SectionWorkflowThatScales from 'components/homepage/sections/workflow-that-scales'
import s from './style.module.css'

export default function HomePage() {
  return (
    <div className={s.homePage}>
      <Hero />
      <Intro />
      <SectionHowItWorks />
      <SectionMonitorAppHealth />
      <SectionExtendPlugins />
      {/* <SectionWorkflowThatScales /> */}
    </div>
  )
}
