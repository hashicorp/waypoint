import Section, { SectionHeading } from '../../section'
import React from 'react'
import s from './style.module.css'

export default function SectionHowItWorks({ children }): JSX.Element {
  return (
    <Section className={s.howItWorks} id="how-it-works">
      <div className={s.container}>
        <SectionHeading>How it works</SectionHeading>
      </div>
      <ol>
        {React.Children.map(children, (child, index) => {
          return <li key={index}>{child}</li>
        })}
      </ol>
    </Section>
  )
}
