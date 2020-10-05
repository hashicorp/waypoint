import styles from './AnimatedStepsList.module.css'
import StepsList from './steps-list'
import StepsIndicator from './steps-indicator'
import AnimatedTerminal from 'components/animated-terminal'
import { useState } from 'react'

export default function AnimatedStepsList({ terminalHeroState, steps }) {
  const [indicatorIndex, setIndicatorIndex] = useState(0)

  const activeTerminalState = terminalHeroState

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
          frameLength={activeTerminalState.frameLength}
          loop={activeTerminalState.loop}
          lines={activeTerminalState.lines}
        />
      </div>
    </div>
  )
}
