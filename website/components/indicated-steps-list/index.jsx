import styles from './IndicatedStepsList.module.css'
import StepsList from './steps-list'
import StepsIndicator from './steps-indicator'
import Terminal from 'components/terminal'
import { useState } from 'react'

export default function IndicatedStepsList({ steps }) {
  const [indicatorIndex, setIndicatorIndex] = useState(0)
  return (
    <div className={styles.indicatedStepsList}>
      <div className={styles.indicatorWrapper}>
        <StepsIndicator steps={steps} activeIndex={indicatorIndex} />
      </div>
      <StepsList
        steps={steps}
        onFocusedIndexChanged={(newStep) => {
          setIndicatorIndex(newStep)
        }}
      />

      <div className={styles.terminalWrapper}>
        <Terminal className={styles.testingTerminal} code="$ waypoint up |" />
      </div>
    </div>
  )
}
