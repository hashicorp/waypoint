import styles from './AnimatedStepsList.module.css'
import StepsList from './steps-list'
import StepsIndicator from './steps-indicator'
import AnimatedTerminal from 'components/animated-terminal'
import { useState } from 'react'

export default function AnimatedStepsList({ steps }) {
  const [indicatorIndex, setIndicatorIndex] = useState(0)
  return (
    <div className={styles.animatedStepsList}>
      <div className={styles.indicatorWrapper}>
        <StepsIndicator steps={steps} activeIndex={indicatorIndex} />
      </div>

      <StepsList
        className={styles.stepsList}
        steps={steps}
        onFocusedIndexChanged={(newStep) => {
          setIndicatorIndex(newStep)
        }}
      />

      <div className={styles.terminalWrapper}>
        <AnimatedTerminal
          frameLength={100}
          loop={true}
          lines={[
            {
              frames: 5,
              code: ['$ waypoint up', '$ waypoint up |'],
            },
          ]}
        />
      </div>
    </div>
  )
}
