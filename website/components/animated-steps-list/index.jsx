import styles from './AnimatedStepsList.module.css'
import StepsList from './steps-list'
import StepsIndicator from './steps-indicator'
import useScrollPosition from 'lib/hooks/useScrollPosition'
import AnimatedTerminal from 'components/animated-terminal'
import FramedTerminal from 'components/animated-terminal/framed-terminal'
import { useState } from 'react'

// The breakpoints where the next step of each animation triggers
const breakpoints = [0, 350, 1258, 2309, 2880]

// The number of pixels before the next breakpoint that the animation should complete
const animationBottomPadding = [0, 400, 250, 0]

function calculateCurrentFrame(terminalSteps, currentIndex, scrollPosition) {
  const percentage = Math.min(
    (scrollPosition - breakpoints[currentIndex]) /
      (breakpoints[currentIndex + 1] -
        breakpoints[currentIndex] -
        animationBottomPadding[currentIndex]),
    1
  )
  const currentLines = terminalSteps[currentIndex].lines
  let totalFrames = 0
  currentLines.forEach((line) => {
    let frames = line.frames ? line.frames : 1
    if (Array.isArray(line.code)) {
      totalFrames += line.code.length * frames
    } else {
      totalFrames += frames
    }
  })
  return Math.max(0, percentage * (totalFrames - 1))
}

export default function AnimatedStepsList({ terminalHeroState, steps }) {
  const scrollPosition = useScrollPosition()
  const [indicatorIndex, setIndicatorIndex] = useState(0)
  const activeTerminalStateIndex =
    scrollPosition <= 350 ? 0 : indicatorIndex + 1
  const terminalSteps = [terminalHeroState].concat(
    steps.map((step) => step.terminal)
  )
  const currentFrame = calculateCurrentFrame(
    terminalSteps,
    activeTerminalStateIndex,
    scrollPosition
  )

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
        {activeTerminalStateIndex === 0 ? (
          <AnimatedTerminal
            frameLength={terminalSteps[activeTerminalStateIndex].frameLength}
            loop={terminalSteps[activeTerminalStateIndex].loop}
            lines={terminalSteps[activeTerminalStateIndex].lines}
          />
        ) : (
          <FramedTerminal
            frame={currentFrame}
            lines={terminalSteps[activeTerminalStateIndex].lines}
          />
        )}
      </div>
    </div>
  )
}
