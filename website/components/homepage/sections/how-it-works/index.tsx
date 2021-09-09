import Section, { SectionHeading } from '../../section'
import ConfigureYourApp from './configure-your-app'
import BuildAndDeploy from './build-and-deploy'
import MonitorAndManage from './monitor-and-manage'
import s from './style.module.css'

export default function SectionHowItWorks() {
  return (
    <Section className={s.howItWorks} id="how-it-works">
      <div className={s.container}>
        <SectionHeading>How it works</SectionHeading>
      </div>
      <ol>
        <li>
          <ConfigureYourApp />
        </li>
        <li>
          <BuildAndDeploy />
        </li>
        <li>
          <MonitorAndManage />
        </li>
      </ol>
    </Section>
  )
}
