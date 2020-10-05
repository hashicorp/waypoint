import styles from './AnimatedStepsList.module.css'
import StepsList from './steps-list'
import StepsIndicator from './steps-indicator'
import useScrollPosition from 'lib/hooks/useScrollPosition'
import SteppedAnimatedTerminal from './stepped-animated-terminal'

import { useState } from 'react'

export default function AnimatedStepsList({ terminalHeroState, steps }) {
  const scrollPosition = useScrollPosition()
  const [indicatorIndex, setIndicatorIndex] = useState(0)
  const activeTerminalStateIndex =
    scrollPosition <= 350 ? 0 : indicatorIndex + 1

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
        <SteppedAnimatedTerminal
          activeIndex={activeTerminalStateIndex}
          steps={[terminalHeroState].concat(steps.map((step) => step.terminal))}
        />
      </div>
    </div>
  )
}
